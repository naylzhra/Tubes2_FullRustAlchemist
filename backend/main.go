package main

import (
	"backend/scraping"
	"backend/search"
	"fmt"
)

func main() {
	// Load recipe data
	recipes, err := scraping.GetScrapedRecipesJSON()
	if err != nil {
		fmt.Printf("Error loading recipes: %v\n", err)
		return
	}

	// Construct recipe graph
	var graph search.RecipeGraph
	err = search.ConstructRecipeGraph(recipes, &graph)
	if err != nil {
		fmt.Printf("Error constructing recipe graph: %v\n", err)
		return
	}
}
