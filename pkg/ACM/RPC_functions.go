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
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_StepIsUnknown, Message: RpcErrorMsg_RegisterUser_StepIsUnknown}
	}
}

func (srv *Server) registerUserStep1(p *am.RegisterUserParams) (result *am.RegisterUserResult, jerr *js.Error) {
	// Is e-mail address valid ?
	if !srv.isUserEmailAddressValid(p.Email) {
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_EmailAddresIsNotValid, Message: RpcErrorMsg_RegisterUser_EmailAddresIsNotValid}
	}

	// Check that client's e-mail address is not used.
	isAddrFree, err := srv.isUserEmailAddressFree(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if !isAddrFree {
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_EmailAddresIsUsed, Message: RpcErrorMsg_RegisterUser_EmailAddresIsUsed}
	}

	// Add the client to a list of pre-registered users.
	err = srv.dbo.InsertPreRegisteredUser(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Send an e-mail message with a verification code.
	var verificationCode string
	verificationCode, err = srv.createVerificationCode()
	if err != nil {
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_VerificationCodeGenerator, Message: RpcErrorMsg_RegisterUser_VerificationCodeGenerator}
	}

	// Attach the verification code to the user.
	err = srv.dbo.AttachVerificationCodeToPreRegUser(p.Email, verificationCode)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.sendVerificationCodeForReg(p.Email, verificationCode)
	if err != nil {
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_SmtpModule, Message: fmt.Sprintf(RpcErrorMsgF_RegisterUser_SmtpModule, err.Error())}
	}

	return &am.RegisterUserResult{NextStep: 2}, nil
}

func (srv *Server) registerUserStep2(p *am.RegisterUserParams) (result *am.RegisterUserResult, jerr *js.Error) {
	// Is e-mail address valid ?
	if !srv.isUserEmailAddressValid(p.Email) {
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_EmailAddresIsNotValid, Message: RpcErrorMsg_RegisterUser_EmailAddresIsNotValid}
	}

	// Is verification code set ?
	if len(p.VerificationCode) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_VerificationCodeIsNotSet, Message: RpcErrorMsg_RegisterUser_VerificationCodeIsNotSet}
	}

	// Check the verification code.
	ok, err := srv.dbo.CheckRegVerificationCode(p.Email, p.VerificationCode)
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

		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_VerificationCodeIsNotValid, Message: RpcErrorMsg_RegisterUser_VerificationCodeIsNotValid}
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
	if !srv.isUserEmailAddressValid(p.Email) {
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_EmailAddresIsNotValid, Message: RpcErrorMsg_RegisterUser_EmailAddresIsNotValid}
	}

	// Is verification code set ?
	if len(p.VerificationCode) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_VerificationCodeIsNotSet, Message: RpcErrorMsg_RegisterUser_VerificationCodeIsNotSet}
	}

	// Is name set ?
	if len(p.Name) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_NameIsNotSet, Message: RpcErrorMsg_RegisterUser_NameIsNotSet}
	}

	// Is password set ?
	if len(p.Password) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_PasswordIsNotSet, Message: RpcErrorMsg_RegisterUser_PasswordIsNotSet}
	}

	// Check name.
	if len([]byte(p.Name)) > int(srv.settings.SystemSettings.UserNameMaxLenInBytes) {
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_UserNameIsTooLong, Message: RpcErrorMsg_RegisterUser_UserNameIsTooLong}
	}

	ok, err := srv.isUserNameFree(p.Name)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if !ok {
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_UserNameIsUsed, Message: RpcErrorMsg_RegisterUser_UserNameIsUsed}
	}

	// Check password.
	err = isPasswordAllowed(p.Password)
	if err != nil {
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_PasswordError, Message: fmt.Sprintf(RpcErrorMsgF_RegisterUser_PasswordError, err.Error())}
	}

	var pwdBytes []byte
	pwdBytes, err = bpp.PackSymbols([]rune(p.Password))
	if err != nil {
		// This error is very unlikely to happen.
		// So if it occurs, then it is an anomaly.
		return nil, &js.Error{Code: RpcErrorCode_BPP_PackSymbols, Message: fmt.Sprintf(RpcErrorMsgF_BPP_PackSymbols, err.Error())}
	}

	if len(pwdBytes) > int(srv.settings.SystemSettings.UserPasswordMaxLenInBytes) {
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_UserPasswordIsTooLong, Message: RpcErrorMsg_RegisterUser_UserPasswordIsTooLong}
	}

	// Check the verification code.
	ok, err = srv.dbo.CheckRegVerificationCode(p.Email, p.VerificationCode)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if !ok {
		srv.incidentManager.ReportIncident(am.IncidentType_VerificationCodeMismatch, p.Email, p.Auth.UserIPAB)
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_VerificationCodeIsNotValid, Message: RpcErrorMsg_RegisterUser_VerificationCodeIsNotValid}
	}

	// Set user's name and password.
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
	if !srv.isUserEmailAddressValid(p.Email) {
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_EmailAddresIsNotValid, Message: RpcErrorMsg_RegisterUser_EmailAddresIsNotValid}
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
		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_StepIsUnknown, Message: RpcErrorMsg_LogUserIn_StepIsUnknown}
	}
}

func (srv *Server) logUserInStep1(p *am.LogUserInParams) (result *am.LogUserInResult, jerr *js.Error) {
	// Is e-mail address valid ?
	if !srv.isUserEmailAddressValid(p.Email) {
		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_EmailAddresIsNotValid, Message: RpcErrorMsg_LogUserIn_EmailAddresIsNotValid}
	}

	var canUserLogIn bool
	var err error
	canUserLogIn, err = srv.canUserLogIn(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if !canUserLogIn {
		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_UserCanNotLogIn, Message: RpcErrorMsg_LogUserIn_UserCanNotLogIn}
	}

	err = srv.dbo.DeleteAbandonedPreSessions()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var areUserActiveSessionsPresent bool
	areUserActiveSessionsPresent, err = srv.hasUserActiveSessions(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if areUserActiveSessionsPresent {
		srv.incidentManager.ReportIncident(am.IncidentType_DoubleLogInAttempt, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_UserIsAlreadyLoggedIn, Message: RpcErrorMsg_LogUserIn_UserIsAlreadyLoggedIn}
	}

	var areUserActivePreSessionsPresent bool
	areUserActivePreSessionsPresent, err = srv.hasUserPreSessions(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if areUserActivePreSessionsPresent {
		srv.incidentManager.ReportIncident(am.IncidentType_DoubleLogInAttempt, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_UserHasAlreadyStartedToLogIn, Message: RpcErrorMsg_LogUserIn_UserHasAlreadyStartedToLogIn}
	}

	var userId uint
	userId, err = srv.dbo.GetUserIdByEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var requestId string
	requestId, err = srv.createLogInRequestId()
	if err != nil {
		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_RequestIdGenerator, Message: RpcErrorMsg_LogUserIn_RequestIdGenerator}
	}

	var pwdSalt []byte
	pwdSalt, err = bpp.GenerateRandomSalt()
	if err != nil {
		return nil, &js.Error{Code: RpcErrorCode_BPP_GenerateRandomSalt, Message: RpcErrorMsg_BPP_GenerateRandomSalt}
	}

	var isCaptchaNeeded bool
	isCaptchaNeeded, err = srv.isCaptchaNeeded(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Create a captcha if needed.
	var captchaData *rm.CreateCaptchaResult
	var captchaId sql.NullString
	if isCaptchaNeeded {
		captchaData, err = srv.createCaptcha()
		if err != nil {
			return nil, &js.Error{Code: RpcErrorCode_RCS_CreateCaptcha, Message: fmt.Sprintf(RpcErrorMsgF_RCS_CreateCaptcha, err.Error())}
		}

		captchaId.Valid = true // => SQL non-NULL value.
		captchaId.String = captchaData.TaskId
	} else {
		captchaId.Valid = false // => SQL NULL value.
	}

	err = srv.dbo.CreatePreliminarySession(userId, requestId, p.Auth.UserIPAB, pwdSalt, isCaptchaNeeded, captchaId)
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
	if !srv.isUserEmailAddressValid(p.Email) {
		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_EmailAddresIsNotValid, Message: RpcErrorMsg_LogUserIn_EmailAddresIsNotValid}
	}

	if len(p.RequestId) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_RequestIdIsNotSet, Message: RpcErrorMsg_LogUserIn_RequestIdIsNotSet}
	}

	if len(p.AuthChallengeResponse) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_AuthChallengeResponseIsNotSet, Message: RpcErrorMsg_LogUserIn_AuthChallengeResponseIsNotSet}
	}

	var canUserLogIn bool
	var err error
	canUserLogIn, err = srv.canUserLogIn(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if !canUserLogIn {
		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_UserCanNotLogIn, Message: RpcErrorMsg_LogUserIn_UserCanNotLogIn}
	}

	err = srv.dbo.DeleteAbandonedPreSessions()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var areUserActiveSessionsPresent bool
	areUserActiveSessionsPresent, err = srv.hasUserActiveSessions(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if areUserActiveSessionsPresent {
		srv.incidentManager.ReportIncident(am.IncidentType_DoubleLogInAttempt, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_UserIsAlreadyLoggedIn, Message: RpcErrorMsg_LogUserIn_UserIsAlreadyLoggedIn}
	}

	var areNoUserActivePreSessionsPresent bool
	areNoUserActivePreSessionsPresent, err = srv.hasNoUserPreSessions(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if areNoUserActivePreSessionsPresent {
		srv.incidentManager.ReportIncident(am.IncidentType_PreSessionHacking, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_UserHasNotStartedToLogIn, Message: RpcErrorMsg_LogUserIn_UserHasNotStartedToLogIn}
	}

	var userId uint
	userId, err = srv.dbo.GetUserIdByEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Check the Request ID and IP address.
	var preSession *am.PreSession
	preSession, err = srv.dbo.GetUserPreSessionByRequestId(p.RequestId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if preSession.UserId != userId {
		srv.incidentManager.ReportIncident(am.IncidentType_PreSessionHacking, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_UserPreSessionIsNotFound, Message: RpcErrorMsg_LogUserIn_UserPreSessionIsNotFound}
	}

	if !p.Auth.UserIPAB.Equal(preSession.UserIPAB) {
		srv.incidentManager.ReportIncident(am.IncidentType_PreSessionHacking, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_UserPreSessionIsNotFound, Message: RpcErrorMsg_LogUserIn_UserPreSessionIsNotFound}
	}

	// Check the captcha answer if needed.
	var ccr *rm.CheckCaptchaResult
	if preSession.IsCaptchaRequired {
		// Get the correct captcha answer by its ID.
		ccr, err = srv.checkCaptcha(preSession.CaptchaId.String, p.CaptchaAnswer)
		if err != nil {
			return nil, &js.Error{Code: RpcErrorCode_RCS_CheckCaptcha, Message: fmt.Sprintf(RpcErrorMsgF_RCS_CheckCaptcha, err.Error())}
		}

		if !ccr.IsSuccess {
			srv.incidentManager.ReportIncident(am.IncidentType_CaptchaAnswerMismatch, p.Email, p.Auth.UserIPAB)

			err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
			if err != nil {
				return nil, srv.databaseError(err)
			}

			// When captcha guess is wrong, we delete the preliminary session
			// to start the process from the first step.
			err = srv.dbo.DeletePreliminarySessionByRequestId(preSession.RequestId)
			if err != nil {
				return nil, srv.databaseError(err)
			}

			return nil, &js.Error{Code: RpcErrorCode_LogUserIn_CaptchaAnswerIsWrong, Message: RpcErrorMsg_LogUserIn_CaptchaAnswerIsWrong}
		}
	}

	// Check the password.
	var ok bool
	ok, err = srv.checkPassword(userId, preSession.AuthDataBytes, p.AuthChallengeResponse)
	if err != nil {
		// The error may occur when not enough data is provided.
		// Such an error is counted as a bad password.
		srv.incidentManager.ReportIncident(am.IncidentType_PasswordMismatch, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// When password is wrong, we delete the preliminary session
		// to start the process from the first step.
		err = srv.dbo.DeletePreliminarySessionByRequestId(preSession.RequestId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_PasswordCheckError, Message: fmt.Sprintf(RpcErrorMsgF_LogUserIn_PasswordCheckError, err.Error())}
	}

	if !ok {
		srv.incidentManager.ReportIncident(am.IncidentType_PasswordMismatch, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// When password is wrong, we delete the preliminary session to start
		// the process from the first step.
		err = srv.dbo.DeletePreliminarySessionByRequestId(preSession.RequestId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_PasswordIsWrong, Message: RpcErrorMsg_LogUserIn_PasswordIsWrong}
	}

	// Set flags.
	//
	// Captcha flag is set to true when either it was not required or it was
	// required and the answer was correct.
	err = srv.dbo.SetUserPreSessionCaptchaFlags(userId, preSession.RequestId, true)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.dbo.SetUserPreSessionPasswordFlag(userId, preSession.RequestId, true)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Create a new Request ID for the next step.
	var step3requestId string
	step3requestId, err = srv.createLogInRequestId()
	if err != nil {
		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_RequestIdGenerator, Message: RpcErrorMsg_LogUserIn_RequestIdGenerator}
	}

	err = srv.dbo.UpdatePreSessionRequestId(userId, preSession.RequestId, step3requestId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Verification by E-mail.
	var verificationCode string
	verificationCode, err = srv.createVerificationCode()
	if err != nil {
		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_VerificationCodeGenerator, Message: RpcErrorMsg_LogUserIn_VerificationCodeGenerator}
	}

	err = srv.dbo.AttachVerificationCodeToPreSession(userId, step3requestId, verificationCode)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.sendVerificationCodeForLogIn(p.Email, verificationCode)
	if err != nil {
		return nil, &js.Error{Code: RpcErrorCode_RegisterUser_SmtpModule, Message: fmt.Sprintf(RpcErrorMsgF_RegisterUser_SmtpModule, err.Error())}
	}

	result = &am.LogUserInResult{NextStep: 3, RequestId: step3requestId}

	return result, nil
}

func (srv *Server) logUserInStep3(p *am.LogUserInParams) (result *am.LogUserInResult, jerr *js.Error) {
	// Is e-mail address valid ?
	if !srv.isUserEmailAddressValid(p.Email) {
		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_EmailAddresIsNotValid, Message: RpcErrorMsg_LogUserIn_EmailAddresIsNotValid}
	}

	if len(p.RequestId) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_RequestIdIsNotSet, Message: RpcErrorMsg_LogUserIn_RequestIdIsNotSet}
	}

	if len(p.VerificationCode) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_VerificationCodeIsNotSet, Message: RpcErrorMsg_LogUserIn_VerificationCodeIsNotSet}
	}

	var canUserLogIn bool
	var err error
	canUserLogIn, err = srv.canUserLogIn(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if !canUserLogIn {
		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_UserCanNotLogIn, Message: RpcErrorMsg_LogUserIn_UserCanNotLogIn}
	}

	err = srv.dbo.DeleteAbandonedPreSessions()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var areUserActiveSessionsPresent bool
	areUserActiveSessionsPresent, err = srv.hasUserActiveSessions(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if areUserActiveSessionsPresent {
		srv.incidentManager.ReportIncident(am.IncidentType_DoubleLogInAttempt, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_UserIsAlreadyLoggedIn, Message: RpcErrorMsg_LogUserIn_UserIsAlreadyLoggedIn}
	}

	var areNoUserActivePreSessionsPresent bool
	areNoUserActivePreSessionsPresent, err = srv.hasNoUserPreSessions(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if areNoUserActivePreSessionsPresent {
		srv.incidentManager.ReportIncident(am.IncidentType_PreSessionHacking, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_UserHasNotStartedToLogIn, Message: RpcErrorMsg_LogUserIn_UserHasNotStartedToLogIn}
	}

	var userId uint
	userId, err = srv.dbo.GetUserIdByEmail(p.Email)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Check the Request ID and IP address.
	var preSession *am.PreSession
	preSession, err = srv.dbo.GetUserPreSessionByRequestId(p.RequestId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if preSession.UserId != userId {
		srv.incidentManager.ReportIncident(am.IncidentType_PreSessionHacking, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_UserPreSessionIsNotFound, Message: RpcErrorMsg_LogUserIn_UserPreSessionIsNotFound}
	}

	if !p.Auth.UserIPAB.Equal(preSession.UserIPAB) {
		srv.incidentManager.ReportIncident(am.IncidentType_PreSessionHacking, p.Email, p.Auth.UserIPAB)

		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_UserPreSessionIsNotFound, Message: RpcErrorMsg_LogUserIn_UserPreSessionIsNotFound}
	}

	// Check the verification code.
	var ok bool
	ok, err = srv.dbo.CheckLogInVerificationCode(p.RequestId, p.VerificationCode)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if !ok {
		// Verification code can not be guessed.
		srv.incidentManager.ReportIncident(am.IncidentType_VerificationCodeMismatch, p.Email, p.Auth.UserIPAB)

		// Delete the pre-session on error to avoid brute force checks.
		err = srv.dbo.UpdateUserLastBadLogInTimeByEmail(p.Email)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// When password is wrong, we delete the preliminary session to start
		// the process from the first step.
		err = srv.dbo.DeletePreliminarySessionByRequestId(preSession.RequestId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_VerificationCodeIsWrong, Message: RpcErrorMsg_LogUserIn_VerificationCodeIsWrong}
	}

	// Set verification flags.
	err = srv.dbo.SetUserPreSessionVerificationFlag(userId, preSession.RequestId, true)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Create a normal session and delete the preliminary one.
	var sessionId int64
	sessionId, err = srv.dbo.CreateUserSession(userId, p.Auth.UserIPAB)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.dbo.DeletePreliminarySessionByRequestId(preSession.RequestId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var tokenString string
	tokenString, err = srv.jwtkm.MakeJWToken(userId, uint(sessionId))
	if err != nil {
		return nil, &js.Error{Code: RpcErrorCode_LogUserIn_JWTCreation, Message: fmt.Sprintf(RpcErrorMsgF_LogUserIn_JWTCreation, err.Error())}
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

	err := srv.dbo.DeleteUserSession(thisUserData.Session.Id, thisUserData.User.Id, thisUserData.Session.UserIPAB)
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

	var err error
	var n int
	n, err = srv.dbo.CountUserSessionsByUserId(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 1 {
		result.IsUserLoggedIn = true
	} else {
		result.IsUserLoggedIn = false
	}

	return result, nil
}

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

	err = srv.dbo.DeleteUserSessionByUserId(p.UserId)
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
