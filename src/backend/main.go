package main

import (
    "net/http"
    "backend/scraping"
    "backend/search"
    "backend/algorithm"
    "strconv"
    "fmt"
    
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
        algo    := strings.ToLower(c.DefaultQuery("algo", "bfs"))

        if element == "" {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": true,
                "type": "missing_parameter",
                "message": "Element parameter is required",
            })
            return
        }

        // find the node
        node, err := search.GetElementByName(&graph, element)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{
                "error": true,
                "type": "element_not_found",
                "message": fmt.Sprintf("Element '%s' not found", element),
            })
            return
        }

        switch algo {
        case "bfs":
            result := algorithm.ReverseBFS(node, 1)
            c.JSON(http.StatusOK, gin.H{
                "error": false,
                "data": result,
            })
        case "dfs":
            result := algorithm.DFS(node, &graph, 1)
            c.JSON(http.StatusOK, result)
        default:
            c.JSON(http.StatusBadRequest, gin.H{
                "error": true,
                "type": "invalid_algorithm",
                "message": "Algorithm must be 'bfs' or 'dfs'",
            })
        }
    })

    r.GET("/api/recipes", func(c *gin.Context) {
    element := c.Query("element")
    algo    := strings.ToLower(c.DefaultQuery("algo", "bfs"))

    if element == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": true,
            "type": "missing_parameter",
            "message": "Element parameter is required",
        })
        return
    }

    max, _ := strconv.Atoi(c.DefaultQuery("max", "5"))
    if max <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": true,
            "type": "invalid_parameter",
            "message": "Max parameter must be greater than 0",
        })
        return
    }

    node, err := search.GetElementByName(&graph, element)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": true,
            "type": "element_not_found",
            "message": fmt.Sprintf("Element '%s' not found", element),
        })
        return
    }

    paths := make([]*algorithm.GraphJSONWithRecipes, 0, max)
    algorithm.ResetCaches()
    
    var p *algorithm.GraphJSONWithRecipes
    switch algo {
        case "bfs":
            for i := 0; i < max; i++ {
                p = algorithm.ReverseBFS(node, i+1)
                if len(p.Recipes) == 0 { // no more unique paths
                    break
                }
                paths = append(paths, p)
            }
            c.JSON(http.StatusOK, gin.H{
                "error": false,
                "data": gin.H{
                    "element": element,
                    "algo": algo,
                    "paths": paths,
                },
            })
        // case "dfs":
        //     p = algorithm.ReverseDFS(node, i+1) // implement if needed
        default:
            c.JSON(http.StatusBadRequest, gin.H{
            "error": true,
            "type": "invalid_algorithm",
            "message": "Algorithm must be 'bfs' or 'dfs'",
            })
            return
    }

    })

    r.Run(":8080")
}
