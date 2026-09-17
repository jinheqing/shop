from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_file=".env", env_file_encoding="utf-8", extra="ignore")

    HOST: str = "0.0.0.0"
    PORT: int = 8090

    WHISPER_MODEL: str = "base"
    WHISPER_DEVICE: str = "cpu"
    WHISPER_COMPUTE_TYPE: str = "int8"

    NLLB_MODEL: str = "facebook/nllb-200-distilled-1.3B"

    MODEL_CACHE_DIR: str = "./models"

    CORS_ORIGINS: list[str] = ["*"]


settings = Settings()
