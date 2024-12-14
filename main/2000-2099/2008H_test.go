// Generated by copypasta/template/generator_test.go
package main

import (
	"github.com/EndlessCheng/codeforces-go/main/testutil"
	"testing"
)

// https://codeforces.com/contest/2008/problem/H
// https://codeforces.com/problemset/status/2008/problem/H
func Test_cf2008H(t *testing.T) {
	testCases := [][2]string{
		{
			`2
5 5
1 2 3 4 5
1
2
3
4
5
6 3
1 2 6 4 1 3
2
1
5`,
			`0 1 1 1 2 
1 0 2`,
		},
	}
	testutil.AssertEqualStringCase(t, testCases, 0, cf2008H)
}
