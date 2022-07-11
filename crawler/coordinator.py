import datetime
import math
import time
import uuid
from dataclasses import dataclass
from typing import *

import requests
import structlog

import urls


@dataclass
class Job:
    id: str
    domain: str
    start: str


class Coordinator:
    crawler_id: str
    address: str
    session: requests.Session

    def __init__(self, host: str, port: str):
        self.crawler_id = str(uuid.uuid4)
        self.address = f"{host}:{port}"
        self.session = requests.Session()

        self.session.headers["X-Crawler-ID"] = self.crawler_id

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

    def add_domain_to_crawl_queue(self, domain: str):
        r = self.session.post(self._make_url(urls.ADD_DOMAIN_TO_CRAWL_QUEUE), data={"domain": domain})

        if r.status_code == 202 or r.status_code == 409:
            return None

        r.raise_for_status()

    def request_job(self, raise_on_bad_id: bool = False) -> Optional[Job]:
        r = self.session.post(self._make_url(urls.CRAWLER_REQUEST_JOB))

        if r.status_code == 409 and not raise_on_bad_id:  # 409 conflict - another crawler already has this ID, which is
            # unlikely but ok well
            self.crawler_id = str(uuid.uuid4())
            return self.request_job(raise_on_bad_id=True)

        r.raise_for_status()

        if r.status_code == 204:
            return None

        assert r.status_code == 201

        response_json = r.json()

        return Job(response_json["id"], response_json["domain"], response_json["start"])

    def preflight_check(self, url: str) -> bool:
        r = self.session.post(self._make_url(urls.REQUEST_PREFLIGHT_CHECK), data={"url": url})
        r.raise_for_status()

        response_data = r.json()
        permission = response_data["permission"]

        if permission == "SKIP":
            return False
        elif permission == "LOAD":
            return True
        else:
            raise ValueError(f"unknown permission {permission}")

    def digest_loaded_page(self, url: str, html: str, loaded_at: datetime.datetime,
                           extra_data: Optional[Dict[str, Any]] = None):
        r = self.session.post(self._make_url(urls.DIGEST_PAGE_LOAD), data={
            "url": url,
            "html": html,
            "loaded_at": loaded_at.isoformat(),
            **extra_data,
        })
        r.raise_for_status()

    def mark_job_completed(self):
        r = self.session.post(self._make_url(urls.COMPLETE_JOB))
        r.raise_for_status()

    def mark_job_cancelled(self):
        r = self.session.post(self._make_url(urls.CANCEL_JOB))
        r.raise_for_status()

    def add_domain_to_blocklist(self, domain: str, reason: str):
        r = self.session.post(self._make_url(urls.BLOCKLIST_ADD), data={
            "domain": domain,
            "reason": reason,
        })
        r.raise_for_status()