package algorithm

import (
	"backend/search"
	"fmt"
)

// PrintCraftingPath prints the crafting steps for a given path
func PrintCraftingPath(path []*search.ElementNode) {
	if len(path) == 0 {
		fmt.Println("No crafting path found.")
		return
	}

	for i := len(path) - 1; i >= 1; i -= 2 {
		if i-1 < 0 || i-2 < 0 {
			break
		}
		if path[i-1].Name == "" && path[i-2].Name == "" {
			continue
		} else {
			fmt.Printf("%s => %s + %s\n", path[i].Name, path[i-1].Name, path[i-2].Name)
		}
	}
}

func DFS(target *search.ElementNode, maxPaths int) [][]*search.ElementNode {
	type StackFrame struct {
		Path []*search.ElementNode
	}

	var stack []StackFrame
	stack = append(stack, StackFrame{Path: []*search.ElementNode{target}})

	var results [][]*search.ElementNode
	visited := make(map[string]bool)

	containsNode := func(path []*search.ElementNode, node *search.ElementNode) bool {
		for _, n := range path {
			if n.ID == node.ID {
				return true
			}
		}
		return false
	}

	for len(stack) > 0 && len(results) < maxPaths {
		// Pop from stack
		lastIdx := len(stack) - 1
		curr := stack[lastIdx]
		stack = stack[:lastIdx]

		currentNode := curr.Path[0]

		if len(currentNode.Recipes) == 0 {
			results = append(results, curr.Path)
			continue
		}

		for _, recipe := range currentNode.Recipes {
			if len(recipe) != 2 {
				continue
			}
			a, b := recipe[0], recipe[1]

			key := fmt.Sprintf("%d+%d->%d", a.ID, b.ID, currentNode.ID)
			if visited[key] {
				continue
			}
			visited[key] = true

			if containsNode(curr.Path, a) || containsNode(curr.Path, b) {
				continue
			}

			newPathA := append([]*search.ElementNode{a}, curr.Path...)
			stack = append(stack, StackFrame{Path: newPathA})
			
			newPathB := append([]*search.ElementNode{b}, curr.Path...)
			stack = append(stack, StackFrame{Path: newPathB})
		}
	}

	return results
}