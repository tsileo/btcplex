// THIS FILE IS AUTO-GENERATED FROM unconfirmed-transactions.tmpl
// DO NOT EDIT.

package templates

import "time"

import "camlistore.org/pkg/fileembed"

func init() {
	Files.Add("unconfirmed-transactions.tmpl", 1867, time.Unix(0, 1393190022006998620), fileembed.String("<h2><span id=\"unconfirmedcnt\"></span> Unconfirmed Transactions</h2>\n"+
		"\n"+
		"<p class=\"lead\">50 latest transactions waiting to be included in a block, <strong"+
		">updated in real time</strong>.</p>\n"+
		"\n"+
		"<div class=\"well\" id=\"waiting\">\n"+
		"  <p><strong>Please wait while loading...</strong></p>\n"+
		"</div>\n"+
		"\n"+
		"<div id=\"txs\">\n"+
		"      {{range .Txs}}\n"+
		"<div class=\"panel panel-default\">\n"+
		"  <!-- Default panel contents -->\n"+
		"  <div class=\"panel-heading hash\">{{.Hash}}</div>\n"+
		"  <div class=\"panel-body\">\n"+
		"    <p>{{.FirstSeenTime | formattime}} (<time datetime=\"{{.FirstSeenTime | format"+
		"iso}}\">{{.FirstSeenTime | formattime}}</time>)</p>\n"+
		"  </div>\n"+
		"\n"+
		"<div class=\"table-responsive\">\n"+
		"<table class=\"table table-striped table-condensed\">\n"+
		"  <thead>\n"+
		"    <tr>\n"+
		"      <th>Transaction</th>\n"+
		"      <th>Fee</th>\n"+
		"      <th>Size (kB)</th>\n"+
		"      <th>From</th>\n"+
		"      <th>To</th>\n"+
		"    </tr>\n"+
		"  </thead>\n"+
		"  <tbody>\n"+
		"      <tr>\n"+
		"        <td style=\"vertical-align:middle\"><a href=\"/tx/{{.Hash}}\" class=\"hash\">{{"+
		"cutmiddle .Hash 15}}</a></td>\n"+
		"        <td style=\"vertical-align:middle\">{{. | computefee}}</td>\n"+
		"        <td style=\"vertical-align:middle\">{{.Size | tokb}}</td>\n"+
		"        <td style=\"vertical-align:middle\">\n"+
		"        <ul class=\"list-unstyled\">\n"+
		"        {{if .TxIns}}\n"+
		"        {{range .TxIns}}\n"+
		"        <li style=\"white-space: nowrap;\"><a href=\"/address/{{.PrevOut.Address}}\" "+
		"class=\"hash\">{{.PrevOut.Address}}</a>: {{.PrevOut.Value |tobtc}}</li>\n"+
		"        {{end}}\n"+
		"        {{else}}\n"+
		"        <li style=\"white-space: nowrap;\">Generation: {{. | generationmsg}}</li>\n"+
		"        {{end}}\n"+
		"        </ul></td>\n"+
		"        \n"+
		"        <td style=\"vertical-align:middle\">\n"+
		"        <ul class=\"list-unstyled\">\n"+
		"        {{range .TxOuts}}\n"+
		"        <li style=\"white-space: nowrap;\"><a href=\"/address/{{.Addr}}\" class=\"hash"+
		"\">{{.Addr}}</a>: {{.Value |tobtc}}</li>\n"+
		"        {{end}}\n"+
		"        </ul>\n"+
		"        </td>\n"+
		"      </tr>\n"+
		"  </tbody>\n"+
		"</table>\n"+
		"</div>\n"+
		"</div>\n"+
		"      {{end}}\n"+
		"</div>"))
}
