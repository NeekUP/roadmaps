package tests

import (
	"sort"
	"testing"

	"github.com/NeekUP/roadmaps/core"
	"github.com/NeekUP/roadmaps/domain"
	"github.com/r3labs/diff"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestNumToStringSuccess(t *testing.T) {
	for i := 0; i < 10000; i++ {
		s := core.EncodeNumToString(i)
		if num, err := core.DecodeStringToNum(s); err != nil || num != i {
			t.Errorf("Fail to convert int to string for url: %d => %s => %d", i, s, num)
		}
	}

}
func TestDeff(t *testing.T) {
	text1 := `You want to contribute to go-diff? GREAT! 
	If you are here because of a bug you want to fix or 
	a feature you want to add, you can just read on. 
	Otherwise we have a list of open issues in the tracker. 
	Just choose something you think you can work 
	on and discuss your plans in the issue by commenting on it.

	Please make sure that every behavioral change is accompanied by test cases. 
	Additionally, every contribution must pass the lint and test 
	Makefile targets which can be run using the following commands 
	in the repository root directory.`
	text2 := `You want to contribute to go-diff? GREAT! 
	If you are here because of a bug you want to fix or 
	a feature you want to add, you can just read on. 
	Otherwise we have a list of open issues in the tracker. 
	Just choose something you think you can work 
	on and discuss your plans in the issue by commenting on it.

	Please make sure that every behavioral change is accompanied by test cases. 
	Additionally, every contribution must pass the lint and test 
	Makefile targets which can be run using the following commands 
	in the repository root directory.`
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text1, text2, false)
	t.Errorf(dmp.DiffPrettyText(diffs))
}

func TestDiff(t *testing.T) {
	from := make([]domain.Step, 3)
	to := make([]domain.Step, 3)

	from[0] = domain.Step{Id: 1, PlanId: 1, ReferenceType: domain.ResourceReference, ReferenceId: 1, Position: 0}
	from[1] = domain.Step{Id: 2, PlanId: 2, ReferenceType: domain.ResourceReference, ReferenceId: 2, Position: 1}
	from[2] = domain.Step{Id: 3, PlanId: 2, ReferenceType: domain.ResourceReference, ReferenceId: 3, Position: 2}

	to[0] = domain.Step{Id: 3, PlanId: 3, ReferenceType: domain.ResourceReference, ReferenceId: 4, Position: 2}
	to[1] = domain.Step{Id: 1, PlanId: 1, ReferenceType: domain.ResourceReference, ReferenceId: 1, Position: 0}
	to[2] = domain.Step{Id: 2, PlanId: 2, ReferenceType: domain.TopicReference, ReferenceId: 2, Position: 1}

	sort.Slice(from, func(i, j int) bool { return from[i].Position < from[j].Position })
	sort.Slice(to, func(i, j int) bool { return to[i].Position < to[j].Position })

	changelog, _ := diff.Diff(from, to)

	t.Errorf("%#v", changelog)
}
