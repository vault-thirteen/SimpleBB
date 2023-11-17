package mm

// RPC handlers.

import (
	"context"
	"encoding/json"
	"time"

	js "github.com/osamingo/jsonrpc/v2"
	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/models"
)

// Ping.

type PingHandler struct {
	Server *Server
}

func (h PingHandler) ServeJSONRPC(_ context.Context, _ *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	result := mm.PingResult{OK: true}
	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}

// Section.

type AddSectionHandler struct {
	Server *Server
}

func (h AddSectionHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.AddSectionParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.AddSectionResult
	result, jerr = h.Server.addSection(&p)
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

type ChangeSectionNameHandler struct {
	Server *Server
}

func (h ChangeSectionNameHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ChangeSectionNameParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ChangeSectionNameResult
	result, jerr = h.Server.changeSectionName(&p)
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

type ChangeSectionParentHandler struct {
	Server *Server
}

func (h ChangeSectionParentHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ChangeSectionParentParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ChangeSectionParentResult
	result, jerr = h.Server.changeSectionParent(&p)
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

type GetSectionHandler struct {
	Server *Server
}

func (h GetSectionHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.GetSectionParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.GetSectionResult
	result, jerr = h.Server.getSection(&p)
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

type DeleteSectionHandler struct {
	Server *Server
}

func (h DeleteSectionHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.DeleteSectionParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.DeleteSectionResult
	result, jerr = h.Server.deleteSection(&p)
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

// Forum.

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

type ChangeForumSectionHandler struct {
	Server *Server
}

func (h ChangeForumSectionHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ChangeForumSectionParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ChangeForumSectionResult
	result, jerr = h.Server.changeForumSection(&p)
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

// Thread.

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

// Message.

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

// Composite objects.

type ListThreadAndMessagesHandler struct {
	Server *Server
}

func (h ListThreadAndMessagesHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ListThreadAndMessagesParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ListThreadAndMessagesResult
	result, jerr = h.Server.listThreadAndMessages(&p)
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

type ListThreadAndMessagesOnPageHandler struct {
	Server *Server
}

func (h ListThreadAndMessagesOnPageHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ListThreadAndMessagesOnPageParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ListThreadAndMessagesOnPageResult
	result, jerr = h.Server.listThreadAndMessagesOnPage(&p)
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

type ListForumAndThreadsHandler struct {
	Server *Server
}

func (h ListForumAndThreadsHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ListForumAndThreadsParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ListForumAndThreadsResult
	result, jerr = h.Server.listForumAndThreads(&p)
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

type ListForumAndThreadsOnPageHandler struct {
	Server *Server
}

func (h ListForumAndThreadsOnPageHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ListForumAndThreadsOnPageParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ListForumAndThreadsOnPageResult
	result, jerr = h.Server.listForumAndThreadsOnPage(&p)
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

type ListSectionsAndForumsHandler struct {
	Server *Server
}

func (h ListSectionsAndForumsHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

	var p mm.ListSectionsAndForumsParams
	jerr := js.Unmarshal(params, &p)
	if jerr != nil {
		return nil, jerr
	}

	var result *mm.ListSectionsAndForumsResult
	result, jerr = h.Server.listSectionsAndForums(&p)
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

// Other.

type ShowDiagnosticDataHandler struct {
	Server *Server
}

func (h ShowDiagnosticDataHandler) ServeJSONRPC(_ context.Context, _ *json.RawMessage) (any, *js.Error) {
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

func (h TestHandler) ServeJSONRPC(_ context.Context, params *json.RawMessage) (any, *js.Error) {
	h.Server.diag.IncTotalRequestsCount()
	var timeStart = time.Now()

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

	h.Server.diag.IncSuccessfulRequestsCount()
	return result, nil
}
