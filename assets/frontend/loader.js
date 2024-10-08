window.onpageshow = function (event) {
	if (event.persisted) {
		// Unfortunately, JavaScript does not reload a page when you click
		// "Go Back" button in your web browser. Every year programmers invent
		// a new "wheel" to fix this bug. And every year old working solutions
		// stop working and new ones are invented. This circus looks infinite,
		// but in reality it will end as soon as this evil programming language
		// dies. Please, do not support JavaScript and its developers in any
		// means possible. Please, let this evil "technology" to die.
		console.info("JavaScript must die. This pseudo language is a big mockery and ridicule of people. This is not a joke. This is truth.");
		window.location.reload();
	}
};

// Path to initial settings for a loader script.
settingsPath = "settings.json";
settingsRootPath = "/";
adminPage = "/admin";
settingsExpirationDuration = 60;
redirectDelay = 5;

// Names of JavaScript storage variables.
varname = {
	PageFirstLoadTime: "pageFirstLoadTime",
	PageCurrentLoadTime: "pageCurrentLoadTime",
	SettingsLoadTime: "settingsLoadTime",
	LogInTime: "logInTime",
	SettingsVersion: "settingsVersion",
	SettingsProductVersion: "settingsProductVersion",
	SettingsSiteName: "settingsSiteName",
	SettingsSiteDomain: "settingsSiteDomain",
	SettingsCaptchaFolder: "settingsCaptchaFolder",
	SettingsSessionMaxDuration: "settingsSessionMaxDuration",
	SettingsMessageEditTime: "settingsMessageEditTime",
	SettingsPageSize: "settingsPageSize",
	SettingsApiFolder: "settingsApiFolder",
	SettingsPublicSettingsFileName: "settingsPublicSettingsFileName",
	SettingsIsFrontEndEnabled: "settingsIsFrontEndEnabled",
	SettingsFrontEndStaticFilesFolder: "settingsFrontEndStaticFilesFolder",
	IsLoggedIn: "isLoggedIn",
	RegistrationEmail: "registrationEmail",
	RegistrationVcode: "registrationVcode",
	LogInEmail: "logInEmail",
	LogInRequestId: "logInRequestId",
	LogInAuthDataBytes: "logInAuthDataBytes",
	LogInIsCaptchaNeeded: "logInIsCaptchaNeeded",
	LogInCaptchaId: "logInCaptchaId",
	ChangeEmailRequestId: "changeEmailRequestId",
	ChangeEmailAuthDataBytes: "changeEmailAuthDataBytes",
	ChangeEmailIsCaptchaNeeded: "changeEmailIsCaptchaNeeded",
	ChangeEmailCaptchaId: "changeEmailCaptchaId",
	ChangePwdRequestId: "changePwdRequestId",
	ChangePwdAuthDataBytes: "changePwdAuthDataBytes",
	ChangePwdIsCaptchaNeeded: "changePwdIsCaptchaNeeded",
	ChangePwdCaptchaId: "changePwdCaptchaId",
}

// Pages and Query Parameters.
qp = {
	Prefix: "?",
	RegistrationStep1: "?reg1",
	RegistrationStep2: "?reg2",
	RegistrationStep3: "?reg3",
	RegistrationStep4: "?reg4",
	LogInStep1: "?login1",
	LogInStep2: "?login2",
	LogInStep3: "?login3",
	LogInStep4: "?login4",
	LogOutStep1: "?logout1",
	LogOutStep2: "?logout2",
	ChangeEmailStep1: "?changeEmail1",
	ChangeEmailStep2: "?changeEmail2",
	ChangeEmailStep3: "?changeEmail3",
	ChangePwdStep1: "?changePwd1",
	ChangePwdStep2: "?changePwd2",
	ChangePwdStep3: "?changePwd3",
	SelfPage: "?selfPage",
	Notifications: "?notifications",
}

qpn = {
	Id: "id",
	Page: "page",
	Section: "section",
	Forum: "forum",
	Thread: "thread",
	Message: "message",
}

// Form Input Elements.
fi = {
	id1: "f1i",
	id2: "f2i",
	id3: "f3i",
	id4: "f4i",
	id4_errflag: "f4ief",
	id5: "f5i",
	id6: "f6i",
	id7: "f7i",
	id7_image: "f7ii",
	id8: "f8i",
	id9: "f9i",
	id10: "f10i",
	id11: "f11i",
	id11_image: "f11ii",
	id12: "f12i",
	id13: "f13i",
	id14: "f14i",
	id15: "f15i",
	id16: "f16i",
	id16_image: "f16ii",
	id17: "f17i",
	id18: "f18i",
	id19: "f19i",
	id20: "f20i",
	id21_tr: "f21tr",
}

// Errors.
err = {
	IdNotSet: "ID is not set",
	IdNotFound: "ID is not found",
	PageNotSet: "page is not set",
	PageNotFound: "page is not found",
	NextStepUnknown: "unknown next step",
	PasswordNotValid: "password is not valid",
	WebTokenIsNotSet: "web token is not set",
	NotOk: "something went wrong",
	Server: "server error",
	Client: "client error",
	Unknown: "unknown error",
	ElementTypeUnsupported: "unsupported element type",
	RootSectionNotFound: "root section is not found",
	SectionNotFound: "section is not found",
	ThreadNotFound: "thread is not found",
	MessageNotFound: "message is not found",
	DuplicateMapKey: "duplicate map key",
	UnknownVariant: "unknown variant",
	PreviousPageDoesNotExist: "previous page does not exist",
	NextPageDoesNotExist: "next page does not exist",
}

// Messages.
msg = {
	Redirecting: "Redirecting. Please wait ...",
	GenericErrorPrefix: "Error: ",
}

// Action names.
actionName = {
	AddMessage: "addMessage",
	AddThread: "addThread",
	ChangeEmail: "changeEmail",
	ChangeMessageText: "changeMessageText",
	ChangePwd: "changePassword",
	CountUnreadNotifications: "countUnreadNotifications",
	DeleteNotification: "deleteNotification",
	GetAllNotifications: "getAllNotifications",
	GetLatestMessageOfThread: "getLatestMessageOfThread",
	GetMessage: "getMessage",
	GetSelfRoles: "getSelfRoles",
	GetUserName: "getUserName",
	ListForumAndThreadsOnPage: "listForumAndThreadsOnPage",
	ListSectionsAndForums: "listSectionsAndForums",
	ListThreadAndMessagesOnPage: "listThreadAndMessagesOnPage",
	MarkNotificationAsRead: "markNotificationAsRead",
	LogUserIn: "logUserIn",
	LogUserOut: "logUserOut",
	RegisterUser: "registerUser",
}

// Section settings.
sectionChildType = {
	None: 0,
	Section: 1,
	Forum: 2,
}
sectionMarginDelta = 10;

// Global variables.
class GlobalVariablesContainer {
	constructor(id, page, pages, unc) {
		this.Id = id;
		this.Page = page;
		this.Pages = pages;
		this.UNC = unc;
	}
}

class UserNameCache {
	constructor() {
		this.m = new Map();
	}

	async GetName(userId) {
		if (this.m.has(userId)) {
			return this.m.get(userId);
		}

		let resp = await getUserName(userId);
		if (resp == null) {
			return null;
		}

		let userName = resp.result.userName;
		this.m.set(userId, userName);
		return userName;
	}
}

mca_gvc = new GlobalVariablesContainer(0, 0, 0, new UserNameCache());

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

class Parameters_ListForumAndThreadsOnPage {
	constructor(forumId, page) {
		this.ForumId = forumId;
		this.Page = page;
	}
}

class Parameters_ListThreadAndMessagesOnPage {
	constructor(threadId, page) {
		this.ThreadId = threadId;
		this.Page = page;
	}
}

class Parameters_GetMessage {
	constructor(messageId) {
		this.MessageId = messageId;
	}
}

class Parameters_GetUserName {
	constructor(userId) {
		this.UserId = userId;
	}
}

class Parameters_AddThread {
	constructor(parent, name) {
		this.ForumId = parent;
		this.Name = name;
	}
}

class Parameters_AddMessage {
	constructor(parent, text) {
		this.ThreadId = parent;
		this.Text = text;
	}
}

class Parameters_ChangeMessageText {
	constructor(messageId, text) {
		this.MessageId = messageId;
		this.Text = text;
	}
}

class Parameters_MarkNotificationAsRead {
	constructor(notificationId) {
		this.NotificationId = notificationId;
	}
}

class Parameters_DeleteNotification {
	constructor(notificationId) {
		this.NotificationId = notificationId;
	}
}

class Parameters_GetLatestMessageOfThread {
	constructor(threadId) {
		this.ThreadId = threadId;
	}
}

class SectionNode {
	constructor(section, level) {
		this.Section = section;
		this.Level = level;
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
async function onPageLoad() {
	let pageFirstLoadTime = getPageFirstLoadTime();
	updatePageCurrentLoadTime();
	await updateSettingsIfNeeded();
	let curPage = window.location.search;

	// Redirect to registration.
	switch (curPage) {
		case qp.RegistrationStep1:
			await showReg1Form();
			return;

		case qp.RegistrationStep2:
			await showReg2Form();
			return;

		case qp.RegistrationStep3:
			await showReg3Form();
			return;

		case qp.RegistrationStep4:
			await showReg4Form();
			return;
	}

	// Redirect to logging in.
	switch (curPage) {
		case qp.LogInStep1:
			await showLogIn1Form();
			return;

		case qp.LogInStep2:
			await showLogIn2Form();
			return;

		case qp.LogInStep3:
			await showLogIn3Form();
			return;

		case qp.LogInStep4:
			await showLogIn4Form();
			await redirectToMainPage(true);
			return;
	}

	let settings = getSettings();
	let isLoggedInB = isLoggedIn(settings);
	if (!isLoggedInB) {
		await showLogIn1Form();
		return;
	}

	// Pages for logged users.
	switch (curPage) {
		case qp.LogOutStep1:
			await showLogOut1Form();
			return;

		case qp.LogOutStep2:
			await showLogOut2Form();
			return;

		case qp.ChangeEmailStep1:
			await showChangeEmail1Form();
			return;

		case qp.ChangeEmailStep2:
			await showChangeEmail2Form();
			return;

		case qp.ChangeEmailStep3:
			await showChangeEmail3Form();
			return;

		case qp.ChangePwdStep1:
			await showChangePwd1Form();
			return;

		case qp.ChangePwdStep2:
			await showChangePwd2Form();
			return;

		case qp.ChangePwdStep3:
			await showChangePwd3Form();
			return;

		case qp.SelfPage:
			await showUserPage();
			return;

		case qp.Notifications:
			await showNotificationsPage();
			return;
	}

	// Show the bulletin board.
	let sp = new URLSearchParams(curPage);
	if (sp.has(qpn.Section)) {
		if (!prepareIdVariable(sp)) {
			return;
		}
		await showSection();
		return;
	}

	if (sp.has(qpn.Forum)) {
		if ((!prepareIdVariable(sp)) || (!preparePageVariable(sp))) {
			return;
		}
		await showForum();
		return;
	}

	if (sp.has(qpn.Thread)) {
		if ((!prepareIdVariable(sp)) || (!preparePageVariable(sp))) {
			return;
		}
		await showThread();
		return;
	}

	if (sp.has(qpn.Message)) {
		if ((!prepareIdVariable(sp))) {
			return;
		}
		await showMessage();
		return;
	}

	await showBB();
}

function getCurrentTimestamp() {
	return Math.floor(Date.now() / 1000);
}

function stringToBoolean(s) {
	if (s == null) {
		return null;
	}

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
	let pageFirstLoadTimeStr = sessionStorage.getItem(varname.PageFirstLoadTime);

	if (pageFirstLoadTimeStr === null) {
		let timeNow = getCurrentTimestamp();
		sessionStorage.setItem(varname.PageFirstLoadTime, timeNow.toString());
		return timeNow;
	} else {
		return pageFirstLoadTimeStr;
	}
}

function updatePageCurrentLoadTime() {
	let timeNow = getCurrentTimestamp();
	sessionStorage.setItem(varname.PageCurrentLoadTime, timeNow.toString());
}

async function updateSettingsIfNeeded() {
	let timeNow = getCurrentTimestamp();
	let settingsLoadTimeStr = sessionStorage.getItem(varname.SettingsLoadTime);

	let settingsAge = 0;
	if (settingsLoadTimeStr != null) {
		settingsAge = timeNow - Number(settingsLoadTimeStr);
	}

	if ((settingsLoadTimeStr == null) || (settingsAge > settingsExpirationDuration)) {
		await updateSettings();
		sessionStorage.setItem(varname.SettingsLoadTime, timeNow.toString());
	}
}

async function updateSettings() {
	let settings = await fetchSettings();
	console.info('Received settings. Version: ' + settings.version.toString() + ".");
	saveSettings(settings);
}

async function fetchSettings() {
	let data = await fetch(settingsRootPath + settingsPath);
	return await data.json();
}

function saveSettings(settings) {
	sessionStorage.setItem(varname.SettingsVersion, settings.version.toString());
	sessionStorage.setItem(varname.SettingsProductVersion, settings.productVersion);
	sessionStorage.setItem(varname.SettingsSiteName, settings.siteName);
	sessionStorage.setItem(varname.SettingsSiteDomain, settings.siteDomain);
	sessionStorage.setItem(varname.SettingsCaptchaFolder, settings.captchaFolder);
	sessionStorage.setItem(varname.SettingsSessionMaxDuration, settings.sessionMaxDuration.toString());
	sessionStorage.setItem(varname.SettingsMessageEditTime, settings.messageEditTime.toString());
	sessionStorage.setItem(varname.SettingsPageSize, settings.pageSize.toString());
	sessionStorage.setItem(varname.SettingsApiFolder, settings.apiFolder);
	sessionStorage.setItem(varname.SettingsPublicSettingsFileName, settings.publicSettingsFileName);
	sessionStorage.setItem(varname.SettingsIsFrontEndEnabled, settings.isFrontEndEnabled.toString());
	sessionStorage.setItem(varname.SettingsFrontEndStaticFilesFolder, settings.frontEndStaticFilesFolder);
}

function getSettings() {
	return new Settings(
		Number(sessionStorage.getItem(varname.SettingsVersion)),
		sessionStorage.getItem(varname.SettingsProductVersion),
		sessionStorage.getItem(varname.SettingsSiteName),
		sessionStorage.getItem(varname.SettingsSiteDomain),
		sessionStorage.getItem(varname.SettingsCaptchaFolder),
		Number(sessionStorage.getItem(varname.SettingsSessionMaxDuration)),
		Number(sessionStorage.getItem(varname.SettingsMessageEditTime)),
		Number(sessionStorage.getItem(varname.SettingsPageSize)),
		sessionStorage.getItem(varname.SettingsApiFolder),
		sessionStorage.getItem(varname.SettingsPublicSettingsFileName),
		stringToBoolean(sessionStorage.getItem(varname.SettingsIsFrontEndEnabled)),
		sessionStorage.getItem(varname.SettingsFrontEndStaticFilesFolder),
	);
}

function isLoggedIn(settings) {
	let isLoggedInStr = localStorage.getItem(varname.IsLoggedIn);
	let isLoggedIn;

	if (isLoggedInStr === null) {
		isLoggedIn = false;
		localStorage.setItem(varname.IsLoggedIn, isLoggedIn.toString());
		return false;
	}

	isLoggedIn = stringToBoolean(isLoggedInStr);
	if (!isLoggedIn) {
		return false;
	}

	// Check if the session is not closed by timeout.
	let logInTime = Number(localStorage.getItem(varname.LogInTime));
	let timeNow = getCurrentTimestamp();
	let sessionAge = timeNow - logInTime;
	if (sessionAge > settings.SessionMaxDuration) {
		isLoggedIn = false;
		localStorage.setItem(varname.IsLoggedIn, isLoggedIn.toString());
		return false;
	}

	return true;
}

function logIn() {
	isLoggedIn = true;
	localStorage.setItem(varname.IsLoggedIn, isLoggedIn.toString());
	let timeNow = getCurrentTimestamp();
	localStorage.setItem(varname.LogInTime, timeNow.toString());
}

function logOut() {
	isLoggedIn = false;
	localStorage.setItem(varname.IsLoggedIn, isLoggedIn.toString());
	localStorage.removeItem(varname.LogInTime);
}

function showBlock(block) {
	block.style.display = "block";
}

async function showReg1Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divReg1");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showReg2Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divReg2");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showReg3Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divReg3");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showReg4Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divReg4");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showLogIn1Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divLogIn1");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showLogIn2Form() {
	// Captcha (optional).
	let isCaptchaNeeded = stringToBoolean(sessionStorage.getItem(varname.LogInIsCaptchaNeeded));
	let captchaId = sessionStorage.getItem(varname.LogInCaptchaId);
	let cptImageTr = document.getElementById("formHolderLogIn2CaptchaImage");
	let cptImage = document.getElementById(fi.id7_image);
	let cptAnswerTr = document.getElementById("formHolderLogIn2CaptchaAnswer");
	let cptAnswer = document.getElementById(fi.id7);
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divLogIn2");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
	setCaptchaInputsVisibility(isCaptchaNeeded, captchaId, cptImageTr, cptImage, cptAnswerTr, cptAnswer);
}

async function showLogIn3Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divLogIn3");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showLogIn4Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divLogIn4");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showLogOut1Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divLogOut1");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showLogOut2Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divLogOut2");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showChangeEmail1Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divChangeEmail1");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showChangeEmail2Form() {
	// Captcha (optional).
	let isCaptchaNeeded = stringToBoolean(sessionStorage.getItem(varname.ChangeEmailIsCaptchaNeeded));
	let captchaId = sessionStorage.getItem(varname.ChangeEmailCaptchaId);
	let cptImageTr = document.getElementById("formHolderChangeEmail2CaptchaImage");
	let cptImage = document.getElementById(fi.id11_image);
	let cptAnswerTr = document.getElementById("formHolderChangeEmail2CaptchaAnswer");
	let cptAnswer = document.getElementById(fi.id11);
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divChangeEmail2");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
	setCaptchaInputsVisibility(isCaptchaNeeded, captchaId, cptImageTr, cptImage, cptAnswerTr, cptAnswer);
}

async function showChangeEmail3Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divChangeEmail3");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showChangePwd1Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divChangePwd1");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showChangePwd2Form() {
	// Captcha (optional).
	let isCaptchaNeeded = stringToBoolean(sessionStorage.getItem(varname.ChangePwdIsCaptchaNeeded));
	let captchaId = sessionStorage.getItem(varname.ChangePwdCaptchaId);
	let cptImageTr = document.getElementById("formHolderChangePwd2CaptchaImage");
	let cptImage = document.getElementById(fi.id16_image);
	let cptAnswerTr = document.getElementById("formHolderChangePwd2CaptchaAnswer");
	let cptAnswer = document.getElementById(fi.id16);
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divChangePwd2");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
	setCaptchaInputsVisibility(isCaptchaNeeded, captchaId, cptImageTr, cptImage, cptAnswerTr, cptAnswer);
}

async function showChangePwd3Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divChangePwd3");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showUserPage() {
	let resp = await getSelfRoles();
	if (resp == null) {
		return;
	}
	let userParams = resp.result;
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divUserPage");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
	fillUserPage(userParams);
}

async function showBB() {
	let resp = await listSectionsAndForums();
	if (resp == null) {
		return;
	}
	let sections = resp.result.saf.sections;
	let forums = resp.result.saf.forums;
	let rootSectionIdx = getRootSectionIdx(sections);
	if (rootSectionIdx == null) {
		console.error(err.RootSectionNotFound);
	}
	let rootSection = sections[rootSectionIdx];
	let sectionsMap = putArrayItemsIntoMap(sections);
	if (sectionsMap == null) {
		return;
	}
	let forumsMap = putArrayItemsIntoMap(forums);
	if (forumsMap == null) {
		return;
	}
	let nodes = [];
	createTreeOfSections(rootSection, sectionsMap, 1, nodes);
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divBB");
	showBlock(p);
	await addPageHead(p, settings.SiteName, false);
	addActionPanel(p, false);
	processSectionNodes(p, nodes, forumsMap);
}

async function showSection() {
	let sectionId = mca_gvc.Id;
	let resp = await listSectionsAndForums();
	if (resp == null) {
		return;
	}
	let sections = resp.result.saf.sections;
	let forums = resp.result.saf.forums;
	let rootSectionIdx = getRootSectionIdx(sections);
	if (rootSectionIdx == null) {
		console.error(err.RootSectionNotFound);
	}
	let rootSection = sections[rootSectionIdx];
	let sectionsMap = putArrayItemsIntoMap(sections);
	if (sectionsMap == null) {
		return;
	}
	let forumsMap = putArrayItemsIntoMap(forums);
	if (forumsMap == null) {
		return;
	}
	let allNodes = [];
	createTreeOfSections(rootSection, sectionsMap, 1, allNodes);
	let nodes = [];
	if (!sectionsMap.has(sectionId)) {
		console.error(err.SectionNotFound);
		return;
	}
	let curSection = sectionsMap.get(sectionId);
	let curLevel = findCurrentNodeLevel(allNodes, sectionId);
	createTreeOfSections(curSection, sectionsMap, curLevel, nodes);
	let settings = getSettings();
	let parentId = curSection.parent;

	// Draw.
	let p = document.getElementById("divBB");
	showBlock(p);
	await addPageHead(p, settings.SiteName, false);
	if (parentId != null) {
		addActionPanel(p, false, "section", parentId);
	} else {
		addActionPanel(p, false);
	}
	processSectionNodes(p, nodes, forumsMap);
}

async function showForum() {
	let forumId = mca_gvc.Id;
	let pageNumber = mca_gvc.Page;
	let resp = await listForumAndThreadsOnPage(forumId, pageNumber);
	if (resp == null) {
		return;
	}
	let pageCount = resp.result.fatop.totalPages;
	if ((pageCount === undefined) || (pageCount === 0)) {
		pageCount = 1;
	}
	mca_gvc.Pages = pageCount;

	// Check page number for overflow.
	if (pageNumber > pageCount) {
		console.error(err.PageNotFound);
		return;
	}

	let forum = resp.result.fatop;
	let threads = resp.result.fatop.threads;
	let threadsMap = putArrayItemsIntoMap(threads);
	if (threadsMap == null) {
		return;
	}
	let settings = getSettings();
	let parentId = forum.forumSectionId;

	// Draw.
	let p = document.getElementById("divBB");
	showBlock(p);
	await addPageHead(p, settings.SiteName, false);
	addActionPanel(p, false, "section", parentId);
	addPaginator(p, pageNumber, pageCount, "forumPagePrev", "forumPageNext");
	processForumAndThreads(p, forum, threadsMap);
	await addBottomActionPanel(p, "forum", forumId, forum);
}

async function showThread() {
	let threadId = mca_gvc.Id;
	let pageNumber = mca_gvc.Page;
	let resp = await listThreadAndMessagesOnPage(threadId, pageNumber);
	if (resp == null) {
		return;
	}
	let pageCount = resp.result.tamop.totalPages;
	if ((pageCount === undefined) || (pageCount === 0)) {
		pageCount = 1;
	}
	mca_gvc.Pages = pageCount;

	// Check page number for overflow.
	if (pageNumber > pageCount) {
		console.error(err.PageNotFound);
		return;
	}

	let thread = resp.result.tamop;
	let messages = resp.result.tamop.messages;
	let messagesMap = putArrayItemsIntoMap(messages);
	if (messagesMap == null) {
		return;
	}
	let settings = getSettings();
	let parentId = thread.threadForumId;

	// Draw.
	let p = document.getElementById("divBB");
	showBlock(p);
	await addPageHead(p, settings.SiteName, false);
	addActionPanel(p, false, "forum", parentId);
	addPaginator(p, pageNumber, pageCount, "threadPagePrev", "threadPageNext");
	await processThreadAndMessages(p, thread, messagesMap);
	await addBottomActionPanel(p, "thread", threadId, thread);
}

async function showMessage() {
	let messageId = mca_gvc.Id;
	let resp = await getMessage(messageId);
	if (resp == null) {
		return;
	}
	let message = resp.result.message;
	let settings = getSettings();
	let parentId = message.threadId;

	// Draw.
	let p = document.getElementById("divBB");
	showBlock(p);
	await addPageHead(p, settings.SiteName, false);
	addActionPanel(p, false, "thread", parentId);
	await processMessage(p, message);
	await addBottomActionPanel(p, "message", messageId, message);
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
	// Send the request.
	let h3Field = document.getElementById("header3TextReg1");
	let errField = document.getElementById("header4TextReg1");
	let email = document.getElementById(fi.id1).value;
	let params = new Parameters_RegisterUser1(1, email);
	let reqData = new ApiRequest(actionName.RegisterUser, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if (nextStep !== 2) {
		errField.innerHTML = composeErrorText(err.NextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	sessionStorage.setItem(varname.RegistrationEmail, email);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msg.Redirecting;
	await redirectToRelativePath(true, qp.RegistrationStep2);
}

async function onReg2Submit(btn) {
	// Send the request.
	let h3Field = document.getElementById("header3TextReg2");
	let errField = document.getElementById("header4TextReg2");
	let email = sessionStorage.getItem(varname.RegistrationEmail);
	let vcode = document.getElementById(fi.id2).value;
	let params = new Parameters_RegisterUser2(2, email, vcode);
	let reqData = new ApiRequest(actionName.RegisterUser, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msg.Redirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if (nextStep !== 3) {
		errField.innerHTML = composeErrorText(err.NextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	sessionStorage.setItem(varname.RegistrationVcode, vcode);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msg.Redirecting;
	await redirectToRelativePath(true, qp.RegistrationStep3);
}

async function onReg3Submit(btn) {
	// Check the input.
	let pwd = document.getElementById(fi.id4).value;
	let pwdErrFlag = document.getElementById(fi.id4_errflag);
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
	let email = sessionStorage.getItem(varname.RegistrationEmail);
	let vcode = sessionStorage.getItem(varname.RegistrationVcode);
	let name = document.getElementById(fi.id3).value;
	let params = new Parameters_RegisterUser3(3, email, vcode, name, pwd);
	let reqData = new ApiRequest(actionName.RegisterUser, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msg.Redirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if ((nextStep !== 4) && (nextStep !== 0)) {
		errField.innerHTML = composeErrorText(err.NextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Clear saved input data from browser.
	sessionStorage.removeItem(varname.RegistrationEmail);
	sessionStorage.removeItem(varname.RegistrationVcode);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msg.Redirecting;
	await redirectToRelativePath(true, qp.RegistrationStep4);
}

async function onLogIn1Submit(btn) {
	// Send the request.
	let h3Field = document.getElementById("header3TextLogIn1");
	let errField = document.getElementById("header4TextLogIn1");
	let email = document.getElementById(fi.id5).value;
	let params = new Parameters_LogIn1(1, email);
	let reqData = new ApiRequest(actionName.LogUserIn, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msg.Redirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if (nextStep !== 2) {
		errField.innerHTML = composeErrorText(err.NextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	sessionStorage.setItem(varname.LogInEmail, email);
	let requestId = resp.JsonObject.result.requestId;
	sessionStorage.setItem(varname.LogInRequestId, requestId);
	let authDataBytes = resp.JsonObject.result.authDataBytes;
	sessionStorage.setItem(varname.LogInAuthDataBytes, authDataBytes);
	let isCaptchaNeeded = resp.JsonObject.result.isCaptchaNeeded;
	sessionStorage.setItem(varname.LogInIsCaptchaNeeded, isCaptchaNeeded.toString());
	let captchaId = resp.JsonObject.result.captchaId;
	sessionStorage.setItem(varname.LogInCaptchaId, captchaId);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msg.Redirecting;
	await redirectToRelativePath(true, qp.LogInStep2);
}

async function onLogIn2Submit(btn) {
	let errField = document.getElementById("header4TextLogIn2");
	let h3Field = document.getElementById("header3TextLogIn2");

	// Captcha (optional).
	let captchaAnswer = document.getElementById(fi.id7).value;

	// Secret.
	let authDataBytes = sessionStorage.getItem(varname.LogInAuthDataBytes);
	let saltBA = base64ToByteArray(authDataBytes);
	let pwd = document.getElementById(fi.id6).value;
	if (!isPasswordAllowed(pwd)) {
		errField.innerHTML = composeErrorText(err.PasswordNotValid);
		return;
	}
	let keyBA = makeHashKey(pwd, saltBA);
	let authChallengeResponse = byteArrayToBase64(keyBA);

	// Send the request.
	let email = sessionStorage.getItem(varname.LogInEmail);
	let requestId = sessionStorage.getItem(varname.LogInRequestId);
	let params = new Parameters_LogIn2(2, email, requestId, captchaAnswer, authChallengeResponse);
	let reqData = new ApiRequest(actionName.LogUserIn, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msg.Redirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if (nextStep !== 3) {
		errField.innerHTML = composeErrorText(err.NextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Clear saved input data from browser.
	sessionStorage.removeItem(varname.LogInAuthDataBytes);
	sessionStorage.removeItem(varname.LogInIsCaptchaNeeded);
	sessionStorage.removeItem(varname.LogInCaptchaId);

	// Save some non-sensitive input data into browser for the next page.
	let newRequestId = resp.JsonObject.result.requestId;
	sessionStorage.setItem(varname.LogInRequestId, newRequestId);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msg.Redirecting;
	await redirectToRelativePath(true, qp.LogInStep3);
}

async function onLogIn3Submit(btn) {
	let errField = document.getElementById("header4TextLogIn3");
	let h3Field = document.getElementById("header3TextLogIn3");

	// Send the request.
	let vcode = document.getElementById(fi.id8).value;
	let email = sessionStorage.getItem(varname.LogInEmail);
	let requestId = sessionStorage.getItem(varname.LogInRequestId);
	let params = new Parameters_LogIn3(3, email, requestId, vcode);
	let reqData = new ApiRequest(actionName.LogUserIn, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msg.Redirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	let isWebTokenSet = resp.JsonObject.result.isWebTokenSet;
	if (nextStep !== 0) {
		errField.innerHTML = composeErrorText(err.NextStepUnknown);
		return;
	}
	if (!isWebTokenSet) {
		errField.innerHTML = composeErrorText(err.WebTokenIsNotSet);
		return;
	}
	errField.innerHTML = "";

	// Clear saved input data from browser.
	sessionStorage.removeItem(varname.LogInEmail);
	sessionStorage.removeItem(varname.LogInRequestId);

	// Save the 'log' flag.
	logIn();

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msg.Redirecting;
	await redirectToMainPage(true);
}

async function onLogOut1Submit(btn) {
	let errField = document.getElementById("header4TextLogOut1");
	let h3Field = document.getElementById("header3TextLogOut1");

	// Send the request.
	let reqData = new ApiRequest(actionName.LogUserOut, {});
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msg.Redirecting;
		await redirectToMainPage(true);
		return;
	}
	let ok = resp.JsonObject.result.ok;
	if (!ok) {
		errField.innerHTML = composeErrorText(err.NotOk);
		return;
	}
	errField.innerHTML = "";

	// Save the 'log' flag.
	logOut();

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msg.Redirecting;
	await redirectToRelativePath(true, qp.LogOutStep2);
}

async function onChangeEmail1Submit(btn) {
	// Send the request.
	let h3Field = document.getElementById("header3TextChangeEmail1");
	let errField = document.getElementById("header4TextChangeEmail1");
	let newEmail = document.getElementById(fi.id9).value;
	let params = new Parameters_ChangeEmail1(1, newEmail);
	let reqData = new ApiRequest(actionName.ChangeEmail, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msg.Redirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if (nextStep !== 2) {
		errField.innerHTML = composeErrorText(err.NextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	let requestId = resp.JsonObject.result.requestId;
	sessionStorage.setItem(varname.ChangeEmailRequestId, requestId);
	let authDataBytes = resp.JsonObject.result.authDataBytes;
	sessionStorage.setItem(varname.ChangeEmailAuthDataBytes, authDataBytes);
	let isCaptchaNeeded = resp.JsonObject.result.isCaptchaNeeded;
	sessionStorage.setItem(varname.ChangeEmailIsCaptchaNeeded, isCaptchaNeeded.toString());
	let captchaId = resp.JsonObject.result.captchaId;
	sessionStorage.setItem(varname.ChangeEmailCaptchaId, captchaId);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msg.Redirecting;
	await redirectToRelativePath(true, qp.ChangeEmailStep2);
}

async function onChangeEmail2Submit(btn) {
	let h3Field = document.getElementById("header3TextChangeEmail2");
	let errField = document.getElementById("header4TextChangeEmail2");

	// Captcha (optional).
	let captchaAnswer = document.getElementById(fi.id11).value;

	// Secret.
	let authDataBytes = sessionStorage.getItem(varname.ChangeEmailAuthDataBytes);
	let saltBA = base64ToByteArray(authDataBytes);
	let pwd = document.getElementById(fi.id10).value;
	if (!isPasswordAllowed(pwd)) {
		errField.innerHTML = composeErrorText(err.PasswordNotValid);
		return;
	}
	let keyBA = makeHashKey(pwd, saltBA);
	let authChallengeResponse = byteArrayToBase64(keyBA);

	// Send the request.
	let requestId = sessionStorage.getItem(varname.ChangeEmailRequestId);
	let vCodeOld = document.getElementById(fi.id12).value;
	let vCodeNew = document.getElementById(fi.id13).value;
	let params = new Parameters_ChangeEmail2(2, requestId, authChallengeResponse, vCodeOld, vCodeNew, captchaAnswer);
	let reqData = new ApiRequest(actionName.ChangeEmail, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msg.Redirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if (nextStep !== 0) {
		errField.innerHTML = composeErrorText(err.NextStepUnknown);
		return;
	}
	let ok = resp.JsonObject.result.ok;
	if (!ok) {
		errField.innerHTML = composeErrorText(err.NotOk);
		return;
	}
	errField.innerHTML = "";

	// Clear saved input data from browser.
	sessionStorage.removeItem(varname.ChangeEmailRequestId);
	sessionStorage.removeItem(varname.ChangeEmailAuthDataBytes);
	sessionStorage.removeItem(varname.ChangeEmailIsCaptchaNeeded);
	sessionStorage.removeItem(varname.ChangeEmailCaptchaId);

	// Save the 'log' flag.
	logOut();

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msg.Redirecting;
	await redirectToRelativePath(true, qp.ChangeEmailStep3);
}

async function onChangePwd1Submit(btn) {
	// Send the request.
	let h3Field = document.getElementById("header3TextChangePwd1");
	let errField = document.getElementById("header4TextChangePwd1");
	let newPwd = document.getElementById(fi.id14).value;
	let params = new Parameters_ChangePwd1(1, newPwd);
	let reqData = new ApiRequest(actionName.ChangePwd, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msg.Redirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if (nextStep !== 2) {
		errField.innerHTML = composeErrorText(err.NextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	let requestId = resp.JsonObject.result.requestId;
	sessionStorage.setItem(varname.ChangePwdRequestId, requestId);
	let authDataBytes = resp.JsonObject.result.authDataBytes;
	sessionStorage.setItem(varname.ChangePwdAuthDataBytes, authDataBytes);
	let isCaptchaNeeded = resp.JsonObject.result.isCaptchaNeeded;
	sessionStorage.setItem(varname.ChangePwdIsCaptchaNeeded, isCaptchaNeeded.toString());
	let captchaId = resp.JsonObject.result.captchaId;
	sessionStorage.setItem(varname.ChangePwdCaptchaId, captchaId);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msg.Redirecting;
	await redirectToRelativePath(true, qp.ChangePwdStep2);
}

async function onChangePwd2Submit(btn) {
	let h3Field = document.getElementById("header3TextChangePwd2");
	let errField = document.getElementById("header4TextChangePwd2");

	// Captcha (optional).
	let captchaAnswer = document.getElementById(fi.id16).value;

	// Secret.
	let authDataBytes = sessionStorage.getItem(varname.ChangePwdAuthDataBytes);
	let saltBA = base64ToByteArray(authDataBytes);
	let pwd = document.getElementById(fi.id15).value;
	if (!isPasswordAllowed(pwd)) {
		errField.innerHTML = composeErrorText(err.PasswordNotValid);
		return;
	}
	let keyBA = makeHashKey(pwd, saltBA);
	let authChallengeResponse = byteArrayToBase64(keyBA);

	// Send the request.
	let requestId = sessionStorage.getItem(varname.ChangePwdRequestId);
	let vcode = document.getElementById(fi.id17).value;
	let params = new Parameters_ChangePwd2(2, requestId, authChallengeResponse, vcode, captchaAnswer);
	let reqData = new ApiRequest(actionName.ChangePwd, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		errField.innerHTML = composeErrorText(resp.ErrorText);
		// Redirect to the main page on error.
		disableButton(btn);
		h3Field.innerHTML = msg.Redirecting;
		await redirectToMainPage(true);
		return;
	}
	let nextStep = resp.JsonObject.result.nextStep;
	if (nextStep !== 0) {
		errField.innerHTML = composeErrorText(err.NextStepUnknown);
		return;
	}
	let ok = resp.JsonObject.result.ok;
	if (!ok) {
		errField.innerHTML = composeErrorText(err.NotOk);
		return;
	}
	errField.innerHTML = "";

	// Clear saved input data from browser.
	sessionStorage.removeItem(varname.ChangePwdRequestId);
	sessionStorage.removeItem(varname.ChangePwdAuthDataBytes);
	sessionStorage.removeItem(varname.ChangePwdIsCaptchaNeeded);
	sessionStorage.removeItem(varname.ChangePwdCaptchaId);

	// Save the 'log' flag.
	logOut();

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msg.Redirecting;
	await redirectToRelativePath(true, qp.ChangePwdStep3);
}

async function sleep(ms) {
	await new Promise(r => setTimeout(r, ms));
}

async function sendApiRequest(data) {
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
		return msg.GenericErrorPrefix + err.Client + " (" + statusCode.toString() + ")";
	}
	if ((statusCode >= 500) && (statusCode <= 599)) {
		return msg.GenericErrorPrefix + err.Server + " (" + statusCode.toString() + ")";
	}
	return msg.GenericErrorPrefix + err.Unknown + " (" + statusCode.toString() + ")";
}

function composeErrorText(errMsg) {
	return msg.GenericErrorPrefix + errMsg.trim() + ".";
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
	switch (btn.tagName) {
		case "INPUT":
			btn.value = "";
			btn.disabled = true;
			btn.style.display = "none";
			return;

		default:
			console.error(err.ElementTypeUnsupported);
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

async function addPageHead(el, text, atTop) {
	let settings = getSettings();
	let isLoggedInB = isLoggedIn(settings);
	let unreadNotificationsCount = -1;
	if (isLoggedInB) {
		let resp = await countUnreadNotifications();
		if (resp == null) {
			return;
		}
		unreadNotificationsCount = resp.result.unc;
	}

	// Draw.
	let cn = "pageHead";
	let div = document.createElement("DIV");
	div.className = cn;
	div.id = cn;

	let tbl = document.createElement("TABLE");
	let tr = document.createElement("TR");
	let tdL = document.createElement("TD");
	tdL.className = cn + "L";
	tdL.id = cn + "L";
	tdL.textContent = text;
	tr.appendChild(tdL);
	let tdR = document.createElement("TD");
	tdR.className = cn + "R";
	tdR.id = cn + "R";
	if (isLoggedInB) {
		let html;
		html = '<table><tr>';
		if (unreadNotificationsCount > 0) {
			html += '<td><input type="button" value=" ☼ " class="btnNotificationsOn" onclick="onBtnNotificationsClick(this)" /></td>';
		} else {
			html += '<td><input type="button" value=" ☼ " class="btnNotificationsOff" onclick="onBtnNotificationsClick(this)" /></td>';
		}
		html += '<td><input type="button" value="Account" class="btnAccount" onclick="onBtnAccountClick(this)" /></td>' + '</tr></table>';
		tdR.innerHTML = html;
	} else {
		tdR.cssText = '';
	}
	tr.appendChild(tdR);
	tbl.appendChild(tr);
	div.appendChild(tbl);

	if (atTop) {
		el.insertBefore(div, el.firstChild);
	} else {
		el.appendChild(div);
	}
}

function addActionPanel(el, atTop, type, parentId) {
	let cn = "actionPanel";
	let div = document.createElement("DIV");
	div.className = cn;
	div.id = cn;

	let tbl = document.createElement("TABLE");
	let tr = document.createElement("TR");
	let td = document.createElement("TD");
	td.innerHTML = '<form><input type="button" value="Go to Index" class="btnGoToIndex" onclick="onBtnGoToIndexClick(this)" /></form>';
	tr.appendChild(td);

	switch (type) {
		case "thread":
			td = document.createElement("TD");
			td.innerHTML = '<span class="spacerA">&nbsp;</span>' +
				'<input type="button" value="Back to Thread" class="btnGoToThread" onclick="onBtnGoToThreadClick(' + parentId + ')" />';
			tr.appendChild(td);
			break;

		case "forum":
			td = document.createElement("TD");
			td.innerHTML = '<span class="spacerA">&nbsp;</span>' +
				'<input type="button" value="Back to Forum" class="btnGoToForum" onclick="onBtnGoToForumClick(' + parentId + ')" />';
			tr.appendChild(td);
			break;

		case "section":
			td = document.createElement("TD");
			td.innerHTML = '<span class="spacerA">&nbsp;</span>' +
				'<input type="button" value="Back to Section" class="btnGoToSection" onclick="onBtnGoToSectionClick(' + parentId + ')" />';
			tr.appendChild(td);
			break;
	}

	tbl.appendChild(tr);
	div.appendChild(tbl);

	if (atTop) {
		el.insertBefore(div, el.firstChild);
	} else {
		el.appendChild(div);
	}
}

async function getSelfRoles() {
	let reqData = new ApiRequest(actionName.GetSelfRoles, {});
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function listSectionsAndForums() {
	let reqData = new ApiRequest(actionName.ListSectionsAndForums, {});
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function listForumAndThreadsOnPage(forumId, page) {
	let params = new Parameters_ListForumAndThreadsOnPage(forumId, page);
	let reqData = new ApiRequest(actionName.ListForumAndThreadsOnPage, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function listThreadAndMessagesOnPage(threadId, page) {
	let params = new Parameters_ListThreadAndMessagesOnPage(threadId, page);
	let reqData = new ApiRequest(actionName.ListThreadAndMessagesOnPage, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function getMessage(messageId) {
	let params = new Parameters_GetMessage(messageId);
	let reqData = new ApiRequest(actionName.GetMessage, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

function getRootSectionIdx(sections) {
	for (let i = 0; i < sections.length; i++) {
		if (sections[i].parent == null) {
			return i;
		}
		return null;
	}
}

function putArrayItemsIntoMap(a) {
	let m = new Map();
	if (a == null) {
		return m;
	}

	let key;
	for (let i = 0; i < a.length; i++) {
		key = a[i].id;
		if (m.has(key)) {
			console.error(err.DuplicateMapKey);
			return null;
		}
		m.set(key, a[i]);
	}
	return m;
}

// createTreeOfSections creates a tree of sections.
// 'nodes' is an output parameter.
function createTreeOfSections(section, sectionsMap, level, nodes) {
	if (section == null) {
		return;
	}

	nodes.push(new SectionNode(section, level));

	if (section.childType !== sectionChildType.Section) {
		return;
	}

	let subSectionIds = section.children;
	let subSection;
	level++;
	for (let i = 0; i < subSectionIds.length; i++) {
		subSection = sectionsMap.get(subSectionIds[i]);
		createTreeOfSections(subSection, sectionsMap, level, nodes);
	}
}

function processSectionNodes(p, nodes, forumsMap) {
	let node, divSection, divForum, ml, url, forumId, forum, sectionForums;
	for (let i = 0; i < nodes.length; i++) {
		node = nodes[i];

		divSection = document.createElement("DIV");
		divSection.className = "section";
		divSection.id = "section_" + node.Section.id;
		ml = sectionMarginDelta * node.Level;
		divSection.style.cssText = "margin-left: " + ml + "px";
		url = composeUrlForSection(node.Section.id);
		divSection.innerHTML = "<a href='" + url + "'>" + node.Section.name + "</a>";
		p.appendChild(divSection);

		if (node.Section.childType === sectionChildType.Forum) {
			sectionForums = node.Section.children;
		} else {
			sectionForums = [];
		}
		for (let j = 0; j < sectionForums.length; j++) {
			forumId = sectionForums[j];
			forum = forumsMap.get(forumId);

			divForum = document.createElement("DIV");
			divForum.className = "forum";
			divForum.id = "forum_" + forumId;
			ml = sectionMarginDelta * (node.Level + 1);
			divForum.style.cssText = "margin-left: " + ml + "px";
			url = composeUrlForForum(forumId);
			divForum.innerHTML = "<a href='" + url + "'>" + forum.name + "</a>";
			p.appendChild(divForum);
		}

	}
}

function processForumAndThreads(p, forum, threadsMap) {
	let divForum, divThread, ml, url, threadIds, threadId, thread;

	divForum = document.createElement("DIV");
	divForum.className = "forum";
	divForum.id = "forum_" + forum.forumId;
	ml = sectionMarginDelta;
	divForum.style.cssText = "margin-left: " + ml + "px";
	url = composeUrlForForum(forum.forumId);
	divForum.innerHTML = "<a href='" + url + "'>" + forum.forumName + "</a>";
	p.appendChild(divForum);

	threadIds = forum.threadIds;
	if (threadIds != null) {
		for (let i = 0; i < threadIds.length; i++) {
			threadId = threadIds[i];
			if (!threadsMap.has(threadId)) {
				console.error(err.ThreadNotFound);
				return;
			}
			thread = threadsMap.get(threadId);

			divThread = document.createElement("DIV");
			divThread.className = "thread";
			divThread.id = "thread_" + thread.id;
			ml = sectionMarginDelta * 2;
			divThread.style.cssText = "margin-left: " + ml + "px";
			url = composeUrlForThread(threadId);
			divThread.innerHTML = "<a href='" + url + "'>" + thread.name + "</a>";
			p.appendChild(divThread);
		}
	}
}

async function processThreadAndMessages(p, thread, messagesMap) {
	let divThread, ml, url, messageIds, messageId, message, txt;

	divThread = document.createElement("DIV");
	divThread.className = "thread";
	divThread.id = "thread_" + thread.threadId;
	ml = sectionMarginDelta;
	divThread.style.cssText = "margin-left: " + ml + "px";
	url = composeUrlForThread(thread.threadId);
	divThread.innerHTML = "<a href='" + url + "'>" + thread.threadName + "</a>";
	p.appendChild(divThread);

	messageIds = thread.messageIds;
	if (messageIds == null) {
		return
	}

	let divMsgHdr, divMsgBody;
	for (let i = 0; i < messageIds.length; i++) {
		messageId = messageIds[i];
		if (!messagesMap.has(messageId)) {
			console.error(err.MessageNotFound);
			return;
		}
		message = messagesMap.get(messageId);

		// Header.
		divMsgHdr = document.createElement("DIV");
		divMsgHdr.className = "messageHeader";
		divMsgHdr.id = "messageHeader_" + message.id;
		ml = sectionMarginDelta * 2;
		divMsgHdr.style.cssText = "margin-left: " + ml + "px";
		txt = await composeMessageHeaderText(message);
		divMsgHdr.innerHTML = txt;
		p.appendChild(divMsgHdr);

		// Body.
		divMsgBody = document.createElement("DIV");
		divMsgBody.className = "messageBody";
		divMsgBody.id = "messageBody_" + message.id;
		ml = sectionMarginDelta * 2;
		divMsgBody.style.cssText = "margin-left: " + ml + "px";
		divMsgBody.innerHTML = processMessageText(message.text);
		p.appendChild(divMsgBody);
	}
}

async function processMessage(p, message) {
	let divMsgHdr, divMsgBody, ml, url, txt, creatorName;

	// Header.
	divMsgHdr = document.createElement("DIV");
	divMsgHdr.className = "messageHeader";
	divMsgHdr.id = "messageHeader_" + message.id;
	ml = sectionMarginDelta * 2;
	divMsgHdr.style.cssText = "margin-left: " + ml + "px";
	txt = await composeMessageHeaderText(message);
	divMsgHdr.innerHTML = txt;
	p.appendChild(divMsgHdr);

	// Body.
	divMsgBody = document.createElement("DIV");
	divMsgBody.className = "messageBody";
	divMsgBody.id = "messageBody_" + message.id;
	ml = sectionMarginDelta * 2;
	divMsgBody.style.cssText = "margin-left: " + ml + "px";
	divMsgBody.innerHTML = processMessageText(message.text);
	p.appendChild(divMsgBody);
}

function findCurrentNodeLevel(allNodes, sectionId) {
	let node;
	for (let i = 0; i < allNodes.length; i++) {
		node = allNodes[i];
		if (node.Section.id === sectionId) {
			return node.Level;
		}
	}
}

function addPaginator(el, pageNumber, pageCount, variantPrev, variantNext) {
	let div = document.createElement("DIV");
	div.className = "paginator";
	div.id = "paginator";

	let s = document.createElement("span");
	s.textContent = "Page " + pageNumber + " of " + pageCount + " ";
	div.appendChild(s);

	let btnPrev = document.createElement("input");
	btnPrev.type = "button";
	btnPrev.className = "btnPrev";
	btnPrev.id = "btnPrev";
	btnPrev.value = "Previous Page";
	addClickEventHandler(btnPrev, variantPrev);
	div.appendChild(btnPrev);

	s = document.createElement("span");
	s.className = "spacerA";
	s.innerHTML = "&nbsp;";
	div.appendChild(s);

	let btnNext = document.createElement("input");
	btnNext.type = "button";
	btnNext.className = "btnNext";
	btnNext.id = "btnNext";
	btnNext.value = "Next Page";
	addClickEventHandler(btnNext, variantNext);
	div.appendChild(btnNext);

	el.appendChild(div);
}

function addClickEventHandler(btn, variant) {
	switch (variant) {
		case "forumPagePrev":
			btn.addEventListener("click", async (e) => {
				await onBtnPrevClick_forumPage(btn);
			});
			return;

		case "forumPageNext":
			btn.addEventListener("click", async (e) => {
				await onBtnNextClick_forumPage(btn);
			});
			return;

		case "threadPagePrev":
			btn.addEventListener("click", async (e) => {
				await onBtnPrevClick_threadPage(btn);
			});
			return;

		case "threadPageNext":
			btn.addEventListener("click", async (e) => {
				await onBtnNextClick_threadPage(btn);
			});
			return;

		default:
			console.error(err.UnknownVariant);
	}
}

async function onBtnPrevClick_forumPage(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(err.PreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = composeUrlForForumPage(mca_gvc.Id, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnNextClick_forumPage(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(err.NextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = composeUrlForForumPage(mca_gvc.Id, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnPrevClick_threadPage(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(err.PreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = composeUrlForThreadPage(mca_gvc.Id, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnNextClick_threadPage(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(err.NextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = composeUrlForThreadPage(mca_gvc.Id, mca_gvc.Page);
	await redirectPage(false, url);
}

function prepareIdVariable(sp) {
	if (!sp.has(qpn.Id)) {
		console.error(err.IdNotSet);
		return false;
	}

	let xId = Number(sp.get(qpn.Id));
	if (xId <= 0) {
		console.error(err.IdNotFound);
		return false;
	}

	mca_gvc.Id = xId;
	return true;
}

function preparePageVariable(sp) {
	let pageNumber;
	if (!sp.has(qpn.Page)) {
		pageNumber = 1;
	} else {
		pageNumber = Number(sp.get(qpn.Page));
		if (pageNumber <= 0) {
			console.error(err.PageNotFound);
			return false;
		}
	}

	mca_gvc.Page = pageNumber;
	return true;
}

function fillUserPage(userParams) {
	document.getElementById(fi.id18).value = userParams.name;
	document.getElementById(fi.id19).value = userParams.email;
	document.getElementById(fi.id20).value = prettyTime(userParams.regTime);

	let roles = [];
	if (userParams.isAdministrator) {
		roles.push("Administrator");
	}
	if (userParams.isModerator) {
		roles.push("Moderator");
	}
	if (userParams.isAuthor) {
		roles.push("Author");
	}
	if (userParams.isWriter) {
		roles.push("Writer");
	}
	if (userParams.isReader) {
		roles.push("Reader");
	}
	let rolesHtml = "";
	for (let i = 0; i < roles.length; i++) {
		if (roles[i] !== "Administrator") {
			rolesHtml += '<span class="userPageRole">' + roles[i] + '</span>';
		} else {
			rolesHtml += '<span class="userPageRole"><a href="' + adminPage + '" target="_blank" rel="noopener noreferrer">' + roles[i] + '</a></span>';
		}
	}
	let tr = document.getElementById(fi.id21_tr);
	tr.children[1].innerHTML = rolesHtml;

}

function prettyTime(timeStr) {
	if (timeStr === null) {
		return "";
	}
	if (timeStr.length === 0) {
		return "";
	}

	let t = new Date(timeStr);
	let monthN = t.getUTCMonth() + 1; // Months in JavaScript start with 0 !

	return t.getUTCDate().toString().padStart(2, '0') + "." +
		monthN.toString().padStart(2, '0') + "." +
		t.getUTCFullYear().toString().padStart(4, '0') + " " +
		t.getUTCHours().toString().padStart(2, '0') + ":" +
		t.getUTCMinutes().toString().padStart(2, '0');
}

async function onBtnChangeEmailClick(btn) {
	let url = qp.ChangeEmailStep1;
	await redirectPage(false, url);
}

async function onBtnChangePwdClick(btn) {
	let url = qp.ChangePwdStep1;
	await redirectPage(false, url);
}

async function onBtnLogOutSelfClick(btn) {
	let url = qp.LogOutStep1;
	await redirectPage(false, url);
}

async function onBtnAccountClick(btn) {
	let url = qp.SelfPage;
	await redirectPage(false, url);
}

async function onBtnNotificationsClick(btn) {
	let url = qp.Notifications;
	await redirectPage(false, url);
}

async function onBtnGoToIndexClick(btn) {
	await redirectToMainPage();
}

async function onBtnGoToThreadClick(parentId) {
	let url = composeUrlForThread(parentId);
	await redirectPage(false, url);
}

async function onBtnGoToForumClick(parentId) {
	let url = composeUrlForForum(parentId);
	await redirectPage(false, url);
}

async function onBtnGoToSectionClick(parentId) {
	let url = composeUrlForSection(parentId);
	await redirectPage(false, url);
}

function composeUrlForSection(sectionId) {
	return composeUrlForEntity(qpn.Section, sectionId);
}

function composeUrlForForum(forumId) {
	return composeUrlForEntity(qpn.Forum, forumId);
}

function composeUrlForThread(threadId) {
	return composeUrlForEntity(qpn.Thread, threadId);
}

function composeUrlForMessage(messageId) {
	return composeUrlForEntity(qpn.Message, messageId);
}

function composeUrlForEntity(entityName, entityId) {
	return qp.Prefix + entityName + "&" + qpn.Id + "=" + entityId;
}

function composeUrlForForumPage(forumId, page) {
	return composeUrlForEntityPage(qpn.Forum, forumId, page);
}

function composeUrlForThreadPage(threadId, page) {
	return composeUrlForEntityPage(qpn.Thread, threadId, page);
}

function composeUrlForEntityPage(entityName, entityId, page) {
	return qp.Prefix + entityName + "&" + qpn.Id + "=" + entityId + "&" + qpn.Page + "=" + page;
}

async function getMessageCreatorName(userId) {
	return await mca_gvc.UNC.GetName(userId);
}

async function getUserName(userId) {
	let params = new Parameters_GetUserName(userId);
	let reqData = new ApiRequest(actionName.GetUserName, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function composeMessageHeaderText(message) {
	let messageId = message.id;
	let url = composeUrlForMessage(messageId);
	let creatorName = await getMessageCreatorName(message.creator.userId);
	let txt = "<a href='" + url + "'>" + prettyTime(message.creator.time) + "</a>" +
		' by <span class="messageCreatorName">' + creatorName + '</span>';

	let editorId = message.editor.userId;
	if (editorId != null) {
		let editorName = await getMessageCreatorName(editorId);
		txt += ', edited by <span class="messageEditorName">' + editorName + '</span>' +
			' on <span class="messageEditorTime">' + prettyTime(message.editor.time) + '</span>';
	}

	return txt;
}

async function addBottomActionPanel(el, type, objectId, object) {
	let resp = await getSelfRoles();
	if (resp == null) {
		return;
	}
	let userParams = resp.result;

	let cn = "bottomActionPanel";
	let div = document.createElement("DIV");
	div.className = cn;
	div.id = cn;

	let tbl = document.createElement("TABLE");
	let tr = document.createElement("TR");

	switch (type) {
		case "forum":
			if (userParams.isAuthor) {
				td = document.createElement("TD");
				td.innerHTML = '<form><input type="button" value="Start a new Thread" class="btnStartNewThread" ' +
					'onclick="onBtnStartNewThreadClick(this, \'' + cn + '\', ' + objectId + ')" /></form>';
				tr.appendChild(td);
			}
			break;

		case "thread":
			let resp = await getLatestMessageOfThread(objectId);
			if (resp == null) {
				return;
			}
			let latestMessageInThread = resp.result.message;
			let canAddMsg = canUserAddMessage(userParams, latestMessageInThread);
			if (canAddMsg) {
				td = document.createElement("TD");
				td.innerHTML = '<form><input type="button" value="Add a Message" class="btnAddMessage" ' +
					'onclick="onBtnAddMessageClick(this, \'' + cn + '\', ' + objectId + ')" /></form>';
				tr.appendChild(td);
			}
			break;

		case "message":
			let canEditMsg = canUserEditMessage(userParams, object);
			if (canEditMsg) {
				td = document.createElement("TD");
				td.innerHTML = '<form><input type="button" value="Edit Message" class="btnEditMessage" ' +
					'onclick="onBtnEditMessageClick(this, \'' + cn + '\', ' + objectId + ')" /></form>';
				tr.appendChild(td);
			}
			break;
	}

	tbl.appendChild(tr);
	div.appendChild(tbl);
	el.appendChild(div);
}

async function onBtnStartNewThreadClick(btn, panelCN, forumId) {
	disableButton(btn);
	let p = document.getElementById(panelCN);
	let div = document.createElement("DIV");
	div.className = "newThreadCreation";
	div.id = "newThreadCreation";
	p.appendChild(div);

	// Draw.
	let d = document.createElement("DIV");
	d.className = "title";
	d.textContent = "New Thread Parameters";
	div.appendChild(d);
	d = document.createElement("DIV");
	d.innerHTML = '<label class="parameter" for="name">Name</label>' +
		'<input type="text" name="name" id="name" value="" class="newThreadName" />';
	div.appendChild(d);
	d = document.createElement("DIV");
	d.innerHTML = '<label class="parameter" for="parent" title="ID of a parent forum" hidden="hidden">Parent</label>' +
		'<input type="text" name="parent" id="parent" value="' + forumId + '" readonly="readonly" hidden="hidden"/>';
	div.appendChild(d);
	d = document.createElement("DIV");
	d.innerHTML = '<input type="button" class="btnConfirmThreadStart" value="Confirm" onclick="onBtnConfirmThreadStartClick(this)">';
	div.appendChild(d);
}

async function onBtnAddMessageClick(btn, panelCN, threadId) {
	disableButton(btn);
	let p = document.getElementById(panelCN);
	let div = document.createElement("DIV");
	div.className = "newMessageCreation";
	div.id = "newMessageCreation";
	p.appendChild(div);

	// Draw.
	let d = document.createElement("DIV");
	d.className = "title";
	d.textContent = "New Message Parameters";
	div.appendChild(d);
	d = document.createElement("DIV");
	d.innerHTML = '<label class="parameter" for="txt">Text</label>' +
		'<textarea name="txt" id="txt" class="newMessageText"></textarea>';
	div.appendChild(d);
	d = document.createElement("DIV");
	d.innerHTML = '<label class="parameter" for="parent" title="ID of a parent thread" hidden="hidden">Parent</label>' +
		'<input type="text" name="parent" id="parent" value="' + threadId + '" readonly="readonly" hidden="hidden"/>';
	div.appendChild(d);
	d = document.createElement("DIV");
	d.innerHTML = '<input type="button" class="btnConfirmMessageCreation" value="Send Message" onclick="onBtnConfirmMessageCreationClick(this)">';
	div.appendChild(d);
}

async function onBtnEditMessageClick(btn, panelCN, messageId) {
	// Get edited message.
	let resp = await getMessage(messageId);
	if (resp == null) {
		return;
	}
	let message = resp.result.message;
	let msgText = message.text;

	disableButton(btn);
	let p = document.getElementById(panelCN);
	let div = document.createElement("DIV");
	div.className = "messageEditing";
	div.id = "messageEditing";
	p.appendChild(div);

	// Draw.
	let d = document.createElement("DIV");
	d.className = "title";
	d.textContent = "Message Editing";
	div.appendChild(d);
	d = document.createElement("DIV");
	d.innerHTML = '<label class="parameter" for="txt">Text</label>' +
		'<textarea name="txt" id="txt" class="newMessageText">' + escapeHtml(msgText) + '</textarea>';
	div.appendChild(d);
	d = document.createElement("DIV");
	d.innerHTML = '<label class="parameter" for="id" title="ID of edited message" hidden="hidden">ID</label>' +
		'<input type="text" name="id" id="id" value="' + messageId + '" readonly="readonly" hidden="hidden"/>';
	div.appendChild(d);
	d = document.createElement("DIV");
	d.innerHTML = '<input type="button" class="btnConfirmMessageEdit" value="Save Message" onclick="onBtnConfirmMessageEditClick(this)">';
	div.appendChild(d);
}

async function onBtnConfirmThreadStartClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let name = pp.childNodes[1].childNodes[1].value;
	if (name.length < 1) {
		console.error(err.NameIsNotSet);
		return;
	}
	let parent = Number(pp.childNodes[2].childNodes[1].value);
	if (parent < 1) {
		console.error(err.ParentIsNotSet);
		return;
	}

	// Work.
	let resp = await addThread(parent, name);
	if (resp == null) {
		return;
	}
	let threadId = resp.result.threadId;
	disableParentForm(btn, pp, false);
	let txt = "A thread was created. ID=" + threadId.toString() + ".";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnConfirmMessageCreationClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let text = pp.childNodes[1].childNodes[1].value;
	if (text.length < 1) {
		console.error(err.TextIsNotSet);
		return;
	}
	let parent = Number(pp.childNodes[2].childNodes[1].value);
	if (parent < 1) {
		console.error(err.ParentIsNotSet);
		return;
	}

	// Work.
	let resp = await addMessage(parent, text);
	if (resp == null) {
		return;
	}
	let messageId = resp.result.messageId;
	disableParentForm(btn, pp, false);
	let txt = "A message was created. ID=" + messageId.toString() + ".";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnConfirmMessageEditClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let newText = pp.childNodes[1].childNodes[1].value;
	if (newText.length < 1) {
		console.error(err.TextIsNotSet);
		return;
	}
	let messageId = Number(pp.childNodes[2].childNodes[1].value);
	if (messageId < 1) {
		console.error(err.IdNotSet);
		return;
	}

	// Work.
	let resp = await changeMessageText(messageId, newText);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Message text was changed.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function addThread(parent, name) {
	let params = new Parameters_AddThread(parent, name);
	let reqData = new ApiRequest(actionName.AddThread, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function addMessage(parent, text) {
	let params = new Parameters_AddMessage(parent, text);
	let reqData = new ApiRequest(actionName.AddMessage, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function changeMessageText(messageId, text) {
	let params = new Parameters_ChangeMessageText(messageId, text);
	let reqData = new ApiRequest(actionName.ChangeMessageText, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

function disableParentForm(btn, pp, ignoreButton) {
	if (!ignoreButton) {
		btn.disabled = true;
	}

	let el;
	for (i = 0; i < pp.childNodes.length; i++) {
		let ch = pp.childNodes[i];
		for (j = 0; j < ch.childNodes.length; j++) {
			el = ch.childNodes[j];

			if (el !== btn) {
				el.disabled = true;
			} else {
				if (!ignoreButton) {
					el.disabled = true;
				}
			}
		}
	}
}

function showActionSuccess(btn, txt) {
	let ppp = btn.parentNode.parentNode.parentNode;
	let d = document.createElement("DIV");
	d.className = "actionSuccess";
	d.textContent = txt;
	ppp.appendChild(d);
}

async function reloadPage(wait) {
	if (wait) {
		await sleep(redirectDelay * 1000);
	}
	location.reload();
}

function escapeHtml(text) {
	let div = document.createElement('div');
	div.textContent = text;
	return div.innerHTML;
}

function processMessageText(msgText) {
	let txt = escapeHtml(msgText);
	txt = txt.replaceAll("\r\n", '<br>');
	txt = txt.replaceAll("\n", '<br>');
	txt = txt.replaceAll("\r", '<br>');
	return txt;
}

async function countUnreadNotifications() {
	let reqData = new ApiRequest(actionName.CountUnreadNotifications, {});
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function showNotificationsPage() {
	let resp = await getAllNotifications();
	if (resp == null) {
		return;
	}
	let notifications = resp.result.notifications;
	notifications.sort(notificationsComparer);
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divBB");
	showBlock(p);
	await addPageHead(p, settings.SiteName, false);
	addActionPanel(p, false);
	addDiv(p, "notificationList");
	fillListOfNotifications("notificationList", notifications);
}

function notificationsComparer(a, b) {
	if (a.toc < b.toc) {
		return 1;
	}
	if (a.toc > b.toc) {
		return -1;
	}
	return 0;
}

function addDiv(el, x) {
	let div = document.createElement("DIV");
	div.className = x;
	div.id = x;
	el.appendChild(div);
}

async function getAllNotifications() {
	let reqData = new ApiRequest(actionName.GetAllNotifications, {});
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

function fillListOfNotifications(elClass, notifications) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";

	// Title.
	let title = document.createElement("DIV");
	title.className = elClass + 'Title';
	title.textContent = 'Notifications';
	div.appendChild(title);

	// Table.
	let tbl = document.createElement("TABLE");
	tbl.className = elClass;

	// Header.
	let tr = document.createElement("TR");
	let ths = ["#", "Time", "Text", "Actions"];
	let th;
	for (let i = 0; i < ths.length; i++) {
		th = document.createElement("TH");
		if (i === 0) {
			th.className = "numCol";
		}
		th.textContent = ths[i];
		tr.appendChild(th);
	}
	tbl.appendChild(tr);
	div.appendChild(tbl);

	let columnsWithHtml = [2, 3];

	// Cells.
	let notification, actions;
	for (let i = 0; i < notifications.length; i++) {
		notification = notifications[i];

		tr = document.createElement("TR");
		let tds = [];
		for (let j = 0; j < ths.length; j++) {
			tds.push("");
		}

		tds[0] = (i + 1).toString();
		tds[1] = prettyTime(notification.toc);
		tds[2] = splitNotificationTextCell(notification.text);

		actions = '<input type="button" class="btnShowFullNotificationU" value="Show" onclick="onBtnShowFullNotificationClick(this)">';
		if (!notification.isRead) {
			actions += '<input type="button" class="btnMarkNotificationAsReadU" value="Read" onclick="onBtnMarkNotificationAsReadClick(this, ' + notification.id + ')">';
		} else {
			actions += '<input type="button" class="btnDeleteNotificationU" value="DEL" onclick="onBtnDeleteNotificationClick(this, ' + notification.id + ')">';
		}
		tds[3] = actions;

		let td;
		let jLast = tds.length - 1;
		for (let j = 0; j < tds.length; j++) {
			td = document.createElement("TD");

			if (j === 0) {
				td.className = "numCol";
			} else {
				if (j === jLast) {
					td.className += "lastCol";
				} else {
					if (!notification.isRead) {
						td.className = "unread";
					}
				}
			}

			if (columnsWithHtml.includes(j)) {
				td.innerHTML = tds[j];
			} else {
				td.textContent = tds[j];
			}
			tr.appendChild(td);
		}

		tbl.appendChild(tr);
	}
}

async function onBtnMarkNotificationAsReadClick(btn, notificationId) {
	let resp = await markNotificationAsRead(notificationId);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}

	// Change the 'unread' class in the table's row.
	let tr = btn.parentNode.parentNode;
	let td;
	for (let i = 0; i < tr.childNodes.length; i++) {
		td = tr.childNodes[i];
		if (td.className === 'unread') {
			td.className = 'read';
		}
	}

	// Hide the button.
	disableButton(btn);
}

async function onBtnDeleteNotificationClick(btn, notificationId) {
	let resp = await deleteNotification(notificationId);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}

	// Change the 'unread' class in the table's row.
	let tr = btn.parentNode.parentNode;
	let td;
	for (let i = 0; i < tr.childNodes.length; i++) {
		td = tr.childNodes[i];
		if (td.className === '') {
			td.className = 'deleted';
		}
	}

	// Hide the button.
	disableButton(btn);
}

async function onBtnShowFullNotificationClick(btn) {
	let tr = btn.parentNode.parentNode;
	let td = tr.childNodes[2];
	let subtableTbody = td.childNodes[0].childNodes[0];
	let childShort = subtableTbody.childNodes[0].childNodes[0];
	let childFull = subtableTbody.childNodes[1].childNodes[0];

	if (childShort.className === 'visible') {
		btn.value = "Hide";
		childShort.className = 'hidden';
		childFull.className = 'visible';
	} else {
		btn.value = "Show";
		childFull.className = 'hidden';
		childShort.className = 'visible';
	}
}

async function markNotificationAsRead(notificationId) {
	let params = new Parameters_MarkNotificationAsRead(notificationId);
	let reqData = new ApiRequest(actionName.MarkNotificationAsRead, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function deleteNotification(notificationId) {
	let params = new Parameters_DeleteNotification(notificationId);
	let reqData = new ApiRequest(actionName.DeleteNotification, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

function composeNotificationShortText(fullText) {
	let wordsCountMax = 8;
	let segmenter = new Intl.Segmenter([], {granularity: 'word'});
	let segmentedText = segmenter.segment(fullText);
	let words = [...segmentedText].filter(s => s.isWordLike).map(s => s.segment);
	if (words.length <= wordsCountMax) {
		return fullText;
	}
	let shortTextWOP = words.slice(1, wordsCountMax + 1).join(" ");
	return shortTextWOP + " ...";
}

function splitNotificationTextCell(fullText) {
	let shortText = composeNotificationShortText(fullText);
	let html = '<table>' +
		'<tr><td class="visible">' + shortText + '</td></tr>' +
		'<tr><td class="hidden">' + fullText + '</td></tr>' +
		'</table>';
	return html;
}

function getMessageMaxEditTime(message, settings) {
	let lastTouchTime = getMessageLastTouchTime(message);
	return addTimeSec(lastTouchTime, settings.MessageEditTime);
}

function getMessageLastTouchTime(message) {
	if (message.editor.time == null) {
		return new Date(message.creator.time);
	}
	return new Date(message.editor.time);
}

function addTimeSec(t, deltaSec) {
	return new Date(t.getTime() + deltaSec * 1000);
}

async function getLatestMessageOfThread(threadId) {
	let params = new Parameters_GetLatestMessageOfThread(threadId);
	let reqData = new ApiRequest(actionName.GetLatestMessageOfThread, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

function canUserEditMessage(userParams, message) {
	if (userParams.isModerator) {
		return true;
	}

	if (!userParams.isWriter) {
		return false;
	}

	if (userParams.userId !== message.creator.userId) {
		return false;
	}

	let messageMaxEditTime = getMessageMaxEditTime(message, getSettings());
	if (Date.now() < messageMaxEditTime) {
		return true
	}

	return false;
}

function canUserAddMessage(userParams, latestMessageInThread) {
	if (!userParams.isWriter) {
		return false;
	}

	if (latestMessageInThread == null) {
		return true;
	}

	if (latestMessageInThread.creator.userId !== userParams.userId) {
		return true;
	}

	let messageMaxEditTime = getMessageMaxEditTime(latestMessageInThread, getSettings());
	if (Date.now() < messageMaxEditTime) {
		return false;
	}

	return true;
}
