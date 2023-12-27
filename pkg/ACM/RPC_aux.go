package acm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/mail"
	"time"

	js "github.com/osamingo/jsonrpc/v2"
	bpp "github.com/vault-thirteen/BytePackedPassword"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	rc "github.com/vault-thirteen/SimpleBB/pkg/RCS/client"
	rm "github.com/vault-thirteen/SimpleBB/pkg/RCS/models"
	sc "github.com/vault-thirteen/SimpleBB/pkg/SMTP/client"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SMTP/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	"github.com/vault-thirteen/SimpleBB/pkg/common/dbo"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	cn "github.com/vault-thirteen/SimpleBB/pkg/common/net"
	num "github.com/vault-thirteen/auxie/number"
)

// Auxiliary functions used in RPC functions.

// logError logs error if debug mode is enabled.
func (srv *Server) logError(err error) {
	if srv.settings.SystemSettings.IsDebugMode {
		log.Println(err)
	}
}

// databaseError processes the database error and returns a JSON RPC error.
func (srv *Server) databaseError(err error) (jerr *js.Error) {
	if c.IsNetworkError(err) {
		log.Println(fmt.Sprintf(c.ErrFDatabaseNetwork, err.Error()))
		*(srv.dbErrors) <- err
	} else {
		srv.logError(err)
	}

	return &js.Error{Code: c.RpcErrorCode_DatabaseError, Message: c.RpcErrorMsg_DatabaseError}
}

// Token-related functions.

// mustBeAuthUserIPA ensures that user's IP address is set. If it is not set,
// an error is returned and the caller of this function must stop and return
// this error.
func (srv *Server) mustBeAuthUserIPA(auth *cmr.Auth) (jerr *js.Error) {
	if auth == nil {
		srv.incidentManager.ReportIncident(am.IncidentType_IllegalAccessAttempt, "", nil)
		return &js.Error{Code: c.RpcErrorCode_MalformedRequest, Message: c.RpcErrorMsg_MalformedRequest}
	}

	if len(auth.UserIPA) == 0 {
		// Report the incident.
		srv.incidentManager.ReportIncident(am.IncidentType_IllegalAccessAttempt, "", nil)
		return &js.Error{Code: c.RpcErrorCode_MalformedRequest, Message: c.RpcErrorMsg_MalformedRequest}
	}

	var err error
	auth.UserIPAB, err = cn.ParseIPA(auth.UserIPA)
	if err != nil {
		srv.logError(err)
		return &js.Error{Code: c.RpcErrorCode_IPAddressError, Message: fmt.Sprintf(c.RpcErrorMsgF_IPAddressError, err.Error())}
	}

	return nil
}

// mustBeNoAuthToken ensures that an authorisation token is not present. If the
// token is present, an error is returned and the caller of this function must
// stop and return this error.
func (srv *Server) mustBeNoAuthToken(auth *cmr.Auth) (jerr *js.Error) {
	jerr = srv.mustBeAuthUserIPA(auth)
	if jerr != nil {
		return jerr
	}

	if len(auth.Token) > 0 {
		srv.incidentManager.ReportIncident(am.IncidentType_IllegalAccessAttempt, "", auth.UserIPAB)
		return &js.Error{Code: RpcErrorCode_ThisActionIsNotForLoggedUsers, Message: RpcErrorMsg_ThisActionIsNotForLoggedUsers}
	}

	return nil
}

// mustBeAnAuthToken ensures that an authorisation token is present and is
// valid. If the token is absent or invalid, an error is returned and the caller
// of this function must stop and return this error. User data is returned when
// token is valid.
func (srv *Server) mustBeAnAuthToken(auth *cmr.Auth) (ud *am.UserData, jerr *js.Error) {
	jerr = srv.mustBeAuthUserIPA(auth)
	if jerr != nil {
		return nil, jerr
	}

	if len(auth.Token) == 0 {
		srv.incidentManager.ReportIncident(am.IncidentType_IllegalAccessAttempt, "", auth.UserIPAB)
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	var err error
	ud, err = srv.getUserDataByAuthToken(auth.Token, auth.UserIPAB)
	if err != nil {
		srv.logError(err)
		srv.incidentManager.ReportIncident(am.IncidentType_FakeToken, "", auth.UserIPAB)
		return nil, &js.Error{Code: c.RpcErrorCode_GetUserDataByAuthToken, Message: fmt.Sprintf(c.RpcErrorMsgF_GetUserDataByAuthToken, err.Error())}
	}

	return ud, nil
}

// getUserDataByAuthToken validates the token and returns information about the
// When 'userData' is set (not null), all its fields are also set (not null).
func (srv *Server) getUserDataByAuthToken(authToken string, userIPAB net.IP) (userData *am.UserData, err error) {
	var userId uint
	var sessionId uint
	userId, sessionId, err = srv.jwtkm.ValidateToken(authToken)
	if err != nil {
		return nil, err
	}

	userData = am.NewUserData()

	userData.User, err = srv.dbo.GetUserById(userId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if userData.User == nil {
		return nil, errors.New(c.ErrFakeJWToken)
	}

	if userData.User.Id != userId {
		return nil, errors.New(c.ErrFakeJWToken)
	}

	if !userData.User.CanLogIn {
		return nil, errors.New(c.ErrFakeJWToken)
	}

	// Attach special user roles from settings.
	userData.User.IsModerator = srv.isUserModerator(userData.User.Id)
	userData.User.IsAdministrator = srv.isUserAdministrator(userData.User.Id)

	userData.Session, err = srv.dbo.GetSessionByUserId(userId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if userData.Session == nil {
		return nil, errors.New(c.ErrFakeJWToken)
	}

	if userData.Session.Id != sessionId {
		return nil, errors.New(c.ErrFakeJWToken)
	}

	if !userIPAB.Equal(userData.Session.UserIPAB) {
		return nil, errors.New(c.ErrFakeJWToken)
	}

	return userData, nil
}

// Other functions.

func (srv *Server) checkCaptcha(captchaId string, answer string) (result *rm.CheckCaptchaResult, jerr *js.Error) {
	var params = rm.CheckCaptchaParams{TaskId: captchaId}
	var err error

	if len(answer) == 0 {
		err = errors.New(c.ErrCaptchaAnswerIsNotSet)
	}
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_RCS_CheckCaptcha, Message: fmt.Sprintf(RpcErrorMsgF_RCS_CheckCaptcha, err.Error())}
	}

	params.Value, err = num.ParseUint(answer)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_RCS_CheckCaptcha, Message: fmt.Sprintf(RpcErrorMsgF_RCS_CheckCaptcha, err.Error())}
	}

	result = new(rm.CheckCaptchaResult)
	err = srv.rcsServiceClient.MakeRequest(context.Background(), result, rc.FuncCheckCaptcha, params)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_RCS_CheckCaptcha, Message: fmt.Sprintf(RpcErrorMsgF_RCS_CheckCaptcha, err.Error())}
	}

	return result, nil
}

func (srv *Server) checkPassword(userId uint, salt []byte, userKey []byte) (ok bool, jerr *js.Error) {
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
		return false, &js.Error{Code: RpcErrorCode_PasswordCheckError, Message: fmt.Sprintf(RpcErrorMsgF_PasswordCheckError, err.Error())}
	}

	ok, err = bpp.CheckHashKey(string(pwdRunes), salt, userKey)
	if err != nil {
		srv.logError(err)
		return false, &js.Error{Code: RpcErrorCode_PasswordCheckError, Message: fmt.Sprintf(RpcErrorMsgF_PasswordCheckError, err.Error())}
	}

	return ok, nil
}

func (srv *Server) countEmailChangesByUserId(userId uint) (emailChangesCount int, err error) {
	emailChangesCount, err = srv.dbo.CountEmailChangesByUserId(userId)
	if err != nil {
		return dbo.CountOnError, err
	}

	return emailChangesCount, nil
}

func (srv *Server) countUsersWithEmail(email string) (usersCount int, err error) {
	usersCount, err = srv.dbo.CountUsersWithEmail(email)
	if err != nil {
		return dbo.CountOnError, err
	}

	return usersCount, nil
}

func (srv *Server) countUsersWithName(name string) (usersCount int, err error) {
	usersCount, err = srv.dbo.CountUsersWithName(name)
	if err != nil {
		return dbo.CountOnError, err
	}

	return usersCount, nil
}

func (srv *Server) countUsersWithEmailAbleToLogIn(email string) (usersCount int, err error) {
	usersCount, err = srv.dbo.CountUsersWithEmailAbleToLogIn(email)
	if err != nil {
		return dbo.CountOnError, err
	}

	return usersCount, nil
}

func (srv *Server) countSessionsByUserEmail(email string) (sessionsCount int, err error) {
	sessionsCount, err = srv.dbo.CountSessionsByUserEmail(email)
	if err != nil {
		return dbo.CountOnError, err
	}

	return sessionsCount, nil
}

func (srv *Server) countSessionsByUserId(userId uint) (sessionsCount int, err error) {
	sessionsCount, err = srv.dbo.CountSessionsByUserId(userId)
	if err != nil {
		return dbo.CountOnError, err
	}

	return sessionsCount, nil
}

func (srv *Server) countPreSessionsByUserEmail(email string) (preSessionsCount int, err error) {
	preSessionsCount, err = srv.dbo.CountPreSessionsByUserEmail(email)
	if err != nil {
		return dbo.CountOnError, err
	}

	return preSessionsCount, nil
}

func (srv *Server) countPasswordChangesByUserId(userId uint) (passwordChangesCount int, err error) {
	passwordChangesCount, err = srv.dbo.CountPasswordChangesByUserId(userId)
	if err != nil {
		return dbo.CountOnError, err
	}

	return passwordChangesCount, nil
}

func (srv *Server) createCaptcha() (result *rm.CreateCaptchaResult, err error) {
	var params = rm.CreateCaptchaParams{}

	result = new(rm.CreateCaptchaResult)
	err = srv.rcsServiceClient.MakeRequest(context.Background(), result, rc.FuncCreateCaptcha, params)
	if err != nil {
		return nil, err
	}

	if result.IsImageDataReturned {
		return nil, errors.New(c.ErrUnsupportedRCSMode)
	}

	return result, nil
}

func (srv *Server) createRequestIdForLogIn() (rid string, err error) {
	var s *string
	s, err = srv.ridg.CreatePassword()
	if err != nil {
		return "", err
	}

	return *s, nil
}

func (srv *Server) createRequestIdForPasswordChange() (rid string, err error) {
	var s *string
	s, err = srv.ridg.CreatePassword()
	if err != nil {
		return "", err
	}

	return *s, nil
}

func (srv *Server) createVerificationCode() (vc string, err error) {
	var s *string
	s, err = srv.vcg.CreatePassword()
	if err != nil {
		return "", err
	}

	return *s, nil
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

func (srv *Server) isEmailOfUserValid(email string) (isEmailValid bool) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	return true
}

func isPasswordAllowed(password string) (err error) {
	_, err = bpp.IsPasswordAllowed(password)
	if err != nil {
		return fmt.Errorf(c.ErrFBadPassword, err.Error())
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

func (srv *Server) sendVerificationCodeForReg(email string, code string) (err error) {
	var subject = fmt.Sprintf(srv.settings.MessageSettings.SubjectTemplate, srv.settings.SystemSettings.SiteName)

	var msg = fmt.Sprintf(srv.settings.MessageSettings.BodyTemplateForReg, srv.settings.SystemSettings.SiteName, code)

	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}

	var result = new(sm.SendMessageResult)
	err = srv.smtpServiceClient.MakeRequest(context.Background(), result, sc.FuncSendMessage, params)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) sendVerificationCodeForLogIn(email string, code string) (err error) {
	var subject = fmt.Sprintf(srv.settings.MessageSettings.SubjectTemplate, srv.settings.SystemSettings.SiteName)

	var msg = fmt.Sprintf(srv.settings.MessageSettings.BodyTemplateForLogIn, code)

	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}

	var result = new(sm.SendMessageResult)
	err = srv.smtpServiceClient.MakeRequest(context.Background(), result, sc.FuncSendMessage, params)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) sendVerificationCodeForEmailChange(email string, code string) (err error) {
	var subject = fmt.Sprintf(srv.settings.MessageSettings.SubjectTemplate, srv.settings.SystemSettings.SiteName)

	var msg = fmt.Sprintf(srv.settings.MessageSettings.BodyTemplateForEmailChange, code)

	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}

	var result = new(sm.SendMessageResult)
	err = srv.smtpServiceClient.MakeRequest(context.Background(), result, sc.FuncSendMessage, params)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) sendVerificationCodeForPwdChange(email string, code string) (err error) {
	var subject = fmt.Sprintf(srv.settings.MessageSettings.SubjectTemplate, srv.settings.SystemSettings.SiteName)

	var msg = fmt.Sprintf(srv.settings.MessageSettings.BodyTemplateForPwdChange, code)

	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}

	var result = new(sm.SendMessageResult)
	err = srv.smtpServiceClient.MakeRequest(context.Background(), result, sc.FuncSendMessage, params)
	if err != nil {
		return err
	}

	return nil
}
