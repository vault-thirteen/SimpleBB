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
	"github.com/vault-thirteen/SimpleBB/pkg/common/dbo"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	cn "github.com/vault-thirteen/SimpleBB/pkg/common/net"
	num "github.com/vault-thirteen/auxie/number"
)

// Auxiliary functions used in RPC functions.

// logError logs error if debug mode is enabled.
func (srv *Server) logError(err error) {
	if err == nil {
		return
	}

	if srv.settings.SystemSettings.IsDebugMode {
		log.Println(err)
	}
}

// databaseError processes the database error and returns a JSON RPC error.
func (srv *Server) databaseError(err error) (re *jrm1.RpcError) {
	if err == nil {
		return nil
	}

	if c.IsNetworkError(err) {
		log.Println(fmt.Sprintf(c.ErrFDatabaseNetwork, err.Error()))
		*(srv.dbErrors) <- err
	} else {
		srv.logError(err)
	}

	return jrm1.NewRpcErrorByUser(c.RpcErrorCode_DatabaseError, c.RpcErrorMsg_DatabaseError, err)
}

// Token-related functions.

// mustBeAuthUserIPA ensures that user's IP address is set. If it is not set,
// an error is returned and the caller of this function must stop and return
// this error.
func (srv *Server) mustBeAuthUserIPA(auth *cmr.Auth) (re *jrm1.RpcError) {
	if auth == nil {
		srv.incidentManager.ReportIncident(am.IncidentType_IllegalAccessAttempt, "", nil)
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_MalformedRequest, c.RpcErrorMsg_MalformedRequest, nil)
	}

	if len(auth.UserIPA) == 0 {
		// Report the incident.
		srv.incidentManager.ReportIncident(am.IncidentType_IllegalAccessAttempt, "", nil)
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_MalformedRequest, c.RpcErrorMsg_MalformedRequest, nil)
	}

	var err error
	auth.UserIPAB, err = cn.ParseIPA(auth.UserIPA)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(c.RpcErrorCode_IPAddressError, fmt.Sprintf(c.RpcErrorMsgF_IPAddressError, err.Error()), nil)
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
		return jrm1.NewRpcErrorByUser(RpcErrorCode_ThisActionIsNotForLoggedUsers, RpcErrorMsg_ThisActionIsNotForLoggedUsers, nil)
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
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_InsufficientPermission, c.RpcErrorMsg_InsufficientPermission, nil)
	}

	ud, re = srv.getUserDataByAuthToken(auth.Token, auth.UserIPAB)
	if re != nil {
		srv.incidentManager.ReportIncident(am.IncidentType_FakeToken, "", auth.UserIPAB)
		return nil, re
	}

	return ud, nil
}

// getUserDataByAuthToken validates the token and returns information about the
// When 'userData' is set (not null), all its fields are also set (not null).
func (srv *Server) getUserDataByAuthToken(authToken string, userIPAB net.IP) (userData *am.UserData, re *jrm1.RpcError) {
	userId, sessionId, err := srv.jwtkm.ValidateToken(authToken)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_InvalidToken, c.RpcErrorMsg_InvalidToken, err)
	}

	userData = am.NewUserData()

	userData.User, err = srv.dbo.GetUserById(userId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if userData.User == nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_InvalidToken, c.RpcErrorMsg_InvalidToken, nil)
	}

	if userData.User.Id != userId {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_InvalidToken, c.RpcErrorMsg_InvalidToken, nil)
	}

	if !userData.User.CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_InvalidToken, c.RpcErrorMsg_InvalidToken, nil)
	}

	// Attach special user roles from settings.
	userData.User.IsModerator = srv.isUserModerator(userData.User.Id)
	userData.User.IsAdministrator = srv.isUserAdministrator(userData.User.Id)

	userData.Session, err = srv.dbo.GetSessionByUserId(userId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if userData.Session == nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_InvalidToken, c.RpcErrorMsg_InvalidToken, nil)
	}

	if userData.Session.Id != sessionId {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_InvalidToken, c.RpcErrorMsg_InvalidToken, nil)
	}

	if !userIPAB.Equal(userData.Session.UserIPAB) {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_InvalidToken, c.RpcErrorMsg_InvalidToken, nil)
	}

	return userData, nil
}

// Other functions.

func (srv *Server) checkCaptcha(captchaId string, answer string) (result *rm.CheckCaptchaResult, re *jrm1.RpcError) {
	var params = rm.CheckCaptchaParams{TaskId: captchaId}
	var err error

	if len(answer) == 0 {
		err = errors.New(c.ErrCaptchaAnswerIsNotSet)
	}
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RCS_CheckCaptcha, fmt.Sprintf(RpcErrorMsgF_RCS_CheckCaptcha, err.Error()), nil)
	}

	params.Value, err = num.ParseUint(answer)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RCS_CheckCaptcha, fmt.Sprintf(RpcErrorMsgF_RCS_CheckCaptcha, err.Error()), nil)
	}

	result = new(rm.CheckCaptchaResult)
	re, err = srv.rcsServiceClient.MakeRequest(context.Background(), rc.FuncCheckCaptcha, params, result)
	if (err == nil) && (re != nil) {
		err = re.AsError()
	}
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RCS_CheckCaptcha, fmt.Sprintf(RpcErrorMsgF_RCS_CheckCaptcha, err.Error()), nil)
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
		return false, jrm1.NewRpcErrorByUser(RpcErrorCode_PasswordCheckError, fmt.Sprintf(RpcErrorMsgF_PasswordCheckError, err.Error()), nil)
	}

	ok, err = bpp.CheckHashKey(string(pwdRunes), salt, userKey)
	if err != nil {
		srv.logError(err)
		return false, jrm1.NewRpcErrorByUser(RpcErrorCode_PasswordCheckError, fmt.Sprintf(RpcErrorMsgF_PasswordCheckError, err.Error()), nil)
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
	var re *jrm1.RpcError
	re, err = srv.rcsServiceClient.MakeRequest(context.Background(), rc.FuncCreateCaptcha, params, result)
	if err != nil {
		return nil, err
	}
	if re != nil {
		return nil, re.AsError()
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
	var re *jrm1.RpcError
	re, err = srv.smtpServiceClient.MakeRequest(context.Background(), sc.FuncSendMessage, params, result)
	if err != nil {
		return err
	}
	if re != nil {
		return re.AsError()
	}

	return nil
}

func (srv *Server) sendVerificationCodeForLogIn(email string, code string) (err error) {
	var subject = fmt.Sprintf(srv.settings.MessageSettings.SubjectTemplate, srv.settings.SystemSettings.SiteName)

	var msg = fmt.Sprintf(srv.settings.MessageSettings.BodyTemplateForLogIn, code)

	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}

	var result = new(sm.SendMessageResult)
	var re *jrm1.RpcError
	re, err = srv.smtpServiceClient.MakeRequest(context.Background(), sc.FuncSendMessage, params, result)
	if err != nil {
		return err
	}
	if re != nil {
		return re.AsError()
	}

	return nil
}

func (srv *Server) sendVerificationCodeForEmailChange(email string, code string) (err error) {
	var subject = fmt.Sprintf(srv.settings.MessageSettings.SubjectTemplate, srv.settings.SystemSettings.SiteName)

	var msg = fmt.Sprintf(srv.settings.MessageSettings.BodyTemplateForEmailChange, code)

	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}

	var result = new(sm.SendMessageResult)
	var re *jrm1.RpcError
	re, err = srv.smtpServiceClient.MakeRequest(context.Background(), sc.FuncSendMessage, params, result)
	if err != nil {
		return err
	}
	if re != nil {
		return re.AsError()
	}

	return nil
}

func (srv *Server) sendVerificationCodeForPwdChange(email string, code string) (err error) {
	var subject = fmt.Sprintf(srv.settings.MessageSettings.SubjectTemplate, srv.settings.SystemSettings.SiteName)

	var msg = fmt.Sprintf(srv.settings.MessageSettings.BodyTemplateForPwdChange, code)

	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}

	var result = new(sm.SendMessageResult)
	var re *jrm1.RpcError
	re, err = srv.smtpServiceClient.MakeRequest(context.Background(), sc.FuncSendMessage, params, result)
	if err != nil {
		return err
	}
	if re != nil {
		return re.AsError()
	}

	return nil
}
