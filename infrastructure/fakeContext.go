package infrastructure

import "time"

type FakeContext struct{}

func (FakeContext) Deadline() (deadline time.Time, ok bool) {
	panic("implement me")
}

func (FakeContext) Done() <-chan struct{} {
	panic("implement me")
}

func (FakeContext) Err() error {
	panic("implement me")
}

func (FakeContext) Value(key interface{}) interface{} {
	return "seed"
}

func (FakeContext) ReqId() string {
	return ""
}
