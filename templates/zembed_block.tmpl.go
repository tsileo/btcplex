// THIS FILE IS AUTO-GENERATED FROM block.tmpl
// DO NOT EDIT.

package templates

import "time"

import "camlistore.org/pkg/fileembed"

func init() {
	Files.Add("block.tmpl", 2262, time.Unix(0, 1393190021882998011), fileembed.String("{{with .Block}}\n"+
		"<h2>Block #{{.Height}}</h2>\n"+
		"\n"+
		"<dl class=\"dl-horizontal\">\n"+
		"  <dt>Hash</dt>\n"+
		"  <dd class=\"hash\">{{.Hash}}</dd>\n"+
		"\n"+
		"  <dt>Previous Block</dt>\n"+
		"  <dd><a href=\"/block/{{.Parent}}\" class=\"hash\">{{.Parent}}</a></dd>\n"+
		"\n"+
		"  {{if .Next}}\n"+
		"  <dt>Next Block</dt>\n"+
		"  <dd><a href=\"/block/{{.Next}}\" class=\"hash\">{{.Next}}</a></dd>\n"+
		"  {{end}}\n"+
		"\n"+
		"  <dt>Merkle Root</dt>\n"+
		"  <dd class=\"hash\">{{.MerkleRoot}}</dd>\n"+
		"\n"+
		"  <dt>Height</dt>\n"+
		"  <dd>{{.Height}}</dd>\n"+
		"\n"+
		"  <dt>Time</dt>\n"+
		"  <dd>{{.BlockTime | formattime}} (<time datetime=\"{{.BlockTime | formatiso}}\">{{"+
		".BlockTime | formattime}}</time>)</dd>\n"+
		"\n"+
		"  <dt>Total BTC</dt>\n"+
		"  <dd>{{.TotalBTC |tobtc}}</dd>\n"+
		"\n"+
		"  <dt>Transactions</dt>\n"+
		"  <dd>{{.TxCnt}}</dd>\n"+
		"\n"+
		"  <dt>Version</dt>\n"+
		"  <dd>{{.Version}}</dd>\n"+
		"\n"+
		"  <dt>Bits</dt>\n"+
		"  <dd>{{.Bits}}</dd>\n"+
		"\n"+
		"  <dt>Nonce</dt>\n"+
		"  <dd>{{.Nonce}}</dd>\n"+
		"\n"+
		"  <dt>Size</dt>\n"+
		"  <dd>{{ .Size | tokb }} KB</dd>\n"+
		"\n"+
		"  <dt class=\"text-muted\">API</dt>\n"+
		"  <dd><a class=\"text-muted\" href=\"/api/v1/block/{{.Hash}}\">JSON</a></dd>\n"+
		"</dl>\n"+
		"\n"+
		"<h3>Transactions</h3>\n"+
		"\n"+
		"<div class=\"table-responsive\">\n"+
		"<table class=\"table table-striped table-condensed\">\n"+
		"  <thead>\n"+
		"    <tr>\n"+
		"      <th>Transaction</th>\n"+
		"      <th>Fee</th>\n"+
		"      <th>Size (KB)</th>\n"+
		"      <th>From</th>\n"+
		"      <th>To</th>\n"+
		"    </tr>\n"+
		"  </thead>\n"+
		"  <tbody>\n"+
		"      {{range .Txs}}\n"+
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
		"      {{end}}\n"+
		"  </tbody>\n"+
		"</table>\n"+
		"</div>\n"+
		"{{end}}"))
}
