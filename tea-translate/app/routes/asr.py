from __future__ import annotations

import os
import time
import uuid

from fastapi import APIRouter, File, Form, UploadFile

from app.core.response import fail, ok
from app.models.asr import get_asr
from app.models.translator import get_translator, to_nllb_lang

router = APIRouter()


@router.post("")
async def asr_translate(
    file: UploadFile = File(..., description="音频文件（PCM/WAV/MP3）"),
    target_lang: str = Form("auto", description="翻译目标语言，auto 时自动对译"),
):
    """上传音频 -> ASR -> 自动翻译为另一种语言。"""
    t0 = time.perf_counter()

    # 1) 落盘
    ext = os.path.splitext(file.filename or "")[1] or ".wav"
    tmp_path = f"/tmp/tea_asr_{uuid.uuid4().hex}{ext}"
    try:
        with open(tmp_path, "wb") as f:
            f.write(await file.read())

        # 2) ASR
        asr = get_asr()
        text, lang = asr.transcribe(tmp_path)

        # 3) 翻译
        translated = ""
        translation_latency = 0
        if text.strip():
            t1 = time.perf_counter()
            # auto -> 另一种主流语言
            src = lang if lang else "en"
            if target_lang == "auto":
                target_lang = "zh" if src.startswith("zh") or src == "zho" else "en"

            translator = get_translator()
            try:
                translated = translator.translate(text, source_lang=src, target_lang=target_lang)
            except ValueError:
                translated = ""
            translation_latency = int((time.perf_counter() - t1) * 1000)

        total_ms = int((time.perf_counter() - t0) * 1000)

        return ok(
            {
                "text": text,
                "language": lang,
                "translation": translated,
                "translation_latency_ms": translation_latency,
                "total_latency_ms": total_ms,
            }
        )
    except Exception as e:
        return fail(message=f"asr failed: {e}", code=500)
    finally:
        if os.path.exists(tmp_path):
            try:
                os.remove(tmp_path)
            except OSError:
                pass
