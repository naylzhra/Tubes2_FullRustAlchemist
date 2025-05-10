"use client";
import { useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import RecipeResult from "../../_components/RecipeResult";

type GraphNode   = { id: number; name: string };
type GraphRecipe = { ingredients: string[]; result: string; step: number };
interface GraphData { nodes: GraphNode[]; recipes: GraphRecipe[] }

const MultiResult = () => {
  const params  = useSearchParams();
  const router  = useRouter();

  const element = params.get("element") || "";
  const algo    = params.get("algo")    || "bfs";
  const max     = params.get("max")     || "5";

  const [paths, setPaths] = useState<GraphData[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!element) return;

    (async () => {
      try {
        const res = await fetch(
          `/api/recipes?element=${encodeURIComponent(element)}` +
          `&algo=${algo}&max=${max}`
        );
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const json = await res.json();
        setPaths(json.paths as GraphData[]);
      } catch (e: any) {
        setError(e.message);
      }
    })();
  }, [element, algo, max]);

  if (error)   return <p className="text-red-400">{error}</p>;
  if (!paths.length) return <p className="text-white">Loading…</p>;

  return (
    <div className="min-h-screen text-white p-8">
      <div className="flex flex-col items-center">
        <h2 className="text-3xl font-semibold mb-4">
          {paths.length} recipe path{paths.length > 1 && "s"} for&nbsp;
          <span className="text-[#d6bd98]">{element}</span> ({algo.toUpperCase()})
        </h2>

        {paths.map((p, i) => (
          <div key={i} className="mb-10">
            <h3 className="text-xl mb-2">Path #{i + 1}</h3>
            <RecipeResult graph={p} />
          </div>
        ))}

        <button
          onClick={() => router.back()}
          className="mt-6 px-6 py-2 bg-[#d6bd98] rounded text-[#1e1e1e]"
        >
          Back
        </button>
      </div>
    </div>
  );
}
export default MultiResult;
