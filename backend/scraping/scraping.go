package scraping

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type RecipeEntry struct {
	Element []string              `json:"element"`
	Recipe  map[string][][]string `json:"recipe"`
}

func ScrapeRecipes() error {
	startTime := time.Now()

	// URL of the page to scrape
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
	doc, err := getDocument(url)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// JSON object to store the recipes
	// recipes["elements"] : a list of all elements
	// recipes["recipes"] : a map of elements (in string) to recipes
	recipesJSON := RecipeEntry{
		Element: make([]string, 0),
		Recipe:  make(map[string][][]string),
	}

	// Scraping
	recipe_tables := doc.Find("table")
	for i := range recipe_tables.Length() {
		if i == 0 {
			continue // this annoying table
		}

		recipe_tables.Eq(i).Find("tr").Each(func(index int, row *goquery.Selection) {
			// Not the expected table. Only read the table on which the first row, first column is "Element"
			// and the second column is "Recipes"
			if index == 0 {
				if row.Find("td").Eq(0).Text() != "Element" || row.Find("td").Eq(1).Text() != "Recipes" {
					return
				}
			}

			columns := row.Find("td")
			// First column is the element
			element := columns.Eq(0).Text()
			if element != "" {
				element = strings.TrimSpace(element)
				recipesJSON.Element = append(recipesJSON.Element, element)
				recipesJSON.Recipe[element] = make([][]string, 0)
			}

			// Second column is the recipe
			recipes := columns.Eq(1).Find("li")
			recipes.Each(func(index int, recipe *goquery.Selection) {
				recipe_text := strings.Split(recipe.Text(), "+")
				if len(recipe_text) > 1 {
					recipe_text[0] = strings.TrimSpace(recipe_text[0])
					recipe_text[1] = strings.TrimSpace(recipe_text[1])

					recipesJSON.Recipe[element] = append(recipesJSON.Recipe[element], recipe_text)
				}
			})
			if recipes.Length() == 0 {
				// If there are no recipes, add an empty recipe
				recipesJSON.Recipe[element] = append(recipesJSON.Recipe[element], []string{})
			}
		})
	}

	// Export the recipes to JSON file
	filename, err := exportJSON(recipesJSON)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	total_recipes := 0
	for _, recipes := range recipesJSON.Recipe {
		total_recipes += len(recipes)
	}
	fmt.Println("Scraping completed. Recipes exported to ", filename)
	fmt.Println("Number of elements:", len(recipesJSON.Element))
	fmt.Println("Number of recipes loaded:", len(recipesJSON.Recipe))
	fmt.Println("Total number of recipes:", total_recipes)
	fmt.Println("Elapsed time:", elapsedTime.Milliseconds(), "ms")

	return nil
}

func getDocument(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", res.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func exportJSON(recipesJSON RecipeEntry) (string, error) {
	// Create the JSON file
	filename := "scraping/recipes.json"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer file.Close()

	// Write the JSON to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(recipesJSON)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	return filename, nil
}
