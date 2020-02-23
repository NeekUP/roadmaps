package infrastructure

import (
	"context"
	"github.com/NeekUP/nptrace"
	"github.com/NeekUP/roadmaps/core"
	"time"
)

func NewContext(ctx context.Context) core.ReqContext {
	return &requestContext{ctx: ctx}
}

type requestContext struct {
	ctx context.Context
}

func (reqCtx *requestContext) Deadline() (deadline time.Time, ok bool) {
	if reqCtx.ctx == nil {
		return time.Now(), false
	}
	return reqCtx.ctx.Deadline()
}

func (reqCtx *requestContext) Done() <-chan struct{} {
	if reqCtx.ctx == nil {
		return nil
	}
	return reqCtx.ctx.Done()
}

func (reqCtx *requestContext) Err() error {
	if reqCtx.ctx == nil {
		return nil
	}
	return reqCtx.ctx.Err()
}

func (reqCtx *requestContext) Value(key interface{}) interface{} {
	if reqCtx.ctx == nil {
		return nil
	}
	return reqCtx.ctx.Value(key)
}

func (reqCtx *requestContext) ReqId() string {
	if reqCtx.ctx == nil {
		return ""
	}
	if reqID, ok := reqCtx.ctx.Value(ReqId).(string); ok {
		return reqID
	}
	return ""
}

func (reqCtx *requestContext) UserId() string {
	if reqCtx.ctx == nil {
		return ""
	}
	if userId, ok := reqCtx.ctx.Value(ReqUserId).(string); ok {
		return userId
	}
	return ""
}

func (reqCtx *requestContext) UserName() string {
	if reqCtx.ctx == nil {
		return ""
	}
	if userName, ok := reqCtx.ctx.Value(ReqUserName).(string); ok {
		return userName
	}
	return ""
}

func (reqCtx *requestContext) StartTrace(name string, args ...interface{}) *nptrace.Trace {
	tr, ok := reqCtx.Value(Tracer).(*nptrace.Task)
	if ok {
		return tr.Start(name, args...)
	}
	return nil
}

func (reqCtx *requestContext) StopTrace(t *nptrace.Trace) {
	if t == nil {
		return
	}
	reqCtx.Value(Tracer).(*nptrace.Task).Stop(t)
}
