import os
from typing import *
import structlog


def event_to_message_log_processor(logger: Any, _: str, event_dict: dict) -> dict:
    # This method is for rough Zerolog compatibility

    if "event" in event_dict:
        event_dict["message"] = event_dict["event"]
        del event_dict["event"]

    return event_dict


structlog.configure(
    processors=[
        structlog.processors.add_log_level,
        structlog.processors.StackInfoRenderer(),
        structlog.dev.set_exc_info,
        structlog.processors.TimeStamper(),
        event_to_message_log_processor,
        structlog.processors.JSONRenderer(),
    ],
)

import config
import coordinator


def main():
    coord = coordinator.Coordinator(
        config.COORDINATOR_SERVE_HOST, config.COORDINATOR_SERVE_PORT
    )
    coord.ping()

    logger = structlog.get_logger()
    logger.info("ok!")


if __name__ == "__main__":
    main()
