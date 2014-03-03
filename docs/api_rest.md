# REST API Documentation

The REST API..

## Path

For this documentation, we will assume every request begins with the above path:

	https://btcplex.com/api/

## Format

All calls are returned in **JSON**.

## HATEOAS Links

Each API calls includes a ``_links`` section containing related links, see the [HAL specification](http://stateless.co/hal_specification.html).
 
## Pagination

Traversing with pagination is easy, you just need to follow HATEOS links.

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

## GET /block/:hash

Return block details along with transactions.

### Example request


	$ curl https://btcplex.com/api/block/000000000000170b01901a691a88d0bc1cde49fe32675d920039540613e3f2d7


### Response

```json
{
  "_links": {
    "next_block": {
      "href": "https://btcplex.com/api/block/00000000000034fa21051368f5a197e65239efb2f99a831615bbdd499429ab94"
    }, 
    "previous_block": {
      "href": "https://btcplex.com/api/block/00000000000005f34c6e10ffeb86f3073c119531629f4e2728204431fd3a6ba7"
    }, 
    "self": {
      "href": "https://btcplex.com/api/block/000000000000170b01901a691a88d0bc1cde49fe32675d920039540613e3f2d7"
    }
  }, 
  "bits": 453023994, 
  "hash": "000000000000170b01901a691a88d0bc1cde49fe32675d920039540613e3f2d7", 
  "height": 121426, 
  "mrkl_root": "71b01258157daeddd7e4b08bf2a149eb0878581e1a108c4ccef801867d105b17", 
  "n_tx": 10, 
  "next_block": "00000000000034fa21051368f5a197e65239efb2f99a831615bbdd499429ab94", 
  "nonce": 702499203, 
  "prev_block": "00000000000005f34c6e10ffeb86f3073c119531629f4e2728204431fd3a6ba7", 
  "size": 3040, 
  "time": 1304344768, 
  "total_out": 23308967646, 
  "tx": [
    {
      "block_hash": "000000000000170b01901a691a88d0bc1cde49fe32675d920039540613e3f2d7", 
      "block_height": 121426, 
      "block_time": 1304344768, 
      "first_seen_height": 0, 
      "first_seen_time": 0, 
      "hash": "cd6f351ba4c9b1d17b367dcb72bfdab0a99dbec2c7a5122f5641ea51b01f08e1", 
      "in": [], 
      "lock_time": 0, 
      "out": [
        {
          "hash": "1BDJGdvEbyy5v53yFspWkLz8f6Un2EYkWz", 
          "n": 0, 
          "spent": {
            "spent": false
          }, 
          "value": 5000000000
        }
      ], 
      "size": 134, 
      "ver": 1, 
      "vin_sz": 0, 
      "vin_total": 0, 
      "vout_sz": 1, 
      "vout_total": 5000000000
    }, 
    {
      "block_hash": "000000000000170b01901a691a88d0bc1cde49fe32675d920039540613e3f2d7", 
      "block_height": 121426, 
      "block_time": 1304344768, 
      "first_seen_height": 0, 
      "first_seen_time": 0, 
      "hash": "407fea16950684221ec202e8b5d51f9883cf6705c6fe980b329a6a6e2bcae090", 
      "in": [
        {
          "n": 0, 
          "prev_out": {
            "address": "1EKc4Go9Hg7AJ85TSuw7HRNN6bHsJPteCq", 
            "hash": "f17fc0337c6a8ca197e6bedd1395c8905326988993e994052566d97b526eb1b3", 
            "n": 1, 
            "value": 25000000
          }
        }, 
        {
          "n": 1, 
          "prev_out": {
            "address": "1MR3oDXdsk2wuj2oP7dfMMFWcDHCcpRakR", 
            "hash": "a5977cafded3d9a976dd1d3eb2796933f884576484d6c5b071f6cec5b41c0a91", 
            "n": 1, 
            "value": 1145000000
          }
        }, 
        {
          "n": 2, 
          "prev_out": {
            "address": "1SgRn6cf1wkZ955stT8LsfTorFuxqpw3G", 
            "hash": "8cd02632d840eb6424df5bd1099a5b4ee3752b8e506c3505f23790541ba94272", 
            "n": 0, 
            "value": 1331000000
          }
        }
      ], 
      "lock_time": 0, 
      "out": [
        {
          "hash": "1NPcjR6vhWz1mJv212cQ1B4dqVrjMsmNZH", 
          "n": 1, 
          "spent": {
            "spent": false
          }, 
          "value": 1000000
        }, 
        {
          "hash": "1AnW35jwT5HLpEZ9FFHoXwzcv1Jo2Y4hBD", 
          "n": 0, 
          "spent": {
            "spent": false
          }, 
          "value": 2500000000
        }
      ], 
      "size": 617, 
      "ver": 1, 
      "vin_sz": 3, 
      "vin_total": 2501000000, 
      "vout_sz": 2, 
      "vout_total": 2501000000
    }, 
    {
      "block_hash": "000000000000170b01901a691a88d0bc1cde49fe32675d920039540613e3f2d7", 
      "block_height": 121426, 
      "block_time": 1304344768, 
      "first_seen_height": 0, 
      "first_seen_time": 0, 
      "hash": "f48f1aec11f556dfea095e4163483a0d053695abdcaef667227a6cec305191f0", 
      "in": [
        {
          "n": 1, 
          "prev_out": {
            "address": "1AUGK1zSUxHZQMhc74e9uhxhj99CjrJwMF", 
            "hash": "202ea2c71cbc65c7d1f49f8d6231167ead5e9fd37e36e5e376f384fac04c84e4", 
            "n": 1, 
            "value": 5000000
          }
        }, 
        {
          "n": 0, 
          "prev_out": {
            "address": "1LwEW32g6P4ZpfnrKwCFrPRtYJyuWAJMHk", 
            "hash": "819ef1d7c583086136ed58017b4d2b7e6de0dab719a6325b94ed3f57074030b2", 
            "n": 0, 
            "value": 95000000
          }
        }
      ], 
      "lock_time": 0, 
      "out": [
        {
          "hash": "13dF2f1joqBh3yHuDGXK8YH8xELYCn8ogG", 
          "n": 0, 
          "spent": {
            "spent": false
          }, 
          "value": 100000000
        }
      ], 
      "size": 404, 
      "ver": 1, 
      "vin_sz": 2, 
      "vin_total": 100000000, 
      "vout_sz": 1, 
      "vout_total": 100000000
    }
  ], 
  "ver": 1
}
```

## GET /tx/:hash

Returns tx data.

### Example request

	$ curl https://btcplex.com/api/tx/79e6e7fa17c8aaf41ca8e0bed7e8bcc27247c50ff6a00bb99c5c3a0da1803412

### Response

```json
{
  "_links": {
    "self": {
      "href": "https://btcplex.com/api/tx/79e6e7fa17c8aaf41ca8e0bed7e8bcc27247c50ff6a00bb99c5c3a0da1803412"
    }
  }, 
  "block_hash": "", 
  "block_height": 0, 
  "block_time": 0, 
  "first_seen_height": 288719, 
  "first_seen_time": 1393839450, 
  "hash": "79e6e7fa17c8aaf41ca8e0bed7e8bcc27247c50ff6a00bb99c5c3a0da1803412", 
  "in": [
    {
      "n": 0, 
      "prev_out": {
        "address": "1G4SFtH12Pq8PsdCQz4DNuRkKNXVn2Q9W8", 
        "hash": "16755d281c7756391c3ebd13f60acbd3bf589053beb34bda27109feb393419c5", 
        "n": 0, 
        "value": 3000000
      }
    }, 
    {
      "n": 0, 
      "prev_out": {
        "address": "17B8x9qberMayEqpAKXQbEWRCv18Ki5jEw", 
        "hash": "a75ae3207959f2e93575dc214b992816b840412d4c26e708b60ea5ba9d6a7062", 
        "n": 2, 
        "value": 137444
      }
    }
  ], 
  "lock_time": 0, 
  "out": [
    {
      "hash": "1PXWmSvrKdnZsZ5io4GDQQAQDUMn2uieKX", 
      "n": 0, 
      "spent": {
        "spent": false
      }, 
      "value": 3000000
    }, 
    {
      "hash": "1LSfUnzxeyymhMc5RM8zhStE5Vqijmxcps", 
      "n": 0, 
      "spent": {
        "spent": false
      }, 
      "value": 117443
    }
  ], 
  "size": 438, 
  "ver": 1, 
  "vin_sz": 2, 
  "vin_total": 3137444, 
  "vout_sz": 2, 
  "vout_total": 3117443
}
```

## GET /address/:address

Returns address summary with relevant transactions.

### Example request

	$ curl https://btcplex.com/api/address/19gzwTuuZDec8JZEddQUZH9kwzqkBfFtDa

### Response

```json
{
  "_links": {
    "self": {
      "href": "https://btcplex.com/api/address/19gzwTuuZDec8JZEddQUZH9kwzqkBfFtDa"
    }
  }, 
  "address": "19gzwTuuZDec8JZEddQUZH9kwzqkBfFtDa", 
  "final_balance": 0, 
  "n_tx": 0, 
  "total_received": 0, 
  "total_sent": 0, 
  "txs": []
}
```
