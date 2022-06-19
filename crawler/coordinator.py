from typing import *
import time
import math
import requests
import structlog


class Coordinator:
    address: str

    def __init__(self, host: str, port: str):
        self.address = f"{host}:{port}"

    def _make_url(self, path: str) -> str:
        return "http://" + self.address + "/" + path.lstrip("/")

    def ping(self):
        logger = structlog.get_logger()
        logger.info(f"pinging coordinator on {self.address}")

        r: Optional[requests.Response] = None

        for i in range(1, 5):
            try:
                r = requests.get(
                    self._make_url("/ok"),
                )
            except requests.exceptions.ConnectionError:
                wait_time = int(math.exp(i))
                logger.warn(f"could not connect to coordinator. waiting {wait_time} seconds then retrying")
                time.sleep(wait_time)
                continue
            break

        assert r is not None, "could not connect to coordinator"
        r.raise_for_status()

        j = r.json()
        assert j.get("status") == "ok"

        logger.info("established connection to coordinator")
