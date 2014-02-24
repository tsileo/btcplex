// THIS FILE IS AUTO-GENERATED FROM address.tmpl
// DO NOT EDIT.

package templates

import "time"

import "camlistore.org/pkg/fileembed"

func init() {
	Files.Add("address.tmpl", 4447, time.Unix(0, 1393190021894998081), fileembed.String("{{with .AddressData}}\n"+
		"  <h2>Address <small class=\"mono\">{{.Address}}</small></h2>\n"+
		"\n"+
		"  {{$addr := .Address}}\n"+
		"\n"+
		"  <dl class=\"dl-horizontal\">\n"+
		"    <dt>Address</dt>\n"+
		"    <dd class=\"hash\">{{.Address}}</dd>\n"+
		"\n"+
		"    <dt>Transactions</dt>\n"+
		"    <dd>{{.TxCnt}}</dd>\n"+
		"\n"+
		"    {{if .TxCnt}}\n"+
		"\n"+
		"    <dt>Received Transactions</dt>\n"+
		"    <dd>{{.ReceivedCnt}}</dd>\n"+
		"\n"+
		"    <dt>Total Received</dt>\n"+
		"    <dd>{{.TotalReceived | tobtc}}</dd>\n"+
		"\n"+
		"    <dt>Sent Transactions</dt>\n"+
		"    <dd>{{.SentCnt}}</dd>\n"+
		"\n"+
		"    <dt>Total Sent</dt>\n"+
		"    <dd>{{.TotalSent | tobtc}}</dd>\n"+
		"\n"+
		"    {{end}}\n"+
		"\n"+
		"    <dt>Final Balance</dt>\n"+
		"    <dd>{{.FinalBalance | tobtc}}</dd>\n"+
		"\n"+
		"    <dt class=\"text-muted\">QR Code</dt>\n"+
		"    <dd><a href=\"\" class=\"text-muted\" data-toggle=\"modal\" data-target=\"#addressQR"+
		"CodeModal\">Display</a></dd>\n"+
		"\n"+
		"    <dt class=\"text-muted\">API</dt>\n"+
		"    <dd><a class=\"text-muted\" href=\"/api/v1/address/{{.Address}}\">JSON</a></dd>\n"+
		"\n"+
		"   </dl>\n"+
		"\n"+
		"{{if .Txs}}\n"+
		"<h3>Transactions</h3>\n"+
		"  <div class=\"table-responsive\">\n"+
		"  <table class=\"table table-striped table-condensed\">\n"+
		"    <thead>\n"+
		"      <tr>\n"+
		"        <th>Transaction</th>\n"+
		"        <th>Block</th>\n"+
		"        <th>Time</th>\n"+
		"        <th>From</th>\n"+
		"        <th>To</th>\n"+
		"        <th>Amount</th>\n"+
		"      </tr>\n"+
		"    </thead>\n"+
		"    <tbody>\n"+
		"        {{range .Txs}}\n"+
		"        <tr>\n"+
		"          <td style=\"vertical-align:middle\"><a href=\"/tx/{{.Hash}}\" class=\"hash\">"+
		"{{cutmiddle .Hash 6}}</a></td>\n"+
		"          <td style=\"vertical-align:middle\"><a href=\"/block/{{.BlockHash}}\">{{.Bl"+
		"ockHeight}}</a></td>\n"+
		"          <td style=\"vertical-align:middle\">{{.BlockTime | formattime}} (<time da"+
		"tetime=\"{{.BlockTime | formatiso}}\">{{.BlockTime | formattime}}</time>)</td>\n"+
		"          <td style=\"vertical-align:middle\">\n"+
		"          \n"+
		"      <ul class=\"list-unstyled\">\n"+
		"\n"+
		"\n"+
		"          {{if .TxAddressInfo.InTxIn}}\n"+
		"        \n"+
		"        <li style=\"white-space: nowrap;\"><span class=\"hash\">{{$addr}}</span></li>\n"+
		"\n"+
		"          {{else}}\n"+
		"          \n"+
		"\n"+
		"          {{if .TxIns}}\n"+
		"          {{range .TxIns}}\n"+
		"          <li style=\"white-space: nowrap;\"><a href=\"/address/{{.PrevOut.Address}}"+
		"\" class=\"hash\">{{.PrevOut.Address}}</a>: {{.PrevOut.Value |tobtc}}</li>\n"+
		"          {{end}}\n"+
		"          {{else}}\n"+
		"          <li style=\"white-space: nowrap;\">Generation: {{. | generationmsg}}</li>\n"+
		"          {{end}}\n"+
		"          </ul>\n"+
		"\n"+
		"          {{end}}\n"+
		"\n"+
		"\n"+
		"          </td>\n"+
		"          \n"+
		"          <td style=\"vertical-align:middle\">\n"+
		"          <ul class=\"list-unstyled\">\n"+
		"\n"+
		"          {{if .TxAddressInfo.InTxOut}}\n"+
		"          \n"+
		"          <li style=\"white-space: nowrap;\"><span class=\"hash\">{{$addr}}</span></l"+
		"i>\n"+
		"\n"+
		"          {{else}}\n"+
		"          \n"+
		"          {{range .TxOuts}}\n"+
		"          <li style=\"white-space: nowrap;\"><a href=\"/address/{{.Addr}}\" class=\"ha"+
		"sh\">{{.Addr}}</a>: {{.Value |tobtc}}</li>\n"+
		"          {{end}}\n"+
		"\n"+
		"\n"+
		"          </ul>\n"+
		"          {{end}}\n"+
		"          </td>\n"+
		"          <td style=\"vertical-align:middle\">{{.TxAddressInfo.Value | inttobtc}}</"+
		"td>\n"+
		"        </tr>\n"+
		"        {{end}}\n"+
		"    </tbody>\n"+
		"  </table>\n"+
		"  </div>\n"+
		"\n"+
		"{{else}}\n"+
		"<p class=\"lead\">This address hasn't been used on the network yet.</p>\n"+
		"{{end}}\n"+
		"\n"+
		"\n"+
		"<!-- Modal -->\n"+
		"<div class=\"modal fade\" id=\"addressQRCodeModal\" tabindex=\"-1\" role=\"dialog\" aria-"+
		"labelledby=\"addressQRCodeModalLabel\" aria-hidden=\"true\">\n"+
		"  <div class=\"modal-dialog\">\n"+
		"    <div class=\"modal-content\">\n"+
		"      <div class=\"modal-header\">\n"+
		"        <button type=\"button\" class=\"close\" data-dismiss=\"modal\" aria-hidden=\"tru"+
		"e\">&times;</button>\n"+
		"        <h4 class=\"modal-title\" id=\"addressQRCodeModalLabel\">{{.Address}} QR Code"+
		"</h4>\n"+
		"      </div>\n"+
		"      <div class=\"modal-body\">\n"+
		"        <div id=\"qrcode\" data-addr=\"{{.Address}}\" style=\"text-align:center;\"></di"+
		"v>\n"+
		"      </div>\n"+
		"      <div class=\"modal-footer\">\n"+
		"        <button type=\"button\" class=\"btn btn-default\" data-dismiss=\"modal\">Close<"+
		"/button>\n"+
		"      </div>\n"+
		"    </div><!-- /.modal-content -->\n"+
		"  </div><!-- /.modal-dialog -->\n"+
		"</div><!-- /.modal -->\n"+
		"\n"+
		"{{end}}\n"+
		"\n"+
		"{{if .AddressData.Txs}}\n"+
		"<div class=\"center-block text-center\">\n"+
		"<ul class=\"pagination \">\n"+
		"{{if .PaginationData.Prev}}\n"+
		"   <li><a href=\"?page={{.PaginationData.Prev}}\">&laquo;</a></li>\n"+
		"{{else}}\n"+
		"  <li class=\"disabled\"><a href=\"#\">&laquo;</a></li>\n"+
		"{{ end }}\n"+
		"{{$cpage := .PaginationData.CurrentPage}}\n"+
		" {{range $index, $tmp := .PaginationData.Pages}}\n"+
		" {{$page := iadd $index 1}}\n"+
		"  <li {{if eq $page $cpage}}class=\"active\"{{end}}><a href=\"?page={{$page}}\">{{$pa"+
		"ge}}</a></li>\n"+
		" {{end}}\n"+
		"\n"+
		"{{if .PaginationData.Next}}\n"+
		"   <li><a href=\"?page={{.PaginationData.Next}}\">&raquo;</a></li>\n"+
		"{{else}}\n"+
		"  <li class=\"disabled\"><a href=\"#\">&raquo;</a></li>\n"+
		"{{ end }}\n"+
		"</ul>\n"+
		"</div>\n"+
		"\n"+
		"{{end}}"))
}
