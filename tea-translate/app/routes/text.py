from __future__ import annotations

import time
from typing import Literal, Optional

from fastapi import APIRouter, Request
from pydantic import BaseModel, Field

from app.core.response import fail, ok
from app.models.asr import get_asr
from app.models.translator import get_translator, to_nllb_lang

router = APIRouter()


class TextTranslateRequest(BaseModel):
    text: str = Field(..., min_length=1, description="待翻译文本")
    source_lang: str = Field(..., description="源语言代码（zh/en/auto/NLLB code）")
    target_lang: str = Field(..., description="目标语言代码（zh/en/NLLB code）")


def _detect_lang_via_whisper(text: str) -> str:
    """用 Whisper 语言检测做兜底（auto 场景）。

    faster-whisper 的 detect_language 需要一段音频，这里采用迂回策略：
    若用户传 auto，则默认按 en/zh 比例经验判定。更稳妥的做法是调用翻译器
    自带的语言识别模块，这里保持简单，直接返回 ``auto`` 占位交由上层处理。
    """
    # 简化：按字符分布判断 zh vs en
    zh_chars = sum(1 for c in text if "\u4e00" <= c <= "\u9fff")
    if zh_chars > len(text) * 0.2:
        return "zh"
    return "en"


@router.post("/text")
async def translate_text(req: TextTranslateRequest, request: Request):
    t0 = time.perf_counter()

    try:
        src = req.source_lang.strip().lower()
        tgt = req.target_lang.strip().lower()

        # auto 处理
        if src == "auto":
            src = _detect_lang_via_whisper(req.text)

        # 基本校验：src != tgt
        if to_nllb_lang(src) == to_nllb_lang(tgt):
            latency_ms = int((time.perf_counter() - t0) * 1000)
            return ok(
                {
                    "translation": req.text,
                    "source_lang": src,
                    "target_lang": tgt,
                    "latency_ms": latency_ms,
                }
            )

        translator = get_translator()
        translated = translator.translate(req.text, source_lang=src, target_lang=tgt)

        latency_ms = int((time.perf_counter() - t0) * 1000)
        return ok(
            {
                "translation": translated,
                "source_lang": src,
                "target_lang": tgt,
                "latency_ms": latency_ms,
            }
        )
    except ValueError as e:
        return fail(message=str(e), code=400)
    except Exception as e:
        return fail(message=f"translate failed: {e}", code=500)
