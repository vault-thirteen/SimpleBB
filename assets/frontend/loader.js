// Path to initial settings for a loader script.
settingsPath = "settings.json";

settingsRootPath = "/";

// Various other settings.
redirectDelay = 5;

// Longevity of cached settings. In seconds.
settingsExpirationDuration = 60;

// Timestamps.
varnamePageFirstLoadTime = "pageFirstLoadTime";
varnamePageCurrentLoadTime = "pageCurrentLoadTime";
varnameSettingsLoadTime = "settingsLoadTime";
varnameLogInTime = "logInTime";

// Names of JavaScript storage variables.
varnameSettingsVersion = "settingsVersion";
varnameSettingsProductVersion = "settingsProductVersion";
varnameSettingsSiteName = "settingsSiteName";
varnameSettingsSiteDomain = "settingsSiteDomain";
varnameSettingsCaptchaFolder = "settingsCaptchaFolder";
varnameSettingsSessionMaxDuration = "settingsSessionMaxDuration";
varnameSettingsMessageEditTime = "settingsMessageEditTime";
varnameSettingsPageSize = "settingsPageSize";
varnameSettingsApiFolder = "settingsApiFolder";
varnameSettingsPublicSettingsFileName = "settingsPublicSettingsFileName";
varnameSettingsIsFrontEndEnabled = "settingsIsFrontEndEnabled";
varnameSettingsFrontEndStaticFilesFolder = "settingsFrontEndStaticFilesFolder";
varnameIsLoggedIn = "isLoggedIn";
varnameRegistrationEmail = "registrationEmail";
varnameRegistrationVcode = "registrationVcode";
varnameLogInEmail = "logInEmail";
varnameLogInRequestId = "logInRequestId";
varnameLogInAuthDataBytes = "logInAuthDataBytes";
varnameLogInIsCaptchaNeeded = "logInIsCaptchaNeeded";
varnameLogInCaptchaId = "logInCaptchaId";
varnameChangeEmailRequestId = "changeEmailRequestId";
varnameChangeEmailAuthDataBytes = "changeEmailAuthDataBytes";
varnameChangeEmailIsCaptchaNeeded = "changeEmailIsCaptchaNeeded";
varnameChangeEmailCaptchaId = "changeEmailCaptchaId";
varnameChangePwdRequestId = "changePwdRequestId";
varnameChangePwdAuthDataBytes = "changePwdAuthDataBytes";
varnameChangePwdIsCaptchaNeeded = "changePwdIsCaptchaNeeded";
varnameChangePwdCaptchaId = "changePwdCaptchaId";

// Pages.
qpRegistrationStep1 = "?reg1"
qpRegistrationStep2 = "?reg2"
qpRegistrationStep3 = "?reg3"
qpRegistrationStep4 = "?reg4"
qpLogInStep1 = "?login1"
qpLogInStep2 = "?login2"
qpLogInStep3 = "?login3"
qpLogInStep4 = "?login4"
qpLogOutStep1 = "?logout1"
qpLogOutStep2 = "?logout2"
qpChangeEmailStep1 = "?changeEmail1";
qpChangeEmailStep2 = "?changeEmail2";
qpChangeEmailStep3 = "?changeEmail3";
qpChangePwdStep1 = "?changePwd1";
qpChangePwdStep2 = "?changePwd2";
qpChangePwdStep3 = "?changePwd3";

// Form Input Elements.
fiid1 = "f1i";
fiid2 = "f2i";
fiid3 = "f3i";
fiid4 = "f4i";
fiid4_errflag = "f4ief";
fiid5 = "f5i";
fiid6 = "f6i";
fiid7 = "f7i";
fiid7_image = "f7ii";
fiid8 = "f8i";
fiid9 = "f9i";
fiid10 = "f10i";
fiid11 = "f11i";
fiid11_image = "f11ii";
fiid12 = "f12i";
fiid13 = "f13i";
fiid14 = "f14i";
fiid15 = "f15i";
fiid16 = "f16i";
fiid16_image = "f16ii";
fiid17 = "f17i";

// Errors.
errNextStepUnknown = "unknown next step";
errPasswordNotValid = "password is not valid";
errWebTokenIsNotSet = "web token is not set";
errNotOk = "something went wrong";
errServer = "server error";
errClient = "client error";
errUnknown = "unknown error";
errElementTypeUnsupported = "unsupported element type";

// Messages.
msgRedirecting = "Redirecting. Please wait ...";
msgGenericErrorPrefix = "Error: ";

// Action names.
actionName_registerUser = "registerUser";
actionName_logUserIn = "logUserIn";
actionName_logUserOut = "logUserOut";
actionName_changeEmail = "changeEmail";
actionName_changePwd = "changePassword";

// Settings class.
class Settings {
	constructor(version, productVersion, siteName, siteDomain, captchaFolder,
				sessionMaxDuration, messageEditTime, pageSize, apiFolder,
				publicSettingsFileName, isFrontEndEnabled, frontEndStaticFilesFolder) {
		this.Version = version; // Number.
		this.ProductVersion = productVersion;
		this.SiteName = siteName;
		this.SiteDomain = siteDomain;
		this.CaptchaFolder = captchaFolder;
		this.SessionMaxDuration = sessionMaxDuration; // Number.
		this.MessageEditTime = messageEditTime; // Number.
		this.PageSize = pageSize; // Number.
		this.ApiFolder = apiFolder;
		this.PublicSettingsFileName = publicSettingsFileName;
		this.IsFrontEndEnabled = isFrontEndEnabled; // Boolean.
		this.FrontEndStaticFilesFolder = frontEndStaticFilesFolder;
	}
}

class ApiRequest {
	constructor(action, parameters) {
		this.Action = action;
		this.Parameters = parameters;
	}
}

class Parameters_RegisterUser1 {
	constructor(stepN, email) {
		this.StepN = stepN;
		this.Email = email;
	}
}

class Parameters_RegisterUser2 {
	constructor(stepN, email, verificationCode) {
		this.StepN = stepN;
		this.Email = email;
		this.VerificationCode = verificationCode;
	}
}

class Parameters_RegisterUser3 {
	constructor(stepN, email, verificationCode, name, pwd) {
		this.StepN = stepN;
		this.Email = email;
		this.VerificationCode = verificationCode;
		this.Name = name;
		this.Password = pwd;
	}
}

class Parameters_LogIn1 {
	constructor(stepN, email) {
		this.StepN = stepN;
		this.Email = email;
	}
}

class Parameters_LogIn2 {
	constructor(stepN, email, requestId, captchaAnswer, authChallengeResponse) {
		this.StepN = stepN;
		this.Email = email;
		this.RequestId = requestId;
		this.CaptchaAnswer = captchaAnswer;
		this.AuthChallengeResponse = authChallengeResponse;
	}
}

class Parameters_LogIn3 {
	constructor(stepN, email, requestId, verificationCode) {
		this.StepN = stepN;
		this.Email = email;
		this.RequestId = requestId;
		this.VerificationCode = verificationCode;
	}
}

class Parameters_ChangeEmail1 {
	constructor(stepN, newEmail) {
		this.StepN = stepN;
		this.NewEmail = newEmail;
	}
}

class Parameters_ChangeEmail2 {
	constructor(stepN, requestId, authChallengeResponse, verificationCodeOld,
				verificationCodeNew, captchaAnswer) {
		this.StepN = stepN;
		this.RequestId = requestId;
		this.AuthChallengeResponse = authChallengeResponse;
		this.VerificationCodeOld = verificationCodeOld;
		this.VerificationCodeNew = verificationCodeNew;
		this.CaptchaAnswer = captchaAnswer;
	}
}

class Parameters_ChangePwd1 {
	constructor(stepN, newPassword) {
		this.StepN = stepN;
		this.NewPassword = newPassword;
	}
}

class Parameters_ChangePwd2 {
	constructor(stepN, requestId, authChallengeResponse, vcode, captchaAnswer) {
		this.StepN = stepN;
		this.RequestId = requestId;
		this.AuthChallengeResponse = authChallengeResponse;
		this.VerificationCode = vcode;
		this.CaptchaAnswer = captchaAnswer;
	}
}

class ApiResponse {
	constructor(isOk, jsonObject, statusCode, errorText) {
		this.IsOk = isOk;
		this.JsonObject = jsonObject;
		this.StatusCode = statusCode;
		this.ErrorText = errorText;
	}
}

// Entry point for loader script.
async function loadPage() {
	let pageFirstLoadTime = getPageFirstLoadTime();
	updatePageCurrentLoadTime();
	await updateSettingsIfNeeded();
	let settings = getSettings();
	console.debug("settings:", settings);
	let isLoggedInB = isLoggedIn(settings);
	console.debug("isLoggedIn:", isLoggedInB);

	// Select a page.
	let curUrl = window.location.search;
	switch (curUrl) {
		case qpRegistrationStep1:
			showReg1Form();
			return;

		case qpRegistrationStep2:
			showReg2Form();
			return;

		case qpRegistrationStep3:
			showReg3Form();
			return;

		case qpRegistrationStep4:
			showReg4Form();
			return;

		case qpLogInStep1:
			showLogIn1Form();
			return;

		case qpLogInStep2:
			showLogIn2Form();
			return;

		case qpLogInStep3:
			showLogIn3Form();
			return;

		case qpLogInStep4:
			showLogIn4Form();
			await redirectToMainPage(true);
			return;

		case qpLogOutStep1:
			showLogOut1Form();
			return;

		case qpLogOutStep2:
			showLogOut2Form();
			return;

		case qpChangeEmailStep1:
			showChangeEmail1Form();
			return;

		case qpChangeEmailStep2:
			showChangeEmail2Form();
			return;

		case qpChangeEmailStep3:
			showChangeEmail3Form();
			return;

		case qpChangePwdStep1:
			showChangePwd1Form();
			return;

		case qpChangePwdStep2:
			showChangePwd2Form();
			return;

		case qpChangePwdStep3:
			showChangePwd3Form();
			return;
	}

	if (!isLoggedInB) {
		showLogIn1Form();
		return;
	}

	// Show the forum.
	console.debug("TODO");//TODO: Show forum.
}

function getCurrentTimestamp() {
	return Math.floor(Date.now() / 1000);
}

function stringToBoolean(s) {
	let x = s.trim().toLowerCase();

	switch (x) {
		case "true":
			return true;

		case "false":
			return false;

		case "yes":
		case "1":
			return true;

		case "no":
		case "0":
			return false;

		default:
			return JSON.parse(x);
	}
}

function getPageFirstLoadTime() {
	let pageFirstLoadTimeStr = sessionStorage.getItem(varnamePageFirstLoadTime);

	if (pageFirstLoadTimeStr === null) {
		let timeNow = getCurrentTimestamp();
		sessionStorage.setItem(varnamePageFirstLoadTime, timeNow.toString());
		return timeNow;
	} else {
		return pageFirstLoadTimeStr;
	}
}

function updatePageCurrentLoadTime() {
	let timeNow = getCurrentTimestamp();
	sessionStorage.setItem(varnamePageCurrentLoadTime, timeNow.toString());
}

function getPageCurrentLoadTime() {
	return sessionStorage.getItem(varnamePageCurrentLoadTime);
}

async function updateSettingsIfNeeded() {
	let timeNow = getCurrentTimestamp();
	let settingsLoadTimeStr = sessionStorage.getItem(varnameSettingsLoadTime);

	let settingsAge = 0;
	if (settingsLoadTimeStr != null) {
		settingsAge = timeNow - Number(settingsLoadTimeStr);
	}
	console.debug("settingsAge:", settingsAge);

	if ((settingsLoadTimeStr == null) || (settingsAge > settingsExpirationDuration)) {
		await updateSettings();
		sessionStorage.setItem(varnameSettingsLoadTime, timeNow.toString());
	}
}

async function updateSettings() {
	let settings = await fetchSettings();
	console.debug('Received settings. Version:', settings.version);
	saveSettings(settings);
}

async function fetchSettings() {
	let data = await fetch(settingsRootPath + settingsPath);
	return await data.json();
}

function saveSettings(settings) {
	sessionStorage.setItem(varnameSettingsVersion, settings.version.toString());
	sessionStorage.setItem(varnameSettingsProductVersion, settings.productVersion);
	sessionStorage.setItem(varnameSettingsSiteName, settings.siteName);
	sessionStorage.setItem(varnameSettingsSiteDomain, settings.siteDomain);
	sessionStorage.setItem(varnameSettingsCaptchaFolder, settings.captchaFolder);
	sessionStorage.setItem(varnameSettingsSessionMaxDuration, settings.sessionMaxDuration.toString());
	sessionStorage.setItem(varnameSettingsMessageEditTime, settings.messageEditTime.toString());
	sessionStorage.setItem(varnameSettingsPageSize, settings.pageSize.toString());
	sessionStorage.setItem(varnameSettingsApiFolder, settings.apiFolder);
	sessionStorage.setItem(varnameSettingsPublicSettingsFileName, settings.publicSettingsFileName);
	sessionStorage.setItem(varnameSettingsIsFrontEndEnabled, settings.isFrontEndEnabled.toString());
	sessionStorage.setItem(varnameSettingsFrontEndStaticFilesFolder, settings.frontEndStaticFilesFolder);
}

function getSettings() {
	return new Settings(
		Number(sessionStorage.getItem(varnameSettingsVersion)),
		sessionStorage.getItem(varnameSettingsProductVersion),
		sessionStorage.getItem(varnameSettingsSiteName),
		sessionStorage.getItem(varnameSettingsSiteDomain),
		sessionStorage.getItem(varnameSettingsCaptchaFolder),
		Number(sessionStorage.getItem(varnameSettingsSessionMaxDuration)),
		Number(sessionStorage.getItem(varnameSettingsMessageEditTime)),
		Number(sessionStorage.getItem(varnameSettingsPageSize)),
		sessionStorage.getItem(varnameSettingsApiFolder),
		sessionStorage.getItem(varnameSettingsPublicSettingsFileName),
		stringToBoolean(sessionStorage.getItem(varnameSettingsIsFrontEndEnabled)),
		sessionStorage.getItem(varnameSettingsFrontEndStaticFilesFolder),
	);
}

function isLoggedIn(settings) {
	let isLoggedInStr = localStorage.getItem(varnameIsLoggedIn);
	let isLoggedIn;

	if (isLoggedInStr === null) {
		isLoggedIn = false;
		localStorage.setItem(varnameIsLoggedIn, isLoggedIn.toString());
		return false;
	}

	isLoggedIn = stringToBoolean(isLoggedInStr);
	if (!isLoggedIn) {
		return false;
	}

	// Check if the session is not closed by timeout.
	let logInTime = Number(localStorage.getItem(varnameLogInTime));
	let timeNow = getCurrentTimestamp();
	let sessionAge = timeNow - logInTime;
	if (sessionAge > settings.SessionMaxDuration) {
		isLoggedIn = false;
		localStorage.setItem(varnameIsLoggedIn, isLoggedIn.toString());
		return false;
	}

	return true;
}

function logIn() {
	isLoggedIn = true;
	localStorage.setItem(varnameIsLoggedIn, isLoggedIn.toString());
	let timeNow = getCurrentTimestamp();
	localStorage.setItem(varnameLogInTime, timeNow.toString());
}

function logOut() {
	isLoggedIn = false;
	localStorage.setItem(varnameIsLoggedIn, isLoggedIn.toString());
	localStorage.removeItem(varnameLogInTime);
}

function showBlock(divId) {
	let divReg = document.getElementById(divId);
	divReg.style.display = "block";
}

function showHeader1(elementId) {
	let header1 = document.getElementById(elementId);
	let settings = getSettings();
	header1.innerHTML = settings.SiteName;
}

function showReg1Form() {
	showBlock("divReg1");
	showHeader1("header1TitleReg1");
}

function showReg2Form() {
	showBlock("divReg2");
	showHeader1("header1TitleReg2");
}

function showReg3Form() {
	showBlock("divReg3");
	showHeader1("header1TitleReg3");
}

function showReg4Form() {
	showBlock("divReg4");
	showHeader1("header1TitleReg4");
}

function showLogIn1Form() {
	showBlock("divLogIn1");
	showHeader1("header1TitleLogIn1");
}

function showLogIn2Form() {
	showBlock("divLogIn2");
	showHeader1("header1TitleLogIn2");

	// Captcha (optional).
	let isCaptchaNeeded = stringToBoolean(sessionStorage.getItem(varnameLogInIsCaptchaNeeded));
	let captchaId = sessionStorage.getItem(varnameLogInCaptchaId);
	let cptImageTr = document.getElementById("formHolderLogIn2CaptchaImage");
	let cptImage = document.getElementById(fiid7_image);
	let cptAnswerTr = document.getElementById("formHolderLogIn2CaptchaAnswer");
	let cptAnswer = document.getElementById(fiid7);
	setCaptchaInputsVisibility(isCaptchaNeeded, captchaId, cptImageTr, cptImage, cptAnswerTr, cptAnswer);
}

function showLogIn3Form() {
	showBlock("divLogIn3");
	showHeader1("header1TitleLogIn3");
}

function showLogIn4Form() {
	showBlock("divLogIn4");
	showHeader1("header1TitleLogIn4");
}

function showLogOut1Form() {
	showBlock("divLogOut1");
	showHeader1("header1TitleLogOut1");
}

function showLogOut2Form() {
	showBlock("divLogOut2");
	showHeader1("header1TitleLogOut2");
}

function showChangeEmail1Form() {
	showBlock("divChangeEmail1");
	showHeader1("header1TitleChangeEmail1");
}

function showChangeEmail2Form() {
	showBlock("divChangeEmail2");
	showHeader1("header1TitleChangeEmail2");

	// Captcha (optional).
	let isCaptchaNeeded = stringToBoolean(sessionStorage.getItem(varnameChangeEmailIsCaptchaNeeded));
	let captchaId = sessionStorage.getItem(varnameChangeEmailCaptchaId);
	let cptImageTr = document.getElementById("formHolderChangeEmail2CaptchaImage");
	let cptImage = document.getElementById(fiid11_image);
	let cptAnswerTr = document.getElementById("formHolderChangeEmail2CaptchaAnswer");
	let cptAnswer = document.getElementById(fiid11);
	setCaptchaInputsVisibility(isCaptchaNeeded, captchaId, cptImageTr, cptImage, cptAnswerTr, cptAnswer);
}

function showChangeEmail3Form() {
	showBlock("divChangeEmail3");
	showHeader1("header1TitleChangeEmail3");
}

function showChangePwd1Form() {
	showBlock("divChangePwd1");
	showHeader1("header1TitleChangePwd1");
}

function showChangePwd2Form() {
	showBlock("divChangePwd2");
	showHeader1("header1TitleChangePwd2");

	// Captcha (optional).
	let isCaptchaNeeded = stringToBoolean(sessionStorage.getItem(varnameChangePwdIsCaptchaNeeded));
	let captchaId = sessionStorage.getItem(varnameChangePwdCaptchaId);
	let cptImageTr = document.getElementById("formHolderChangePwd2CaptchaImage");
	let cptImage = document.getElementById(fiid16_image);
	let cptAnswerTr = document.getElementById("formHolderChangePwd2CaptchaAnswer");
	let cptAnswer = document.getElementById(fiid16);
	setCaptchaInputsVisibility(isCaptchaNeeded, captchaId, cptImageTr, cptImage, cptAnswerTr, cptAnswer);
}

function showChangePwd3Form() {
	showBlock("divChangePwd3");
	showHeader1("header1TitleChangePwd3");
}

function setCaptchaInputsVisibility(isCaptchaNeeded, captchaId, cptImageTr, cptImage, cptAnswerTr, cptAnswer) {
	if (isCaptchaNeeded) {
		cptImageTr.style.display = "table-row";
		cptAnswerTr.style.display = "table-row";
		if (captchaId.length > 0) {
			cptImage.src = makeCaptchaImageUrl(captchaId);
		}
		cptAnswer.enabled = true;
	} else {
		cptImageTr.style.display = "none";
		cptAnswerTr.style.display = "none";
		cptAnswer.enabled = false;
	}
}

async function onReg1Submit(btn) {
	console.debug("onReg1Submit");

	// Send the request.
	let h3Field = document.getElementById("header3TextReg1");
	let errField = document.getElementById("header4TextReg1");
	let email = document.getElementById(fiid1).value;
	let params = new Parameters_RegisterUser1(1, email);
	let reqData = new ApiRequest(actionName_registerUser, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if (nextStep !== 2) {
		errField.innerHTML = composeErrorText(errNextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	sessionStorage.setItem(varnameRegistrationEmail, email);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msgRedirecting;
	await redirectToRelativePath(true, qpRegistrationStep2);
}

async function onReg2Submit(btn) {
	console.debug("onReg2Submit");

	// Send the request.
	let h3Field = document.getElementById("header3TextReg2");
	let errField = document.getElementById("header4TextReg2");
	let email = sessionStorage.getItem(varnameRegistrationEmail);
	let vcode = document.getElementById(fiid2).value;
	let params = new Parameters_RegisterUser2(2, email, vcode);
	let reqData = new ApiRequest(actionName_registerUser, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msgRedirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if (nextStep !== 3) {
		errField.innerHTML = composeErrorText(errNextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	sessionStorage.setItem(varnameRegistrationVcode, vcode);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msgRedirecting;
	await redirectToRelativePath(true, qpRegistrationStep3);
}

async function onReg3Submit(btn) {
	console.debug("onReg3Submit");

	// Check the input.
	let pwd = document.getElementById(fiid4).value;
	let pwdErrFlag = document.getElementById(fiid4_errflag);
	let ok = checkPwd(pwd);
	if (ok) {
		pwdErrFlag.className = "flag_none";
	} else {
		pwdErrFlag.className = "flag_error";
		return;
	}

	// Send the request.
	let h3Field = document.getElementById("header3TextReg3");
	let errField = document.getElementById("header4TextReg3");
	let email = sessionStorage.getItem(varnameRegistrationEmail);
	let vcode = sessionStorage.getItem(varnameRegistrationVcode);
	let name = document.getElementById(fiid3).value;
	let params = new Parameters_RegisterUser3(3, email, vcode, name, pwd);
	let reqData = new ApiRequest(actionName_registerUser, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msgRedirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if ((nextStep !== 4) && (nextStep !== 0)) {
		errField.innerHTML = composeErrorText(errNextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Clear saved input data from browser.
	sessionStorage.removeItem(varnameRegistrationEmail);
	sessionStorage.removeItem(varnameRegistrationVcode);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msgRedirecting;
	await redirectToRelativePath(true, qpRegistrationStep4);
}

async function onLogIn1Submit(btn) {
	console.debug("onLogIn1Submit");

	// Send the request.
	let h3Field = document.getElementById("header3TextLogIn1");
	let errField = document.getElementById("header4TextLogIn1");
	let email = document.getElementById(fiid5).value;
	let params = new Parameters_LogIn1(1, email);
	let reqData = new ApiRequest(actionName_logUserIn, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msgRedirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if (nextStep !== 2) {
		errField.innerHTML = composeErrorText(errNextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	sessionStorage.setItem(varnameLogInEmail, email);
	let requestId = resp.JsonObject.result.requestId;
	sessionStorage.setItem(varnameLogInRequestId, requestId);
	let authDataBytes = resp.JsonObject.result.authDataBytes;
	sessionStorage.setItem(varnameLogInAuthDataBytes, authDataBytes);
	let isCaptchaNeeded = resp.JsonObject.result.isCaptchaNeeded;
	sessionStorage.setItem(varnameLogInIsCaptchaNeeded, isCaptchaNeeded.toString());
	let captchaId = resp.JsonObject.result.captchaId;
	sessionStorage.setItem(varnameLogInCaptchaId, captchaId);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msgRedirecting;
	await redirectToRelativePath(true, qpLogInStep2);
}

async function onLogIn2Submit(btn) {
	console.debug("onLogIn2Submit");

	let errField = document.getElementById("header4TextLogIn2");
	let h3Field = document.getElementById("header3TextLogIn2");

	// Captcha (optional).
	let captchaAnswer = document.getElementById(fiid7).value;

	// Secret.
	let authDataBytes = sessionStorage.getItem(varnameLogInAuthDataBytes);
	let saltBA = base64ToByteArray(authDataBytes);
	let pwd = document.getElementById(fiid6).value;
	if (!isPasswordAllowed(pwd)) {
		errField.innerHTML = composeErrorText(errPasswordNotValid);
		return;
	}
	let keyBA = makeHashKey(pwd, saltBA);
	let authChallengeResponse = byteArrayToBase64(keyBA);

	// Send the request.
	let email = sessionStorage.getItem(varnameLogInEmail);
	let requestId = sessionStorage.getItem(varnameLogInRequestId);
	let params = new Parameters_LogIn2(2, email, requestId, captchaAnswer, authChallengeResponse);
	let reqData = new ApiRequest(actionName_logUserIn, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msgRedirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if (nextStep !== 3) {
		errField.innerHTML = composeErrorText(errNextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Clear saved input data from browser.
	sessionStorage.removeItem(varnameLogInAuthDataBytes);
	sessionStorage.removeItem(varnameLogInIsCaptchaNeeded);
	sessionStorage.removeItem(varnameLogInCaptchaId);

	// Save some non-sensitive input data into browser for the next page.
	let newRequestId = resp.JsonObject.result.requestId;
	sessionStorage.setItem(varnameLogInRequestId, newRequestId);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msgRedirecting;
	await redirectToRelativePath(true, qpLogInStep3);
}

async function onLogIn3Submit(btn) {
	console.log("onLogIn3Submit");

	let errField = document.getElementById("header4TextLogIn3");
	let h3Field = document.getElementById("header3TextLogIn3");

	// Send the request.
	let vcode = document.getElementById(fiid8).value;
	let email = sessionStorage.getItem(varnameLogInEmail);
	let requestId = sessionStorage.getItem(varnameLogInRequestId);
	let params = new Parameters_LogIn3(3, email, requestId, vcode);
	let reqData = new ApiRequest(actionName_logUserIn, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msgRedirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	let isWebTokenSet = resp.JsonObject.result.isWebTokenSet;
	if (nextStep !== 0) {
		errField.innerHTML = composeErrorText(errNextStepUnknown);
		return;
	}
	if (!isWebTokenSet) {
		errField.innerHTML = composeErrorText(errWebTokenIsNotSet);
		return;
	}
	errField.innerHTML = "";

	// Clear saved input data from browser.
	sessionStorage.removeItem(varnameLogInEmail);
	sessionStorage.removeItem(varnameLogInRequestId);

	// Save the 'log' flag.
	logIn();

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msgRedirecting;
	await redirectToMainPage(true);
}

async function onLogOut1Submit(btn) {
	console.debug("onLogOut1Submit");

	let errField = document.getElementById("header4TextLogOut1");
	let h3Field = document.getElementById("header3TextLogOut1");

	// Send the request.
	let reqData = new ApiRequest(actionName_logUserOut, {});
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msgRedirecting;
		await redirectToMainPage(true);
		return;
	}
	let ok = resp.JsonObject.result.ok;
	if (!ok) {
		errField.innerHTML = composeErrorText(errNotOk);
		return;
	}
	errField.innerHTML = "";

	// Save the 'log' flag.
	logOut();

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msgRedirecting;
	await redirectToRelativePath(true, qpLogOutStep2);
}

async function onChangeEmail1Submit(btn) {
	console.debug("onChangeEmail1Submit");

	// Send the request.
	let h3Field = document.getElementById("header3TextChangeEmail1");
	let errField = document.getElementById("header4TextChangeEmail1");
	let newEmail = document.getElementById(fiid9).value;
	let params = new Parameters_ChangeEmail1(1, newEmail);
	let reqData = new ApiRequest(actionName_changeEmail, params);
	let resp = await sendApiRequest(reqData);
	console.debug("resp.JsonObject:", resp.JsonObject);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msgRedirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if (nextStep !== 2) {
		errField.innerHTML = composeErrorText(errNextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	let requestId = resp.JsonObject.result.requestId;
	sessionStorage.setItem(varnameChangeEmailRequestId, requestId);
	let authDataBytes = resp.JsonObject.result.authDataBytes;
	sessionStorage.setItem(varnameChangeEmailAuthDataBytes, authDataBytes);
	let isCaptchaNeeded = resp.JsonObject.result.isCaptchaNeeded;
	sessionStorage.setItem(varnameChangeEmailIsCaptchaNeeded, isCaptchaNeeded.toString());
	let captchaId = resp.JsonObject.result.captchaId;
	sessionStorage.setItem(varnameChangeEmailCaptchaId, captchaId);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msgRedirecting;
	await redirectToRelativePath(true, qpChangeEmailStep2);
}

async function onChangeEmail2Submit(btn) {
	console.debug("onChangeEmail2Submit");

	let h3Field = document.getElementById("header3TextChangeEmail2");
	let errField = document.getElementById("header4TextChangeEmail2");

	// Captcha (optional).
	let captchaAnswer = document.getElementById(fiid11).value;

	// Secret.
	let authDataBytes = sessionStorage.getItem(varnameChangeEmailAuthDataBytes);
	let saltBA = base64ToByteArray(authDataBytes);
	let pwd = document.getElementById(fiid10).value;
	if (!isPasswordAllowed(pwd)) {
		errField.innerHTML = composeErrorText(errPasswordNotValid);
		return;
	}
	let keyBA = makeHashKey(pwd, saltBA);
	let authChallengeResponse = byteArrayToBase64(keyBA);

	// Send the request.
	let requestId = sessionStorage.getItem(varnameChangeEmailRequestId);
	let verificationCodeOld = document.getElementById(fiid12).value;
	let verificationCodeNew = document.getElementById(fiid13).value;
	let params = new Parameters_ChangeEmail2(2, requestId, authChallengeResponse, verificationCodeOld, verificationCodeNew, captchaAnswer);
	let reqData = new ApiRequest(actionName_changeEmail, params);
	let resp = await sendApiRequest(reqData);
	console.debug("resp.JsonObject:", resp.JsonObject);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msgRedirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if (nextStep !== 0) {
		errField.innerHTML = composeErrorText(errNextStepUnknown);
		return;
	}
	let ok = resp.JsonObject.result.ok;
	if (!ok) {
		errField.innerHTML = composeErrorText(errNotOk);
		return;
	}
	errField.innerHTML = "";

	// Clear saved input data from browser.
	sessionStorage.removeItem(varnameChangeEmailRequestId);
	sessionStorage.removeItem(varnameChangeEmailAuthDataBytes);
	sessionStorage.removeItem(varnameChangeEmailIsCaptchaNeeded);
	sessionStorage.removeItem(varnameChangeEmailCaptchaId);

	// Save the 'log' flag.
	logOut();

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msgRedirecting;
	await redirectToRelativePath(true, qpChangeEmailStep3);
}

async function onChangePwd1Submit(btn) {
	console.debug("onChangePwd1Submit");

	// Send the request.
	let h3Field = document.getElementById("header3TextChangePwd1");
	let errField = document.getElementById("header4TextChangePwd1");
	let newPwd = document.getElementById(fiid14).value;
	let params = new Parameters_ChangePwd1(1, newPwd);
	let reqData = new ApiRequest(actionName_changePwd, params);
	let resp = await sendApiRequest(reqData);
	console.debug("resp.JsonObject:", resp.JsonObject);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msgRedirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if (nextStep !== 2) {
		errField.innerHTML = composeErrorText(errNextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	let requestId = resp.JsonObject.result.requestId;
	sessionStorage.setItem(varnameChangePwdRequestId, requestId);
	let authDataBytes = resp.JsonObject.result.authDataBytes;
	sessionStorage.setItem(varnameChangePwdAuthDataBytes, authDataBytes);
	let isCaptchaNeeded = resp.JsonObject.result.isCaptchaNeeded;
	sessionStorage.setItem(varnameChangePwdIsCaptchaNeeded, isCaptchaNeeded.toString());
	let captchaId = resp.JsonObject.result.captchaId;
	sessionStorage.setItem(varnameChangePwdCaptchaId, captchaId);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msgRedirecting;
	await redirectToRelativePath(true, qpChangePwdStep2);
}

async function onChangePwd2Submit(btn) {
	console.debug("onChangePwd2Submit");

	let h3Field = document.getElementById("header3TextChangePwd2");
	let errField = document.getElementById("header4TextChangePwd2");

	// Captcha (optional).
	let captchaAnswer = document.getElementById(fiid16).value;

	// Secret.
	let authDataBytes = sessionStorage.getItem(varnameChangePwdAuthDataBytes);
	let saltBA = base64ToByteArray(authDataBytes);
	let pwd = document.getElementById(fiid15).value;
	if (!isPasswordAllowed(pwd)) {
		errField.innerHTML = composeErrorText(errPasswordNotValid);
		return;
	}
	let keyBA = makeHashKey(pwd, saltBA);
	let authChallengeResponse = byteArrayToBase64(keyBA);

	// Send the request.
	let requestId = sessionStorage.getItem(varnameChangePwdRequestId);
	let vcode = document.getElementById(fiid17).value;
	let params = new Parameters_ChangePwd2(2, requestId, authChallengeResponse, vcode, captchaAnswer);
	let reqData = new ApiRequest(actionName_changePwd, params);
	let resp = await sendApiRequest(reqData);
	console.debug("resp.JsonObject:", resp.JsonObject);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msgRedirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if (nextStep !== 0) {
		errField.innerHTML = composeErrorText(errNextStepUnknown);
		return;
	}
	let ok = resp.JsonObject.result.ok;
	if (!ok) {
		errField.innerHTML = composeErrorText(errNotOk);
		return;
	}
	errField.innerHTML = "";

	// Clear saved input data from browser.
	sessionStorage.removeItem(varnameChangePwdRequestId);
	sessionStorage.removeItem(varnameChangePwdAuthDataBytes);
	sessionStorage.removeItem(varnameChangePwdIsCaptchaNeeded);
	sessionStorage.removeItem(varnameChangePwdCaptchaId);

	// Save the 'log' flag.
	logOut();

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msgRedirecting;
	await redirectToRelativePath(true, qpChangePwdStep3);
}

async function sleep(ms) {
	await new Promise(r => setTimeout(r, ms));
}

async function sendApiRequest(data) {
	console.debug("sendApiRequest", "data:", data);//TODO:DEBUG.
	let settings = getSettings();
	let url = settingsRootPath + settings.ApiFolder;
	let resp;
	let result;
	let ri = {
		method: "POST",
		body: JSON.stringify(data)
	};
	resp = await fetch(url, ri);
	if (resp.status === 200) {
		result = new ApiResponse(true, await resp.json(), resp.status, null);
		return result;
	} else {
		result = new ApiResponse(false, null, resp.status, await resp.text());
		if (result.ErrorText.length === 0) {
			result.ErrorText = createErrorTextByStatusCode(result.StatusCode);
		}
		console.error(result.ErrorText);
		return result;
	}
}

function createErrorTextByStatusCode(statusCode) {
	if ((statusCode >= 400) && (statusCode <= 499)) {
		return msgGenericErrorPrefix + errClient + " (" + statusCode.toString() + ")";
	}
	if ((statusCode >= 500) && (statusCode <= 599)) {
		return msgGenericErrorPrefix + errServer + " (" + statusCode.toString() + ")";
	}
	return msgGenericErrorPrefix + errUnknown + " (" + statusCode.toString() + ")";
}

function composeErrorText(errMsg) {
	return msgGenericErrorPrefix + errMsg + ".";
}

async function redirectPage(wait, url) {
	if (wait) {
		await sleep(redirectDelay * 1000);
	}

	document.location.href = url;
}

// redirectToRelativePath redirects to a page with a relative path.
// E.g., if a relative path is 'x', then URL is '/x', where '/' is a front end
// root. Front end root is taken from settings stored in browser's JavaScript
// session storage.
async function redirectToRelativePath(wait, relPath) {
	let url = settingsRootPath + relPath;
	await redirectPage(wait, url);
}

async function redirectToMainPage(wait) {
	await redirectPage(wait, settingsRootPath);
}

function disableButton(btn) {
	//console.debug(btn);
	switch (btn.tagName) {
		case "INPUT":
			btn.value = "";
			btn.disabled = true;
			btn.style.display = "none";
			return;

		default:
			console.error(errElementTypeUnsupported);
	}
}

function checkPwd(pwd) {
	if (pwd.length < 16) {
		return false;
	}

	if ((pwd.length % 4) !== 0) {
		return false;
	}

	let symbol;
	for (let i = 0; i < pwd.length; i++) {
		symbol = pwd.charAt(i);
		if (!checkPwdSymbol(symbol)) {
			return false;
		}
	}

	return true
}

function checkPwdSymbol(symbol) {
	switch (symbol) {
		case '0':
		case '1':
		case '2':
		case '3':
		case '4':
		case '5':
		case '6':
		case '7':
		case '8':
		case '9':
			return true;

		case 'A':
		case 'B':
		case 'C':
		case 'D':
		case 'E':
		case 'F':
		case 'G':
		case 'H':
		case 'I':
		case 'J':
		case 'K':
		case 'L':
		case 'M':
		case 'N':
		case 'O':
		case 'P':
		case 'Q':
		case 'R':
		case 'S':
		case 'T':
		case 'U':
		case 'V':
		case 'W':
		case 'X':
		case 'Y':
		case 'Z':
			return true;

		case ' ':
		case '!':
		case '"':
		case '#':
		case '$':
		case '%':
		case '&':
		case "'":
		case '(':
		case ')':
		case '*':
		case '+':
		case ',':
		case '-':
		case '.':
		case '/':
		case ':':
		case ';':
		case '<':
		case '=':
		case '>':
		case '?':
		case '@':
		case '[':
		case "\\":
		case ']':
		case '^':
		case '_':
			return true;
	}

	return false;
}

function makeCaptchaImageUrl(captchaId) {
	let settings = getSettings();
	let captchaPath = settingsRootPath + settings.CaptchaFolder;
	return captchaPath + "?id=" + captchaId;
}
