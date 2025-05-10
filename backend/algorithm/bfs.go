package main

import (
	"backend/scraping"
	"backend/search"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var visited = make(map[int]bool)
var usedRecipes = make(map[string]bool)

// Fungsi untuk membuat kunci unik dari recipe
func makeRecipeKeyWithDepth(elem1, elem2, result string, depth int) string {
	elements := []string{elem1, elem2}
	sort.Strings(elements)
	return fmt.Sprintf("%s+%s=%s:%d", elements[0], elements[1], result, depth)
}

type JSONNode struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type BFSState struct {
	Node   *search.ElementNode
	Parent *BFSState
	Depth  int
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

func ReverseBFS(target *search.ElementNode, pathNumber int) *GraphJSONWithRecipes {
	var nodes []JSONNode
	var recipes []JSONRecipe

	processedNodes := make(map[int]bool)
	nodesToInclude := make(map[int]bool)
	pendingNodes := make(map[int]*search.ElementNode)

	nodesToInclude[target.ID] = true
	nodes = append(nodes, JSONNode{
		ID:   target.ID,
		Name: target.Name,
	})

	if !isBaseElement(target) {
		pendingNodes[target.ID] = target
	}

	maxIterations := 1000
	iteration := 0

	nodeStep := make(map[int]int)
	nodeStep[target.ID] = 0

	baseElementsFound := map[string]bool{
		"Air":   false,
		"Earth": false,
		"Fire":  false,
		"Water": false,
	}

	for len(pendingNodes) > 0 && iteration < maxIterations {
		iteration++

		var currentNode *search.ElementNode
		for _, node := range pendingNodes {
			currentNode = node
			delete(pendingNodes, currentNode.ID)
			break
		}

		if processedNodes[currentNode.ID] || isNoRecipe(currentNode) {
			fmt.Print(currentNode.Name)
			fmt.Println("SKIPPED PROCESSED NODES")
			continue
		}

		processedNodes[currentNode.ID] = true

		recipeFound := false
		for _, recipe := range currentNode.Recipes {
			if len(recipe) != 2 || recipe[0] == nil || recipe[1] == nil {
				continue
			}

			if recipe[0].Tier >= currentNode.Tier || recipe[1].Tier >= currentNode.Tier {
				continue
			}

			recipeKey := makeRecipeKeyWithDepth(recipe[0].Name, recipe[1].Name, currentNode.Name, nodeStep[currentNode.ID])
			if usedRecipes[recipeKey] {
				fmt.Print(currentNode.Name)
				fmt.Println("SKIPPED USED RECIPES")
				continue
			}

			// Recipe valid! Tandai sebagai digunakan
			recipeFound = true
			usedRecipes[recipeKey] = true

			// Tambahkan recipe ke hasil
			recipes = append(recipes, JSONRecipe{
				Ingredients: []string{recipe[0].Name, recipe[1].Name},
				Result:      currentNode.Name,
				Step:        nodeStep[currentNode.ID],
			})

			// Tandai node sebagai dikunjungi secara global
			//visited[currentNode.ID] = true

			// Proses ingredient nodes
			for _, ingredient := range recipe {
				if ingredient == nil {
					continue
				}

				// Tambahkan ke nodes jika belum ada
				if !nodesToInclude[ingredient.ID] {
					nodesToInclude[ingredient.ID] = true
					nodes = append(nodes, JSONNode{
						ID:   ingredient.ID,
						Name: ingredient.Name,
					})
				}

				// Update step untuk ingredient
				newStep := nodeStep[currentNode.ID] + 1
				if existingStep, exists := nodeStep[ingredient.ID]; !exists || newStep > existingStep {
					nodeStep[ingredient.ID] = newStep
				}

				// Cek apakah base element
				if isBaseElement(ingredient) {
					baseElementsFound[ingredient.Name] = true
				} else if !processedNodes[ingredient.ID] {
					// Tambahkan ke pending jika bukan base element
					pendingNodes[ingredient.ID] = ingredient
				}
			}

			// Hanya gunakan satu recipe valid
			break
		}

		// Jika tidak menemukan recipe valid dan belum dikunjungi secara global,
		// tambahkan kembali ke pending untuk coba lagi nanti
		if !recipeFound && !visited[currentNode.ID] {
			pendingNodes[currentNode.ID] = currentNode
		}
	}

	// Cek apakah semua base element ditemukan
	allBasesFound := baseElementsFound["Air"] && baseElementsFound["Earth"] &&
		baseElementsFound["Fire"] && baseElementsFound["Water"]

	if !allBasesFound {
		fmt.Printf("Warning: Not all base elements found in path %d. Found: Air=%v, Earth=%v, Fire=%v, Water=%v\n",
			pathNumber, baseElementsFound["Air"], baseElementsFound["Earth"],
			baseElementsFound["Fire"], baseElementsFound["Water"])
	}

	if iteration >= maxIterations {
		fmt.Printf("Warning: Reached max iterations (%d) for path %d\n", maxIterations, pathNumber)
	}

	return &GraphJSONWithRecipes{
		Nodes:   nodes,
		Recipes: recipes,
	}
}

func main() {
	err := scraping.ScrapeRecipes(false)
	if err != nil {
		log.Fatal("Error while scraping recipes:", err)
	}

	recipes, err := scraping.GetScrapedRecipesJSON()
	if err != nil {
		log.Fatal("Error loading recipes from JSON:", err)
	}

	var graph search.RecipeGraph
	err = search.ConstructRecipeGraph(recipes, &graph)
	if err != nil {
		log.Fatal("Error constructing recipe graph:", err)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter target element name: ")
	targetName, _ := reader.ReadString('\n')
	targetName = strings.TrimSpace(targetName)

	target, err := search.GetElementByName(&graph, targetName)
	if err != nil {
		log.Fatalf("Error: element '%s' not found.\n", targetName)
	}

	fmt.Print("Enter number of paths to find: ")
	inputMax, _ := reader.ReadString('\n')
	inputMax = strings.TrimSpace(inputMax)
	maxPaths, err := strconv.Atoi(inputMax)
	if err != nil || maxPaths <= 0 {
		log.Fatalf("Invalid number: %v\n", inputMax)
	}

	usedRecipes = make(map[string]bool)

	for i := 0; i < maxPaths; i++ {
		fmt.Printf("Finding path %d...\n", i+1)
		visited = make(map[int]bool)
		result := ReverseBFS(target, i+1) // Selalu cari satu path per panggilan dengan nomor path

		if len(result.Recipes) == 0 {
			fmt.Println("No more paths found.")
			break
		}

		filename := fmt.Sprintf("graph_output_%d.json", i+1)
		jsonBytes, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %v", err)
		}
		err = os.WriteFile(filename, jsonBytes, 0644)
		if err != nil {
			log.Fatalf("Failed to write file: %v", err)
		}
		fmt.Printf("Saved path %d to '%s'\n", i+1, filename)

		// Print debug info
		fmt.Printf("Path %d has %d nodes and %d recipes\n",
			i+1, len(result.Nodes), len(result.Recipes))
	}
}
