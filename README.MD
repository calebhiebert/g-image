# GFile

A http file management utility

- Store files on a S3 bucket
- **Tested Services**
  - DigitalOcean Spaces
  - Google Cloud Storage
- Store files in a local cache on disk
- Customizable cache size
- Use webhooks for storage to any database or system
- API key generation with variable permissions (Create, Read, Delete, Admin)

## Quickstart

### Environment Variables

- `DATA_DIR` The file cache directory (defaults to ./data/)
- `BUCKET_NAME` S3 Bucket name
- `S3_ENDPOINT` S3 Endpoint
- `S3_ACCESS_KEY` S3 Access Key
- `S3_SECRET` S3 Secret Key
- `S3_SSL` If S3 should use SSL (default false)
- `WEBHOOK_URL` Webhook url
- `CACHE_SIZE` The size of the cache to store locally in mb (default 50)

When starting the program, a master api key will be printed to the console. It will look like this:

```
Admin API Key: h~9PjzdKQh3B7Y~S_XtGZrsle8ztfcry4wEH0FcViI5WjIPy
```

You can use this api key to make any request.

### Requests

#### Get File

`GET /:id`

Returns a file that matches `id`

#### Get File Info

`GET /:id/info`

Returns info about a file

#### Upload File

`POST /upload`

Accepts only multipart/form-data. The file should be at parameter `file`

Returns file data:

```json
{
  "id": "bchaka2md6sdpllvkvl0",
  "filename": "before-submit.png",
  "mime": "image/png",
  "size": 78529,
  "hash": "493ac2f6d6f009415be5653d8125ef9969379502ddd595e7dde59c22a79b432c"
}
```

#### Delete File

`DELETE /:id`

Deletes the file at `id`

Returns a blank 200 upon success

## Docker

### Plain ol' Docker

You can run a regular docker container with the following command

```shell
docker run -d \
  -p 8080:8080 --name gfile \
  -e BUCKET_NAME=example \
  -e S3_ENDPOINT=s3.endpoint.here \
  -e S3_ACCESS_KEY=<ACCESSKEY> \
  -e S3_SECRET=<SECRET> \
  -e WEBHOOK_URL=http://localhost:8008/hookr \
  -v "$PWD/data":/data \
  panchem/gfile
```

### Docker Compose

Include gfile in your docker compose file like this:

```yaml
version: '3'
services:
  gfile:
    image: panchem/gfile
    ports:
      - "8080:8080"
    environment:
      BUCKET_NAME: example
      S3_ENDPOINT: s3.endpoint.here
      S3_ACCESS_KEY: <ACCESSKEY>
      S3_SECRET: <SECRET>
      WEBHOOK_URL: http://localhost:8008/hookr
    volumes:
      - "$PWD/data:/data"
```

## Webhooks

GFile uses webhooks to allow users to store file data in a custom way on their own servers

#### Upload File

Webhook will send a application/json POST request with the following format:

```json
{
  "filename": "My awesome File.png",
  "hash": "10558753b0a010c4be0dd9698dddc90987242f0e218b5111d8194a2e6f6e5266",
  "id": "bch9ivamd6serk2kp5tg",
  "mime": "image/png",
  "size": 2967905
}
```

- `filename` The name of the submitted file
- `hash` The sha256 hash of the file
- `id` The file id (mongo compatible)
- `mime` File type
- `size` File size in bytes

The webhook expects an `application/json` response of exactly:

```json
{
  "ok": true
}
```

#### Get File

Webhook will send a GET request. This get request will contain a query parameter `id`, which is the id of the info to retrieve

eg: `GET http://examplehook.com/?id=bcgor7qmd6sb45tct3pg`

The webhook expects an `application/json` response in the format of:

```json
{
  "filename": "My awesome File.png",
  "hash": "10558753b0a010c4be0dd9698dddc90987242f0e218b5111d8194a2e6f6e5266",
  "id": "bch9ivamd6serk2kp5tg",
  "mime": "image/png",
  "size": 2967905
}
```

#### Delete File 

Webhook will send a DELETE request. This request will contain a query parameter `id`, which is the id of the file to delete

eg: `DELETE http://examplehook.com/?id=bcgotrqmd6sb45tct3qg`

The webhook expects an `application/json` response of exactly:

```json
{
  "ok": true
}
```
