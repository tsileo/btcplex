// THIS FILE IS AUTO-GENERATED FROM layout.tmpl
// DO NOT EDIT.

package templates

import "time"

import "camlistore.org/pkg/fileembed"

func init() {
	Files.Add("layout.tmpl", 6160, time.Unix(0, 1393190021858997892), fileembed.String("<!DOCTYPE html>\n"+
		"<html lang=\"en\">\n"+
		"  <head>\n"+
		"    <meta charset=\"utf-8\">\n"+
		"    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n"+
		"    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n"+
		"    <meta name=\"description\" content=\"{{.Description}}\">\n"+
		"    <meta name=\"author\" content=\"Thomas Sileo\">\n"+
		"    <link rel=\"shortcut icon\" href=\"https://btcplex.s3.amazonaws.com/favicon.ico\""+
		">\n"+
		"\n"+
		"    <title>{{.Title}} - Bitcoin Block Chain Explorer - BTCPlex</title>\n"+
		"\n"+
		"    <link href=\"//netdna.bootstrapcdn.com/bootswatch/3.0.3/yeti/bootstrap.min.css"+
		"\" rel=\"stylesheet\">\n"+
		"    <link href=\"//fonts.googleapis.com/css?family=Inconsolata:400,700\" rel=\"style"+
		"sheet\">\n"+
		"    <link rel=\"stylesheet\" href=\"//cdnjs.cloudflare.com/ajax/libs/jquery-jgrowl/1"+
		".2.12/jquery.jgrowl.min.css\">\n"+
		"  \n"+
		"    <style type=\"text/css\">\n"+
		"    #main {\n"+
		"    	margin-top: 60px;\n"+
		"    }\n"+
		"    a:target {\n"+
		"      color:#001f3f;\n"+
		"    }\n"+
		"\n"+
		"    dl {\n"+
		"      margin: 30px 0;\n"+
		"    }\n"+
		"\n"+
		"    span.address, .hash {\n"+
		"      font-family: \"Inconsolata\";\n"+
		"      font-size: 1.15em;\n"+
		"    }\n"+
		"\n"+
		"    .mono {\n"+
		"      font-family: \"Inconsolata\";\n"+
		"    }\n"+
		"\n"+
		"    .anchor:before {\n"+
		"   content: \"\";\n"+
		"   display: block;\n"+
		"   height: 50px;\n"+
		"   margin: -30px 0 0;\n"+
		"}\n"+
		"\n"+
		"#footer {\n"+
		"  margin-top:10px;color:#999;\n"+
		"}\n"+
		"\n"+
		"#footer > div {\n"+
		"  padding:20px 0;\n"+
		"}\n"+
		"\n"+
		"\n"+
		"@media (min-width: 979px) {\n"+
		"    #apiendpoints {\n"+
		"      position:fixed;width:260px;\n"+
		"    }\n"+
		"  }\n"+
		"\n"+
		"    </style>\n"+
		"    <!-- HTML5 shim and Respond.js IE8 support of HTML5 elements and media querie"+
		"s -->\n"+
		"    <!--[if lt IE 9]>\n"+
		"      <script src=\"https://oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js\"></sc"+
		"ript>\n"+
		"      <script src=\"https://oss.maxcdn.com/libs/respond.js/1.3.0/respond.min.js\"><"+
		"/script>\n"+
		"    <![endif]-->\n"+
		"  </head>\n"+
		"\n"+
		"  <body>\n"+
		"\n"+
		"    <div class=\"navbar navbar-inverse navbar-fixed-top\" role=\"navigation\">\n"+
		"      <div class=\"container\">\n"+
		"        <div class=\"navbar-header\">\n"+
		"          <button type=\"button\" class=\"navbar-toggle\" data-toggle=\"collapse\" data"+
		"-target=\".navbar-collapse\">\n"+
		"            <span class=\"sr-only\">Toggle navigation</span>\n"+
		"            <span class=\"icon-bar\"></span>\n"+
		"            <span class=\"icon-bar\"></span>\n"+
		"            <span class=\"icon-bar\"></span>\n"+
		"          </button>\n"+
		"          <a class=\"navbar-brand\" href=\"/\">BTCplex</a>\n"+
		"        </div>\n"+
		"        <div class=\"collapse navbar-collapse\">\n"+
		"          <ul class=\"nav navbar-nav\">\n"+
		"            <li{{if eq .Menu \"latest_blocks\"}} class=\"active\"{{end}}><a href=\"/\">"+
		"Latest blocks</a></li>\n"+
		"            <li{{if eq .Menu \"utxs\"}} class=\"active\"{{end}}><a href=\"/unconfirmed"+
		"-transactions\">Unconfirmed transactions</a></li>\n"+
		"        \n"+
		"            <li class=\"dropdown {{if eq .Menu \"api\"}}active{{end}}\">\n"+
		"          <a href=\"#\" class=\"dropdown-toggle\" data-toggle=\"dropdown\">API <b class"+
		"=\"caret\"></b></a>\n"+
		"          <ul class=\"dropdown-menu\">\n"+
		"            <li><a href=\"/api\">Overview</a></li>\n"+
		"            <li><a href=\"/docs/query_api\">Query API</a></li>\n"+
		"            <li><a href=\"/docs/rest_api\">REST API</a></li>\n"+
		"            <li><a href=\"/docs/sse_api\">Server-Sent Events API</a></li>\n"+
		"          </ul>\n"+
		"        </li>\n"+
		"            <li{{if eq .Menu \"about\"}} class=\"active\"{{end}}><a href=\"/about\">Abo"+
		"ut</a></li>\n"+
		"          </ul>\n"+
		"\n"+
		"        <form class=\"navbar-form navbar-right\" role=\"search\" method=\"post\" action"+
		"=\"/search\">\n"+
		"  <div class=\"form-group\">\n"+
		"    <input type=\"text\" class=\"form-control\" placeholder=\"Search block, tx or addr"+
		"ess\" name=\"q\">\n"+
		"  </div>\n"+
		"</form>\n"+
		"\n"+
		"\n"+
		"<ul class=\"nav navbar-nav navbar-right\">\n"+
		"            {{if .Price}}\n"+
		"            <li><p class=\"navbar-text\">1BTC = <span id=\"price\">{{.Price}}</span>U"+
		"SD</p></li>\n"+
		"            {{end}}\n"+
		"          </ul>\n"+
		"\n"+
		"        </div><!--/.nav-collapse -->\n"+
		"\n"+
		"      </div>\n"+
		"    </div>\n"+
		"\n"+
		"    <div class=\"container\" id=\"main\">\n"+
		"\n"+
		"    {{if .Error}}\n"+
		"    <div class=\"alert alert-danger\">\n"+
		"      <strong>Oops!</strong> {{.Error}}.\n"+
		"    </div>\n"+
		"    {{end}}\n"+
		"\n"+
		"      {{ yield }}\n"+
		"\n"+
		"      <footer id=\"footer\">\n"+
		"      <div class=\"pull-left\">\n"+
		"        \xc2\xa9 2014 Thomas Sileo <a href=\"https://twitter.com/trucsdedev\">@trucsdedev"+
		"</a> / <a href=\"http://thomassileo.com/\">thomassileo.com</a>\n"+
		"      </div>\n"+
		"\n"+
		"      <div class=\"pull-right\">\n"+
		"        <a href=\"https://github.com/tsileo/btcplex\">BTCplex on GitHub</a> | <a hr"+
		"ef=\"mailto:contact@btcplex.com\">Feedback</a> | Donations: <a href=\"/address/19gzw"+
		"TuuZDec8JZEddQUZH9kwzqkBfFtDa\" class=\"hash\">19gzwTuuZDec8JZEddQUZH9kwzqkBfFtDa</a"+
		">  \n"+
		"      </div>\n"+
		"      \n"+
		"      </footer>\n"+
		"\n"+
		"    </div><!-- /.container -->\n"+
		"\n"+
		"    <script src=\"//cdnjs.cloudflare.com/ajax/libs/jquery/2.0.3/jquery.min.js\"></s"+
		"cript>\n"+
		"    <script src=\"//cdnjs.cloudflare.com/ajax/libs/jquery-jgrowl/1.2.12/jquery.jgr"+
		"owl.min.js\"></script>\n"+
		"    <script src=\"//netdna.bootstrapcdn.com/bootstrap/3.0.3/js/bootstrap.min.js\"><"+
		"/script>\n"+
		"    <script src=\"//cdnjs.cloudflare.com/ajax/libs/jquery-timeago/1.1.0/jquery.tim"+
		"eago.min.js\"></script>\n"+
		"    <script type=\"text/javascript\">\n"+
		"    $(function()\xc2\xa0{\n"+
		"      if ($('#apiendpoints').length == 1) {\n"+
		"        $(document).scroll(function() {\n"+
		"          if ($(document).scrollTop() > 50) { $('#apiendpoints').css('top', '60px"+
		"'); } else { $('#apiendpoints').css('top', '') };\n"+
		"        })\n"+
		"      };\n"+
		"      $(\"time\").timeago();\n"+
		"      if ($('#qrcode').length == 1) {\n"+
		"        $('#addressQRCodeModal').on('show.bs.modal', function(e) {\n"+
		"          $.getScript(\"https://cdnjs.cloudflare.com/ajax/libs/jquery.qrcode/1.0/j"+
		"query.qrcode.min.js\", function() {\n"+
		"            $('#qrcode').qrcode({text: $('#qrcode').data('addr')});\n"+
		"          });\n"+
		"        });\n"+
		"      };\n"+
		"\n"+
		"\n"+
		"    var source = new EventSource('/events');\n"+
		"    source.onmessage = function(e) {\n"+
		"      var data = JSON.parse(e.data);\n"+
		"      if (data.t == \"price\") {\n"+
		"        $.jGrowl(\"Price updated!\");\n"+
		"        $('#price').html(data.price);\n"+
		"      };\n"+
		"      if (data.t == \"height\") {\n"+
		"        $.jGrowl(\"New block found!\");\n"+
		"      };\n"+
		"    };\n"+
		"\n"+
		"    if ($(\"#unconfirmedcnt\").length == 1) {\n"+
		"      var source2 = new EventSource('/events_unconfirmed');\n"+
		"      source2.onmessage = function(e) {\n"+
		"        $('#waiting').hide();\n"+
		"        var data = JSON.parse(e.data);\n"+
		"        $('#unconfirmedcnt').html(data.cnt);\n"+
		"        $('#txs').prepend(data.tmpl);\n"+
		"        if ($('#txs').size() > 50) {\n"+
		"          $('#txs').last().remove();\n"+
		"        }\n"+
		"        $(\"time\").timeago();\n"+
		"      };\n"+
		"    };\n"+
		"\n"+
		"    });\n"+
		"    </script>\n"+
		"  </body>\n"+
		"</html>"))
}
