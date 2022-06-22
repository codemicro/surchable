# Coordinator HTTP protocol

In addition to any codes listed below, you may encounter errors such as `400 Bad Request` if malformed JSON bodies are sent.

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

## Crawler request preflight check

**`POST /page/preflight`**

Check to see if the page was requested recently and be told if the page should be loaded again.

Request body:

```json
{
    "url": "https://www.example.com/cheesecake"
}
```

Possible responses:

* **200 OK** - not loaded recently, go ahead
    ```json
    {
        "permission": "LOAD"
    }
    ```
* **200 OK** - loaded recently, move on to the next, please
    ```json
    {
        "permission": "SKIP"
    }
    ```
  
## Digest loaded page

**`POST /page/digest`**

Request body:

```json
{
  "url": "https://www.popularmechanics.com/home/interior-projects/a30743006/how-to-clean-mold-mildew/", // required
  "title": "How to Battle Mould and Mildew in Your Home",
  "description": "These unwelcome fungi cling to any damp area. You must destroy them.",
  "content": "One of the most common and challenging problems facing homeowners today...",
  "html": "<!DOCTYPE html><html><head>...", // required
  "notLoadBefore": 5, // wait 5 minutes before loading again, default 60. probably derived from headers.
  "outboundLinks": [
    "https://www.popularmechanics.com/author/6836/joseph-truini/",
    "https://www.popularmechanics.com/space/a40230192/dyson-sphere-immortality/",
    "..."
  ],
  "loadedAt": "1985-04-12T23:20:50.52Z" // required
}
```

Possible responses:

* **204 No Content** - all's good.