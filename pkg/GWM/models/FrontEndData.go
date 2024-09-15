package models

// FrontEndData is auxiliary data for front end.
type FrontEndData struct {
	ArgonJs       FrontEndFileData
	ArgonWasm     FrontEndFileData
	BppJs         FrontEndFileData
	IndexHtmlPage FrontEndFileData
	LoaderScript  FrontEndFileData
	CssStyles     FrontEndFileData
}
