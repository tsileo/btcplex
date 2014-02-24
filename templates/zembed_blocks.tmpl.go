// THIS FILE IS AUTO-GENERATED FROM blocks.tmpl
// DO NOT EDIT.

package templates

import "time"

import "camlistore.org/pkg/fileembed"

func init() {
	Files.Add("blocks.tmpl", 1199, time.Unix(0, 1393190022046998829), fileembed.String("\n"+
		"{{$fblock := index .Blocks 0}}\n"+
		"{{$lblock := index .Blocks 29}}\n"+
		"\n"+
		"{{define \"pagination\"}}\n"+
		"{{$fblock := index .Blocks 0}}\n"+
		"{{$lblock := index .Blocks 29}}\n"+
		"\n"+
		"<ul class=\"pager\">\n"+
		"\n"+
		"<li class=\"next{{if eq $fblock.Height .LastHeight}} disabled{{end}}\"><a href=\"{{i"+
		"f eq $fblock.Height .LastHeight}}{{else}}/blocks/{{add $fblock.Height 30}}{{end}}"+
		"\" class=\"pull-right\">Next</a></li>\n"+
		"\n"+
		"<li class=\"previous\"><a href=\"/blocks/{{sub $lblock.Height 1}}\" class=\"pull-left\""+
		">Previous</a></li>\n"+
		"</ul>\n"+
		"{{end}}\n"+
		"\n"+
		"<h2>Blocks #{{$fblock.Height}} to #{{$lblock.Height}}</h2>\n"+
		"\n"+
		"{{template \"pagination\" .}}\n"+
		"\n"+
		"{{with .Blocks}}\n"+
		"<div class=\"table-responsive\">\n"+
		"<table class=\"table table-striped table-condensed\">\n"+
		"<thead>\n"+
		"<tr>\n"+
		"	<th>Height</th>\n"+
		"	<th>Hash</th>\n"+
		"	<th>Time</th>\n"+
		"	<th>Transactions</th>\n"+
		"	<th>Total BTC</th>\n"+
		"	<th>Size (kB)</th>\n"+
		"</tr>\n"+
		"</thead>\n"+
		"<tbody>\n"+
		"{{range .}}\n"+
		"<tr>\n"+
		"	<td>{{.Height}}</td>\n"+
		"	<td><a href=\"/block/{{.Hash}}\" class=\"hash\">{{.Hash}}</a></td>\n"+
		"	<td>{{.BlockTime | formattime}} (<time datetime=\"{{.BlockTime | formatiso}}\">{{."+
		"BlockTime | formattime}}</time>)</td>\n"+
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
		"{{template \"pagination\" .}}\n"+
		""))
}
