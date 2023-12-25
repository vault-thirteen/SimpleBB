package mm

import (
	"fmt"
	"math"
	"sync"
	"time"

	js "github.com/osamingo/jsonrpc/v2"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/models"
	ul "github.com/vault-thirteen/SimpleBB/pkg/UidList"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
)

// RPC functions.

// Section.

// addSection inserts a new section as a root section or as a sub-section.
func (srv *Server) addSection(p *mm.AddSectionParams) (result *mm.AddSectionResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if len(p.Name) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_SectionNameIsNotSet, Message: RpcErrorMsg_SectionNameIsNotSet}
	}

	// If parent is not set, the new section is a root section.
	// Only a single root section may exist.
	var err error
	var n int
	if p.Parent == nil {
		n, err = srv.dbo.CountRootSections()
		if err != nil {
			return nil, srv.databaseError(err)
		}

		if n > 0 {
			return nil, &js.Error{Code: RpcErrorCode_RootSectionAlreadyExists, Message: RpcErrorMsg_RootSectionAlreadyExists}
		}

		var insertedSectionId int64
		insertedSectionId, err = srv.dbo.InsertNewSection(p.Parent, p.Name, userRoles.UserId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		result = &mm.AddSectionResult{
			SectionId: uint(insertedSectionId),
		}

		return result, nil
	}

	// Insert a sub-section.
	// Ensure that a parent exists.
	n, err = srv.dbo.CountSectionsById(*p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_SectionIsNotFound, Message: RpcErrorMsg_SectionIsNotFound}
	}

	// Check compatibility.
	var childType byte
	childType, err = srv.dbo.GetSectionChildTypeById(*p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if childType == mm.ChildTypeForum {
		return nil, &js.Error{Code: RpcErrorCode_IncompatibleChildType, Message: RpcErrorMsg_IncompatibleChildType}
	}

	if childType == mm.ChildTypeNone {
		err = srv.dbo.SetSectionChildTypeById(*p.Parent, mm.ChildTypeSection)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	// Insert a section and link it with its parent.
	var parentChildren *ul.UidList
	parentChildren, err = srv.dbo.GetSectionChildrenById(*p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var insertedSectionId int64
	insertedSectionId, err = srv.dbo.InsertNewSection(p.Parent, p.Name, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = parentChildren.AddItem(uint(insertedSectionId), false)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.dbo.SetSectionChildrenById(*p.Parent, parentChildren)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.AddSectionResult{
		SectionId: uint(insertedSectionId),
	}

	return result, nil
}

// changeSectionName renames a section.
func (srv *Server) changeSectionName(p *mm.ChangeSectionNameParams) (result *mm.ChangeSectionNameResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.SectionId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_SectionIdIsNotSet, Message: RpcErrorMsg_SectionIdIsNotSet}
	}

	if len(p.Name) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_SectionNameIsNotSet, Message: RpcErrorMsg_SectionNameIsNotSet}
	}

	var err error
	err = srv.dbo.SetSectionNameById(p.SectionId, p.Name, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ChangeSectionNameResult{
		OK: true,
	}

	return result, nil
}

// changeSectionParent moves a section from an old parent to a new parent.
func (srv *Server) changeSectionParent(p *mm.ChangeSectionParentParams) (result *mm.ChangeSectionParentResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.SectionId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_SectionIdIsNotSet, Message: RpcErrorMsg_SectionIdIsNotSet}
	}

	if p.Parent == 0 {
		return nil, &js.Error{Code: RpcErrorCode_SectionIdIsNotSet, Message: RpcErrorMsg_SectionIdIsNotSet}
	}

	// Ensure that an old parent exists.
	var oldParent *uint
	var err error
	oldParent, err = srv.dbo.GetSectionParentById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if oldParent == nil {
		return nil, &js.Error{Code: RpcErrorCode_RootSectionCanNotBeMoved, Message: RpcErrorMsg_RootSectionCanNotBeMoved}
	}

	var n int
	n, err = srv.dbo.CountSectionsById(*oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_SectionIsNotFound, Message: RpcErrorMsg_SectionIsNotFound}
	}

	// Ensure that a new parent exists.
	n, err = srv.dbo.CountSectionsById(p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_SectionIsNotFound, Message: RpcErrorMsg_SectionIsNotFound}
	}

	// Check compatibility.
	var childType byte
	childType, err = srv.dbo.GetSectionChildTypeById(p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if childType == mm.ChildTypeForum {
		return nil, &js.Error{Code: RpcErrorCode_IncompatibleChildType, Message: RpcErrorMsg_IncompatibleChildType}
	}

	if childType == mm.ChildTypeNone {
		err = srv.dbo.SetSectionChildTypeById(p.Parent, mm.ChildTypeSection)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	// Update the moved section.
	err = srv.dbo.SetSectionParentById(p.SectionId, p.Parent, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the new link.
	var childrenR *ul.UidList
	childrenR, err = srv.dbo.GetSectionChildrenById(p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = childrenR.AddItem(p.SectionId, false)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.dbo.SetSectionChildrenById(p.Parent, childrenR)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the old link.
	var childrenL *ul.UidList
	childrenL, err = srv.dbo.GetSectionChildrenById(*oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = childrenL.RemoveItem(p.SectionId)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.dbo.SetSectionChildrenById(*oldParent, childrenL)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Clear the child type if the old parent becomes empty.
	if childrenL.Size() == 0 {
		err = srv.dbo.SetSectionChildTypeById(*oldParent, mm.ChildTypeNone)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	result = &mm.ChangeSectionParentResult{
		OK: true,
	}

	return result, nil
}

// getSection reads a section.
func (srv *Server) getSection(p *mm.GetSectionParams) (result *mm.GetSectionResult, jerr *js.Error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsReader {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.SectionId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_SectionIdIsNotSet, Message: RpcErrorMsg_SectionIdIsNotSet}
	}

	// Read the section.
	var section *mm.Section
	var err error
	section, err = srv.dbo.GetSectionById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.GetSectionResult{
		Section: section,
	}

	return result, nil
}

// deleteSection removes a section.
func (srv *Server) deleteSection(p *mm.DeleteSectionParams) (result *mm.DeleteSectionResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.SectionId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_SectionIdIsNotSet, Message: RpcErrorMsg_SectionIdIsNotSet}
	}

	// Read the section.
	var section *mm.Section
	var err error
	section, err = srv.dbo.GetSectionById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var isRootSection = false
	if section.Parent == nil {
		isRootSection = true
	}

	// Check for children.
	if section.Children.Size() > 0 {
		return nil, &js.Error{Code: RpcErrorCode_SectionHasChildren, Message: RpcErrorMsg_SectionHasChildren}
	}

	// Update the link.
	if !isRootSection {
		var linkSections *ul.UidList
		linkSections, err = srv.dbo.GetSectionChildrenById(*section.Parent)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		err = linkSections.RemoveItem(p.SectionId)
		if err != nil {
			srv.logError(err)
			return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
		}

		err = srv.dbo.SetSectionChildrenById(*section.Parent, linkSections)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// Clear the child type if the old parent becomes empty.
		if linkSections.Size() == 0 {
			err = srv.dbo.SetSectionChildTypeById(*section.Parent, mm.ChildTypeNone)
			if err != nil {
				return nil, srv.databaseError(err)
			}
		}
	}

	// Delete the section.
	err = srv.dbo.DeleteSectionById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.DeleteSectionResult{
		OK: true,
	}

	return result, nil
}

// Forum.

// addForum inserts a new forum into a section.
func (srv *Server) addForum(p *mm.AddForumParams) (result *mm.AddForumResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.SectionId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_SectionIdIsNotSet, Message: RpcErrorMsg_SectionIdIsNotSet}
	}

	if len(p.Name) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumNameIsNotSet, Message: RpcErrorMsg_ForumNameIsNotSet}
	}

	// Ensure that a section exists.
	n, err := srv.dbo.CountSectionsById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_SectionIsNotFound, Message: RpcErrorMsg_SectionIsNotFound}
	}

	// Check compatibility.
	var childType byte
	childType, err = srv.dbo.GetSectionChildTypeById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if childType == mm.ChildTypeSection {
		return nil, &js.Error{Code: RpcErrorCode_IncompatibleChildType, Message: RpcErrorMsg_IncompatibleChildType}
	}

	if childType == mm.ChildTypeNone {
		err = srv.dbo.SetSectionChildTypeById(p.SectionId, mm.ChildTypeForum)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	// Insert a forum and link it with its section.
	var parentChildren *ul.UidList
	parentChildren, err = srv.dbo.GetSectionChildrenById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var insertedForumId int64
	insertedForumId, err = srv.dbo.InsertNewForum(p.SectionId, p.Name, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = parentChildren.AddItem(uint(insertedForumId), false)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.dbo.SetSectionChildrenById(p.SectionId, parentChildren)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.AddForumResult{
		ForumId: uint(insertedForumId),
	}

	return result, nil
}

// changeForumName renames a forum.
func (srv *Server) changeForumName(p *mm.ChangeForumNameParams) (result *mm.ChangeForumNameResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ForumId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIdIsNotSet, Message: RpcErrorMsg_ForumIdIsNotSet}
	}

	if len(p.Name) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumNameIsNotSet, Message: RpcErrorMsg_ForumNameIsNotSet}
	}

	var err error
	err = srv.dbo.SetForumNameById(p.ForumId, p.Name, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ChangeForumNameResult{
		OK: true,
	}

	return result, nil
}

// changeForumSection moves a forum from an old section to a new section.
func (srv *Server) changeForumSection(p *mm.ChangeForumSectionParams) (result *mm.ChangeForumSectionResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ForumId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIdIsNotSet, Message: RpcErrorMsg_ForumIdIsNotSet}
	}

	if p.SectionId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_SectionIdIsNotSet, Message: RpcErrorMsg_SectionIdIsNotSet}
	}

	// Ensure that an old section exists.
	var oldParent uint
	var err error
	oldParent, err = srv.dbo.GetForumSectionById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var n int
	n, err = srv.dbo.CountSectionsById(oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_SectionIsNotFound, Message: RpcErrorMsg_SectionIsNotFound}
	}

	// Ensure that a new section exists.
	n, err = srv.dbo.CountSectionsById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_SectionIsNotFound, Message: RpcErrorMsg_SectionIsNotFound}
	}

	// Check compatibility.
	var childType byte
	childType, err = srv.dbo.GetSectionChildTypeById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if childType == mm.ChildTypeSection {
		return nil, &js.Error{Code: RpcErrorCode_IncompatibleChildType, Message: RpcErrorMsg_IncompatibleChildType}
	}

	if childType == mm.ChildTypeNone {
		err = srv.dbo.SetSectionChildTypeById(p.SectionId, mm.ChildTypeForum)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	// Update the moved forum.
	err = srv.dbo.SetForumSectionById(p.ForumId, p.SectionId, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the new link.
	var childrenR *ul.UidList
	childrenR, err = srv.dbo.GetSectionChildrenById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = childrenR.AddItem(p.ForumId, false)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.dbo.SetSectionChildrenById(p.SectionId, childrenR)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the old link.
	var childrenL *ul.UidList
	childrenL, err = srv.dbo.GetSectionChildrenById(oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = childrenL.RemoveItem(p.ForumId)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.dbo.SetSectionChildrenById(oldParent, childrenL)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Clear the child type if the old section becomes empty.
	if childrenL.Size() == 0 {
		err = srv.dbo.SetSectionChildTypeById(oldParent, mm.ChildTypeNone)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	result = &mm.ChangeForumSectionResult{
		OK: true,
	}

	return result, nil
}

// getForum reads a forum.
func (srv *Server) getForum(p *mm.GetForumParams) (result *mm.GetForumResult, jerr *js.Error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsReader {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ForumId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIdIsNotSet, Message: RpcErrorMsg_ForumIdIsNotSet}
	}

	// Read the forum.
	var forum *mm.Forum
	var err error
	forum, err = srv.dbo.GetForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.GetForumResult{
		Forum: forum,
	}

	return result, nil
}

// deleteForum removes a forum.
func (srv *Server) deleteForum(p *mm.DeleteForumParams) (result *mm.DeleteForumResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ForumId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIdIsNotSet, Message: RpcErrorMsg_ForumIdIsNotSet}
	}

	// Read the forum.
	var forum *mm.Forum
	var err error
	forum, err = srv.dbo.GetForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Check for threads.
	if forum.Threads.Size() > 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumHasThreads, Message: RpcErrorMsg_ForumHasThreads}
	}

	// Update the link.
	var linkChildren *ul.UidList
	linkChildren, err = srv.dbo.GetSectionChildrenById(forum.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = linkChildren.RemoveItem(p.ForumId)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.dbo.SetSectionChildrenById(forum.SectionId, linkChildren)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Clear the child type if the old parent becomes empty.
	if linkChildren.Size() == 0 {
		err = srv.dbo.SetSectionChildTypeById(forum.SectionId, mm.ChildTypeNone)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	// Delete the forum.
	err = srv.dbo.DeleteForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.DeleteForumResult{
		OK: true,
	}

	return result, nil
}

// Thread.

// addThread inserts a new thread into a forum.
func (srv *Server) addThread(p *mm.AddThreadParams) (result *mm.AddThreadResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAuthor {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ForumId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIdIsNotSet, Message: RpcErrorMsg_ForumIdIsNotSet}
	}

	if len(p.Name) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadNameIsNotSet, Message: RpcErrorMsg_ThreadNameIsNotSet}
	}

	// Ensure that a forum exists.
	var err error
	var n int
	n, err = srv.dbo.CountForumsById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIsNotFound, Message: RpcErrorMsg_ForumIsNotFound}
	}

	// Insert a thread and link it with its forum.
	var parentThreads *ul.UidList
	parentThreads, err = srv.dbo.GetForumThreadsById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var insertedThreadId int64
	insertedThreadId, err = srv.dbo.InsertNewThread(p.ForumId, p.Name, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = parentThreads.AddItem(uint(insertedThreadId), srv.settings.SystemSettings.NewThreadsAtTop)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.dbo.SetForumThreadsById(p.ForumId, parentThreads)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.AddThreadResult{
		ThreadId: uint(insertedThreadId),
	}

	return result, nil
}

// changeThreadName renames a thread.
func (srv *Server) changeThreadName(p *mm.ChangeThreadNameParams) (result *mm.ChangeThreadNameResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ThreadId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIdIsNotSet, Message: RpcErrorMsg_ThreadIdIsNotSet}
	}

	if len(p.Name) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadNameIsNotSet, Message: RpcErrorMsg_ThreadNameIsNotSet}
	}

	var err error
	err = srv.dbo.SetThreadNameById(p.ThreadId, p.Name, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ChangeThreadNameResult{
		OK: true,
	}

	return result, nil
}

// changeThreadForum moves a thread from an old forum to a new forum.
func (srv *Server) changeThreadForum(p *mm.ChangeThreadForumParams) (result *mm.ChangeThreadForumResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ThreadId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIdIsNotSet, Message: RpcErrorMsg_ThreadIdIsNotSet}
	}

	if p.ForumId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIdIsNotSet, Message: RpcErrorMsg_ForumIdIsNotSet}
	}

	// Ensure that an old parent exists.
	var oldParent uint
	var err error
	oldParent, err = srv.dbo.GetThreadForumById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var n int
	n, err = srv.dbo.CountForumsById(oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIsNotFound, Message: RpcErrorMsg_ForumIsNotFound}
	}

	// Ensure that a new parent exists.
	n, err = srv.dbo.CountForumsById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIsNotFound, Message: RpcErrorMsg_ForumIsNotFound}
	}

	// Update the moved thread.
	err = srv.dbo.SetThreadForumById(p.ThreadId, p.ForumId, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the new link.
	var threadsR *ul.UidList
	threadsR, err = srv.dbo.GetForumThreadsById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = threadsR.AddItem(p.ThreadId, false)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.dbo.SetForumThreadsById(p.ForumId, threadsR)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the old link.
	var threadsL *ul.UidList
	threadsL, err = srv.dbo.GetForumThreadsById(oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = threadsL.RemoveItem(p.ThreadId)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.dbo.SetForumThreadsById(oldParent, threadsL)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ChangeThreadForumResult{
		OK: true,
	}

	return result, nil
}

// getThread reads a thread.
func (srv *Server) getThread(p *mm.GetThreadParams) (result *mm.GetThreadResult, jerr *js.Error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsReader {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ThreadId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIdIsNotSet, Message: RpcErrorMsg_ThreadIdIsNotSet}
	}

	// Read the thread.
	var thread *mm.Thread
	var err error
	thread, err = srv.dbo.GetThreadById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.GetThreadResult{
		Thread: thread,
	}

	return result, nil
}

// deleteThread removes a thread.
func (srv *Server) deleteThread(p *mm.DeleteThreadParams) (result *mm.DeleteThreadResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ThreadId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIdIsNotSet, Message: RpcErrorMsg_ThreadIdIsNotSet}
	}

	// Read the thread.
	var thread *mm.Thread
	var err error
	thread, err = srv.dbo.GetThreadById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Check for children.
	if thread.Messages.Size() > 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIsNotEmpty, Message: RpcErrorMsg_ThreadIsNotEmpty}
	}

	// Update the link.
	var linkThreads *ul.UidList
	linkThreads, err = srv.dbo.GetForumThreadsById(thread.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = linkThreads.RemoveItem(p.ThreadId)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.dbo.SetForumThreadsById(thread.ForumId, linkThreads)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Delete the thread.
	err = srv.dbo.DeleteThreadById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.DeleteThreadResult{
		OK: true,
	}

	return result, nil
}

// Message.

// addMessage inserts a new message into a thread.
func (srv *Server) addMessage(p *mm.AddMessageParams) (result *mm.AddMessageResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsWriter {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ThreadId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIdIsNotSet, Message: RpcErrorMsg_ThreadIdIsNotSet}
	}

	if len(p.Text) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_MessageTextIsNotSet, Message: RpcErrorMsg_MessageTextIsNotSet}
	}

	// Ensure that a parent exists.
	var err error
	var n int
	n, err = srv.dbo.CountThreadsById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIsNotFound, Message: RpcErrorMsg_ThreadIsNotFound}
	}

	// Insert a message and link it with its thread.
	var parentMessages *ul.UidList
	parentMessages, err = srv.dbo.GetThreadMessagesById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	messageTextChecksum := srv.getMessageTextChecksum(p.Text)

	var insertedMessageId int64
	insertedMessageId, err = srv.dbo.InsertNewMessage(p.ThreadId, p.Text, messageTextChecksum, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = parentMessages.AddItem(uint(insertedMessageId), false)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.dbo.SetThreadMessagesById(p.ThreadId, parentMessages)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update thread's position if needed.
	if srv.settings.SystemSettings.NewThreadsAtTop {
		var messageThread *mm.Thread
		messageThread, err = srv.dbo.GetThreadById(p.ThreadId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		var threads *ul.UidList
		threads, err = srv.dbo.GetForumThreadsById(messageThread.ForumId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		if threads.Size() > 1 {
			err = threads.RaiseItem(p.ThreadId)
			if err != nil {
				srv.logError(err)
				return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
			}

			err = srv.dbo.SetForumThreadsById(messageThread.ForumId, threads)
			if err != nil {
				return nil, srv.databaseError(err)
			}
		}
	}

	result = &mm.AddMessageResult{
		MessageId: uint(insertedMessageId),
	}

	return result, nil
}

// changeMessageText changes text of a message.
func (srv *Server) changeMessageText(p *mm.ChangeMessageTextParams) (result *mm.ChangeMessageTextResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	var ok = false
	if userRoles.IsModerator {
		ok = true
	} else {
		// User can edit its own messages if they are not too old.
		if userRoles.IsWriter {
			var creatorUserId uint
			var ToC time.Time
			var ToE *time.Time
			var err error
			creatorUserId, ToC, ToE, err = srv.dbo.GetMessageCreatorAndTimeById(p.MessageId)
			if err != nil {
				return nil, srv.databaseError(err)
			}

			if userRoles.UserId == creatorUserId {
				var lastTouchTime time.Time
				if ToE != nil {
					lastTouchTime = *ToE
				} else {
					lastTouchTime = ToC
				}

				if time.Now().Before(lastTouchTime.Add(time.Second * time.Duration(srv.settings.SystemSettings.MessageEditTime))) {
					ok = true
				}
			}
		}
	}

	if !ok {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.MessageId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_MessageIdIsNotSet, Message: RpcErrorMsg_MessageIdIsNotSet}
	}

	if len(p.Text) == 0 {
		return nil, &js.Error{Code: RpcErrorCode_MessageTextIsNotSet, Message: RpcErrorMsg_MessageTextIsNotSet}
	}

	messageTextChecksum := srv.getMessageTextChecksum(p.Text)

	var err error
	err = srv.dbo.SetMessageTextById(p.MessageId, p.Text, messageTextChecksum, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ChangeMessageTextResult{
		OK: true,
	}

	return result, nil
}

// changeMessageThread moves a message from an old thread to a new thread.
func (srv *Server) changeMessageThread(p *mm.ChangeMessageThreadParams) (result *mm.ChangeMessageThreadResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.MessageId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_MessageIdIsNotSet, Message: RpcErrorMsg_MessageIdIsNotSet}
	}

	if p.ThreadId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIdIsNotSet, Message: RpcErrorMsg_ThreadIdIsNotSet}
	}

	// Ensure that an old parent exists.
	var oldParent uint
	var err error
	oldParent, err = srv.dbo.GetMessageThreadById(p.MessageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var n int
	n, err = srv.dbo.CountThreadsById(oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIsNotFound, Message: RpcErrorMsg_ThreadIsNotFound}
	}

	// Ensure that a new parent exists.
	n, err = srv.dbo.CountThreadsById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIsNotFound, Message: RpcErrorMsg_ThreadIsNotFound}
	}

	// Update the moved message.
	err = srv.dbo.SetMessageThreadById(p.MessageId, p.ThreadId, userRoles.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the new link.
	var messagesR *ul.UidList
	messagesR, err = srv.dbo.GetThreadMessagesById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = messagesR.AddItem(p.MessageId, false)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.dbo.SetThreadMessagesById(p.ThreadId, messagesR)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the old link.
	var messagesL *ul.UidList
	messagesL, err = srv.dbo.GetThreadMessagesById(oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = messagesL.RemoveItem(p.MessageId)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.dbo.SetThreadMessagesById(oldParent, messagesL)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ChangeMessageThreadResult{
		OK: true,
	}

	return result, nil
}

// getMessage reads a message.
func (srv *Server) getMessage(p *mm.GetMessageParams) (result *mm.GetMessageResult, jerr *js.Error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsReader {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.MessageId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_MessageIdIsNotSet, Message: RpcErrorMsg_MessageIdIsNotSet}
	}

	// Read the message.
	var message *mm.Message
	var err error
	message, err = srv.dbo.GetMessageById(p.MessageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.GetMessageResult{
		Message: message,
	}

	return result, nil
}

// deleteMessage removes a message.
func (srv *Server) deleteMessage(p *mm.DeleteMessageParams) (result *mm.DeleteMessageResult, jerr *js.Error) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsAdministrator {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.MessageId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_MessageIdIsNotSet, Message: RpcErrorMsg_MessageIdIsNotSet}
	}

	// Read the message.
	var message *mm.Message
	var err error
	message, err = srv.dbo.GetMessageById(p.MessageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the link.
	var linkMessages *ul.UidList
	linkMessages, err = srv.dbo.GetThreadMessagesById(message.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = linkMessages.RemoveItem(p.MessageId)
	if err != nil {
		srv.logError(err)
		return nil, &js.Error{Code: c.RpcErrorCode_UidList, Message: fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error())}
	}

	err = srv.dbo.SetThreadMessagesById(message.ThreadId, linkMessages)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Delete the message.
	err = srv.dbo.DeleteMessageById(p.MessageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.DeleteMessageResult{
		OK: true,
	}

	return result, nil
}

// Composite objects.

// listThreadAndMessages reads a thread and all its messages.
func (srv *Server) listThreadAndMessages(p *mm.ListThreadAndMessagesParams) (result *mm.ListThreadAndMessagesResult, jerr *js.Error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsReader {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ThreadId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIdIsNotSet, Message: RpcErrorMsg_ThreadIdIsNotSet}
	}

	// Read the thread.
	var thread *mm.Thread
	var err error
	thread, err = srv.dbo.GetThreadById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	tam := mm.NewThreadAndMessages(thread)
	tam.MessageIds = thread.Messages

	// Read all the messages.
	if tam.MessageIds != nil {
		tam.Messages, err = srv.dbo.ReadMessagesById(*tam.MessageIds)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	result = &mm.ListThreadAndMessagesResult{
		ThreadAndMessages: tam,
	}

	return result, nil
}

// listThreadAndMessagesOnPage reads a thread and its messages on a selected page.
func (srv *Server) listThreadAndMessagesOnPage(p *mm.ListThreadAndMessagesOnPageParams) (result *mm.ListThreadAndMessagesOnPageResult, jerr *js.Error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsReader {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ThreadId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ThreadIdIsNotSet, Message: RpcErrorMsg_ThreadIdIsNotSet}
	}

	if p.Page == 0 {
		return nil, &js.Error{Code: RpcErrorCode_PageIsNotSet, Message: RpcErrorMsg_PageIsNotSet}
	}

	// Read the thread.
	var thread *mm.Thread
	var err error
	thread, err = srv.dbo.GetThreadById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	tamop := mm.NewThreadAndMessages(thread)
	tamop.MessageIds = thread.Messages

	if tamop.MessageIds != nil {
		// Total numbers before pagination.
		tm := uint(tamop.MessageIds.Size())
		tamop.TotalMessages = &tm
		tp := uint(math.Ceil(float64(tm) / float64(srv.settings.SystemSettings.PageSize)))
		tamop.TotalPages = &tp

		// Read messages of a specified page.
		tamop.Page = &p.Page
		tamop.MessageIds = tamop.MessageIds.OnPage(p.Page, srv.settings.SystemSettings.PageSize)
		if tamop.MessageIds.Size() > 0 {
			tamop.Messages, err = srv.dbo.ReadMessagesById(*tamop.MessageIds)
			if err != nil {
				return nil, srv.databaseError(err)
			}
		}
	}

	result = &mm.ListThreadAndMessagesOnPageResult{
		ThreadAndMessagesOnPage: tamop,
	}

	return result, nil
}

// listForumAndThreads reads a forum and all its threads.
func (srv *Server) listForumAndThreads(p *mm.ListForumAndThreadsParams) (result *mm.ListForumAndThreadsResult, jerr *js.Error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsReader {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ForumId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIdIsNotSet, Message: RpcErrorMsg_ForumIdIsNotSet}
	}

	// Read the forum.
	var forum *mm.Forum
	var err error
	forum, err = srv.dbo.GetForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	fat := mm.NewForumAndThreads(forum)
	fat.ThreadIds = forum.Threads

	// Read all the threads.
	if fat.ThreadIds != nil {
		fat.Threads, err = srv.dbo.ReadThreadsById(*forum.Threads)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	result = &mm.ListForumAndThreadsResult{
		ForumAndThreads: fat,
	}

	return result, nil
}

// listForumAndThreadsOnPage reads a forum and its threads on a selected page.
func (srv *Server) listForumAndThreadsOnPage(p *mm.ListForumAndThreadsOnPageParams) (result *mm.ListForumAndThreadsOnPageResult, jerr *js.Error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsReader {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Check parameters.
	if p.ForumId == 0 {
		return nil, &js.Error{Code: RpcErrorCode_ForumIdIsNotSet, Message: RpcErrorMsg_ForumIdIsNotSet}
	}

	if p.Page == 0 {
		return nil, &js.Error{Code: RpcErrorCode_PageIsNotSet, Message: RpcErrorMsg_PageIsNotSet}
	}

	// Read the forum.
	var forum *mm.Forum
	var err error
	forum, err = srv.dbo.GetForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	fatop := mm.NewForumAndThreads(forum)
	fatop.ThreadIds = forum.Threads

	if fatop.ThreadIds != nil {
		// Total numbers before pagination.
		tt := uint(fatop.ThreadIds.Size())
		fatop.TotalThreads = &tt
		tp := uint(math.Ceil(float64(tt) / float64(srv.settings.SystemSettings.PageSize)))
		fatop.TotalPages = &tp

		// Read threads of a specified page.
		fatop.Page = &p.Page
		fatop.ThreadIds = fatop.ThreadIds.OnPage(p.Page, srv.settings.SystemSettings.PageSize)
		if fatop.ThreadIds.Size() > 0 {
			fatop.Threads, err = srv.dbo.ReadThreadsById(*fatop.ThreadIds)
			if err != nil {
				return nil, srv.databaseError(err)
			}
		}
	}

	result = &mm.ListForumAndThreadsOnPageResult{
		ForumAndThreadsOnPage: fatop,
	}

	return result, nil
}

// listSectionsAndForums reads all sections and forums.
func (srv *Server) listSectionsAndForums(p *mm.ListSectionsAndForumsParams) (result *mm.ListSectionsAndForumsResult, jerr *js.Error) {
	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	var userRoles *am.GetSelfRolesResult
	userRoles, jerr = srv.mustBeAnAuthToken(p.Auth)
	if jerr != nil {
		return nil, jerr
	}

	// Check permissions.
	if !userRoles.IsReader {
		return nil, &js.Error{Code: c.RpcErrorCode_InsufficientPermission, Message: c.RpcErrorMsg_InsufficientPermission}
	}

	// Read all the sections.
	var sections []mm.Section
	var err error
	sections, err = srv.dbo.ReadSections()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Read all the forums.
	var forums []mm.Forum
	forums, err = srv.dbo.ReadForums()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ListSectionsAndForumsResult{
		SectionsAndForums: &mm.SectionsAndForums{
			Sections: sections,
			Forums:   forums,
		},
	}

	return result, nil
}

// Other.

func (srv *Server) showDiagnosticData() (result *mm.ShowDiagnosticDataResult, jerr *js.Error) {
	result = &mm.ShowDiagnosticDataResult{
		TotalRequestsCount:      srv.diag.GetTotalRequestsCount(),
		SuccessfulRequestsCount: srv.diag.GetSuccessfulRequestsCount(),
	}

	return result, nil
}

func (srv *Server) doTest(p *mm.TestParams) (result *mm.TestResult, jerr *js.Error) {
	result = &mm.TestResult{}

	var wg = new(sync.WaitGroup)
	var errChan = make(chan error, p.N)

	for i := uint(1); i <= p.N; i++ {
		wg.Add(1)
		go srv.doTestA(wg, errChan)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			srv.logError(err)
			return nil, &js.Error{Code: RpcErrorCode_TestError, Message: fmt.Sprintf(RpcErrorMsgF_TestError, err.Error())}
		}
	}

	return result, nil
}
