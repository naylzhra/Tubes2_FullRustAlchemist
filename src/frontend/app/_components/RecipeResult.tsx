"use client";
import React, { useEffect, useRef } from "react";
import * as d3 from "d3";
//import { HierarchyPointNode, HierarchyPointLink } from "d3-hierarchy";

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
    // Increase width and height for better visualization
    const width = 800;
    const height = 600;

    if (!svgRef.current) return;

    const svg = d3.select(svgRef.current)
      .attr("viewBox", `0 0 ${width} ${height}`)

    const rootData = buildTree(graph.recipes[0]?.result ?? "", graph.recipes);
    const root = d3.hierarchy(rootData);

    // Adjust tree layout size and node separation
    const treeLayout = d3.tree<any>()
      .size([width - 120, height - 120])
      .separation((a, b) => (a.parent === b.parent ? 1.2 : 2)); // Reduce separation between nodes

    treeLayout(root);

    // Add margins by translating the entire visualization
    const g = svg.selectAll("*").remove() // Bersihkan isi svg sebelum gambar baru
      .append("g")
      .attr("transform", `translate(60, 60)`); // Add margins

    // Garis antar node
    const linkGenerator = d3.linkVertical<any, any>()
      .x((d: any) => d.x)
      .y((d: any) => d.y);

    svg.append("g")
      .selectAll("path")
      .data(root.links())
      .join("path")
      .attr("fill", "none")
      .attr("stroke", "#555")
      .attr("stroke-width", 2)
      .attr("d", d => linkGenerator(d));

    // Node
    const node = svg.append("g")
      .selectAll("g")
      .data(root.descendants())
      .join("g")
      .attr("transform", d => `translate(${d.x},${d.y})`);

    const rectWidth = 60;
    const rectHeight = 30;

    node.append("rect")
      .attr("x", -rectWidth / 2)
      .attr("y", -rectHeight / 2)
      .attr("width", rectWidth)
      .attr("height", rectHeight)
      .attr("fill", "#677D6A")
      .attr("rx", 10) // untuk sudut membulat, opsional
      .attr("ry", 10); // untuk sudut membulat, opsional

    node.append("text")
      .attr("dy", ".35em")
      .attr("text-anchor", "middle")
      .text(d => d.data.name)
      .style("font-size", "px");
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
        <svg ref={svgRef} style={{ width: "100%", height: "500px" }}></svg>
      </div>
    </div>
  );
};

export default RecipeResult;