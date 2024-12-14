// Generated by copypasta/template/generator_test.go
package main

import (
	"github.com/EndlessCheng/codeforces-go/main/testutil"
	"testing"
)

// https://codeforces.com/problemset/problem/379/F
// https://codeforces.com/problemset/status/379/problem/F?friends=on
func Test_cf379F(t *testing.T) {
	testCases := [][2]string{
		{
			`5
2
3
4
8
5`,
			`3
4
4
5
6`,
		},
	}
	testutil.AssertEqualStringCase(t, testCases, 0, cf379F)
}
