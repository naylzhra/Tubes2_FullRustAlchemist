package algorithm

import (
	"backend/search"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sync"
)

type SearchStatus struct {
	chan0           chan int
	chan1           chan int
	continueSearch0 chan struct{}
	continueSearch1 chan struct{}
}

type ResultTree struct {
	mu             sync.Mutex
	elements       map[string][]*search.ElementNode
	remainingPaths int
}

type RecipeJSON struct {
	Element string   `json:"element"`
	Recipe  []string `json:"recipe"`
}

func DFS(target *search.ElementNode, graph *search.RecipeGraph, maxPaths int) []string {
	foundPaths := 0
	completePathSignal := make(chan int, 1)
	continueSearch := make(chan struct{}, 1)
	result := ResultTree{
		mu:             sync.Mutex{},
		elements:       make(map[string][]*search.ElementNode),
		remainingPaths: maxPaths,
	}
	resultJSONs := []string{}
	// Start the DFS for the target element
	go findPath(target, graph, &result, completePathSignal, continueSearch, []*search.ElementNode{})

	for {
		status := <-completePathSignal
		if status == 0 {
			break
		} else if status == 1 {
			foundPaths++
			result.mu.Lock()
			result.remainingPaths--
			resultJSONs = append(resultJSONs, ParseCraftingPath(&result, foundPaths))
			result.mu.Unlock()

			continueSearch <- struct{}{}
		}
	}
	close(continueSearch)

	return resultJSONs
}

func findPath(target *search.ElementNode, graph *search.RecipeGraph, result *ResultTree, completePathSignal chan int, continueSearch chan struct{}, prevs []*search.ElementNode) {
	defer func() {
		completePathSignal <- 0
		close(completePathSignal)
	}()

	// Base case: if the target is a base element, return
	if slices.Contains(graph.BaseElements, target) {
		completePathSignal <- 1
		return
	}

	fmt.Println("Searching for:", target.Name)

	// Generate thread for the two parents for a recipe (if not already created)
	// Wait until the recipe is fully complete before checking another recipe of this node
	for _, recipe := range target.Recipes {
		result.mu.Lock()
		if result.remainingPaths == 0 {
			result.mu.Unlock()
			return
		}
		result.mu.Unlock()

		if slices.Contains(graph.BaseElements, recipe[0]) && slices.Contains(graph.BaseElements, recipe[1]) {
			completePathSignal <- 1
			result.mu.Lock()
			result.elements[target.Name] = []*search.ElementNode{}
			result.elements[target.Name] = append(result.elements[target.Name], recipe[0])
			result.elements[target.Name] = append(result.elements[target.Name], recipe[1])
			result.mu.Unlock()
			<-continueSearch
			continue // No need to generate thread for base elements
		}

		if recipe[0].Tier >= target.Tier || recipe[1].Tier >= target.Tier {
			continue
		}

		// Decide whether to search for the parents of this element
		if pathAlreadyContains(prevs, recipe[0]) || pathAlreadyContains(prevs, recipe[1]) {
			continue
		}

		// Create a new thread for each parent of a recipe
		currentPath := slices.Clone(prevs)
		currentPath = append(currentPath, target)
		searchStatus := SearchStatus{
			chan0:           make(chan int, 1),
			chan1:           make(chan int, 1),
			continueSearch0: make(chan struct{}, 1),
			continueSearch1: make(chan struct{}, 1),
		}

		// Let the first parent backtrack, then the second
		// This is to enumerate all paths to recipe[0] x all paths to recipe[1]
		go findPath(recipe[0], graph, result, searchStatus.chan0, searchStatus.continueSearch0, currentPath)
		go findPath(recipe[1], graph, result, searchStatus.chan1, searchStatus.continueSearch1, currentPath)
		status0 := <-searchStatus.chan0
		status1 := <-searchStatus.chan1
		for status1 != 0 {
			result.mu.Lock()
			stop := result.remainingPaths == 0
			result.mu.Unlock()

			if stop {
				if status0 != 0 {
					searchStatus.continueSearch0 <- struct{}{}
				}
				searchStatus.continueSearch1 <- struct{}{}
				break // No more paths to find
			}

			if status0 == 0 { // First parent is done, so backtrack the second parent and restart the first parent thread
				searchStatus.continueSearch1 <- struct{}{}
				status1 = <-searchStatus.chan1

				if status1 == 0 {
					break // Second parent is done, so no other possible path to find
				}
				<-continueSearch
				searchStatus.chan0 = make(chan int)
				go findPath(recipe[0], graph, result, searchStatus.chan0, searchStatus.continueSearch0, currentPath)
				status0 = <-searchStatus.chan0
			} else if status0 == 1 && status1 == 1 {
				completePathSignal <- 1
				result.mu.Lock()
				result.elements[target.Name] = []*search.ElementNode{}
				result.elements[target.Name] = append(result.elements[target.Name], recipe[0])
				result.elements[target.Name] = append(result.elements[target.Name], recipe[1])
				result.mu.Unlock()

				<-continueSearch
				searchStatus.continueSearch0 <- struct{}{}
			}
		}

		close(searchStatus.continueSearch0)
		close(searchStatus.continueSearch1)
		// Kill this routine
	}
}

func pathAlreadyContains(prevs []*search.ElementNode, elem *search.ElementNode) bool {
	return slices.Contains(prevs, elem)
}

func ParseCraftingPath(result *ResultTree, foundPaths int) string {
	// ResultTree is locked in the caller side
	resultFile := "result" + fmt.Sprintf("%03d", foundPaths) + ".json"

	recipies := []RecipeJSON{}
	for elementName, recipes := range result.elements {
		recipeEntry := RecipeJSON{
			Element: elementName,
			Recipe:  make([]string, 0),
		}
		for _, recipe := range recipes {
			recipeEntry.Recipe = append(recipeEntry.Recipe, recipe.Name)
		}
		recipies = append(recipies, recipeEntry)
	}

	jsonData, err := json.MarshalIndent(recipies, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return ""
	}

	file, err := os.Create(resultFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return ""
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return ""
	}

	return resultFile
}
