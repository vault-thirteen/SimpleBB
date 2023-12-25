package acm

import (
	"database/sql"
	"fmt"
	"time"

	js "github.com/osamingo/jsonrpc/v2"
	bpp "github.com/vault-thirteen/BytePackedPassword"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	rm "github.com/vault-thirteen/SimpleBB/pkg/RCS/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
)

// RPC functions.

// User registration.

func (srv *Server) registerUser(p *am.RegisterUserParams) (result *am.RegisterUserResult, jerr *js.Error) {
	jerr = srv.mustBeNoAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	switch p.StepN {
	case 1:
		return srv.registerUserStep1(p)
	case 2:
		return srv.registerUserStep2(p)
	case 3:
		return srv.registerUserStep3(p)
	default:
		// Step is not supported.
		return nil, &js.Error{Code: RpcErrorCode_StepIsUnknown, Message: RpcErrorMsg_StepIsUnknown}
	}
}

func (srv *Server) registerUserStep1(p *am.RegisterUserParams) (result *am.RegisterUserResult, jerr *js.Error) {
	// Is e-mail address valid ?
	if !srv.isEmailOfUserValid(p.Email) {
		return nil, &js.Error{Code: RpcErrorCode_EmailAddressIsNotValid, Message: RpcErrorMsg_EmailAddressIsNotValid}
	}

	// Check that client's e-mail address is not used.
	usersCount, err := srv.countUsersWithEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	isAddrFree := usersCount == 0
	if !isAddrFree {
		return nil, &js.Error{Code: RpcErrorCode_EmailAddressIsUsed, Message: RpcErrorMsg_EmailAddressIsUsed}
	}

	// Add the client to a list of pre-registered users.
	// Database should contain restrictions against the following duplicate
	// properties of a user and a pre-user: name, email.
	err = srv.dbo.InsertPreRegisteredUser(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Create a verification code.
	var verificationCode string
	verificationCode, err = srv.createVerificationCode()
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_VerificationCodeGenerator, Message: RpcErrorMsg_VerificationCodeGenerator}
	}

	// Attach the verification code to the user.
	err = srv.dbo.AttachVerificationCodeToPreRegUser(p.Email, verificationCode)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Send an e-mail message with a verification code.
	err = srv.sendVerificationCodeForReg(p.Email, verificationCode)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_SmtpModule, Message: fmt.Sprintf(RpcErrorMsgF_SmtpModule, err.Error())}
	}

	// Confirm the email message send.
	err = srv.dbo.SetPreRegUserEmailSendStatus(true, p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return &am.RegisterUserResult{NextStep: 2}, nil
}

func (srv *Server) registerUserStep2(p *am.RegisterUserParams) (result *am.RegisterUserResult, jerr *js.Error) {
	// Is e-mail address valid ?
	if !srv.isEmailOfUserValid(p.Email) {
		return nil, &js.Error{Code: RpcErrorCode_EmailAddressIsNotValid, Message: RpcErrorMsg_EmailAddressIsNotValid}
	}

	// Is verification code set ?
	if len(p.VerificationCode) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_VerificationCodeIsNotSet, Message: RpcErrorMsg_VerificationCodeIsNotSet}
	}

	// Check the verification code.
	ok, err := srv.dbo.CheckVerificationCodeForPreReg(p.Email, p.VerificationCode)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if !ok {
		// Verification code can not be guessed.
		srv.incidentManager.ReportIncident(am.IncidentType_VerificationCodeMismatch, p.Email, p.Auth.UserIPAB)

		// Delete the pre-registered user.
		err = srv.dbo.DeletePreRegUserIfNotApprovedByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_VerificationCodeIsWrong, Message: RpcErrorMsg_VerificationCodeIsWrong}
	}

	// Mark the client as verified by e-mail address.
	err = srv.dbo.ApproveUserByEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return &am.RegisterUserResult{NextStep: 3}, nil
}

func (srv *Server) registerUserStep3(p *am.RegisterUserParams) (result *am.RegisterUserResult, jerr *js.Error) {
	// Is e-mail address valid ?
	if !srv.isEmailOfUserValid(p.Email) {
		return nil, &js.Error{Code: RpcErrorCode_EmailAddressIsNotValid, Message: RpcErrorMsg_EmailAddressIsNotValid}
	}

	// Is verification code set ?
	if len(p.VerificationCode) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_VerificationCodeIsNotSet, Message: RpcErrorMsg_VerificationCodeIsNotSet}
	}

	// Is name set ?
	if len(p.Name) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_NameIsNotSet, Message: RpcErrorMsg_NameIsNotSet}
	}

	// Is password set ?
	if len(p.Password) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_PasswordIsNotSet, Message: RpcErrorMsg_PasswordIsNotSet}
	}

	// Check name.
	if len([]byte(p.Name)) > int(srv.settings.SystemSettings.UserNameMaxLenInBytes) {
		return nil, &js.Error{Code: RpcErrorCode_NameIsTooLong, Message: RpcErrorMsg_NameIsTooLong}
	}

	usersCount, err := srv.countUsersWithName(p.Name)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	isUserNameFree := usersCount == 0
	if !isUserNameFree {
		return nil, &js.Error{Code: RpcErrorCode_NameIsUsed, Message: RpcErrorMsg_NameIsUsed}
	}

	// Check password.
	err = isPasswordAllowed(p.Password)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_PasswordError, Message: fmt.Sprintf(RpcErrorMsgF_PasswordError, err.Error())}
	}

	var pwdBytes []byte
	pwdBytes, err = bpp.PackSymbols([]rune(p.Password))
	if err != nil {
		// This error is very unlikely to happen.
		// So if it occurs, then it is an anomaly.
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_BPP_PackSymbols, Message: fmt.Sprintf(RpcErrorMsgF_BPP_PackSymbols, err.Error())}
	}

	if len(pwdBytes) > int(srv.settings.SystemSettings.UserPasswordMaxLenInBytes) {
		return nil, &js.Error{Code: RpcErrorCode_PasswordIsTooLong, Message: RpcErrorMsg_PasswordIsTooLong}
	}

	// Check the verification code.
	var ok bool
	ok, err = srv.dbo.CheckVerificationCodeForPreReg(p.Email, p.VerificationCode)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if !ok {
		srv.incidentManager.ReportIncident(am.IncidentType_VerificationCodeMismatch, p.Email, p.Auth.UserIPAB)
		return nil, &js.Error{Code: RpcErrorCode_VerificationCodeIsWrong, Message: RpcErrorMsg_VerificationCodeIsWrong}
	}

	// Set user's name and password.
	// Database should contain restrictions against the following duplicate
	// properties of a user and a pre-user: name, email.
	err = srv.dbo.SetPreRegUserData(p.Email, p.VerificationCode, p.Name, pwdBytes)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if srv.settings.SystemSettings.IsAdminApprovalRequired {
		return &am.RegisterUserResult{NextStep: 4}, nil
	}

	err = srv.dbo.ApprovePreRegUser(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.dbo.RegisterPreRegUser(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return &am.RegisterUserResult{NextStep: 0}, nil
}

func (srv *Server) approveAndRegisterUser(p *am.ApproveAndRegisterUserParams) (result *am.ApproveAndRegisterUserResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData *am.UserData
	thisUserData, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !thisUserData.User.IsAdministrator {
		srv.incidentManager.ReportIncident(am.IncidentType_IllegalAccessAttempt, p.Email, p.Auth.UserIPAB)
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Is e-mail address valid ?
	if !srv.isEmailOfUserValid(p.Email) {
		return nil, &js.Error{Code: RpcErrorCode_EmailAddressIsNotValid, Message: RpcErrorMsg_EmailAddressIsNotValid}
	}

	err := srv.dbo.ApprovePreRegUser(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.dbo.RegisterPreRegUser(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return &am.ApproveAndRegisterUserResult{OK: true}, nil
}

// Logging in and out.

func (srv *Server) logUserIn(p *am.LogUserInParams) (result *am.LogUserInResult, jerr *js.Error) {
	jerr = srv.mustBeNoAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	switch p.StepN {
	case 1:
		return srv.logUserInStep1(p)
	case 2:
		return srv.logUserInStep2(p)
	case 3:
		return srv.logUserInStep3(p)
	default:
		// Step is not supported.
		return nil, &js.Error{Code: RpcErrorCode_StepIsUnknown, Message: RpcErrorMsg_StepIsUnknown}
	}
}

func (srv *Server) logUserInStep1(p *am.LogUserInParams) (result *am.LogUserInResult, jerr *js.Error) {
	// Is e-mail address valid ?
	if !srv.isEmailOfUserValid(p.Email) {
		return nil, &js.Error{Code: RpcErrorCode_EmailAddressIsNotValid, Message: RpcErrorMsg_EmailAddressIsNotValid}
	}

	usersCount, err := srv.countUsersWithEmailAbleToLogIn(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	canUserLogIn := usersCount == 1
	if !canUserLogIn {
		return nil, &js.Error{Code: RpcErrorCode_UserCanNotLogIn, Message: RpcErrorMsg_UserCanNotLogIn}
	}

	err = srv.dbo.DeleteAbandonedPreSessions()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var sessionsCount int
	sessionsCount, err = srv.countSessionsByUserEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	areUserSessionsPresent := sessionsCount > 0
	if areUserSessionsPresent {
		srv.incidentManager.ReportIncident(am.IncidentType_DoubleLogInAttempt, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_UserIsAlreadyLoggedIn, Message: RpcErrorMsg_UserIsAlreadyLoggedIn}
	}

	var preSessionsCount int
	preSessionsCount, err = srv.countPreSessionsByUserEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	areUserPreSessionsPresent := preSessionsCount > 0
	if areUserPreSessionsPresent {
		srv.incidentManager.ReportIncident(am.IncidentType_PreSessionHacking, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_UserHasAlreadyStartedToLogIn, Message: RpcErrorMsg_UserHasAlreadyStartedToLogIn}
	}

	var userId uint
	userId, err = srv.dbo.GetUserIdByEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var requestId string
	requestId, err = srv.createRequestIdForLogIn()
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_RequestIdGenerator, Message: RpcErrorMsg_RequestIdGenerator}
	}

	var pwdSalt []byte
	pwdSalt, err = bpp.GenerateRandomSalt()
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_BPP_GenerateRandomSalt, Message: RpcErrorMsg_BPP_GenerateRandomSalt}
	}

	// Create a captcha if needed.
	var isCaptchaNeeded bool
	isCaptchaNeeded, err = srv.isCaptchaNeededForLogIn(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var captchaData *rm.CreateCaptchaResult
	var captchaId sql.NullString
	if isCaptchaNeeded {
		captchaData, err = srv.createCaptcha()
		if err != nil {
			srv.logError(err)
			return nil, &js.Error{Code: RpcErrorCode_RCS_CreateCaptcha, Message: fmt.Sprintf(RpcErrorMsgF_RCS_CreateCaptcha, err.Error())}
		}

		captchaId.Valid = true // => SQL non-NULL value.
		captchaId.String = captchaData.TaskId
	} else {
		captchaId.Valid = false // => SQL NULL value.
	}

	err = srv.dbo.CreatePreSession(userId, requestId, p.Auth.UserIPAB, pwdSalt, isCaptchaNeeded, captchaId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &am.LogUserInResult{NextStep: 2, RequestId: requestId, AuthDataBytes: pwdSalt}

	if isCaptchaNeeded {
		result.IsCaptchaNeeded = true
		result.CaptchaId = captchaId.String
	} else {
		result.IsCaptchaNeeded = false
	}

	return result, nil
}

func (srv *Server) logUserInStep2(p *am.LogUserInParams) (result *am.LogUserInResult, jerr *js.Error) {
	// Is e-mail address valid ?
	if !srv.isEmailOfUserValid(p.Email) {
		return nil, &js.Error{Code: RpcErrorCode_EmailAddressIsNotValid, Message: RpcErrorMsg_EmailAddressIsNotValid}
	}

	if len(p.RequestId) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_RequestIdIsNotSet, Message: RpcErrorMsg_RequestIdIsNotSet}
	}

	if len(p.AuthChallengeResponse) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_AuthChallengeResponseIsNotSet, Message: RpcErrorMsg_AuthChallengeResponseIsNotSet}
	}

	usersCount, err := srv.countUsersWithEmailAbleToLogIn(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	canUserLogIn := usersCount == 1
	if !canUserLogIn {
		return nil, &js.Error{Code: RpcErrorCode_UserCanNotLogIn, Message: RpcErrorMsg_UserCanNotLogIn}
	}

	err = srv.dbo.DeleteAbandonedPreSessions()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var sessionsCount int
	sessionsCount, err = srv.countSessionsByUserEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	areUserSessionsPresent := sessionsCount > 0
	if areUserSessionsPresent {
		srv.incidentManager.ReportIncident(am.IncidentType_DoubleLogInAttempt, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_UserIsAlreadyLoggedIn, Message: RpcErrorMsg_UserIsAlreadyLoggedIn}
	}

	var preSessionsCount int
	preSessionsCount, err = srv.countPreSessionsByUserEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if preSessionsCount != 1 {
		srv.incidentManager.ReportIncident(am.IncidentType_PreSessionHacking, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_UserHasNotStartedToLogIn, Message: RpcErrorMsg_UserHasNotStartedToLogIn}
	}

	var userId uint
	userId, err = srv.dbo.GetUserIdByEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Get the request to log in.
	var preSession *am.PreSession
	preSession, err = srv.dbo.GetPreSessionByRequestId(p.RequestId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if preSession == nil {
		// Request Id code can not be guessed.
		srv.incidentManager.ReportIncident(am.IncidentType_PreSessionHacking, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_UserPreSessionIsNotFound, Message: RpcErrorMsg_UserPreSessionIsNotFound}
	}

	// Check the Request ID (indirectly).
	if preSession.UserId != userId {
		srv.incidentManager.ReportIncident(am.IncidentType_PreSessionHacking, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_UserPreSessionIsNotFound, Message: RpcErrorMsg_UserPreSessionIsNotFound}
	}

	// Check the IP address.
	if !p.Auth.UserIPAB.Equal(preSession.UserIPAB) {
		srv.incidentManager.ReportIncident(am.IncidentType_PreSessionHacking, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_UserPreSessionIsNotFound, Message: RpcErrorMsg_UserPreSessionIsNotFound}
	}

	// Check the captcha answer if needed.
	var ccr *rm.CheckCaptchaResult
	if preSession.IsCaptchaRequired {
		if len(p.CaptchaAnswer) == 0 {
			return nil, &js.Error{Code: RpcErrorCode_CaptchaAnswerIsNotSet, Message: RpcErrorMsg_CaptchaAnswerIsNotSet}
		}

		// Get the correct captcha answer by its ID.
		ccr, jerr = srv.checkCaptcha(preSession.CaptchaId.String, p.CaptchaAnswer)
		if jerr != nil {
			return nil, jerr
		}

		if !ccr.IsSuccess {
			srv.incidentManager.ReportIncident(am.IncidentType_CaptchaAnswerMismatch, p.Email, p.Auth.UserIPAB)

			err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
			if err != nil {
				return nil, srv.databaseError(err)
			}

			// When captcha guess is wrong, we delete the preliminary session
			// to start the process from the first step.
			err = srv.dbo.DeletePreSessionByRequestId(preSession.RequestId)
			if err != nil {
				return nil, srv.databaseError(err)
			}

			return nil, &js.Error{Code: RpcErrorCode_CaptchaAnswerIsWrong, Message: RpcErrorMsg_CaptchaAnswerIsWrong}
		}
	}

	// Check the password.
	var ok bool
	ok, jerr = srv.checkPassword(userId, preSession.AuthDataBytes, p.AuthChallengeResponse)
	if jerr != nil {
		return nil, jerr
	}

	if !ok {
		srv.incidentManager.ReportIncident(am.IncidentType_PasswordMismatch, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// When password is wrong, we delete the preliminary session to start
		// the process from the first step.
		err = srv.dbo.DeletePreSessionByRequestId(preSession.RequestId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_PasswordIsWrong, Message: RpcErrorMsg_PasswordIsWrong}
	}

	// Set flags.
	//
	// Captcha flag is set to true when either it was not required or it was
	// required and the answer was correct.
	err = srv.dbo.SetPreSessionCaptchaFlags(userId, preSession.RequestId, true)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.dbo.SetPreSessionPasswordFlag(userId, preSession.RequestId, true)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Create a new Request ID for the next step.
	var step3requestId string
	step3requestId, err = srv.createRequestIdForLogIn()
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_RequestIdGenerator, Message: RpcErrorMsg_RequestIdGenerator}
	}

	err = srv.dbo.UpdatePreSessionRequestId(userId, preSession.RequestId, step3requestId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Verification by E-mail.
	var verificationCode string
	verificationCode, err = srv.createVerificationCode()
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_VerificationCodeGenerator, Message: RpcErrorMsg_VerificationCodeGenerator}
	}

	err = srv.dbo.AttachVerificationCodeToPreSession(userId, step3requestId, verificationCode)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.sendVerificationCodeForLogIn(p.Email, verificationCode)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_SmtpModule, Message: fmt.Sprintf(RpcErrorMsgF_SmtpModule, err.Error())}
	}

	err = srv.dbo.SetPreSessionEmailSendStatus(userId, step3requestId, true)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &am.LogUserInResult{NextStep: 3, RequestId: step3requestId}

	return result, nil
}

func (srv *Server) logUserInStep3(p *am.LogUserInParams) (result *am.LogUserInResult, jerr *js.Error) {
	// Is e-mail address valid ?
	if !srv.isEmailOfUserValid(p.Email) {
		return nil, &js.Error{Code: RpcErrorCode_EmailAddressIsNotValid, Message: RpcErrorMsg_EmailAddressIsNotValid}
	}

	if len(p.RequestId) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_RequestIdIsNotSet, Message: RpcErrorMsg_RequestIdIsNotSet}
	}

	if len(p.VerificationCode) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_VerificationCodeIsNotSet, Message: RpcErrorMsg_VerificationCodeIsNotSet}
	}

	usersCount, err := srv.countUsersWithEmailAbleToLogIn(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	canUserLogIn := usersCount == 1
	if !canUserLogIn {
		return nil, &js.Error{Code: RpcErrorCode_UserCanNotLogIn, Message: RpcErrorMsg_UserCanNotLogIn}
	}

	err = srv.dbo.DeleteAbandonedPreSessions()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var sessionsCount int
	sessionsCount, err = srv.countSessionsByUserEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	areUserSessionsPresent := sessionsCount > 0
	if areUserSessionsPresent {
		srv.incidentManager.ReportIncident(am.IncidentType_DoubleLogInAttempt, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_UserIsAlreadyLoggedIn, Message: RpcErrorMsg_UserIsAlreadyLoggedIn}
	}

	var preSessionsCount int
	preSessionsCount, err = srv.countPreSessionsByUserEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if preSessionsCount != 1 {
		srv.incidentManager.ReportIncident(am.IncidentType_PreSessionHacking, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_UserHasNotStartedToLogIn, Message: RpcErrorMsg_UserHasNotStartedToLogIn}
	}

	var userId uint
	userId, err = srv.dbo.GetUserIdByEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Get the request to log in.
	var preSession *am.PreSession
	preSession, err = srv.dbo.GetPreSessionByRequestId(p.RequestId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if preSession == nil {
		// Request Id code can not be guessed.
		srv.incidentManager.ReportIncident(am.IncidentType_PreSessionHacking, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_UserPreSessionIsNotFound, Message: RpcErrorMsg_UserPreSessionIsNotFound}
	}

	// Check the Request ID (indirectly).
	if preSession.UserId != userId {
		srv.incidentManager.ReportIncident(am.IncidentType_PreSessionHacking, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_UserPreSessionIsNotFound, Message: RpcErrorMsg_UserPreSessionIsNotFound}
	}

	// Check the IP address.
	if !p.Auth.UserIPAB.Equal(preSession.UserIPAB) {
		srv.incidentManager.ReportIncident(am.IncidentType_PreSessionHacking, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_UserPreSessionIsNotFound, Message: RpcErrorMsg_UserPreSessionIsNotFound}
	}

	// Check the verification code.
	var ok bool
	ok, err = srv.dbo.CheckVerificationCodeForLogIn(p.RequestId, p.VerificationCode)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if !ok {
		// Verification code can not be guessed.
		srv.incidentManager.ReportIncident(am.IncidentType_VerificationCodeMismatch, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// Delete the pre-session on error to avoid brute force checks.
		err = srv.dbo.DeletePreSessionByRequestId(preSession.RequestId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_VerificationCodeIsWrong, Message: RpcErrorMsg_VerificationCodeIsWrong}
	}

	// Set verification flags.
	err = srv.dbo.SetPreSessionVerificationFlag(userId, preSession.RequestId, true)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Create a normal session and delete the preliminary one.
	var sessionId int64
	sessionId, err = srv.dbo.CreateSession(userId, p.Auth.UserIPAB)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.dbo.DeletePreSessionByRequestId(preSession.RequestId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var tokenString string
	tokenString, err = srv.jwtkm.MakeJWToken(userId, uint(sessionId))
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_JWTCreation, Message: fmt.Sprintf(RpcErrorMsgF_JWTCreation, err.Error())}
	}

	// Journaling.
	logEvent := &am.LogEvent{
		Type:     am.LogEventType_LogIn,
		UserId:   userId,
		Email:    p.Email,
		UserIPAB: p.Auth.UserIPAB,
	}

	err = srv.dbo.SaveLogEvent(logEvent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &am.LogUserInResult{NextStep: 0, IsWebTokenSet: true, WTS: tokenString}

	return result, nil
}

func (srv *Server) logUserOut(p *am.LogUserOutParams) (result *am.LogUserOutResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData *am.UserData
	thisUserData, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	err := srv.dbo.DeleteSession(thisUserData.Session.Id, thisUserData.User.Id, thisUserData.Session.UserIPAB)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Journaling.
	logEvent := &am.LogEvent{
		Type:     am.LogEventType_LogOut,
		UserId:   thisUserData.User.Id,
		Email:    thisUserData.User.Email,
		UserIPAB: p.Auth.UserIPAB,
	}

	err = srv.dbo.SaveLogEvent(logEvent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return &am.LogUserOutResult{OK: true}, nil
}

func (srv *Server) getListOfLoggedUsers(p *am.GetListOfLoggedUsersParams) (result *am.GetListOfLoggedUsersResult, jerr *js.Error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	_, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	result = &am.GetListOfLoggedUsersResult{}

	var err error
	result.LoggedUserIds, err = srv.dbo.GetListOfLoggedUsers()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return result, nil
}

func (srv *Server) isUserLoggedIn(p *am.IsUserLoggedInParams) (result *am.IsUserLoggedInResult, jerr *js.Error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	_, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	if p.UserId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_UserIdIsNotSet, Message: RpcErrorMsg_UserIdIsNotSet}
	}

	result = &am.IsUserLoggedInResult{
		UserId: p.UserId,
	}

	sessionsCount, err := srv.countSessionsByUserId(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if sessionsCount == 1 {
		result.IsUserLoggedIn = true
	} else {
		result.IsUserLoggedIn = false
	}

	return result, nil
}

// Various actions.

func (srv *Server) changeEmail(p *am.ChangeEmailParams) (result *am.ChangeEmailResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData *am.UserData
	thisUserData, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	switch p.StepN {
	case 1:
		return srv.changeEmailStep1(p, thisUserData)
	case 2:
		return srv.changeEmailStep2(p, thisUserData)
	default:
		// Step is not supported.
		return nil, &js.Error{Code: RpcErrorCode_StepIsUnknown, Message: RpcErrorMsg_StepIsUnknown}
	}
}

func (srv *Server) changeEmailStep1(p *am.ChangeEmailParams, ud *am.UserData) (result *am.ChangeEmailResult, jerr *js.Error) {
	var ec = &am.EmailChange{
		UserId:   ud.User.Id,
		UserIPAB: p.Auth.UserIPAB,
		NewEmail: p.NewEmail,
	}

	// Check for duplicate request.
	emailChangesCount, err := srv.countEmailChangesByUserId(ud.User.Id)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if emailChangesCount > 0 {
		srv.incidentManager.ReportIncident(am.IncidentType_EmailChangeHacking, ud.User.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.User.Id)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_UserAlreadyStartedToChangeEmail, Message: RpcErrorMsg_UserAlreadyStartedToChangeEmail}
	}

	// Is new e-mail address set ?
	if len(ec.NewEmail) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_NewEmailIsNotSet, Message: RpcErrorMsg_NewEmailIsNotSet}
	}

	// Check the new e-mail address.
	// Is e-mail address valid ?
	if !srv.isEmailOfUserValid(ec.NewEmail) {
		return nil, &js.Error{Code: RpcErrorCode_EmailAddressIsNotValid, Message: RpcErrorMsg_EmailAddressIsNotValid}
	}

	// Check that the new e-mail address is not used.
	var usersCount int
	usersCount, err = srv.countUsersWithEmail(ec.NewEmail)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	isAddrFree := usersCount == 0
	if !isAddrFree {
		return nil, &js.Error{Code: RpcErrorCode_EmailAddressIsUsed, Message: RpcErrorMsg_EmailAddressIsUsed}
	}

	// Request ID.
	ec.RequestId, err = srv.createRequestIdForPasswordChange()
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_RequestIdGenerator, Message: RpcErrorMsg_RequestIdGenerator}
	}

	// Password salt.
	ec.AuthDataBytes, err = bpp.GenerateRandomSalt()
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_BPP_GenerateRandomSalt, Message: RpcErrorMsg_BPP_GenerateRandomSalt}
	}

	// Verification code for the old e-mail address.
	ec.VerificationCodeOld, err = srv.createVerificationCode()
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_VerificationCodeGenerator, Message: RpcErrorMsg_VerificationCodeGenerator}
	}

	err = srv.sendVerificationCodeForEmailChange(ud.User.Email, ec.VerificationCodeOld)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_SmtpModule, Message: fmt.Sprintf(RpcErrorMsgF_SmtpModule, err.Error())}
	}

	// Verification code for the new e-mail address.
	ec.VerificationCodeNew, err = srv.createVerificationCode()
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_VerificationCodeGenerator, Message: RpcErrorMsg_VerificationCodeGenerator}
	}

	err = srv.sendVerificationCodeForEmailChange(ec.NewEmail, ec.VerificationCodeNew)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_SmtpModule, Message: fmt.Sprintf(RpcErrorMsgF_SmtpModule, err.Error())}
	}

	// Captcha.
	ec.IsCaptchaRequired, err = srv.isCaptchaNeededForEmailChange(ud.User.Id)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var captchaData *rm.CreateCaptchaResult
	if ec.IsCaptchaRequired {
		captchaData, err = srv.createCaptcha()
		if err != nil {
			srv.logError(err)
			return nil, &js.Error{Code: RpcErrorCode_RCS_CreateCaptcha, Message: fmt.Sprintf(RpcErrorMsgF_RCS_CreateCaptcha, err.Error())}
		}

		ec.CaptchaId.Valid = true // => SQL non-NULL value.
		ec.CaptchaId.String = captchaData.TaskId
	} else {
		ec.CaptchaId.Valid = false // => SQL NULL value.
	}

	// Save the request.
	err = srv.dbo.CreateEmailChangeRequest(ec)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Response.
	result = &am.ChangeEmailResult{
		NextStep:      2,
		RequestId:     ec.RequestId,
		AuthDataBytes: ec.AuthDataBytes,
	}

	if ec.IsCaptchaRequired {
		result.IsCaptchaNeeded = true
		result.CaptchaId = ec.CaptchaId.String
	} else {
		result.IsCaptchaNeeded = false
	}

	return result, nil
}

func (srv *Server) changeEmailStep2(p *am.ChangeEmailParams, ud *am.UserData) (result *am.ChangeEmailResult, jerr *js.Error) {
	// Check for duplicate request.
	emailChangesCount, err := srv.countEmailChangesByUserId(ud.User.Id)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if emailChangesCount != 1 {
		srv.incidentManager.ReportIncident(am.IncidentType_EmailChangeHacking, ud.User.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.User.Id)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_UserAlreadyStartedToChangeEmail, Message: RpcErrorMsg_UserAlreadyStartedToChangeEmail}
	}

	// Check input parameters.
	if len(p.RequestId) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_RequestIdIsNotSet, Message: RpcErrorMsg_RequestIdIsNotSet}
	}
	if len(p.AuthChallengeResponse) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_AuthChallengeResponseIsNotSet, Message: RpcErrorMsg_AuthChallengeResponseIsNotSet}
	}
	if len(p.VerificationCodeOld) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_VerificationCodeIsNotSet, Message: RpcErrorMsg_VerificationCodeIsNotSet}
	}
	if len(p.VerificationCodeNew) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_VerificationCodeIsNotSet, Message: RpcErrorMsg_VerificationCodeIsNotSet}
	}

	// Get the request for e-mail address change.
	var ecr *am.EmailChange
	ecr, err = srv.dbo.GetEmailChangeByRequestId(p.RequestId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if ecr == nil {
		// Request Id code can not be guessed.
		srv.incidentManager.ReportIncident(am.IncidentType_EmailChangeHacking, ud.User.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.User.Id)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_EmailChangeIsNotFound, Message: RpcErrorMsg_EmailChangeIsNotFound}
	}

	// Check the Request ID (indirectly).
	if ecr.UserId != ud.User.Id {
		srv.incidentManager.ReportIncident(am.IncidentType_EmailChangeHacking, ud.User.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.User.Id)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_EmailChangeIsNotFound, Message: RpcErrorMsg_EmailChangeIsNotFound}
	}

	// Check the IP address.
	if !p.Auth.UserIPAB.Equal(ecr.UserIPAB) {
		srv.incidentManager.ReportIncident(am.IncidentType_EmailChangeHacking, ud.User.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.User.Id)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_EmailChangeIsNotFound, Message: RpcErrorMsg_EmailChangeIsNotFound}
	}

	// Check the captcha answer if needed.
	var ccr *rm.CheckCaptchaResult
	if ecr.IsCaptchaRequired {
		if len(p.CaptchaAnswer) == 0 {
			return nil, &js.Error{Code: RpcErrorCode_CaptchaAnswerIsNotSet, Message: RpcErrorMsg_CaptchaAnswerIsNotSet}
		}

		// Get the correct captcha answer by its ID.
		ccr, jerr = srv.checkCaptcha(ecr.CaptchaId.String, p.CaptchaAnswer)
		if jerr != nil {
			return nil, jerr
		}

		if !ccr.IsSuccess {
			srv.incidentManager.ReportIncident(am.IncidentType_CaptchaAnswerMismatch, ud.User.Email, p.Auth.UserIPAB)

			err = srv.dbo.UpdateUserLastBadActionTimeById(ud.User.Id)
			if err != nil {
				return nil, srv.databaseError(err)
			}

			// When captcha guess is wrong, we delete the e-mail address change
			// request to start the process from the first step.
			err = srv.dbo.DeleteEmailChangeByRequestId(ecr.RequestId)
			if err != nil {
				return nil, srv.databaseError(err)
			}

			return nil, &js.Error{Code: RpcErrorCode_CaptchaAnswerIsWrong, Message: RpcErrorMsg_CaptchaAnswerIsWrong}
		}
	}

	// Check the password.
	var ok bool
	ok, jerr = srv.checkPassword(ud.User.Id, ecr.AuthDataBytes, p.AuthChallengeResponse)
	if jerr != nil {
		return nil, jerr
	}

	if !ok {
		srv.incidentManager.ReportIncident(am.IncidentType_PasswordMismatch, ud.User.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.User.Id)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// When password is wrong, we delete the password change request to
		// start the process from the first step.
		err = srv.dbo.DeleteEmailChangeByRequestId(ecr.RequestId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_PasswordIsWrong, Message: RpcErrorMsg_PasswordIsWrong}
	}

	// Check the verification codes.
	ok, err = srv.dbo.CheckVerificationCodesForEmailChange(p.RequestId, p.VerificationCodeOld, p.VerificationCodeNew)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if !ok {
		// Verification codes can not be guessed.
		srv.incidentManager.ReportIncident(am.IncidentType_VerificationCodeMismatch, ud.User.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.User.Id)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// Delete the e-mail address change request on error to avoid brute
		// force checks.
		err = srv.dbo.DeleteEmailChangeByRequestId(ecr.RequestId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_VerificationCodeIsWrong, Message: RpcErrorMsg_VerificationCodeIsWrong}
	}

	// Update request's flags.
	ecvf := &am.EmailChangeVerificationFlags{
		IsVerifiedByPassword: true,
		IsVerifiedByOldEmail: true,
		IsVerifiedByNewEmail: true,
	}
	if ecr.IsCaptchaRequired {
		ecvf.IsVerifiedByCaptcha.Valid = true // => SQL non-NULL value.
		ecvf.IsVerifiedByCaptcha.Bool = true
	} else {
		ecvf.IsVerifiedByCaptcha.Valid = false // => SQL NULL value.
	}

	err = srv.dbo.SetEmailChangeVFlags(ud.User.Id, ecr.RequestId, ecvf)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Change the e-mail address.
	err = srv.dbo.SetUserEmail(ud.User.Id, ud.User.Email, ecr.NewEmail)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Close the session to force re-logging.
	err = srv.dbo.DeleteSession(ud.Session.Id, ud.User.Id, p.Auth.UserIPAB)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// N.B.: We do not immediately delete the request for e-mail address change
	// here to avoid spamming with these requests. It will be deleted later by
	// the scheduler.

	// Response.
	result = &am.ChangeEmailResult{
		NextStep: 0,
	}

	return result, nil
}

func (srv *Server) changePassword(p *am.ChangePasswordParams) (result *am.ChangePasswordResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData *am.UserData
	thisUserData, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	switch p.StepN {
	case 1:
		return srv.changePasswordStep1(p, thisUserData)
	case 2:
		return srv.changePasswordStep2(p, thisUserData)
	default:
		// Step is not supported.
		return nil, &js.Error{Code: RpcErrorCode_StepIsUnknown, Message: RpcErrorMsg_StepIsUnknown}
	}
}

func (srv *Server) changePasswordStep1(p *am.ChangePasswordParams, ud *am.UserData) (result *am.ChangePasswordResult, jerr *js.Error) {
	var pc = &am.PasswordChange{
		UserId:   ud.User.Id,
		UserIPAB: p.Auth.UserIPAB,
	}

	// Check for duplicate request.
	passwordChangesCount, err := srv.countPasswordChangesByUserId(ud.User.Id)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if passwordChangesCount > 0 {
		srv.incidentManager.ReportIncident(am.IncidentType_PasswordChangeHacking, ud.User.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.User.Id)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_UserAlreadyStartedToChangePassword, Message: RpcErrorMsg_UserAlreadyStartedToChangePassword}
	}

	// Is new password set ?
	if len(p.NewPassword) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_NewPasswordIsNotSet, Message: RpcErrorMsg_NewPasswordIsNotSet}
	}

	// Check the new password.
	err = isPasswordAllowed(p.NewPassword)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_PasswordError, Message: fmt.Sprintf(RpcErrorMsgF_PasswordError, err.Error())}
	}

	pc.NewPassword, err = bpp.PackSymbols([]rune(p.NewPassword))
	if err != nil {
		// This error is very unlikely to happen.
		// So if it occurs, then it is an anomaly.
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_BPP_PackSymbols, Message: fmt.Sprintf(RpcErrorMsgF_BPP_PackSymbols, err.Error())}
	}

	if len(pc.NewPassword) > int(srv.settings.SystemSettings.UserPasswordMaxLenInBytes) {
		return nil, &js.Error{Code: RpcErrorCode_PasswordIsTooLong, Message: RpcErrorMsg_PasswordIsTooLong}
	}

	// Request ID.
	pc.RequestId, err = srv.createRequestIdForPasswordChange()
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_RequestIdGenerator, Message: RpcErrorMsg_RequestIdGenerator}
	}

	// Password salt.
	pc.AuthDataBytes, err = bpp.GenerateRandomSalt()
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_BPP_GenerateRandomSalt, Message: RpcErrorMsg_BPP_GenerateRandomSalt}
	}

	// Verification code.
	pc.VerificationCode, err = srv.createVerificationCode()
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_VerificationCodeGenerator, Message: RpcErrorMsg_VerificationCodeGenerator}
	}

	err = srv.sendVerificationCodeForPwdChange(ud.User.Email, pc.VerificationCode)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: RpcErrorCode_SmtpModule, Message: fmt.Sprintf(RpcErrorMsgF_SmtpModule, err.Error())}
	}

	// Captcha.
	pc.IsCaptchaRequired, err = srv.isCaptchaNeededForPasswordChange(ud.User.Id)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var captchaData *rm.CreateCaptchaResult
	if pc.IsCaptchaRequired {
		captchaData, err = srv.createCaptcha()
		if err != nil {
			srv.logError(err)
			return nil, &js.Error{Code: RpcErrorCode_RCS_CreateCaptcha, Message: fmt.Sprintf(RpcErrorMsgF_RCS_CreateCaptcha, err.Error())}
		}

		pc.CaptchaId.Valid = true // => SQL non-NULL value.
		pc.CaptchaId.String = captchaData.TaskId
	} else {
		pc.CaptchaId.Valid = false // => SQL NULL value.
	}

	// Save the request.
	err = srv.dbo.CreatePasswordChangeRequest(pc)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Response.
	result = &am.ChangePasswordResult{
		NextStep:      2,
		RequestId:     pc.RequestId,
		AuthDataBytes: pc.AuthDataBytes,
	}

	if pc.IsCaptchaRequired {
		result.IsCaptchaNeeded = true
		result.CaptchaId = pc.CaptchaId.String
	} else {
		result.IsCaptchaNeeded = false
	}

	return result, nil
}

func (srv *Server) changePasswordStep2(p *am.ChangePasswordParams, ud *am.UserData) (result *am.ChangePasswordResult, jerr *js.Error) {
	// Check for duplicate request.
	passwordChangesCount, err := srv.countPasswordChangesByUserId(ud.User.Id)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if passwordChangesCount != 1 {
		srv.incidentManager.ReportIncident(am.IncidentType_PasswordChangeHacking, ud.User.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.User.Id)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_UserAlreadyStartedToChangePassword, Message: RpcErrorMsg_UserAlreadyStartedToChangePassword}
	}

	// Check input parameters.
	if len(p.RequestId) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_RequestIdIsNotSet, Message: RpcErrorMsg_RequestIdIsNotSet}
	}
	if len(p.AuthChallengeResponse) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_AuthChallengeResponseIsNotSet, Message: RpcErrorMsg_AuthChallengeResponseIsNotSet}
	}
	if len(p.VerificationCode) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_VerificationCodeIsNotSet, Message: RpcErrorMsg_VerificationCodeIsNotSet}
	}

	// Get the request for password change.
	var pcr *am.PasswordChange
	pcr, err = srv.dbo.GetPasswordChangeByRequestId(p.RequestId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if pcr == nil {
		// Request Id code can not be guessed.
		srv.incidentManager.ReportIncident(am.IncidentType_PasswordChangeHacking, ud.User.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.User.Id)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_PasswordChangeIsNotFound, Message: RpcErrorMsg_PasswordChangeIsNotFound}
	}

	// Check the Request ID (indirectly).
	if pcr.UserId != ud.User.Id {
		srv.incidentManager.ReportIncident(am.IncidentType_PasswordChangeHacking, ud.User.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.User.Id)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_PasswordChangeIsNotFound, Message: RpcErrorMsg_PasswordChangeIsNotFound}
	}

	// Check the IP address.
	if !p.Auth.UserIPAB.Equal(pcr.UserIPAB) {
		srv.incidentManager.ReportIncident(am.IncidentType_PasswordChangeHacking, ud.User.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.User.Id)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_PasswordChangeIsNotFound, Message: RpcErrorMsg_PasswordChangeIsNotFound}
	}

	// Check the captcha answer if needed.
	var ccr *rm.CheckCaptchaResult
	if pcr.IsCaptchaRequired {
		if len(p.CaptchaAnswer) == 0 {
			return nil, &js.Error{Code: RpcErrorCode_CaptchaAnswerIsNotSet, Message: RpcErrorMsg_CaptchaAnswerIsNotSet}
		}

		// Get the correct captcha answer by its ID.
		ccr, jerr = srv.checkCaptcha(pcr.CaptchaId.String, p.CaptchaAnswer)
		if jerr != nil {
			return nil, jerr
		}

		if !ccr.IsSuccess {
			srv.incidentManager.ReportIncident(am.IncidentType_CaptchaAnswerMismatch, ud.User.Email, p.Auth.UserIPAB)

			err = srv.dbo.UpdateUserLastBadActionTimeById(ud.User.Id)
			if err != nil {
				return nil, srv.databaseError(err)
			}

			// When captcha guess is wrong, we delete the password change
			// request to start the process from the first step.
			err = srv.dbo.DeletePasswordChangeByRequestId(pcr.RequestId)
			if err != nil {
				return nil, srv.databaseError(err)
			}

			return nil, &js.Error{Code: RpcErrorCode_CaptchaAnswerIsWrong, Message: RpcErrorMsg_CaptchaAnswerIsWrong}
		}
	}

	// Check the password.
	var ok bool
	ok, jerr = srv.checkPassword(ud.User.Id, pcr.AuthDataBytes, p.AuthChallengeResponse)
	if jerr != nil {
		return nil, jerr
	}

	if !ok {
		srv.incidentManager.ReportIncident(am.IncidentType_PasswordMismatch, ud.User.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.User.Id)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// When password is wrong, we delete the password change request to
		// start the process from the first step.
		err = srv.dbo.DeletePasswordChangeByRequestId(pcr.RequestId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_PasswordIsWrong, Message: RpcErrorMsg_PasswordIsWrong}
	}

	// Check the verification code.
	ok, err = srv.dbo.CheckVerificationCodeForPwdChange(p.RequestId, p.VerificationCode)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if !ok {
		// Verification code can not be guessed.
		srv.incidentManager.ReportIncident(am.IncidentType_VerificationCodeMismatch, ud.User.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadActionTimeById(ud.User.Id)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// Delete the password change request on error to avoid brute force
		// checks.
		err = srv.dbo.DeletePasswordChangeByRequestId(pcr.RequestId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_VerificationCodeIsWrong, Message: RpcErrorMsg_VerificationCodeIsWrong}
	}

	// Update request's flags.
	pcvf := &am.PasswordChangeVerificationFlags{
		IsVerifiedByPassword: true,
		IsVerifiedByEmail:    true,
	}
	if pcr.IsCaptchaRequired {
		pcvf.IsVerifiedByCaptcha.Valid = true // => SQL non-NULL value.
		pcvf.IsVerifiedByCaptcha.Bool = true
	} else {
		pcvf.IsVerifiedByCaptcha.Valid = false // => SQL NULL value.
	}

	err = srv.dbo.SetPasswordChangeVFlags(ud.User.Id, pcr.RequestId, pcvf)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Change the password.
	err = srv.dbo.SetUserPassword(ud.User.Id, ud.User.Email, pcr.NewPassword)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Close the session to force re-logging.
	err = srv.dbo.DeleteSession(ud.Session.Id, ud.User.Id, p.Auth.UserIPAB)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// N.B.: We do not immediately delete the request for password change here
	// to avoid spamming with these requests. It will be deleted later by the
	// scheduler.

	// Response.
	result = &am.ChangePasswordResult{
		NextStep: 0,
	}

	return result, nil
}

// User properties.

func (srv *Server) getUserRoles(p *am.GetUserRolesParams) (result *am.GetUserRolesResult, jerr *js.Error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	_, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	if p.UserId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_UserIdIsNotSet, Message: RpcErrorMsg_UserIdIsNotSet}
	}

	result = &am.GetUserRolesResult{
		UserId: p.UserId,
	}

	var err error
	var roles *cm.UserRoles
	roles, err = srv.dbo.GetUserRolesById(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if roles == nil {
		return nil, &js.Error{Code: RpcErrorCode_UserIsNotFound, Message: RpcErrorMsg_UserIsNotFound}
	}

	roles.IsModerator = srv.isUserModerator(p.UserId)
	roles.IsAdministrator = srv.isUserAdministrator(p.UserId)

	result.UserRoles = *roles

	return result, nil
}

func (srv *Server) viewUserParameters(p *am.ViewUserParametersParams) (result *am.ViewUserParametersResult, jerr *js.Error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var thisUserData *am.UserData
	thisUserData, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !thisUserData.User.IsAdministrator {
		srv.incidentManager.ReportIncident(am.IncidentType_IllegalAccessAttempt, "", p.Auth.UserIPAB)
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	if p.UserId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_UserIdIsNotSet, Message: RpcErrorMsg_UserIdIsNotSet}
	}

	result = &am.ViewUserParametersResult{
		UserId: p.UserId,
	}

	var err error
	var userParameters *cm.UserParameters
	userParameters, err = srv.dbo.ViewUserParametersById(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if userParameters == nil {
		return nil, &js.Error{Code: RpcErrorCode_UserIsNotFound, Message: RpcErrorMsg_UserIsNotFound}
	}

	userParameters.IsModerator = srv.isUserModerator(p.UserId)
	userParameters.IsAdministrator = srv.isUserAdministrator(p.UserId)

	result.UserParameters = *userParameters

	return result, nil
}

func (srv *Server) setUserRoleAuthor(p *am.SetUserRoleAuthorParams) (result *am.SetUserRoleAuthorResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData *am.UserData
	thisUserData, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !thisUserData.User.IsAdministrator {
		srv.incidentManager.ReportIncident(am.IncidentType_IllegalAccessAttempt, "", p.Auth.UserIPAB)
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	if p.UserId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_UserIdIsNotSet, Message: RpcErrorMsg_UserIdIsNotSet}
	}

	err := srv.dbo.SetUserRoleAuthor(p.UserId, p.IsRoleEnabled)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return &am.SetUserRoleAuthorResult{OK: true}, nil
}

func (srv *Server) setUserRoleWriter(p *am.SetUserRoleWriterParams) (result *am.SetUserRoleWriterResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData *am.UserData
	thisUserData, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !thisUserData.User.IsAdministrator {
		srv.incidentManager.ReportIncident(am.IncidentType_IllegalAccessAttempt, "", p.Auth.UserIPAB)
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	if p.UserId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_UserIdIsNotSet, Message: RpcErrorMsg_UserIdIsNotSet}
	}

	err := srv.dbo.SetUserRoleWriter(p.UserId, p.IsRoleEnabled)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return &am.SetUserRoleWriterResult{OK: true}, nil
}

func (srv *Server) setUserRoleReader(p *am.SetUserRoleReaderParams) (result *am.SetUserRoleReaderResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData *am.UserData
	thisUserData, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !thisUserData.User.IsAdministrator {
		srv.incidentManager.ReportIncident(am.IncidentType_IllegalAccessAttempt, "", p.Auth.UserIPAB)
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	if p.UserId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_UserIdIsNotSet, Message: RpcErrorMsg_UserIdIsNotSet}
	}

	err := srv.dbo.SetUserRoleReader(p.UserId, p.IsRoleEnabled)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return &am.SetUserRoleReaderResult{OK: true}, nil
}

func (srv *Server) getSelfRoles(p *am.GetSelfRolesParams) (result *am.GetSelfRolesResult, jerr *js.Error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var thisUserData *am.UserData
	thisUserData, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	result = &am.GetSelfRolesResult{
		UserId: thisUserData.User.Id,
	}

	result.UserRoles = cm.UserRoles{
		IsAdministrator: thisUserData.User.IsAdministrator,
		IsModerator:     thisUserData.User.IsModerator,
		IsAuthor:        thisUserData.User.IsAuthor,
		IsWriter:        thisUserData.User.IsWriter,
		IsReader:        thisUserData.User.IsReader,
		CanLogIn:        thisUserData.User.CanLogIn,
	}

	return result, nil
}

// User banning.

func (srv *Server) banUser(p *am.BanUserParams) (result *am.BanUserResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData *am.UserData
	thisUserData, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !thisUserData.User.IsModerator {
		srv.incidentManager.ReportIncident(am.IncidentType_IllegalAccessAttempt, "", p.Auth.UserIPAB)
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	if p.UserId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_UserIdIsNotSet, Message: RpcErrorMsg_UserIdIsNotSet}
	}

	err := srv.dbo.SetUserRoleCanLogIn(p.UserId, false)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.dbo.UpdateUserBanTime(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.dbo.DeleteSessionByUserId(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return &am.BanUserResult{OK: true}, nil
}

func (srv *Server) unbanUser(p *am.UnbanUserParams) (result *am.UnbanUserResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var thisUserData *am.UserData
	thisUserData, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !thisUserData.User.IsModerator {
		srv.incidentManager.ReportIncident(am.IncidentType_IllegalAccessAttempt, "", p.Auth.UserIPAB)
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	if p.UserId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_UserIdIsNotSet, Message: RpcErrorMsg_UserIdIsNotSet}
	}

	err := srv.dbo.SetUserRoleCanLogIn(p.UserId, true)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return &am.UnbanUserResult{OK: true}, nil
}

// Other.

func (srv *Server) showDiagnosticData() (result *am.ShowDiagnosticDataResult, jerr *js.Error) {
	result = &am.ShowDiagnosticDataResult{
		TotalRequestsCount:      srv.diag.GetTotalRequestsCount(),
		SuccessfulRequestsCount: srv.diag.GetSuccessfulRequestsCount(),
	}

	return result, nil
}

func (srv *Server) doTest(_ *am.TestParams) (result *am.TestResult, jerr *js.Error) {
	result = &am.TestResult{}

	time.Sleep(time.Second * 10)

	return result, nil
}
