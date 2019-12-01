package infrastructure

import (
	"context"
	"github.com/NeekUP/roadmaps/core"
	"time"
)

func NewContext(ctx context.Context) core.ReqContext {
	return &requestContext{ctx: ctx}
}

type requestContext struct {
	ctx context.Context
}

func (this *requestContext) Deadline() (deadline time.Time, ok bool) {
	if this.ctx == nil {
		return time.Now(), false
	}
	return this.ctx.Deadline()
}

func (this *requestContext) Done() <-chan struct{} {
	if this.ctx == nil {
		return nil
	}
	return this.ctx.Done()
}

func (this *requestContext) Err() error {
	if this.ctx == nil {
		return nil
	}
	return this.ctx.Err()
}

func (this *requestContext) Value(key interface{}) interface{} {
	if this.ctx == nil {
		return nil
	}
	return this.ctx.Value(key)
}

func (this *requestContext) ReqId() string {
	if this.ctx == nil {
		return ""
	}
	if reqID, ok := this.ctx.Value(ReqId).(string); ok {
		return reqID
	}
	return ""
}

func (this *requestContext) UserId() string {
	if this.ctx == nil {
		return ""
	}
	if userId, ok := this.ctx.Value(ReqUserId).(string); ok {
		return userId
	}
	return ""
}

func (this *requestContext) UserName() string {
	if this.ctx == nil {
		return ""
	}
	if userName, ok := this.ctx.Value(ReqUserName).(string); ok {
		return userName
	}
	return ""
}
