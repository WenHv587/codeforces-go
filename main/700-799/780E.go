package main

import (
	"bufio"
	. "fmt"
	"io"
)

// github.com/EndlessCheng/codeforces-go
func CF780E(_r io.Reader, _w io.Writer) {
	in := bufio.NewReader(_r)
	out := bufio.NewWriter(_w)
	defer out.Flush()

	var n, m, k, v, w int
	Fscan(in, &n, &m, &k)
	g := make([][]int, n+1)
	for ; m > 0; m-- {
		Fscan(in, &v, &w)
		g[v] = append(g[v], w)
		g[w] = append(g[w], v)
	}
	vs := []interface{}{}
	vis := make([]bool, n+1)
	var f func(int)
	f = func(v int) {
		vs = append(vs, v)
		vis[v] = true
		for _, w := range g[v] {
			if !vis[w] {
				f(w)
				vs = append(vs, v)
			}
		}
	}
	f(1)
	q, c := (2*n-1)/k, (2*n-1)%k
	for i := 0; i < c; i++ {
		Fprint(out, q+1, " ")
		Fprintln(out, vs[:q+1]...)
		vs = vs[q+1:]
	}
	for i := c; i < k; i++ {
		Fprint(out, q, " ")
		Fprintln(out, vs[:q]...)
		vs = vs[q:]
	}
}

//func main() { CF780E(os.Stdin, os.Stdout) }
