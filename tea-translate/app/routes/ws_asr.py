from __future__ import annotations

import asyncio
import json
import time
from typing import Optional

import numpy as np
from fastapi import APIRouter, WebSocket, WebSocketDisconnect

from app.core.response import fail
from app.models.asr import get_asr
from app.models.translator import get_translator

router = APIRouter()


# ============ 延时关键参数（改这些就能调性能/精度平衡） ============
# 首次 flush 窗口：用户开始说话后多久出第一个 partial 字幕
#   800ms = 最快出字幕，短窗口 Whisper 识别率略降
FIRST_FLUSH_WINDOW_SEC = 0.8
# 稳定 partial 窗口：之后每个 partial 字幕的音频窗口
PARTIAL_WINDOW_SEC = 1.2
# 定时器兜底：即使音频不够窗口，多久强制 flush 一次
FLUSH_INTERVAL_SEC = 1.5
# 静音检测阈值 (RMS)：低于此值视为静音
VAD_RMS_THRESHOLD = 0.02
# 静音帧判定（连续多少个 20ms 静音帧 → 触发 final flush）
#   5 帧 = 100ms 静音窗口
VAD_SILENCE_FRAMES = 5
# final flush 最短音频（秒）：太短就跳过，不值得跑 NLLB
MIN_FINAL_SEC = 0.5

SAMPLE_RATE = 16000
BYTES_PER_SAMPLE = 2  # int16
FRAME_MS = 20  # VAD 每帧 20ms
FRAME_SAMPLES = SAMPLE_RATE * FRAME_MS // 1000  # 320
FRAME_BYTES = FRAME_SAMPLES * BYTES_PER_SAMPLE  # 640


def _has_speech(pcm: bytes) -> bool:
    """RMS VAD: 用 20ms 帧检测 PCM 是否有语音活动。

    把累积音频按 20ms 切片，任一块 RMS > 阈值就认为有语音。
    这个函数 < 1ms，零依赖，完全替代 Whisper 的 vad_filter。
    """
    if len(pcm) < FRAME_BYTES:
        return False
    audio = np.frombuffer(pcm, dtype=np.int16).astype(np.float32) / 32768.0
    for i in range(0, len(audio) - FRAME_SAMPLES + 1, FRAME_SAMPLES):
        chunk = audio[i : i + FRAME_SAMPLES]
        rms = float(np.sqrt(np.mean(chunk**2)))
        if rms > VAD_RMS_THRESHOLD:
            return True
    return False


def _count_silence_tail(pcm: bytes) -> int:
    """数尾部连续多少帧是静音的（用于触发 final flush）"""
    if len(pcm) < FRAME_BYTES:
        return 0
    audio = np.frombuffer(pcm, dtype=np.int16).astype(np.float32) / 32768.0
    count = 0
    # 从尾部往前扫
    for i in range(len(audio) - FRAME_SAMPLES, -1, -FRAME_SAMPLES):
        chunk = audio[i : i + FRAME_SAMPLES]
        rms = float(np.sqrt(np.mean(chunk**2)))
        if rms > VAD_RMS_THRESHOLD:
            break
        count += 1
    return count


@router.websocket("/asr-stream")
async def asr_stream(websocket: WebSocket):
    await websocket.accept()

    # 客户端首帧可选发 JSON {"target_lang": "zh"}
    target_lang = "auto"
    audio_buffer = bytearray()
    try:
        first = await asyncio.wait_for(websocket.receive(), timeout=1.0)
        if first.get("type") == "websocket.receive" and "text" in first:
            cfg = json.loads(first["text"])
            target_lang = cfg.get("target_lang", "auto")
        elif "bytes" in first:
            audio_buffer = bytearray(first["bytes"])
    except (asyncio.TimeoutError, json.JSONDecodeError):
        pass

    asr = get_asr()
    translator = get_translator()
    asr._ensure_loaded()
    translator._ensure_loaded()

    # 首次 flush 用小窗口
    first_flush_size = int(SAMPLE_RATE * BYTES_PER_SAMPLE * FIRST_FLUSH_WINDOW_SEC)
    steady_flush_size = int(SAMPLE_RATE * BYTES_PER_SAMPLE * PARTIAL_WINDOW_SEC)
    min_partial_bytes = int(SAMPLE_RATE * BYTES_PER_SAMPLE * 0.4)  # 400ms 以下不 flush

    last_flush = time.perf_counter()
    first_partial_sent = False
    ever_had_speech = False

    try:
        while True:
            try:
                msg = await asyncio.wait_for(websocket.receive(), timeout=FLUSH_INTERVAL_SEC)
            except asyncio.TimeoutError:
                # 定时器到点：如果缓冲区有东西 + 检测到语音 → flush partial
                if len(audio_buffer) >= min_partial_bytes and _has_speech(bytes(audio_buffer)):
                    await _send_partial(websocket, bytes(audio_buffer), asr)
                    audio_buffer.clear()
                    first_partial_sent = True
                    last_flush = time.perf_counter()
                continue

            if msg.get("type") == "websocket.disconnect":
                break
            if "bytes" in msg:
                audio_buffer.extend(msg["bytes"])
            elif "text" in msg:
                try:
                    obj = json.loads(msg["text"])
                    if obj.get("cmd") == "flush":
                        break
                except json.JSONDecodeError:
                    pass

            # ===== VAD：有语音才累积计数 =====
            has_speech = _has_speech(bytes(audio_buffer))
            if has_speech:
                ever_had_speech = True

            # ===== 尾部连续静音帧 → final flush =====
            silence_tail = _count_silence_tail(bytes(audio_buffer))
            if ever_had_speech and silence_tail >= VAD_SILENCE_FRAMES and len(audio_buffer) >= min_partial_bytes:
                # 先把有效语音部分送 final（去掉尾部静音帧）
                trimmed_len = len(audio_buffer) - silence_tail * FRAME_BYTES
                if trimmed_len >= min_partial_bytes:
                    await _send_final(websocket, bytes(audio_buffer[:trimmed_len]), asr, translator, target_lang)
                    audio_buffer.clear()
                    first_partial_sent = True  # final 后下次用 steady
                    last_flush = time.perf_counter()
                ever_had_speech = False  # 重置，等下一段语音

            # ===== partial flush =====
            now = time.perf_counter()
            flush_size = steady_flush_size if first_partial_sent else first_flush_size
            if len(audio_buffer) >= flush_size and has_speech:
                await _send_partial(websocket, bytes(audio_buffer), asr)
                audio_buffer.clear()
                first_partial_sent = True
                last_flush = now

    except WebSocketDisconnect:
        pass
    finally:
        # 退出时 final flush
        min_final_bytes = int(SAMPLE_RATE * BYTES_PER_SAMPLE * MIN_FINAL_SEC)
        if len(audio_buffer) >= min_final_bytes and ever_had_speech:
            await _send_final(websocket, bytes(audio_buffer), asr, translator, target_lang)


async def _send_partial(ws: WebSocket, pcm: bytes, asr) -> None:
    """partial: 只 Whisper 识别，不翻译（翻译留给 final）"""
    text = await asyncio.to_thread(asr.transcribe_realtime, pcm, SAMPLE_RATE)
    if not text.strip():
        return
    payload = {"type": "partial", "text": text}
    try:
        await ws.send_text(json.dumps(payload, ensure_ascii=False))
    except Exception:
        pass


async def _send_final(
    ws: WebSocket, pcm: bytes, asr, translator, target_lang: str
) -> None:
    """final: Whisper 识别 + NLLB-200 翻译，返回完整句子"""
    text = await asyncio.to_thread(asr.transcribe_realtime, pcm, SAMPLE_RATE)
    translated = ""
    try:
        if text.strip():
            if target_lang == "auto":
                zh_chars = sum(1 for c in text if "\u4e00" <= c <= "\u9fff")
                target_lang = "en" if zh_chars > len(text) * 0.2 else "zh"
            src_lang = "zh" if target_lang == "en" else "en"
            translated = await asyncio.to_thread(
                translator.translate, text, src_lang, target_lang
            )
    except Exception:
        translated = ""

    payload = {"type": "final", "text": text, "translation": translated}
    try:
        await ws.send_text(json.dumps(payload, ensure_ascii=False))
    except Exception:
        pass
