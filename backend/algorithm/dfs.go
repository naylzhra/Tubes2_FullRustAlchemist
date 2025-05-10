package algorithm

import (
	"backend/search"
	"fmt"
	"slices"
	"sync"
)

type SearchStatus struct {
	mu             sync.Mutex
	routines       map[string]chan struct{}
	foundPathCount map[string]int
	prevElems		 map[string][]*search.ElementNode
	complete       map[string]bool
	target         string
	maxPaths       int
	resultGraph    *search.RecipeGraph
}

type ResultTree struct {
	mu sync.Mutex
	elements map[string][]*search.ElementNode
}

// PrintCraftingPath prints the crafting steps for a given path
func PrintCraftingPath(result *ResultTree) string {
	if len(path) == 0 {
		fmt.Println("No crafting path found.")
		return
	}
	fmt.Printf("Crafting path for %s:\n", path[0].Name)

	for i := len(path) - 1; i >= 1; i -= 2 {''
		if i-1 < 0 || i-2 < 0 {
			break
		}
		if path[i-1].Name == "" && path[i-2].Name == "" {
			continue
		}

		// normal triple
		product, ing1, ing2 := path[i], path[i+1], path[i+2]
		sfx := func(n *search.ElementNode) string {
			if isBaseElement(n) { return " (base)" }
			return ""
		}
		fmt.Printf("%s <= %s%s + %s%s\n",
			product.Name, ing1.Name, sfx(ing1), ing2.Name, sfx(ing2))
		i += 3
	}
}

func DFS(target *search.ElementNode, maxPaths int) []string {
	// status := &SearchStatus{
	// 	mu:             sync.Mutex{},
	// 	routines:       make(map[string]chan struct{}),
	// 	foundPathCount: make(map[string]int),
	// 	complete:       make(map[string]bool),
	// 	target:         target.Name,
	// 	maxPaths:       maxPaths,
	// 	resultGraph:    &search.RecipeGraph{},
	// }
	// search.ConstructElementsGraph(graph, status.resultGraph)
	// // Create channel for each possible elements
	// for _, element := range graph.Elements {
	// 	status.routines[element.Name] = make(chan struct{})
	// 	status.complete[element.Name] = false
	// }
	// // Mark base elements as compplete
	// for _, element := range graph.BaseElements {
	// 	status.foundPathCount[element.Name] = 1
	// 	status.complete[element.Name] = true
	// }

	// status.mu.Lock()
	// status.foundPathCount[target.Name] = 0
	// status.mu.Unlock()

	foundPaths := 0
	completePathSignal := make(chan int)
	result := ResultTree{sync.Mutex{}, make(map[string][]*search.ElementNode)}
	resultJSONs := []string{}
	// Start the DFS for the target element
	go findPath(target, &result, completePathSignal)

	for {
		status := <-completePathSignal
		if status == 0 {
			break
		} else if status == 1 {
			foundPaths = foundPaths + 1
			resultJSONs = append(resultJSONs, PrintCraftingPath(&result))
		}
	}

	return resultJSONs
}

func updateCompletePath(node *search.ElementNode, status *SearchStatus, direct bool) {
	
}

func findPath(target *search.ElementNode, graph *search.RecipeGraph, result *ResultTree, completePathSignal chan int, continueSearch chan struct{}, prevs []*search.ElementNode) {
	defer func() {
		completePathSignal <- 0
		close(completePathSignal)
	}()

	// Generate thread for the two parents for a recipe (if not already created)
	// Wait until the recipe is fully complete before checking another recipe of this node
	for _, recipe := range target.Recipes {
		if slices.Contains(graph.BaseElements, recipe[0]) && slices.Contains(graph.BaseElements, recipe[1]) {
			completePathSignal <- 1
			<-continueSearch
			continue // No need to generate thread for base elements
		}

		// Decide whether to search for the parents of this element
		if pathAlreadyContains(prevs, recipe[0]) || pathAlreadyContains(prevs, recipe[1]) {
			continue
		}

		// Create a new thread for each parent of a recipe
		currentPath := append([]*search.ElementNode{}, target)
		currentPath = append([]*search.ElementNode{}, prevs...)
		chan2 := make(chan int)
		continueSearch1 := make(chan struct{})
		continueSearch2 := make(chan struct{})
		status1 := -1
		status2 := -1
		
		// Let 
		go findPath(recipe[1], graph, result, chan2, continueSearch2, currentPath)
		while (status2 != 0 && status1 != 0) {
			chan1 := make(chan int)
			go findPath(recipe[0], graph, result, chan1, continueSearch1, currentPath)

			status1 = <-chan1
			status2 = <-chan2
			if ()
			<-continueSearch
		}




		// Wait until 

		<-continueSearch
		continueSearch1 <- struct{}{}
		continueSearch2 <- struct{}{}
	}
}


func pathAlreadyContains(prevs []*search.ElementNode, elem *search.ElementNode) bool {
	for _, prev := range prevs {
		if elem == prev {
			return true
		} 
	}
	return false
}
