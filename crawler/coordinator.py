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

        r = requests.get(
            self._make_url("/ok"),
        )
        r.raise_for_status()

        j = r.json()
        assert j.get("status") == "ok"
