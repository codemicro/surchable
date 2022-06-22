import os
import sys
import structlog

logger = structlog.get_logger()

_hasErrored: bool = False


def _required_env_var(key: str) -> str:
    global _hasErrored

    x = os.getenv(key)
    if x is None:
        logger.error(f"missing required environment variable {key}")
        _hasErrored = True
    return x


def _env_var_default(key: str, default: str) -> str:
    x = os.getenv(key)
    if x is None:
        return default
    return x


"""[[[cog
import cog
from generateMappings import *
cog.outl(
	generate_python_configuration(
		parse_configuration(
			load_raw_configuration(),
		),
	),
)
]]]"""
# The below was generated. Do not edit.
# Modify mappings/urls instead.
DB_DATABASE_NAME = _required_env_var("SURCHABLE_DB_DATABASE_NAME")
DB_USER = _required_env_var("SURCHABLE_DB_USER")
DB_PASSWORD = _required_env_var("SURCHABLE_DB_PASSWORD")
DB_HOST = _required_env_var("SURCHABLE_DB_HOST")
COORDINATOR_SERVE_PORT = _env_var_default("SURCHABLE_COORDINATOR_SERVE_PORT", "7200")
COORDINATOR_SERVE_HOST = _env_var_default("SURCHABLE_COORDINATOR_SERVE_HOST", "0.0.0.0")

# [[[end]]]

if _hasErrored:
    sys.exit(1)
