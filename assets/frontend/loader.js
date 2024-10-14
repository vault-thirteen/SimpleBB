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
rootPath = "/";
adminPage = "/admin";
settingsExpirationDuration = 60;
redirectDelay = 2;

// Names of JavaScript storage variables.
varname = {
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
	id22: "f22i",
}

// Section settings.
sectionChildType = {
	None: 0,
	Section: 1,
	Forum: 2,
}
sectionMarginDelta = 10;

buttonName = {
	BackToRoot: "🏠",
	BackToSection: "Go back",
	BackToForum: "Go back",
	BackToThread: "Go back",
	StartNewThread: "Start a new Thread",
	SubscribeToThread: "Subscribe",
	AddMessage: "Add a Message",
	EditMessage: "Edit Message",
}

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

function isSettingsUpdateNeeded() {
	let settingsLoadTimeStr = sessionStorage.getItem(varname.SettingsLoadTime);
	if (settingsLoadTimeStr == null) {
		return true;
	}

	let timeNow = getCurrentTimestamp();
	let settingsAge = timeNow - Number(settingsLoadTimeStr);
	if (settingsAge >= settingsExpirationDuration) {
		return true;
	}

	return false;
}

function saveSettings(s) {
	sessionStorage.setItem(varname.SettingsVersion, s.Version.toString());
	sessionStorage.setItem(varname.SettingsProductVersion, s.ProductVersion);
	sessionStorage.setItem(varname.SettingsSiteName, s.SiteName);
	sessionStorage.setItem(varname.SettingsSiteDomain, s.SiteDomain);
	sessionStorage.setItem(varname.SettingsCaptchaFolder, s.CaptchaFolder);
	sessionStorage.setItem(varname.SettingsSessionMaxDuration, s.SessionMaxDuration.toString());
	sessionStorage.setItem(varname.SettingsMessageEditTime, s.MessageEditTime.toString());
	sessionStorage.setItem(varname.SettingsPageSize, s.PageSize.toString());
	sessionStorage.setItem(varname.SettingsApiFolder, s.ApiFolder);
	sessionStorage.setItem(varname.SettingsPublicSettingsFileName, s.PublicSettingsFileName);
	sessionStorage.setItem(varname.SettingsIsFrontEndEnabled, s.IsFrontEndEnabled.toString());
	sessionStorage.setItem(varname.SettingsFrontEndStaticFilesFolder, s.FrontEndStaticFilesFolder);

	let timeNow = getCurrentTimestamp();
	sessionStorage.setItem(varname.SettingsLoadTime, timeNow.toString());
}

function getSettings() {
	let settingsLoadTime = sessionStorage.getItem(varname.SettingsLoadTime);
	if (settingsLoadTime == null) {
		console.error(err.Settings);
		return null;
	}

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

// Entry point.
async function onPageLoad() {
	// Settings initialisation.
	let ok = await updateSettingsIfNeeded();
	if (!ok) {
		return;
	}
	let settings = getSettings();

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
	}

	let sp = new URLSearchParams(curPage);

	// Notifications.
	if (sp.has(qpn.Notifications)) {
		if (!preparePageVariable(sp)) {
			return;
		}
		await showPage_Notifications();
		return;
	}

	// Subscriptions.
	if (sp.has(qpn.SubscriptionsPage)) {
		if (!preparePageVariable(sp)) {
			return;
		}
		await showPage_Subscriptions();
		return;
	}

	// Show the bulletin board.
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

function saveLogInStatus() {
	isLoggedIn = true;
	localStorage.setItem(varname.IsLoggedIn, isLoggedIn.toString());
	let timeNow = getCurrentTimestamp();
	localStorage.setItem(varname.LogInTime, timeNow.toString());
}

function saveLogOutStatus() {
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
	await fillUserPage(userParams);
}

async function showBB() {
	// Prepare data.
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
	drawSectionsAndForums(p, nodes, forumsMap);
}

async function showSection() {
	// Prepare data.
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
	drawSectionsAndForums(p, nodes, forumsMap);
}

async function showForum() {
	// Prepare data.
	let forumId = mca_gvc.Id;
	let pageNumber = mca_gvc.Page;
	let resp = await listForumAndThreadsOnPage(forumId, pageNumber);
	if (resp == null) {
		return;
	}
	let pageCount = resp.result.fatop.totalPages;
	pageCount = repairUndefinedPageCount(pageCount);
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
	drawForumAndThreads(p, forum, threadsMap);
	await addBottomActionPanel(p, "forum", forumId, forum);
}

async function showThread() {
	// Prepare data.
	let threadId = mca_gvc.Id;
	let pageNumber = mca_gvc.Page;
	let resp = await listThreadAndMessagesOnPage(threadId, pageNumber);
	if (resp == null) {
		return;
	}
	let pageCount = resp.result.tamop.totalPages;
	pageCount = repairUndefinedPageCount(pageCount);
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
	await drawThreadAndMessages(p, thread, messagesMap);
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
			cptImage.src = composeCaptchaImageUrl(captchaId);
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
	let resp = await registerUser1(1, email);
	if (resp == null) {
		return;
	}
	let nextStep = resp.result.nextStep;
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
	let resp = await registerUser2(2, email, vcode);
	if (resp == null) {
		return;
	}
	let nextStep = resp.result.nextStep;
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
	let resp = await registerUser3(3, email, vcode, name, pwd);
	if (resp == null) {
		return;
	}
	let nextStep = resp.result.nextStep;
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
	let resp = await logUserIn1(1, email);
	if (resp == null) {
		return;
	}
	let nextStep = resp.result.nextStep;
	if (nextStep !== 2) {
		errField.innerHTML = composeErrorText(err.NextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	sessionStorage.setItem(varname.LogInEmail, email);
	let requestId = resp.result.requestId;
	sessionStorage.setItem(varname.LogInRequestId, requestId);
	let authDataBytes = resp.result.authDataBytes;
	sessionStorage.setItem(varname.LogInAuthDataBytes, authDataBytes);
	let isCaptchaNeeded = resp.result.isCaptchaNeeded;
	sessionStorage.setItem(varname.LogInIsCaptchaNeeded, isCaptchaNeeded.toString());
	let captchaId = resp.result.captchaId;
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
	let resp = await logUserIn2(2, email, requestId, captchaAnswer, authChallengeResponse);
	if (resp == null) {
		return;
	}
	let nextStep = resp.result.nextStep;
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
	let newRequestId = resp.result.requestId;
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
	let resp = await logUserIn3(3, email, requestId, vcode);
	if (resp == null) {
		return;
	}
	let nextStep = resp.result.nextStep;
	let isWebTokenSet = resp.result.isWebTokenSet;
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
	saveLogInStatus();

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msg.Redirecting;
	await redirectToMainPage(true);
}

async function onLogOut1Submit(btn) {
	let errField = document.getElementById("header4TextLogOut1");
	let h3Field = document.getElementById("header3TextLogOut1");
	let resp = await logUserOut1();
	if (resp == null) {
		return;
	}
	let ok = resp.result.ok;
	if (!ok) {
		errField.innerHTML = composeErrorText(err.NotOk);
		return;
	}
	errField.innerHTML = "";

	// Save the 'log' flag.
	saveLogOutStatus();

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
	let resp = changeEmail1(1, newEmail);
	if (resp == null) {
		return;
	}
	let nextStep = resp.result.nextStep;
	if (nextStep !== 2) {
		errField.innerHTML = composeErrorText(err.NextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	let requestId = resp.result.requestId;
	sessionStorage.setItem(varname.ChangeEmailRequestId, requestId);
	let authDataBytes = resp.result.authDataBytes;
	sessionStorage.setItem(varname.ChangeEmailAuthDataBytes, authDataBytes);
	let isCaptchaNeeded = resp.result.isCaptchaNeeded;
	sessionStorage.setItem(varname.ChangeEmailIsCaptchaNeeded, isCaptchaNeeded.toString());
	let captchaId = resp.result.captchaId;
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
	let resp = changeEmail2(2, requestId, authChallengeResponse, vCodeOld, vCodeNew, captchaAnswer);
	if (resp == null) {
		return;
	}
	let nextStep = resp.result.nextStep;
	if (nextStep !== 0) {
		errField.innerHTML = composeErrorText(err.NextStepUnknown);
		return;
	}
	let ok = resp.result.ok;
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
	saveLogOutStatus();

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
	let resp = await changePwd1(1, newPwd);
	if (resp == null) {
		return;
	}
	let nextStep = resp.result.nextStep;
	if (nextStep !== 2) {
		errField.innerHTML = composeErrorText(err.NextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	let requestId = resp.result.requestId;
	sessionStorage.setItem(varname.ChangePwdRequestId, requestId);
	let authDataBytes = resp.result.authDataBytes;
	sessionStorage.setItem(varname.ChangePwdAuthDataBytes, authDataBytes);
	let isCaptchaNeeded = resp.result.isCaptchaNeeded;
	sessionStorage.setItem(varname.ChangePwdIsCaptchaNeeded, isCaptchaNeeded.toString());
	let captchaId = resp.result.captchaId;
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
	let resp = await changePwd2(2, requestId, authChallengeResponse, vcode, captchaAnswer);
	if (resp == null) {
		return;
	}
	let nextStep = resp.result.nextStep;
	if (nextStep !== 0) {
		errField.innerHTML = composeErrorText(err.NextStepUnknown);
		return;
	}
	let ok = resp.result.ok;
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
	saveLogOutStatus();

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = msg.Redirecting;
	await redirectToRelativePath(true, qp.ChangePwdStep3);
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
	td.innerHTML = '<form><input type="button" value="' + buttonName.BackToRoot + '" class="btnGoToIndex" onclick="onBtnGoToIndexClick(this)" /></form>';
	tr.appendChild(td);

	switch (type) {
		case "thread":
			td = document.createElement("TD");
			td.innerHTML = '<span class="spacerA">&nbsp;</span>' +
				'<input type="button" value="' + buttonName.BackToThread + '" class="btnGoToThread" onclick="onBtnGoToThreadClick(' + parentId + ')" />';
			tr.appendChild(td);
			break;

		case "forum":
			td = document.createElement("TD");
			td.innerHTML = '<span class="spacerA">&nbsp;</span>' +
				'<input type="button" value="' + buttonName.BackToForum + '" class="btnGoToForum" onclick="onBtnGoToForumClick(' + parentId + ')" />';
			tr.appendChild(td);
			break;

		case "section":
			td = document.createElement("TD");
			td.innerHTML = '<span class="spacerA">&nbsp;</span>' +
				'<input type="button" value="' + buttonName.BackToSection + '" class="btnGoToSection" onclick="onBtnGoToSectionClick(' + parentId + ')" />';
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

function drawSectionsAndForums(p, nodes, forumsMap) {
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

function drawForumAndThreads(p, forum, threadsMap) {
	drawPageTitle(p, forum.forumName);
	addPaginator(p, mca_gvc.Page, mca_gvc.Pages, "forumPagePrev", "forumPageNext");
	drawForumThreads(p, forum, threadsMap);
}

function drawForumName(p, forum) {
	let divForum = document.createElement("DIV");
	divForum.className = "forum";
	divForum.id = "forum_" + forum.forumId;
	let ml = sectionMarginDelta;
	divForum.style.cssText = "margin-left: " + ml + "px";
	let url = composeUrlForForum(forum.forumId);
	divForum.innerHTML = "<a href='" + url + "'>" + forum.forumName + "</a>";
	p.appendChild(divForum);
}

function drawForumThreads(p, forum, threadsMap) {
	let divThread, ml, url, threadId, thread;
	let threadIds = forum.threadIds;
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

async function drawThreadAndMessages(p, thread, messagesMap) {
	drawPageTitle(p, thread.threadName);
	addPaginator(p, mca_gvc.Page, mca_gvc.Pages, "threadPagePrev", "threadPageNext");
	await drawThreadMessages(p, thread, messagesMap);
}

function drawThreadName(p, thread) {
	let divThread = document.createElement("DIV");
	divThread.className = "thread";
	divThread.id = "thread_" + thread.threadId;
	let ml = sectionMarginDelta;
	divThread.style.cssText = "margin-left: " + ml + "px";
	let url = composeUrlForThread(thread.threadId);
	divThread.innerHTML = "<a href='" + url + "'>" + thread.threadName + "</a>";
	p.appendChild(divThread);
}

async function drawThreadMessages(p, thread, messagesMap) {
	let messageIds = thread.messageIds;
	if (messageIds == null) {
		return
	}

	let divMsgHdr, divMsgBody, ml, messageId, message, txt;
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
	btnPrev.value = "<";
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
	btnNext.value = ">";
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

		case "subscriptionsPrev":
			btn.addEventListener("click", async (e) => {
				await onBtnPrevClick_subscriptionsPage(btn);
			});
			return;

		case "subscriptionsNext":
			btn.addEventListener("click", async (e) => {
				await onBtnNextClick_subscriptionsPage(btn);
			});
			return;

		case "notificationsPagePrev":
			btn.addEventListener("click", async (e) => {
				await onBtnPrevClick_notificationsPage(btn);
			});
			return;

		case "notificationsPageNext":
			btn.addEventListener("click", async (e) => {
				await onBtnNextClick_notificationsPage(btn);
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

async function onBtnPrevClick_subscriptionsPage(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(err.PreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = composeUrlForSubscriptionsPage(mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnNextClick_subscriptionsPage(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(err.NextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = composeUrlForSubscriptionsPage(mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnPrevClick_notificationsPage(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(err.PreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = composeUrlForNotificationsPage(mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnNextClick_notificationsPage(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(err.NextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = composeUrlForNotificationsPage(mca_gvc.Page);
	await redirectPage(false, url);
}

async function fillUserPage(userParams) {
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

	let resp = await countSelfSubscriptions();
	if (resp == null) {
		return;
	}
	let subscriptionsCount = resp.result.userSubscriptionsCount;
	document.getElementById(fi.id22).value = subscriptionsCount.toString();
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

async function onBtnManageSubscriptionsClick(btn) {
	let url = qp.Prefix + qpn.SubscriptionsPage;
	await redirectPage(false, url);
}

async function onBtnAccountClick(btn) {
	let url = qp.SelfPage;
	await redirectPage(false, url);
}

async function onBtnNotificationsClick(btn) {
	let url = qp.Prefix + qpn.Notifications;
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

async function getMessageCreatorName(userId) {
	return await mca_gvc.UNC.GetName(userId);
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
	tbl.appendChild(tr);
	let tdL = document.createElement("TD");
	tdL.id = "bottomActionPanelLeft";
	tdL.className = "bottomActionPanelLeft";
	tr.appendChild(tdL);
	let tdR = document.createElement("TD");
	tdR.id = "bottomActionPanelRight";
	tdR.className = "bottomActionPanelRight";
	tr.appendChild(tdR);
	div.appendChild(tbl);
	el.appendChild(div);

	switch (type) {
		case "forum":
			await drawBottomActionPanelButtonsForForum(objectId, object, tdL, tdR, userParams, cn);
			break;

		case "thread":
			await drawBottomActionPanelButtonsForThread(objectId, object, tdL, tdR, userParams, cn);
			break;

		case "message":
			await drawBottomActionPanelButtonsForMessage(objectId, object, tdL, tdR, userParams, cn);
			break;
	}
}

async function drawBottomActionPanelButtonsForForum(forumId, forum, tdL, tdR, userParams, panelClass) {
	if (userParams.isAuthor) {
		tdR.innerHTML = '<form><input type="button" value="' + buttonName.StartNewThread + '" class="btnStartNewThread" ' +
			'onclick="onBtnStartNewThreadClick(this, \'' + panelClass + '\', ' + forumId + ')" /></form>';
	}
}

async function drawBottomActionPanelButtonsForThread(threadId, thread, tdL, tdR, userParams, panelClass) {
	let resp = await getLatestMessageOfThread(threadId);
	if (resp == null) {
		return;
	}
	let latestMessageInThread = resp.result.message;

	resp = await isSelfSubscribed(threadId);
	if (resp == null) {
		return;
	}
	let userId = resp.result.userId;
	let isUserSubscribed = resp.result.isSubscribed;

	let tbl = document.createElement("TABLE");
	let tr = document.createElement("TR");

	if (!isUserSubscribed) {
		let td = document.createElement("TD");
		td.innerHTML += '<form><input type="button" value="' + buttonName.SubscribeToThread + '" class="btnSubscribe" ' +
			'onclick="onBtnSubscribeClick(this, ' + threadId + ', ' + userId + ')" /></form>';
		tr.appendChild(td);
	}

	let canAddMsg = canUserAddMessage(userParams, latestMessageInThread);
	if (canAddMsg) {
		let td = document.createElement("TD");
		td.innerHTML += '<form><input type="button" value="' + buttonName.AddMessage + '" class="btnAddMessage" ' +
			'onclick="onBtnAddMessageClick(this, \'' + panelClass + '\', ' + threadId + ')" /></form>';
		tr.appendChild(td);
	}

	tbl.appendChild(tr);
	tdR.appendChild(tbl);
}

async function drawBottomActionPanelButtonsForMessage(messageId, message, tdL, tdR, userParams, panelClass) {
	if (userParams.isAuthor) {
		let canEditMsg = canUserEditMessage(userParams, message);
		if (canEditMsg) {
			tdR.innerHTML = '<form><input type="button" value="' + buttonName.EditMessage + '" class="btnEditMessage" ' +
				'onclick="onBtnEditMessageClick(this, \'' + panelClass + '\', ' + messageId + ')" /></form>';
		}
	}
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

async function onBtnSubscribeClick(btn, threadId, userId) {
	let resp = await addSubscription(threadId, userId);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableButton(btn);
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

function showActionSuccess(btn, txt) {
	let ppp = btn.parentNode.parentNode.parentNode;
	let d = document.createElement("DIV");
	d.className = "actionSuccess";
	d.textContent = txt;
	ppp.appendChild(d);
}

function processMessageText(msgText) {
	let txt = escapeHtml(msgText);
	txt = txt.replaceAll("\r\n", '<br>');
	txt = txt.replaceAll("\n", '<br>');
	txt = txt.replaceAll("\r", '<br>');
	return txt;
}

async function showPage_Notifications() {
	// Prepare data.
	let pageNumber = mca_gvc.Page;
	let resp = await getNotificationsOnPage(pageNumber);
	if (resp == null) {
		return;
	}
	let pageCount = resp.result.nop.totalPages;
	pageCount = repairUndefinedPageCount(pageCount);
	mca_gvc.Pages = pageCount;

	// Check page number for overflow.
	if (pageNumber > pageCount) {
		console.error(err.PageNotFound);
		return;
	}

	let nop = resp.result.nop;
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divBB");
	showBlock(p);
	await addPageHead(p, settings.SiteName, false);
	addActionPanel(p, false);
	drawNotificationsOnPage(p, nop);
}

function drawNotificationsOnPage(p, nop) {
	drawPageTitle(p, 'Notifications');
	addPaginator(p, mca_gvc.Page, mca_gvc.Pages, "notificationsPagePrev", "notificationsPageNext");
	drawNotifications(p, 'notificationList', nop);
}

function drawNotifications(p, elClass, nop) {
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
	p.appendChild(tbl);

	let columnsWithHtml = [2, 3];

	// Cells.
	let notification, actions;
	for (let i = 0; i < nop.notifications.length; i++) {
		notification = nop.notifications[i];

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
			} else if (j === jLast) {
				td.className = "lastCol";
			} else {
				if (!notification.isRead) {
					td.className = "unread";
				}
				if (j === 1) {
					td.className += " col2";
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

function drawPageTitle(p, title) {
	let d = document.createElement("DIV");
	d.className = "pageTitle";
	let ml = sectionMarginDelta;
	d.style.cssText = "margin-left: " + ml + "px";
	d.textContent = title;
	p.appendChild(d);
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

async function showPage_Subscriptions() {
	// Prepare data.
	let pageNumber = mca_gvc.Page;
	let sop = await getSelfSubscriptionsPaginated(pageNumber);
	let pageCount = sop.TotalPages;
	pageCount = repairUndefinedPageCount(pageCount);
	mca_gvc.Pages = pageCount;

	// Check page number for overflow.
	if (pageNumber > pageCount) {
		console.error(err.PageNotFound);
		return;
	}

	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divBB");
	showBlock(p);
	await addPageHead(p, settings.SiteName, false);
	addActionPanel(p, false);
	drawSubscriptionsOnPage(p, sop);
}

function drawSubscriptionsOnPage(p, sop) {
	drawPageTitle(p, 'Subscriptions');
	addPaginator(p, mca_gvc.Page, mca_gvc.Pages, "subscriptionsPrev", "subscriptionsNext");
	drawSubscriptions(p, 'subscriptionsList', sop);
}

function drawSubscriptions(p, elClass, sop) {
	let tbl = document.createElement("TABLE");
	tbl.className = elClass;

	// Header.
	let tr = document.createElement("TR");
	let ths = ["#", "Thread ID", "Thread Name", "Actions"];
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
	p.appendChild(tbl);

	let columnsWithLink = [1, 2];
	let columnsWithHtml = [3];

	if (sop.ThreadIds == null) {
		return;
	}

	// Cells.
	let threadId, threadName;
	for (let i = 0; i < sop.ThreadIds.length; i++) {
		threadId = sop.ThreadIds[i];
		threadName = sop.ThreadNames[i];

		// Fill data.
		tr = document.createElement("TR");
		let tds = [];
		for (let j = 0; j < ths.length; j++) {
			tds.push("");
		}

		tds[0] = (i + 1).toString();
		tds[1] = threadId.toString();
		tds[2] = threadName;
		tds[3] = '<input type="button" class="btnUnsubscribe" value="Unsubscribe" onclick="onBtnUnsubscribeClick(this)">';

		let td, url;
		for (let j = 0; j < tds.length; j++) {
			url = composeUrlForThread(threadId);
			td = document.createElement("TD");

			if (j === 0) {
				td.className = "numCol";
			}

			if (columnsWithLink.includes(j)) {
				td.innerHTML = '<a href="' + url + '">' + tds[j] + '</a>';
			} else {
				if (columnsWithHtml.includes(j)) {
					td.innerHTML = tds[j];
				} else {
					td.textContent = tds[j];
				}
			}
			tr.appendChild(td);
		}

		tbl.appendChild(tr);
	}
}

async function onBtnUnsubscribeClick(btn) {
	let tr = btn.parentElement.parentElement;
	let threadId = Number(tr.children[1].textContent);
	let resp = await deleteSelfSubscription(threadId);
	if (resp == null) {
		return;
	}
	if (!resp.result.ok) {
		return;
	}
	tr.style.display = "none";
}

function repairUndefinedPageCount(pageCount) {
	// Unfortunately JavaScript can compare a number with 'undefined' !
	if (pageCount === undefined) {
		return 1;
	}
	if (pageCount === 0) {
		return 1;
	}
	return pageCount;
}

//TODO: Message model.

function getMessageLastTouchTime(message) {
	if (message.editor.time == null) {
		return new Date(message.creator.time);
	}
	return new Date(message.editor.time);
}

function getMessageMaxEditTime(message, settings) {
	let lastTouchTime = getMessageLastTouchTime(message);
	return addTimeSec(lastTouchTime, settings.MessageEditTime);
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
