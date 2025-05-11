"use client";
import React, { useEffect, useRef, useState } from "react";
import * as d3 from "d3";
import Image from "next/image";

/* ------------ type definitions ------------ */
export type GraphNode = { id: number; name: string };
export type GraphRecipe = { ingredients: string[]; result: string; step: number };
export interface GraphData {
  nodes: GraphNode[];
  recipes: GraphRecipe[];
  elapsed?: string;
}

/* props */
interface RecipeResultProps {
  graph: GraphData;
}

const RecipeResult: React.FC<RecipeResultProps> = ({ graph }) => {
  const svgRef = useRef<SVGSVGElement>(null);
  const [uniqueElements, setUniqueElements] = useState<string[]>([]);

  // fungsi untuk bikin tree dari graph
  function buildTree(target: string, recipes: GraphRecipe[]): any {
    const recipe = recipes.find(r => r.result === target);
    if (!recipe) {
      return { name: target };
    }
    return {
      name: target,
      children: recipe.ingredients.map(ingredient => buildTree(ingredient, recipes))
    };
  }

  // Function to get unique elements from recipes
  const getUniqueElements = (recipes: GraphRecipe[]): string[] => {
    const elements = new Set<string>();
    
    // Add all results
    recipes.forEach(recipe => {
      elements.add(recipe.result);
      // Add all ingredients
      recipe.ingredients.forEach(ing => elements.add(ing));
    });
    
    return Array.from(elements);
  };

  useEffect(() => {
    if (!graph.recipes || graph.recipes.length === 0) return;
    
    // Extract unique elements and set state
    const elements = getUniqueElements(graph.recipes);
    setUniqueElements(elements);
    
    // Increase width and height for better visualization
    const width = 1600;
    const height = 1200; // Increased height to accommodate icon legend
    if (!svgRef.current) return;
    
    // Clear existing SVG content
    const svg = d3.select(svgRef.current)
      .attr("viewBox", `0 0 ${width} ${height}`)
      .selectAll("*").remove();
      
    // Create a fresh SVG container
    const container = d3.select(svgRef.current)
      .attr("viewBox", `0 0 ${width} ${height}`);
    
    // Define larger margins to allow more space - added extra bottom margin for legend
    const margin = { top: 150, right: 200, bottom: 300, left: 200 };
    
    // Create main group with translation for margins
    const g = container.append("g")
      .attr("transform", `translate(${margin.left}, ${margin.top})`);
    
    const rootData = buildTree(graph.recipes[0]?.result ?? "", graph.recipes);
    const root = d3.hierarchy(rootData);
    
    // Adjust tree layout with proper dimensions accounting for margins
    const treeLayout = d3.tree<any>()
      .size([width - margin.left - margin.right, height - margin.top - margin.bottom - 150]) // Reduced height to make room for legend
      .separation((a, b) => {
        // Dynamically increase separation based on depth
        const baseMultiplier = 6; // Increased base multiplier
        const depthFactor = Math.pow(2, Math.max(a.depth, b.depth) * 0.3); // Exponential scaling based on depth
        return (a.parent === b.parent ? baseMultiplier : baseMultiplier * 1.5) * depthFactor;
      }); // Drastically increased separation for deeper levels
    
    treeLayout(root);
    
    // Garis antar node - add to the translated group
    const linkGenerator = d3.linkVertical<any, any>()
      .x((d: any) => d.x)
      .y((d: any) => d.y);
    
    g.append("g")
      .selectAll("path")
      .data(root.links())
      .join("path")
      .attr("fill", "none")
      .attr("stroke", "#555")
      .attr("stroke-width", 2)
      .attr("d", d => linkGenerator(d));
    
    // Node - add to the translated group with circular nodes
    const node = g.append("g")
      .selectAll("g")
      .data(root.descendants())
      .join("g")
      .attr("transform", d => `translate(${d.x},${d.y})`);
    
    // Size for circular nodes
    const circleRadius = 40;
    
    // Add circular background for nodes
    node.append("circle")
      .attr("r", circleRadius)
      .attr("fill", "#677D6A")
      .attr("stroke", "white")
      .attr("stroke-width", 2);
    
    // Add foreignObject to hold images
    node.append("foreignObject")
      .attr("x", -circleRadius * 0.7)
      .attr("y", -circleRadius * 0.7)
      .attr("width", circleRadius * 1.4)
      .attr("height", circleRadius * 1.4)
      .append("xhtml:div")
      .style("display", "flex")
      .style("align-items", "center")
      .style("justify-content", "center")
      .style("width", "100%")
      .style("height", "100%")
      .html(d => {
        const elementName = d.data.name;
        return `<img src="/icons/${elementName}.webp" alt="${elementName}" style="width:100%; height:100%; object-fit:contain;"/>`;
      });
      
    // Add text label below circle
    node.append("text")
      .attr("dy", circleRadius + 20)
      .attr("text-anchor", "middle")
      .text(d => d.data.name)
      .style("font-size", "12px")
      .attr("fill", "white")
      .style("filter", "drop-shadow(1px 1px 1px rgba(0, 0, 0, 0.8))");
      
    // Add legend for icons
    const legendG = container.append("g")
      .attr("transform", `translate(${margin.left}, ${height - margin.bottom + 50})`);
    
    // Title for legend
    legendG.append("text")
      .attr("x", 0)
      .attr("y", 0)
      .text("Icon Legend")
      .style("font-size", "18px")
      .style("font-weight", "bold")
      .attr("fill", "white");
      
    // Calculate number of icons per row
    const iconsPerRow = 8;
    const iconSize = 40;
    const iconMargin = 10;
    const rowHeight = 80;
    
    // Get unique elements to display in legend
    const uniqueItems = Array.from(new Set(uniqueElements));
    
    // Create legend items
    uniqueItems.forEach((item, i) => {
      const row = Math.floor(i / iconsPerRow);
      const col = i % iconsPerRow;
      const x = col * (iconSize + 60);
      const y = row * rowHeight + 30;
      
      // Icon background
      legendG.append("circle")
        .attr("cx", x + iconSize/2)
        .attr("cy", y + iconSize/2)
        .attr("r", iconSize/2)
        .attr("fill", "#677D6A")
        .attr("stroke", "white")
        .attr("stroke-width", 1);
        
      // Icon image
      legendG.append("foreignObject")
        .attr("x", x)
        .attr("y", y)
        .attr("width", iconSize)
        .attr("height", iconSize)
        .append("xhtml:div")
        .style("display", "flex")
        .style("align-items", "center")
        .style("justify-content", "center")
        .style("width", "100%")
        .style("height", "100%")
        .html(`<img src="/icons/${item}.webp" alt="${item}" style="width:100%; height:100%; object-fit:contain;"/>`);
        
      // Icon text
      legendG.append("text")
        .attr("x", x + iconSize/2)
        .attr("y", y + iconSize + 15)
        .text(item)
        .style("font-size", "12px")
        .attr("text-anchor", "middle")
        .attr("fill", "white")
        .style("filter", "drop-shadow(1px 1px 1px rgba(0, 0, 0, 0.8))");
    });
    
  }, [graph]);
  
  /* render list + tree */
  return (
    <div className="flex flex-col gap-4">
      <div className="border rounded p-4 w-full">
        <h3 className="font-semibold mb-2">Recipe steps</h3>
        <ul className="text-sm list-disc pl-5 space-y-1">
          {graph.recipes.map((r, i) => (
            <li key={i}>
              <span className="text-gray-600">{r.ingredients.join(" + ")}</span>{" "}
              ➜ <span className="font-medium">{r.result}</span>
            </li>
          ))}
        </ul>
      </div>
      {/* SVG buat tree */}
      <div className="border rounded p-4">
        <h3 className="font-semibold mb-2">Recipe Tree</h3>
        <div className="w-full overflow-auto max-h-screen">
          <svg ref={svgRef} style={{ width: "100%", height: "1100px" }}></svg>
        </div>
      </div>
    </div>
  );
};

export default RecipeResult;