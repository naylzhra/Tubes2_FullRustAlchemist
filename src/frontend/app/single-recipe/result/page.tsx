"use client";
import { useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import RecipeResult from "../../_components/RecipeResult";

type ErrorResponse = {
  error: true;
  type: string;
  message: string;
};

type SuccessResponse = {
  error: false;
  data: GraphData;
};

type ApiResponse = ErrorResponse | SuccessResponse;

type GraphData = { 
  nodes: any[]; 
  recipes: any[]; 
  elapsed?: string; 
  visitedNodes?: number;
};

const Result = () => {
  const params = useSearchParams();
  const router = useRouter();
  const element = params.get("element") || "";
  const algo    = params.get("algo")?.toLowerCase()    || "bfs";

  const [data, setData]   = useState<GraphData | null>(null);
  const [error, setError] = useState<string | null>('error test');
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (!element) {
      setError("No element specified");
      setIsLoading(false);
      return;
    }

    (async () => {
      try {
        setIsLoading(true);
        const t0 = performance.now();
        const res = await fetch(
          `/api/recipe?element=${encodeURIComponent(element)}&algo=${algo}`
        );
        const json = await res.json() as ApiResponse;
        if (json.error) {
          const errorResponse = json as ErrorResponse;
          throw new Error(errorResponse.message || `Error: ${errorResponse.type}`);
        }

        const successResponse = json as SuccessResponse;
        const graphData = successResponse.data;

        graphData.elapsed = (performance.now() - t0).toFixed(2)
        setData(graphData);
        setError(null);

      } catch (e: any) {
        console.error("Error fetching recipe:", e);
        setError(e.message || "An unknown error occurred");
        setData(null);
      } finally {
        setIsLoading(false);
      }
    })();
  }, [element, algo]);

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
          className="m-[10px] p-[10px] w-[199px] h-[44px] border
                    border-[#d6bd98] rounded-[12px] bg-[#d6bd98] text-[#1E1E1E]"
          onClick={() => router.back()}
        >
          Back
        </button>
      </div>
    );
  }  
  
  return (
    <div className="flex flex-col items-center p-[2%]">
      <p className="w-[510px] h-[58px] m-[5px] p-4 border
                   border-[var(--foreground)] bg-[var(--foreground)]
                   rounded-[12px] text-white text-center">
        {element}
      </p>
      <div className="flex justify-between w-[510px] text-[#b3b3b3] m-[5px]">
        <p>Time execution: {data.elapsed} ms</p>
        <p>Visited nodes: {data?.visitedNodes}</p>
      </div>

      <RecipeResult graph={data} />

      <button
        className="m-[10px] p-[10px] w-[199px] h-[44px] border
                   border-[#d6bd98] rounded-[12px] bg-[#d6bd98]"
        onClick={() => router.back()}
      >
        Back
      </button>
    </div>
  );
}

export default Result;
