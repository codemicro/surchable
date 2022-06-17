# Coordinator HTTP protocol

## Add domain to crawl queue

**`POST /job/add`**

Example body:

```json
{
    "domain": "tdpain.net", // required - does not include subdomains
}
```

Possible responses:

* **202 Accepted** - domain has been accepted into the queue
    ```json
    {
        "id": "44151d5b-088c-4ba9-99dc-0c925dbd9cb9"
    }
    ```
* **409 Conflict** - domain+subdomain combination has already been enqueued
    ```json
    {
        "message": "domain already queued",
        "status": "error"
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