package main

import (
	"backend/scraping"
	"fmt"
)

func main() {
	err := scraping.ScrapeRecipes()
	if err != nil {
		panic(err)
	}
	recipes, err := scraping.GetScrapedRecipesJSON()
	if err != nil {
		panic(err)
	}

	// Initialize the recipe graph
	graph := RecipeGraph{Elements: make([]ElementNode, 0)}
	err = ConstructRecipeGraph(recipes, &graph)
	if err != nil {
		panic(err)
	}

	fmt.Println("Recipe graph constructed successfully!")
	fmt.Println("Length of element nodes:", len(graph.Elements))
	total_edges := 0
	for i, element := range graph.Elements {
		total_edges += len(element.Recipes)
		if i != 0 {
			continue
		}
		fmt.Printf("Element %d: %s\n", i, element.Name)
		fmt.Printf("Number of children: %d\n", len(element.Children))
		fmt.Printf("Number of recipes: %d\n", len(element.Recipes))
	}
	fmt.Println("Total number of edges:", total_edges)
}
