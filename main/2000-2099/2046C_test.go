// Generated by copypasta/template/generator_test.go
package main

import (
	"github.com/EndlessCheng/codeforces-go/main/testutil"
	"testing"
)

// https://codeforces.com/contest/2046/problem/C
// https://codeforces.com/problemset/status/2046/problem/C?friends=on
func Test_cf2046C(t *testing.T) {
	testCases := [][2]string{
		{
			`4
4
1 1
1 2
2 1
2 2
4
0 0
0 0
0 0
0 0
8
1 2
2 1
2 -1
1 -2
-1 -2
-2 -1
-2 1
-1 2
7
1 1
1 2
1 3
1 4
2 1
3 1
4 1`,
			`1
2 2
0
0 0
2
1 0
0
0 0`,
		},
		{
			`1
14
-3 -7
6 -4
5 -8
10 6
-1 -10
10 -1
2 -5
-9 -7
-10 -4
4 -7
5 7
-8 4
6 -5
-1 8`,
			`3
2 -4`,
		},
	}
	testutil.AssertEqualStringCase(t, testCases, 0, cf2046C)
}
