package mm

import (
	"fmt"
	"sync"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/models"
	c "github.com/vault-thirteen/SimpleBB/pkg/common"
	"github.com/vault-thirteen/SimpleBB/pkg/common/UidList"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

// RPC functions.

// Section.

// addSection inserts a new section as a root section or as a sub-section.
func (srv *Server) addSection(p *mm.AddSectionParams) (result *mm.AddSectionResult, re *jrm1.RpcError) {
	// Check parameters.
	if len(p.Name) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionNameIsNotSet, RpcErrorMsg_SectionNameIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// If parent is not set, the new section is a root section.
	// Only a single root section may exist.
	var err error
	var n cmb.Count
	if p.Parent == nil {
		n, err = srv.dbo.CountRootSections()
		if err != nil {
			return nil, srv.databaseError(err)
		}

		if n > 0 {
			return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RootSectionAlreadyExists, RpcErrorMsg_RootSectionAlreadyExists, nil)
		}

		var insertedSectionId cmb.Id
		insertedSectionId, err = srv.dbo.InsertNewSection(p.Parent, p.Name, userRoles.User.Id)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		result = &mm.AddSectionResult{
			SectionId: insertedSectionId,
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
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Check compatibility.
	var childType mm.SectionChildType
	childType, err = srv.dbo.GetSectionChildTypeById(*p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if childType.GetValue() == mm.SectionChildType_Forum {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IncompatibleChildType, RpcErrorMsg_IncompatibleChildType, nil)
	}

	if childType.GetValue() == 0 {
		err = srv.dbo.SetSectionChildTypeById(*p.Parent, mm.SectionChildType_Section)
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

	var insertedSectionId cmb.Id
	insertedSectionId, err = srv.dbo.InsertNewSection(p.Parent, p.Name, userRoles.User.Id)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = parentChildren.AddItem(insertedSectionId, false)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(*p.Parent, parentChildren)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.AddSectionResult{
		SectionId: insertedSectionId,
	}

	return result, nil
}

// changeSectionName renames a section.
func (srv *Server) changeSectionName(p *mm.ChangeSectionNameParams) (result *mm.ChangeSectionNameResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.SectionId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	if len(p.Name) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionNameIsNotSet, RpcErrorMsg_SectionNameIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var n cmb.Count
	var err error
	n, err = srv.dbo.CountSectionsById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	err = srv.dbo.SetSectionNameById(p.SectionId, p.Name, userRoles.User.Id)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ChangeSectionNameResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// changeSectionParent moves a section from an old parent to a new parent.
func (srv *Server) changeSectionParent(p *mm.ChangeSectionParentParams) (result *mm.ChangeSectionParentResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.SectionId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	if p.Parent == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var n cmb.Count
	var err error
	n, err = srv.dbo.CountSectionsById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Ensure that an old parent exists.
	var oldParent *cmb.Id
	oldParent, err = srv.dbo.GetSectionParentById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if oldParent == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RootSectionCanNotBeMoved, RpcErrorMsg_RootSectionCanNotBeMoved, nil)
	}

	n, err = srv.dbo.CountSectionsById(*oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Ensure that a new parent exists.
	n, err = srv.dbo.CountSectionsById(p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Check compatibility.
	var childType mm.SectionChildType
	childType, err = srv.dbo.GetSectionChildTypeById(p.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if childType.GetValue() == mm.SectionChildType_Forum {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IncompatibleChildType, RpcErrorMsg_IncompatibleChildType, nil)
	}

	if childType.GetValue() == 0 {
		err = srv.dbo.SetSectionChildTypeById(p.Parent, mm.SectionChildType_Section)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	// Update the moved section.
	err = srv.dbo.SetSectionParentById(p.SectionId, p.Parent, userRoles.User.Id)
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
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
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
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(*oldParent, childrenL)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Clear the child type if the old parent becomes empty.
	if childrenL.Size() == 0 {
		err = srv.dbo.SetSectionChildTypeById(*oldParent, 0)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	result = &mm.ChangeSectionParentResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// getSection reads a section.
func (srv *Server) getSection(p *mm.GetSectionParams) (result *mm.GetSectionResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.SectionId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the section.
	var section *mm.Section
	var err error
	section, err = srv.dbo.GetSectionById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if section == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	result = &mm.GetSectionResult{
		Section: section,
	}

	return result, nil
}

// moveSectionUp moves a section up by one position if possible.
func (srv *Server) moveSectionUp(p *mm.MoveSectionUpParams) (result *mm.MoveSectionUpResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.SectionId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Check existence of the moved section.
	var n cmb.Count
	var err error
	n, err = srv.dbo.CountSectionsById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Get the section which is being moved.
	var section *mm.Section
	section, err = srv.dbo.GetSectionById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if section == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}
	if section.Parent == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RootSectionCanNotBeMoved, RpcErrorMsg_RootSectionCanNotBeMoved, nil)
	}

	// Get the parent section.
	var parent *mm.Section
	parent, err = srv.dbo.GetSectionById(*section.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if parent == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}
	if parent.Children == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Check compatibility.
	if parent.ChildType.GetValue() != mm.SectionChildType_Section {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IncompatibleChildType, RpcErrorMsg_IncompatibleChildType, nil)
	}

	err = parent.Children.MoveItemUp(p.SectionId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(parent.Id, parent.Children)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.MoveSectionUpResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// moveSectionDown moves a section down by one position if possible.
func (srv *Server) moveSectionDown(p *mm.MoveSectionDownParams) (result *mm.MoveSectionDownResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.SectionId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Check existence of the moved section.
	var n cmb.Count
	var err error
	n, err = srv.dbo.CountSectionsById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Get the section which is being moved.
	var section *mm.Section
	section, err = srv.dbo.GetSectionById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if section == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}
	if section.Parent == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_RootSectionCanNotBeMoved, RpcErrorMsg_RootSectionCanNotBeMoved, nil)
	}

	// Get the parent section.
	var parent *mm.Section
	parent, err = srv.dbo.GetSectionById(*section.Parent)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if parent == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}
	if parent.Children == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Check compatibility.
	if parent.ChildType.GetValue() != mm.SectionChildType_Section {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IncompatibleChildType, RpcErrorMsg_IncompatibleChildType, nil)
	}

	err = parent.Children.MoveItemDown(p.SectionId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(parent.Id, parent.Children)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.MoveSectionDownResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// deleteSection removes a section.
func (srv *Server) deleteSection(p *mm.DeleteSectionParams) (result *mm.DeleteSectionResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.SectionId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Read the section.
	var section *mm.Section
	var err error
	section, err = srv.dbo.GetSectionById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if section == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	var isRootSection = false
	if section.Parent == nil {
		isRootSection = true
	}

	// Check for children.
	if section.Children.Size() > 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionHasChildren, RpcErrorMsg_SectionHasChildren, nil)
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
			return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
		}

		err = srv.dbo.SetSectionChildrenById(*section.Parent, linkSections)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		// Clear the child type if the old parent becomes empty.
		if linkSections.Size() == 0 {
			err = srv.dbo.SetSectionChildTypeById(*section.Parent, 0)
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
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// Forum.

// addForum inserts a new forum into a section.
func (srv *Server) addForum(p *mm.AddForumParams) (result *mm.AddForumResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.SectionId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	if len(p.Name) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumNameIsNotSet, RpcErrorMsg_ForumNameIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Ensure that a section exists.
	n, err := srv.dbo.CountSectionsById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Check compatibility.
	var childType mm.SectionChildType
	childType, err = srv.dbo.GetSectionChildTypeById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if childType.GetValue() == mm.SectionChildType_Section {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IncompatibleChildType, RpcErrorMsg_IncompatibleChildType, nil)
	}

	if childType.GetValue() == 0 {
		err = srv.dbo.SetSectionChildTypeById(p.SectionId, mm.SectionChildType_Forum)
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

	var insertedForumId cmb.Id
	insertedForumId, err = srv.dbo.InsertNewForum(p.SectionId, p.Name, userRoles.User.Id)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = parentChildren.AddItem(insertedForumId, false)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(p.SectionId, parentChildren)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.AddForumResult{
		ForumId: insertedForumId,
	}

	return result, nil
}

// changeForumName renames a forum.
func (srv *Server) changeForumName(p *mm.ChangeForumNameParams) (result *mm.ChangeForumNameResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	if len(p.Name) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumNameIsNotSet, RpcErrorMsg_ForumNameIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var n cmb.Count
	var err error
	n, err = srv.dbo.CountForumsById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	err = srv.dbo.SetForumNameById(p.ForumId, p.Name, userRoles.User.Id)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ChangeForumNameResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// changeForumSection moves a forum from an old section to a new section.
func (srv *Server) changeForumSection(p *mm.ChangeForumSectionParams) (result *mm.ChangeForumSectionResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	if p.SectionId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIdIsNotSet, RpcErrorMsg_SectionIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var n cmb.Count
	var err error
	n, err = srv.dbo.CountForumsById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Ensure that an old section exists.
	var oldParent cmb.Id
	oldParent, err = srv.dbo.GetForumSectionById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	n, err = srv.dbo.CountSectionsById(oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Ensure that a new section exists.
	n, err = srv.dbo.CountSectionsById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}

	// Check compatibility.
	var childType mm.SectionChildType
	childType, err = srv.dbo.GetSectionChildTypeById(p.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if childType.GetValue() == mm.SectionChildType_Section {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IncompatibleChildType, RpcErrorMsg_IncompatibleChildType, nil)
	}

	if childType.GetValue() == 0 {
		err = srv.dbo.SetSectionChildTypeById(p.SectionId, mm.SectionChildType_Forum)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	// Update the moved forum.
	err = srv.dbo.SetForumSectionById(p.ForumId, p.SectionId, userRoles.User.Id)
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
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
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
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(oldParent, childrenL)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Clear the child type if the old section becomes empty.
	if childrenL.Size() == 0 {
		err = srv.dbo.SetSectionChildTypeById(oldParent, 0)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	result = &mm.ChangeForumSectionResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// getForum reads a forum.
func (srv *Server) getForum(p *mm.GetForumParams) (result *mm.GetForumResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the forum.
	var forum *mm.Forum
	var err error
	forum, err = srv.dbo.GetForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if forum == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	result = &mm.GetForumResult{
		Forum: forum,
	}

	return result, nil
}

// moveForumUp moves a forum up by one position if possible.
func (srv *Server) moveForumUp(p *mm.MoveForumUpParams) (result *mm.MoveForumUpResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Check existence of the moved forum.
	var n cmb.Count
	var err error
	n, err = srv.dbo.CountForumsById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Get the forum which is being moved.
	var forum *mm.Forum
	forum, err = srv.dbo.GetForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if forum == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Get the parent section.
	var parent *mm.Section
	parent, err = srv.dbo.GetSectionById(forum.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if parent == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}
	if parent.Children == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Check compatibility.
	if parent.ChildType.GetValue() != mm.SectionChildType_Forum {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IncompatibleChildType, RpcErrorMsg_IncompatibleChildType, nil)
	}

	err = parent.Children.MoveItemUp(p.ForumId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(parent.Id, parent.Children)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.MoveForumUpResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// moveForumDown moves a forum down by one position if possible.
func (srv *Server) moveForumDown(p *mm.MoveForumDownParams) (result *mm.MoveForumDownResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Check existence of the moved forum.
	var n cmb.Count
	var err error
	n, err = srv.dbo.CountForumsById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Get the forum which is being moved.
	var forum *mm.Forum
	forum, err = srv.dbo.GetForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if forum == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Get the parent section.
	var parent *mm.Section
	parent, err = srv.dbo.GetSectionById(forum.SectionId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if parent == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SectionIsNotFound, RpcErrorMsg_SectionIsNotFound, nil)
	}
	if parent.Children == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Check compatibility.
	if parent.ChildType.GetValue() != mm.SectionChildType_Forum {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_IncompatibleChildType, RpcErrorMsg_IncompatibleChildType, nil)
	}

	err = parent.Children.MoveItemDown(p.ForumId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(parent.Id, parent.Children)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.MoveForumDownResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// deleteForum removes a forum.
func (srv *Server) deleteForum(p *mm.DeleteForumParams) (result *mm.DeleteForumResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Read the forum.
	var forum *mm.Forum
	var err error
	forum, err = srv.dbo.GetForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if forum == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Check for threads.
	if forum.Threads.Size() > 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumHasThreads, RpcErrorMsg_ForumHasThreads, nil)
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
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetSectionChildrenById(forum.SectionId, linkChildren)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Clear the child type if the old parent becomes empty.
	if linkChildren.Size() == 0 {
		err = srv.dbo.SetSectionChildTypeById(forum.SectionId, 0)
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
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// Thread.

// addThread inserts a new thread into a forum.
func (srv *Server) addThread(p *mm.AddThreadParams) (result *mm.AddThreadResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	if len(p.Name) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadNameIsNotSet, RpcErrorMsg_ThreadNameIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAuthor {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Ensure that a forum exists.
	var err error
	var n cmb.Count
	n, err = srv.dbo.CountForumsById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Insert a thread and link it with its forum.
	var parentThreads *ul.UidList
	parentThreads, err = srv.dbo.GetForumThreadsById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var insertedThreadId cmb.Id
	insertedThreadId, err = srv.dbo.InsertNewThread(p.ForumId, p.Name, userRoles.User.Id)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = parentThreads.AddItem(insertedThreadId, srv.settings.SystemSettings.NewThreadsAtTop.AsBool())
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetForumThreadsById(p.ForumId, parentThreads)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.AddThreadResult{
		ThreadId: insertedThreadId,
	}

	return result, nil
}

// changeThreadName renames a thread.
func (srv *Server) changeThreadName(p *mm.ChangeThreadNameParams) (result *mm.ChangeThreadNameResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	if len(p.Name) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadNameIsNotSet, RpcErrorMsg_ThreadNameIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	se, err := cm.NewSystemEventWithData(
		cm.SystemEventData{
			Type:     cm.NewSystemEventTypeWithValue(cm.SystemEventType_ThreadNameChange),
			ThreadId: &p.ThreadId,
			UserId:   &userRoles.User.Id,
		},
	)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	re = srv.changeThreadNameH(p.ThreadId, p.Name, userRoles.User.Id)
	if re != nil {
		return nil, re
	}

	result = &mm.ChangeThreadNameResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// changeThreadForum moves a thread from an old forum to a new forum.
func (srv *Server) changeThreadForum(p *mm.ChangeThreadForumParams) (result *mm.ChangeThreadForumResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	se, err := cm.NewSystemEventWithData(
		cm.SystemEventData{
			Type:     cm.NewSystemEventTypeWithValue(cm.SystemEventType_ThreadParentChange),
			ThreadId: &p.ThreadId,
			UserId:   &userRoles.User.Id,
		},
	)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	re = srv.changeThreadForumH(p.ThreadId, p.ForumId, userRoles.User.Id)
	if re != nil {
		return nil, re
	}

	result = &mm.ChangeThreadForumResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// getThread reads a thread.
func (srv *Server) getThread(p *mm.GetThreadParams) (result *mm.GetThreadResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the thread.
	var thread *mm.Thread
	var err error
	thread, err = srv.dbo.GetThreadById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if thread == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	result = &mm.GetThreadResult{
		Thread: thread,
	}

	return result, nil
}

// getThreadNamesByIds reads names of threads specified by their IDs.
func (srv *Server) getThreadNamesByIds(p *mm.GetThreadNamesByIdsParams) (result *mm.GetThreadNamesByIdsResult, re *jrm1.RpcError) {
	// Check parameters.
	if len(p.ThreadIds) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read thread names.
	var threadNames []cm.Name
	var err error
	threadNames, err = srv.dbo.ReadThreadNamesByIds(p.ThreadIds)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.GetThreadNamesByIdsResult{
		ThreadIds:   p.ThreadIds,
		ThreadNames: threadNames,
	}

	return result, nil
}

// moveThreadUp moves a thread up by one position if possible.
func (srv *Server) moveThreadUp(p *mm.MoveThreadUpParams) (result *mm.MoveThreadUpResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Check existence of the moved thread.
	var n cmb.Count
	var err error
	n, err = srv.dbo.CountThreadsById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Get the thread which is being moved.
	var thread *mm.Thread
	thread, err = srv.dbo.GetThreadById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if thread == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Get the parent forum.
	var parent *mm.Forum
	parent, err = srv.dbo.GetForumById(thread.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if parent == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}
	if parent.Threads == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	err = parent.Threads.MoveItemUp(p.ThreadId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetForumThreadsById(parent.Id, parent.Threads)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.MoveThreadUpResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// moveThreadDown moves a thread down by one position if possible.
func (srv *Server) moveThreadDown(p *mm.MoveThreadDownParams) (result *mm.MoveThreadDownResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Check existence of the moved thread.
	var n cmb.Count
	var err error
	n, err = srv.dbo.CountThreadsById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Get the thread which is being moved.
	var thread *mm.Thread
	thread, err = srv.dbo.GetThreadById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if thread == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Get the parent forum.
	var parent *mm.Forum
	parent, err = srv.dbo.GetForumById(thread.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if parent == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}
	if parent.Threads == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	err = parent.Threads.MoveItemDown(p.ThreadId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetForumThreadsById(parent.Id, parent.Threads)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.MoveThreadDownResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// deleteThread removes a thread.
func (srv *Server) deleteThread(p *mm.DeleteThreadParams) (result *mm.DeleteThreadResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	se, err := cm.NewSystemEventWithData(
		cm.SystemEventData{
			Type:     cm.NewSystemEventTypeWithValue(cm.SystemEventType_ThreadDeletion),
			ThreadId: &p.ThreadId,
			UserId:   &userRoles.User.Id,
		},
	)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	re = srv.deleteThreadH(p.ThreadId)
	if re != nil {
		return nil, re
	}

	result = &mm.DeleteThreadResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// threadExistsS checks whether the specified thread exists or not. This method
// is used by the system.
func (srv *Server) threadExistsS(p *mm.ThreadExistsSParams) (result *mm.ThreadExistsSResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check the DKey.
	if !srv.dKeyI.CheckString(p.DKey.ToString()) {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Count threads.
	var n cmb.Count
	var err error
	n, err = srv.dbo.CountThreadsById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ThreadExistsSResult{
		Exists: n == 1,
	}

	return result, nil
}

// Message.

// addMessage inserts a new message into a thread.
func (srv *Server) addMessage(p *mm.AddMessageParams) (result *mm.AddMessageResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	if len(p.Text) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageTextIsNotSet, RpcErrorMsg_MessageTextIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions (Part I).
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	result, re = srv.addMessageH(p.ThreadId, p.Text, userRoles)
	if re != nil {
		return nil, re
	}

	se, err := cm.NewSystemEventWithData(
		cm.SystemEventData{
			Type:      cm.NewSystemEventTypeWithValue(cm.SystemEventType_ThreadNewMessage),
			ThreadId:  &p.ThreadId,
			MessageId: &result.MessageId,
			UserId:    &userRoles.User.Id,
		},
	)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	return result, nil
}

// changeMessageText changes text of a message.
func (srv *Server) changeMessageText(p *mm.ChangeMessageTextParams) (result *mm.ChangeMessageTextResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.MessageId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageIdIsNotSet, RpcErrorMsg_MessageIdIsNotSet, nil)
	}

	if len(p.Text) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageTextIsNotSet, RpcErrorMsg_MessageTextIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	var initialMessage *mm.Message
	initialMessage, re = srv.changeMessageTextH(p.MessageId, p.Text, userRoles)
	if re != nil {
		return nil, re
	}

	se, err := cm.NewSystemEventWithData(
		cm.SystemEventData{
			Type:      cm.NewSystemEventTypeWithValue(cm.SystemEventType_ThreadMessageEdit),
			ThreadId:  &initialMessage.ThreadId,
			MessageId: &p.MessageId,
			UserId:    &userRoles.User.Id,
		},
	)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	se, err = cm.NewSystemEventWithData(
		cm.SystemEventData{
			Type:      cm.NewSystemEventTypeWithValue(cm.SystemEventType_MessageTextEdit),
			ThreadId:  &initialMessage.ThreadId,
			MessageId: &p.MessageId,
			UserId:    &userRoles.User.Id,
			Creator:   &initialMessage.Creator.UserId,
		},
	)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	result = &mm.ChangeMessageTextResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// changeMessageThread moves a message from an old thread to a new thread.
func (srv *Server) changeMessageThread(p *mm.ChangeMessageThreadParams) (result *mm.ChangeMessageThreadResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.MessageId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageIdIsNotSet, RpcErrorMsg_MessageIdIsNotSet, nil)
	}

	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var initialMessage *mm.Message
	initialMessage, re = srv.changeMessageThreadH(p.MessageId, p.ThreadId, userRoles)
	if re != nil {
		return nil, re
	}

	se, err := cm.NewSystemEventWithData(
		cm.SystemEventData{
			Type:      cm.NewSystemEventTypeWithValue(cm.SystemEventType_MessageParentChange),
			ThreadId:  &initialMessage.ThreadId,
			MessageId: &p.MessageId,
			UserId:    &userRoles.User.Id,
			Creator:   &initialMessage.Creator.UserId,
		},
	)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	result = &mm.ChangeMessageThreadResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// getMessage reads a message.
func (srv *Server) getMessage(p *mm.GetMessageParams) (result *mm.GetMessageResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.MessageId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageIdIsNotSet, RpcErrorMsg_MessageIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the message.
	var message *mm.Message
	var err error
	message, err = srv.dbo.GetMessageById(p.MessageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if message == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageIsNotFound, RpcErrorMsg_MessageIsNotFound, nil)
	}

	result = &mm.GetMessageResult{
		Message: message,
	}

	return result, nil
}

// getLatestMessageOfThread reads the latest message of a thread.
func (srv *Server) getLatestMessageOfThread(p *mm.GetLatestMessageOfThreadParams) (result *mm.GetLatestMessageOfThreadResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	result = &mm.GetLatestMessageOfThreadResult{}

	result.Message, re = srv.getLatestMessageOfThreadH(p.ThreadId)
	if re != nil {
		return nil, re
	}

	return result, nil
}

// deleteMessage removes a message.
func (srv *Server) deleteMessage(p *mm.DeleteMessageParams) (result *mm.DeleteMessageResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.MessageId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageIdIsNotSet, RpcErrorMsg_MessageIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var initialMessage *mm.Message
	initialMessage, re = srv.deleteMessageH(p.MessageId)
	if re != nil {
		return nil, re
	}

	se, err := cm.NewSystemEventWithData(
		cm.SystemEventData{
			Type:      cm.NewSystemEventTypeWithValue(cm.SystemEventType_ThreadMessageDeletion),
			ThreadId:  &initialMessage.ThreadId,
			MessageId: &p.MessageId,
			UserId:    &userRoles.User.Id,
		},
	)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	se, err = cm.NewSystemEventWithData(
		cm.SystemEventData{
			Type:      cm.NewSystemEventTypeWithValue(cm.SystemEventType_MessageDeletion),
			ThreadId:  &initialMessage.ThreadId,
			MessageId: &p.MessageId,
			UserId:    &userRoles.User.Id,
			Creator:   &initialMessage.Creator.UserId,
		},
	)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_SystemEvent, c.RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.reportSystemEvent(se)
	if re != nil {
		return nil, re
	}

	result = &mm.DeleteMessageResult{
		Success: cmr.Success{
			OK: true,
		},
	}
	return result, nil
}

// Composite objects.

// listThreadAndMessages reads a thread and all its messages.
func (srv *Server) listThreadAndMessages(p *mm.ListThreadAndMessagesParams) (result *mm.ListThreadAndMessagesResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the thread.
	var thread *mm.Thread
	var err error
	thread, err = srv.dbo.GetThreadById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if thread == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Read messages.
	var allMessageIds = thread.Messages

	var allMessages []mm.Message
	allMessages, err = srv.dbo.ReadMessagesById(allMessageIds)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ListThreadAndMessagesResult{
		ThreadAndMessages: &mm.ThreadAndMessages{
			Thread:   thread,
			Messages: allMessages,
		},
	}

	return result, nil
}

// listThreadAndMessagesOnPage reads a thread and its messages on a selected page.
func (srv *Server) listThreadAndMessagesOnPage(p *mm.ListThreadAndMessagesOnPageParams) (result *mm.ListThreadAndMessagesOnPageResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	if p.Page == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PageIsNotSet, RpcErrorMsg_PageIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the thread.
	var thread *mm.Thread
	var err error
	thread, err = srv.dbo.GetThreadById(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if thread == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Read messages.
	var allMessageIds = thread.Messages
	var messageIdsOnPage = allMessageIds.OnPage(p.Page, srv.settings.SystemSettings.PageSize)

	var messagesOnPage []mm.Message
	messagesOnPage, err = srv.dbo.ReadMessagesById(messageIdsOnPage)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ListThreadAndMessagesOnPageResult{
		ThreadAndMessagesOnPage: &mm.ThreadAndMessages{
			Thread:   thread,
			Messages: messagesOnPage,
			PageData: &cmr.PageData{
				PageNumber:  p.Page,
				TotalPages:  cmb.CalculateTotalPages(allMessageIds.Size(), srv.settings.SystemSettings.PageSize),
				PageSize:    srv.settings.SystemSettings.PageSize,
				ItemsOnPage: messageIdsOnPage.Size(),
				TotalItems:  allMessageIds.Size(),
			},
		},
	}
	thread.Messages = messageIdsOnPage

	return result, nil
}

// listForumAndThreads reads a forum and all its threads.
func (srv *Server) listForumAndThreads(p *mm.ListForumAndThreadsParams) (result *mm.ListForumAndThreadsResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the forum.
	var forum *mm.Forum
	var err error
	forum, err = srv.dbo.GetForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if forum == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Read threads.
	var allThreadIds = forum.Threads

	var allThreads []mm.Thread
	allThreads, err = srv.dbo.ReadThreadsById(allThreadIds)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ListForumAndThreadsResult{
		ForumAndThreads: &mm.ForumAndThreads{
			Forum:   forum,
			Threads: allThreads,
		},
	}

	return result, nil
}

// listForumAndThreadsOnPage reads a forum and its threads on a selected page.
func (srv *Server) listForumAndThreadsOnPage(p *mm.ListForumAndThreadsOnPageParams) (result *mm.ListForumAndThreadsOnPageResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ForumId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIdIsNotSet, RpcErrorMsg_ForumIdIsNotSet, nil)
	}

	if p.Page == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PageIsNotSet, RpcErrorMsg_PageIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the forum.
	var forum *mm.Forum
	var err error
	forum, err = srv.dbo.GetForumById(p.ForumId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if forum == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Read threads.
	var allThreadIds = forum.Threads
	var threadIdsOnPage = allThreadIds.OnPage(p.Page, srv.settings.SystemSettings.PageSize)

	var threadsOnPage []mm.Thread
	threadsOnPage, err = srv.dbo.ReadThreadsById(threadIdsOnPage)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &mm.ListForumAndThreadsOnPageResult{
		ForumAndThreadsOnPage: &mm.ForumAndThreads{
			Forum:   forum,
			Threads: threadsOnPage,
			PageData: &cmr.PageData{
				PageNumber:  p.Page,
				TotalPages:  cmb.CalculateTotalPages(allThreadIds.Size(), srv.settings.SystemSettings.PageSize),
				PageSize:    srv.settings.SystemSettings.PageSize,
				ItemsOnPage: threadIdsOnPage.Size(),
				TotalItems:  allThreadIds.Size(),
			},
		},
	}
	forum.Threads = threadIdsOnPage

	return result, nil
}

// listSectionsAndForums reads all sections and forums.
func (srv *Server) listSectionsAndForums(p *mm.ListSectionsAndForumsParams) (result *mm.ListSectionsAndForumsResult, re *jrm1.RpcError) {
	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.Roles.IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

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

func (srv *Server) getDKey(p *mm.GetDKeyParams) (result *mm.GetDKeyResult, re *jrm1.RpcError) {
	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	result = &mm.GetDKeyResult{
		DKey: cmb.Text(srv.dKeyI.GetString()),
	}

	return result, nil
}

func (srv *Server) showDiagnosticData() (result *mm.ShowDiagnosticDataResult, re *jrm1.RpcError) {
	trc, src := srv.js.GetRequestsCount()

	result = &mm.ShowDiagnosticDataResult{
		RequestsCount: cmr.RequestsCount{
			TotalRequestsCount:      cmb.Text(trc),
			SuccessfulRequestsCount: cmb.Text(src),
		},
	}

	return result, nil
}

func (srv *Server) test(p *mm.TestParams) (result *mm.TestResult, re *jrm1.RpcError) {
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
			return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_TestError, fmt.Sprintf(RpcErrorMsgF_TestError, err.Error()), nil)
		}
	}

	return result, nil
}
