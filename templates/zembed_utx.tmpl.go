// THIS FILE IS AUTO-GENERATED FROM utx.tmpl
// DO NOT EDIT.

package templates

import "time"

import "camlistore.org/pkg/fileembed"

func init() {
	Files.Add("utx.tmpl", 1529, time.Unix(0, 1393190021874997961), fileembed.String("\n"+
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
		"</div>"))
}
