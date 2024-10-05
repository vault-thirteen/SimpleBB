package nm

// RPC handlers.

import (
	"encoding/json"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	nm "github.com/vault-thirteen/SimpleBB/pkg/NM/models"
	cs "github.com/vault-thirteen/SimpleBB/pkg/common/settings"
)

func (srv *Server) initRpc() (err error) {
	rpcDurationFieldName := cs.RpcDurationFieldName
	rpcRequestIdFieldName := cs.RpcRequestIdFieldName

	ps := &jrm1.ProcessorSettings{
		CatchExceptions:    true,
		LogExceptions:      true,
		CountRequests:      true,
		DurationFieldName:  &rpcDurationFieldName,
		RequestIdFieldName: &rpcRequestIdFieldName,
	}

	srv.js, err = jrm1.NewProcessor(ps)
	if err != nil {
		return err
	}

	fns := []jrm1.RpcFunction{
		srv.Ping,
		srv.AddNotification,
		srv.AddNotificationS,
		srv.GetNotification,
		srv.GetAllNotifications,
		srv.GetUnreadNotifications,
		srv.CountUnreadNotifications,
		srv.MarkNotificationAsRead,
		srv.DeleteNotification,
		srv.GetDKey,
		srv.ShowDiagnosticData,
		srv.Test,
	}

	for _, fn := range fns {
		err = srv.js.AddFunc(fn)
		if err != nil {
			return err
		}
	}

	return nil
}

// Ping.

func (srv *Server) Ping(_ *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	return nm.PingResult{OK: true}, nil
}

// Notification.

func (srv *Server) AddNotification(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.AddNotificationParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.AddNotificationResult
	r, re = srv.addNotification(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) AddNotificationS(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.AddNotificationSParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.AddNotificationSResult
	r, re = srv.addNotificationS(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetNotification(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.GetNotificationParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.GetNotificationResult
	r, re = srv.getNotification(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetAllNotifications(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.GetAllNotificationsParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.GetAllNotificationsResult
	r, re = srv.getAllNotifications(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) GetUnreadNotifications(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.GetUnreadNotificationsParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.GetUnreadNotificationsResult
	r, re = srv.getUnreadNotifications(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) CountUnreadNotifications(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.CountUnreadNotificationsParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.CountUnreadNotificationsResult
	r, re = srv.countUnreadNotifications(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) MarkNotificationAsRead(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.MarkNotificationAsReadParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.MarkNotificationAsReadResult
	r, re = srv.markNotificationAsRead(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) DeleteNotification(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.DeleteNotificationParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.DeleteNotificationResult
	r, re = srv.deleteNotification(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

// Other.

func (srv *Server) GetDKey(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.GetDKeyParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.GetDKeyResult
	r, re = srv.getDKey(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) ShowDiagnosticData(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.ShowDiagnosticDataParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.ShowDiagnosticDataResult
	r, re = srv.showDiagnosticData()
	if re != nil {
		return nil, re
	}

	return r, nil
}

func (srv *Server) Test(params *json.RawMessage, _ *jrm1.ResponseMetaData) (result any, re *jrm1.RpcError) {
	var p *nm.TestParams
	re = jrm1.ParseParameters(params, &p)
	if re != nil {
		return nil, re
	}

	var r *nm.TestResult
	r, re = srv.test(p)
	if re != nil {
		return nil, re
	}

	return r, nil
}
