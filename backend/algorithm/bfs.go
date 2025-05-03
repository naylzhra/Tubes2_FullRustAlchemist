package algorithm

import (
	"backend/search"
	"fmt"
)

type Queue[T any] struct {
	data []T
}

func (q *Queue[T]) Enqueue(value T) {
	q.data = append(q.data, value)
}

func (q *Queue[T]) Dequeue() T {
	val := q.data[0]
	q.data = q.data[1:]
	return val
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.data) == 0
}

// Kalo shortest path, maxPathsnya = 1
func BFS(target *search.ElementNode, maxPaths int) [][]*search.ElementNode {
	type PathNode struct {
		Path []*search.ElementNode
	}

	queue := Queue[PathNode]{}
	queue.Enqueue(PathNode{Path: []*search.ElementNode{target}})

	var results [][]*search.ElementNode
	visited := make(map[string]bool)

	for !queue.IsEmpty() && len(results) < maxPaths {
		curr := queue.Dequeue()
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

			newPath := append([]*search.ElementNode{a, b}, curr.Path...)
			queue.Enqueue(PathNode{Path: newPath})
		}
	}

	return results
}

func PrintCraftingPath(path []*search.ElementNode) {
	if len(path) == 0 {
		fmt.Println("No crafting path found.")
		return
	}

	for i := len(path) - 1; i >= 1; i -= 2 {
		if path[i-1].Name == "" && path[i-2].Name == "" {
			continue
		} else {
			fmt.Printf("%s => %s + %s\n", path[i].Name, path[i-1].Name, path[i-2].Name)
		}
	}
}

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
// 	paths := BFS(target, maxPaths)
// 	fmt.Printf("Found %d crafting paths:\n", len(paths))
// 	for i, path := range paths {
// 		fmt.Printf("Path #%d:\n", i+1)
// 		PrintCraftingPath(path)
// 		fmt.Println()
// 	}

// }
