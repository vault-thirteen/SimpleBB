package gwm

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/vault-thirteen/SimpleBB/pkg/GWM/models/api"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	hh "github.com/vault-thirteen/auxie/http-helper"
)

const (
	// General methods of API.
	ApiFunctionName_GetProductVersion = "getProductVersion"
	ApiFunctionName_GetSettings       = "getSettings"

	// ACM.
	ApiFunctionName_RegisterUser           = "registerUser"
	ApiFunctionName_ApproveAndRegisterUser = "approveAndRegisterUser"
	ApiFunctionName_LogUserIn              = "logUserIn"
	ApiFunctionName_LogUserOut             = "logUserOut"
	ApiFunctionName_GetListOfLoggedUsers   = "getListOfLoggedUsers"
	ApiFunctionName_IsUserLoggedIn         = "isUserLoggedIn"
	ApiFunctionName_ChangePassword         = "changePassword"
	ApiFunctionName_ChangeEmail            = "changeEmail"
	ApiFunctionName_GetUserRoles           = "getUserRoles"
	ApiFunctionName_ViewUserParameters     = "viewUserParameters"
	ApiFunctionName_SetUserRoleAuthor      = "setUserRoleAuthor"
	ApiFunctionName_SetUserRoleWriter      = "setUserRoleWriter"
	ApiFunctionName_SetUserRoleReader      = "setUserRoleReader"
	ApiFunctionName_GetSelfRoles           = "getSelfRoles"
	ApiFunctionName_BanUser                = "banUser"
	ApiFunctionName_UnbanUser              = "unbanUser"

	// MM.
	ApiFunctionName_AddSection                  = "addSection"
	ApiFunctionName_ChangeSectionName           = "changeSectionName"
	ApiFunctionName_ChangeSectionParent         = "changeSectionParent"
	ApiFunctionName_GetSection                  = "getSection"
	ApiFunctionName_DeleteSection               = "deleteSection"
	ApiFunctionName_AddForum                    = "addForum"
	ApiFunctionName_ChangeForumName             = "changeForumName"
	ApiFunctionName_ChangeForumSection          = "changeForumSection"
	ApiFunctionName_GetForum                    = "getForum"
	ApiFunctionName_DeleteForum                 = "deleteForum"
	ApiFunctionName_AddThread                   = "addThread"
	ApiFunctionName_ChangeThreadName            = "changeThreadName"
	ApiFunctionName_ChangeThreadForum           = "changeThreadForum"
	ApiFunctionName_GetThread                   = "getThread"
	ApiFunctionName_DeleteThread                = "deleteThread"
	ApiFunctionName_AddMessage                  = "addMessage"
	ApiFunctionName_ChangeMessageText           = "changeMessageText"
	ApiFunctionName_ChangeMessageThread         = "changeMessageThread"
	ApiFunctionName_GetMessage                  = "getMessage"
	ApiFunctionName_DeleteMessage               = "deleteMessage"
	ApiFunctionName_ListThreadAndMessages       = "listThreadAndMessages"
	ApiFunctionName_ListThreadAndMessagesOnPage = "listThreadAndMessagesOnPage"
	ApiFunctionName_ListForumAndThreads         = "listForumAndThreads"
	ApiFunctionName_ListForumAndThreadsOnPage   = "listForumAndThreadsOnPage"
	ApiFunctionName_ListSectionsAndForums       = "listSectionsAndForums"
)

func (srv *Server) handlePublicApi(rw http.ResponseWriter, req *http.Request, clientIPA string) {
	if req.Method != http.MethodPost {
		srv.respondMethodNotAllowed(rw)
		return
	}

	// Check accepted MIME types.
	ok, err := hh.CheckBrowserSupportForJson(req)
	if err != nil {
		srv.respondBadRequest(rw)
		return
	}

	if !ok {
		srv.respondNotAcceptable(rw)
		return
	}

	var reqBody []byte
	reqBody, err = io.ReadAll(req.Body)
	if err != nil {
		srv.processInternalServerError(rw, err)
		return
	}

	// Check the action.
	var arwoa api.RequestWithOnlyAction
	err = json.Unmarshal(reqBody, &arwoa)
	if err != nil {
		srv.respondBadRequest(rw)
		return
	}

	if (arwoa.Action == nil) ||
		(arwoa.Parameters == nil) {
		srv.respondBadRequest(rw)
		return
	}

	var handler api.RequestHandler
	handler, ok = srv.apiHandlers[*arwoa.Action]
	if !ok {
		srv.respondNotFound(rw)
		return
	}

	var token *string
	token, err = cm.GetToken(req)
	if err != nil {
		srv.respondBadRequest(rw)
		return
	}

	var ar = &api.Request{
		Action:     arwoa.Action,
		Parameters: arwoa.Parameters,
		Authorisation: &cmr.Auth{
			UserIPA: clientIPA,
		},
	}

	if token != nil {
		ar.Authorisation.Token = *token
	}

	handler(ar, req, rw)
	return
}
