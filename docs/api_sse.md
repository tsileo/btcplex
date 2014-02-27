# Server-Sent Events API Documentation

The [Server-Sent Events](https://developer.mozilla.org/en-US/docs/Server-sent_events/Using_server-sent_events) API allows you to get real time update on events.

You can consume server-sent events using [Javascript](https://developer.mozilla.org/en-US/docs/Server-sent_events) or any language (with Python using [sseclient](https://pypi.python.org/pypi/sseclient)...).

## Path

For this documentation, we will assume every request begins with the above path:

  https://btcplex.com/api/v1/

## Format

All calls are returned in **JSON**.

## Rate limiting

The rate limit allows you to make **3600 requests per hour** and implements the standard ``X-RateLimit-*`` headers in every API response

- ``X-RateLimit-Limit`` The number of requests allowed per hour.
- ``X-RateLimit-Remaining`` The number of requests remaining in the current window.
- ``X-RateLimit-Reset`` The time (in UTC epoch seconds) at which the rate limit window resets.

## Status Codes

- **200 OK** Response to a successful request.
- **429 Too many requests** Request aborted due to rate-limiting.
- **500 Internal server error** Something bad happened.

## Resources

All endpoints are listed here:

## GET /getblockcount

Returns the current block height / number of blocks in the longest chain.
	
### Example request

  $ curl https://btcplex.com/api/v1/getblockcount

### Response

  275120
