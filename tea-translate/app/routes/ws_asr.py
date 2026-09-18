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


# ============ 延时关键参数 ============
# PARTIAL_WINDOW_SEC: partial 字幕的音频窗口（秒）
#   2.0s → 1.2s: 每个 partial 字幕提前 ~800ms 出，代价是识别准确率略降（1.2s 短窗口）
#   如果需要精准识别可改回 2.0s，但字幕会 ~800ms 延迟
PARTIAL_WINDOW_SEC = 1.2
# FLUSH_INTERVAL_SEC: 定时器 flush 间隔（音频不够窗口时兜底）
FLUSH_INTERVAL_SEC = 1.5
SAMPLE_RATE = 16000
BYTES_PER_SAMPLE = 2  # int16

PARTIAL_CHUNK_SIZE = int(SAMPLE_RATE * BYTES_PER_SAMPLE * PARTIAL_WINDOW_SEC)


@router.websocket("/asr-stream")
async def asr_stream(websocket: WebSocket):
    await websocket.accept()

    # 可选：客户端首帧发 JSON 配置 {"target_lang": "zh"}
    target_lang = "auto"
    try:
        first = await asyncio.wait_for(websocket.receive(), timeout=1.0)
        if first.get("type") == "websocket.receive" and "text" in first:
            cfg = json.loads(first["text"])
            target_lang = cfg.get("target_lang", "auto")
        elif "bytes" in first:
            # 首帧就是音频，塞回去让主循环处理
            audio_buffer = bytearray(first["bytes"])
        else:
            audio_buffer = bytearray()
    except (asyncio.TimeoutError, json.JSONDecodeError):
        audio_buffer = bytearray()

    asr = get_asr()
    translator = get_translator()

    # 确保模型已加载
    asr._ensure_loaded()
    translator._ensure_loaded()

    last_flush = time.perf_counter()
    running = True

    try:
        while running:
            try:
                msg = await asyncio.wait_for(websocket.receive(), timeout=FLUSH_INTERVAL_SEC)
            except asyncio.TimeoutError:
                # 定时 flush partial
                if len(audio_buffer) >= SAMPLE_RATE * BYTES_PER_SAMPLE:
                    await _send_partial(websocket, bytes(audio_buffer), asr)
                continue

            if msg.get("type") == "websocket.disconnect":
                break
            if "bytes" in msg:
                audio_buffer.extend(msg["bytes"])
            elif "text" in msg:
                # 控制指令：{"cmd": "flush"} -> final
                try:
                    obj = json.loads(msg["text"])
                    if obj.get("cmd") == "flush":
                        break
                except json.JSONDecodeError:
                    pass

            # partial 条件：够一个窗口 或 定时器到点
            now = time.perf_counter()
            if len(audio_buffer) >= PARTIAL_CHUNK_SIZE or (
                now - last_flush >= FLUSH_INTERVAL_SEC
                and len(audio_buffer) >= SAMPLE_RATE * BYTES_PER_SAMPLE
            ):
                await _send_partial(websocket, bytes(audio_buffer), asr)
                audio_buffer.clear()
                last_flush = now

    except WebSocketDisconnect:
        pass
    finally:
        # final flush
        if len(audio_buffer) >= SAMPLE_RATE * BYTES_PER_SAMPLE:
            await _send_final(websocket, bytes(audio_buffer), asr, translator, target_lang)


async def _send_partial(ws: WebSocket, pcm: bytes, asr) -> None:
    # asyncio.to_thread: 把 CPU 密集的 Whisper 推理放到线程池，不阻塞 event loop
    # event loop 被阻塞 → 收不到新音频包 → 积压 → 字幕延迟
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
    # Whisper: to_thread 避免阻塞
    text = await asyncio.to_thread(asr.transcribe_realtime, pcm, SAMPLE_RATE)
    translated = ""
    try:
        if text.strip():
            if target_lang == "auto":
                zh_chars = sum(1 for c in text if "\u4e00" <= c <= "\u9fff")
                target_lang = "en" if zh_chars > len(text) * 0.2 else "zh"
            src_lang = "zh" if target_lang == "en" else "en"
            # NLLB-200: to_thread 避免阻塞
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
