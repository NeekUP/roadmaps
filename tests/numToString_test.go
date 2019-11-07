package tests

import (
	"roadmaps/core"
	"testing"
)

func TestNumToStringSuccess(t *testing.T) {
	for i := 0; i < 10000; i++ {
		s := core.EncodeNumToString(i)
		if num, err := core.DecodeStringToNum(s); err != nil || num != i {
			t.Errorf("Fail to convert int to string for url: %d => %s => %d", i, s, num)
		}
	}

}
