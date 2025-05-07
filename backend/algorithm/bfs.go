package algorithm

import (
	"backend/search"
)

type JSONNode struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type JSONEdge struct {
	From int32 `json:"from"`
	To   int32 `json:"to"`
}

type GraphJSON struct {
	Nodes []JSONNode `json:"nodes"`
	Edges []JSONEdge `json:"edges"`
}

func isBaseElement(node *search.ElementNode) bool {
	for _, recipe := range node.Recipes {
		if len(recipe) == 2 && recipe[0].Name == "" && recipe[1].Name == "" {
			return true
		}
	}
	return false
}

type JSONRecipe struct {
	Ingredients []int32 `json:"ingredients"`
	Result      int32   `json:"result"`
}

type GraphJSONWithRecipes struct {
	Nodes   []JSONNode   `json:"nodes"`
	Recipes []JSONRecipe `json:"recipes"`
}

func ReverseBFS(target *search.ElementNode, maxPaths int) *GraphJSONWithRecipes {
	type pathItem struct {
		Node      *search.ElementNode
		PathDepth int
	}

	visited := make(map[int32]bool)
	nodes := make(map[int32]JSONNode)
	var recipes []JSONRecipe

	queue := []pathItem{{Node: target, PathDepth: 0}}
	pathsUsed := 0

	for len(queue) > 0 && pathsUsed < maxPaths {
		curr := queue[0]
		queue = queue[1:]

		if visited[curr.Node.ID] {
			continue
		}
		visited[curr.Node.ID] = true
		nodes[curr.Node.ID] = JSONNode{ID: curr.Node.ID, Name: curr.Node.Name}

		// Jika node ini adalah base element, langsung tambahkan resep [0, 0]
		if isBaseElement(curr.Node) {
			recipes = append(recipes, JSONRecipe{
				Ingredients: []int32{0, 0},
				Result:      curr.Node.ID,
			})
			continue
		}

		// Cek semua resep yang dimiliki node ini
		for _, recipe := range curr.Node.Recipes {
			if pathsUsed >= maxPaths {
				break
			}
			if len(recipe) != 2 || recipe[0] == nil || recipe[1] == nil {
				continue
			}

			// Simpan resep valid
			recipes = append(recipes, JSONRecipe{
				Ingredients: []int32{recipe[0].ID, recipe[1].ID},
				Result:      curr.Node.ID,
			})
			pathsUsed++

			// Tambahkan parent nodes ke antrian untuk ditelusuri
			for _, parent := range recipe {
				if parent != nil {
					nodes[parent.ID] = JSONNode{ID: parent.ID, Name: parent.Name}
					if !visited[parent.ID] {
						queue = append(queue, pathItem{Node: parent, PathDepth: curr.PathDepth + 1})
					}
				}
			}
			break // hanya satu path per node sesuai BFS
		}
	}

	// Susun nodeList dari map
	var nodeList []JSONNode
	for _, node := range nodes {
		nodeList = append(nodeList, node)
	}

	return &GraphJSONWithRecipes{
		Nodes:   nodeList,
		Recipes: recipes,
	}
}

// func PrintCraftingRecipesFromJSON(jsonData []byte, targetName string) {
// 	var graph GraphJSON
// 	err := json.Unmarshal(jsonData, &graph)
// 	if err != nil {
// 		fmt.Println("Error parsing JSON:", err)
// 		return
// 	}

// 	// Buat map ID -> Node dan Name -> ID
// 	idToNode := make(map[int32]JSONNode)
// 	nameToID := make(map[string]int32)
// 	for _, node := range graph.Nodes {
// 		idToNode[node.ID] = node
// 		nameToID[node.Name] = node.ID
// 	}

// 	targetID, ok := nameToID[targetName]
// 	if !ok {
// 		fmt.Println("Target not found in nodes:", targetName)
// 		return
// 	}

// 	// Bangun reverse map To → list of From
// 	childToParents := make(map[int32][]int32)
// 	for _, edge := range graph.Edges {
// 		childToParents[edge.To] = append(childToParents[edge.To], edge.From)
// 	}

// 	// DFS recursive untuk kumpulkan kombinasi
// 	var result [][]int32
// 	var dfs func(currID int32, path []int32)

// 	dfs = func(currID int32, path []int32) {
// 		path = append([]int32{currID}, path...)

// 		parents, exists := childToParents[currID]
// 		if !exists || len(parents) == 0 {
// 			result = append(result, path)
// 			return
// 		}

// 		for _, parentID := range parents {
// 			dfs(parentID, path)
// 		}
// 	}

// 	dfs(targetID, []int32{})

// 	// Cetak kombinasi resep
// 	fmt.Printf("Crafting recipes to form: %s\n", targetName)
// 	for i, recipe := range result {
// 		fmt.Printf("Path %d: ", i+1)
// 		for j, id := range recipe {
// 			if j > 0 {
// 				fmt.Print(" + ")
// 			}
// 			fmt.Print(idToNode[id].Name)
// 		}
// 		fmt.Println()
// 	}
// }

// func main() {
// 	err := scraping.ScrapeRecipes()
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

// 	fmt.Print("Enter target element name: ")
// 	reader := bufio.NewReader(os.Stdin)
// 	targetName, _ := reader.ReadString('\n')
// 	targetName = strings.TrimSpace(targetName)
// 	fmt.Print("Enter maximum number of crafting paths to show: ")
// 	inputMax, _ := reader.ReadString('\n')
// 	inputMax = strings.TrimSpace(inputMax)
// 	maxPaths, err := strconv.Atoi(inputMax)
// 	if err != nil || maxPaths <= 0 {
// 		log.Fatalf("Invalid number: %v\n", inputMax)
// 	}

// 	target, err := search.GetElementByName(&graph, targetName)
// 	if err != nil {
// 		log.Fatalf("Error: element '%s' not found.\n", targetName)
// 	}
// 	fmt.Printf("Crafting path to: %s\n", targetName)
// 	result := ReverseBFS(target, maxPaths)

// 	jsonBytes, err := json.MarshalIndent(result, "", "  ")
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = os.WriteFile("graph_output.json", jsonBytes, 0644)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// Mencetak resep crafting untuk target
// 	//PrintCraftingRecipesFromJSON(jsonBytes, target.Name)
// }
