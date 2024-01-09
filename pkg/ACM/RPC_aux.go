package acm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/mail"
	"time"

	bpp "github.com/vault-thirteen/BytePackedPassword"
	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	rc "github.com/vault-thirteen/SimpleBB/pkg/RCS/client"
	rm "github.com/vault-thirteen/SimpleBB/pkg/RCS/models"
	sc "github.com/vault-thirteen/SimpleBB/pkg/SMTP/client"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SMTP/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	cn "github.com/vault-thirteen/SimpleBB/pkg/common/net"
	num "github.com/vault-thirteen/auxie/number"
)

// Auxiliary functions used in RPC functions.

const (
	ErrAuthData = "authorisation data error"
)

// logError logs error if debug mode is enabled.
func (srv *Server) logError(err error) {
	if err == nil {
		return
	}

	if srv.settings.SystemSettings.IsDebugMode {
		log.Println(err)
	}
}

// processDatabaseError processes a database error.
func (srv *Server) processDatabaseError(err error) {
	if err == nil {
		return
	}

	if c.IsNetworkError(err) {
		log.Println(fmt.Sprintf(c.ErrFDatabaseNetwork, err.Error()))
		*(srv.dbErrors) <- err
	} else {
		srv.logError(err)
	}

	return
}

// databaseError processes the database error and returns an RPC error.
func (srv *Server) databaseError(err error) (re *jrm1.RpcError) {
	srv.processDatabaseError(err)
	return jrm1.NewRpcErrorByUser(c.RpcErrorCode_Database, c.RpcErrorMsg_Database, err)
}

// Token-related functions.

// mustBeAuthUserIPA ensures that user's IP address is set. If it is not set,
// an error is returned and the caller of this function must stop and return
// this error.
func (srv *Server) mustBeAuthUserIPA(auth *cmr.Auth) (re *jrm1.RpcError) {
	if (auth == nil) ||
		(len(auth.UserIPA) == 0) {
		srv.incidentManager.ReportIncident(am.IncidentType_IllegalAccessAttempt, "", nil)
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_Authorisation, c.RpcErrorMsg_Authorisation, nil)
	}

	var err error
	auth.UserIPAB, err = cn.ParseIPA(auth.UserIPA)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_Authorisation, c.RpcErrorMsg_Authorisation, nil)
	}

	return nil
}

// mustBeNoAuthToken ensures that an authorisation token is not present. If the
// token is present, an error is returned and the caller of this function must
// stop and return this error.
func (srv *Server) mustBeNoAuthToken(auth *cmr.Auth) (re *jrm1.RpcError) {
	re = srv.mustBeAuthUserIPA(auth)
	if re != nil {
		return re
	}

	if len(auth.Token) > 0 {
		srv.incidentManager.ReportIncident(am.IncidentType_IllegalAccessAttempt, "", auth.UserIPAB)
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	return nil
}

// mustBeAnAuthToken ensures that an authorisation token is present and is
// valid. If the token is absent or invalid, an error is returned and the caller
// of this function must stop and return this error. User data is returned when
// token is valid.
func (srv *Server) mustBeAnAuthToken(auth *cmr.Auth) (ud *am.UserData, re *jrm1.RpcError) {
	re = srv.mustBeAuthUserIPA(auth)
	if re != nil {
		return nil, re
	}

	if len(auth.Token) == 0 {
		srv.incidentManager.ReportIncident(am.IncidentType_IllegalAccessAttempt, "", auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Authorisation, c.RpcErrorMsg_Authorisation, nil)
	}

	var err error
	ud, err = srv.getUserDataByAuthToken(auth.Token, auth.UserIPAB)
	if err != nil {
		srv.incidentManager.ReportIncident(am.IncidentType_FakeToken, "", auth.UserIPAB)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Authorisation, c.RpcErrorMsg_Authorisation, nil)
	}

	return ud, nil
}

// getUserDataByAuthToken validates the token and returns information about the
// When 'userData' is set (not null), all its fields are also set (not null).
func (srv *Server) getUserDataByAuthToken(authToken string, userIPAB net.IP) (userData *am.UserData, err error) {
	var userId, sessionId uint
	userId, sessionId, err = srv.jwtkm.ValidateToken(authToken)
	if err != nil {
		return nil, err
	}

	userData = am.NewUserData()

	userData.User, err = srv.dbo.GetUserById(userId)
	if err != nil {
		srv.processDatabaseError(err)
		return nil, err
	}

	if (userData.User == nil) ||
		(userData.User.Id != userId) ||
		(!userData.User.CanLogIn) {
		return nil, errors.New(ErrAuthData)
	}

	// Attach special user roles from settings.
	userData.User.IsModerator = srv.isUserModerator(userData.User.Id)
	userData.User.IsAdministrator = srv.isUserAdministrator(userData.User.Id)

	userData.Session, err = srv.dbo.GetSessionByUserId(userId)
	if err != nil {
		srv.processDatabaseError(err)
		return nil, err
	}

	if (userData.Session == nil) ||
		(userData.Session.Id != sessionId) ||
		(!userIPAB.Equal(userData.Session.UserIPAB)) {
		return nil, errors.New(ErrAuthData)
	}

	return userData, nil
}

// Other functions.

func (srv *Server) checkCaptcha(captchaId string, answer string) (result *rm.CheckCaptchaResult, re *jrm1.RpcError) {
	var params = rm.CheckCaptchaParams{TaskId: captchaId}
	var err error

	if len(answer) == 0 {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Captcha, c.RpcErrorMsg_Captcha, nil)
	}

	params.Value, err = num.ParseUint(answer)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Captcha, c.RpcErrorMsg_Captcha, nil)
	}

	result = new(rm.CheckCaptchaResult)
	re, err = srv.rcsServiceClient.MakeRequest(context.Background(), rc.FuncCheckCaptcha, params, result)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_RPCCall, c.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Captcha, c.RpcErrorMsg_Captcha, nil)
	}

	return result, nil
}

func (srv *Server) checkPassword(userId uint, salt []byte, userKey []byte) (ok bool, re *jrm1.RpcError) {
	var pwd *[]byte
	var err error
	pwd, err = srv.dbo.GetUserPasswordById(userId)
	if err != nil {
		return false, srv.databaseError(err)
	}
	if pwd == nil {
		return false, nil
	}

	var pwdRunes []rune
	pwdRunes, err = bpp.UnpackBytes(*pwd)
	if err != nil {
		srv.logError(err)
		return false, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Password, c.RpcErrorMsg_Password, nil)
	}

	ok, err = bpp.CheckHashKey(string(pwdRunes), salt, userKey)
	if err != nil {
		srv.logError(err)
		return false, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Password, c.RpcErrorMsg_Password, nil)
	}

	return ok, nil
}

func (srv *Server) createCaptcha() (result *rm.CreateCaptchaResult, re *jrm1.RpcError) {
	var params = rm.CreateCaptchaParams{}

	result = new(rm.CreateCaptchaResult)
	var err error
	re, err = srv.rcsServiceClient.MakeRequest(context.Background(), rc.FuncCreateCaptcha, params, result)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_RPCCall, c.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Captcha, c.RpcErrorMsg_Captcha, nil)
	}
	if result.IsImageDataReturned {
		// We do not return an image in a JSON message.
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Captcha, c.RpcErrorMsg_Captcha, nil)
	}

	return result, nil
}

func (srv *Server) createRequestIdForLogIn() (rid *string, re *jrm1.RpcError) {
	s, err := srv.ridg.CreatePassword()
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RequestIdGenerator, RpcErrorMsg_RequestIdGenerator, nil)
	}

	return s, nil
}

func (srv *Server) createRequestIdForPasswordChange() (rid *string, re *jrm1.RpcError) {
	s, err := srv.ridg.CreatePassword()
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RequestIdGenerator, RpcErrorMsg_RequestIdGenerator, nil)
	}

	return s, nil
}

func (srv *Server) createVerificationCode() (vc *string, re *jrm1.RpcError) {
	var err error
	vc, err = srv.vcg.CreatePassword()
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_VerificationCodeGenerator, RpcErrorMsg_VerificationCodeGenerator, nil)
	}

	return vc, nil
}

func (srv *Server) isCaptchaNeededForLogIn(email string) (isCaptchaNeeded bool, err error) {
	var lastBadLogInTime *time.Time
	lastBadLogInTime, err = srv.dbo.GetUserLastBadLogInTimeByEmail(email)
	if err != nil {
		return false, err
	}
	if lastBadLogInTime == nil {
		return false, nil
	}

	if time.Now().Before(lastBadLogInTime.Add(time.Duration(srv.settings.SystemSettings.LogInTryTimeout) * time.Second)) {
		return true, nil
	}

	return false, nil
}

func (srv *Server) isCaptchaNeededForEmailChange(userId uint) (isCaptchaNeeded bool, err error) {
	var lastBadActionTime *time.Time
	lastBadActionTime, err = srv.dbo.GetUserLastBadActionTimeById(userId)
	if err != nil {
		return false, err
	}
	if lastBadActionTime == nil {
		return false, nil
	}

	if time.Now().Before(lastBadActionTime.Add(time.Duration(srv.settings.SystemSettings.ActionTryTimeout) * time.Second)) {
		return true, nil
	}

	return false, nil
}

func (srv *Server) isCaptchaNeededForPasswordChange(userId uint) (isCaptchaNeeded bool, err error) {
	var lastBadActionTime *time.Time
	lastBadActionTime, err = srv.dbo.GetUserLastBadActionTimeById(userId)
	if err != nil {
		return false, err
	}
	if lastBadActionTime == nil {
		return false, nil
	}

	if time.Now().Before(lastBadActionTime.Add(time.Duration(srv.settings.SystemSettings.ActionTryTimeout) * time.Second)) {
		return true, nil
	}

	return false, nil
}

func (srv *Server) isEmailOfUserValid(email string) (re *jrm1.RpcError) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return jrm1.NewRpcErrorByUser(RpcErrorCode_EmailAddressIsNotValid, RpcErrorMsg_EmailAddressIsNotValid, nil)

	}

	return nil
}

func isPasswordAllowed(password string) (re *jrm1.RpcError) {
	_, err := bpp.IsPasswordAllowed(password)
	if err != nil {
		return jrm1.NewRpcErrorByUser(RpcErrorCode_PasswordIsNotValid, RpcErrorMsg_PasswordIsNotValid, nil)
	}

	return nil
}

func (srv *Server) isUserAdministrator(userId uint) (isAdministrator bool) {
	// While system has only few administrators, the simple array look-up is
	// faster than access to a map.

	for _, id := range srv.settings.UserRoleSettings.AdministratorIds {
		if id == userId {
			return true
		}
	}

	return false
}

func (srv *Server) isUserModerator(userId uint) (isModerator bool) {
	// While system has only few moderators, the simple array look-up is
	// faster than access to a map.

	for _, id := range srv.settings.UserRoleSettings.ModeratorIds {
		if id == userId {
			return true
		}
	}

	return false
}

func (srv *Server) sendVerificationCodeCommon(params sm.SendMessageParams) (re *jrm1.RpcError) {
	var result = new(sm.SendMessageResult)

	var err error
	re, err = srv.smtpServiceClient.MakeRequest(context.Background(), sc.FuncSendMessage, params, result)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_RPCCall, c.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return jrm1.NewRpcErrorByUser(RpcErrorCode_SmtpModule, RpcErrorMsg_SmtpModule, nil)
	}

	return nil
}

func (srv *Server) sendVerificationCodeForReg(email string, code string) (re *jrm1.RpcError) {
	var subject = fmt.Sprintf(srv.settings.MessageSettings.SubjectTemplate, srv.settings.SystemSettings.SiteName)
	var msg = fmt.Sprintf(srv.settings.MessageSettings.BodyTemplateForReg, srv.settings.SystemSettings.SiteName, code)
	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}
	return srv.sendVerificationCodeCommon(params)
}

func (srv *Server) sendVerificationCodeForLogIn(email string, code string) (re *jrm1.RpcError) {
	var subject = fmt.Sprintf(srv.settings.MessageSettings.SubjectTemplate, srv.settings.SystemSettings.SiteName)
	var msg = fmt.Sprintf(srv.settings.MessageSettings.BodyTemplateForLogIn, code)
	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}
	return srv.sendVerificationCodeCommon(params)
}

func (srv *Server) sendVerificationCodeForEmailChange(email string, code string) (re *jrm1.RpcError) {
	var subject = fmt.Sprintf(srv.settings.MessageSettings.SubjectTemplate, srv.settings.SystemSettings.SiteName)
	var msg = fmt.Sprintf(srv.settings.MessageSettings.BodyTemplateForEmailChange, code)
	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}
	return srv.sendVerificationCodeCommon(params)
}

func (srv *Server) sendVerificationCodeForPwdChange(email string, code string) (re *jrm1.RpcError) {
	var subject = fmt.Sprintf(srv.settings.MessageSettings.SubjectTemplate, srv.settings.SystemSettings.SiteName)
	var msg = fmt.Sprintf(srv.settings.MessageSettings.BodyTemplateForPwdChange, code)
	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}
	return srv.sendVerificationCodeCommon(params)
}
