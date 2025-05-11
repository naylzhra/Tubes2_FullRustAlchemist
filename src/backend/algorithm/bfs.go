package algorithm

import (
	"backend/search"
	"fmt"
	"sort"
)

var visited = make(map[int]bool)
var usedElemComb = make(map[string]map[string]bool)

type AncestryChain struct {
	Element string
	Parents *AncestryChain
}

func getAncSignature(chain *AncestryChain) string {
	if chain == nil {
		return ""
	}

	if chain.Parents == nil {
		return chain.Element
	}

	return fmt.Sprintf("%s<-%s", chain.Element, getAncSignature(chain.Parents))
}

func isElemCombUsed(result string, elem1ID, elem2ID int, ancestryChain *AncestryChain) bool {
	ids := []int{elem1ID, elem2ID}
	sort.Ints(ids)
	combKey := fmt.Sprintf("%d+%d", ids[0], ids[1])
	ancSignature := getAncSignature(ancestryChain)
	combMapKey := fmt.Sprintf("%s:%s", result, ancSignature)

	if _, exists := usedElemComb[combMapKey]; !exists {
		usedElemComb[combMapKey] = make(map[string]bool)
		return false
	}

	return usedElemComb[combMapKey][combKey]
}

func markElemCombUsed(result string, elem1ID, elem2ID int, ancestryChain *AncestryChain) {
	ids := []int{elem1ID, elem2ID}
	sort.Ints(ids)
	combKey := fmt.Sprintf("%d+%d", ids[0], ids[1])

	ancSignature := getAncSignature(ancestryChain)
	combMapKey := fmt.Sprintf("%s:%s", result, ancSignature)

	if _, exists := usedElemComb[combMapKey]; !exists {
		usedElemComb[combMapKey] = make(map[string]bool)
	}

	usedElemComb[combMapKey][combKey] = true
}

type JSONNode struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type BFSState struct {
	Node          *search.ElementNode
	Parent        *BFSState
	Depth         int
	AncestryChain *AncestryChain
}

type JSONEdge struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type GraphJSON struct {
	Nodes []JSONNode `json:"nodes"`
	Edges []JSONEdge `json:"edges"`
}

func isNoRecipe(node *search.ElementNode) bool {
	for _, recipe := range node.Recipes {
		if len(recipe) == 2 && recipe[0].Name == "" && recipe[1].Name == "" {
			return true
		}
	}
	return false
}

func isBaseElement(node *search.ElementNode) bool {
	if node.Name == "Air" || node.Name == "Earth" || node.Name == "Fire" || node.Name == "Water" {
		return true
	}
	return false
}

type JSONRecipe struct {
	Ingredients []string `json:"ingredients"`
	Result      string   `json:"result"`
	Step        int      `json:"step"`
}

type GraphJSONWithRecipes struct {
	Nodes   []JSONNode   `json:"nodes"`
	Recipes []JSONRecipe `json:"recipes"`
}

func ReverseBFS(target *search.ElementNode, pathNumber int) (*GraphJSONWithRecipes, int) {
	visitedNodes := 0

	if isBaseElement(target) {
		return nil, 0
	}

	var nodes []JSONNode
	var recipes []JSONRecipe

	//processedNodes := make(map[int]bool)
	nodesToInclude := make(map[int]bool)

	type QueueItem struct {
		Node          *search.ElementNode
		AncestryChain *AncestryChain
		Depth         int
	}

	// Start with target node and empty ancestry
	queue := []QueueItem{
		{
			Node: target,
			AncestryChain: &AncestryChain{
				Element: target.Name,
				Parents: nil,
			},
			Depth: 0,
		},
	}

	// Add target to results
	nodesToInclude[target.ID] = true
	nodes = append(nodes, JSONNode{
		ID:   target.ID,
		Name: target.Name,
	})

	nodeStep := make(map[int]int)
	nodeStep[target.ID] = 0

	maxIterations := 1000
	iteration := 0

	for len(queue) > 0 && iteration < maxIterations {
		iteration++
		visitedNodes++
		current := queue[0]
		queue = queue[1:]

		if isBaseElement(current.Node) {
			//baseElementsFound[current.Node.Name] = true
			continue
		}

		if isNoRecipe(current.Node) {
			continue
		}

		recipeFound := false
		for _, recipe := range current.Node.Recipes {
			if len(recipe) != 2 || recipe[0] == nil || recipe[1] == nil {
				continue
			}

			if (isNoRecipe(recipe[0]) && !isBaseElement(recipe[0])) || (isNoRecipe(recipe[1]) && !isBaseElement(recipe[1])) {
				continue
			}

			if recipe[0].Tier >= current.Node.Tier || recipe[1].Tier >= current.Node.Tier {
				continue
			}

			if isElemCombUsed(current.Node.Name, recipe[0].ID, recipe[1].ID, current.AncestryChain) {
				continue
			}

			markElemCombUsed(current.Node.Name, recipe[0].ID, recipe[1].ID, current.AncestryChain)

			recipeFound = true
			recipes = append(recipes, JSONRecipe{
				Ingredients: []string{recipe[0].Name, recipe[1].Name},
				Result:      current.Node.Name,
				Step:        current.Depth,
			})

			for _, ingredient := range recipe {
				if ingredient == nil {
					continue
				}

				if !nodesToInclude[ingredient.ID] {
					nodesToInclude[ingredient.ID] = true
					nodes = append(nodes, JSONNode{
						ID:   ingredient.ID,
						Name: ingredient.Name,
					})
				}

				newAncestry := &AncestryChain{
					Element: ingredient.Name,
					Parents: current.AncestryChain,
				}

				nodeStep[ingredient.ID] = current.Depth + 1

				if !isBaseElement(ingredient) {
					queue = append(queue, QueueItem{
						Node:          ingredient,
						AncestryChain: newAncestry,
						Depth:         current.Depth + 1,
					})
				} else {
					visitedNodes++
				}
			}
			break
		}

		if !recipeFound && !visited[current.Node.ID] {
			queue = append(queue, current)
		}
	}

	if iteration >= maxIterations {
		fmt.Printf("Warning: Reached max iterations (%d) for path %d\n", maxIterations, pathNumber)
	}

	return &GraphJSONWithRecipes{
		Nodes:   nodes,
		Recipes: recipes,
	}, visitedNodes
}

func ResetCaches() {
	visited = make(map[int]bool)
	usedElemComb = make(map[string]map[string]bool)
}

// func main() {
// 	err := scraping.ScrapeRecipes(false)
// 	if err != nil {
// 		log.Fatal("Error while scraping recipes:", err)
// 	}

// 	recipes, err := scraping.GetScrapedRecipesJSON()
// 	if err != nil {
// 		log.Fatal("Error loading recipes from JSON:", err)
// 	}

// 	var graph search.RecipeGraph
// 	err = search.ConstructRecipeGraph(recipes, &graph)
// 	if err != nil {
// 		log.Fatal("Error constructing recipe graph:", err)
// 	}

// 	reader := bufio.NewReader(os.Stdin)

// 	fmt.Print("Enter target element name: ")
// 	targetName, _ := reader.ReadString('\n')
// 	targetName = strings.TrimSpace(targetName)

// 	target, err := search.GetElementByName(&graph, targetName)
// 	if err != nil {
// 		log.Fatalf("Error: element '%s' not found.\n", targetName)
// 	}

// 	fmt.Print("Enter number of paths to find: ")
// 	inputMax, _ := reader.ReadString('\n')
// 	inputMax = strings.TrimSpace(inputMax)
// 	maxPaths, err := strconv.Atoi(inputMax)
// 	if err != nil || maxPaths <= 0 {
// 		log.Fatalf("Invalid number: %v\n", inputMax)
// 	}

// 	usedElemComb = make(map[string]map[string]bool)
// 	pathsFound := 0
// 	consecutiveFailures := 0
// 	maxConsecutiveFailures := 3 // Number of empty path results before giving up

// 	for i := 0; i < maxPaths; i++ {
// 		fmt.Printf("Finding path %d...\n", i+1)
// 		startTime := time.Now()

// 		visited = make(map[int]bool)
// 		result, tmp := ReverseBFS(target, i+1)

// 		// If no recipes found in this path
// 		if len(result.Recipes) == 0 {
// 			consecutiveFailures++
// 			fmt.Printf("No recipes found for path %d (attempt %d of %d).\n",
// 				i+1, consecutiveFailures, maxConsecutiveFailures)

// 			// If we've had multiple failures in a row, assume no more paths exist
// 			if consecutiveFailures >= maxConsecutiveFailures {
// 				fmt.Printf("No more paths found after %d consecutive empty results.\n", consecutiveFailures)
// 				break
// 			}

// 			// Skip saving this empty result and try again
// 			continue
// 		}

// 		// Reset the failure counter since we found a valid path
// 		consecutiveFailures = 0
// 		pathsFound++

// 		filename := fmt.Sprintf("graph_output_%d.json", pathsFound)
// 		jsonBytes, err := json.MarshalIndent(result, "", "  ")
// 		if err != nil {
// 			log.Fatalf("Failed to marshal JSON: %v", err)
// 		}
// 		err = os.WriteFile(filename, jsonBytes, 0644)
// 		if err != nil {
// 			log.Fatalf("Failed to write file: %v", err)
// 		}

// 		elapsedTime := time.Since(startTime)
// 		fmt.Printf("Saved path %d to '%s' (took %.2f seconds)\n",
// 			pathsFound, filename, elapsedTime.Seconds())

// 		fmt.Printf("Path %d has %d nodes and %d recipes\n",
// 			pathsFound, len(result.Nodes), len(result.Recipes))

// 		// if len(result.Recipes) > 0 {
// 		// 	fmt.Println("Cek ancestry signatures:")
// 		// 	count := 0
// 		// 	for key := range usedElemComb {
// 		// 		if count >= 3 {
// 		// 			break
// 		// 		}
// 		// 		fmt.Printf("  %s\n", key)
// 		// 		count++
// 		// 	}
// 		// }
// 	}

// 	if pathsFound < maxPaths {
// 		fmt.Printf("\nFound %d paths out of %d requested. No more unique paths available.\n",
// 			pathsFound, maxPaths)
// 	} else {
// 		fmt.Printf("\nSuccessfully found all %d requested paths.\n", maxPaths)
// 	}
// }
