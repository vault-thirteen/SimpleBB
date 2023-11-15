package mm

// RPC handlers.

import (
	"context"
	"encoding/json"
	"time"

	js "github.com/osamingo/jsonrpc/v2"
	mc "github.com/vault-thirteen/SimpleBB/pkg/MM/client"
	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/models"
)

func (srv *Server) initJsonRpcHandlers() (err error) {
	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncPing, PingHandler{Server: srv}, mm.PingParams{}, mm.PingResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncAddForum, AddForumHandler{Server: srv}, mm.AddForumParams{}, mm.AddForumResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncChangeForumName, ChangeForumNameHandler{Server: srv}, mm.ChangeForumNameParams{}, mm.ChangeForumNameResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncChangeForumParent, ChangeForumParentHandler{Server: srv}, mm.ChangeForumParentParams{}, mm.ChangeForumParentResult{})
	if err != nil {
		return err
	}

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

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncDeleteMessage, DeleteMessageHandler{Server: srv}, mm.DeleteMessageParams{}, mm.DeleteMessageResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncGetMessage, GetMessageHandler{Server: srv}, mm.GetMessageParams{}, mm.GetMessageResult{})
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

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncGetForum, GetForumHandler{Server: srv}, mm.GetForumParams{}, mm.GetForumResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncDeleteForum, DeleteForumHandler{Server: srv}, mm.DeleteForumParams{}, mm.DeleteForumResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncListThreadMessages, ListThreadMessagesHandler{Server: srv}, mm.ListThreadMessagesParams{}, mm.ListThreadMessagesResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncListThreadMessagesOnPage, ListThreadMessagesOnPageHandler{Server: srv}, mm.ListThreadMessagesOnPageParams{}, mm.ListThreadMessagesOnPageResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncListForumThreads, ListForumThreadsHandler{Server: srv}, mm.ListForumThreadsParams{}, mm.ListForumThreadsResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncListForumThreadsOnPage, ListForumThreadsOnPageHandler{Server: srv}, mm.ListForumThreadsOnPageParams{}, mm.ListForumThreadsOnPageResult{})
	if err != nil {
		return err
	}

	err = srv.jsonRpcHandlers.RegisterMethod(mc.FuncListForums, ListForumsHandler{Server: srv}, mm.ListForumsParams{}, mm.ListForumsResult{})
	if err != nil {
		return err
	}

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

type PingHandler struct {
	Server *Server
}

func (h PingHandler) ServeJSONRPC(_ context.Context, _ *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	result := mm.PingResult{OK: true}
	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type AddForumHandler struct {
	Server *Server
}

func (h AddForumHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.AddForumParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.AddForumResult
	result, jerr = h.Server.addForum(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type ChangeForumNameHandler struct {
	Server *Server
}

func (h ChangeForumNameHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ChangeForumNameParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ChangeForumNameResult
	result, jerr = h.Server.changeForumName(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type ChangeForumParentHandler struct {
	Server *Server
}

func (h ChangeForumParentHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ChangeForumParentParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ChangeForumParentResult
	result, jerr = h.Server.changeForumParent(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type AddThreadHandler struct {
	Server *Server
}

func (h AddThreadHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.AddThreadParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.AddThreadResult
	result, jerr = h.Server.addThread(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type ChangeThreadNameHandler struct {
	Server *Server
}

func (h ChangeThreadNameHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ChangeThreadNameParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ChangeThreadNameResult
	result, jerr = h.Server.changeThreadName(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type ChangeThreadForumHandler struct {
	Server *Server
}

func (h ChangeThreadForumHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ChangeThreadForumParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ChangeThreadForumResult
	result, jerr = h.Server.changeThreadForum(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type AddMessageHandler struct {
	Server *Server
}

func (h AddMessageHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.AddMessageParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.AddMessageResult
	result, jerr = h.Server.addMessage(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type ChangeMessageTextHandler struct {
	Server *Server
}

func (h ChangeMessageTextHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ChangeMessageTextParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ChangeMessageTextResult
	result, jerr = h.Server.changeMessageText(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type ChangeMessageThreadHandler struct {
	Server *Server
}

func (h ChangeMessageThreadHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ChangeMessageThreadParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ChangeMessageThreadResult
	result, jerr = h.Server.changeMessageThread(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type DeleteMessageHandler struct {
	Server *Server
}

func (h DeleteMessageHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.DeleteMessageParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.DeleteMessageResult
	result, jerr = h.Server.deleteMessage(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type GetMessageHandler struct {
	Server *Server
}

func (h GetMessageHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.GetMessageParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.GetMessageResult
	result, jerr = h.Server.getMessage(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type GetThreadHandler struct {
	Server *Server
}

func (h GetThreadHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.GetThreadParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.GetThreadResult
	result, jerr = h.Server.getThread(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type DeleteThreadHandler struct {
	Server *Server
}

func (h DeleteThreadHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.DeleteThreadParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.DeleteThreadResult
	result, jerr = h.Server.deleteThread(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type GetForumHandler struct {
	Server *Server
}

func (h GetForumHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.GetForumParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.GetForumResult
	result, jerr = h.Server.getForum(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type DeleteForumHandler struct {
	Server *Server
}

func (h DeleteForumHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.DeleteForumParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.DeleteForumResult
	result, jerr = h.Server.deleteForum(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type ListThreadMessagesHandler struct {
	Server *Server
}

func (h ListThreadMessagesHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ListThreadMessagesParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ListThreadMessagesResult
	result, jerr = h.Server.listThreadMessages(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type ListThreadMessagesOnPageHandler struct {
	Server *Server
}

func (h ListThreadMessagesOnPageHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ListThreadMessagesOnPageParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ListThreadMessagesOnPageResult
	result, jerr = h.Server.listThreadMessagesOnPage(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type ListForumThreadsHandler struct {
	Server *Server
}

func (h ListForumThreadsHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ListForumThreadsParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ListForumThreadsResult
	result, jerr = h.Server.listForumThreads(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type ListForumThreadsOnPageHandler struct {
	Server *Server
}

func (h ListForumThreadsOnPageHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ListForumThreadsOnPageParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ListForumThreadsOnPageResult
	result, jerr = h.Server.listForumThreadsOnPage(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type ListForumsHandler struct {
	Server *Server
}

func (h ListForumsHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ListForumsParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ListForumsResult
	result, jerr = h.Server.listForums(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type ShowDiagnosticDataHandler struct {
	Server *Server
}

func (h ShowDiagnosticDataHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	result, jerr := h.Server.showDiagnosticData()
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

type TestHandler struct {
	Server *Server
}

func (h TestHandler) ServeJSONRPC(ctx context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	//requestId, jerr := c.GetRequestId(ctx)
	//if jerr != nil {
	//	return nil, jerr
	//}
	//log.Println("Test. Request start. ID:", requestId) //TODO

	var p mm.TestParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.TestResult
	result, jerr = h.Server.doTest(&p)
	if jerr != nil {
		return nil, jerr
	}

	var taskDuration = time.Now().Sub(timeStart).Milliseconds()
	if result != nil {
		result.TimeSpent = taskDuration
	}

	//log.Println("Test. Request end. ID:", requestId) //TODO

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}
