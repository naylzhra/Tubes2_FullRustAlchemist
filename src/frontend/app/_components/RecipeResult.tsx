"use client";
import React, { useEffect, useRef } from "react";
import * as d3 from "d3";

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

  useEffect(() => {
    if (!graph.recipes || graph.recipes.length === 0) return;
    
    // Increase width and height for better visualization
    const width = 1600;
    const height = 1000;
    if (!svgRef.current) return;
    
    // Clear existing SVG content
    const svg = d3.select(svgRef.current)
      .attr("viewBox", `0 0 ${width} ${height}`)
      .selectAll("*").remove();
      
    // Create a fresh SVG container
    const container = d3.select(svgRef.current)
      .attr("viewBox", `0 0 ${width} ${height}`);
    
    // Define larger margins to allow more space
    const margin = { top: 150, right: 200, bottom: 150, left: 200 };
    
    // Create main group with translation for margins
    const g = container.append("g")
      .attr("transform", `translate(${margin.left}, ${margin.top})`);
    
    const rootData = buildTree(graph.recipes[0]?.result ?? "", graph.recipes);
    const root = d3.hierarchy(rootData);
    
    // Adjust tree layout with proper dimensions accounting for margins
    const treeLayout = d3.tree<any>()
      .size([width - margin.left - margin.right, height - margin.top - margin.bottom])
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
    
    // Node - add to the translated group with larger nodes
    const node = g.append("g")
      .selectAll("g")
      .data(root.descendants())
      .join("g")
      .attr("transform", d => `translate(${d.x},${d.y})`);
    
    // Increase node sizes for better visibility
    const rectWidth = 100;
    const rectHeight = 45;
    
    node.append("rect")
      .attr("x", -rectWidth / 2)
      .attr("y", -rectHeight / 2)
      .attr("width", rectWidth)
      .attr("height", rectHeight)
      .attr("fill", "#677D6A")
      .attr("rx", 10)
      .attr("ry", 10);
    
    node.append("text")
      .attr("dy", ".35em")
      .attr("text-anchor", "middle")
      .text(d => d.data.name)
      .style("font-size", "14px") // Increased font size
      .attr("fill", "white");
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
          <svg ref={svgRef} style={{ width: "100%", height: "900px" }}></svg>
        </div>
      </div>
    </div>
  );
};

export default RecipeResult;