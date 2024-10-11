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

// Names of Query Parameters.
qp = {
	Prefix: "?",
}

qpn = {
	Id: "id",
	Page: "page",
	ListOfUsers: "listOfUsers",
	ListOfLoggedUsers: "listOfLoggedUsers",
	RegistrationsReadyForApproval: "registrationsReadyForApproval",
	UserPage: "userPage",
	ManagerOfSections: "manageSections",
	ManagerOfForums: "manageForums",
	ManagerOfThreads: "managerOfThreads",
	ManagerOfMessages: "managerOfMessages",
	ManagerOfNotifications: "managerOfNotifications",
}

// Action names.
actionName = {
	AddForum: "addForum",
	AddMessage: "addMessage",
	AddNotification: "addNotification",
	AddSection: "addSection",
	AddThread: "addThread",
	ApproveAndRegisterUser: "approveAndRegisterUser",
	BanUser: "banUser",
	ChangeForumName: "changeForumName",
	ChangeForumSection: "changeForumSection",
	ChangeMessageText: "changeMessageText",
	ChangeMessageThread: "changeMessageThread",
	ChangeSectionName: "changeSectionName",
	ChangeSectionParent: "changeSectionParent",
	ChangeThreadName: "changeThreadName",
	ChangeThreadForum: "changeThreadForum",
	DeleteForum: "deleteForum",
	DeleteMessage: "deleteMessage",
	DeleteNotification: "deleteNotification",
	DeleteSection: "deleteSection",
	DeleteThread: "deleteThread",
	GetListOfAllUsers: "getListOfAllUsers",
	GetListOfLoggedUsers: "getListOfLoggedUsers",
	GetListOfRegistrationsReadyForApproval: "getListOfRegistrationsReadyForApproval",
	GetUserSession: "getUserSession",
	IsUserLoggedIn: "isUserLoggedIn",
	LogUserOutA: "logUserOutA",
	MoveForumDown: "moveForumDown",
	MoveForumUp: "moveForumUp",
	MoveSectionDown: "moveSectionDown",
	MoveSectionUp: "moveSectionUp",
	MoveThreadDown: "moveThreadDown",
	MoveThreadUp: "moveThreadUp",
	RejectRegistrationRequest: "rejectRegistrationRequest",
	SetUserRoleAuthor: "setUserRoleAuthor",
	SetUserRoleReader: "setUserRoleReader",
	SetUserRoleWriter: "setUserRoleWriter",
	UnbanUser: "unbanUser",
	ViewUserParameters: "viewUserParameters",
}

// Messages.
msg = {
	GenericErrorPrefix: "Error: ",
}

// Errors.
err = {
	IdNotSet: "ID is not set",
	IdNotFound: "ID is not found",
	PageNotSet: "page is not set",
	PageNotFound: "page is not found",
	Settings: "settings error",
	NotOk: "something went wrong",
	Server: "server error",
	Client: "client error",
	Unknown: "unknown error",
	PreviousPageDoesNotExist: "previous page does not exist",
	NextPageDoesNotExist: "next page does not exist",
	UnknownVariant: "unknown variant",
	NameIsNotSet: "name is not set",
	ParentIsNotSet: "parent is not set",
	TextIsNotSet: "text is not set",
}

// User role names.
userRole = {
	Author: "author",
	Writer: "writer",
	Reader: "reader",
	Logging: "logging",
}

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
	constructor(registrationRequestId) {
		this.RegistrationRequestId = registrationRequestId;
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

class Parameters_AddSection {
	constructor(parent, name) {
		this.Parent = parent;
		this.Name = name;
	}
}

class Parameters_AddForum {
	constructor(parent, name) {
		this.SectionId = parent;
		this.Name = name;
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

class Parameters_ChangeSectionName {
	constructor(sectionId, name) {
		this.SectionId = sectionId;
		this.Name = name;
	}
}

class Parameters_ChangeForumName {
	constructor(forumId, name) {
		this.ForumId = forumId;
		this.Name = name;
	}
}

class Parameters_ChangeThreadName {
	constructor(threadId, name) {
		this.ThreadId = threadId;
		this.Name = name;
	}
}

class Parameters_ChangeMessageText {
	constructor(messageId, text) {
		this.MessageId = messageId;
		this.Text = text;
	}
}

class Parameters_ChangeSectionParent {
	constructor(sectionId, parent) {
		this.SectionId = sectionId;
		this.Parent = parent;
	}
}

class Parameters_ChangeForumSection {
	constructor(forumId, sectionId) {
		this.ForumId = forumId;
		this.SectionId = sectionId; // Parent.
	}
}

class Parameters_ChangeThreadForum {
	constructor(threadId, forumId) {
		this.ThreadId = threadId;
		this.ForumId = forumId; // Parent.
	}
}

class Parameters_ChangeMessageThread {
	constructor(messageId, threadId) {
		this.MessageId = messageId;
		this.ThreadId = threadId; // Parent.
	}
}

class Parameters_MoveSectionUp {
	constructor(sectionId) {
		this.SectionId = sectionId;
	}
}

class Parameters_MoveSectionDown {
	constructor(sectionId) {
		this.SectionId = sectionId;
	}
}

class Parameters_MoveForumUp {
	constructor(forumId) {
		this.ForumId = forumId;
	}
}

class Parameters_MoveForumDown {
	constructor(forumId) {
		this.ForumId = forumId;
	}
}

class Parameters_MoveThreadUp {
	constructor(threadId) {
		this.ThreadId = threadId;
	}
}

class Parameters_MoveThreadDown {
	constructor(threadId) {
		this.ThreadId = threadId;
	}
}

class Parameters_DeleteSection {
	constructor(sectionId) {
		this.SectionId = sectionId;
	}
}

class Parameters_DeleteForum {
	constructor(forumId) {
		this.ForumId = forumId;
	}
}

class Parameters_DeleteThread {
	constructor(threadId) {
		this.ThreadId = threadId;
	}
}

class Parameters_DeleteMessage {
	constructor(messageId) {
		this.MessageId = messageId;
	}
}

class Parameters_AddNotification {
	constructor(userId, text) {
		this.UserId = userId;
		this.Text = text;
	}
}

class Parameters_DeleteNotification {
	constructor(notificationId) {
		this.NotificationId = notificationId;
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
		console.error(err.Settings);
		return null;
	}

	return x;
}

async function fetchSettings() {
	let data = await fetch(rootPath + settingsPath);
	return await data.json();
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
	let pageCount = Math.ceil(userCount / mca_gvc.Settings.PageSize);
	pageCount = repairUndefinedPageCount(pageCount);
	mca_gvc.Pages = pageCount;

	// Check page number for overflow.
	if (pageNumber > pageCount) {
		console.error(err.PageNotFound);
		return;
	}

	let userIdsOnPage = calculateItemsOnPage(userIds, pageNumber, mca_gvc.Settings.PageSize);

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

async function getListOfAllUsers(pageN) {
	let params = new Parameters_GetListOfAllUsers(pageN);
	let reqData = new ApiRequest(actionName.GetListOfAllUsers, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function getListOfRegistrationsReadyForApproval(pageN) {
	let params = new Parameters_GetListOfRegistrationsReadyForApproval(pageN);
	let reqData = new ApiRequest(actionName.GetListOfRegistrationsReadyForApproval, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
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

async function viewUserParameters(userId) {
	let params = new Parameters_ViewUserParameters(userId);
	let reqData = new ApiRequest(actionName.ViewUserParameters, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function isUserLoggedIn(userId) {
	let params = new Parameters_IsUserLoggedIn(userId);
	let reqData = new ApiRequest(actionName.IsUserLoggedIn, params);
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

async function approveAndRegisterUser(email) {
	let params = new Parameters_ApproveAndRegisterUser(email);
	let reqData = new ApiRequest(actionName.ApproveAndRegisterUser, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function rejectRegistrationRequest(registrationRequestId) {
	let params = new Parameters_RejectRegistrationRequest(registrationRequestId);
	let reqData = new ApiRequest(actionName.RejectRegistrationRequest, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

function calculateItemsOnPage(items, pageN, pageSize) {
	let x = Math.min(pageN * pageSize, items.length);
	return items.slice((pageN - 1) * pageSize, x);
}

async function logUserOutA(userId) {
	let params = new Parameters_LogUserOutA(userId);
	let reqData = new ApiRequest(actionName.LogUserOutA, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function setUserRoleAuthor(userId, roleValue) {
	let params = new Parameters_SetUserRoleAuthor(userId, roleValue);
	let reqData = new ApiRequest(actionName.SetUserRoleAuthor, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function setUserRoleWriter(userId, roleValue) {
	let params = new Parameters_SetUserRoleWriter(userId, roleValue);
	let reqData = new ApiRequest(actionName.SetUserRoleWriter, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function setUserRoleReader(userId, roleValue) {
	let params = new Parameters_SetUserRoleReader(userId, roleValue);
	let reqData = new ApiRequest(actionName.SetUserRoleReader, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function banUser(userId) {
	let params = new Parameters_BanUser(userId);
	let reqData = new ApiRequest(actionName.BanUser, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function unbanUser(userId) {
	let params = new Parameters_UnbanUser(userId);
	let reqData = new ApiRequest(actionName.UnbanUser, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function getUserSession(userId) {
	let params = new Parameters_GetUserSession(userId);
	let reqData = new ApiRequest(actionName.GetUserSession, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

function composeUserPageLink(userId) {
	return qp.Prefix + qpn.UserPage + "&" + qpn.Id + "=" + userId;
}

async function reloadPage(wait) {
	if (wait) {
		await sleep(redirectDelay * 1000);
	}
	location.reload();
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

function composeUrlForAdminPage(func, page) {
	return qp.Prefix + func + "&" + qpn.Page + "=" + page;
}

async function addSection(parent, name) {
	let params = new Parameters_AddSection(parent, name);
	let reqData = new ApiRequest(actionName.AddSection, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function addForum(parent, name) {
	let params = new Parameters_AddForum(parent, name);
	let reqData = new ApiRequest(actionName.AddForum, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
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

async function changeSectionName(sectionId, name) {
	let params = new Parameters_ChangeSectionName(sectionId, name);
	let reqData = new ApiRequest(actionName.ChangeSectionName, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function changeForumName(forumId, name) {
	let params = new Parameters_ChangeForumName(forumId, name);
	let reqData = new ApiRequest(actionName.ChangeForumName, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function changeThreadName(threadId, name) {
	let params = new Parameters_ChangeThreadName(threadId, name);
	let reqData = new ApiRequest(actionName.ChangeThreadName, params);
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

async function changeSectionParent(sectionId, newParent) {
	let params = new Parameters_ChangeSectionParent(sectionId, newParent);
	let reqData = new ApiRequest(actionName.ChangeSectionParent, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function changeForumSection(forumId, newParent) {
	let params = new Parameters_ChangeForumSection(forumId, newParent);
	let reqData = new ApiRequest(actionName.ChangeForumSection, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function changeThreadForum(threadId, newParent) {
	let params = new Parameters_ChangeThreadForum(threadId, newParent);
	let reqData = new ApiRequest(actionName.ChangeThreadForum, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function changeMessageThread(messageId, newParent) {
	let params = new Parameters_ChangeMessageThread(messageId, newParent);
	let reqData = new ApiRequest(actionName.ChangeMessageThread, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function moveSectionUp(sectionId) {
	let params = new Parameters_MoveSectionUp(sectionId);
	let reqData = new ApiRequest(actionName.MoveSectionUp, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function moveSectionDown(sectionId) {
	let params = new Parameters_MoveSectionDown(sectionId);
	let reqData = new ApiRequest(actionName.MoveSectionDown, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function moveForumUp(forumId) {
	let params = new Parameters_MoveForumUp(forumId);
	let reqData = new ApiRequest(actionName.MoveForumUp, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function moveForumDown(forumId) {
	let params = new Parameters_MoveForumDown(forumId);
	let reqData = new ApiRequest(actionName.MoveForumDown, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function moveThreadUp(threadId) {
	let params = new Parameters_MoveThreadUp(threadId);
	let reqData = new ApiRequest(actionName.MoveThreadUp, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function moveThreadDown(threadId) {
	let params = new Parameters_MoveThreadDown(threadId);
	let reqData = new ApiRequest(actionName.MoveThreadDown, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function deleteSection(sectionId) {
	let params = new Parameters_DeleteSection(sectionId);
	let reqData = new ApiRequest(actionName.DeleteSection, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function deleteForum(forumId) {
	let params = new Parameters_DeleteForum(forumId);
	let reqData = new ApiRequest(actionName.DeleteForum, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function deleteThread(threadId) {
	let params = new Parameters_DeleteThread(threadId);
	let reqData = new ApiRequest(actionName.DeleteThread, params);
	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}
	return resp.JsonObject;
}

async function deleteMessage(messageId) {
	let params = new Parameters_DeleteMessage(messageId);
	let reqData = new ApiRequest(actionName.DeleteMessage, params);
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

async function addNotification(userId, text) {
	let params = new Parameters_AddNotification(userId, text);
	let reqData = new ApiRequest(actionName.AddNotification, params);
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
