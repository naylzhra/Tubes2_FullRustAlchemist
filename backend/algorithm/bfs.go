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
func makeRecipeKey(elem1, elem2, result string) string {
	elements := []string{elem1, elem2}
	sort.Strings(elements)
	return fmt.Sprintf("%s+%s=%s", elements[0], elements[1], result)
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
	type QueueNode struct {
		Node *search.ElementNode
		Step int
	}

	var nodes []JSONNode
	var recipes []JSONRecipe
	var queue []QueueNode

	nodesToInclude := make(map[int]bool)
	localVisited := make(map[int]bool)

	queue = append(queue, QueueNode{Node: target, Step: 0})

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if localVisited[curr.Node.ID] && !isBaseElement(curr.Node) {
			continue
		}
		localVisited[curr.Node.ID] = true

		visited[curr.Node.ID] = true

		if isNoRecipe(curr.Node) {
			continue
		}

		if !nodesToInclude[curr.Node.ID] {
			nodesToInclude[curr.Node.ID] = true
			nodes = append(nodes, JSONNode{
				ID:   curr.Node.ID,
				Name: curr.Node.Name,
			})
		}

		if isBaseElement(curr.Node) {
			continue
		}

		for _, recipe := range curr.Node.Recipes {
			if len(recipe) != 2 || recipe[0] == nil || recipe[1] == nil {
				continue
			}

			if recipe[0].Tier >= curr.Node.Tier || recipe[1].Tier >= curr.Node.Tier {
				continue
			}

			if recipe[0] != nil && recipe[1] != nil {
				recipeKey := makeRecipeKey(recipe[0].Name, recipe[1].Name, curr.Node.Name)
				if usedRecipes[recipeKey] {
					continue
				}
				usedRecipes[recipeKey] = true
			}

			// Tambahkan kedua parent ke queue
			recipeIngredients := []string{}
			validRecipe := true

			for _, parent := range recipe {
				if parent != nil {
					recipeIngredients = append(recipeIngredients, parent.Name)
					if !localVisited[parent.ID] {
						queue = append(queue, QueueNode{
							Node: parent,
							Step: curr.Step + 1,
						})
					}
				} else {
					validRecipe = false
				}
			}
			if validRecipe && len(recipeIngredients) == 2 {
				recipes = append(recipes, JSONRecipe{
					Ingredients: recipeIngredients,
					Result:      curr.Node.Name,
					Step:        curr.Step,
				})
			}
			break
		}
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

	// TIDAK mereset visited global agar pencarian berikutnya
	// akan menghindari path yang sudah dikunjungi sebelumnya

	// Bersihkan usedRecipes sebelum memulai pencarian
	usedRecipes = make(map[string]bool)

	for i := 0; i < maxPaths; i++ {
		fmt.Printf("Finding path %d...\n", i+1)
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
