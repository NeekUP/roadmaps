package tests

import "time"

type fakeContext struct{}

func (fakeContext) Deadline() (deadline time.Time, ok bool) {
	panic("implement me")
}

func (fakeContext) Done() <-chan struct{} {
	panic("implement me")
}

func (fakeContext) Err() error {
	panic("implement me")
}

func (fakeContext) Value(key interface{}) interface{} {
	return "seed"
}

func (fakeContext) ReqId() string {
	return ""
}
