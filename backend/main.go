package main

import (
	"backend/scraping"
	"backend/search"
	"fmt"
)

func main() {
	err := scraping.ScrapeRecipes(true)
	if err != nil {
		panic(err)
	}
	recipes, err := scraping.GetScrapedRecipesJSON()
	if err != nil {
		panic(err)
	}

	// Initialize the recipe graph
	graph := search.RecipeGraph{}
	err = search.ConstructRecipeGraph(recipes, &graph)
	if err != nil {
		panic(err)
	}

	fmt.Println("Recipe graph constructed successfully!")
	fmt.Println("Length of element nodes:", len(graph.Elements))
	total_edges := 0
	for _, element := range graph.Elements {
		total_edges += len(element.Recipes)
	}
	fmt.Println("Total number of edges:", total_edges)

	elem, err := search.GetElementByName(&graph, "Treasure")
	if err == nil {
		fmt.Println(search.GetName(elem))
		fmt.Println("Number of children:", len(elem.Children))
		fmt.Println("Number of recipes:", len(elem.Recipes))
		fmt.Println("Tier:", elem.Tier)
	}
}
