package copypasta

import (
	. "fmt"
	"io"
	"math"
	"math/bits"
	"math/rand"
	"sort"
)

// 从数据范围找思路：
// 1e9~1e18 √n logn 1     二分 二进制
// 1e5~1e6  nlogn nαn n   二分 RMQ 并查集
// 1e3~1e4  n^2 n√n       RMQ DP 分块
// 300~500  n^3           DP 二分图

// General ideas https://codeforces.ml/blog/entry/48417
// 从特殊到一般：尝试修改条件或缩小题目的数据范围，先研究某个特殊情况下的思路，然后再逐渐扩大数据范围来思考怎么改进算法

// 异类双变量：固定某变量统计另一变量的 [0,n)
// 同类双变量①：固定 i 统计 [0,n)
// 同类双变量②：固定 i 统计 [0,i-1]
// 套路：预处理数据（按照某种顺序排序/优先队列/BST/...），或者边遍历边维护，
//      然后固定变量 i，用均摊 O(1)~O(logn) 的复杂度统计范围内的另一变量 j
// 这样可以将复杂度从 O(n^2) 降低到 O(n) 或 O(nlogn)

// NOTE: 正难则反。 all => any, any => all https://codeforces.ml/problemset/problem/621/C
// NOTE: 子区间和为 0 => 出现了两个同样的前缀和。这种题目建议下标从 1 开始，见 https://codeforces.ml/problemset/problem/1333/C
// NOTE: 和式的另一视角：若每一项的值都在一个范围，不妨考虑另一个问题 - 值为 x 的项有多少个？https://atcoder.jp/contests/abc162/tasks/abc162_e
// NOTE: 变换考察角度：对所有排列考察所有子区间的性质，可以转换成对所有子区间考察所有排列，将子区间内部的排列和区间外部的排列进行区分，内部的性质单独研究，外部的当作 (n-(r-l))! 个排列 https://codeforces.ml/problemset/problem/1284/C

// NOTE: 若不止两个数相加，要特别注意 inf 的选择
// 一个 Golang 的注意事项：forr array 时，遍历 i 时修改 i 后面的元素的值是不影响 ai 的，只能用 for+a[i] 获取
func commonCollection() {
	// HELPER
	const mod int64 = 1e9 + 7 // 998244353
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	pow2 := [...]int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144}
	pow10 := [...]int{1, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9} // math.Pow10
	factorial := [...]int{1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880, 3628800 /*10!*/, 39916800, 479001600}
	// TIPS: dir4[i] 和 dir4[i^1] 互为相反方向
	type pair struct{ x, y int }
	dir4 := [...]pair{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} // 上下左右
	dir4C := [...]pair{ // 西东南北
		'W': {-1, 0},
		'E': {1, 0},
		'S': {0, -1},
		'N': {0, 1},
	}
	dir4c := [...]pair{ // 左右下上
		'L': {-1, 0},
		'R': {1, 0},
		'D': {0, -1},
		'U': {0, 1},
	}
	dir4R := [...]pair{{1, 1}, {-1, 1}, {-1, -1}, {1, -1}}
	dir8 := [...]pair{{1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}, {0, -1}, {1, -1}}
	orderP3 := [6][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}

	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	mins := func(vals ...int) int {
		ans := vals[0]
		for _, val := range vals[1:] {
			if val < ans {
				ans = val
			}
		}
		return ans
	}
	maxs := func(vals ...int) int {
		ans := vals[0]
		for _, val := range vals[1:] {
			if val > ans {
				ans = val
			}
		}
		return ans
	}
	sort3 := func(a ...int) (x, y, z int) { sort.Ints(a); return a[0], a[1], a[2] }
	// 用堆求前 k 小
	smallK := func(a []int, k int) []int {
		k++
		q := hp{} // 最大堆
		for _, v := range a {
			if q.Len() < k || v < q.top() {
				q.push(v)
			}
			if q.Len() > k {
				q.pop() // 不断弹出更大的元素，留下的就是较小的
			}
		}
		return q.IntSlice // 注意返回的不是有序数组
	}
	ternaryI := func(cond bool, r1, r2 int) int {
		if cond {
			return r1
		}
		return r2
	}
	ternaryS := func(cond bool, r1, r2 string) string {
		if cond {
			return r1
		}
		return r2
	}
	toInts := func(s []byte) []int {
		ints := make([]int, len(s))
		for i, b := range s {
			ints[i] = int(b)
		}
		return ints
	}
	xor := func(b1, b2 bool) bool { return b1 && !b2 || !b1 && b2 }
	zip := func(a, b []int) {
		n := len(a)
		type pair struct{ x, y int }
		ps := make([]pair, n)
		for i := range ps {
			ps[i] = pair{a[i], b[i]}
		}
	}
	zipI := func(a []int) {
		n := len(a)
		type pair struct{ x, y int }
		ps := make([]pair, n)
		for i := range ps {
			ps[i] = pair{a[i], i}
		}
	}
	getCol := func(mat [][]int, j int) (col []int) {
		for _, row := range mat {
			col = append(col, row[j])
		}
		return
	}
	minString := func(a, b string) string {
		if len(a) != len(b) {
			if len(a) < len(b) {
				return a
			}
			return b
		}
		if a < b {
			return a
		}
		return b
	}
	removeLeadingZero := func(s string) string {
		for i, b := range s {
			if b > '0' {
				return s[i:]
			}
		}
		return "0"
	}
	// END HELPER

	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}
	absAll := func(a []int) {
		for i, v := range a {
			if v < 0 {
				a[i] = -v
			}
		}
	}

	// https://en.wikipedia.org/wiki/Exponentiation_by_squaring
	pow := func(x int64, n int, mod int64) int64 {
		x %= mod
		res := int64(1) % mod
		for ; n > 0; n >>= 1 {
			if n&1 == 1 {
				res = res * x % mod
			}
			x = x * x % mod
		}
		return res
	}

	calcFactorial := func(n int) int64 {
		ans := int64(1)
		for i := 2; i <= n; i++ {
			ans *= int64(i)
		}
		return ans
	}

	// 从低位到高位
	toAnyBase := func(x, base int) (res []int) {
		for ; x > 0; x /= base {
			res = append(res, x%base)
		}
		return
	}
	digits := func(x int) (res []int) {
		for ; x > 0; x /= 10 {
			res = append(res, x%10)
		}
		return
	}

	var sum2d [][]int
	initSum2D := func(mat [][]int) {
		n, m := len(mat), len(mat[0])
		sum2d = make([][]int, n+1)
		sum2d[0] = make([]int, m+1)
		for i, row := range mat {
			sum2d[i+1] = make([]int, m+1)
			for j, v := range row {
				sum2d[i+1][j+1] = sum2d[i+1][j] + sum2d[i][j+1] - sum2d[i][j] + v
			}
		}
	}
	// r1<=r<=r2 && c1<=c<=c2
	querySum2D := func(r1, c1, r2, c2 int) int {
		r2++
		c2++
		return sum2d[r2][c2] - sum2d[r2][c1] - sum2d[r1][c2] + sum2d[r1][c1]
	}

	// 启发式合并：map 版
	mergeMap := func(a, b map[int]int) map[int]int {
		if len(a) < len(b) {
			a, b = b, a
		}
		for k, v := range b {
			a[k] += v
		}
		return a
	}

	//

	copyMat := func(mat [][]int) [][]int {
		n, m := len(mat), len(mat[0])
		dst := make([][]int, n)
		for i, row := range mat {
			dst[i] = make([]int, m)
			copy(dst[i], row)
		}
		return dst
	}

	hash01Mat := func(mat [][]int) int {
		hash := 0
		cnt := 0
		for _, row := range mat {
			for _, v := range row {
				hash |= v << cnt
				cnt++
			}
		}
		return hash
	}

	reverse := func(a []byte) []byte {
		n := len(a)
		r := make([]byte, n)
		for i, v := range a {
			r[n-1-i] = v
		}
		return r
	}
	reverseSelf := func(s []byte) {
		for i, j := 0, len(s)-1; i < j; {
			s[i], s[j] = s[j], s[i]
			i++
			j--
		}
	}

	equals := func(a, b []int) bool {
		// assert len(a) == len(b)
		for i, v := range a {
			if v != b[i] {
				return false
			}
		}
		return true
	}

	// 合并有序数组，保留重复元素
	// a b 必须是有序的（可以为空）
	merge := func(a, b []int) []int {
		i, n := 0, len(a)
		j, m := 0, len(b)
		res := make([]int, 0, n+m)
		for {
			if i == n {
				return append(res, b[j:]...)
			}
			if j == m {
				return append(res, a[i:]...)
			}
			if a[i] < b[j] { // 改成 > 为降序
				res = append(res, a[i])
				i++
			} else {
				res = append(res, b[j])
				j++
			}
		}
	}

	// 求差集 A-B, B-A 和交集 A∩B
	// EXTRA: 求并集 union: A∪B = A-B+A∩B = merge(differenceA, intersection) 或 merge(differenceB, intersection)
	// EXTRA: 求对称差 symmetric_difference: A▲B = A-B ∪ B-A = merge(differenceA, differenceB)
	// a b 必须是有序的（可以为空）
	// 与图论结合 https://codeforces.ml/problemset/problem/243/B
	splitDifferenceAndIntersection := func(a, b []int) (differenceA, differenceB, intersection []int) {
		i, n := 0, len(a)
		j, m := 0, len(b)
		for {
			if i == n {
				differenceB = append(differenceB, b[j:]...)
				return
			}
			if j == m {
				differenceA = append(differenceA, a[i:]...)
				return
			}
			x, y := a[i], b[j]
			if x < y { // 改成 > 为降序
				differenceA = append(differenceA, x)
				i++
			} else if x > y { // 改成 < 为降序
				differenceB = append(differenceB, y)
				j++
			} else {
				intersection = append(intersection, x)
				i++
				j++
			}
		}
	}

	// a 是否为 b 的子集（相当于 differenceA 为空）
	// a b 需要是有序的
	isSubset := func(a, b []int) bool {
		i, n := 0, len(a)
		j, m := 0, len(b)
		for {
			if i == n {
				return true
			}
			if j == m {
				return false
			}
			x, y := a[i], b[j]
			if x < y { // 改成 > 为降序
				return false
			} else if x > y { // 改成 < 为降序
				j++
			} else {
				i++
				j++
			}
		}
	}

	// 是否为不相交集合（相当于 intersection 为空）
	// a b 需要是有序的
	isDisjoint := func(a, b []int) bool {
		i, n := 0, len(a)
		j, m := 0, len(b)
		for {
			if i == n || j == m {
				return true
			}
			x, y := a[i], b[j]
			if x < y { // 改成 > 为降序
				i++
			} else if x > y { // 改成 < 为降序
				j++
			} else {
				return false
			}
		}
	}

	// a 必须是有序的
	unique := func(a []int) (res []int) {
		n := len(a)
		if n == 0 {
			return
		}
		res = make([]int, 1, n)
		res[0] = a[0]
		for i := 1; i < n; i++ {
			if a[i] != a[i-1] {
				res = append(res, a[i])
			}
		}
		//n = len(res)
		return
	}

	uniqueInPlace := func(a []int) []int {
		n := len(a)
		if n == 0 {
			return nil
		}
		j := 0
		for i := 1; i < n; i++ {
			if a[j] != a[i] {
				j++
				a[j] = a[i]
			}
		}
		//n = j + 1
		return a[:j+1]
	}

	// 离散化 discrete([]int{100,20,50,50}, 1) => []int{3,1,2,2}
	// 相当于转换成第几小
	// 若允许修改原数组，可以先将其排序去重后，再调用 discrete，注意去重后 n 需要重新赋值
	discrete := func(a []int, startIndex int) (kth []int) {
		// assert len(a) > 0
		type pair struct{ v, i int }
		n := len(a)
		ps := make([]pair, n)
		for i, v := range a {
			ps[i] = pair{v, i}
		}
		sort.Slice(ps, func(i, j int) bool { return ps[i].v < ps[j].v }) // or SliceStable
		kth = make([]int, n)

		// a 有重复元素
		k := startIndex
		kth[ps[0].i] = k
		for i := 1; i < n; i++ {
			if ps[i].v != ps[i-1].v {
				k++
			}
			kth[ps[i].i] = k
		}

		// a 无重复元素
		for i, p := range ps {
			kth[p.i] = i + startIndex
		}

		return
	}

	// 离散化 discreteMap([]int{100,20,50,50}, 1) => map[int]int{100:3, 20:1, 50:2}
	// 若允许修改原数组，可以先将其排序去重后，再调用 discreteMap，注意去重后 n 需要重新赋值
	discreteMap := func(a []int, startIndex int) (kth map[int]int) {
		// assert len(a) > 0
		n := len(a)
		b := make([]int, n)
		copy(b, a)
		sort.Ints(b)

		// 有重复元素
		k := startIndex
		kth = map[int]int{b[0]: k}
		for i := 1; i < n; i++ {
			if b[i] != b[i-1] {
				k++
				kth[b[i]] = k
			}
		}

		// 无重复元素
		kth = make(map[int]int, n)
		for i, v := range b {
			kth[v] = i + startIndex
		}

		return
	}

	// 哈希编号，也可以理解成另一种离散化（无序）
	// 编号从 0 开始
	indexMap := func(a []string) map[string]int {
		mp := map[string]int{}
		for _, v := range a {
			if _, ok := mp[v]; !ok {
				mp[v] = len(mp)
			}
		}
		return mp
	}

	allSame := func(a ...int) bool {
		for _, v := range a[1:] {
			if v != a[0] {
				return false
			}
		}
		return true
	}

	// a 相对于 [0,n) 的补集
	// a 必须是升序且无重复元素
	complement := func(n int, a []int) (res []int) {
		j := 0
		for i := 0; i < n; i++ {
			if j == len(a) || i < a[j] {
				res = append(res, i)
			} else {
				j++
			}
		}
		return
	}

	// 逆序数
	var mergeCount func([]int) int64
	mergeCount = func(a []int) int64 {
		n := len(a)
		if n <= 1 {
			return 0
		}
		b := make([]int, n/2)
		c := make([]int, n-n/2)
		copy(b, a[:n/2])
		copy(c, a[n/2:])
		cnt := mergeCount(b) + mergeCount(c)
		ai, bi, ci := 0, 0, 0
		for ai < n {
			// 归并排序的同时计算逆序数
			if bi < len(b) && (ci == len(c) || b[bi] <= c[ci]) {
				a[ai] = b[bi]
				bi++
			} else {
				cnt += int64(n/2 - bi)
				a[ai] = c[ci]
				ci++
			}
			ai++
		}
		return cnt
	}

	// 数组第 k 小 (Quick Select)
	// 0 <= k < len(a)
	// 调用会改变数组中元素顺序
	// 代码实现参考算法第四版 p.221
	// 算法的平均比较次数为 ~2n+2kln(n/k)+2(n-k)ln(n/(n-k))
	// https://en.wikipedia.org/wiki/Quickselect
	// https://www.geeksforgeeks.org/quickselect-algorithm/
	// 模板题 https://leetcode-cn.com/problems/kth-largest-element-in-an-array/
	// 模板题 https://codeforces.ml/contest/977/problem/C
	quickSelect := func(a []int, k int) int {
		//k = len(a) - 1 - k // 求第 k 大
		rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
		for l, r := 0, len(a)-1; l < r; {
			v := a[l] // 切分元素
			i, j := l, r+1
			for {
				for i++; i < r && a[i] < v; i++ {
				}
				for j--; j > l && a[j] > v; j-- {
				}
				if i >= j {
					break
				}
				a[i], a[j] = a[j], a[i]
			}
			a[l], a[j] = a[j], v
			if j == k {
				break
			} else if j < k {
				l = j + 1
			} else {
				r = j - 1
			}
		}
		return a[k] //  a[:k+1]  a[k:]
	}

	contains := func(a []int, x int) bool {
		for _, v := range a {
			if v == x {
				return true
			}
		}
		return false
	}

	// x 是否包含 y 中的所有元素，且顺序一致
	containsAll := func(x, y []int) bool {
		for len(y) < len(x) {
			if len(y) == 0 {
				return true
			}
			if x[0] == y[0] {
				y = y[1:]
			}
			x = x[1:]
		}
		return false
	}

	//

	// 判环
	// 1<=next[i]<=n
	getCycle := func(next []int, n, st int) (beforeCycle, cycle []int) {
		vis := make([]int8, n+1)
		for v := st; vis[v] < 2; v = next[v] {
			if vis[v] == 1 {
				cycle = append(cycle, v)
			}
			vis[v]++
		}
		for v := 1; vis[v] == 1; v = next[v] {
			beforeCycle = append(beforeCycle, v)
		}
		return
	}

	// 算法导论 练习4.1-5
	maxSubArraySum := func(a []int) int {
		curSum, maxSum := a[0], a[0]
		for _, v := range a[1:] {
			curSum = max(curSum+v, v)
			maxSum = max(maxSum, curSum)
		}
		return maxSum
	}

	maxSubArrayAbsSum := func(a []int) int {
		//min, max, abs := math.Min, math.Max, math.Abs
		curMaxSum, maxSum := a[0], a[0]
		curMinSum, minSum := a[0], a[0]
		for _, v := range a[1:] {
			curMaxSum = max(curMaxSum+v, v)
			maxSum = max(maxSum, curMaxSum)
			curMinSum = min(curMinSum+v, v)
			minSum = min(minSum, curMinSum)
		}
		return max(abs(maxSum), abs(minSum))
	}

	// 扫描线
	// https://cses.fi/book/book.pdf 30.1
	// TODO 窗口的星星 https://www.luogu.com.cn/problem/P1502
	// 天际线问题 https://leetcode-cn.com/problems/the-skyline-problem/
	// TODO 矩形面积并 https://leetcode-cn.com/problems/rectangle-area-ii/ 《算法与实现》5.4.3
	// 经典题 https://codeforces.ml/problemset/problem/1000/C
	// LC 套题 https://leetcode-cn.com/tag/line-sweep/
	sweepLine := func(in io.Reader, n int) {
		type event struct{ pos, delta int }
		events := make([]event, 0, 2*n)
		for i := 0; i < n; i++ {
			var l, r int
			Fscan(in, &l, &r)
			events = append(events, event{l, 1}, event{r, -1})
		}
		sort.Slice(events, func(i, j int) bool {
			a, b := events[i], events[j]
			return a.pos < b.pos || a.pos == b.pos && a.delta < b.delta // < 先出后进；> 先进后出
		})

		for _, e := range events {
			_ = e
		}
	}

	// 悬线法
	// 求一最大子矩形，矩形内部元素均相同
	// todo https://oi-wiki.org/misc/largest-matrix/

	// 从 st 出发，步长为 gap，不超过 upper 的最大值
	// st <= upper, gap > 0
	maxValueStepToUpper := func(st, upper, gap int) int {
		upper -= st
		return st + upper - upper%gap
	}

	// 二维离散化
	// 代码来源 https://atcoder.jp/contests/abc168/tasks/abc168_f
	discrete2D := func(n, m int) (ans int) {
		type line struct{ a, b, c int }
		lr := make([]line, n)
		du := make([]line, m)
		// read ...

		xs := []int{-2e9, 0, 2e9}
		ys := []int{-2e9, 0, 2e9}
		for _, l := range lr {
			a, b, c := l.a, l.b, l.c
			xs = append(xs, a, b)
			ys = append(ys, c)
		}
		for _, l := range du {
			a, b, c := l.a, l.b, l.c
			xs = append(xs, a)
			ys = append(ys, b, c)
		}
		xs = unique(xs)
		xi := discreteMap(xs, 0)
		ys = unique(ys)
		yi := discrete(ys, 0)

		lx, ly := len(xi), len(yi)
		glr := make([][]int, lx)
		gdu := make([][]int, lx)
		vis := make([][]bool, lx)
		for i := range glr {
			glr[i] = make([]int, ly)
			gdu[i] = make([]int, ly)
			vis[i] = make([]bool, ly)
		}
		for _, p := range lr {
			glr[xi[p.a]][yi[p.c]]++
			glr[xi[p.b]][yi[p.c]]--
		}
		for _, p := range du {
			gdu[xi[p.a]][yi[p.b]]++
			gdu[xi[p.a]][yi[p.c]]--
		}
		for i := 1; i < lx-1; i++ {
			for j := 1; j < ly-1; j++ {
				glr[i][j] += glr[i-1][j]
				gdu[i][j] += gdu[i][j-1]
			}
		}

		type pair struct{ x, y int }
		q := []pair{{xi[0], yi[0]}}
		for len(q) > 0 {
			p := q[0]
			q = q[1:]
			x, y := p.x, p.y
			if x == 0 || x == lx-1 || y == 0 || y == ly-1 {
				return -1
			} // 无穷大
			if !vis[x][y] {
				vis[x][y] = true
				ans += (xs[x+1] - xs[x]) * (ys[y+1] - ys[y])
				if glr[x][y] == 0 {
					q = append(q, pair{x, y - 1})
				}
				if glr[x][y+1] == 0 {
					q = append(q, pair{x, y + 1})
				}
				if gdu[x][y] == 0 {
					q = append(q, pair{x - 1, y})
				}
				if gdu[x+1][y] == 0 {
					q = append(q, pair{x + 1, y})
				}
			}
		}
		return
	}

	// 括号拼接
	// 代码来源 https://codeforces.ml/gym/101341/problem/A
	// 类似题目 https://atcoder.jp/contests/abc167/tasks/abc167_f
	//         https://codeforces.ml/problemset/problem/1203/F1
	concatBrackets := func(ss [][]byte) (ids []int) {
		type pair struct{ x, y, i int }

		d := 0
		var ls, rs []pair
		for i, s := range ss {
			l, r := 0, 0
			for _, b := range s {
				if b == '(' {
					l++
				} else if l > 0 {
					l--
				} else {
					r++
				}
			}
			if r < l {
				ls = append(ls, pair{r, l, i})
			} else {
				rs = append(rs, pair{l, r, i})
			}
			d += l - r
		}

		sort.Slice(ls, func(i, j int) bool { return ls[i].x < ls[j].x })
		sort.Slice(rs, func(i, j int) bool { return rs[i].x < rs[j].x })
		f := func(ps []pair) []int {
			_ids := []int{}
			s := 0
			for _, p := range ps {
				if s < p.x {
					return nil
				}
				s += p.y - p.x
				_ids = append(_ids, p.i)
			}
			return _ids
		}
		idsL := f(ls)
		idsR := f(rs)
		if d != 0 || idsL == nil || idsR == nil {
			return
		}
		for _, id := range idsL {
			ids = append(ids, id)
		}
		for i := len(idsR) - 1; i >= 0; i-- {
			ids = append(ids, idsR[i])
		}
		return
	}

	_ = []interface{}{
		pow2, pow10, dir4, dir4C, dir4c, dir4R, dir8, orderP3, factorial,
		min, mins, max, maxs, ternaryI, ternaryS, toInts, xor, zip, zipI, getCol, minString, removeLeadingZero,
		abs, absAll, pow, calcFactorial, toAnyBase, digits, initSum2D, querySum2D, mergeMap,
		copyMat, hash01Mat, sort3, smallK, reverse, reverseSelf, equals,
		merge, splitDifferenceAndIntersection, isSubset, isDisjoint,
		unique, uniqueInPlace, discrete, discreteMap, indexMap, allSame, complement, quickSelect, contains, containsAll,
		getCycle, maxSubArraySum, maxSubArrayAbsSum, sweepLine,
		maxValueStepToUpper,
		discrete2D,
		concatBrackets,
	}
}

// https://cp-algorithms.com/sequences/rmq.html
func rmqCollection() {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	// 预处理 log
	logInit := func() {
		const mx int = 1e6
		log := make([]int, mx+1)
		for i := 2; i <= mx; i++ {
			log[i] = log[i>>1] + 1
		}
	}

	// Sparse Table
	// st[i][j] 对应的区间是 [i, i+1<<j)
	// https://oi-wiki.org/ds/sparse-table/
	// https://codeforces.ml/blog/entry/66643
	// 模板中的核心函数 max 可以换成其他具有区间合并性质的函数（允许区间重叠），如 gcd 等
	// 模板题 https://www.luogu.com.cn/problem/P3865
	// 题目推荐 https://cp-algorithms.com/data_structures/sparse-table.html#toc-tgt-5
	const mx = 17 // 131072, 262144, 524288, 1048576
	var st [][mx]int
	stInit := func(a []int) {
		n := len(a)
		st = make([][mx]int, n)
		for i, v := range a {
			st[i][0] = v
		}
		for j := 1; 1<<j <= n; j++ {
			for i := 0; i+1<<j <= n; i++ {
				st[i][j] = max(st[i][j-1], st[i+1<<(j-1)][j-1])
			}
		}
	}
	// [l,r) 注意 l r 是从 0 开始算的
	stQuery := func(l, r int) int { k := bits.Len(uint(r-l)) - 1; return max(st[l][k], st[r-1<<k][k]) }

	// Sparse Table 下标版本，查询返回的是区间最值的下标
	{
		type pair struct{ v, i int }
		const mx = 17
		var st [][mx]pair
		stInit := func(a []int) {
			n := len(a)
			st = make([][mx]pair, n)
			for i, v := range a {
				st[i][0] = pair{v, i}
			}
			for j := 1; 1<<j <= n; j++ {
				for i := 0; i+1<<j <= n; i++ {
					if a, b := st[i][j-1], st[i+1<<(j-1)][j-1]; a.v <= b.v { // 最小值，相等时下标取左侧
						st[i][j] = a
					} else {
						st[i][j] = b
					}
				}
			}
		}
		stQuery := func(l, r int) int { // [l,r) 注意 l r 是从 0 开始算的
			k := bits.Len(uint(r-l)) - 1
			a, b := st[l][k], st[r-1<<k][k]
			if a.v <= b.v { // 最小值，相等时下标取左侧
				return a.i
			}
			return b.i
		}
		_, _ = stInit, stQuery
	}

	// 分块 Sqrt Decomposition
	// https://oi-wiki.org/ds/decompose/
	// https://oi-wiki.org/ds/block-array/
	// 题目推荐 https://cp-algorithms.com/data_structures/sqrt_decomposition.html#toc-tgt-8
	// TODO: 台湾的《根號算法》https://www.csie.ntu.edu.tw/~sprout/algo2018/ppt_pdf/root_methods.pdf
	type block struct {
		l, r           int // [l,r]
		origin, sorted []int
		//lazyAdd int
	}
	var blocks []block
	sqrtInit := func(a []int) {
		n := len(a)
		blockSize := int(math.Sqrt(float64(n)))
		//blockSize := int(math.Sqrt(float64(n) * math.Log2(float64(n+1))))
		blockNum := (n-1)/blockSize + 1
		blocks = make([]block, blockNum)
		for i, v := range a {
			j := i / blockSize
			if i%blockSize == 0 {
				blocks[j] = block{l: i, origin: make([]int, 0, blockSize)}
			}
			blocks[j].origin = append(blocks[j].origin, v)
		}
		for i := range blocks {
			b := &blocks[i]
			b.r = b.l + len(b.origin) - 1
			b.sorted = make([]int, len(b.origin))
			copy(b.sorted, b.origin)
			sort.Ints(b.sorted)
		}
	}
	sqrtOp := func(l, r int) { // [l,r], starts at 0
		for i := range blocks {
			b := &blocks[i]
			if b.r < l {
				continue
			}
			if b.l > r {
				break
			}
			if l <= b.l && b.r <= r {
				// do op on full block
			} else {
				// do op on part block
				bl := max(b.l, l)
				br := min(b.r, r)
				for j := bl - b.l; j <= br-b.l; j++ {
					// do b.origin[j]...
				}
			}
		}
	}

	_ = []interface{}{
		logInit,
		stInit, stQuery,
		sqrtInit, sqrtOp,
	}
}

/* 平方根算法：组合两种算法从而降低复杂度 O(n^2) -> O(n√n)
参考 Competitive Programmer’s Handbook Ch.27

有 n 个对象，每个对象有一个「关于其他对象的统计量」ci（一个数、一个集合的元素个数，等等）
为方便起见，假设 ∑ci 的数量级和 n 一样，下面用 n 表示 ∑ci
当 ci > √n 时，这样的对象不超过 √n 个，暴力枚举这些对象之间的关系（或者，该对象与其他所有对象的关系），时间复杂度为 O(n) 或 O(n√n)。此乃算法一
当 ci ≤ √n 时，这样的对象有 O(n) 个，由于统计量 ci 很小，暴力枚举当前对象的统计量，时间复杂度为 O(n√n)。此乃算法二
这样，以 √n 为界，我们将所有对象划分成了两组，并用两个不同的算法处理
这两种算法是看待同一个问题的两种不同方式，通过恰当地组合这两个算法，复杂度由 O(n^2) 降至 O(n√n)
注意：**枚举时要做到不重不漏**

另一种题型是注意到 n 的整数分拆中，不同数字的个数至多有 O(√n) 种

好题 https://leetcode-cn.com/problems/you-le-yuan-de-you-lan-ji-hua/
*/

// 莫队算法：对询问分块
// 分块，将左端点分配在一个较小的范围，并且按照右端点从小到大排序，
// 这样对于每一块，指针移动的次数为 O(√n*√n+n) = O(n)，从而整体复杂度为 O(n√n)
// 此外，记录的是 [l,r)，这样能简化处理查询结果的代码
// https://oi-wiki.org/misc/mo-algo/
// 模板题 https://www.luogu.com.cn/problem/P1494
// 题目推荐 https://cp-algorithms.com/data_structures/sqrt_decomposition.html#toc-tgt-8
func moAlgorithm() {
	mo := func(in io.Reader, a []int, q int) []int {
		n := len(a)
		type query struct{ blockIdx, l, r, idx int }
		qs := make([]query, q)
		blockSize := int(math.Round(math.Sqrt(float64(n))))
		for i := range qs {
			var l, r int
			Fscan(in, &l, &r)
			qs[i] = query{l / blockSize, l, r + 1, i}
		}
		sort.Slice(qs, func(i, j int) bool {
			qi, qj := qs[i], qs[j]
			if qi.blockIdx != qj.blockIdx {
				return qi.blockIdx < qj.blockIdx
			}
			// 奇偶化排序
			if qi.blockIdx&1 == 0 {
				return qi.r < qj.r
			}
			return qi.r > qj.r
		})

		cnt := 0
		l, r := 1, 1 // 区间从 1 开始，方便 debug
		update := func(idx, delta int) {
			// NOTE: 有些题目在 delta 为 1 和 -1 时逻辑的顺序是严格对称的
			// v := a[idx-1]
			// ...
			if delta == 1 {
				cnt++
			} else {
				cnt--
			}
		}
		getAns := func(q query) int {
			// 提醒：q.r 是加一后的，计算时需要注意
			// sz := q.r - q.l
			// ...
			return cnt
		}
		ans := make([]int, q)
		for _, q := range qs {
			// prepare
			// NOTE: 有些题目需要维护差分值，由于 [l,r] 的差分是 s(r)-s(l-1)，此时 update 传入的应为 l-1
			for ; r < q.r; r++ {
				update(r, 1)
			}
			for ; l < q.l; l++ {
				update(l, -1)
			}
			for l > q.l {
				l--
				update(l, 1)
			}
			for r > q.r {
				r--
				update(r, -1)
			}
			ans[q.idx] = getAns(q)
		}
		return ans
	}

	// TODO: 带修改的莫队
	// https://www.luogu.com.cn/blog/deco/qian-tan-ji-chu-gen-hao-suan-fa-fen-kuai

	// TODO: 树上莫队

	_ = mo
}

func monotoneCollection() {
	// 推荐 https://cp-algorithms.com/data_structures/stack_queue_modification.html

	// 单调栈
	// 举例：返回每个元素两侧严格大于它的元素位置（不存在则为 -1 或 n）
	// 如何理解：把数组想象成一列山峰，站在 a[i] 的山顶仰望两侧的山峰，是看不到高山背后的矮山的，只能看到一座座更高的山峰
	//          这就启发我们引入一个底大顶小的单调栈，入栈时不断比较栈顶元素直到找到一个比当前元素大的
	// 技巧：事先压入一个边界元素到栈底，这样保证循环时栈一定不会为空，从而简化逻辑
	// https://oi-wiki.org/ds/monotonous-stack/
	// 模板题 https://www.luogu.com.cn/problem/P5788
	//       https://leetcode.com/problems/next-greater-element-i/
	//       https://leetcode.com/problems/next-greater-element-ii/
	// 柱状图中最大的矩形 https://leetcode-cn.com/problems/largest-rectangle-in-histogram/
	// 与 DP 结合 https://codeforces.ml/problemset/problem/1313/C2
	monotoneStack := func(a []int) ([]int, []int) {
		n := len(a)
		const border int = -1 // 2e9
		type pair struct{ v, i int }
		posL := make([]int, n)
		stack := []pair{{border, -1}}
		for i, v := range a {
			for {
				if top := stack[len(stack)-1]; top.v < v { // 严格小于
					posL[i] = top.i //+ 1
					break
				}
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, pair{v, i})
		}
		posR := make([]int, n)
		stack = []pair{{border, n}}
		for i := n - 1; i >= 0; i-- {
			v := a[i]
			for {
				if top := stack[len(stack)-1]; top.v < v { // 严格小于
					posR[i] = top.i //- 1
					break
				}
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, pair{v, i})
		}

		return posL, posR
	}

	/* 单调队列
	需要不断维护队列的单调性，即保证队列(指向)元素从大到小或从小到大
	为简单起见，这里用数组+双下标模拟双端队列
	为保证有足够空间，队列初始大小应和原数组相同
	队列存储的是元素的下标
	l == r 表示队列为空
	l < r 表示队列指向元素为 a[idQ[l]], ..., a[idQ[r-1]]（注意：这不等同于考察的区间就是 [idQ[l], idQ[r-1]]，但至少包含这一区间）
	一般的操作流程是「弹右-插右-更新答案-弹左」：
	    「弹右-插右」在维护队列单调性的同时，向右扩大了考察的区间范围
			for ; l < r && a[idQ[r-1]] <= v; r-- { // <= 为从大到小的单调队列
			}
			idQ[r] = i
			r++
	    「(检查是否满足条件)-更新答案-(弹左)」在更新答案之后，若队首在下一个循环中无用，则弹出
			具体写法随问题不同而不同，参见下面的例子
	https://oi-wiki.org/ds/monotonous-queue/
	*/

	// 单调队列模板题 - 固定区间大小的区间最值（滑动窗口）
	// https://www.luogu.com.cn/problem/P1886
	monotoneQueue := func(a []int, fixedSize int) ([]int, []int) {
		n := len(a)
		mins := make([]int, n) // mins[i] 表示 min{a[i],...,a[i+fixedSize-1]}
		idQ := make([]int, n)
		l, r := 0, 0
		for i, v := range a {
			for ; l < r && a[idQ[r-1]] >= v; r-- { // >= 意味着相等的元素取靠右的，若改成 > 表示相等的元素取靠左的
			}
			idQ[r] = i
			r++
			if i+1 >= fixedSize {
				mins[i+1-fixedSize] = a[idQ[l]]
				if idQ[l] == i+1-fixedSize {
					l++
				}
			}
		}
		maxs := make([]int, n)
		idQ = make([]int, n)
		l, r = 0, 0
		for i, v := range a {
			for ; l < r && a[idQ[r-1]] <= v; r-- {
			}
			idQ[r] = i
			r++
			if i+1 >= fixedSize {
				maxs[i+1-fixedSize] = a[idQ[l]]
				if idQ[l] == i+1-fixedSize {
					l++
				}
			}
		}
		return mins, maxs
	}

	// 单调队列应用 - 和至少为 k 的最短子数组长度
	// https://leetcode-cn.com/problems/shortest-subarray-with-sum-at-least-k/
	shortestSubarray := func(a []int, k int) int {
		n := len(a)
		const inf int = 1e9
		ans := inf
		sum := make([]int, n+1)
		for i, v := range a {
			sum[i+1] = sum[i] + v
		}
		idQ := make([]int, n+1)
		l, r := 0, 0
		for i, s := range sum {
			for ; l < r && sum[idQ[r-1]] >= s; r-- { // 贪心：相等时也弹出右侧
			}
			idQ[r] = i
			r++
			for ; l < r && s-sum[idQ[l]] >= k; l++ { // 不断弹出左侧直到队列为空或不满足要求
				// 满足要求，更新答案
				if i-idQ[l] < ans {
					ans = i - idQ[l]
				}
			}
		}
		if ans == inf {
			ans = -1
		}
		return ans
	}

	// https://codeforces.ml/problemset/problem/1237/D
	cf1237d := func(a []int, n int) (_ans []int) {
		a = append(append(a, a...), a...)
		idQ := make([]int, 3*n)
		l, r := 0, 0
		for i, j := 0, 0; i < n; i++ {
			// 维护的是从大到小的单调队列，即，队首为区间最值
			// 检查当前元素与队首(区间最值)的关系是否满足题目要求，满足则弹右插右
			for ; j < 3*n && (l == r || 2*a[j] >= a[idQ[l]]); j++ {
				for ; l < r && a[idQ[r-1]] <= a[j]; r-- {
				}
				idQ[r] = j
				r++
			}
			// 更新答案
			ans := j - i
			if ans > 2*n {
				ans = -1
			}
			_ans = append(_ans, ans)
			// 准备：若下一个循环中的队首不在考察区间内，则弹左
			if l < r && idQ[l] == i {
				l++
			}
		}

		{
			idQ := []int{}
			for i, j := 0, 0; i < n; i++ {
				for ; j < 3*n && (len(idQ) == 0 || 2*a[j] >= a[idQ[0]]); j++ {
					for ; len(idQ) > 0 && a[idQ[len(idQ)-1]] <= a[j]; r-- {
					}
					idQ = append(idQ, j)
				}
				ans := j - i
				if ans > 2*n {
					ans = -1
				}
				_ans = append(_ans, ans)
				if len(idQ) > 0 && idQ[0] == i {
					idQ = idQ[1:]
				}
			}
		}
		return
	}

	_ = []interface{}{monotoneStack, monotoneQueue, shortestSubarray, cf1237d}
}

func loopCollection() {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	// 枚举 {0,1,...,n-1} 的全部子集
	loopSet := func(arr []int) {
		n := len(arr)
		//outer:
		for sub := 0; sub < 1<<n; sub++ { // sub repr a subset which elements are in range [0,n)
			// do(sub)
			for i := 0; i < n; i++ {
				if sub>>i&1 == 1 { // choose i in sub
					_ = arr[i]
					// do(arr[i]) or continue outer
				}
			}
		}
	}

	// 枚举 subset 的全部子集
	// 作为结束条件，处理完 0 之后，会有 -1&subset == subset
	loopSubset := func(n, subset int) {
		sub := subset
		for ok := true; ok; ok = sub != subset {
			// do(sub)
			sub = (sub - 1) & subset
		}
	}

	// 枚举大小为 n 的集合的大小为 k 的子集（按字典序）
	// 参考《挑战程序设计竞赛》p.156-158
	// 比如在 n 个数中求满足某种性质的最大子集，则可以从 n 开始倒着枚举子集大小，直到找到一个符合性质的子集
	// 例题（TS1）https://codingcompetitions.withgoogle.com/codejam/round/0000000000007706/0000000000045875
	loopSubsetK := func(arr []int, k int) {
		n := len(arr)
		for sub := 1<<k - 1; sub < 1<<n; {
			// do(arr, sub) ...
			x := sub & -sub
			y := sub + x
			sub = sub&^y/x>>1 | y
		}
	}

	/*
		遍历以 (centerI, centerJ) 为中心的欧几里得距离为 dis 范围内的格点
		例如 dis=2 时：
		  #
		 # #
		# @ #
		 # #
		  #
	*/
	type pair struct{ x, y int }
	dir4 := [...]pair{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} // 上下左右
	searchDir4 := func(maxI, maxJ, centerI, centerJ, dis int) {
		for i, d := range dir4 {
			d2 := dir4[(i+1)%4]
			dx := d2.x - d.x
			dy := d2.y - d.y
			x := centerI + d.x*dis
			y := centerJ + d.y*dis
			for _i := 0; _i < dis; _i++ {
				if x >= 0 && x < maxI && y >= 0 && y < maxJ {
					// do
				}
				x += dx
				y += dy
			}
		}
	}

	/*
		#####
		#   #
		# @ #
		#   #
		#####
	*/
	searchDir4R := func(maxI, maxJ, centerI, centerJ, dis int) {
		// 上下
		for _, x := range [...]int{centerI - dis, centerI + dis} {
			if x >= 0 && x < maxI {
				for y := max(centerJ-dis, 0); y < min(centerJ+dis, maxJ); y++ {
					// do
				}
			}
		}
		// 左右
		for _, y := range [...]int{centerJ - dis, centerJ + dis} {
			if y >= 0 && y < maxJ {
				for x := max(centerI-dis, 0); x < min(centerI+dis, maxI); x++ {
					// do
				}
			}
		}
	}

	loopDiagonal := func(mat [][]int) {
		n, m := len(mat), len(mat[0])
		for j := 0; j < m; j++ {
			for i := 0; i < n; i++ {
				if i > j {
					break
				}
				_ = mat[i][j-i]
			}
		}
		for i := 1; i < n; i++ {
			for j := m - 1; j >= 0; j-- {
				if i+m-1-j >= n {
					break
				}
				_ = mat[i+m-1-j][j]
			}
		}
	}

	loopDiagonal2 := func(n int) {
		for sum := 0; sum < 2*n-1; sum++ {
			for x := 0; x <= sum; x++ {
				y := sum - x
				if x >= n || y >= n {
					continue
				}
				Println(x, y)
			}
			Println()
		}
	}

	_ = []interface{}{
		loopSet, loopSubset, loopSubsetK,
		searchDir4, searchDir4R, loopDiagonal, loopDiagonal2,
	}
}
