package http

import (
	mime "github.com/vault-thirteen/auxie/MIME"
)

const (
	ContentType_PlainText  = mime.TypeTextPlain
	ContentType_Json       = mime.TypeApplicationJson
	ContentType_JavaScript = mime.TypeApplicationJavascript
	ContentType_Wasm       = mime.TypeApplicationWasm
	ContentType_HtmlPage   = mime.TypeTextHtml
	ContentType_CssStyle   = mime.TypeTextCss
)
