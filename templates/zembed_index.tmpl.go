// THIS FILE IS AUTO-GENERATED FROM index.tmpl
// DO NOT EDIT.

package templates

import "time"

import "camlistore.org/pkg/fileembed"

func init() {
	Files.Add("index.tmpl", 663, time.Unix(0, 1393190021954998359), fileembed.String("{{with .Blocks}}\n"+
		"<div class=\"table-responsive\">\n"+
		"<table class=\"table table-striped table-condensed\">\n"+
		"<thead>\n"+
		"<tr>\n"+
		"	<th>Height</th>\n"+
		"	<th>Hash</th>\n"+
		"	<th>Time</th>\n"+
		"	<th>Transactions</th>\n"+
		"	<th>Total BTC</th>\n"+
		"	<th>Size (KB)</th>\n"+
		"</tr>\n"+
		"</thead>\n"+
		"<tbody>\n"+
		"{{range .}}\n"+
		"<tr>\n"+
		"	<td>{{.Height}}</td>\n"+
		"	<td><a href=\"/block/{{.Hash}}\" class=\"hash\">{{.Hash}}</a></td>\n"+
		"	<td>{{.BlockTime | formattime}} (<time datetime=\"{{.BlockTime | formatiso}}\"></t"+
		"ime>)</td>\n"+
		"	<td>{{.TxCnt}}</td>\n"+
		"	<td>{{.TotalBTC | tobtc}}</td>\n"+
		"	<td>{{.Size | tokb}}</td>\n"+
		"</tr>\n"+
		"{{end}}\n"+
		"</tbody>\n"+
		"</table>\n"+
		"</div>\n"+
		"{{end}}\n"+
		"\n"+
		"<ul class=\"pager\">\n"+
		"<li class=\"next\">\n"+
		"<a href=\"/blocks/{{.LastHeight}}\">More...</a>\n"+
		"</li>\n"+
		"</ul>"))
}
