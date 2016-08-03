package main

import (
	"bytes"
	"github.com/GeoNet/weft"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"net/http"
)

var renderer blackfriday.Renderer
var extensions int
var docsIndex []byte

func init() {
	htmlFlags := 0
	// htmlFlags |= blackfriday.HTML_TOC

	renderer = blackfriday.HtmlRenderer(htmlFlags, "", "")

	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_DEFINITION_LISTS
	extensions |= blackfriday.EXTENSION_HEADER_IDS

	input, _ := ioutil.ReadFile("assets/docs/index.md")

	var b bytes.Buffer
	b.WriteString(h)
	b.Write(blackfriday.Markdown(input, renderer, extensions))
	b.WriteString(f)
	docsIndex = b.Bytes()
}

func docs(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	if res := weft.CheckQuery(r, []string{}, []string{}); !res.Ok {
		return res
	}

	if r.URL.Path != "/" {
		return weft.BadRequest("invalid path")
	}

	b.Write(docsIndex)
	return &weft.StatusOK
}

const h = `<html>
	<head>
	<meta charset="utf-8"/>
	<meta http-equiv="X-UA-Compatible" content="IE=edge"/>
	<meta name="viewport" content="width=device-width, initial-scale=1"/>
	<title>GeoNet API</title>
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap.min.css">
	<style>
	body { padding-top: 60px; }
	a.anchor { 
		display: block; position: relative; top: -60px; visibility: hidden; 
	}

	.panel-height {
		height: 150px; 
		overflow-y: scroll;
	}

	.footer {
		margin-top: 20px;
		padding: 20px 0 20px;
		border-top: 1px solid #e5e5e5;
	}

	.footer p {
		text-align: center;
	}

	#logo{position:relative;}
	#logo li{margin:0;padding:0;list-style:none;position:absolute;top:0;}
	#logo li a span
	{
		position: absolute;
		left: -10000px;
	}

	#gns li, #gns a
	{
		float: left;
		display:block;
		height: 90px;
		width: 54px;
	}

	#gns{left:-20px;height:90px;width:54px;}
	#gns{background:url('http://static.geonet.org.nz/geonet-2.0.2/images/logos.png') -0px -0px;}

	#eqc li, #eqc a
	{
		display:block;
		height: 61px;
		width: 132px;
	}

	#eqc{right:0px;height:79px;width:132px;}
	#eqc{background:url('http://static.geonet.org.nz/geonet-2.0.2/images/logos.png') -0px -312px;}

	#ccby li, #ccby a
	{
		display:block;
		height: 15px;
		width: 80px;
	}
	#ccby{left:15px;height:15px;width:80px; }
	#ccby{background:url('http://static.geonet.org.nz/geonet-2.0.2/images/logos.png') -0px -100px;}

	#geonet{
		background:url('http://static.geonet.org.nz/geonet-2.0.2/images/logos.png') 0px -249px; 
		width:137px; 
		height:53px;
		display:block;
	}


	</style>
	</head>
	<body>
	<div class="navbar navbar-inverse navbar-fixed-top" role="navigation">
	<div class="container">
	<div class="navbar-header">
	<a class="navbar-brand" href="http://geonet.org.nz">GeoNet</a>
	</div>
	</div>
	</div>

	<div class="container-fluid">`

const f = `<div id="footer" class="footer">
	<div class="row">
	<div class="col-sm-3 hidden-xs">
	<ul id="logo">
	<li id="geonet"><a target="_blank" href="http://www.geonet.org.nz"><span>GeoNet</span></a></li>
	</ul>            
	</div>

	<div class="col-sm-6">
	<p>GeoNet is a collaboration between the <a target="_blank" href="http://www.eqc.govt.nz">Earthquake Commission</a> and <a target="_blank" href="http://www.gns.cri.nz/">GNS Science</a>.</p>
	<p><a target="_blank" href="http://info.geonet.org.nz/x/loYh">about</a> | <a target="_blank" href="http://info.geonet.org.nz/x/JYAO">contact</a> | <a target="_blank" href="http://info.geonet.org.nz/x/RYAo">privacy</a> | <a target="_blank" href="http://info.geonet.org.nz/x/EIIW">disclaimer</a> </p>
	<p>GeoNet content is copyright <a target="_blank" href="http://www.gns.cri.nz/">GNS Science</a> and is licensed under a <a rel="license" target="_blank" href="http://creativecommons.org/licenses/by/3.0/nz/">Creative Commons Attribution 3.0 New Zealand License</a></p>
	</div>

	<div  class="col-sm-2 hidden-xs">
	<ul id="logo">
	<li id="eqc"><a target="_blank" href="http://www.eqc.govt.nz" ><span>EQC</span></a></li>
	</ul>
	</div>
	<div  class="col-sm-1 hidden-xs">
	<ul id="logo">
	<li id="gns"><a target="_blank" href="http://www.gns.cri.nz"><span>GNS Science</span></a></li>
	</ul>  
	</div>
	</div>

	<div class="row">
	<div class="col-sm-1 col-sm-offset-5 hidden-xs">
	<ul id="logo">
	<li id="ccby"><a href="http://creativecommons.org/licenses/by/3.0/nz/" ><span>CC-BY</span></a></li>
	</ul>
	</div>
	</div>

	</div>
	</div>
	</body>
	</html>`
