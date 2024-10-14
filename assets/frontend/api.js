// Pages and Query Parameters.
qp = {
	ChangeEmailStep1: "?changeEmail1",
	ChangeEmailStep2: "?changeEmail2",
	ChangeEmailStep3: "?changeEmail3",
	ChangePwdStep1: "?changePwd1",
	ChangePwdStep2: "?changePwd2",
	ChangePwdStep3: "?changePwd3",
	LogInStep1: "?login1",
	LogInStep2: "?login2",
	LogInStep3: "?login3",
	LogInStep4: "?login4",
	LogOutStep1: "?logout1",
	LogOutStep2: "?logout2",
	Prefix: "?",
	RegistrationStep1: "?reg1",
	RegistrationStep2: "?reg2",
	RegistrationStep3: "?reg3",
	RegistrationStep4: "?reg4",
	SelfPage: "?selfPage",
}

qpn = {
	Forum: "forum",
	Id: "id",
	ListOfUsers: "listOfUsers",
	ListOfLoggedUsers: "listOfLoggedUsers",
	ManagerOfSections: "manageSections",
	ManagerOfForums: "manageForums",
	ManagerOfThreads: "managerOfThreads",
	ManagerOfMessages: "managerOfMessages",
	ManagerOfNotifications: "managerOfNotifications",
	Message: "message",
	Notifications: "notifications",
	Page: "page",
	RegistrationsReadyForApproval: "registrationsReadyForApproval",
	Section: "section",
	SubscriptionsPage: "subscriptions",
	Thread: "thread",
	UserPage: "userPage",
}

// Action names.
actionName = {
	AddForum: "addForum",
	AddMessage: "addMessage",
	AddNotification: "addNotification",
	AddSection: "addSection",
	AddSubscription: "addSubscription",
	AddThread: "addThread",
	ApproveAndRegisterUser: "approveAndRegisterUser",
	BanUser: "banUser",
	ChangeEmail: "changeEmail",
	ChangeForumName: "changeForumName",
	ChangeForumSection: "changeForumSection",
	ChangeMessageText: "changeMessageText",
	ChangeMessageThread: "changeMessageThread",
	ChangePwd: "changePassword",
	ChangeSectionName: "changeSectionName",
	ChangeSectionParent: "changeSectionParent",
	ChangeThreadName: "changeThreadName",
	ChangeThreadForum: "changeThreadForum",
	CountSelfSubscriptions: "countSelfSubscriptions",
	CountUnreadNotifications: "countUnreadNotifications",
	DeleteForum: "deleteForum",
	DeleteMessage: "deleteMessage",
	DeleteNotification: "deleteNotification",
	DeleteSection: "deleteSection",
	DeleteSelfSubscription: "deleteSelfSubscription",
	DeleteThread: "deleteThread",
	GetAllNotifications: "getAllNotifications",
	GetLatestMessageOfThread: "getLatestMessageOfThread",
	GetListOfAllUsers: "getListOfAllUsers",
	GetListOfLoggedUsers: "getListOfLoggedUsers",
	GetListOfRegistrationsReadyForApproval: "getListOfRegistrationsReadyForApproval",
	GetMessage: "getMessage",
	GetNotificationsOnPage: "getNotificationsOnPage",
	GetSelfRoles: "getSelfRoles",
	GetSelfSubscriptions: "getSelfSubscriptions",
	GetSelfSubscriptionsOnPage: "getSelfSubscriptionsOnPage",
	GetThread: "getThread",
	GetThreadNamesByIds: "getThreadNamesByIds",
	GetUserName: "getUserName",
	GetUserSession: "getUserSession",
	IsSelfSubscribed: "isSelfSubscribed",
	IsUserLoggedIn: "isUserLoggedIn",
	ListForumAndThreadsOnPage: "listForumAndThreadsOnPage",
	ListSectionsAndForums: "listSectionsAndForums",
	ListThreadAndMessagesOnPage: "listThreadAndMessagesOnPage",
	LogUserIn: "logUserIn",
	LogUserOut: "logUserOut",
	LogUserOutA: "logUserOutA",
	MoveForumDown: "moveForumDown",
	MoveForumUp: "moveForumUp",
	MoveSectionDown: "moveSectionDown",
	MoveSectionUp: "moveSectionUp",
	MoveThreadDown: "moveThreadDown",
	MoveThreadUp: "moveThreadUp",
	MarkNotificationAsRead: "markNotificationAsRead",
	RegisterUser: "registerUser",
	RejectRegistrationRequest: "rejectRegistrationRequest",
	SetUserRoleAuthor: "setUserRoleAuthor",
	SetUserRoleReader: "setUserRoleReader",
	SetUserRoleWriter: "setUserRoleWriter",
	UnbanUser: "unbanUser",
	ViewUserParameters: "viewUserParameters",
}

// Errors.
err = {
	Client: "client error",
	DuplicateMapKey: "duplicate map key",
	ElementTypeUnsupported: "unsupported element type",
	IdNotSet: "ID is not set",
	IdNotFound: "ID is not found",
	MessageNotFound: "message is not found",
	NameIsNotSet: "name is not set",
	NextPageDoesNotExist: "next page does not exist",
	NextStepUnknown: "unknown next step",
	NotOk: "something went wrong",
	PageNotSet: "page is not set",
	PageNotFound: "page is not found",
	ParentIsNotSet: "parent is not set",
	PasswordNotValid: "password is not valid",
	PreviousPageDoesNotExist: "previous page does not exist",
	RootSectionNotFound: "root section is not found",
	SectionNotFound: "section is not found",
	Server: "server error",
	Settings: "settings error",
	TextIsNotSet: "text is not set",
	ThreadNotFound: "thread is not found",
	Unknown: "unknown error",
	UnknownVariant: "unknown variant",
	WebTokenIsNotSet: "web token is not set",
}

// Messages.
msg = {
	GenericErrorPrefix: "Error: ",
	Redirecting: "Redirecting. Please wait ...",
}

// User role names.
userRole = {
	Author: "author",
	Writer: "writer",
	Reader: "reader",
	Logging: "logging",
}

class ApiRequest {
	constructor(action, parameters) {
		this.Action = action;
		this.Parameters = parameters;
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

// Request parameters.

class Parameters_AddForum {
	constructor(parent, name) {
		this.SectionId = parent;
		this.Name = name;
	}
}

class Parameters_AddMessage {
	constructor(parent, text) {
		this.ThreadId = parent;
		this.Text = text;
	}
}

class Parameters_AddNotification {
	constructor(userId, text) {
		this.UserId = userId;
		this.Text = text;
	}
}

class Parameters_AddSection {
	constructor(parent, name) {
		this.Parent = parent;
		this.Name = name;
	}
}

class Parameters_AddSubscription {
	constructor(threadId, userId) {
		this.ThreadId = threadId;
		this.UserId = userId;
	}
}

class Parameters_AddThread {
	constructor(parent, name) {
		this.ForumId = parent;
		this.Name = name;
	}
}

class Parameters_ApproveAndRegisterUser {
	constructor(email) {
		this.Email = email;
	}
}

class Parameters_BanUser {
	constructor(userId) {
		this.UserId = userId;
	}
}

class Parameters_ChangeEmail1 {
	constructor(stepN, newEmail) {
		this.StepN = stepN;
		this.NewEmail = newEmail;
	}
}

class Parameters_ChangeEmail2 {
	constructor(stepN, requestId, authChallengeResponse, verificationCodeOld, verificationCodeNew, captchaAnswer) {
		this.StepN = stepN;
		this.RequestId = requestId;
		this.AuthChallengeResponse = authChallengeResponse;
		this.VerificationCodeOld = verificationCodeOld;
		this.VerificationCodeNew = verificationCodeNew;
		this.CaptchaAnswer = captchaAnswer;
	}
}

class Parameters_ChangeForumName {
	constructor(forumId, name) {
		this.ForumId = forumId;
		this.Name = name;
	}
}

class Parameters_ChangeForumSection {
	constructor(forumId, sectionId) {
		this.ForumId = forumId;
		this.SectionId = sectionId; // Parent.
	}
}

class Parameters_ChangeMessageText {
	constructor(messageId, text) {
		this.MessageId = messageId;
		this.Text = text;
	}
}

class Parameters_ChangeMessageThread {
	constructor(messageId, threadId) {
		this.MessageId = messageId;
		this.ThreadId = threadId; // Parent.
	}
}

class Parameters_ChangePwd1 {
	constructor(stepN, newPassword) {
		this.StepN = stepN;
		this.NewPassword = newPassword;
	}
}

class Parameters_ChangePwd2 {
	constructor(stepN, requestId, authChallengeResponse, verificationCode, captchaAnswer) {
		this.StepN = stepN;
		this.RequestId = requestId;
		this.AuthChallengeResponse = authChallengeResponse;
		this.VerificationCode = verificationCode;
		this.CaptchaAnswer = captchaAnswer;
	}
}

class Parameters_ChangeSectionName {
	constructor(sectionId, name) {
		this.SectionId = sectionId;
		this.Name = name;
	}
}

class Parameters_ChangeSectionParent {
	constructor(sectionId, parent) {
		this.SectionId = sectionId;
		this.Parent = parent;
	}
}

class Parameters_ChangeThreadForum {
	constructor(threadId, forumId) {
		this.ThreadId = threadId;
		this.ForumId = forumId; // Parent.
	}
}

class Parameters_ChangeThreadName {
	constructor(threadId, name) {
		this.ThreadId = threadId;
		this.Name = name;
	}
}

class Parameters_DeleteForum {
	constructor(forumId) {
		this.ForumId = forumId;
	}
}

class Parameters_DeleteMessage {
	constructor(messageId) {
		this.MessageId = messageId;
	}
}

class Parameters_DeleteNotification {
	constructor(notificationId) {
		this.NotificationId = notificationId;
	}
}

class Parameters_DeleteSection {
	constructor(sectionId) {
		this.SectionId = sectionId;
	}
}

class Parameters_DeleteSelfSubscription {
	constructor(threadId) {
		this.ThreadId = threadId;
	}
}

class Parameters_DeleteThread {
	constructor(threadId) {
		this.ThreadId = threadId;
	}
}

class Parameters_GetLatestMessageOfThread {
	constructor(threadId) {
		this.ThreadId = threadId;
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

class Parameters_GetMessage {
	constructor(messageId) {
		this.MessageId = messageId;
	}
}

class Parameters_GetNotificationsOnPage {
	constructor(page) {
		this.Page = page;
	}
}

class Parameters_GetSelfSubscriptionsOnPage {
	constructor(page) {
		this.Page = page;
	}
}

class Parameters_GetThread {
	constructor(threadId) {
		this.ThreadId = threadId;
	}
}

class Parameters_GetThreadNamesByIds {
	constructor(threadIds) {
		this.ThreadIds = threadIds;
	}
}

class Parameters_GetUserName {
	constructor(userId) {
		this.UserId = userId;
	}
}

class Parameters_GetUserSession {
	constructor(userId) {
		this.UserId = userId;
	}
}

class Parameters_IsSelfSubscribed {
	constructor(threadId) {
		this.ThreadId = threadId;
	}
}

class Parameters_IsUserLoggedIn {
	constructor(userId) {
		this.UserId = userId;
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

class Parameters_LogUserOutA {
	constructor(userId) {
		this.UserId = userId;
	}
}

class Parameters_MarkNotificationAsRead {
	constructor(notificationId) {
		this.NotificationId = notificationId;
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

class Parameters_RejectRegistrationRequest {
	constructor(registrationRequestId) {
		this.RegistrationRequestId = registrationRequestId;
	}
}

class Parameters_SetUserRoleAuthor {
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

class Parameters_SetUserRoleWriter {
	constructor(userId, isRoleEnabled) {
		this.UserId = userId;
		this.IsRoleEnabled = isRoleEnabled;
	}
}

class Parameters_UnbanUser {
	constructor(userId) {
		this.UserId = userId;
	}
}

class Parameters_ViewUserParameters {
	constructor(userId) {
		this.UserId = userId;
	}
}

// Various models.

class Settings {
	constructor(version, productVersion, siteName, siteDomain, captchaFolder,
				sessionMaxDuration, messageEditTime, pageSize, apiFolder,
				publicSettingsFileName, isFrontEndEnabled,
				frontEndStaticFilesFolder) {
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

class SectionNode {
	constructor(section, level) {
		this.Section = section;
		this.Level = level;
	}
}

class Subscriptions {
	constructor(userId, threadIds, threadNames, pageNumber, totalPages,
				totalSubscriptions) {
		this.UserId = userId;
		this.ThreadIds = threadIds;
		this.ThreadNames = threadNames;
		this.PageNumber = pageNumber;
		this.TotalPages = totalPages;
		this.TotalSubscriptions = totalSubscriptions;
	}
}

// API methods.

async function sendApiRequest(data) {
	let settings = getSettings();
	let url = rootPath + settings.ApiFolder;
	let ri = {
		method: "POST",
		body: JSON.stringify(data)
	};
	let resp = await fetch(url, ri);
	let result;
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

async function sendApiRequestAndReturnJson(reqData) {
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
	return await sendApiRequestAndReturnJson(reqData);
}

async function addMessage(parent, text) {
	let params = new Parameters_AddMessage(parent, text);
	let reqData = new ApiRequest(actionName.AddMessage, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function addNotification(userId, text) {
	let params = new Parameters_AddNotification(userId, text);
	let reqData = new ApiRequest(actionName.AddNotification, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function addSection(parent, name) {
	let params = new Parameters_AddSection(parent, name);
	let reqData = new ApiRequest(actionName.AddSection, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function addSubscription(threadId, userId) {
	let params = new Parameters_AddSubscription(threadId, userId);
	let reqData = new ApiRequest(actionName.AddSubscription, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function addThread(parent, name) {
	let params = new Parameters_AddThread(parent, name);
	let reqData = new ApiRequest(actionName.AddThread, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function approveAndRegisterUser(email) {
	let params = new Parameters_ApproveAndRegisterUser(email);
	let reqData = new ApiRequest(actionName.ApproveAndRegisterUser, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function banUser(userId) {
	let params = new Parameters_BanUser(userId);
	let reqData = new ApiRequest(actionName.BanUser, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function changeEmail1(stepN, newEmail) {
	let params = new Parameters_ChangeEmail1(stepN, newEmail);
	let reqData = new ApiRequest(actionName.ChangeEmail, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function changeEmail2(stepN, requestId, authChallengeResponse, verificationCodeOld, verificationCodeNew, captchaAnswer) {
	let params = new Parameters_ChangeEmail2(stepN, requestId, authChallengeResponse, verificationCodeOld, verificationCodeNew, captchaAnswer);
	let reqData = new ApiRequest(actionName.ChangeEmail, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function changeForumName(forumId, name) {
	let params = new Parameters_ChangeForumName(forumId, name);
	let reqData = new ApiRequest(actionName.ChangeForumName, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function changeForumSection(forumId, newParent) {
	let params = new Parameters_ChangeForumSection(forumId, newParent);
	let reqData = new ApiRequest(actionName.ChangeForumSection, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function changeMessageText(messageId, text) {
	let params = new Parameters_ChangeMessageText(messageId, text);
	let reqData = new ApiRequest(actionName.ChangeMessageText, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function changeMessageThread(messageId, newParent) {
	let params = new Parameters_ChangeMessageThread(messageId, newParent);
	let reqData = new ApiRequest(actionName.ChangeMessageThread, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function changePwd1(stepN, newPassword) {
	let params = new Parameters_ChangePwd1(stepN, newPassword);
	let reqData = new ApiRequest(actionName.ChangePwd, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function changePwd2(stepN, requestId, authChallengeResponse, verificationCode, captchaAnswer) {
	let params = new Parameters_ChangePwd2(stepN, requestId, authChallengeResponse, verificationCode, captchaAnswer);
	let reqData = new ApiRequest(actionName.ChangePwd, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function changeSectionName(sectionId, name) {
	let params = new Parameters_ChangeSectionName(sectionId, name);
	let reqData = new ApiRequest(actionName.ChangeSectionName, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function changeSectionParent(sectionId, newParent) {
	let params = new Parameters_ChangeSectionParent(sectionId, newParent);
	let reqData = new ApiRequest(actionName.ChangeSectionParent, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function changeThreadForum(threadId, newParent) {
	let params = new Parameters_ChangeThreadForum(threadId, newParent);
	let reqData = new ApiRequest(actionName.ChangeThreadForum, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function changeThreadName(threadId, name) {
	let params = new Parameters_ChangeThreadName(threadId, name);
	let reqData = new ApiRequest(actionName.ChangeThreadName, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function countSelfSubscriptions() {
	let reqData = new ApiRequest(actionName.CountSelfSubscriptions, {});
	return await sendApiRequestAndReturnJson(reqData);
}

async function countUnreadNotifications() {
	let reqData = new ApiRequest(actionName.CountUnreadNotifications, {});
	return await sendApiRequestAndReturnJson(reqData);
}

async function deleteForum(forumId) {
	let params = new Parameters_DeleteForum(forumId);
	let reqData = new ApiRequest(actionName.DeleteForum, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function deleteMessage(messageId) {
	let params = new Parameters_DeleteMessage(messageId);
	let reqData = new ApiRequest(actionName.DeleteMessage, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function deleteNotification(notificationId) {
	let params = new Parameters_DeleteNotification(notificationId);
	let reqData = new ApiRequest(actionName.DeleteNotification, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function deleteSection(sectionId) {
	let params = new Parameters_DeleteSection(sectionId);
	let reqData = new ApiRequest(actionName.DeleteSection, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function deleteSelfSubscription(threadId) {
	let params = new Parameters_DeleteSelfSubscription(threadId);
	let reqData = new ApiRequest(actionName.DeleteSelfSubscription, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function deleteThread(threadId) {
	let params = new Parameters_DeleteThread(threadId);
	let reqData = new ApiRequest(actionName.DeleteThread, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function getAllNotifications() {
	let reqData = new ApiRequest(actionName.GetAllNotifications, {});
	return await sendApiRequestAndReturnJson(reqData);
}

async function getLatestMessageOfThread(threadId) {
	let params = new Parameters_GetLatestMessageOfThread(threadId);
	let reqData = new ApiRequest(actionName.GetLatestMessageOfThread, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function getListOfAllUsers(pageN) {
	let params = new Parameters_GetListOfAllUsers(pageN);
	let reqData = new ApiRequest(actionName.GetListOfAllUsers, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function getListOfRegistrationsReadyForApproval(pageN) {
	let params = new Parameters_GetListOfRegistrationsReadyForApproval(pageN);
	let reqData = new ApiRequest(actionName.GetListOfRegistrationsReadyForApproval, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function getMessage(messageId) {
	let params = new Parameters_GetMessage(messageId);
	let reqData = new ApiRequest(actionName.GetMessage, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function getNotificationsOnPage(page) {
	let params = new Parameters_GetNotificationsOnPage(page);
	let reqData = new ApiRequest(actionName.GetNotificationsOnPage, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function getSelfRoles() {
	let reqData = new ApiRequest(actionName.GetSelfRoles, {});
	return await sendApiRequestAndReturnJson(reqData);
}

async function getSelfSubscriptions() {
	let reqData = new ApiRequest(actionName.GetSelfSubscriptions, {});
	return await sendApiRequestAndReturnJson(reqData);
}

async function getSelfSubscriptionsOnPage(page) {
	let params = new Parameters_GetSelfSubscriptionsOnPage(page);
	let reqData = new ApiRequest(actionName.GetSelfSubscriptionsOnPage, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function getThread(threadId) {
	let params = new Parameters_GetThread(threadId);
	let reqData = new ApiRequest(actionName.GetThread, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function getThreadNamesByIds(threadIds) {
	let params = new Parameters_GetThreadNamesByIds(threadIds);
	let reqData = new ApiRequest(actionName.GetThreadNamesByIds, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function getUserName(userId) {
	let params = new Parameters_GetUserName(userId);
	let reqData = new ApiRequest(actionName.GetUserName, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function getUserSession(userId) {
	let params = new Parameters_GetUserSession(userId);
	let reqData = new ApiRequest(actionName.GetUserSession, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function isSelfSubscribed(threadId) {
	let params = new Parameters_IsSelfSubscribed(threadId);
	let reqData = new ApiRequest(actionName.IsSelfSubscribed, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function isUserLoggedIn(userId) {
	let params = new Parameters_IsUserLoggedIn(userId);
	let reqData = new ApiRequest(actionName.IsUserLoggedIn, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function listForumAndThreadsOnPage(forumId, page) {
	let params = new Parameters_ListForumAndThreadsOnPage(forumId, page);
	let reqData = new ApiRequest(actionName.ListForumAndThreadsOnPage, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function listSectionsAndForums() {
	let reqData = new ApiRequest(actionName.ListSectionsAndForums, {});
	return await sendApiRequestAndReturnJson(reqData);
}

async function listThreadAndMessagesOnPage(threadId, page) {
	let params = new Parameters_ListThreadAndMessagesOnPage(threadId, page);
	let reqData = new ApiRequest(actionName.ListThreadAndMessagesOnPage, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function logUserIn1(stepN, email) {
	let params = new Parameters_LogIn1(stepN, email);
	let reqData = new ApiRequest(actionName.LogUserIn, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function logUserIn2(stepN, email, requestId, captchaAnswer, authChallengeResponse) {
	let params = new Parameters_LogIn2(stepN, email, requestId, captchaAnswer, authChallengeResponse);
	let reqData = new ApiRequest(actionName.LogUserIn, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function logUserIn3(stepN, email, requestId, verificationCode) {
	let params = new Parameters_LogIn3(stepN, email, requestId, verificationCode);
	let reqData = new ApiRequest(actionName.LogUserIn, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function logUserOut1() {
	let reqData = new ApiRequest(actionName.LogUserOut, {});
	return await sendApiRequestAndReturnJson(reqData);
}

async function logUserOutA(userId) {
	let params = new Parameters_LogUserOutA(userId);
	let reqData = new ApiRequest(actionName.LogUserOutA, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function markNotificationAsRead(notificationId) {
	let params = new Parameters_MarkNotificationAsRead(notificationId);
	let reqData = new ApiRequest(actionName.MarkNotificationAsRead, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function moveForumUp(forumId) {
	let params = new Parameters_MoveForumUp(forumId);
	let reqData = new ApiRequest(actionName.MoveForumUp, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function moveForumDown(forumId) {
	let params = new Parameters_MoveForumDown(forumId);
	let reqData = new ApiRequest(actionName.MoveForumDown, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function moveSectionUp(sectionId) {
	let params = new Parameters_MoveSectionUp(sectionId);
	let reqData = new ApiRequest(actionName.MoveSectionUp, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function moveSectionDown(sectionId) {
	let params = new Parameters_MoveSectionDown(sectionId);
	let reqData = new ApiRequest(actionName.MoveSectionDown, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function moveThreadUp(threadId) {
	let params = new Parameters_MoveThreadUp(threadId);
	let reqData = new ApiRequest(actionName.MoveThreadUp, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function moveThreadDown(threadId) {
	let params = new Parameters_MoveThreadDown(threadId);
	let reqData = new ApiRequest(actionName.MoveThreadDown, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function registerUser1(stepN, email) {
	let params = new Parameters_RegisterUser1(stepN, email);
	let reqData = new ApiRequest(actionName.RegisterUser, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function registerUser2(stepN, email, verificationCode) {
	let params = new Parameters_RegisterUser2(stepN, email, verificationCode);
	let reqData = new ApiRequest(actionName.RegisterUser, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function registerUser3(stepN, email, verificationCode, name, pwd) {
	let params = new Parameters_RegisterUser3(stepN, email, verificationCode, name, pwd);
	let reqData = new ApiRequest(actionName.RegisterUser, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function rejectRegistrationRequest(registrationRequestId) {
	let params = new Parameters_RejectRegistrationRequest(registrationRequestId);
	let reqData = new ApiRequest(actionName.RejectRegistrationRequest, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function setUserRoleAuthor(userId, roleValue) {
	let params = new Parameters_SetUserRoleAuthor(userId, roleValue);
	let reqData = new ApiRequest(actionName.SetUserRoleAuthor, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function setUserRoleReader(userId, roleValue) {
	let params = new Parameters_SetUserRoleReader(userId, roleValue);
	let reqData = new ApiRequest(actionName.SetUserRoleReader, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function setUserRoleWriter(userId, roleValue) {
	let params = new Parameters_SetUserRoleWriter(userId, roleValue);
	let reqData = new ApiRequest(actionName.SetUserRoleWriter, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function unbanUser(userId) {
	let params = new Parameters_UnbanUser(userId);
	let reqData = new ApiRequest(actionName.UnbanUser, params);
	return await sendApiRequestAndReturnJson(reqData);
}

async function viewUserParameters(userId) {
	let params = new Parameters_ViewUserParameters(userId);
	let reqData = new ApiRequest(actionName.ViewUserParameters, params);
	return await sendApiRequestAndReturnJson(reqData);
}

// Various API helpers.

async function getSelfSubscriptionsPaginated(pageNumber) {
	let resp = await getSelfSubscriptionsOnPage(pageNumber);
	if (resp == null) {
		return null;
	}
	let result = new Subscriptions(
		resp.result.sop.subscriber,
		resp.result.sop.subscriptions,
		[],
		pageNumber,
		resp.result.sop.totalPages,
		resp.result.sop.totalSubscriptions,
	)

	if ((result.ThreadIds != null) && (result.ThreadIds.length > 0)) {
		resp = await getThreadNamesByIds(result.ThreadIds);
		if (resp == null) {
			return null;
		}
		result.ThreadNames = resp.result.threadNames;
	}
	return result;
}

// Common methods for settings.

async function fetchSettings() {
	let data = await fetch(rootPath + settingsPath);
	return await data.json();
}

function respToSettings(resp) {
	return new Settings(
		resp.version,
		resp.productVersion,
		resp.siteName,
		resp.siteDomain,
		resp.captchaFolder,
		resp.sessionMaxDuration,
		resp.messageEditTime,
		resp.pageSize,
		resp.apiFolder,
		resp.publicSettingsFileName,
		resp.isFrontEndEnabled,
		resp.frontEndStaticFilesFolder,
	);
}

function checkSettings(s) {
	if ((s.PublicSettingsFileName !== settingsPath)) {
		console.error(err.Settings);
		return false;
	}
	return true;
}

async function updateSettingsIfNeeded() {
	if (isSettingsUpdateNeeded()) {
		return await updateSettings();
	}
	return true;
}

async function updateSettings() {
	let resp = await fetchSettings();
	let s = respToSettings(resp);
	if (!checkSettings(s)) {
		return false;
	}
	console.info('Received settings. Version: ' + s.Version.toString() + ".");

	// Save the settings for future usage.
	saveSettings(s);
	return true;
}

// Various common functions.

function booleanToString(b) {
	if (b === true) {
		return "Yes";
	}
	if (b === false) {
		return "No";
	}
	console.error("boolToText:", b);
	return null;
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

function getCurrentTimestamp() {
	return Math.floor(Date.now() / 1000);
}

async function sleep(ms) {
	await new Promise(r => setTimeout(r, ms));
}

function addTimeSec(t, deltaSec) {
	return new Date(t.getTime() + deltaSec * 1000);
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

async function reloadPage(wait) {
	if (wait) {
		await sleep(redirectDelay * 1000);
	}
	location.reload();
}

// redirectToRelativePath redirects to a page with a relative path.
// E.g., if a relative path is 'x', then URL is '/x', where '/' is a front end
// root.
async function redirectToRelativePath(wait, relPath) {
	let url = rootPath + relPath;
	await redirectPage(wait, url);
}

async function redirectToMainPage(wait) {
	await redirectPage(wait, rootPath);
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

function escapeHtml(text) {
	let div = document.createElement('div');
	div.textContent = text;
	return div.innerHTML;
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

// URL composition.

function composeCaptchaImageUrl(captchaId) {
	let settings = getSettings();
	let captchaPath = rootPath + settings.CaptchaFolder;
	return captchaPath + qp.Prefix + qpn.Id + "=" + captchaId;
}

function composeUrlForEntity(entityName, entityId) {
	return qp.Prefix + entityName + "&" + qpn.Id + "=" + entityId;
}

function composeUrlForEntityPage(entityName, entityId, page) {
	return qp.Prefix + entityName + "&" + qpn.Id + "=" + entityId + "&" + qpn.Page + "=" + page;
}

function composeUrlForForum(forumId) {
	return composeUrlForEntity(qpn.Forum, forumId);
}

function composeUrlForForumPage(forumId, page) {
	return composeUrlForEntityPage(qpn.Forum, forumId, page);
}

function composeUrlForFuncPage(func, page) {
	return qp.Prefix + func + "&" + qpn.Page + "=" + page;
}

function composeUrlForMessage(messageId) {
	return composeUrlForEntity(qpn.Message, messageId);
}

function composeUrlForNotificationsPage(page) {
	return composeUrlForFuncPage(qpn.Notifications, page);
}

function composeUrlForSection(sectionId) {
	return composeUrlForEntity(qpn.Section, sectionId);
}

function composeUrlForSubscriptionsPage(page) {
	return composeUrlForFuncPage(qpn.SubscriptionsPage, page);
}

function composeUrlForThread(threadId) {
	return composeUrlForEntity(qpn.Thread, threadId);
}

function composeUrlForThreadPage(threadId, page) {
	return composeUrlForEntityPage(qpn.Thread, threadId, page);
}
