package gwm

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	"github.com/vault-thirteen/SimpleBB/pkg/GWM/models"
	"github.com/vault-thirteen/SimpleBB/pkg/GWM/models/api"
	ch "github.com/vault-thirteen/SimpleBB/pkg/common/http"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	"github.com/vault-thirteen/auxie/header"
	hh "github.com/vault-thirteen/auxie/http-helper"
)

const (
	// ACM.
	ApiFunctionName_RegisterUser                           = "registerUser"
	ApiFunctionName_GetListOfRegistrationsReadyForApproval = "getListOfRegistrationsReadyForApproval"
	ApiFunctionName_RejectRegistrationRequest              = "rejectRegistrationRequest"
	ApiFunctionName_ApproveAndRegisterUser                 = "approveAndRegisterUser"
	ApiFunctionName_LogUserIn                              = "logUserIn"
	ApiFunctionName_LogUserOut                             = "logUserOut"
	ApiFunctionName_GetListOfLoggedUsers                   = "getListOfLoggedUsers"
	ApiFunctionName_IsUserLoggedIn                         = "isUserLoggedIn"
	ApiFunctionName_ChangePassword                         = "changePassword"
	ApiFunctionName_ChangeEmail                            = "changeEmail"
	ApiFunctionName_GetListOfAllUsers                      = "getListOfAllUsers"
	ApiFunctionName_GetUserRoles                           = "getUserRoles"
	ApiFunctionName_ViewUserParameters                     = "viewUserParameters"
	ApiFunctionName_SetUserRoleAuthor                      = "setUserRoleAuthor"
	ApiFunctionName_SetUserRoleWriter                      = "setUserRoleWriter"
	ApiFunctionName_SetUserRoleReader                      = "setUserRoleReader"
	ApiFunctionName_GetSelfRoles                           = "getSelfRoles"
	ApiFunctionName_BanUser                                = "banUser"
	ApiFunctionName_UnbanUser                              = "unbanUser"

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

func (srv *Server) handlePublicSettings(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
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

	if srv.settings.SystemSettings.IsDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, srv.settings.SystemSettings.DevModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.Header().Set(header.HttpHeaderContentType, ch.ContentType_Json)

	_, err = rw.Write(srv.publicSettingsFileData)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

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

func (srv *Server) handleCaptcha(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		srv.respondMethodNotAllowed(rw)
		return
	}

	srv.captchaProxy.ServeHTTP(rw, req)
}

// handleFrontEnd serves static files for the front end part.
func (srv *Server) handleFrontEnd(rw http.ResponseWriter, req *http.Request, clientIPA string) {
	// While the number of cases is less than 10..20, the "switch" branching
	// works faster than other methods.
	switch req.URL.Path {
	case srv.frontEnd.ArgonJs.UrlPath:
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.ArgonJs)
		return

	case srv.frontEnd.ArgonWasm.UrlPath:
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.ArgonWasm)
		return

	case srv.frontEnd.BppJs.UrlPath:
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.BppJs)
		return

	case srv.frontEnd.IndexHtmlPage.UrlPath:
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.IndexHtmlPage)
		return

	case srv.frontEnd.LoaderScript.UrlPath:
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.LoaderScript)
		return

	case srv.frontEnd.CssStyles.UrlPath:
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.CssStyles)
		return

	case srv.frontEnd.AdminJs.UrlPath:
		srv.handleAdminFrontEnd(rw, req, clientIPA)
		return

	default:
		srv.respondNotFound(rw)
		return
	}
}

// handleAdminFrontEnd serves static files for the front end part for
// administrator users.
func (srv *Server) handleAdminFrontEnd(rw http.ResponseWriter, req *http.Request, clientIPA string) {
	// Step 1. Filter for fake URL.
	switch req.URL.Path {
	case srv.frontEnd.AdminHtmlPage.UrlPath: // <- /admin.html
	case srv.frontEnd.AdminJs.UrlPath: // <- /fe/admin.js
	default:
		srv.respondNotFound(rw)
		return
	}

	// Step 2. Filter for administrator role.
	ok, err := srv.checkForAdministrator(rw, req, clientIPA)
	if err != nil {
		srv.processInternalServerError(rw, err)
		return
	}
	if !ok {
		srv.respondForbidden(rw)
		return
	}

	// Step 3. Page selection.
	switch req.URL.Path {
	case srv.frontEnd.AdminHtmlPage.UrlPath: // <- /admin.html
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.AdminHtmlPage)
		return

	case srv.frontEnd.AdminJs.UrlPath: // <- /fe/admin.js
		srv.handleFrontEndStaticFile(rw, req, srv.frontEnd.AdminJs)
		return

	default:
		srv.respondNotFound(rw)
		return
	}
}

func (srv *Server) handleFrontEndStaticFile(rw http.ResponseWriter, req *http.Request, fedf models.FrontEndFileData) {
	if req.Method != http.MethodGet {
		srv.respondMethodNotAllowed(rw)
		return
	}

	if srv.settings.SystemSettings.IsDeveloperMode {
		rw.Header().Set(header.HttpHeaderAccessControlAllowOrigin, srv.settings.SystemSettings.DevModeHttpHeaderAccessControlAllowOrigin)
	}
	rw.Header().Set(header.HttpHeaderContentType, fedf.ContentType)

	_, err := rw.Write(fedf.CachedFile)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (srv *Server) checkForAdministrator(rw http.ResponseWriter, req *http.Request, clientIPA string) (ok bool, err error) {
	var token *string
	token, err = cm.GetToken(req)
	if err != nil {
		srv.respondBadRequest(rw)
		return false, err
	}

	if token == nil {
		srv.respondForbidden(rw)
		return
	}

	// Make a 'GetSelfRoles' RPC request to verify user's roles.
	action := ApiFunctionName_GetSelfRoles
	var params json.RawMessage = []byte("{}")
	var ar = &api.Request{
		Action:     &action,
		Parameters: &params,
		Authorisation: &cmr.Auth{
			UserIPA: clientIPA,
			Token:   *token,
		},
	}

	var response *api.Response
	response, err = srv.getSelfRoles(ar)
	if err != nil {
		srv.processInternalServerError(rw, err)
		return false, err
	}

	var result *am.GetSelfRolesResult
	result, ok = response.Result.(*am.GetSelfRolesResult)
	if !ok {
		err = errors.New(ErrTypeCast)
		srv.processInternalServerError(rw, err)
		return false, err
	}

	if !result.IsAdministrator {
		srv.respondForbidden(rw)
		return false, nil
	}

	return true, nil
}
