// Code generated by copypasta/template/leetcode/generator_test.go
package main

import (
	"github.com/EndlessCheng/codeforces-go/leetcode/testutil"
	"testing"
)

func Test(t *testing.T) {
	t.Log("Current test is [d]")
	examples := [][]string{
		{
			`"327"`, 
			`2`,
		},
		{
			`"094"`, 
			`0`,
		},
		{
			`"0"`, 
			`0`,
		},
		{
			`"9999999999999"`, 
			`101`,
		},
		{
			`"412"`,
			`2`,
		},
		{
			`"417"`,
			`2`,
		},
	}
	targetCaseNum := -1
	if err := testutil.RunLeetCodeFuncWithExamples(t, numberOfCombinations, examples, targetCaseNum); err != nil {
		t.Fatal(err)
	}
}
// https://leetcode-cn.com/contest/biweekly-contest-59/problems/number-of-ways-to-separate-numbers/
