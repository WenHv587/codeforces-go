思路和 [70. 爬楼梯](https://leetcode-cn.com/problems/climbing-stairs/) 类似。

把相同字符分为一组，每组内只有一种字符，考虑如下 DP：

对于字符不为 $\texttt{7}$ 或 $\texttt{9}$ 的情况，定义 $f[i]$ 表示长为 $i$ 的只有一种字符的字符串对应的文字信息种类数，我们可以将末尾的 $1$ 个、$2$ 个或 $3$ 个字符单独视作一个字母，那么有转移方程

$$
f[i] = f[i-1]+f[i-2]+f[i-3]
$$

对于字符为 $\texttt{7}$ 或 $\texttt{9}$ 的情况，定义 $g[i]$ 表示长为 $i$ 的只有一种字符的字符串对应的文字信息种类数，可以得到类似的转移方程

$$
g[i] = g[i-1]+g[i-2]+g[i-3]+g[i-4]
$$

这样能算出每组字符串的文字信息种类数。

由于不同组之间互不影响，根据乘法原理，把不同组的文字信息种类数相乘，得到答案。

```python [sol-Python3]
MOD = 1_000_000_007
f = [1, 1, 2, 4]
g = [1, 1, 2, 4]
for _ in range(10 ** 5 - 3):  # 预处理所有长度的结果
    f.append((f[-1] + f[-2] + f[-3]) % MOD)
    g.append((g[-1] + g[-2] + g[-3] + g[-4]) % MOD)

class Solution:
    def countTexts(self, pressedKeys: str) -> int:
        ans = 1
        for ch, s in groupby(pressedKeys):
            m = len(list(s))
            ans = ans * (g[m] if ch in "79" else f[m]) % MOD
        return ans
```

```java [sol-Java]
class Solution {
    private static final int MOD = 1_000_000_007;
    private static final int MX = 100_001;
    private static final long[] f = new long[MX];
    private static final long[] g = new long[MX];

    static {
        f[0] = g[0] = 1;
        f[1] = g[1] = 1;
        f[2] = g[2] = 2;
        f[3] = g[3] = 4;
        for (int i = 4; i < MX; i++) {
            f[i] = (f[i - 1] + f[i - 2] + f[i - 3]) % MOD;
            g[i] = (g[i - 1] + g[i - 2] + g[i - 3] + g[i - 4]) % MOD;
        }
    }

    public int countTexts(String s) {
        long ans = 1;
        int cnt = 0;
        for (int i = 0; i < s.length(); i++) {
            char c = s.charAt(i);
            cnt++;
            if (i == s.length() - 1 || c != s.charAt(i + 1)) { // 找到一个完整的组
                ans = ans * (c != '7' && c != '9' ? f[cnt] : g[cnt]) % MOD;
                cnt = 0;
            }
        }
        return (int) ans;
    }
}
```

```cpp [sol-C++]
const int MOD = 1'000'000'007, MX = 100'001;
long long f[MX], g[MX];

int init = []() {
    f[0] = g[0] = 1;
    f[1] = g[1] = 1;
    f[2] = g[2] = 2;
    f[3] = g[3] = 4;
    for (int i = 4; i < MX; ++i) {
        f[i] = (f[i - 1] + f[i - 2] + f[i - 3]) % MOD;
        g[i] = (g[i - 1] + g[i - 2] + g[i - 3] + g[i - 4]) % MOD;
    }
    return 0;
}();

class Solution {
public:
    int countTexts(string& s) {
        long long ans = 1;
        int cnt = 0;
        for (int i = 0; i < s.length(); i++) {
            char c = s[i];
            cnt++;
            if (i == s.length() - 1 || c != s[i + 1]) { // 找到一个完整的组
                ans = ans * (c != '7' && c != '9' ? f[cnt] : g[cnt]) % MOD;
                cnt = 0;
            }
        }
        return ans;
    }
};
```

```go [sol-Go]
const mod = 1_000_000_007
const mx = 100_001

var f = [mx]int{1, 1, 2, 4}
var g = f

func init() {
	for i := 4; i < mx; i++ {
		f[i] = (f[i-1] + f[i-2] + f[i-3]) % mod
		g[i] = (g[i-1] + g[i-2] + g[i-3] + g[i-4]) % mod
	}
}

func countTexts(s string) int {
	ans, cnt := 1, 0
	for i, c := range s {
		cnt++
		if i == len(s)-1 || s[i+1] != byte(c) { // 找到一个完整的组
			if c != '7' && c != '9' {
				ans = ans * f[cnt] % mod
			} else {
				ans = ans * g[cnt] % mod
			}
			cnt = 0
		}
	}
	return ans
}
```

#### 复杂度分析

- 时间复杂度：$\mathcal{O}(n)$，其中 $n$ 是 $\textit{pressedKeys}$ 的长度。忽略预处理的时间和空间。
- 空间复杂度：$\mathcal{O}(1)$。

## 分类题单

[如何科学刷题？](https://leetcode.cn/circle/discuss/RvFUtj/)

1. [滑动窗口与双指针（定长/不定长/单序列/双序列/三指针）](https://leetcode.cn/circle/discuss/0viNMK/)
2. [二分算法（二分答案/最小化最大值/最大化最小值/第K小）](https://leetcode.cn/circle/discuss/SqopEo/)
3. [单调栈（基础/矩形面积/贡献法/最小字典序）](https://leetcode.cn/circle/discuss/9oZFK9/)
4. [网格图（DFS/BFS/综合应用）](https://leetcode.cn/circle/discuss/YiXPXW/)
5. [位运算（基础/性质/拆位/试填/恒等式/思维）](https://leetcode.cn/circle/discuss/dHn9Vk/)
6. [图论算法（DFS/BFS/拓扑排序/最短路/最小生成树/二分图/基环树/欧拉路径）](https://leetcode.cn/circle/discuss/01LUak/)
7. [动态规划（入门/背包/状态机/划分/区间/状压/数位/数据结构优化/树形/博弈/概率期望）](https://leetcode.cn/circle/discuss/tXLS3i/)
8. [常用数据结构（前缀和/差分/栈/队列/堆/字典树/并查集/树状数组/线段树）](https://leetcode.cn/circle/discuss/mOr1u6/)
9. [数学算法（数论/组合/概率期望/博弈/计算几何/随机算法）](https://leetcode.cn/circle/discuss/IYT3ss/)
10. [贪心与思维（基本贪心策略/反悔/区间/字典序/数学/思维/脑筋急转弯/构造）](https://leetcode.cn/circle/discuss/g6KTKL/)
11. [链表、二叉树与一般树（前后指针/快慢指针/DFS/BFS/直径/LCA）](https://leetcode.cn/circle/discuss/K0n2gO/)

[我的题解精选（已分类）](https://github.com/EndlessCheng/codeforces-go/blob/master/leetcode/SOLUTIONS.md)

欢迎关注 [B站@灵茶山艾府](https://space.bilibili.com/206214)
