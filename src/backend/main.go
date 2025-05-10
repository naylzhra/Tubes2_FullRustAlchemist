package main

import (
	"net/http"
	"backend/scraping"
	"backend/search"
	"backend/algorithm"
	
	"github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"

	"strings"
)

func main() {
    // 1.  Ensure the JSON is available.  Run the scraper once if needed.
    if err := scraping.ScrapeRecipes(false); err != nil {
        panic(err)
    }
    recipes, err := scraping.GetScrapedRecipesJSON()
    if err != nil {
        panic(err)
    }

    // 2.  Build the in-memory graph once at start-up.
    var graph search.RecipeGraph
    if err := search.ConstructRecipeGraph(recipes, &graph); err != nil {
        panic(err)
    }

    // 3.  Start Gin.
    r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
    r.Use(cors.New(cors.Config{
        AllowOrigins: []string{"http://localhost:3000"},
        AllowMethods: []string{"GET"},
        AllowHeaders: []string{"Content-Type"},
    }))

    // 4.  /api/recipe?element=Acid%20Rain&algo=bfs|dfs
    r.GET("/api/recipe", func(c *gin.Context) {
        algorithm.ResetCaches()
		element := c.Query("element")
        algo    := strings.ToLower(c.DefaultQuery("algo", "bfs"))

        // find the node
        node, err := search.GetElementByName(&graph, element)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "element not found"})
            return
        }

        switch algo {
        case "bfs":
            result := algorithm.ReverseBFS(node, 1)
            c.JSON(http.StatusOK, result)
        // case "dfs":
        //     result := algorithm.ReverseDFS(node, 1)
        //     c.JSON(http.StatusOK, result)
        default:
            c.JSON(http.StatusBadRequest, gin.H{"error": "algo must be bfs or dfs"})
        }
    })

    // 5.  Listen on :8080
    r.Run(":8080")
}
