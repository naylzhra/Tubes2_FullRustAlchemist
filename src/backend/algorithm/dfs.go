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
	result         chan int
	continueSignal chan struct{}
}

type ResultTree struct {
	mu             sync.Mutex
	path           []*Recipe
	remainingPaths int
}

type Recipe struct {
	element     *search.ElementNode
	composition []*Recipe
}

func DFS(target *search.ElementNode, graph *search.RecipeGraph, maxPaths int, nodeVisited *int) []string {
	if maxPaths == 1 {
		result := &ResultTree{path: make([]*Recipe, 0)}
		findSinglePath(target, graph, result, nodeVisited)

		return []string{ParseCraftingPath(result, 1, graph)}
	}

	// Multithreaded multiple recipe DFS
	resultJSONs := make([]string, 0)
	return resultJSONs
}

/* ----------------------------------------- Single Recipe DFS ----------------------------------------------- */

func findSinglePath(target *search.ElementNode, graph *search.RecipeGraph, result *ResultTree, nodeVisited *int) *Recipe {
	*nodeVisited++

	if slices.Contains(graph.BaseElements, target) {
		*result = ResultTree{path: make([]*Recipe, 0)}
		baseElem := &Recipe{element: target}
		baseElem.composition = []*Recipe{baseElem, baseElem}
		result.path = append(result.path, baseElem)
		return baseElem
	}
	if target.Name == "Time" {
		return nil
	}

	// Try each recipe
	for _, recipe := range target.Recipes {
		if recipe[0].Tier >= target.Tier || recipe[1].Tier >= target.Tier {
			continue
		}

		result0 := &ResultTree{path: make([]*Recipe, 0)}
		component0 := findSinglePath(recipe[0], graph, result0, nodeVisited)
		if component0 == nil {
			continue
		}
		result1 := &ResultTree{path: make([]*Recipe, 0)}
		component1 := findSinglePath(recipe[1], graph, result1, nodeVisited)
		if component1 == nil {
			continue
		}

		mergeTree(result0, result1, result)
		validRecipe := &Recipe{
			element:     target,
			composition: []*Recipe{component0, component1},
		}
		result.path = append(result.path, validRecipe)
		return validRecipe
	}

	return nil
}

func mergeTree(tree0 *ResultTree, tree1 *ResultTree, resulto *ResultTree) {
	resulto.path = append(resulto.path, tree0.path...)
	resulto.path = append(resulto.path, tree1.path...)
}

/* ----------------------------------------- Multiple Recipe DFS ----------------------------------------------- */

func findPath(target *search.ElementNode, graph *search.RecipeGraph, result *ResultTree, status SearchStatus, prevs []*search.ElementNode) {
	// defer func() {
	// 	completePathSignal <- 0
	// }()

	// // Base case: if the target is a base element, return
	// if slices.Contains(graph.BaseElements, target) {
	// 	result.mu.Lock()
	// 	for _, elem := range result.paths {
	// 		if elem.Name == target.Name {
	// 			completePathSignal <- elem.ID
	// 		}
	// 	}
	// 	result.mu.Unlock()
	// 	<-completePathSignal
	// 	return
	// }

	// fmt.Println("Searching for:", target.Name)

	// // Generate thread for the two parents for a recipe (if not already created)
	// // Wait until the recipe is fully complete before checking another recipe of this node
	// for _, recipe := range target.Recipes {
	// 	if slices.Contains(graph.BaseElements, recipe[0]) && slices.Contains(graph.BaseElements, recipe[1]) {
	// 		result.mu.Lock()
	// 		firstID := -1
	// 		secondID := -1
	// 		for _, elem := range result.paths {
	// 			if elem.Name == recipe[0].Name {
	// 				firstID = elem.ID
	// 			}
	// 			if elem.Name == recipe[1].Name {
	// 				secondID = elem.ID
	// 			}
	// 		}
	// 		newID := addResult(result, target, firstID, secondID)
	// 		result.mu.Unlock()
	// 		completePathSignal <- newID
	// 		continueSearch := <-completePathSignal
	// 		if continueSearch == 0 {
	// 			break
	// 		} else {
	// 			continue
	// 		}
	// 	}

	// 	if recipe[0].Tier >= target.Tier || recipe[1].Tier >= target.Tier {
	// 		continue
	// 	}

	// 	// Decide whether to search for the parents of this element
	// 	if pathAlreadyContains(prevs, recipe[0]) || pathAlreadyContains(prevs, recipe[1]) {
	// 		continue
	// 	}

	// 	// Create a new thread for each parent of a recipe
	// 	currentPath := slices.Clone(prevs)
	// 	currentPath = append(currentPath, target)
	// 	searchStatus := SearchStatus{
	// 		pathSignal0: make(chan int),
	// 		pathSignal1: make(chan int),
	// 	}

	// 	// Let the first parent backtrack, then the second
	// 	// This is to enumerate all paths to recipe[0] x all paths to recipe[1]
	// 	go findPath(recipe[0], graph, result, searchStatus.pathSignal0, currentPath)
	// 	go findPath(recipe[1], graph, result, searchStatus.pathSignal1, currentPath)
	// 	status0 := <-searchStatus.pathSignal0
	// 	status1 := <-searchStatus.pathSignal1
	// 	for status1 != 0 && status0 != 0 {
	// 		if status0 == 0 { // First parent is done, so backtrack the second parent and restart the first parent thread
	// 			searchStatus.pathSignal1 <- 1
	// 			go findPath(recipe[0], graph, result, searchStatus.pathSignal0, currentPath)
	// 			status0 = <-searchStatus.pathSignal0
	// 			status1 = <-searchStatus.pathSignal1
	// 		} else if status0 == 1 && status1 == 1 {
	// 			result.mu.Lock()
	// 			newID := addResult(result, target, status0, status1)
	// 			result.mu.Unlock()
	// 			completePathSignal <- newID

	// 			continueSearch := <-completePathSignal
	// 			if continueSearch == 0 {
	// 				searchStatus.pathSignal0 <- 0
	// 				searchStatus.pathSignal1 <- 0
	// 				break
	// 			} else {
	// 				// Continue searching for the next recipe
	// 				searchStatus.pathSignal0 <- 1
	// 				status0 = <-searchStatus.pathSignal0
	// 			}
	// 		}
	// 	}
	// 	// Kill this routine
	// }
}

func pathAlreadyContains(prevs []*search.ElementNode, elem *search.ElementNode) bool {
	return slices.Contains(prevs, elem)
}

/* ----------------------------------------- Parse Search Output ----------------------------------------------- */

type NodeJSON struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RecipeJSON struct {
	Element string   `json:"element"`
	Recipe  []string `json:"recipe"`
}

type ResultJSON struct {
	Recipes []RecipeJSON
}

func ParseCraftingPath(result *ResultTree, counter int, graph *search.RecipeGraph) string {
	// ResultTree is locked in the caller side
	resultFile := "result_" + fmt.Sprintf("%03d", counter) + ".json"

	recipeToID := make(map[*Recipe]int)
	for i, recipe := range result.path {
		recipeToID[recipe] = i
	}
	pathJSON := make(map[string]RecipeJSON)
	for _, recipe := range result.path {
		if slices.Contains(graph.BaseElements, recipe.element) {
			continue
		}

		pathJSON[fmt.Sprintf("%d", recipeToID[recipe])] = RecipeJSON{
			Element: recipe.element.Name,
			Recipe:  make([]string, len(recipe.composition)),
		}
		for i, comp := range recipe.composition {
			pathJSON[fmt.Sprintf("%d", recipeToID[recipe])].Recipe[i] = fmt.Sprintf("%d", recipeToID[comp])
		}
	}

	jsonData, err := json.MarshalIndent(pathJSON, "", " ")
	if err != nil {
		fmt.Println("Error encoding JSON: ", err)
		return ""
	}

	// Write JSON to file
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
