# Query API Documentation

The Query API implements some of [Bitcoin Block Explorer Query API](https://blockexplorer.com/q) endpoints (also compatible with [Blockchain.info API](https://blockchain.info/api/blockchain_api)).

## Path

For this documentation, we will assume every request begins with the above path:

	https://btcplex.com/api/

## Format

All calls are returned in **plain text**.

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

## GET /getblockcount

Returns the current block height / number of blocks in the longest chain.
	
### Example request

	$ curl https://btcplex.com/api/getblockcount

### Response

	275120

## GET /getblockhash/:height

Returns the hash of a block for the given height.

### Example request

	$ curl https://btcplex.com/api/getblockhash/233943

### Response

	00000000000001a5e28bb1e9aad40f846d66fb7911c0db4dc32a039bb168a237


## GET /latesthash

Returns the hash of the latest block.

### Example request

	$ curl https://btcplex.com/api/latesthash

### Response

	00000000000001a5e28bb1e9aad40f846d66fb7911c0db4dc32a039bb168a237

## GET /getreceivedbyaddress/:address

Returns the total received for the given address.

### Example request

	$ curl https://btcplex.com/api/getreceivedbyaddress/1CjPR7Z5ZSyWk6WtXvSFgkptmpoi4UM9BC

### Response

	22508511200317

## GET /getsentbyaddress/:address

Returns the total sent for the given address.

### Example request

	$ curl https://btcplex.com/api/getsentbyaddress/1CjPR7Z5ZSyWk6WtXvSFgkptmpoi4UM9BC

### Response

	21892586119180

## GET /addressbalance/:address

Returns getreceivedbyaddress minus getsentbyaddress for the given address.

### Example request

	$ curl https://btcplex.com/api/addressbalance/1CjPR7Z5ZSyWk6WtXvSFgkptmpoi4UM9BC

### Response

	615925081137

## GET /checkaddress/:address

Check a Bitcoin address for validity.

### Example request

	$ curl https://btcplex.com/api/checkaddress/19gzwTuuZDec8JZEddQUZH9kwzqkBfFtDa

### Response

	true
