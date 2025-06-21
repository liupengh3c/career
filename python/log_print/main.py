import os
import sys
import logging
import logging.config

# 日志目录配置
LOG_DIR = "./logs"
os.makedirs(LOG_DIR, exist_ok=True)

LOG_CONFIG = {
    "version": 1,
    "disable_existing_loggers": False,  # 不禁用已有日志器
    "formatters": {
        "standard": {
            # "format": "%(asctime)s [%(levelname)s] [%(name)s]: %(message)s",
            "format": "%(asctime)s [%(levelname)s] %(pathname)s line %(lineno)d in [%(funcName)s]: %(message)s",
            "datefmt": "%Y-%m-%d %H:%M:%S",
        },
    },
    "handlers": {
        "console": {
            "class": "logging.StreamHandler",
            "formatter": "standard",
            "level": "DEBUG",
            "stream": sys.stdout,
        },
        "file": {
            "class": "logging.handlers.TimedRotatingFileHandler",
            "formatter": "standard",
            "level": "DEBUG",
            "filename": os.path.join(LOG_DIR, "app.log"),
            "when": "midnight",
            "interval": 1,
            "backupCount": 10,
            "encoding": "utf-8",
        },
    },
    "loggers": {
        "": {  # root logger
            "handlers": ["console", "file"],
            "level": "DEBUG",
            "propagate": False,
        },
    },
}

def division():
    try:
        logger.info("division func start")
        logger.debug("调试信息")
        x = 1 / 0
    except Exception:
        logger.exception("出现异常：")  # 打印堆栈信息

if __name__ == "__main__":
    logging.config.dictConfig(LOG_CONFIG)
    logger = logging.getLogger(__name__)
    division()