package common_test

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	"github.com/opslevel/kubectl-opslevel/common"
)

func TestGetSample(t *testing.T) {
	tests := []struct {
		n    int
		data []string
		want []string
	}{
		{-1, []string{"A", "B"}, []string{"A", "B"}}, // edge case N < 0
		{0, []string{"A", "B"}, []string{"A", "B"}},  // edge case N < 1
		{2, []string{"A", "B"}, []string{"A", "B"}},  // edge case N == len(data)
		{0, []string{}, []string{}},                  // edge case N == len(data) == 0
		{1, []string{"A"}, []string{"A"}},            // edge case N == len(data) == 1
		{3, []string{"A", "B"}, []string{"A", "B"}},  // edge case N > len(data)

		// test case slices must be written in ascending alphabetical order without
		// any duplicates otherwise verifying correctness will be impossible.
		// if the function is correct, C could never be before A or B.
		{5, []string{"A", "B", "C", "D", "E", "F", "G"}, nil},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("n is %d slice is %v want", tt.n, tt.data)
		if tt.want != nil {
			testname += fmt.Sprintf(" %v", tt.want)
		} else {
			testname += " ordered output."
		}

		t.Run(testname, func(t *testing.T) {
			result := common.GetSample(tt.n, tt.data)

			if len(result) != tt.n && tt.n > 1 && tt.n < len(result) {
				t.Error("incorrect length")
			}

			if tt.want != nil {
				if !reflect.DeepEqual(tt.want, result) {
					t.Error("edge case: result != want")
				}
			} else {
				// applies to test cases only - should not have duplicate elements
				entries := map[string]struct{}{}
				for i, elem := range result {
					if _, ok := entries[elem]; ok {
						t.Errorf("result elem %v pos %d not unique", elem, i)
					}
					entries[elem] = struct{}{}
				}
				deepCopy := append([]string(nil), result...)
				slices.Sort(deepCopy)
				for i := range result {
					a, b := deepCopy[i], result[i]
					if a != b {
						t.Errorf("result is not sorted, got sorted value '%s' vs result '%s", a, b)
					}
				}
			}

			// addresses should not match (returns a shallow pointer)
			if &tt.data == &result {
				t.Errorf("data pointer %p should not match result pointer %p", &tt.data, &result)
			}
		})
	}
}
