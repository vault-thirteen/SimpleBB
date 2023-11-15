package c

import (
	"fmt"
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
)

// List of supported functions.
const (
	FuncPing                     = "Ping"
	FuncAddForum                 = "AddForum"
	FuncChangeForumName          = "ChangeForumName"
	FuncChangeForumParent        = "ChangeForumParent"
	FuncAddThread                = "AddThread"
	FuncChangeThreadName         = "ChangeThreadName"
	FuncChangeThreadForum        = "ChangeThreadForum"
	FuncAddMessage               = "AddMessage"
	FuncChangeMessageText        = "ChangeMessageText"
	FuncChangeMessageThread      = "ChangeMessageThread"
	FuncDeleteMessage            = "DeleteMessage"
	FuncGetMessage               = "GetMessage"
	FuncGetThread                = "GetThread"
	FuncDeleteThread             = "DeleteThread"
	FuncGetForum                 = "GetForum"
	FuncDeleteForum              = "DeleteForum"
	FuncListThreadMessages       = "ListThreadMessages"
	FuncListThreadMessagesOnPage = "ListThreadMessagesOnPage"
	FuncListForumThreads         = "ListForumThreads"
	FuncListForumThreadsOnPage   = "ListForumThreadsOnPage"
	FuncListForums               = "ListForums"
	FuncShowDiagnosticData       = "ShowDiagnosticData"
	FuncTest                     = "Test" //TODO
)

func NewClient(host string, port uint16, path string, enableSelfSignedCertificate bool) (c *cc.Client, err error) {
	dsn := fmt.Sprintf("https://%s:%d%s", host, port, path)
	return cc.NewClient(dsn, enableSelfSignedCertificate)
}
