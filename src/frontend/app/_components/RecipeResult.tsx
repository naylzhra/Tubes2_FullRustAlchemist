import React from "react";

/* ------------  type definitions  ------------ */
export type GraphNode   = { id: number; name: string };
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
  /* render however you like */
  return (
    <div className="border rounded p-4 w-[510px]">
      <h3 className="font-semibold mb-2">Recipe steps</h3>
      <ul className="text-sm list-disc pl-5 space-y-1">
        {graph.recipes.map((r, i) => (
          <li key={i}>
            <span className="text-[#d6bd98]">{r.ingredients.join(" + ")}</span>{" "}
            ➜ <span className="text-white">{r.result}</span>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default RecipeResult;

// const RecipeResult = () => {
//     return (
//         <img src="graph.png" alt="Graph" className="p-[50px] m-[20px] w-[80%] h-[500px] bg-[var(--foreground)]"></img>
//     );
// };

// export default RecipeResult;