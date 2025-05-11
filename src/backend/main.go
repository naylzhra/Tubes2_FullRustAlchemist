package main

import (
	"backend/algorithm"
	"backend/scraping"
	"backend/search"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"strings"
)

func main() {
	if err := scraping.ScrapeRecipes(false); err != nil {
		panic(err)
	}
	recipes, err := scraping.GetScrapedRecipesJSON()
	if err != nil {
		panic(err)
	}

	var graph search.RecipeGraph
	if err := search.ConstructRecipeGraph(recipes, &graph); err != nil {
		panic(err)
	}

	// target, err := search.GetElementByName(&graph, "Picnic")
	// if err == nil {
	// 	algorithm.DFS(target, &graph, 1)
	// }
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET"},
		AllowHeaders: []string{"Content-Type"},
	}))

	// http://localhost:8080/api/recipe?element=Acid%20Rain&algo=bfs|dfs
	r.GET("/api/recipe", func(c *gin.Context) {
		algorithm.ResetCaches()
		element := c.Query("element")
		algo := strings.ToLower(c.DefaultQuery("algo", "bfs"))

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

	r.GET("/api/recipes", func(c *gin.Context) {
		element := c.Query("element")
		algo := strings.ToLower(c.DefaultQuery("algo", "bfs"))

		max, _ := strconv.Atoi(c.DefaultQuery("max", "5"))
		if max <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "max must be > 0"})
			return
		}

		node, err := search.GetElementByName(&graph, element)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "element not found"})
			return
		}

		paths := make([]*algorithm.GraphJSONWithRecipes, 0, max)
		algorithm.ResetCaches()
		for i := 0; i < max; i++ {

			var p *algorithm.GraphJSONWithRecipes
			switch algo {
			case "bfs":
				p = algorithm.ReverseBFS(node, i+1)
			// case "dfs":
			//     p = algorithm.ReverseDFS(node, i+1) // implement if needed
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "algo must be bfs or dfs"})
				return
			}

			if len(p.Recipes) == 0 { // no more unique paths
				break
			}
			paths = append(paths, p)
		}

		c.JSON(http.StatusOK, gin.H{
			"element": element,
			"algo":    algo,
			"paths":   paths,
		})
	})

	r.Run(":8080")
}
