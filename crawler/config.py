import os
import sys
import structlog

logger = structlog.get_logger()

_hasErrored: bool = False


def _requiredEnvVar(key: str) -> str:
    global _hasErrored

    x = os.getenv(key)
    if x is None:
        logger.fatal(f"missing required environment variable {key}")
        _hasErrored = True
    return x


COORDINATOR_SERVE_PORT = _requiredEnvVar("SURCHABLE_COORDINATOR_SERVE_PORT")
COORDINATOR_SERVE_HOST = _requiredEnvVar("SURCHABLE_COORDINATOR_SERVE_HOST")

if _hasErrored:
    sys.exit(1)
