package c

import (
	"fmt"
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
)

// List of supported functions.
const (
	// Ping.
	FuncPing = "Ping"

	// Section.
	FuncAddSection          = "AddSection"
	FuncChangeSectionName   = "ChangeSectionName"
	FuncChangeSectionParent = "ChangeSectionParent"
	FuncGetSection          = "GetSection"
	FuncDeleteSection       = "DeleteSection"

	// Forum.
	FuncAddForum           = "AddForum"
	FuncChangeForumName    = "ChangeForumName"
	FuncChangeForumSection = "ChangeForumSection"
	FuncGetForum           = "GetForum"
	FuncDeleteForum        = "DeleteForum"

	// Thread.
	FuncAddThread         = "AddThread"
	FuncChangeThreadName  = "ChangeThreadName"
	FuncChangeThreadForum = "ChangeThreadForum"
	FuncGetThread         = "GetThread"
	FuncDeleteThread      = "DeleteThread"

	// Message.
	FuncAddMessage          = "AddMessage"
	FuncChangeMessageText   = "ChangeMessageText"
	FuncChangeMessageThread = "ChangeMessageThread"
	FuncGetMessage          = "GetMessage"
	FuncDeleteMessage       = "DeleteMessage"

	// Composite objects.
	FuncListThreadAndMessages       = "ListThreadAndMessages"
	FuncListThreadAndMessagesOnPage = "ListThreadAndMessagesOnPage"
	FuncListForumAndThreads         = "ListForumAndThreads"
	FuncListForumAndThreadsOnPage   = "ListForumAndThreadsOnPage"
	FuncListSectionsAndForums       = "ListSectionsAndForums"

	// Other.
	FuncShowDiagnosticData = "ShowDiagnosticData"
	FuncTest               = "Test"
)

func NewClient(host string, port uint16, path string, enableSelfSignedCertificate bool) (c *cc.Client, err error) {
	dsn := fmt.Sprintf("https://%s:%d%s", host, port, path)
	return cc.NewClient(dsn, enableSelfSignedCertificate)
}
