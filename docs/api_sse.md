# Server-Sent Events API Documentation

The [Server-Sent Events](https://developer.mozilla.org/en-US/docs/Server-sent_events/Using_server-sent_events) API allows you to get real time update on events.

You can consume server-sent events using [Javascript](https://developer.mozilla.org/en-US/docs/Server-sent_events) or any language (with Python using [sseclient](https://pypi.python.org/pypi/sseclient)...).

## Path

For this documentation, we will assume every request begins with the above path:

	https://btcplex.com/api/

## Format

All calls are returned in **JSON**.

## Rate limiting

The rate limit allows you to make **3600 requests per hour** and implements the standard ``X-RateLimit-*`` headers in every API response

- ``X-RateLimit-Limit`` The number of requests allowed per hour.
- ``X-RateLimit-Remaining`` The number of requests remaining in the current window.
- ``X-RateLimit-Reset`` The time (in UTC epoch seconds) at which the rate limit window resets.

## Cross Origin Resource Sharing

The API supports Cross Origin Resource Sharing (CORS) allowing you to make AJAX requests from anywhere.

## Status Codes

- **200 OK** Response to a successful request.
- **429 Too many requests** Request aborted due to rate-limiting.
- **500 Internal server error** Something bad happened.

## Resources

All endpoints are listed here:

## GET /blocknotify

Get the new best block hash when it changes.

### Example

```javascript
var blocknotify = new EventSource("https://btcplex.com/api/blocknotify");
blocknotify.onmessage = function(e) {
	console.log("New best block hash: " + e.data);
}
```

## GET /utxs

Get the unconfirmed transactions stream.

### Example

```javascript
var blocknotify = new EventSource("https://btcplex.com/api/utxs");
blocknotify.onmessage = function(e) {
	var data = JSON.parse(e.data);
	console.log("New unconfirmed tx: " + data);
}
```

## GET /utxs/:address

Get the unconfirmed transactions stream involving the given address.

### Example

```javascript
var address = "1dice6gJgPDYz8PLQyJb8cgPBnmWqCSuF";
var blocknotify = new EventSource("https://btcplex.com/api/utxs/" + address);
blocknotify.onmessage = function(e) {
	var data = JSON.parse(e.data);
	console.log("New unconfirmed tx involving " + address + ": " + data);
}
```
