"use client";
import { useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import RecipeResult from "../../_components/RecipeResult";
import Navbar from "../../_components/Navbar";
import config from "@/config";

type GraphNode   = { id: number; name: string };
type GraphRecipe = { ingredients: string[]; result: string; step: number };
interface GraphData { nodes: GraphNode[]; recipes: GraphRecipe[] }

type ErrorResponse = { error: true;  type: string; message: string };
type SuccessResponseMR = {
  data: { algo: string; element: string; paths: GraphData[]; visitedNodes: number };
  error: false;
};
type ApiResponse = ErrorResponse | SuccessResponseMR;

const MultiResult = () => {
  const params  = useSearchParams();
  const router  = useRouter();

  const element = params.get("element") || "";
  const algo    = params.get("algo")    || "bfs";
  const max     = params.get("max")     || "5";
  const [mode, setMode] = useState(2);

  const [paths,        setPaths]        = useState<GraphData[]>([]);
  const [error,        setError]        = useState<string | null>(null);
  const [isLoading,    setIsLoading]    = useState<boolean>(true);
  const [elapsed,      setElapsed]      = useState<number | null>(null)
  const [visited,      setVisited]      = useState<number | null>(null)

  useEffect(() => {
    if (!element) {
      setError("No element specified");
      setIsLoading(false);
      return;
    }

    (async () => {
      try {
        setIsLoading(true);
        const t0  = performance.now();
        const res = await fetch(
          `${config.API_URL}/api/recipes?element=${encodeURIComponent(element)}` +
          `&algo=${algo}&max=${max}`
        );
        const t1  = performance.now();

        const json = (await res.json()) as ApiResponse;

        if (json.error) {
          const { message, type } = json as ErrorResponse;
          throw new Error(message || `Error: ${type}`);
        }

        const { paths, visitedNodes } = (json as SuccessResponseMR).data;
        setPaths(paths);
        setElapsed(Math.round(t1 - t0));
        setVisited(visitedNodes);        
        setError(null);
      } catch (e: any) {
        setError(e.message);
        setPaths([]);
        setElapsed(null);
        setVisited(null);
      } finally {
        setIsLoading(false);
      }
    })();
  }, [element, algo, max]);

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-[50vh]">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-[#D6BD98]" />
      </div>
    );
  }
  
  if (error) {
    return (
      <div className="max-h-screen flex flex-col bg-[var(--background)]">
        <Navbar variant="multiple" currentRecipeMode={mode} setRecipeMode={setMode} />
        <div className="flex flex-col items-center p-[2%]">
          <p className="text-white mt-10 text-center mb-6">{error}</p>
          <button className="px-6 py-2 bg-[#d6bd98] rounded-md text-[#1e1e1e]"
                  onClick={() => router.back()}>
            Back
          </button>
        </div>
      </div>
    );
  }
  
  if (paths.length === 0) {
    return (
      <div className="max-h-screen flex flex-col bg-[var(--background)]">
        <Navbar variant="multiple" currentRecipeMode={mode} setRecipeMode={setMode} />
        <div className="text-white p-8 flex flex-col items-center">
          <h2 className="text-3xl font-semibold mb-4">
            No recipes found for <span className="text-[#d6bd98]">{element}</span>
          </h2>
          <p className="text-gray-400 mb-6">Try a different element or algorithm.</p>
          <button className="px-6 py-2 bg-[#d6bd98] rounded text-[#1e1e1e] transition duration-200 hover:scale-105"
                  onClick={() => router.back()}>
            Back
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="max-h-screen flex flex-col bg-[var(--background)]">
      <Navbar variant="multiple" currentRecipeMode={mode} setRecipeMode={setMode} />
      <div className="flex flex-col items-center p-[2%]">
        <p className="w-[510px] h-[58px] m-[5px] p-4 border
                    border-[var(--foreground)] bg-[var(--foreground)]
                    rounded-[12px] text-white text-center">
          {element}
        </p>

        <div className="flex justify-between w-[510px] text-[#b3b3b3] m-[5px]">
          <p>Execution time:&nbsp;{elapsed ?? "--"}&nbsp;ms</p>
          <p>Visited nodes:&nbsp;{visited ?? "--"}</p>
        </div>

        {paths.map((p, i) => (
          <div key={i} className="mb-10">
            <h3 className="text-xl mb-2">Path #{i + 1}</h3>
            <RecipeResult graph={p} />
          </div>
        ))}

        <button className="mt-6 px-6 py-2 bg-[#d6bd98] rounded text-[#1e1e1e]"
                onClick={() => router.back()}>
          Back
        </button>
      </div>
    </div>
  );
};

export default MultiResult;
