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
adminPage = "/admin";
redirectDelay = 3;

// Global variables.
class GlobalVariablesContainer {
	constructor(areSettingsReady, settings, id, page, pages) {
		this.AreSettingsReady = areSettingsReady;
		this.Settings = settings;
		this.Id = id;
		this.Page = page;
		this.Pages = pages;
	}
}

mca_gvc = new GlobalVariablesContainer(false, null, 0, 0);

function isSettingsUpdateNeeded() {
	return true;
}

function saveSettings(s) {
	mca_gvc.Settings = s;
	mca_gvc.AreSettingsReady = true;
}

function getSettings() {
	if (!mca_gvc.AreSettingsReady) {
		console.error(err.Settings);
		return null;
	}

	return mca_gvc.Settings;
}

// Entry point.
async function onPageLoad() {
	// Settings initialisation.
	let ok = await updateSettingsIfNeeded();
	if (!ok) {
		return;
	}
	//let settings = getSettings();

	// Select a page.
	let curPage = window.location.search;
	let sp = new URLSearchParams(curPage);

	if (sp.has(qpn.ListOfUsers)) {
		if (!preparePageVariable(sp)) {
			return;
		}
		await showPage_ListOfUsers();
		return;
	}

	if (sp.has(qpn.ListOfLoggedUsers)) {
		if (!preparePageVariable(sp)) {
			return;
		}
		await showPage_ListOfLoggedUsers();
		return;
	}

	if (sp.has(qpn.RegistrationsReadyForApproval)) {
		if (!preparePageVariable(sp)) {
			return;
		}
		await showPage_RegistrationsReadyForApproval();
		return;
	}

	if (sp.has(qpn.UserPage)) {
		if (!prepareIdVariable(sp)) {
			return;
		}
		await showPage_UserPage();
		return;
	}

	if (sp.has(qpn.ManagerOfSections)) {
		await showPage_ManagerOfSections();
		return;
	}

	if (sp.has(qpn.ManagerOfForums)) {
		await showPage_ManagerOfForums();
		return;
	}

	if (sp.has(qpn.ManagerOfThreads)) {
		await showPage_ManagerOfThreads();
		return;
	}

	if (sp.has(qpn.ManagerOfMessages)) {
		await showPage_ManagerOfMessages();
		return;
	}

	if (sp.has(qpn.ManagerOfNotifications)) {
		await showPage_ManagerOfNotifications();
		return;
	}

	showPage_MainMenu();
}

async function onGoRegApprovalClick(btn) {
	await redirectToSubPage(false, qp.Prefix + qpn.RegistrationsReadyForApproval);
}

async function onGoLoggedUsersClick(btn) {
	await redirectToSubPage(false, qp.Prefix + qpn.ListOfLoggedUsers);
}

async function onGoListAllUsersClick(btn) {
	await redirectToSubPage(false, qp.Prefix + qpn.ListOfUsers);
}

async function onGoManageSectionsClick(btn) {
	await redirectToSubPage(false, qp.Prefix + qpn.ManagerOfSections);
}

async function onGoManageForumsClick(btn) {
	await redirectToSubPage(false, qp.Prefix + qpn.ManagerOfForums);
}

async function onGoManageThreadsClick(btn) {
	await redirectToSubPage(false, qp.Prefix + qpn.ManagerOfThreads);
}

async function onGoManageMessagesClick(btn) {
	await redirectToSubPage(false, qp.Prefix + qpn.ManagerOfMessages);
}

async function onGoManageNotificationsClick(btn) {
	await redirectToSubPage(false, qp.Prefix + qpn.ManagerOfNotifications);
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
	pageCount = repairUndefinedPageCount(pageCount);
	mca_gvc.Pages = pageCount;

	// Check page number for overflow.
	if (pageNumber > pageCount) {
		console.error(err.PageNotFound);
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
	let settings = getSettings();
	let pageCount = Math.ceil(userCount / settings.PageSize);
	pageCount = repairUndefinedPageCount(pageCount);
	mca_gvc.Pages = pageCount;

	// Check page number for overflow.
	if (pageNumber > pageCount) {
		console.error(err.PageNotFound);
		return;
	}

	let userIdsOnPage = calculateItemsOnPage(userIds, pageNumber, settings.PageSize);

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
	let resp = await getListOfRegistrationsReadyForApproval(pageNumber);
	if (resp == null) {
		return;
	}
	let pageCount = resp.result.totalPages;
	pageCount = repairUndefinedPageCount(pageCount);
	mca_gvc.Pages = pageCount;

	// Check page number for overflow.
	if (pageNumber > pageCount) {
		console.error(err.PageNotFound);
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

async function showPage_ManagerOfSections() {
	// Draw.
	let p = document.getElementById("subpage");
	p.style.display = "block";
	addBtnBack(p);
	addTitle(p, "Management of Sections");
	addDiv(p, "sectionManager");
	fillSectionManager("sectionManager");
}

async function showPage_ManagerOfForums() {
	// Draw.
	let p = document.getElementById("subpage");
	p.style.display = "block";
	addBtnBack(p);
	addTitle(p, "Management of Forums");
	addDiv(p, "forumManager");
	fillForumManager("forumManager");
}

async function showPage_ManagerOfThreads() {
	// Draw.
	let p = document.getElementById("subpage");
	p.style.display = "block";
	addBtnBack(p);
	addTitle(p, "Management of Threads");
	addDiv(p, "threadManager");
	fillThreadManager("threadManager");
}

async function showPage_ManagerOfMessages() {
	// Draw.
	let p = document.getElementById("subpage");
	p.style.display = "block";
	addBtnBack(p);
	addTitle(p, "Management of Messages");
	addDiv(p, "messageManager");
	fillMessageManager("messageManager");
}

async function showPage_ManagerOfNotifications() {
	// Draw.
	let p = document.getElementById("subpage");
	p.style.display = "block";
	addBtnBack(p);
	addTitle(p, "Management of Notifications");
	addDiv(p, "notificationManager");
	fillNotificationManager("notificationManager");
}

function addBtnBack(el) {
	let btn = document.createElement("INPUT");
	btn.type = "button";
	btn.className = "btnBack";
	btn.value = "Go Back";
	btn.addEventListener("click", async (e) => {
		await redirectToMainMenu(false);
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
			console.error(err.UnknownVariant);
	}
}

async function onBtnPrevClick_userList(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(err.PreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = composeUrlForAdminPage(qpn.ListOfUsers, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnNextClick_userList(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(err.NextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = composeUrlForAdminPage(qpn.ListOfUsers, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnPrevClick_logged(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(err.PreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = composeUrlForAdminPage(qpn.ListOfLoggedUsers, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnNextClick_logged(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(err.NextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = composeUrlForAdminPage(qpn.ListOfLoggedUsers, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnPrevClick_rrfa(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(err.PreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = composeUrlForAdminPage(qpn.RegistrationsReadyForApproval, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnNextClick_rrfa(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(err.NextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = composeUrlForAdminPage(qpn.RegistrationsReadyForApproval, mca_gvc.Page);
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
			url = composeUrlForUserPage(userId.toString());
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
	let ths = ["#", "ID", "PreRegTime", "E-Mail", "Name", "Actions"];
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

async function getListOfLoggedInUsers() {
	let reqData = new ApiRequest(actionName.GetListOfLoggedUsers, {});
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
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
	await reloadPage(false);
}

async function onBtnEnableRoleUPClick(role, userId) {
	let resp;
	switch (role) {
		case userRole.Author:
			resp = await setUserRoleAuthor(userId, true);
			break;

		case userRole.Writer:
			resp = await setUserRoleWriter(userId, true);
			break;

		case userRole.Reader:
			resp = await setUserRoleReader(userId, true);
			break;

		case userRole.Logging:
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
	await reloadPage(false);
}

async function onBtnDisableRoleUPClick(role, userId) {
	let resp;
	switch (role) {
		case userRole.Author:
			resp = await setUserRoleAuthor(userId, false);
			break;

		case userRole.Writer:
			resp = await setUserRoleWriter(userId, false);
			break;

		case userRole.Reader:
			resp = await setUserRoleReader(userId, false);
			break;

		case userRole.Logging:
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
	await reloadPage(false);
}

function calculateItemsOnPage(items, pageN, pageSize) {
	let x = Math.min(pageN * pageSize, items.length);
	return items.slice((pageN - 1) * pageSize, x);
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
		tds[2] = booleanToString(isUserLoggedIn);
		tds[3] = userParams.email;
		tds[4] = userParams.name;
		tds[5] = prettyTime(userParams.regTime);
		tds[6] = prettyTime(userParams.approvalTime);
		tds[7] = prettyTime(userParams.lastBadLogInTime);
		tds[8] = prettyTime(userParams.lastBadActionTime);
		tds[9] = prettyTime(userParams.banTime);
		tds[10] = booleanToString(userParams.canLogIn);
		tds[11] = booleanToString(userParams.isReader);
		tds[12] = booleanToString(userParams.isWriter);
		tds[13] = booleanToString(userParams.isAuthor);
		tds[14] = booleanToString(userParams.isModerator);
		tds[15] = booleanToString(userParams.isAdministrator);

		let td, url;
		for (let j = 0; j < tds.length; j++) {
			url = composeUrlForUserPage(userId.toString());
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
		booleanToString(userParams.isAdministrator),
		booleanToString(userParams.isModerator),
		booleanToString(userParams.isAuthor),
		booleanToString(userParams.isWriter),
		booleanToString(userParams.isReader),
		booleanToString(userParams.canLogIn),
		prettyTime(userParams.lastBadLogInTime),
		prettyTime(userParams.banTime),
		prettyTime(userParams.lastBadActionTime),
		booleanToString(userLogInState),
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
						'onclick="onBtnDisableRoleUPClick(\'' + userRole.Author + '\',' + userId + ')">';
				} else {
					actions = '<input type="button" class="btnEnableRole" value="Enable Role" ' +
						'onclick="onBtnEnableRoleUPClick(\'' + userRole.Author + '\',' + userId + ')">';
				}
				break;

			case "IsWriter":
				if (userParams.isWriter) {
					actions = '<input type="button" class="btnDisableRole" value="Disable Role" ' +
						'onclick="onBtnDisableRoleUPClick(\'' + userRole.Writer + '\',' + userId + ')">';
				} else {
					actions = '<input type="button" class="btnEnableRole" value="Enable Role" ' +
						'onclick="onBtnEnableRoleUPClick(\'' + userRole.Writer + '\',' + userId + ')">';
				}
				break;

			case "IsReader":
				if (userParams.isReader) {
					actions = '<input type="button" class="btnDisableRole" value="Disable Role" ' +
						'onclick="onBtnDisableRoleUPClick(\'' + userRole.Reader + '\',' + userId + ')">';
				} else {
					actions = '<input type="button" class="btnEnableRole" value="Enable Role" ' +
						'onclick="onBtnEnableRoleUPClick(\'' + userRole.Reader + '\',' + userId + ')">';
				}
				break;

			case "CanLogIn":
				if (userParams.canLogIn) {
					actions = '<input type="button" class="btnDisableRole" value="Disable Role" ' +
						'onclick="onBtnDisableRoleUPClick(\'' + userRole.Logging + '\',' + userId + ')">';
				} else {
					actions = '<input type="button" class="btnEnableRole" value="Enable Role" ' +
						'onclick="onBtnEnableRoleUPClick(\'' + userRole.Logging + '\',' + userId + ')">';
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

function fillSectionManager(elClass) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let fs = document.createElement("FIELDSET");
	div.appendChild(fs);

	let actionNames = ["Select an action", "Create a root section", "Create a normal section",
		"Change section's name", "Change section's parent", "Move section up & down", "Delete a section"];
	createRadioButtonsForActions(fs, actionNames);
	let d = document.createElement("DIV");
	d.innerHTML = '<input type="button" class="btnProceed" value="Proceed" onclick="onSectionManagerBtnProceedClick(this)">';
	fs.appendChild(d);
}

function fillForumManager(elClass) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let fs = document.createElement("FIELDSET");
	div.appendChild(fs);

	let actionNames = ["Select an action", "Create a forum", "Change forums's name",
		"Change forums's parent", "Move forum up & down", "Delete a forum"];
	createRadioButtonsForActions(fs, actionNames);
	let d = document.createElement("DIV");
	d.innerHTML = '<input type="button" class="btnProceed" value="Proceed" onclick="onForumManagerBtnProceedClick(this)">';
	fs.appendChild(d);
}

function fillThreadManager(elClass) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let fs = document.createElement("FIELDSET");
	div.appendChild(fs);

	let actionNames = ["Select an action", "Create a thread", "Change thread's name",
		"Change thread's parent", "Move thread up & down", "Delete a thread"];
	createRadioButtonsForActions(fs, actionNames);
	let d = document.createElement("DIV");
	d.innerHTML = '<input type="button" class="btnProceed" value="Proceed" onclick="onThreadManagerBtnProceedClick(this)">';
	fs.appendChild(d);
}

function fillMessageManager(elClass) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let fs = document.createElement("FIELDSET");
	div.appendChild(fs);

	let actionNames = ["Select an action",
		"Create a message", "Change message's text", "Change message's parent", "Delete a message"];
	createRadioButtonsForActions(fs, actionNames);
	let d = document.createElement("DIV");
	d.innerHTML = '<input type="button" class="btnProceed" value="Proceed" onclick="onMessageManagerBtnProceedClick(this)">';
	fs.appendChild(d);
}

function fillNotificationManager(elClass) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let fs = document.createElement("FIELDSET");
	div.appendChild(fs);

	let actionNames = ["Select an action", "Create a notification", "Delete a notification"];
	createRadioButtonsForActions(fs, actionNames);
	let d = document.createElement("DIV");
	d.innerHTML = '<input type="button" class="btnProceed" value="Proceed" onclick="onNotificationManagerBtnProceedClick(this)">';
	fs.appendChild(d);
}

function onSectionManagerBtnProceedClick(btn) {
	let selectedActionIdx = getSelectedActionIdxBPC(btn);
	if (selectedActionIdx == null) {
		return;
	}

	btn.disabled = true;
	disableParentFormBPC(btn);

	// Draw.
	let sm = document.getElementById("sectionManager");
	let fs = document.createElement("FIELDSET");
	sm.appendChild(fs);

	let d = document.createElement("DIV");
	d.className = "title";
	d.textContent = "Section Parameters";
	fs.appendChild(d);

	switch (selectedActionIdx) {
		case 1: // Create a root section.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="name">Name</label>' +
				'<input type="text" name="name" id="name" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnCreateRootSection" value="Create" onclick="onBtnCreateRootSectionClick(this)">';
			fs.appendChild(d);
			break;

		case 2: // Create a normal section.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="name">Name</label>' +
				'<input type="text" name="name" id="name" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="parent" title="ID of a parent section">Parent</label>' +
				'<input type="text" name="parent" id="parent" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnCreateNormalSection" value="Create" onclick="onBtnCreateNormalSectionClick(this)">';
			fs.appendChild(d);
			break;

		case 3: // Change section's name.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="id" title="ID of the changed section">ID</label>' +
				'<input type="text" name="id" id="id" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="name" title="New name of the section">New Name</label>' +
				'<input type="text" name="name" id="name" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnChangeSectionName" value="Change Name" onclick="onBtnChangeSectionNameClick(this)">';
			fs.appendChild(d);
			break;

		case 4: // Change section's parent.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="id" title="ID of the changed section">ID</label>' +
				'<input type="text" name="id" id="id" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="newParent" title="ID of the new parent">New Parent</label>' +
				'<input type="text" name="newParent" id="newParent" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnChangeSectionParent" value="Change Parent" onclick="onBtnChangeSectionParentClick(this)">';
			fs.appendChild(d);
			break;

		case 5: // Move section up & down.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="id" title="ID of the moved section">ID</label>' +
				'<input type="text" name="id" id="id" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnMoveSection" value="Move Up" onclick="onBtnMoveSectionUpClick(this)">' +
				'<span class="subpageSpacerA">&nbsp;</span>' +
				'<input type="button" class="btnMoveSection" value="Move Down" onclick="onBtnMoveSectionDownClick(this)">';
			fs.appendChild(d);
			break;

		case 6: // Delete a section.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="id" title="ID of the section to delete">ID</label>' +
				'<input type="text" name="id" id="id" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnDeleteSection" value="Delete" onclick="onBtnDeleteSectionClick(this)">';
			fs.appendChild(d);
			break;
	}
}

async function onBtnCreateRootSectionClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let name = pp.childNodes[1].childNodes[1].value;
	if (name.length < 1) {
		console.error(err.NameIsNotSet);
		return;
	}

	// Work.
	let resp = await addSection(null, name);
	if (resp == null) {
		return;
	}
	let sectionId = resp.result.sectionId;
	disableParentForm(btn, pp, false);
	let txt = "A root section was created. ID=" + sectionId.toString() + ".";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnCreateNormalSectionClick(btn) {
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
	let resp = await addSection(parent, name);
	if (resp == null) {
		return;
	}
	let sectionId = resp.result.sectionId;
	disableParentForm(btn, pp, false);
	let txt = "A normal section was created. ID=" + sectionId.toString() + ".";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnChangeSectionNameClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let sectionId = Number(pp.childNodes[1].childNodes[1].value);
	if (sectionId < 1) {
		console.error(err.IdNotSet);
		return;
	}
	let newName = pp.childNodes[2].childNodes[1].value;
	if (newName.length < 1) {
		console.error(err.NameIsNotSet);
		return;
	}

	// Work.
	let resp = await changeSectionName(sectionId, newName);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Section name was changed.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnChangeSectionParentClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let sectionId = Number(pp.childNodes[1].childNodes[1].value);
	if (sectionId < 1) {
		console.error(err.IdNotSet);
		return;
	}
	let newParent = Number(pp.childNodes[2].childNodes[1].value);
	if (newParent < 1) {
		console.error(err.IdNotSet);
		return;
	}

	// Work.
	let resp = await changeSectionParent(sectionId, newParent);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Section was moved to a new parent.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnMoveSectionUpClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let sectionId = Number(pp.childNodes[1].childNodes[1].value);
	if (sectionId < 1) {
		console.error(err.IdNotSet);
		return;
	}

	// Work.
	let resp = await moveSectionUp(sectionId);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, true);
	let txt = "Section was moved up.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnMoveSectionDownClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let sectionId = Number(pp.childNodes[1].childNodes[1].value);
	if (sectionId < 1) {
		console.error(err.IdNotSet);
		return;
	}

	// Work.
	let resp = await moveSectionDown(sectionId);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, true);
	let txt = "Section was moved down.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnDeleteSectionClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let sectionId = Number(pp.childNodes[1].childNodes[1].value);
	if (sectionId < 1) {
		console.error(err.IdNotSet);
		return;
	}

	// Work.
	let resp = await deleteSection(sectionId);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Section was deleted.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

function onForumManagerBtnProceedClick(btn) {
	let selectedActionIdx = getSelectedActionIdxBPC(btn);
	if (selectedActionIdx == null) {
		return;
	}

	btn.disabled = true;
	disableParentFormBPC(btn);

	// Draw.
	let fm = document.getElementById("forumManager");
	let fs = document.createElement("FIELDSET");
	fm.appendChild(fs);

	let d = document.createElement("DIV");
	d.className = "title";
	d.textContent = "Forum Parameters";
	fs.appendChild(d);

	switch (selectedActionIdx) {
		case 1: // Create a forum.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="name">Name</label>' +
				'<input type="text" name="name" id="name" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="parent" title="ID of a parent section">Parent</label>' +
				'<input type="text" name="parent" id="parent" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnCreateForum" value="Create" onclick="onBtnCreateForumClick(this)">';
			fs.appendChild(d);
			break;

		case 2: // Change forum's name.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="id" title="ID of the changed forum">ID</label>' +
				'<input type="text" name="id" id="id" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="name" title="New name of the forum">New Name</label>' +
				'<input type="text" name="name" id="name" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnChangeForumName" value="Change Name" onclick="onBtnChangeForumNameClick(this)">';
			fs.appendChild(d);
			break;

		case 3: // Change forum's parent.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="id" title="ID of the changed forum">ID</label>' +
				'<input type="text" name="id" id="id" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="newParent" title="ID of the new parent">New Parent</label>' +
				'<input type="text" name="newParent" id="newParent" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnChangeForumParent" value="Change Parent" onclick="onBtnChangeForumParentClick(this)">';
			fs.appendChild(d);
			break;

		case 4: // Move forum up & down.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="id" title="ID of the moved forum">ID</label>' +
				'<input type="text" name="id" id="id" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnMoveForum" value="Move Up" onclick="onBtnMoveForumUpClick(this)">' +
				'<span class="subpageSpacerA">&nbsp;</span>' +
				'<input type="button" class="btnMoveForum" value="Move Down" onclick="onBtnMoveForumDownClick(this)">';
			fs.appendChild(d);
			break;

		case 5: // Delete a forum.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="id" title="ID of the forum to delete">ID</label>' +
				'<input type="text" name="id" id="id" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnDeleteForum" value="Delete" onclick="onBtnDeleteForumClick(this)">';
			fs.appendChild(d);
			break;
	}
}

async function onBtnCreateForumClick(btn) {
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
	let resp = await addForum(parent, name);
	if (resp == null) {
		return;
	}
	let forumId = resp.result.forumId;
	disableParentForm(btn, pp, false);
	let txt = "A forum was created. ID=" + forumId.toString() + ".";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnChangeForumNameClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let forumId = Number(pp.childNodes[1].childNodes[1].value);
	if (forumId < 1) {
		console.error(err.IdNotSet);
		return;
	}
	let newName = pp.childNodes[2].childNodes[1].value;
	if (newName.length < 1) {
		console.error(err.NameIsNotSet);
		return;
	}

	// Work.
	let resp = await changeForumName(forumId, newName);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Forum name was changed.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnChangeForumParentClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let forumId = Number(pp.childNodes[1].childNodes[1].value);
	if (forumId < 1) {
		console.error(err.IdNotSet);
		return;
	}
	let newParent = Number(pp.childNodes[2].childNodes[1].value);
	if (newParent < 1) {
		console.error(err.IdNotSet);
		return;
	}

	// Work.
	let resp = await changeForumSection(forumId, newParent);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Forum was moved to a new parent.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnMoveForumUpClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let forumId = Number(pp.childNodes[1].childNodes[1].value);
	if (forumId < 1) {
		console.error(err.IdNotSet);
		return;
	}

	// Work.
	let resp = await moveForumUp(forumId);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, true);
	let txt = "Forum was moved up.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnMoveForumDownClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let forumId = Number(pp.childNodes[1].childNodes[1].value);
	if (forumId < 1) {
		console.error(err.IdNotSet);
		return;
	}

	// Work.
	let resp = await moveForumDown(forumId);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, true);
	let txt = "Forum was moved down.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnDeleteForumClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let forumId = Number(pp.childNodes[1].childNodes[1].value);
	if (forumId < 1) {
		console.error(err.IdNotSet);
		return;
	}

	// Work.
	let resp = await deleteForum(forumId);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Forum was deleted.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

function onThreadManagerBtnProceedClick(btn) {
	let selectedActionIdx = getSelectedActionIdxBPC(btn);
	if (selectedActionIdx == null) {
		return;
	}

	btn.disabled = true;
	disableParentFormBPC(btn);

	// Draw.
	let tm = document.getElementById("threadManager");
	let fs = document.createElement("FIELDSET");
	tm.appendChild(fs);

	let d = document.createElement("DIV");
	d.className = "title";
	d.textContent = "Thread Parameters";
	fs.appendChild(d);

	switch (selectedActionIdx) {
		case 1: // Create a thread.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="name">Name</label>' +
				'<input type="text" name="name" id="name" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="parent" title="ID of a parent forum">Parent</label>' +
				'<input type="text" name="parent" id="parent" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnCreateThread" value="Create" onclick="onBtnCreateThreadClick(this)">';
			fs.appendChild(d);
			break;

		case 2: // Change thread's name.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="id" title="ID of the changed thread">ID</label>' +
				'<input type="text" name="id" id="id" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="name" title="New name of the thread">New Name</label>' +
				'<input type="text" name="name" id="name" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnChangeThreadName" value="Change Name" onclick="onBtnChangeThreadNameClick(this)">';
			fs.appendChild(d);
			break;

		case 3: // Change thread's parent.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="id" title="ID of the changed thread">ID</label>' +
				'<input type="text" name="id" id="id" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="newParent" title="ID of the new parent">New Parent</label>' +
				'<input type="text" name="newParent" id="newParent" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnChangeThreadParent" value="Change Parent" onclick="onBtnChangeThreadParentClick(this)">';
			fs.appendChild(d);
			break;

		case 4: // Move thread up & down.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="id" title="ID of the moved thread">ID</label>' +
				'<input type="text" name="id" id="id" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnMoveThread" value="Move Up" onclick="onBtnMoveThreadUpClick(this)">' +
				'<span class="subpageSpacerA">&nbsp;</span>' +
				'<input type="button" class="btnMoveThread" value="Move Down" onclick="onBtnMoveThreadDownClick(this)">';
			fs.appendChild(d);
			break;

		case 5: // Delete a thread.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="id" title="ID of the thread to delete">ID</label>' +
				'<input type="text" name="id" id="id" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnDeleteThread" value="Delete" onclick="onBtnDeleteThreadClick(this)">';
			fs.appendChild(d);
			break;
	}
}

async function onBtnCreateThreadClick(btn) {
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

async function onBtnChangeThreadNameClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let threadId = Number(pp.childNodes[1].childNodes[1].value);
	if (threadId < 1) {
		console.error(err.IdNotSet);
		return;
	}
	let newName = pp.childNodes[2].childNodes[1].value;
	if (newName.length < 1) {
		console.error(err.NameIsNotSet);
		return;
	}

	// Work.
	let resp = await changeThreadName(threadId, newName);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Thread name was changed.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnChangeThreadParentClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let threadId = Number(pp.childNodes[1].childNodes[1].value);
	if (threadId < 1) {
		console.error(err.IdNotSet);
		return;
	}
	let newParent = Number(pp.childNodes[2].childNodes[1].value);
	if (newParent < 1) {
		console.error(err.IdNotSet);
		return;
	}

	// Work.
	let resp = await changeThreadForum(threadId, newParent);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Thread was moved to a new parent.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnMoveThreadUpClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let threadId = Number(pp.childNodes[1].childNodes[1].value);
	if (threadId < 1) {
		console.error(err.IdNotSet);
		return;
	}

	// Work.
	let resp = await moveThreadUp(threadId);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, true);
	let txt = "Thread was moved up.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnMoveThreadDownClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let threadId = Number(pp.childNodes[1].childNodes[1].value);
	if (threadId < 1) {
		console.error(err.IdNotSet);
		return;
	}

	// Work.
	let resp = await moveThreadDown(threadId);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, true);
	let txt = "Thread was moved down.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnDeleteThreadClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let threadId = Number(pp.childNodes[1].childNodes[1].value);
	if (threadId < 1) {
		console.error(err.IdNotSet);
		return;
	}

	// Work.
	let resp = await deleteThread(threadId);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Thread was deleted.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

function onMessageManagerBtnProceedClick(btn) {
	let selectedActionIdx = getSelectedActionIdxBPC(btn);
	if (selectedActionIdx == null) {
		return;
	}

	btn.disabled = true;
	disableParentFormBPC(btn);

	// Draw.
	let mm = document.getElementById("messageManager");
	let fs = document.createElement("FIELDSET");
	mm.appendChild(fs);

	let d = document.createElement("DIV");
	d.className = "title";
	d.textContent = "Message Parameters";
	fs.appendChild(d);

	switch (selectedActionIdx) {
		case 1: // Create a message.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="txt">Text</label>' +
				'<input type="text" name="txt" id="txt" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="parent" title="ID of a parent thread">Parent</label>' +
				'<input type="text" name="parent" id="parent" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnCreateMessage" value="Create" onclick="onBtnCreateMessageClick(this)">';
			fs.appendChild(d);
			break;

		case 2: // Change message's text.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="id" title="ID of the changed message">ID</label>' +
				'<input type="text" name="id" id="id" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="txt" title="New text of the message">New Text</label>' +
				'<input type="text" name="txt" id="txt" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnChangeMessageText" value="Change Text" onclick="onBtnChangeMessageTextClick(this)">';
			fs.appendChild(d);
			break;

		case 3: // Change message's parent.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="id" title="ID of the changed message">ID</label>' +
				'<input type="text" name="id" id="id" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="newParent" title="ID of the new parent">New Parent</label>' +
				'<input type="text" name="newParent" id="newParent" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnChangeMessageParent" value="Change Parent" onclick="onBtnChangeMessageParentClick(this)">';
			fs.appendChild(d);
			break;

		case 4: // Delete a message.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="id" title="ID of the message to delete">ID</label>' +
				'<input type="text" name="id" id="id" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnDeleteMessage" value="Delete" onclick="onBtnDeleteMessageClick(this)">';
			fs.appendChild(d);
			break;
	}
}

function onNotificationManagerBtnProceedClick(btn) {
	let selectedActionIdx = getSelectedActionIdxBPC(btn);
	if (selectedActionIdx == null) {
		return;
	}

	btn.disabled = true;
	disableParentFormBPC(btn);

	// Draw.
	let nm = document.getElementById("notificationManager");
	let fs = document.createElement("FIELDSET");
	nm.appendChild(fs);

	let d = document.createElement("DIV");
	d.className = "title";
	d.textContent = "Notification Parameters";
	fs.appendChild(d);

	switch (selectedActionIdx) {
		case 1: // Create a notification.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="txt">Text</label>' +
				'<input type="text" name="txt" id="txt" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="user" title="ID of a user">User</label>' +
				'<input type="text" name="user" id="user" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnCreateNotification" value="Create" onclick="onBtnCreateNotificationClick(this)">';
			fs.appendChild(d);
			break;

		case 2: // Delete a notification.
			d = document.createElement("DIV");
			d.innerHTML = '<label class="parameter" for="id" title="ID of the notification to delete">ID</label>' +
				'<input type="text" name="id" id="id" value="" />';
			fs.appendChild(d);
			d = document.createElement("DIV");
			d.innerHTML = '<input type="button" class="btnDeleteNotification" value="Delete" onclick="onBtnDeleteNotificationClick(this)">';
			fs.appendChild(d);
			break;
	}
}

async function onBtnCreateMessageClick(btn) {
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

async function onBtnCreateNotificationClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let text = pp.childNodes[1].childNodes[1].value;
	if (text.length < 1) {
		console.error(err.TextIsNotSet);
		return;
	}
	let userId = Number(pp.childNodes[2].childNodes[1].value);
	if (userId < 1) {
		console.error(err.IdNotSet);
		return;
	}

	// Work.
	let resp = await addNotification(userId, text);
	if (resp == null) {
		return;
	}
	let notificationId = resp.result.notificationId;
	disableParentForm(btn, pp, false);
	let txt = "A notification was created. ID=" + notificationId.toString() + ".";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnChangeMessageTextClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let messageId = Number(pp.childNodes[1].childNodes[1].value);
	if (messageId < 1) {
		console.error(err.IdNotSet);
		return;
	}
	let newText = pp.childNodes[2].childNodes[1].value;
	if (newText.length < 1) {
		console.error(err.TextIsNotSet);
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

async function onBtnChangeMessageParentClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let messageId = Number(pp.childNodes[1].childNodes[1].value);
	if (messageId < 1) {
		console.error(err.IdNotSet);
		return;
	}
	let newParent = Number(pp.childNodes[2].childNodes[1].value);
	if (newParent < 1) {
		console.error(err.IdNotSet);
		return;
	}

	// Work.
	let resp = await changeMessageThread(messageId, newParent);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Message was moved to a new parent.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnDeleteMessageClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let messageId = Number(pp.childNodes[1].childNodes[1].value);
	if (messageId < 1) {
		console.error(err.IdNotSet);
		return;
	}

	// Work.
	let resp = await deleteMessage(messageId);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Message was deleted.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnDeleteNotificationClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let notificationId = Number(pp.childNodes[1].childNodes[1].value);
	if (notificationId < 1) {
		console.error(err.IdNotSet);
		return;
	}

	// Work.
	let resp = await deleteNotification(notificationId);
	if (resp == null) {
		return;
	}
	if (resp.result.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Notification was deleted.";
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

function disableParentFormBPC(btn) {
	let pp = btn.parentNode.parentNode;
	for (let i = 0; i < pp.childNodes.length; i++) {
		let ch = pp.childNodes[i];
		ch.childNodes[0].disabled = true;
	}
}

function getSelectedActionIdxBPC(btn) {
	let selectedActionIdx = 0;
	let pp = btn.parentNode.parentNode;
	for (let i = 0; i < pp.childNodes.length; i++) {
		let ch = pp.childNodes[i];
		if (ch.childNodes[0].checked === true) {
			selectedActionIdx = i;
			break;
		}
	}
	if (selectedActionIdx < 1) {
		return null;
	}
	return selectedActionIdx;
}

function createRadioButtonsForActions(fs, actionNames) {
	for (let i = 0; i < actionNames.length; i++) {
		let d = document.createElement("DIV");
		if (i === 0) {
			d.className = "title";
			d.textContent = actionNames[i];
		} else {
			d.innerHTML = '<input type="radio" name="action" id="action_' + i + '" value="' + actionNames[i] + '" />' +
				'<label class="action" for="action_' + i + '">' + actionNames[i] + '</label>';

		}
		fs.appendChild(d);
	}
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

async function redirectToSubPage(wait, qp) {
	let url = adminPage + qp;
	await redirectPage(wait, url);
}

async function redirectToMainMenu(wait) {
	let url = adminPage;
	await redirectPage(wait, url);
}

function composeUrlForUserPage(userId) {
	return qp.Prefix + qpn.UserPage + "&" + qpn.Id + "=" + userId;
}

function composeUrlForAdminPage(func, page) {
	return qp.Prefix + func + "&" + qpn.Page + "=" + page;
}
