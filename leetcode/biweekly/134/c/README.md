本题是环形数组，请先完成普通数组的版本：[3101. 交替子数组计数](https://leetcode.cn/problems/count-alternating-subarrays/)（[我的题解](https://leetcode.cn/problems/count-alternating-subarrays/solution/jian-ji-xie-fa-pythonjavacgo-by-endlessc-tcc9/)）

把数组复制一份拼接起来，和 3101 题一样，遍历数组的同时，维护以 $i$ 为右端点的交替子数组的长度 $\textit{cnt}$。

如果 $i\ge n$ 且 $\textit{cnt}\ge k$，那么 $i$ 就是一个长为 $k$ 的交替子数组的右端点，答案加一。注意这里要判断 $i\ge n$，从而避免重复统计。

代码实现时，不需要复制数组，而是用 $i\bmod n$ 的方式取到对应的值。

具体请看 [视频讲解](https://www.bilibili.com/video/BV1Yz421q7dD/) 第三题，欢迎点赞关注！

```py [sol-Python3]
class Solution:
    def numberOfAlternatingGroups(self, colors: List[int], k: int) -> int:
        n = len(colors)
        ans = cnt = 0
        for i in range(n * 2):
            if i > 0 and colors[i % n] == colors[(i - 1) % n]:
                cnt = 0
            cnt += 1
            if i >= n and cnt >= k:
                ans += 1
        return ans
```

```java [sol-Java]
public class Solution {
    public int numberOfAlternatingGroups(int[] colors, int k) {
        int n = colors.length;
        int ans = 0;
        int cnt = 0;
        for (int i = 0; i < n * 2; i++) {
            if (i > 0 && colors[i % n] == colors[(i - 1) % n]) {
                cnt = 0;
            }
            cnt++;
            if (i >= n && cnt >= k) {
                ans++;
            }
        }
        return ans;
    }
}
```

```cpp [sol-C++]
class Solution {
public:
    int numberOfAlternatingGroups(vector<int>& colors, int k) {
        int n = colors.size();
        int ans = 0, cnt = 0;
        for (int i = 0; i < n * 2; i++) {
            if (i > 0 && colors[i % n] == colors[(i - 1) % n]) {
                cnt = 0;
            }
            cnt++;
            ans += i >= n && cnt >= k;
        }
        return ans;
    }
};
```

```go [sol-Go]
func numberOfAlternatingGroups(colors []int, k int) (ans int) {
	n := len(colors)
	cnt := 0
	for i := 0; i < n*2; i++ {
		if i > 0 && colors[i%n] == colors[(i-1)%n] {
			cnt = 0
		}
		cnt++
		if i >= n && cnt >= k {
			ans++
		}
	}
	return
}
```

#### 复杂度分析

- 时间复杂度：$\mathcal{O}(n)$，其中 $n$ 是 $\textit{colors}$ 的长度。
- 空间复杂度：$\mathcal{O}(1)$。

## 相似题目

- [503. 下一个更大元素 II](https://leetcode.cn/problems/next-greater-element-ii/)

## 分类题单

以下题单没有特定的顺序，可以按照个人喜好刷题。

1. [滑动窗口（定长/不定长/多指针）](https://leetcode.cn/circle/discuss/0viNMK/)
2. [二分算法（二分答案/最小化最大值/最大化最小值/第K小）](https://leetcode.cn/circle/discuss/SqopEo/)
3. [单调栈（基础/矩形面积/贡献法/最小字典序）](https://leetcode.cn/circle/discuss/9oZFK9/)
4. [网格图（DFS/BFS/综合应用）](https://leetcode.cn/circle/discuss/YiXPXW/)
5. [位运算（基础/性质/拆位/试填/恒等式/思维）](https://leetcode.cn/circle/discuss/dHn9Vk/)
6. [图论算法（DFS/BFS/拓扑排序/最短路/最小生成树/二分图/基环树/欧拉路径）](https://leetcode.cn/circle/discuss/01LUak/)
7. [动态规划（入门/背包/状态机/划分/区间/状压/数位/数据结构优化/树形/博弈/概率期望）](https://leetcode.cn/circle/discuss/tXLS3i/)
8. [常用数据结构（前缀和/差分/栈/队列/堆/字典树/并查集/树状数组/线段树）](https://leetcode.cn/circle/discuss/mOr1u6/)
9. [数学算法（数论/组合/概率期望/博弈/计算几何/随机算法）](https://leetcode.cn/circle/discuss/IYT3ss/)
10. [贪心算法（基本贪心策略/反悔/区间/字典序/数学/思维/脑筋急转弯/构造）](https://leetcode.cn/circle/discuss/g6KTKL/)

[我的题解精选（已分类）](https://github.com/EndlessCheng/codeforces-go/blob/master/leetcode/SOLUTIONS.md)
