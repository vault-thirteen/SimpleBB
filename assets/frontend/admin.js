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

// Settings.
settingsPath = "settings.json";
rootPath = "/";
redirectDelay = 0;
adminPage = "admin.html";

// Names of Query Parameters.
qpPrefix = "?"
qpnId = "id";
qpnPage = "page";
qpnListOfUsers = "listOfUsers";
qpnListOfLoggedUsers = "listOfLoggedUsers"
qpnRegistrationsReadyForApproval = "registrationsReadyForApproval";
qpnUserPage = "userPage";

// Action names.
actionName_GetListOfAllUsers = "getListOfAllUsers";
actionName_GetListOfRegistrationsReadyForApproval = "getListOfRegistrationsReadyForApproval";
actionName_ApproveAndRegisterUser = "approveAndRegisterUser";
actionName_RejectRegistrationRequest = "rejectRegistrationRequest";
actionName_GetListOfLoggedUsers = "getListOfLoggedUsers";
actionName_ViewUserParameters = "viewUserParameters";
actionName_LogUserOutA = "logUserOutA";
actionName_IsUserLoggedIn = "isUserLoggedIn";
actionName_SetUserRoleAuthor = "setUserRoleAuthor";
actionName_SetUserRoleWriter = "setUserRoleWriter";
actionName_SetUserRoleReader = "setUserRoleReader";
actionName_BanUser = "banUser";
actionName_UnbanUser = "unbanUser";
actionName_GetUserSession = "getUserSession";

// Messages.
msgGenericErrorPrefix = "Error: ";

// Errors.
errIdNotSet = "ID is not set";
errIdNotFound = "ID is not found";
errPageNotSet = "page is not set";
errPageNotFound = "page is not found";
errSettings = "settings error";
errNotOk = "something went wrong";
errServer = "server error";
errClient = "client error";
errUnknown = "unknown error";
errPreviousPageDoesNotExist = "previous page does not exist";
errNextPageDoesNotExist = "next page does not exist";
errUnknownVariant = "unknown variant";

// User role names.
userRoleAuthor = "author";
userRoleWriter = "writer";
userRoleReader = "reader";
userRoleLogging = "logging";

// Global variables.
class GlobalVariablesContainer {
	constructor(settings, id, page, pages) {
		this.Settings = settings;
		this.Id = id;
		this.Page = page;
		this.Pages = pages;
	}
}

mca_gvc = new GlobalVariablesContainer(0, 0, 0);

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

class Parameters_GetListOfAllUsers {
	constructor(page) {
		this.Page = page;
	}
}

class Parameters_GetListOfRegistrationsReadyForApproval {
	constructor(page) {
		this.Page = page;
	}
}

class Parameters_ViewUserParameters {
	constructor(userId) {
		this.UserId = userId;
	}
}

class Parameters_ApproveAndRegisterUser {
	constructor(email) {
		this.Email = email;
	}
}

class Parameters_RejectRegistrationRequest {
	constructor(id) {
		this.Id = id;
	}
}

class Parameters_LogUserOutA {
	constructor(userId) {
		this.UserId = userId;
	}
}

class Parameters_IsUserLoggedIn {
	constructor(userId) {
		this.UserId = userId;
	}
}

class Parameters_SetUserRoleAuthor {
	constructor(userId, isRoleEnabled) {
		this.UserId = userId;
		this.IsRoleEnabled = isRoleEnabled;
	}
}

class Parameters_SetUserRoleWriter {
	constructor(userId, isRoleEnabled) {
		this.UserId = userId;
		this.IsRoleEnabled = isRoleEnabled;
	}
}

class Parameters_SetUserRoleReader {
	constructor(userId, isRoleEnabled) {
		this.UserId = userId;
		this.IsRoleEnabled = isRoleEnabled;
	}
}

class Parameters_BanUser {
	constructor(userId) {
		this.UserId = userId;
	}
}

class Parameters_UnbanUser {
	constructor(userId) {
		this.UserId = userId;
	}
}

class Parameters_GetUserSession {
	constructor(userId) {
		this.UserId = userId;
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

async function onPageLoad() {
	let settings = await getSettings();
	if (settings === null) {
		return;
	}
	mca_gvc.Settings = settings;
	console.info('Received settings. Version: ' + settings.Version.toString() + ".");

	// Select a page.
	let curPage = window.location.search;
	let sp = new URLSearchParams(curPage);

	if (sp.has(qpnListOfUsers)) {
		if (!preparePageVariable(sp)) {
			return;
		}
		await showPage_ListOfUsers();
		return;
	}

	if (sp.has(qpnListOfLoggedUsers)) {
		if (!preparePageVariable(sp)) {
			return;
		}
		await showPage_ListOfLoggedUsers();
		return;
	}

	if (sp.has(qpnRegistrationsReadyForApproval)) {
		if (!preparePageVariable(sp)) {
			return;
		}
		await showPage_RegistrationsReadyForApproval();
		return;
	}

	if (sp.has(qpnUserPage)) {
		if (!prepareIdVariable(sp)) {
			return;
		}
		await showPage_UserPage();
		return;
	}

	showPage_MainMenu();
}

async function getSettings() {
	let data = await fetchSettings();

	let x = new Settings(
		data.version,
		data.productVersion,
		data.siteName,
		data.siteDomain,
		data.captchaFolder,
		data.sessionMaxDuration,
		data.messageEditTime,
		data.pageSize,
		data.apiFolder,
		data.publicSettingsFileName,
		data.isFrontEndEnabled,
		data.frontEndStaticFilesFolder,
	);

	// Self-check.
	if ((x.PublicSettingsFileName !== settingsPath)) {
		console.error(errSettings);
		return null;
	}

	return x;
}

async function fetchSettings() {
	let data = await fetch(rootPath + settingsPath);
	return await data.json();
}

async function onGoRegApprovalClick(btn) {
	await redirectToSubPage(true, qpPrefix + qpnRegistrationsReadyForApproval);
}

async function onGoLoggedUsersClick(btn) {
	await redirectToSubPage(true, qpPrefix + qpnListOfLoggedUsers);
}

async function onGoListAllUsersClick(btn) {
	await redirectToSubPage(true, qpPrefix + qpnListOfUsers);
}

async function redirectToSubPage(wait, qp) {
	let url = adminPage + qp;
	await redirectPage(wait, url);
}

async function redirectToMainMenu(wait) {
	let url = adminPage;
	await redirectPage(wait, url);
}

async function redirectPage(wait, url) {
	if (wait) {
		await sleep(redirectDelay * 1000);
	}

	document.location.href = url;
}

async function sleep(ms) {
	await new Promise(r => setTimeout(r, ms));
}

function showPage_MainMenu() {
	document.getElementById("acpMenu").style.display = "table";
}

async function showPage_ListOfUsers() {
	let pageNumber = mca_gvc.Page;
	let resp = await getListOfAllUsers(pageNumber);
	if (resp == null) {
		return;
	}
	let pageCount = resp.result.totalPages;
	mca_gvc.Pages = pageCount;

	// Check page number for overflow.
	if (pageNumber > pageCount) {
		console.error(errPageNotFound);
		return;
	}

	let userIds = resp.result.userIds;

	// Draw.
	let p = document.getElementById("subpage");
	p.style.display = "block";
	addBtnBack(p);
	addTitle(p, "List of All Users");
	addPaginator(p, pageNumber, pageCount, "userListPrev", "userListNext");
	addDiv(p, "subpageListOfUsers");
	await fillListOfUsers("subpageListOfUsers", userIds);
}

async function showPage_ListOfLoggedUsers() {
	let pageNumber = mca_gvc.Page;
	let resp = await getListOfLoggedInUsers();
	if (resp == null) {
		return;
	}
	let userIds = resp.result.loggedUserIds;
	let userCount = userIds.length;
	let pageCount = Math.ceil(userCount / mca_gvc.Settings.PageSize);
	mca_gvc.Pages = pageCount;

	// Check page number for overflow.
	if (pageNumber > pageCount) {
		console.error(errPageNotFound);
		return;
	}

	let userIdsOnPage = calculateUserIdsOnPage(userIds, pageNumber, mca_gvc.Settings.PageSize);

	// Draw.
	let p = document.getElementById("subpage");
	p.style.display = "block";
	addBtnBack(p);
	addTitle(p, "List of logged-in Users");
	addPaginator(p, pageNumber, pageCount, "loggedUserListPrev", "loggedUserListNext");
	addDiv(p, "subpageListOfLoggedUsers");
	await fillListOfLoggedUsers("subpageListOfLoggedUsers", userIdsOnPage);
}

async function showPage_RegistrationsReadyForApproval() {
	let pageNumber = mca_gvc.Page;
	let resp = await getListOfRegistrationsReadyForApproval(mca_gvc.Page);
	if (resp == null) {
		return;
	}
	let pageCount = resp.result.totalPages;
	mca_gvc.Pages = pageCount;

	// Check page number for overflow.
	if (pageNumber > pageCount) {
		console.error(errPageNotFound);
		return;
	}

	let rrfas = resp.result.rrfa;

	// Draw.
	let p = document.getElementById("subpage");
	p.style.display = "block";
	addBtnBack(p);
	addTitle(p, "List of Registrations ready for Approval");
	addPaginator(p, pageNumber, pageCount, "rrfaListPrev", "rrfaListNext");
	addDiv(p, "subpageListOfRRFA");
	await fillListOfRRFA("subpageListOfRRFA", rrfas);
}

async function showPage_UserPage() {
	let userId = mca_gvc.Id;
	let resp = await viewUserParameters(userId);
	if (resp == null) {
		return;
	}
	let userParams = resp.result;

	// Get additional information.
	resp = await isUserLoggedIn(userId);
	if (resp == null) {
		return;
	}
	if (resp.result.userId !== userId) {
		return;
	}
	let userLogInState = resp.result.isUserLoggedIn;

	// Draw.
	let p = document.getElementById("subpage");
	p.style.display = "block";
	addBtnBack(p);
	addTitle(p, "User Page");
	addDiv(p, "subpageUserPage");
	fillUserPage("subpageUserPage", userParams, userLogInState);
}

function addBtnBack(el) {
	let btn = document.createElement("INPUT");
	btn.type = "button";
	btn.className = "btnBack";
	btn.value = "Go Back";
	btn.addEventListener("click", async (e) => {
		await redirectToMainMenu(true);
	})
	el.appendChild(btn);
}

function addTitle(el, text) {
	let div = document.createElement("DIV");
	div.className = "subpageTitle";
	div.id = "subpageTitle";
	div.textContent = text;
	el.appendChild(div);
}

function addPaginator(el, pageNumber, pageCount, variantPrev, variantNext) {
	let div = document.createElement("DIV");
	div.className = "subpagePaginator";
	div.id = "subpagePaginator";

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
	s.className = "subpageSpacerA";
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
		case "userListPrev":
			btn.addEventListener("click", async (e) => {
				await onBtnPrevClick_userList(btn);
			});
			return;

		case "userListNext":
			btn.addEventListener("click", async (e) => {
				await onBtnNextClick_userList(btn);
			});
			return;

		case "loggedUserListPrev":
			btn.addEventListener("click", async (e) => {
				await onBtnPrevClick_logged(btn);
			});
			return;

		case "loggedUserListNext":
			btn.addEventListener("click", async (e) => {
				await onBtnNextClick_logged(btn);
			});
			return;

		case "rrfaListPrev":
			btn.addEventListener("click", async (e) => {
				await onBtnPrevClick_rrfa(btn);
			});
			return;

		case "rrfaListNext":
			btn.addEventListener("click", async (e) => {
				await onBtnNextClick_rrfa(btn);
			});
			return;

		default:
			console.error(errUnknownVariant);
	}
}

async function onBtnPrevClick_userList(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(errPreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = qpPrefix + qpnListOfUsers + "&" + qpnPage + "=" + mca_gvc.Page;
	await redirectPage(false, url);
}

async function onBtnNextClick_userList(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(errNextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = qpPrefix + qpnListOfUsers + "&" + qpnPage + "=" + mca_gvc.Page;
	await redirectPage(false, url);
}

async function onBtnPrevClick_logged(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(errPreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = qpPrefix + qpnListOfLoggedUsers + "&" + qpnPage + "=" + mca_gvc.Page;
	await redirectPage(false, url);
}

async function onBtnNextClick_logged(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(errNextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = qpPrefix + qpnListOfLoggedUsers + "&" + qpnPage + "=" + mca_gvc.Page;
	await redirectPage(false, url);
}

async function onBtnPrevClick_rrfa(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(errPreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = qpPrefix + qpnRegistrationsReadyForApproval + "&" + qpnPage + "=" + mca_gvc.Page;
	await redirectPage(false, url);
}

async function onBtnNextClick_rrfa(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(errNextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = qpPrefix + qpnRegistrationsReadyForApproval + "&" + qpnPage + "=" + mca_gvc.Page;
	await redirectPage(false, url);
}

function addDiv(el, x) {
	let div = document.createElement("DIV");
	div.className = x;
	div.id = x;
	el.appendChild(div);
}

async function fillListOfLoggedUsers(elClass, userIdsOnPage) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let tbl = document.createElement("TABLE");
	tbl.className = elClass;

	// Header.
	let tr = document.createElement("TR");
	let ths = [
		"#", "ID", "E-Mail", "Name", "IP Address", "Log Time", "Actions"];
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

	let columnsWithLink = [1, 2, 3];

	// Cells.
	let userId, userParams, userSession, resp;
	for (let i = 0; i < userIdsOnPage.length; i++) {
		userId = userIdsOnPage[i];

		// Get user parameters.
		resp = await viewUserParameters(userId);
		if (resp == null) {
			return;
		}
		userParams = resp.result;

		// Get user session.
		resp = await getUserSession(userId);
		if (resp == null) {
			return;
		}
		userSession = resp.result.session;

		// Fill data.
		tr = document.createElement("TR");
		let tds = [];
		for (let j = 0; j < ths.length; j++) {
			tds.push("");
		}

		tds[0] = (i + 1).toString();
		tds[1] = userId.toString();
		tds[2] = userParams.email;
		tds[3] = userParams.name;
		tds[4] = userSession.userIPA;
		tds[5] = prettyTime(userSession.startTime);
		tds[6] = '<input type="button" class="btnLogOut" value="Log Out" onclick="onBtnLogOutClick(this)">';

		let td, url;
		for (let j = 0; j < tds.length; j++) {
			url = composeUserPageLink(userId.toString());
			td = document.createElement("TD");

			if (j === 0) {
				td.className = "numCol";
			}

			if (columnsWithLink.includes(j)) {
				td.innerHTML = '<a href="' + url + '">' + tds[j] + '</a>';
			} else if (j === 6) {
				td.innerHTML = tds[j];
			} else {
				td.textContent = tds[j];
			}
			tr.appendChild(td);
		}

		tbl.appendChild(tr);
	}

	div.appendChild(tbl);
}

async function fillListOfRRFA(elClass, rrfas) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let tbl = document.createElement("TABLE");
	tbl.className = elClass;

	// Header.
	let tr = document.createElement("TR");
	let ths = [
		"#", "ID", "PreRegTime", "E-Mail", "Name", "Actions"];
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

	// Cells.
	let rrfa;
	for (let i = 0; i < rrfas.length; i++) {
		rrfa = rrfas[i];

		// Fill data.
		tr = document.createElement("TR");
		let tds = [];
		for (let j = 0; j < ths.length; j++) {
			tds.push("");
		}

		tds[0] = (i + 1).toString();
		tds[1] = rrfa.id.toString();
		tds[2] = prettyTime(rrfa.preRegTime);
		tds[3] = rrfa.email;
		tds[4] = rrfa.name;
		tds[5] = '<input type="button" class="btnAccept" value="Accept" onclick="onBtnAcceptClick(this)">' +
			'<span class="subpageSpacerA">&nbsp;</span>' +
			'<input type="button" class="btnReject" value="Reject" onclick="onBtnRejectClick(this)">';

		let td;
		for (let j = 0; j < tds.length; j++) {
			td = document.createElement("TD");

			if (j === 0) {
				td.className = "numCol";
			}

			if (j !== 5) {
				td.textContent = tds[j];
			} else {
				td.innerHTML = tds[j];
			}
			tr.appendChild(td);
		}

		tbl.appendChild(tr);
	}

	div.appendChild(tbl);
}

async function sendApiRequest(data) {
	let result;
	let ri = {
		method: "POST",
		body: JSON.stringify(data)
	};
	let resp = await fetch(rootPath + mca_gvc.Settings.ApiFolder, ri);
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
	return msgGenericErrorPrefix + errMsg.trim() + ".";
}

async function getListOfAllUsers(pageN) {
	let params = new Parameters_GetListOfAllUsers(pageN);
	let reqData = new ApiRequest(actionName_GetListOfAllUsers, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function getListOfRegistrationsReadyForApproval(pageN) {
	let params = new Parameters_GetListOfRegistrationsReadyForApproval(pageN);
	let reqData = new ApiRequest(actionName_GetListOfRegistrationsReadyForApproval, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function getListOfLoggedInUsers() {
	let reqData = new ApiRequest(actionName_GetListOfLoggedUsers, {});
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function viewUserParameters(userId) {
	let params = new Parameters_ViewUserParameters(userId);
	let reqData = new ApiRequest(actionName_ViewUserParameters, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function isUserLoggedIn(userId) {
	let params = new Parameters_IsUserLoggedIn(userId);
	let reqData = new ApiRequest(actionName_IsUserLoggedIn, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

function boolToText(b) {
	if (b === true) {
		return "Yes";
	}
	if (b === false) {
		return "No";
	}
	console.error("boolToText:", b);
	return null;
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

async function onBtnAcceptClick(btn) {
	let tr = btn.parentElement.parentElement;
	let reqEmail = tr.children[3].textContent;
	let resp = await approveAndRegisterUser(reqEmail);
	if (resp == null) {
		return;
	}
	if (!resp.result.ok) {
		return;
	}
	tr.style.display = "none";
}

async function onBtnRejectClick(btn) {
	let tr = btn.parentElement.parentElement;
	let reqId = Number(tr.children[1].textContent);
	let resp = await rejectRegistrationRequest(reqId);
	if (resp == null) {
		return;
	}
	if (!resp.result.ok) {
		return;
	}
	tr.style.display = "none";
}

async function onBtnLogOutClick(btn) {
	let tr = btn.parentElement.parentElement;
	let userId = Number(tr.children[1].textContent);
	let resp = await logUserOutA(userId);
	if (resp == null) {
		return;
	}
	if (!resp.result.ok) {
		return;
	}
	tr.style.display = "none";
}

async function onBtnLogOutUPClick(userId) {
	let resp = await logUserOutA(userId);
	if (resp == null) {
		return;
	}
	if (!resp.result.ok) {
		return;
	}
	reloadPage();
}

async function onBtnEnableRoleUPClick(role, userId) {
	let resp;
	switch (role) {
		case userRoleAuthor:
			resp = await setUserRoleAuthor(userId, true);
			break;

		case userRoleWriter:
			resp = await setUserRoleWriter(userId, true);
			break;

		case userRoleReader:
			resp = await setUserRoleReader(userId, true);
			break;

		case userRoleLogging:
			resp = await unbanUser(userId);
			break;

		default:
			return;
	}

	if (resp == null) {
		return;
	}
	if (!resp.result.ok) {
		return;
	}
	reloadPage();
}

async function onBtnDisableRoleUPClick(role, userId) {
	let resp;
	switch (role) {
		case userRoleAuthor:
			resp = await setUserRoleAuthor(userId, false);
			break;

		case userRoleWriter:
			resp = await setUserRoleWriter(userId, false);
			break;

		case userRoleReader:
			resp = await setUserRoleReader(userId, false);
			break;

		case userRoleLogging:
			resp = await banUser(userId);
			break;

		default:
			return;
	}

	if (resp == null) {
		return;
	}
	if (!resp.result.ok) {
		return;
	}
	reloadPage();
}

async function approveAndRegisterUser(email) {
	let params = new Parameters_ApproveAndRegisterUser(email);
	let reqData = new ApiRequest(actionName_ApproveAndRegisterUser, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function rejectRegistrationRequest(id) {
	let params = new Parameters_RejectRegistrationRequest(id);
	let reqData = new ApiRequest(actionName_RejectRegistrationRequest, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

function calculateUserIdsOnPage(allIds, pageN, pageSize) {
	let x = Math.min(pageN * pageSize, allIds.length);
	return allIds.slice((pageN - 1) * pageSize, x);
}

async function logUserOutA(userId) {
	let params = new Parameters_LogUserOutA(userId);
	let reqData = new ApiRequest(actionName_LogUserOutA, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function setUserRoleAuthor(userId, roleValue) {
	let params = new Parameters_SetUserRoleAuthor(userId, roleValue);
	let reqData = new ApiRequest(actionName_SetUserRoleAuthor, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function setUserRoleWriter(userId, roleValue) {
	let params = new Parameters_SetUserRoleWriter(userId, roleValue);
	let reqData = new ApiRequest(actionName_SetUserRoleWriter, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function setUserRoleReader(userId, roleValue) {
	let params = new Parameters_SetUserRoleReader(userId, roleValue);
	let reqData = new ApiRequest(actionName_SetUserRoleReader, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function banUser(userId) {
	let params = new Parameters_BanUser(userId);
	let reqData = new ApiRequest(actionName_BanUser, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function unbanUser(userId) {
	let params = new Parameters_UnbanUser(userId);
	let reqData = new ApiRequest(actionName_UnbanUser, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function getUserSession(userId) {
	let params = new Parameters_GetUserSession(userId);
	let reqData = new ApiRequest(actionName_GetUserSession, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

function composeUserPageLink(userId) {
	return qpPrefix + qpnUserPage + "&" + qpnId + "=" + userId;
}

function reloadPage() {
	location.reload();
}

function prepareIdVariable(sp) {
	if (!sp.has(qpnId)) {
		console.error(errIdNotSet);
		return false;
	}

	let xId = Number(sp.get(qpnId));
	if (xId <= 0) {
		console.error(errIdNotFound);
		return false;
	}

	mca_gvc.Id = xId;
	return true;
}

function preparePageVariable(sp) {
	let pageNumber;
	if (!sp.has(qpnPage)) {
		pageNumber = 1;
	} else {
		pageNumber = Number(sp.get(qpnPage));
		if (pageNumber <= 0) {
			console.error(errPageNotFound);
			return false;
		}
	}

	mca_gvc.Page = pageNumber;
	return true;
}

async function fillListOfUsers(elClass, userIds) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let tbl = document.createElement("TABLE");
	tbl.className = elClass;

	// Header.
	let tr = document.createElement("TR");
	let ths = [
		"#", "ID", "IsLoggedIn", "E-Mail", "Name", "RegTime", "ApprovalTime",
		"LastBadLogInTime", "LastBadActionTime", "BanTime", "CanLogIn",
		"IsReader", "IsWriter", "IsAuthor", "IsModerator", "IsAdministrator"];
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

	// Get list of logged-in users (for flags).
	let resp = await getListOfLoggedInUsers();
	if (resp == null) {
		return;
	}
	let loggedUserIds = resp.result.loggedUserIds;
	let columnsWithLink = [1, 3, 4];

	// Cells.
	let isUserLoggedIn, userId, userParams;
	for (let i = 0; i < userIds.length; i++) {
		userId = userIds[i];

		// Get user parameters.
		resp = await viewUserParameters(userId);
		if (resp == null) {
			return;
		}
		userParams = resp.result;

		// Fill data.
		tr = document.createElement("TR");
		let tds = [];
		for (let j = 0; j < ths.length; j++) {
			tds.push("");
		}

		tds[0] = (i + 1).toString();
		tds[1] = userId.toString();
		isUserLoggedIn = loggedUserIds.includes(userId);
		tds[2] = boolToText(isUserLoggedIn);
		tds[3] = userParams.email;
		tds[4] = userParams.name;
		tds[5] = prettyTime(userParams.regTime);
		tds[6] = prettyTime(userParams.approvalTime);
		tds[7] = prettyTime(userParams.lastBadLogInTime);
		tds[8] = prettyTime(userParams.lastBadActionTime);
		tds[9] = prettyTime(userParams.banTime);
		tds[10] = boolToText(userParams.canLogIn);
		tds[11] = boolToText(userParams.isReader);
		tds[12] = boolToText(userParams.isWriter);
		tds[13] = boolToText(userParams.isAuthor);
		tds[14] = boolToText(userParams.isModerator);
		tds[15] = boolToText(userParams.isAdministrator);

		let td, url;
		for (let j = 0; j < tds.length; j++) {
			url = composeUserPageLink(userId.toString());
			td = document.createElement("TD");

			if (j === 0) {
				td.className = "numCol";
			}

			if (columnsWithLink.includes(j)) {
				td.innerHTML = '<a href="' + url + '">' + tds[j] + '</a>';
			} else {
				td.textContent = tds[j];
			}
			tr.appendChild(td);
		}

		tbl.appendChild(tr);
	}

	div.appendChild(tbl);
}

function fillUserPage(elClass, userParams, userLogInState) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let tbl = document.createElement("TABLE");
	tbl.className = elClass;

	// Header.
	let tr = document.createElement("TR");
	let ths = ["#", "Field Name", "Value", "Actions"];
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

	let userId = mca_gvc.Id;
	let fieldNames = [
		"ID", "E-Mail", "Name", "PreRegTime", "RegTime", "ApprovalTime",
		"IsAdministrator", "IsModerator", "IsAuthor", "IsWriter", "IsReader",
		"CanLogIn", "LastBadLogInTime", "BanTime", "LastBadActionTime",
		"IsLoggedIn"
	];
	let fieldValues = [
		userId.toString(),
		userParams.email,
		userParams.name,
		prettyTime(userParams.preRegTime),
		prettyTime(userParams.regTime),
		prettyTime(userParams.approvalTime),
		boolToText(userParams.isAdministrator),
		boolToText(userParams.isModerator),
		boolToText(userParams.isAuthor),
		boolToText(userParams.isWriter),
		boolToText(userParams.isReader),
		boolToText(userParams.canLogIn),
		prettyTime(userParams.lastBadLogInTime),
		prettyTime(userParams.banTime),
		prettyTime(userParams.lastBadActionTime),
		boolToText(userLogInState),
	];

	// Rows.
	let tds, td, actions;
	for (let i = 0; i < fieldNames.length; i++) {
		tr = document.createElement("TR");

		tds = [];
		for (let j = 0; j < ths.length; j++) {
			tds.push("");
		}

		tds[0] = (i + 1).toString();
		tds[1] = fieldNames[i];
		tds[2] = fieldValues[i];

		switch (fieldNames[i]) {
			case "IsAuthor":
				if (userParams.isAuthor) {
					actions = '<input type="button" class="btnDisableRole" value="Disable Role" ' +
						'onclick="onBtnDisableRoleUPClick(\'' + userRoleAuthor + '\',' + userId + ')">';
				} else {
					actions = '<input type="button" class="btnEnableRole" value="Enable Role" ' +
						'onclick="onBtnEnableRoleUPClick(\'' + userRoleAuthor + '\',' + userId + ')">';
				}
				break;

			case "IsWriter":
				if (userParams.isWriter) {
					actions = '<input type="button" class="btnDisableRole" value="Disable Role" ' +
						'onclick="onBtnDisableRoleUPClick(\'' + userRoleWriter + '\',' + userId + ')">';
				} else {
					actions = '<input type="button" class="btnEnableRole" value="Enable Role" ' +
						'onclick="onBtnEnableRoleUPClick(\'' + userRoleWriter + '\',' + userId + ')">';
				}
				break;

			case "IsReader":
				if (userParams.isReader) {
					actions = '<input type="button" class="btnDisableRole" value="Disable Role" ' +
						'onclick="onBtnDisableRoleUPClick(\'' + userRoleReader + '\',' + userId + ')">';
				} else {
					actions = '<input type="button" class="btnEnableRole" value="Enable Role" ' +
						'onclick="onBtnEnableRoleUPClick(\'' + userRoleReader + '\',' + userId + ')">';
				}
				break;

			case "CanLogIn":
				if (userParams.canLogIn) {
					actions = '<input type="button" class="btnDisableRole" value="Disable Role" ' +
						'onclick="onBtnDisableRoleUPClick(\'' + userRoleLogging + '\',' + userId + ')">';
				} else {
					actions = '<input type="button" class="btnEnableRole" value="Enable Role" ' +
						'onclick="onBtnEnableRoleUPClick(\'' + userRoleLogging + '\',' + userId + ')">';
				}
				break;

			case "IsLoggedIn":
				if (userLogInState) {
					actions = '<input type="button" class="btnLogOut" value="Log Out" onclick="onBtnLogOutUPClick(' + userId + ')">';
				} else {
					actions = "";
				}
				break;

			default:
				actions = "";
		}
		tds[3] = actions;

		let jLast = tds.length - 1;
		for (let j = 0; j < tds.length; j++) {
			td = document.createElement("TD");
			if (j === 0) {
				td.className = "numCol";
			}
			if (j === jLast) {
				td.innerHTML = tds[j];
			} else {
				td.textContent = tds[j];
			}
			tr.appendChild(td);
		}

		tbl.appendChild(tr);
	}

	div.appendChild(tbl);
}
