// Generated by copypasta/template/generator_test.go
package main

import (
	"github.com/EndlessCheng/codeforces-go/main/testutil"
	"testing"
)

// https://codeforces.com/problemset/problem/1967/B1
// https://codeforces.com/problemset/status/1967/problem/B1?friends=on
func Test_cf1967B1(t *testing.T) {
	testCases := [][2]string{
		{
			`6
1 1
2 3
3 5
10 8
100 1233
1000000 1145141`,
			`1
3
4
14
153
1643498`,
		},
	}
	testutil.AssertEqualStringCase(t, testCases, 0, cf1967B1)
}
