# Coordinator HTTP protocol

## Add domain to crawl queue

**`POST /job/add`**

Example body:

```json
{
    "domain": "tdpain.net", // required
}
```

Possible responses:

* **202 Accepted** - domain has been accepted into the queue
    ```json
    {
        "id": "44151d5b-088c-4ba9-99dc-0c925dbd9cb9"
    }
    ```
* **400 Bad Request** - validation failed
    ```json
    {
    "detail": [
        {
            "field": "Domain",
            "param": "",
            "tag": "required"
        }
    ],
    "message": "failed validation",
    "status": "error"
    }
    ```
* **409 Conflict** - domain+subdomain combination has already been enqueued
    ```json
    {
        "message": "domain already queued",
        "status": "error"
    }
    ```

## Crawler request job

**`GET /job/request`**

For a crawler to request a job to run.

The `X-Crawler-ID` header must be set to a unique ID to that crawler. The best option is a UUID that resets per crawler session.

Possible responses:

* **201 Created** - new job created for this crawler.
    ```json
    {
        "id": "44151d5b-088c-4ba9-99dc-0c925dbd9cb9",
        "domain": "sub.example.com",
        "start": "/", // path to start working at
    }
    ```
* **204 No Content** - there's nothing in the queue for processing at this time
* **409 Conflict** - another crawler with the same ID is already assigned to a job.
    ```json
    {
        "message": "crawler ID in use",
        "status": "error"
    }
    ```