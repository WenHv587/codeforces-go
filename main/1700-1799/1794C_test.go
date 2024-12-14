// Generated by copypasta/template/generator_test.go
package main

import (
	"github.com/EndlessCheng/codeforces-go/main/testutil"
	"testing"
)

// https://codeforces.com/problemset/problem/1794/C
// https://codeforces.com/problemset/status/1794/problem/C
func Test_cf1794C(t *testing.T) {
	testCases := [][2]string{
		{
			`3
3
1 2 3
2
1 1
5
5 5 5 5 5`,
			`1 1 2 
1 1 
1 2 3 4 5`,
		},
	}
	testutil.AssertEqualStringCase(t, testCases, 0, cf1794C)
}
