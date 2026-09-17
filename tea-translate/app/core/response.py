from __future__ import annotations

import uuid
from typing import Any

from fastapi import Request
from fastapi.responses import JSONResponse, Response
from starlette.middleware.base import BaseHTTPMiddleware, RequestResponseEndpoint


def ok(data: Any = None, message: str = "ok", code: int = 0) -> dict[str, Any]:
    """统一成功响应体。"""
    return {"code": code, "message": message, "data": data}


def fail(message: str = "error", code: int = 1, data: Any = None) -> dict[str, Any]:
    """统一失败响应体。"""
    return {"code": code, "message": message, "data": data}


class TraceIdMiddleware(BaseHTTPMiddleware):
    """X-Trace-Id 生成 / 透传中间件。

    - 若请求头携带 ``X-Trace-Id`` 则透传；否则生成新的 UUID。
    - 写入 ``request.state.trace_id``，同时在响应头回写。
    """

    HEADER_NAME = "X-Trace-Id"

    async def dispatch(self, request: Request, call_next: RequestResponseEndpoint) -> Response:
        trace_id = request.headers.get(self.HEADER_NAME) or uuid.uuid4().hex
        request.state.trace_id = trace_id

        response: Response = await call_next(request)
        response.headers[self.HEADER_NAME] = trace_id
        return response
