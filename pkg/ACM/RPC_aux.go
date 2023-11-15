package acm

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"net/mail"
	"strings"
	"time"

	js "github.com/osamingo/jsonrpc/v2"
	bpp "github.com/vault-thirteen/BytePackedPassword"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	rc "github.com/vault-thirteen/SimpleBB/pkg/RCS/client"
	rm "github.com/vault-thirteen/SimpleBB/pkg/RCS/models"
	sc "github.com/vault-thirteen/SimpleBB/pkg/SMTP/client"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SMTP/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	n "github.com/vault-thirteen/auxie/number"
)

// Auxiliary functions used in RPC functions.

const (
	EmailAndVerificationCodeSeparator = ":"
)

const (
	ErrDecodeEmailWithVerificationCode = "e-mail with verification code decoding error"
	ErrCaptchaAnswerIsNotSet           = "answer is not set"
	ErrFPasswordIsNotAllowed           = "password is not allowed: %s" // Template.
	ErrUnsupportedRCSMode              = "unsupported RCS mode"
	ErrFakeJWToken                     = "fake JWT token"
	ErrFDBError                        = "DB error: %s"
	ErrFMakeVerificationLink           = "makeVerificationLink: %v. E-mail: %s, VC: %s."
)

// databaseError checks the database error and returns the JSON RPC error.
func (srv *Server) databaseError(err error) (jerr *js.Error) {
	srv.checkForNetworkError(err)
	log.Println(fmt.Sprintf(ErrFDBError, err.Error()))
	return &js.Error{Code: c.RpcErrorCode_DatabaseError, Message: c.RpcErrorMsg_DatabaseError}
}

// checkForNetworkError informs the system about database network errors.
func (srv *Server) checkForNetworkError(err error) {
	var nerr net.Error
	isNetworkError := errors.As(err, &nerr)

	if isNetworkError {
		srv.dbErrors <- nerr
	}
}

// mustBeAuthUserIPA ensures that user's IP address is set. If it is not set,
// an error is returned and the caller of this function must stop and return
// this error.
func (srv *Server) mustBeAuthUserIPA(auth *cm.Auth) (jerr *js.Error) {
	if auth == nil {
		// Report the incident.
		if srv.settings.SystemSettings.IsTableOfIncidentsUsed {
			err := srv.saveIncidentWithoutUserIPAM(IncidentType_IllegalAccessAttempt, "")
			if err != nil {
				return srv.databaseError(err)
			}
		}

		return &js.Error{Code: c.RpcErrorCode_MalformedRequest, Message: c.RpcErrorMsg_MalformedRequest}
	}

	if len(auth.UserIPA) == 0 {
		// Report the incident.
		if srv.settings.SystemSettings.IsTableOfIncidentsUsed {
			err := srv.saveIncidentWithoutUserIPAM(IncidentType_IllegalAccessAttempt, "")
			if err != nil {
				return srv.databaseError(err)
			}
		}

		return &js.Error{Code: c.RpcErrorCode_MalformedRequest, Message: c.RpcErrorMsg_MalformedRequest}
	}

	auth.UserIPAB = net.ParseIP(auth.UserIPA)
	if auth.UserIPAB == nil {
		return &js.Error{Code: c.RpcErrorCode_IPAddressError, Message: fmt.Sprintf(c.RpcErrorMsgF_IPAddressError, c.ErrNetParseIP)}
	}

	return nil
}

// mustBeNoAuthToken ensures that an authorization token is not present. If the
// token is present, an error is returned and the caller of this function must
// stop and return this error.
func (srv *Server) mustBeNoAuthToken(auth *cm.Auth) (jerr *js.Error) {
	jerr = srv.mustBeAuthUserIPA(auth)
	if jerr != nil {
		return jerr
	}

	if len(auth.Token) > 0 {
		// Report the incident.
		if srv.settings.SystemSettings.IsTableOfIncidentsUsed {
			err := srv.saveIncidentM(IncidentType_IllegalAccessAttempt, "", auth.UserIPAB)
			if err != nil {
				return srv.databaseError(err)
			}
		}

		return &js.Error{Code: RpcErrorCode_ThisActionIsNotForLoggedUsers, Message: RpcErrorMsg_ThisActionIsNotForLoggedUsers}
	}

	return nil
}

// mustBeAnAuthToken ensures that an authorization token is present and is
// valid. If the token is absent or invalid, an error is returned and the caller
// of this function must stop and return this error. User data is returned when
// token is valid.
func (srv *Server) mustBeAnAuthToken(auth *cm.Auth) (ud *am.UserData, jerr *js.Error) {
	jerr = srv.mustBeAuthUserIPA(auth)
	if jerr != nil {
		return nil, jerr
	}

	if len(auth.Token) == 0 {
		// Report the incident.
		if srv.settings.SystemSettings.IsTableOfIncidentsUsed {
			err := srv.saveIncidentM(IncidentType_IllegalAccessAttempt, "", auth.UserIPAB)
			if err != nil {
				return nil, srv.databaseError(err)
			}
		}

		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	var err error
	ud, err = srv.getUserDataByAuthToken(auth.Token, auth.UserIPAB)
	if err != nil {
		jerr = &js.Error{Code: c.RpcErrorCode_GetUserDataByAuthToken, Message: fmt.Sprintf(c.RpcErrorMsgF_GetUserDataByAuthToken, err.Error())}

		// Report the incident.
		if srv.settings.SystemSettings.IsTableOfIncidentsUsed {
			err = srv.saveIncidentM(IncidentType_FakeToken, "", auth.UserIPAB)
			if err != nil {
				return nil, srv.databaseError(err)
			}
		}

		return nil, jerr
	}

	return ud, nil
}

func (srv *Server) isUserEmailAddressFree(addr string) (isFree bool, err error) {
	var n int
	n, err = srv.countUsersWithEmailAddrM(addr)
	if err != nil {
		return false, err
	}

	if n > 0 {
		return false, nil
	}

	return true, nil
}

func (srv *Server) isUserEmailAddressValid(addr string) (ok bool) {
	_, err := mail.ParseAddress(addr)
	if err != nil {
		return false
	}

	return true
}

func (srv *Server) createVerificationCode() (code string, err error) {
	var s *string
	s, err = srv.vcg.CreatePassword()
	if err != nil {
		return "", err
	}

	return *s, nil
}

func (srv *Server) sendVerificationCodeForReg(email string, code string) (err error) {
	var subject = fmt.Sprintf(srv.settings.SmtpModuleSettings.MessageSubjectTemplate, srv.settings.SystemSettings.SiteName)

	var msg = fmt.Sprintf(srv.settings.SmtpModuleSettings.MessageBodyTemplateForReg, srv.settings.SystemSettings.SiteName, srv.makeVerificationLink(email, code), code)

	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}

	var result sm.SendMessageResult

	err = srv.smtpModuleClient.MakeRequest(context.Background(), &result, sc.FuncSendMessage, params)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) sendVerificationCodeForLogIn(email string, code string) (err error) {
	var subject = fmt.Sprintf(srv.settings.SmtpModuleSettings.MessageSubjectTemplate, srv.settings.SystemSettings.SiteName)

	var msg = fmt.Sprintf(srv.settings.SmtpModuleSettings.MessageBodyTemplateForLogIn, code)

	var params = sm.SendMessageParams{Recipient: email, Subject: subject, Message: msg}

	var result sm.SendMessageResult

	err = srv.smtpModuleClient.MakeRequest(context.Background(), &result, sc.FuncSendMessage, params)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) makeVerificationLink(email string, code string) (href string) {
	x := encodeEmailWithVerificationCode(email, code)

	xemail, xcode, err := decodeEmailWithVerificationCode(x)
	if (err != nil) || (xemail != email) || (xcode != code) {
		log.Println(fmt.Sprintf(ErrFMakeVerificationLink, err, email, code))
	}

	return "https://" + srv.settings.SystemSettings.SiteDomain + srv.settings.SystemSettings.EmailVerificationUrlPath + "/" + "?x=" + x
}

func encodeEmailWithVerificationCode(email string, code string) (x string) {
	return base64.RawURLEncoding.EncodeToString([]byte(email + EmailAndVerificationCodeSeparator + code))
}

func decodeEmailWithVerificationCode(x string) (email string, code string, err error) {
	var buf []byte
	buf, err = base64.RawURLEncoding.DecodeString(x)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(string(buf), EmailAndVerificationCodeSeparator)
	if len(parts) != 2 {
		return "", "", errors.New(ErrDecodeEmailWithVerificationCode)
	}

	return parts[0], parts[1], nil
}

func isPasswordAllowed(password string) (err error) {
	_, err = bpp.IsPasswordAllowed(password)
	if err != nil {
		return fmt.Errorf(ErrFPasswordIsNotAllowed, err.Error())
	}

	return nil
}

func (srv *Server) isUserNameFree(name string) (ok bool, err error) {
	var n int
	n, err = srv.countUsersWithNameM(name)
	if err != nil {
		return false, err
	}

	if n > 0 {
		return false, nil
	}

	return true, nil
}

func (srv *Server) canUserLogIn(email string) (ok bool, err error) {
	var n int
	n, err = srv.countUsersWithEmailAbleToLogInM(email)
	if err != nil {
		return false, err
	}

	if n == 1 {
		return true, nil
	}

	return false, nil
}

func (srv *Server) hasUserActiveSessions(email string) (has bool, err error) {
	var n int
	n, err = srv.countStartedUserSessionsByEmailM(email)
	if err != nil {
		return false, err
	}

	if n > 0 {
		return true, nil
	}

	return false, nil
}

func (srv *Server) hasUserPreSessions(email string) (has bool, err error) {
	var n int
	n, err = srv.countCreatedUserPreSessionsByEmailM(email)
	if err != nil {
		return false, err
	}

	if n > 0 {
		return true, nil
	}

	return false, nil
}

func (srv *Server) hasNoUserPreSessions(email string) (hasNo bool, err error) {
	var n int
	n, err = srv.countCreatedUserPreSessionsByEmailM(email)
	if err != nil {
		return false, err
	}

	if n == 0 {
		return true, nil
	}

	return false, nil
}

func (srv *Server) isCaptchaNeeded(email string) (isCaptchaNeeded bool, err error) {
	var lastBadLogInTime *time.Time
	lastBadLogInTime, err = srv.getUserLastBadLogInTimeByEmailM(email)
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

func (srv *Server) createCaptcha() (result *rm.CreateCaptchaResult, err error) {
	var params = rm.CreateCaptchaParams{}
	result = &rm.CreateCaptchaResult{}

	err = srv.captchaServiceClient.MakeRequest(context.Background(), result, rc.FuncCreateCaptcha, params)
	if err != nil {
		return nil, err
	}

	if result.IsImageDataReturned {
		return nil, errors.New(ErrUnsupportedRCSMode)
	}

	return result, nil
}

func (srv *Server) createLogInRequestId() (rid string, err error) {
	var s *string
	s, err = srv.ridg.CreatePassword()
	if err != nil {
		return "", err
	}

	return *s, nil
}

func (srv *Server) checkCaptcha(captchaId string, answer string) (result *rm.CheckCaptchaResult, err error) {
	var params = rm.CheckCaptchaParams{TaskId: captchaId}

	if len(answer) == 0 {
		return nil, errors.New(ErrCaptchaAnswerIsNotSet)
	}

	params.Value, err = n.ParseUint(answer)
	if err != nil {
		return nil, err
	}

	result = &rm.CheckCaptchaResult{}

	err = srv.captchaServiceClient.MakeRequest(context.Background(), result, rc.FuncCheckCaptcha, params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (srv *Server) checkPassword(userId uint, salt []byte, userKey []byte) (ok bool, err error) {
	var pwd []byte
	pwd, err = srv.getUserPasswordByIdM(userId)
	if err != nil {
		return false, err
	}

	var pwdRunes []rune
	pwdRunes, err = bpp.UnpackBytes(pwd)
	if err != nil {
		return false, err
	}

	return bpp.CheckHashKey(string(pwdRunes), salt, userKey)
}

// getUserDataByAuthToken validates the token and returns information about the
// user and its session.
// When 'userData' is set (not null), all its fields are also set (not null).
func (srv *Server) getUserDataByAuthToken(authToken string, userIPAB net.IP) (userData *am.UserData, err error) {
	var userId uint
	var sessionId uint
	userId, sessionId, err = srv.jwtkm.ValidateToken(authToken)
	if err != nil {
		return nil, err
	}

	userData = am.NewUserData()

	userData.User, err = srv.getUserByIdM(userId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if userData.User == nil {
		return nil, errors.New(ErrFakeJWToken)
	}

	if userData.User.Id != userId {
		return nil, errors.New(ErrFakeJWToken)
	}

	if !userData.User.CanLogIn {
		return nil, errors.New(ErrFakeJWToken)
	}

	// Attach special user roles from settings.
	userData.User.IsModerator = srv.isUserModerator(userData.User.Id)
	userData.User.IsAdministrator = srv.isUserAdministrator(userData.User.Id)

	userData.Session, err = srv.getSessionByUserIdM(userId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if userData.Session == nil {
		return nil, errors.New(ErrFakeJWToken)
	}

	if userData.Session.Id != sessionId {
		return nil, errors.New(ErrFakeJWToken)
	}

	if !userIPAB.Equal(userData.Session.UserIPAB) {
		return nil, errors.New(ErrFakeJWToken)
	}

	return userData, nil
}

func (srv *Server) isUserAdministrator(userId uint) bool {
	// While system has only few administrators, the simple array look-up is
	// faster than access to a map.

	for _, id := range srv.settings.UserRoleSettings.AdministratorIds {
		if id == userId {
			return true
		}
	}

	return false
}

func (srv *Server) isUserModerator(userId uint) bool {
	// While system has only few moderators, the simple array look-up is
	// faster than access to a map.

	for _, id := range srv.settings.UserRoleSettings.ModeratorIds {
		if id == userId {
			return true
		}
	}

	return false
}
