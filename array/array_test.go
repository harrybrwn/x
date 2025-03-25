package array

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestMove(t *testing.T) {
	s := []int{1, 2, 3, 4}
	Move(s, 1, 3)
	exp := []int{1, 3, 4, 2}
	assert(t, slices.Equal(s, exp), "got %v, want %v", s, exp)
}

func assert(t *testing.T, exp bool, msg ...any) {
	t.Helper()
	if !exp {
		if len(msg) > 1 {
			t.Errorf(msg[0].(string), msg[1:]...)
		} else if len(msg) == 1 {
			t.Errorf(
				"failed assertion: %s",
				strings.Join(Map(msg, func(m any) string {
					return fmt.Sprintf("%v", m)
				}), " "),
			)
		} else {
			t.Error("failed assertion")
		}
	}
}
