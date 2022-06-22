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
        """
        add_domain_to_crawl_queue adds the specified domain to the queue to crawl. It can be a combination of a
        subdomain and a domain or just a domain.

        :param domain: the domain to be added to the queue.
        :return: None
        """
        r = self.session.post(self._make_url(urls.ADD_DOMAIN_TO_CRAWL_QUEUE), data={"domain": domain})

        if r.status_code == 202 or r.status_code == 409:
            return None

        r.raise_for_status()

    def request_job(self, raise_on_bad_id: bool = False) -> Optional[Job]:
        """
        request_job requests a job to be run from the coordinator.

        :param raise_on_bad_id: if the UUID selected for the crawler is in use, raise an exception
        :return: None of there's no job in the queue, or a Job
        """
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
        """
        preflight_check checks that the provided URL is allowed to be loaded or not.

        :param url: URL to query against
        :return: True if load can go ahead, False if URL should be skipped.
        """
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
        """
        digest_loaded_page returns the result of a page load to the coordinator for storage.

        :param url: URL of the loaded page
        :param html: raw HTML of the loaded page
        :param loaded_at: time of page load
        :param extra_data: extra data to send to the coordinator. This will be piped directly into the request.
        :return: None
        """
        r = self.session.post(self._make_url(urls.DIGEST_PAGE_LOAD), data={
            "url": url,
            "html": html,
            "loaded_at": loaded_at.isoformat(),
            **extra_data,
        })
        r.raise_for_status()
