package main

import (
	"backend/algorithm"
	"backend/scraping"
	"backend/search"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	
	fmt.Println("Recipe graph constructed successfully!")
	fmt.Printf("Number of elements: %d\n", len(graph.Elements)-1) // Exclude sentinel

	scanner := bufio.NewScanner(os.Stdin)
	
	// Get the element to search for
	fmt.Print("Enter the element you want to find recipes for: ")
	scanner.Scan()
	elementName := scanner.Text()
	
	// Find the element in the graph
	element, err := search.GetElementByName(&graph, elementName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Print recipes for this element
	fmt.Printf("\nRecipes in JSON for %s:\n", elementName)
	for i, recipe := range element.Recipes {
		if len(recipe) == 2 {
			fmt.Printf("  %d. %s + %s\n", i+1, recipe[0].Name, recipe[1].Name)
		}
	}
	
	fmt.Printf("\nRecipes in graph for %s:\n", elementName)
	for i, recipe := range element.Recipes {
		if len(recipe) == 2 {
			fmt.Printf("  %d. %s + %s\n", i+1, recipe[0].Name, recipe[1].Name)
		}
	}
	
	// Find crafting paths
	fmt.Print("\nHow many paths do you want to find? ")
	scanner.Scan()
	maxPathsStr := scanner.Text()
	maxPaths, err := strconv.Atoi(strings.TrimSpace(maxPathsStr))
	if err != nil {
		fmt.Printf("Invalid number: %v. Using default of 5.\n", err)
		maxPaths = 5
	}
	
	// Find and print crafting paths
	paths := algorithm.FindCraftingPaths(element, maxPaths)
	if len(paths) == 0 {
		fmt.Println("No path found.")
		return
	}
	for i, path := range paths {
		fmt.Printf("\nPath %d:\n", i+1)
		algorithm.PrintCraftingPath(path)
	}
	jsonFile, err := algorithm.PathsToNestedJSON(element.Name, paths)
	if err != nil {
		fmt.Printf("Error writing JSON: %v\n", err)
		return
	}
	fmt.Printf("✓ Paths saved to %s\n", jsonFile)
}