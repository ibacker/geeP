package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  //待匹配路由
	part     string  //路由中的一部分
	children []*node //子节点
	isWild   bool    //是否模糊匹配
}

// 匹配一个子节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 匹配所有子节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	fmt.Println("trie-matchChildren ", n.children)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 递归插入路由节点信息
func (n *node) insert(pattern string, parts []string, height int) {
	// 结束递归
	// 待插入节点数等于递归次数（高度）
	// 叶子节点保存 pattern
	if len(parts) == height {
		n.pattern = pattern

		fmt.Println("insert n", n)
		return
	}

	fmt.Println("insert", parts)
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		wild := part[0] == ':' || part[0] == '*'
		child = &node{part: part, isWild: wild}
		// 将子节点添加到 children 中
		n.children = append(n.children, child)
	}
	fmt.Println("insert", child)
	child.insert(pattern, parts, height+1)
}

// 递归查询 parts：请求路径
// height：保存深度
func (n *node) search(parts []string, height int) *node {
	// 结束递归条件：遍历完所有路径 或 命中通配符路径
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		fmt.Println("trie.search-", n)
		// 叶子节点中无请求路径，返回空
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}
