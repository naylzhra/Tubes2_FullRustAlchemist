package algorithm

import (
	"backend/search"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sync"
)

type ResultTree struct {
	mu   sync.Mutex
	path []*Recipe
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

	return findMultiplePaths(target, graph, maxPaths, nodeVisited)
}

func mergeTree(tree0 *ResultTree, tree1 *ResultTree, resulto *ResultTree) {
	resulto.path = append(resulto.path, tree0.path...)
	resulto.path = append(resulto.path, tree1.path...)
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

/* ----------------------------------------- Multiple Recipe DFS ----------------------------------------------- */

type SearchStatus struct {
	result         chan int
	continueSignal chan int
}

type SearchStatistic struct {
	mu             sync.Mutex
	remainingPaths int
	nodeVisited    int
}

func findMultiplePaths(target *search.ElementNode, graph *search.RecipeGraph, maxPaths int, nodeVisited *int) []string {
	resultJSONs := make([]string, 0, maxPaths)

	status := SearchStatus{
		result:         make(chan int),
		continueSignal: make(chan int),
	}
	stats := &SearchStatistic{
		mu:             sync.Mutex{},
		remainingPaths: maxPaths,
		nodeVisited:    0,
	}
	result := &ResultTree{path: make([]*Recipe, 0)}

	go findPath(target, graph, result, status, stats)

	counter := 0
	condition := <-status.result
	for condition != 0 {
		counter++

		newFilename := ParseCraftingPath(result, counter, graph)
		resultJSONs = append(resultJSONs, newFilename)

		if counter >= maxPaths {
			status.continueSignal <- 0
			<-status.result
			break
		} else {
			status.continueSignal <- 1
			condition = <-status.result
		}
	}

	stats.mu.Lock()
	*nodeVisited = stats.nodeVisited
	stats.mu.Unlock()
	return resultJSONs
}

func findPath(target *search.ElementNode, graph *search.RecipeGraph, result *ResultTree, status SearchStatus, stats *SearchStatistic) {
	stats.mu.Lock()
	stats.nodeVisited++
	stats.mu.Unlock()

	// Base case: if the target is a base element, return
	if slices.Contains(graph.BaseElements, target) {
		baseRecipe := &Recipe{element: target}
		baseRecipe.composition = []*Recipe{baseRecipe, baseRecipe}
		result.mu.Lock()
		result.path = append(result.path, baseRecipe)
		result.mu.Unlock()

		status.result <- 1
		<-status.continueSignal
		status.result <- 0
		return
	}
	if target.Name == "Time" {
		status.result <- 0
		return
	}

	for _, recipe := range target.Recipes {
		if recipe[0].Tier >= target.Tier || recipe[1].Tier >= target.Tier {
			continue
		}

		status0 := SearchStatus{result: make(chan int), continueSignal: make(chan int)}
		result0 := &ResultTree{path: make([]*Recipe, 0)}
		go findPath(recipe[0], graph, result0, status0, stats)

		status1 := SearchStatus{result: make(chan int), continueSignal: make(chan int)}
		result1 := &ResultTree{path: make([]*Recipe, 0)}
		go findPath(recipe[1], graph, result1, status1, stats)

		condition0 := <-status0.result
		condition1 := <-status1.result
		for condition0 != 0 && condition1 != 0 {
			if condition0 == 1 && condition1 == 1 {
				recipe := &Recipe{
					element:     target,
					composition: []*Recipe{result0.path[0], result1.path[0]},
				}
				result.mu.Lock()
				result.path = make([]*Recipe, 0)
				result.path = append(result.path, recipe)
				mergeTree(result0, result1, result)
				result.mu.Unlock()

				status.result <- 1
				continueSearch := <-status.continueSignal

				if continueSearch == 0 {
					status0.continueSignal <- 0
					status1.continueSignal <- 0
					<-status0.result
					<-status1.result

					break
				} else {
					status0.continueSignal <- 1
					condition0 = <-status0.result
				}
			} else if condition0 == 0 {
				status1.continueSignal <- 0

				status0 = SearchStatus{result: make(chan int), continueSignal: make(chan int)}
				result0 = &ResultTree{path: make([]*Recipe, 0)}
				go findPath(recipe[0], graph, result0, status0, stats)

				condition0 = <-status0.result
				condition1 = <-status1.result
			}
		}
	}

	status.result <- 0
	// Kill this routine
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
