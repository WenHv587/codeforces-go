// Generated by copypasta/template/generator_test.go
package main

import (
	"github.com/EndlessCheng/codeforces-go/main/testutil"
	"testing"
)

// https://codeforces.com/problemset/problem/1875/C
// https://codeforces.com/problemset/status/1875/problem/C
func Test_cf1875C(t *testing.T) {
	testCases := [][2]string{
		{
			`4
10 5
1 5
10 4
3 4`,
			`0
-1
2
5`,
		},
	}
	testutil.AssertEqualStringCase(t, testCases, 0, cf1875C)
}