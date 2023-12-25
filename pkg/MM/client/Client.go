package c

import (
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/client"
)

// List of supported functions.
const (
	// Ping.
	FuncPing = cc.FuncPing

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
	FuncShowDiagnosticData = cc.FuncShowDiagnosticData
	FuncTest               = "Test"
)
