from __future__ import annotations

import os
import tempfile
from typing import Optional

import numpy as np
from faster_whisper import WhisperModel

from app.core.config import settings


class WhisperASR:
    """Faster-Whisper base 单例 ASR。

    懒加载；同时暴露标准 transcribe 与简易实时片段识别。
    """

    _instance: Optional["WhisperASR"] = None

    def __new__(cls) -> "WhisperASR":
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._instance._loaded = False
            cls._instance._model = None
        return cls._instance

    def load(self) -> None:
        if self._loaded:
            return

        cache_dir = os.environ.get("TRANSFORMERS_CACHE", settings.MODEL_CACHE_DIR)
        os.makedirs(cache_dir, exist_ok=True)

        print(
            f"[Whisper] Loading model='{settings.WHISPER_MODEL}' "
            f"device='{settings.WHISPER_DEVICE}' compute='{settings.WHISPER_COMPUTE_TYPE}' ..."
        )
        self._model = WhisperModel(
            settings.WHISPER_MODEL,
            device=settings.WHISPER_DEVICE,
            compute_type=settings.WHISPER_COMPUTE_TYPE,
            download_root=cache_dir,
        )
        self._loaded = True
        print("[Whisper] Model loaded.")

    def _ensure_loaded(self) -> None:
        if not self._loaded:
            self.load()

    @property
    def is_loaded(self) -> bool:
        return self._loaded

    def transcribe(self, audio_path: str, language: Optional[str] = None) -> tuple[str, str]:
        """转写完整音频文件。

        Args:
            audio_path: 音频文件路径（WAV/MP3/PCM 等）。
            language: 强制指定语言（可选），不传则自动检测。

        Returns:
            ``(text, language)`` — 转写文本与检测到的短语言代码（zh/en/...）。
        """
        self._ensure_loaded()

        segments, info = self._model.transcribe(
            audio_path,
            language=language,
            beam_size=5,
            vad_filter=True,
        )
        text = "".join(seg.text for seg in segments).strip()
        lang = (info.language or "").lower()
        return text, lang

    def transcribe_realtime(self, audio_bytes: bytes, sample_rate: int = 16000) -> str:
        """实时片段识别。

        将传入的 PCM 原始字节（int16, mono）写成临时文件后调用 ``transcribe``。
        这是一个简化实现，生产环境可替换为 streaming VAD + 增量解码。
        """
        self._ensure_loaded()

        # 16-bit PCM -> float32 numpy
        if len(audio_bytes) == 0:
            return ""

        audio = np.frombuffer(audio_bytes, dtype=np.int16).astype(np.float32) / 32768.0

        # faster_whisper 也接受 numpy 数组 + 采样率
        segments, _info = self._model.transcribe(
            audio,
            language=None,
            beam_size=1,
            vad_filter=True,
        )
        return "".join(seg.text for seg in segments).strip()


def get_asr() -> WhisperASR:
    return WhisperASR()
