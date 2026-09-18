from __future__ import annotations

import os
from typing import Optional

import ctranslate2
from transformers import AutoTokenizer

from app.core.config import settings


# 支持的语言代码映射（小写 -> NLLB 代码）
NLLB_LANG_MAP: dict[str, str] = {
    "zh": "zho_Hans",
    "zh-cn": "zho_Hans",
    "zho": "zho_Hans",
    "en": "eng_Latn",
    "eng": "eng_Latn",
    "en-us": "eng_Latn",
    "en-gb": "eng_Latn",
}

# NLLB -> 我们对外暴露的短代码（用于回传）
NLLB_TO_SHORT: dict[str, str] = {
    "zho_Hans": "zh",
    "eng_Latn": "en",
}


def to_nllb_lang(code: str) -> str:
    """将用户传入的语言代码归一化为 NLLB 代码。"""
    key = code.strip().lower()
    if key in NLLB_LANG_MAP:
        return NLLB_LANG_MAP[key]
    # 允许直接传 NLLB 代码
    if "_" in code and code not in NLLB_LANG_MAP:
        return code
    raise ValueError(f"Unsupported language code: {code}")


class NLLBTranslator:
    """NLLB-200 + CTranslate2 单例翻译器。

    懒加载：首次调用 ``load()`` 或 ``translate()`` 时才真正加载模型。
    """

    _instance: Optional["NLLBTranslator"] = None

    def __new__(cls) -> "NLLBTranslator":
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._instance._loaded = False
            cls._instance._translator = None
            cls._instance._tokenizer = None
        return cls._instance

    def load(self) -> None:
        """显式加载模型（适合在 lifespan 中预热）。"""
        if self._loaded:
            return

        cache_dir = os.environ.get("TRANSFORMERS_CACHE", settings.MODEL_CACHE_DIR)
        os.makedirs(cache_dir, exist_ok=True)

        print(f"[NLLB] Loading tokenizer from {settings.NLLB_MODEL} (cache={cache_dir}) ...")
        self._tokenizer = AutoTokenizer.from_pretrained(settings.NLLB_MODEL, cache_dir=cache_dir)

        print(f"[NLLB] Loading CTranslate2 translator from {settings.NLLB_MODEL} ...")
        # compute_type="int8" 适合 CPU
        self._translator = ctranslate2.Translator(
            model_path=settings.NLLB_MODEL,
            device="cpu",
            compute_type="int8",
        )

        self._loaded = True
        print("[NLLB] Model loaded.")

    def _ensure_loaded(self) -> None:
        if not self._loaded:
            self.load()

    @property
    def is_loaded(self) -> bool:
        return self._loaded

    def translate(self, text: str, source_lang: str, target_lang: str) -> str:
        """翻译文本。

        Args:
            text: 待翻译文本。
            source_lang: 源语言代码（zh / en / NLLB code）。
            target_lang: 目标语言代码（zh / en / NLLB code）。

        Returns:
            译文文本。
        """
        self._ensure_loaded()

        src = to_nllb_lang(source_lang)
        tgt = to_nllb_lang(target_lang)

        if not text or not text.strip():
            return ""

        # NLLB tokenizer 需要在句首加上源语言 BOS token
        src_lang_token = self._tokenizer.convert_tokens_to_ids(src)
        tgt_lang_token = self._tokenizer.convert_tokens_to_ids(tgt)

        # 某些 tokenizer 返回 None，需要兜底
        if src_lang_token is None or tgt_lang_token is None:
            # 退回使用 tokenizer 的 language_code 参数（新版 transformers 支持）
            self._tokenizer.src_lang = src
            encoded = self._tokenizer(text, return_tensors="np")
            source_tokens = self._tokenizer.convert_ids_to_tokens(encoded["input_ids"][0])
        else:
            encoded = self._tokenizer(text, return_tensors="np")
            source_tokens = [self._tokenizer.convert_ids_to_tokens(src_lang_token)]
            source_tokens += self._tokenizer.convert_ids_to_tokens(encoded["input_ids"][0].tolist())

        # CTranslate2 translate_batch
        results = self._translator.translate_batch(
            [source_tokens],
            target_prefix=[[self._tokenizer.convert_ids_to_tokens(tgt_lang_token)]],
        )
        target_tokens = results[0].hypotheses[0][1:]  # 跳过目标语言 BOS
        output_ids = [self._tokenizer.convert_tokens_to_ids(t) for t in target_tokens]
        translated = self._tokenizer.decode(output_ids, skip_special_tokens=True)
        return translated.strip()


def get_translator() -> NLLBTranslator:
    return NLLBTranslator()
