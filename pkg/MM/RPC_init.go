package mm

import (
	mc "github.com/vault-thirteen/SimpleBB/pkg/MM/client"
	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/models"
)

func (srv *Server) initJsonRpcHandlers() (err error) {
	// Ping.
	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncPing, PingHandler{Server: srv}, mm.PingParams{}, mm.PingResult{})
	if err != nil {
		return err
	}

	// Section.
	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncAddSection, AddSectionHandler{Server: srv}, mm.AddSectionParams{}, mm.AddSectionResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncChangeSectionName, ChangeSectionNameHandler{Server: srv}, mm.ChangeSectionNameParams{}, mm.ChangeSectionNameResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncChangeSectionParent, ChangeSectionParentHandler{Server: srv}, mm.ChangeSectionParentParams{}, mm.ChangeSectionParentResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncGetSection, GetSectionHandler{Server: srv}, mm.GetSectionParams{}, mm.GetSectionResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncDeleteSection, DeleteSectionHandler{Server: srv}, mm.DeleteSectionParams{}, mm.DeleteSectionResult{})
	if err != nil {
		return err
	}

	// Forum.
	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncAddForum, AddForumHandler{Server: srv}, mm.AddForumParams{}, mm.AddForumResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncChangeForumName, ChangeForumNameHandler{Server: srv}, mm.ChangeForumNameParams{}, mm.ChangeForumNameResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncChangeForumSection, ChangeForumSectionHandler{Server: srv}, mm.ChangeForumSectionParams{}, mm.ChangeForumSectionResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncGetForum, GetForumHandler{Server: srv}, mm.GetForumParams{}, mm.GetForumResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncDeleteForum, DeleteForumHandler{Server: srv}, mm.DeleteForumParams{}, mm.DeleteForumResult{})
	if err != nil {
		return err
	}

	// Thread.
	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncAddThread, AddThreadHandler{Server: srv}, mm.AddThreadParams{}, mm.AddThreadResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncChangeThreadName, ChangeThreadNameHandler{Server: srv}, mm.ChangeThreadNameParams{}, mm.ChangeThreadNameResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncChangeThreadForum, ChangeThreadForumHandler{Server: srv}, mm.ChangeThreadForumParams{}, mm.ChangeThreadForumResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncGetThread, GetThreadHandler{Server: srv}, mm.GetThreadParams{}, mm.GetThreadResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncDeleteThread, DeleteThreadHandler{Server: srv}, mm.DeleteThreadParams{}, mm.DeleteThreadResult{})
	if err != nil {
		return err
	}

	// Message.
	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncAddMessage, AddMessageHandler{Server: srv}, mm.AddMessageParams{}, mm.AddMessageResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncChangeMessageText, ChangeMessageTextHandler{Server: srv}, mm.ChangeMessageTextParams{}, mm.ChangeMessageTextResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncChangeMessageThread, ChangeMessageThreadHandler{Server: srv}, mm.ChangeMessageThreadParams{}, mm.ChangeMessageThreadResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncGetMessage, GetMessageHandler{Server: srv}, mm.GetMessageParams{}, mm.GetMessageResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncDeleteMessage, DeleteMessageHandler{Server: srv}, mm.DeleteMessageParams{}, mm.DeleteMessageResult{})
	if err != nil {
		return err
	}

	// Composite objects.
	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncListThreadAndMessages, ListThreadAndMessagesHandler{Server: srv}, mm.ListThreadAndMessagesParams{}, mm.ListThreadAndMessagesResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncListThreadAndMessagesOnPage, ListThreadAndMessagesOnPageHandler{Server: srv}, mm.ListThreadAndMessagesOnPageParams{}, mm.ListThreadAndMessagesOnPageResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncListForumAndThreads, ListForumAndThreadsHandler{Server: srv}, mm.ListForumAndThreadsParams{}, mm.ListForumAndThreadsResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncListForumAndThreadsOnPage, ListForumAndThreadsOnPageHandler{Server: srv}, mm.ListForumAndThreadsOnPageParams{}, mm.ListForumAndThreadsOnPageResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncListSectionsAndForums, ListSectionsAndForumsHandler{Server: srv}, mm.ListSectionsAndForumsParams{}, mm.ListSectionsAndForumsResult{})
	if err != nil {
		return err
	}

	// Other.
	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncShowDiagnosticData, ShowDiagnosticDataHandler{Server: srv}, mm.ShowDiagnosticDataParams{}, mm.ShowDiagnosticDataResult{})
	if err != nil {
		return err
	}

	if srv.settings.SystemSettings.IsDebugMode {
		err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncTest, TestHandler{Server: srv}, mm.TestParams{}, mm.TestResult{})
		if err != nil {
			return err
		}
	}

	// Template.
	//err = srv.jsonRpcHandlers.RegisterMethod("Abc", AbcHandler{Server: srv}, models.AbcParams{}, models.AbcResult{})
	//if err != nil {
	//	return err
	//}

	return nil
}
