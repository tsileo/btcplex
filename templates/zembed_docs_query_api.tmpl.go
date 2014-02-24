// THIS FILE IS AUTO-GENERATED FROM docs_query_api.tmpl
// DO NOT EDIT.

package templates

import "time"

import "camlistore.org/pkg/fileembed"

func init() {
	Files.Add("docs_query_api.tmpl", 3412, time.Unix(0, 1393190021266994961), fileembed.String("<h2>Query API Documentation</h2>\n"+
		"<div class=\"row\" style=\"margin-top:30px;\">\n"+
		"<div class=\"col-md-3\">\n"+
		"	<div class=\"list-group\" id=\"apiendpoints\">\n"+
		"  <a href=\"#getblockcount\" class=\"list-group-item\">GET /getblockcount</a>\n"+
		"  <a href=\"#getblockhash_height\" class=\"list-group-item\">GET /getblockhash/:heigh"+
		"t</a>\n"+
		"  <a href=\"#latesthash\" class=\"list-group-item\">GET /latesthash</a>\n"+
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
		"<div class=\"col-md-9\">\n"+
		"<p class=\"lead\">The Query API implements some of <a href=\"https://blockexplorer.c"+
		"om/q\">Bitcoin Block Explorer Query API</a> endpoints (also compatible with <a hre"+
		"f=\"https://blockchain.info/api/blockchain_api\">Blockchain.info API</a>).</p>\n"+
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
		"<p>All calls are returned in <strong>plain text</strong>.</p>\n"+
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
