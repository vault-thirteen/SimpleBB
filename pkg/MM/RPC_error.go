package mm

// RPC errors.

// Error codes must not exceed 999.

// Codes.
const (
	RpcErrorCode_Exception                = 1
	RpcErrorCode_SectionNameIsNotSet      = 2
	RpcErrorCode_RootSectionAlreadyExists = 3
	RpcErrorCode_SectionIsNotFound        = 4
	RpcErrorCode_SectionIdIsNotSet        = 5
	RpcErrorCode_SectionHasChildren       = 6
	RpcErrorCode_RootSectionCanNotBeMoved = 7
	RpcErrorCode_ForumNameIsNotSet        = 8
	RpcErrorCode_ForumIsNotFound          = 9
	RpcErrorCode_ForumIdIsNotSet          = 10
	RpcErrorCode_ForumHasThreads          = 11
	RpcErrorCode_ThreadNameIsNotSet       = 12
	RpcErrorCode_ThreadIdIsNotSet         = 13
	RpcErrorCode_ThreadIsNotFound         = 14
	RpcErrorCode_ThreadIsNotEmpty         = 15
	RpcErrorCode_MessageTextIsNotSet      = 16
	RpcErrorCode_MessageIdIsNotSet        = 17
	RpcErrorCode_IncompatibleChildType    = 18
	RpcErrorCode_PageIsNotSet             = 19
	RpcErrorCode_TestError                = 20
)

// Messages.
const (
	RpcErrorMsg_Exception                = "exception"
	RpcErrorMsg_SectionNameIsNotSet      = "section name is not set"
	RpcErrorMsg_RootSectionAlreadyExists = "root section already exists"
	RpcErrorMsg_SectionIsNotFound        = "section is not found"
	RpcErrorMsg_SectionIdIsNotSet        = "section ID is not set"
	RpcErrorMsg_SectionHasChildren       = "section has children"
	RpcErrorMsg_RootSectionCanNotBeMoved = "root section can not be moved"
	RpcErrorMsg_ForumNameIsNotSet        = "forum name is not set"
	RpcErrorMsg_ForumIsNotFound          = "forum is not found"
	RpcErrorMsg_ForumIdIsNotSet          = "forum ID is not set"
	RpcErrorMsg_ForumHasThreads          = "forum has threads"
	RpcErrorMsg_ThreadNameIsNotSet       = "thread name is not set"
	RpcErrorMsg_ThreadIdIsNotSet         = "thread ID is not set"
	RpcErrorMsg_ThreadIsNotFound         = "thread is not found"
	RpcErrorMsg_ThreadIsNotEmpty         = "thread is not empty"
	RpcErrorMsg_MessageTextIsNotSet      = "message text is not set"
	RpcErrorMsg_MessageIdIsNotSet        = "message ID is not set"
	RpcErrorMsg_IncompatibleChildType    = "incompatible child type"
	RpcErrorMsg_PageIsNotSet             = "page is not set"
	RpcErrorMsgF_TestError               = "test error: %s"
)
