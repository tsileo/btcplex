// THIS FILE IS AUTO-GENERATED FROM api_docs.tmpl
// DO NOT EDIT.

package templates

import "time"

import "camlistore.org/pkg/fileembed"

func init() {
	Files.Add("api_docs.tmpl", 6880, time.Unix(0, 1393190021254994903), fileembed.String("<h2>API Documentation</h2>\n"+
		"<div class=\"row\" style=\"margin-top:30px;\">\n"+
		"<div class=\"col-md-3\">\n"+
		"	<div class=\"list-group\" id=\"apiendpoints\">\n"+
		"  <a href=\"#getblockcount\" class=\"list-group-item\">GET /getblockcount</a>\n"+
		"  <a href=\"#getblockhash_height\" class=\"list-group-item\">GET /getblockhash/:heigh"+
		"t</a>\n"+
		"  <a href=\"#latesthash\" class=\"list-group-item\">GET /latesthash</a>\n"+
		"  <a href=\"#block_hash\" class=\"list-group-item\">GET /block/:hash</a>\n"+
		"  <a href=\"#tx_hash\" class=\"list-group-item\">GET /tx/:hash</a>\n"+
		"  <a href=\"#address_hash\" class=\"list-group-item\">GET /address/:hash</a>\n"+
		"  <a href=\"#getreceivedbyaddress_address\" class=\"list-group-item\">GET /getreceive"+
		"dbyaddress/:address</a>\n"+
		"  <a href=\"#block-height_height\" class=\"list-group-item\">GET /getsentbyaddress/:a"+
		"ddress</a>\n"+
		"  <a href=\"#block-height_height\" class=\"list-group-item\">GET /addressbalance/:add"+
		"ress</a>\n"+
		"  <a href=\"#block-height_height\" class=\"list-group-item\">GET /addressfirstseen/:a"+
		"ddress</a>\n"+
		"  <a href=\"#checkaddress_address\" class=\"list-group-item\">GET /checkaddress/:addr"+
		"ess</a>\n"+
		"</div>\n"+
		"\n"+
		"</div>\n"+
		"\n"+
		"\n"+
		"\n"+
		"<div class=\"col-md-9\">\n"+
		"<p class=\"lead\">The API is similar to <a href=\"https://blockchain.info/api/blockc"+
		"hain_api\">Blockchain.info API</a> but is <strong>not</strong> designed as a drop-"+
		"in replacement. It also implements some of <a href=\"https://blockexplorer.com/q\">"+
		"Bitcoin Block Explorer Query API</a> endpoints.</p>\n"+
		"\n"+
		"<p>You can get almost any page as <strong>JSON</strong>, just prepend <code>/api/"+
		"v1</code> before the first slash, e.g. <code>https://btcplex.com/block/0000000000"+
		"0034fa21051368f5a197e65239efb2f99a831615bbdd499429ab94</code> become <code>https:"+
		"//btcplex.com/api/v1/block/00000000000034fa21051368f5a197e65239efb2f99a831615bbdd"+
		"499429ab94</code>. It works for blocks, transactions and addresses.</p>\n"+
		"\n"+
		"<h3>Path</h3>\n"+
		"\n"+
		"<p>For this documentation, we will assume every request begins with the above pat"+
		"h.</p>\n"+
		"\n"+
		"<pre>https://btcplex.com/api/v1/</pre>\n"+
		"\n"+
		"<h3>Format</h3>\n"+
		"\n"+
		"<p>All calls are returned in <strong>JSON</strong>.</p>\n"+
		"\n"+
		"<h3>Rate limiting</h3>\n"+
		"\n"+
		"<p>The rate limit allows you to make <strong>3600 requests per hour</strong> and "+
		"implements the standard <code>X-RateLimit-*</code> headers in every API response."+
		"</p>\n"+
		"\n"+
		"<ul class=\"list-unstyled\">\n"+
		"  <li><code>X-RateLimit-Limit</code> The number of requests allowed per hour.</li"+
		">\n"+
		"  <li><code>X-RateLimit-Remaining</code> The number of requests remaining in the "+
		"current window.</li>\n"+
		"  <li><code>X-RateLimit-Reset</code> The time (in UTC epoch seconds) at which the"+
		" rate limit window resets.</li>\n"+
		"</ul>\n"+
		"\n"+
		"\n"+
		"<h3>Status Codes</h3>\n"+
		"\n"+
		"<ul class=\"list-unstyled\">\n"+
		"<li><strong>200 OK</strong> Response to a successful request.</li>\n"+
		"<li><strong>429 Too many requests</strong> Request aborted due to rate-limiting.<"+
		"/li>\n"+
		"<li><strong>500 Internal server error</strong> Something bad happened.</li>\n"+
		"</ul>\n"+
		"\n"+
		"\n"+
		"<h3>Resources</h3>\n"+
		"\n"+
		"<p style=\"margin-bottom:20px;\">All endpoints are listed here:</p>\n"+
		"\n"+
		"<h4 id=\"getblockcount\" class=\"anchor\">GET /getblockcount</h4>\n"+
		"\n"+
		"<p>Returns the current block height / number of blocks in the longest chain.</p>\n"+
		"	\n"+
		"<h5>Example request</h5>\n"+
		"\n"+
		"<pre>$ curl https://btcplex.com/api/v1/getblockcount</pre>\n"+
		"\n"+
		"<h5>Response</h5>\n"+
		"<pre>\n"+
		"275120\n"+
		"</pre>\n"+
		"\n"+
		"<h4 id=\"getblockhash_height\" class=\"anchor\">GET /getblockhash/:height</h4>\n"+
		"\n"+
		"<p>Returns the hash of a block for the given height.</p>\n"+
		"\n"+
		"<h5>Example request</h5>\n"+
		"\n"+
		"<pre>$ curl https://btcplex.com/api/v1/getblockhash/233943</pre>\n"+
		"\n"+
		"<h5>Response</h5>\n"+
		"<pre>\n"+
		"00000000000001a5e28bb1e9aad40f846d66fb7911c0db4dc32a039bb168a237\n"+
		"</pre>\n"+
		"\n"+
		"<h4 id=\"latesthash\" class=\"anchor\">GET /latesthash</h4>\n"+
		"\n"+
		"<p>Returns the hash of the latest block.</p>\n"+
		"\n"+
		"<h5>Example request</h5>\n"+
		"\n"+
		"<pre>$ curl https://btcplex.com/api/v1/latesthash</pre>\n"+
		"\n"+
		"<h5>Response</h5>\n"+
		"<pre>\n"+
		"00000000000001a5e28bb1e9aad40f846d66fb7911c0db4dc32a039bb168a237\n"+
		"</pre>\n"+
		"\n"+
		"<h4 id=\"block_hash\" class=\"anchor\">GET /block/:hash</h4>\n"+
		"\n"+
		"<p>Return block details along with transactions.</p>\n"+
		"\n"+
		"<h5>Example request</h5>\n"+
		"\n"+
		"<pre>$ curl https://btcplex.com/api/v1/block/000000000000170b01901a691a88d0bc1cde"+
		"49fe32675d920039540613e3f2d7</pre>\n"+
		"\n"+
		"<h5>Response</h5>\n"+
		"<pre class=\"pre-scrollable\">\n"+
		"{\n"+
		"  \"bits\": 453023994, \n"+
		"  \"hash\": \"000000000000170b01901a691a88d0bc1cde49fe32675d920039540613e3f2d7\", \n"+
		"  \"height\": 121426, \n"+
		"  \"mrkl_root\": \"71b01258157daeddd7e4b08bf2a149eb0878581e1a108c4ccef801867d105b17\""+
		", \n"+
		"  \"n_tx\": 10, \n"+
		"  \"next_block\": \"00000000000034fa21051368f5a197e65239efb2f99a831615bbdd499429ab94"+
		"\", \n"+
		"  \"nonce\": 702499203, \n"+
		"  \"prev_block\": \"00000000000005f34c6e10ffeb86f3073c119531629f4e2728204431fd3a6ba7"+
		"\", \n"+
		"  \"size\": 3040, \n"+
		"  \"time\": 1304344768, \n"+
		"  \"total_out\": 23308967646, \n"+
		"  \"tx\": [\n"+
		"    {\n"+
		"      \"block_hash\": \"000000000000170b01901a691a88d0bc1cde49fe32675d920039540613e3"+
		"f2d7\", \n"+
		"      \"block_height\": 121426, \n"+
		"      \"block_time\": 0, \n"+
		"      \"hash\": \"cd6f351ba4c9b1d17b367dcb72bfdab0a99dbec2c7a5122f5641ea51b01f08e1\","+
		" \n"+
		"      \"in\": null, \n"+
		"      \"lock_time\": 0, \n"+
		"      \"out\": [\n"+
		"        {\n"+
		"          \"hash\": \"1BDJGdvEbyy5v53yFspWkLz8f6Un2EYkWz\", \n"+
		"          \"value\": 5000000000\n"+
		"        }\n"+
		"      ], \n"+
		"      \"size\": 134, \n"+
		"      \"ver\": 1, \n"+
		"      \"vin_sz\": 0, \n"+
		"      \"vin_total\": 0, \n"+
		"      \"vout_sz\": 1, \n"+
		"      \"vout_total\": 5000000000\n"+
		"    }, \n"+
		"    {\n"+
		"      \"block_hash\": \"000000000000170b01901a691a88d0bc1cde49fe32675d920039540613e3"+
		"f2d7\", \n"+
		"      \"block_height\": 121426, \n"+
		"      \"block_time\": 0, \n"+
		"      \"hash\": \"407fea16950684221ec202e8b5d51f9883cf6705c6fe980b329a6a6e2bcae090\","+
		" \n"+
		"      \"in\": [\n"+
		"        {\n"+
		"          \"prev_out\": {\n"+
		"            \"address\": \"1EKc4Go9Hg7AJ85TSuw7HRNN6bHsJPteCq\", \n"+
		"            \"hash\": \"f17fc0337c6a8ca197e6bedd1395c8905326988993e994052566d97b526e"+
		"b1b3\", \n"+
		"            \"n\": 1, \n"+
		"            \"value\": 25000000\n"+
		"          }\n"+
		"        }, \n"+
		"        {\n"+
		"          \"prev_out\": {\n"+
		"            \"address\": \"1MR3oDXdsk2wuj2oP7dfMMFWcDHCcpRakR\", \n"+
		"            \"hash\": \"a5977cafded3d9a976dd1d3eb2796933f884576484d6c5b071f6cec5b41c"+
		"0a91\", \n"+
		"            \"n\": 1, \n"+
		"            \"value\": 1145000000\n"+
		"          }\n"+
		"        }, \n"+
		"        {\n"+
		"          \"prev_out\": {\n"+
		"            \"address\": \"1SgRn6cf1wkZ955stT8LsfTorFuxqpw3G\", \n"+
		"            \"hash\": \"8cd02632d840eb6424df5bd1099a5b4ee3752b8e506c3505f23790541ba9"+
		"4272\", \n"+
		"            \"n\": 0, \n"+
		"            \"value\": 1331000000\n"+
		"          }\n"+
		"        }\n"+
		"      ], \n"+
		"      \"lock_time\": 0, \n"+
		"      \"out\": [\n"+
		"        {\n"+
		"          \"hash\": \"1AnW35jwT5HLpEZ9FFHoXwzcv1Jo2Y4hBD\", \n"+
		"          \"value\": 2500000000\n"+
		"        }, \n"+
		"        {\n"+
		"          \"hash\": \"1NPcjR6vhWz1mJv212cQ1B4dqVrjMsmNZH\", \n"+
		"          \"value\": 1000000\n"+
		"        }\n"+
		"      ], \n"+
		"      \"size\": 617, \n"+
		"      \"ver\": 1, \n"+
		"      \"vin_sz\": 3, \n"+
		"      \"vin_total\": 2501000000, \n"+
		"      \"vout_sz\": 2, \n"+
		"      \"vout_total\": 2501000000\n"+
		"    }, \n"+
		"  ], \n"+
		"  \"ver\": 1\n"+
		"}\n"+
		"</pre>\n"+
		"\n"+
		"\n"+
		"<h4 id=\"checkaddress_address\" class=\"anchor\">GET /checkaddress/:address</h4>\n"+
		"\n"+
		"<p>Check a Bitcoin address for validity.</p>\n"+
		"\n"+
		"<h5>example request</h5>\n"+
		"\n"+
		"<pre>$ curl https://btcplex.com/api/v1/checkaddress/19gzwTuuZDec8JZEddQUZH9kwzqkB"+
		"fFtDa</pre>\n"+
		"\n"+
		"<h5>Response</h5>\n"+
		"\n"+
		"<pre>true</pre>\n"+
		"\n"+
		"</div>\n"+
		"</div>"))
}
