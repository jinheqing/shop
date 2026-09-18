from __future__ import annotations

import asyncio
import contextlib
import sys
import time

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.core.config import settings
from app.core.response import TraceIdMiddleware, ok
from app.models.asr import get_asr
from app.models.translator import get_translator
from app.routes import asr as asr_route
from app.routes import text as text_route
from app.routes import ws_asr as ws_asr_route


@contextlib.asynccontextmanager
async def lifespan(app: FastAPI):
    # ---- 启动：预加载模型 ----
    t0 = time.perf_counter()
    print("[lifespan] Preloading Whisper + NLLB models ...", flush=True)

    def _load_all():
        get_asr().load()
        get_translator().load()

    # 放线程里不阻塞事件循环
    await asyncio.to_thread(_load_all)

    dt = int((time.perf_counter() - t0) * 1000)
    print(f"[lifespan] Models ready in {dt} ms.", flush=True)

    yield

    # ---- 关闭 ----
    print("[lifespan] Shutting down.", flush=True)


app = FastAPI(
    title="Tea Translate Engine",
    version="0.1.0",
    description="NLLB-200 (CTranslate2) + Faster-Whisper 翻译/ASR 服务",
    lifespan=lifespan,
)

# 中间件
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
app.add_middleware(TraceIdMiddleware)

# 路由
app.include_router(text_route.router, prefix="/translate", tags=["text"])
app.include_router(asr_route.router, prefix="/translate", tags=["asr"])
app.include_router(ws_asr_route.router, prefix="/translate", tags=["ws"])


@app.get("/health")
async def health():
    return ok(
        {
            "status": "ok",
            "whisper_loaded": get_asr().is_loaded,
            "nllb_loaded": get_translator().is_loaded,
        }
    )


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "app.main:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=False,
    )
