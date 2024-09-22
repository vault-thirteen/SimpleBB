// Settings.
settingsPath = "settings.json";
rootPath = "/";
let settings;
redirectDelay = 0;

// Names of Query Parameters.
qpPrefix = "?"
qpnId = "id";
qpnListOfUsers = "listOfUsers";
qpnListOfLoggedUsers = "listOfLoggedUsers"
qpnRegistrationsReadyForApproval = "registrationsReadyForApproval";
qpnUserPage = "userPage";

// Pages.
adminPage = "admin.html";
qpListOfUsers = qpPrefix + qpnListOfUsers;
qpListOfLoggedUsers = qpPrefix + qpnListOfLoggedUsers;
qpRegistrationsReadyForApproval = qpPrefix + qpnRegistrationsReadyForApproval;
qpUserPage = qpPrefix + qpnUserPage;

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
errSettings = "settings error";
errNotOk = "something went wrong";
errServer = "server error";
errClient = "client error";
errUnknown = "unknown error";
errPreviousPageDoesNotExist = "previous page does not exist";
errNextPageDoesNotExist = "next page does not exist";
errPageNotFound = "page is not found";
errIdNotSet = "ID is not set";

// User role names.
userRoleAuthor = "author";
userRoleWriter = "writer";
userRoleReader = "reader";
userRoleLogging = "logging";

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
	settings = await getSettings();
	if (settings === null) {
		return;
	}
	console.info('Received settings. Version:', settings.Version);

	// Select a page.
	let curPage = window.location.search;
	let sp = new URLSearchParams(curPage);

	switch (curPage) {
		case qpListOfUsers:
			showPage_ListOfUsers();
			return;

		case qpListOfLoggedUsers:
			showPage_ListOfLoggedUsers();
			return;

		case qpRegistrationsReadyForApproval:
			showPage_RegistrationsReadyForApproval();
			return;
	}

	if (sp.has(qpnUserPage)) {
		if (!sp.has(qpnId)) {
			console.error(errIdNotSet);
			return;
		}

		let userId = Number(sp.get(qpnId));
		showPage_UserPage(userId);
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
	await redirectToSubPage(true, qpRegistrationsReadyForApproval);
}

async function onGoLoggedUsersClick(btn) {
	await redirectToSubPage(true, qpListOfLoggedUsers);
}

async function onGoListAllUsersClick(btn) {
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

async function showPage_ListOfLoggedUsers() {
	let sp = document.getElementById("subpage");
	sp.style.display = "block";
	addBtnBack(sp);
	addTitle(sp, "List of logged-in Users");
	page = 1;
	let resp = await getListOfLoggedInUsers();
	if (resp == null) {
		return;
	}
	let userIds = resp.result.loggedUserIds;
	let userCount = userIds.length;
	let pageNumber = page;
	let pageCount = Math.ceil(userCount / settings.PageSize);
	pages = pageCount;
	let userIdsOnPage = calculateUserIdsOnPage(userIds, page, settings.PageSize);
	addPaginator(sp, pageNumber, pageCount, "loggedUserListPrev", "loggedUserListNext");
	addListOfLoggedUsers(sp);
	refreshListOfLoggedUsers("subpageListOfLoggedUsers", userIdsOnPage);
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

async function showPage_UserPage(userId) {
	let sp = document.getElementById("subpage");
	sp.style.display = "block";
	addBtnBack(sp);
	addTitle(sp, "User Page");

	// Get information.
	let resp = await viewUserParameters(userId);
	if (resp == null) {
		return;
	}
	let userParams = resp.result;

	resp = await isUserLoggedIn(userId);
	if (resp == null) {
		return;
	}
	if (resp.result.userId !== userId) {
		return;
	}
	let userLogInState = resp.result.isUserLoggedIn;

	// Draw the information.
	let divUserPage = document.createElement("DIV");
	divUserPage.className = "subpageUserPage";
	divUserPage.id = "subpageUserPage";
	sp.appendChild(divUserPage);
	let tbl = document.createElement("TABLE");
	tbl.className = "subpageUserPage";

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

	// Header.
	let tr = document.createElement("TR");
	let th;
	let ths = ["#", "Field Name", "Value", "Actions"];

	for (i = 0; i < ths.length; i++) {
		th = document.createElement("TH");
		if (i === 0) {
			th.className = "numCol";
		}
		th.textContent = ths[i];
		tr.appendChild(th);
	}
	tbl.appendChild(tr);

	// Rows.
	let tds;
	let td;
	let actions;
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
		for (j = 0; j < tds.length; j++) {
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

	divUserPage.appendChild(tbl);
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
				onBtnPrevClick(btn);
			});
			return;

		case "userListNext":
			btn.addEventListener("click", async (e) => {
				onBtnNextClick(btn);
			});
			return;

		case "loggedUserListPrev":
			btn.addEventListener("click", async (e) => {
				onBtnPrevClick_logged(btn);
			});
			return;

		case "loggedUserListNext":
			btn.addEventListener("click", async (e) => {
				onBtnNextClick_logged(btn);
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
	ch1.textContent = "Page " + pageNumber + " of " + pageCount + " ";
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

async function onBtnPrevClick_logged(btn) {
	if (page <= 1) {
		console.error(errPreviousPageDoesNotExist);
		return;
	}

	page--;
	await updateTableOfLoggedUsers();
}

async function onBtnNextClick_logged(btn) {
	if (page >= pages) {
		console.error(errNextPageDoesNotExist);
		return;
	}

	page++;
	await updateTableOfLoggedUsers();
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

async function updateTableOfLoggedUsers() {
	let resp = await getListOfLoggedInUsers();
	if (resp == null) {
		return;
	}
	let userIds = resp.result.loggedUserIds;
	let userCount = userIds.length;
	let pageNumber = page;
	let pageCount = Math.ceil(userCount / settings.PageSize);
	pages = pageCount;
	let userIdsOnPage = calculateUserIdsOnPage(userIds, page, settings.PageSize);
	refreshPaginator("subpagePaginator", pageNumber, pageCount);
	refreshListOfLoggedUsers("subpageListOfLoggedUsers", userIdsOnPage);
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

function addListOfLoggedUsers(el) {
	let div = document.createElement("DIV");
	div.className = "subpageListOfLoggedUsers";
	div.id = "subpageListOfLoggedUsers";
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
	let columnsWithLink = [1, 3, 4];

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
		let url;
		for (j = 0; j < tds.length; j++) {
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

async function refreshListOfLoggedUsers(id, userIdsOnPage) {
	let div = document.getElementById(id);
	div.innerHTML = "";
	let tbl = document.createElement("TABLE");
	tbl.className = "subpageListOfLoggedUsers";

	// Header.
	let tr = document.createElement("TR");
	let th;
	let ths = [
		"#", "ID", "E-Mail", "Name", "IP Address", "Log Time", "Actions"];
	for (i = 0; i < ths.length; i++) {
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
	let userId;
	let userParams;
	let userSession;
	let resp;
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

		let td;
		for (j = 0; j < tds.length; j++) {
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
	let result;
	let ri = {
		method: "POST",
		body: JSON.stringify(data)
	};
	let resp = await fetch(rootPath + settings.ApiFolder, ri);
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
	return qpUserPage + "&" + qpnId + "=" + userId;
}

function reloadPage() {
	location.reload();
}
