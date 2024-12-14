设 $s$ 的长度为 $n$，$t$ 的长度为 $k$。

由于 $s$ 是由若干长度为 $k$ 的字符串拼接而成，所以 $k$ 一定是 $n$ 的因子。

由于 $10^5$ 以内的因子个数至多为 $128$（$83160$ 的因子个数），所以我们可以暴力枚举 $n$ 的因子作为 $k$。

然后比较所有首字母下标为 $0,k,2k,3k,\cdots,n-k$ 的长为 $k$ 的子串，所包含的字母及其个数是否一样（同位字符串）。

注意只需枚举小于 $n$ 的因子，如果这些因子都不满足要求，答案一定是 $n$（如示例 2）。

请看 [视频讲解](https://www.bilibili.com/video/BV1Nf421U7em/) 第三题，欢迎点赞关注~

```py [sol-Py]
class Solution:
    def minAnagramLength(self, s: str) -> int:
        n = len(s)
        for k in range(1, n // 2 + 1):
            if n % k:
                continue
            cnt0 = Counter(s[:k])
            for i in range(k * 2, n + 1, k):
                if Counter(s[i - k: i]) != cnt0:
                    break
            else:
                return k
        return n
```

```py [sol-Py 写法二]
class Solution:
    def minAnagramLength(self, s: str) -> int:
        n = len(s)
        for k in range(1, n // 2 + 1):
            if n % k:
                continue
            cnt0 = Counter(s[:k])
            if all(Counter(s[i - k: i]) == cnt0 for i in range(k * 2, n + 1, k)):
                return k
        return n
```

```py [sol-Py 写法三]
class Solution:
    def minAnagramLength(self, s: str) -> int:
        n = len(s)
        for k in range(1, n // 2 + 1):
            if n % k:
                continue
            t = sorted(s[:k])
            if all(sorted(s[i - k: i]) == t for i in range(k * 2, n + 1, k)):
                return k
        return n
```

```java [sol-Java]
class Solution {
    public int minAnagramLength(String S) {
        char[] s = S.toCharArray();
        int n = s.length;
        next:
        for (int k = 1; k <= n / 2; k++) {
            if (n % k > 0) {
                continue;
            }
            int[] cnt0 = new int[26];
            for (int j = 0; j < k; j++) {
                cnt0[s[j] - 'a']++;
            }
            for (int i = k * 2; i <= n; i += k) {
                int[] cnt = new int[26];
                for (int j = i - k; j < i; j++) {
                    cnt[s[j] - 'a']++;
                }
                if (!Arrays.equals(cnt, cnt0)) {
                    continue next;
                }
            }
            return k;
        }
        return n;
    }
}
```

```cpp [sol-C++]
class Solution {
public:
    int minAnagramLength(string s) {
        int n = s.length();
        for (int k = 1; k <= n / 2; k++) {
            if (n % k) {
                continue;
            }
            array<int, 26> cnt0{};
            for (int j = 0; j < k; j++) {
                cnt0[s[j] - 'a']++;
            }
            for (int i = k * 2; i <= n; i += k) {
                array<int, 26> cnt{};
                for (int j = i - k; j < i; j++) {
                    cnt[s[j] - 'a']++;
                }
                if (cnt != cnt0) {
                    goto next;
                }
            }
            return k;
            next:;
        }
        return n;
    }
};
```

```go [sol-Go]
func minAnagramLength(s string) int {
	n := len(s)
next:
	for k := 1; k <= n/2; k++ {
		if n%k > 0 {
			continue
		}
		cnt0 := [26]int{}
		for _, b := range s[:k] {
			cnt0[b-'a']++
		}
		for i := k * 2; i <= len(s); i += k {
			cnt := [26]int{}
			for _, b := range s[i-k : i] {
				cnt[b-'a']++
			}
			if cnt != cnt0 {
				continue next
			}
		}
		return k
	}
	return n
}
```

#### 复杂度分析

- 时间复杂度：$\mathcal{O}(A\cdot n)$，其中 $n$ 为 $s$ 的长度，$A$ 为 $n$ 的因子个数。
- 空间复杂度：$\mathcal{O}(|\Sigma|)$。其中 $|\Sigma|$ 为字符集合的大小，本题字符均为小写字母，所以 $|\Sigma|=26$。

## 分类题单

[如何科学刷题？](https://leetcode.cn/circle/discuss/RvFUtj/)

1. [滑动窗口与双指针（定长/不定长/单序列/双序列/三指针/分组循环）](https://leetcode.cn/circle/discuss/0viNMK/)
2. [二分算法（二分答案/最小化最大值/最大化最小值/第K小）](https://leetcode.cn/circle/discuss/SqopEo/)
3. [单调栈（基础/矩形面积/贡献法/最小字典序）](https://leetcode.cn/circle/discuss/9oZFK9/)
4. [网格图（DFS/BFS/综合应用）](https://leetcode.cn/circle/discuss/YiXPXW/)
5. [位运算（基础/性质/拆位/试填/恒等式/思维）](https://leetcode.cn/circle/discuss/dHn9Vk/)
6. [图论算法（DFS/BFS/拓扑排序/最短路/最小生成树/二分图/基环树/欧拉路径）](https://leetcode.cn/circle/discuss/01LUak/)
7. [动态规划（入门/背包/状态机/划分/区间/状压/数位/数据结构优化/树形/博弈/概率期望）](https://leetcode.cn/circle/discuss/tXLS3i/)
8. [常用数据结构（前缀和/差分/栈/队列/堆/字典树/并查集/树状数组/线段树）](https://leetcode.cn/circle/discuss/mOr1u6/)
9. [数学算法（数论/组合/概率期望/博弈/计算几何/随机算法）](https://leetcode.cn/circle/discuss/IYT3ss/)
10. [贪心与思维（基本贪心策略/反悔/区间/字典序/数学/思维/脑筋急转弯/构造）](https://leetcode.cn/circle/discuss/g6KTKL/)
11. [链表、二叉树与回溯（前后指针/快慢指针/DFS/BFS/直径/LCA/一般树）](https://leetcode.cn/circle/discuss/K0n2gO/)
12. [字符串（KMP/Z函数/Manacher/字符串哈希/AC自动机/后缀数组/子序列自动机）](https://leetcode.cn/circle/discuss/SJFwQI/)

[我的题解精选（已分类）](https://github.com/EndlessCheng/codeforces-go/blob/master/leetcode/SOLUTIONS.md)

欢迎关注 [B站@灵茶山艾府](https://space.bilibili.com/206214)
