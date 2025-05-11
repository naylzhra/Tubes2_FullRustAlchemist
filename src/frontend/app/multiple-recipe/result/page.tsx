"use client";
import { useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import RecipeResult from "../../_components/RecipeResult";

type GraphNode   = { id: number; name: string };
type GraphRecipe = { ingredients: string[]; result: string; step: number };
interface GraphData { nodes: GraphNode[]; recipes: GraphRecipe[] }

type ErrorResponse = {
  error: true;
  type: string;
  message: string;
};

type SuccessResponse = {
  error: false;
  data: {
    element: string;
    algo: string;
    paths: GraphData[];
  };
};

type ApiResponse = ErrorResponse | SuccessResponse;

const MultiResult = () => {
  const params  = useSearchParams();
  const router  = useRouter();

  const element = params.get("element") || "";
  const algo    = params.get("algo")    || "bfs";
  const max     = params.get("max")     || "5";

  const [paths, setPaths] = useState<GraphData[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  useEffect(() => {
    if (!element) {
      setError("No element specified");
      setIsLoading(false);
      return;
    }

    (async () => {
      try {
        setIsLoading(true);
        const res = await fetch(
          `/api/recipes?element=${encodeURIComponent(element)}` +
          `&algo=${algo}&max=${max}`
        );
        
        const json = await res.json() as ApiResponse;
        
        if (json.error) {
          const errorResponse = json as ErrorResponse;
          throw new Error(errorResponse.message || `Error: ${errorResponse.type}`);
        }
        
        const successResponse = json as SuccessResponse;
        setPaths(successResponse.data.paths);
        setError(null);
      } catch (e: any) {
        setError(e.message);
        setPaths([]);
      } finally {
        setIsLoading(false);
      }
    })();
  }, [element, algo, max]);
  
  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-[50vh]">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-[#D6BD98]"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center p-[2%]">
        <div className="bg-red-500 bg-opacity-20 border border-red-500 rounded-md p-4 max-w-md w-full mb-6">
          <p className="text-red-300 text-center">{error}</p>
        </div>
        
        <button
          className="mt-6 px-6 py-2 bg-[#d6bd98] rounded text-[#1e1e1e]"
          onClick={() => router.back()}
        >
          Go Back
        </button>
      </div>
    );
  }
  
  if (!paths || paths.length === 0) {
    return (
      <div className="min-h-screen text-white p-8">
        <div className="flex flex-col items-center">
          <h2 className="text-3xl font-semibold mb-4">
            No recipes found for <span className="text-[#d6bd98]">{element}</span>
          </h2>
          <p className="text-amber-300 mb-6">Try a different element or algorithm.</p>
          
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
