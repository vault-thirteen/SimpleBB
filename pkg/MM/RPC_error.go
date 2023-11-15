package mm

// RPC errors.

// Error codes must not exceed 999.

// Codes.
const (
	RpcErrorCode_ForumNameIsNotSet      = 1
	RpcErrorCode_RootForumAlreadyExists = 2
	RpcErrorCode_ForumIsNotFound        = 3
	RpcErrorCode_ForumIdIsNotSet        = 4
	RpcErrorCode_RootForumCanNotBeMoved = 5
	RpcErrorCode_ThreadNameIsNotSet     = 6
	RpcErrorCode_ThreadIdIsNotSet       = 7
	RpcErrorCode_MessageTextIsNotSet    = 8
	RpcErrorCode_ThreadIsNotFound       = 9
	RpcErrorCode_MessageIdIsNotSet      = 10
	RpcErrorCode_ThreadIsNotEmpty       = 11
	RpcErrorCode_ForumHasChildren       = 12
	RpcErrorCode_ForumHasThreads        = 13
	RpcErrorCode_PageIsNotSet           = 14
	RpcErrorCode_TestError              = 15
)

// Messages.
const (
	RpcErrorMsg_ForumNameIsNotSet      = "forum name is not set"
	RpcErrorMsg_RootForumAlreadyExists = "root forum already exists"
	RpcErrorMsg_ForumIsNotFound        = "forum is not found"
	RpcErrorMsg_ForumIdIsNotSet        = "forum ID is not set"
	RpcErrorMsg_RootForumCanNotBeMoved = "root forum can not be moved"
	RpcErrorMsg_ThreadNameIsNotSet     = "thread name is not set"
	RpcErrorMsg_ThreadIdIsNotSet       = "thread ID is not set"
	RpcErrorMsg_MessageTextIsNotSet    = "message text is not set"
	RpcErrorMsg_ThreadIsNotFound       = "thread is not found"
	RpcErrorMsg_MessageIdIsNotSet      = "message ID is not set"
	RpcErrorMsg_ThreadIsNotEmpty       = "thread is not empty"
	RpcErrorMsg_ForumHasChildren       = "forum has children"
	RpcErrorMsg_ForumHasThreads        = "forum has threads"
	RpcErrorMsg_PageIsNotSet           = "page is not set"
	RpcErrorMsgF_TestError             = "test error: %s"
)
