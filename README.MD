# GFile

A http file management utility

- Store files on a S3 bucket (amazon, google, digitalocean, ovh, etc...)
- Store files in a local cache on disk
- Customizable cache size
- Use webhooks for storage to any database or system
- API key generation with variable permissions (Create, Read, Delete, Admin)

## Quickstart

When starting the program, a master api key will be printed to the console. It will look like this:

```
Admin API Key: h~9PjzdKQh3B7Y~S_XtGZrsle8ztfcry4wEH0FcViI5WjIPy
```

You can use this api key to make any request.

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
