// THIS FILE IS AUTO-GENERATED FROM tx.tmpl
// DO NOT EDIT.

package templates

import "time"

import "camlistore.org/pkg/fileembed"

func init() {
	Files.Add("tx.tmpl", 2780, time.Unix(0, 1393190021702997130), fileembed.String("{{$unconfirmed := .TxUnconfirmed}}\n"+
		"{{with .Tx}}\n"+
		"<h2>Transaction <small class=\"mono\">{{cutmiddle .Hash 15}}</small></h2>\n"+
		"\n"+
		"<dl class=\"dl-horizontal\">\n"+
		"  <dt>Hash</dt>\n"+
		"  <dd class=\"hash\">{{.Hash}}</dd>\n"+
		"\n"+
		"  {{if $unconfirmed}}\n"+
		"\n"+
		"  <dt>Confirmations</dt>\n"+
		"  <dd class=\"text-danger\"><span class=\"glyphicon glyphicon-warning-sign\"></span> "+
		"<strong>Unconfirmed transaction</strong></dd>\n"+
		"\n"+
		"  <dt>Time</dt>\n"+
		"  <dd>{{.FirstSeenTime | formattime}} (<time datetime=\"{{.FirstSeenTime | formati"+
		"so}}\">{{.FirstSeenTime | formattime}}</time>)</dd>\n"+
		"\n"+
		"  {{else}}\n"+
		"\n"+
		"  <dt>Block Hash</dt>\n"+
		"  <dd class=\"hash\"><a href=\"/block/{{.BlockHash}}\">{{.BlockHash}}</a></dd>\n"+
		"\n"+
		"  <dt>Block Height</dt>\n"+
		"  <dd><a href=\"/block/{{.BlockHash}}\">{{.BlockHeight}}</a></dd>\n"+
		"\n"+
		"  <dt>Block Time</dt>\n"+
		"  <dd>{{.BlockTime | formattime}} (<time datetime=\"{{.BlockTime | formatiso}}\">{{"+
		".BlockTime | formattime}}</time>)</dd>\n"+
		"\n"+
		"  <dt>Confirmations</dt>\n"+
		"  <dd>{{confirmation .BlockHeight}}</dd>\n"+
		"\n"+
		"  {{end}}\n"+
		"\n"+
		"  {{if .TxIns}}\n"+
		"  <dt>Number of Input</dt>\n"+
		"  <dd>{{.TxInCnt}}</dd>\n"+
		"\n"+
		"  <dt>Total Input</dt>\n"+
		"  <dd>{{.TotalIn | tobtc}}</dd>\n"+
		"  {{else}}\n"+
		"\n"+
		"  <dt>Reward</dt>\n"+
		"  <dd>{{ . | generationmsg}}</dd>\n"+
		"\n"+
		"  {{end}}\n"+
		"  <dt>Number of Output</dt>\n"+
		"  <dd>{{.TxOutCnt}}</dd>\n"+
		"\n"+
		"  <dt>Total Output</dt>\n"+
		"  <dd>{{.TotalOut | tobtc}}</dd>\n"+
		"\n"+
		"  <dt>Fee</dt>\n"+
		"  <dd>{{. | computefee}}</dd>\n"+
		"\n"+
		"  <dt>Size</dt>\n"+
		"  <dd>{{.Size |tokb}} KB</dd>\n"+
		"\n"+
		"  <dt class=\"text-muted\">API</dt>\n"+
		"  <dd><a class=\"text-muted\" href=\"/api/v1/tx/{{.Hash}}\">JSON</a></dd>\n"+
		"\n"+
		"</dl>\n"+
		"\n"+
		"<h3>Inputs</h3>\n"+
		"\n"+
		"{{if .TxIns}}\n"+
		"<div class=\"table-responsive\">\n"+
		"<table class=\"table table-striped table-condensed\">\n"+
		"  <thead>\n"+
		"    <tr>\n"+
		"      <th>Index</th>\n"+
		"      <th>Previous output</th>\n"+
		"      <th>From</th>\n"+
		"      <th>Amount</th>\n"+
		"    </tr>\n"+
		"  </thead>\n"+
		"  <tbody>\n"+
		"{{range $index, $txi := .TxIns}}\n"+
		"<tr>\n"+
		"<td>{{$index}}</td>\n"+
		"<td class=\"hash\"><a href=\"/tx/{{$txi.PrevOut.Hash}}#out{{$txi.PrevOut.Vout}}\">{{$"+
		"txi.PrevOut | formatprevout}}</a></td>\n"+
		"<td class=\"hash\"><a href=\"/address/{{$txi.PrevOut.Address}}\" name=\"in{{$index}}\">"+
		"{{$txi.PrevOut.Address}}</a></td>\n"+
		"<td>{{$txi.PrevOut.Value | tobtc}}</td>\n"+
		"</tr>\n"+
		"{{end}}\n"+
		"  </tbody>\n"+
		"</table>\n"+
		"</div>\n"+
		"{{else}}\n"+
		"<p>Generation: {{ . | generationmsg}}</p>\n"+
		"{{end}}\n"+
		"<h3>Outputs</h3>\n"+
		"\n"+
		"<div class=\"table-responsive\">\n"+
		"<table class=\"table table-striped table-condensed\">\n"+
		"  <thead>\n"+
		"    <tr>\n"+
		"      <th>Index</th>\n"+
		"      <th>To</th>\n"+
		"      <th>Amount</th>\n"+
		"      <th>Spent</th>\n"+
		"    </tr>\n"+
		"  </thead>\n"+
		"{{range $index, $txo := .TxOuts}}\n"+
		"<tr>\n"+
		"<td>{{$index}}</td>\n"+
		"<td><a href=\"/address/{{$txo.Addr}}\" name=\"out{{$index}}\" class=\"hash\">{{$txo.Add"+
		"r}}</a></td>\n"+
		"<td>{{$txo.Value | tobtc}}</td>	\n"+
		"<td>{{if $txo.Spent.Spent}}\n"+
		"\n"+
		"<a href=\"/tx/{{$txo.Spent.InputHash}}#in{{$txo.Spent.InputIndex}}\">Spent at block"+
		" {{$txo.Spent.BlockHeight}}</a>\n"+
		"\n"+
		"{{else}}\n"+
		"Unspent\n"+
		"{{end}}</td>\n"+
		"</tr>\n"+
		"{{end}}\n"+
		"  </tbody>\n"+
		"</table>\n"+
		"</div>\n"+
		"{{end}}"))
}
