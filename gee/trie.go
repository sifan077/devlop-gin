package gee

import "strings"

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	// 遍历字节点
	for _, child := range n.children {
		// 如果路由中的一部分相等，或者是通配即可返回子节点
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	// 声明一个节点切片
	nodes := make([]*node, 0)
	// 遍历当前节点的子节点
	for _, child := range n.children {
		// 如果路由一部分能够匹配或者是通配符，加入切片
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	// 如果高度和子节点数组的长度相等，则当前的节点的代匹配路由等于传入的路由
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	// 获取所有路由的第height个路由
	part := parts[height]
	// 获取当前节点的第一个匹配的子节点
	child := n.matchChild(part)
	// 如果当前节点没有匹配的子节点，则创建一个节点，添加到子节点的切片中
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	// 递归插入
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	// 如果路由切片到结尾或者匹配的路由以通配符开头，如果当前待匹配路由不为空字符串则返回当前节点
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	// 获取所有路由的第height个路由
	part := parts[height]
	// 获取当前节点的第一个匹配的子节点
	children := n.matchChildren(part)
	// 遍历当前节点的所有子节点,并且在子节点中继续递归search寻找节点,如果找到返回目标节点
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
