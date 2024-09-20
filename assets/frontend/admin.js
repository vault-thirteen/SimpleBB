// Settings.
settingsPath = "settings.json";
rootPath = "/";
let settings;
redirectDelay = 0;

// Pages.
adminPage = "admin.html";
qpListOfUsers = "?listOfUsers";
qpRegistrationsReadyForApproval = "?registrationsReadyForApproval";

// Action names.
actionName_GetListOfAllUsers = "getListOfAllUsers";
actionName_GetListOfRegistrationsReadyForApproval = "getListOfRegistrationsReadyForApproval";
actionName_ApproveAndRegisterUser = "approveAndRegisterUser";
actionName_RejectRegistrationRequest = "rejectRegistrationRequest";
actionName_GetListOfLoggedUsers = "getListOfLoggedUsers";
actionName_ViewUserParameters = "viewUserParameters";

// Messages.
msgGenericErrorPrefix = "Error: ";

// Errors.
errSettings = "settings error";
errNotOk = "something went wrong";
errServer = "server error";
errClient = "client error";
errUnknown = "unknown error";
errPreviousPageDoesNotExist = "previous page does not exist";
errNextPageDoesNotExist = "next page does not exist";

// Global variables.
page = 0;
pages = 0;

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

class ApiResponse {
	constructor(isOk, jsonObject, statusCode, errorText) {
		this.IsOk = isOk;
		this.JsonObject = jsonObject;
		this.StatusCode = statusCode;
		this.ErrorText = errorText;
	}
}

async function onPageLoad() {
	console.debug("onPageLoad");

	settings = await getSettings();
	if (settings === null) {
		return;
	}
	console.debug('Received settings. Version:', settings.Version);

	// Select a page.
	let curPage = window.location.search;
	switch (curPage) {
		case qpListOfUsers:
			showPage_ListOfUsers();
			return;

		case qpRegistrationsReadyForApproval:
			showPage_RegistrationsReadyForApproval();
			return;

		default:
			showPage_MainMenu();
			return;
	}
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
	console.debug("onGoRegApprovalClick");
	await redirectToSubPage(true, qpRegistrationsReadyForApproval);
}

function onGoLoggedUsersClick(btn) {
	console.debug("onGoLoggedUsersClick");
}

async function onGoListAllUsersClick(btn) {
	console.debug("onGoListAllUsersClick");
	await redirectToSubPage(true, qpListOfUsers);
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
	let sp = document.getElementById("subpage");
	sp.style.display = "block";
	addBtnBack(sp);
	addTitle(sp, "List of All Users");
	page = 1;
	let resp = await getListOfAllUsers(page);
	if (resp == null) {
		return;
	}
	let pageNumber = resp.result.page;
	let pageCount = resp.result.totalPages;
	pages = pageCount;
	addPaginator(sp, pageNumber, pageCount, "userListPrev", "userListNext");
	addListOfUsers(sp);
	refreshListOfUsers("subpageListOfUsers", resp.result.userIds);
}

async function showPage_RegistrationsReadyForApproval() {
	let sp = document.getElementById("subpage");
	sp.style.display = "block";
	addBtnBack(sp);
	addTitle(sp, "List of Registrations ready for Approval");
	page = 1;
	let resp = await getListOfRegistrationsReadyForApproval(page);
	if (resp == null) {
		return;
	}
	let pageNumber = resp.result.page;
	let pageCount = resp.result.totalPages;
	pages = pageCount;
	addPaginator(sp, pageNumber, pageCount, "rrfaListPrev", "rrfaListNext");
	addListOfRRFA(sp);
	refreshListOfRRFA("subpageListOfRRFA", resp.result.rrfa);
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
	s.textContent = "Page " + pageNumber + " of " + pageCount + ". ";
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
				onBtnPrevClick(btn);
			});
			return;

		case "userListNext":
			btn.addEventListener("click", async (e) => {
				onBtnNextClick(btn);
			});
			return;

		case "rrfaListPrev":
			btn.addEventListener("click", async (e) => {
				onBtnPrevClick_rrfa(btn);
			});
			return;

		case "rrfaListNext":
			btn.addEventListener("click", async (e) => {
				onBtnNextClick_rrfa(btn);
			});
			return;

		default:
			console.error(x);
	}
}

function refreshPaginator(id, pageNumber, pageCount) {
	let div = document.getElementById(id);
	let ch1 = div.children[0];
	ch1.textContent = "Page " + pageNumber + " of " + pageCount + ". ";
}

async function onBtnPrevClick(btn) {
	if (page <= 1) {
		console.error(errPreviousPageDoesNotExist);
		return;
	}

	page--;
	await updateTableOfUsers();
}

async function onBtnNextClick(btn) {
	if (page >= pages) {
		console.error(errNextPageDoesNotExist);
		return;
	}

	page++;
	await updateTableOfUsers();
}

async function onBtnPrevClick_rrfa(btn) {
	if (page <= 1) {
		console.error(errPreviousPageDoesNotExist);
		return;
	}

	page--;
	await updateTableOfRRFA();
}

async function onBtnNextClick_rrfa(btn) {
	if (page >= pages) {
		console.error(errNextPageDoesNotExist);
		return;
	}

	page++;
	await updateTableOfRRFA();
}

async function updateTableOfUsers() {
	let resp = await getListOfAllUsers(page);
	if (resp == null) {
		return;
	}
	let pageNumber = resp.result.page;
	let pageCount = resp.result.totalPages;
	refreshPaginator("subpagePaginator", pageNumber, pageCount);
	refreshListOfUsers("subpageListOfUsers", resp.result.userIds);
}

async function updateTableOfRRFA() {
	let resp = await getListOfRegistrationsReadyForApproval(page);
	if (resp == null) {
		return;
	}
	let pageNumber = resp.result.page;
	let pageCount = resp.result.totalPages;
	refreshPaginator("subpagePaginator", pageNumber, pageCount);
	refreshListOfRRFA("subpageListOfRRFA", resp.result.rrfa);
}

function addListOfUsers(el) {
	let div = document.createElement("DIV");
	div.className = "subpageListOfUsers";
	div.id = "subpageListOfUsers";
	el.appendChild(div);
}

function addListOfRRFA(el) {
	let div = document.createElement("DIV");
	div.className = "subpageListOfRRFA";
	div.id = "subpageListOfRRFA";
	el.appendChild(div);
}

async function refreshListOfUsers(id, userIds) {
	let div = document.getElementById(id);
	div.innerHTML = "";
	let tbl = document.createElement("TABLE");
	tbl.className = "subpageListOfUsers";

	// Header.
	let tr = document.createElement("TR");
	let th;
	let ths = [
		"#", "ID", "IsLoggedIn", "E-Mail", "Name", "RegTime", "ApprovalTime",
		"LastBadLogInTime", "LastBadActionTime", "BanTime", "CanLogIn",
		"IsReader", "IsWriter", "IsAuthor", "IsModerator", "IsAdministrator"];
	for (i = 0; i < ths.length; i++) {
		th = document.createElement("TH");
		if (i === 0) {
			th.className = "numCol";
		}
		th.textContent = ths[i];
		tr.appendChild(th);
	}
	tbl.appendChild(tr);

	// Get list of logged-in users.
	let resp = await getListOfLoggedInUsers();
	if (resp == null) {
		return;
	}
	let loggedUserIds = resp.result.loggedUserIds;

	// Cells.
	let isUserLoggedIn;
	let userId;
	let userParams;
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

		let td;
		for (j = 0; j < tds.length; j++) {
			td = document.createElement("TD");

			if (j === 0) {
				td.className = "numCol";
			}

			td.textContent = tds[j];
			tr.appendChild(td);
		}

		tbl.appendChild(tr);
	}

	div.appendChild(tbl);
}

async function refreshListOfRRFA(id, rrfas) {
	let div = document.getElementById(id);
	div.innerHTML = "";
	let tbl = document.createElement("TABLE");
	tbl.className = "subpageListOfRRFA";

	// Header.
	let tr = document.createElement("TR");
	let th;
	let ths = [
		"#", "ID", "PreRegTime", "E-Mail", "Name", "Actions"];
	for (i = 0; i < ths.length; i++) {
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
		for (j = 0; j < tds.length; j++) {
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
	let resp;
	let result;
	let ri = {
		method: "POST",
		body: JSON.stringify(data)
	};
	resp = await fetch(rootPath + settings.ApiFolder, ri);
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

	return t.getUTCDay().toString().padStart(2, '0') + "." +
		t.getUTCMonth().toString().padStart(2, '0') + "." +
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
