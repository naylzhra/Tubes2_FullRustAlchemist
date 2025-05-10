"use client";
import { useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import RecipeResult from "../_components/RecipeResult";

type GraphData = { nodes: any[]; recipes: any[]; elapsed?: string };

export default function Result() {
  const params = useSearchParams();
  const router = useRouter();
  const element = params.get("element") || "";
  const algo    = params.get("algo")    || "bfs";

  const [data, setData]   = useState<GraphData | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!element) return;

    (async () => {
      try {
        const t0 = performance.now();
        const res = await fetch(
          `/api/recipe?element=${encodeURIComponent(element)}&algo=${algo}`
        );
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const json: GraphData = await res.json();
        json.elapsed = (performance.now() - t0).toString();
        setData(json);
      } catch (e: any) {
        setError(e.message);
      }
    })();
  }, [element, algo]);

  if (error) return <p className="text-red-400">{error}</p>;
  if (!data)  return <p>Loading…</p>;

  return (
    <div className="flex flex-col items-center p-[2%]">
      <p className="w-[510px] h-[58px] m-[5px] p-[10px] border
                   border-[var(--foreground)] bg-[var(--foreground)]
                   rounded-[12px]">
        {element}
      </p>
      <div className="flex justify-between w-[510px] text-[#b3b3b3] m-[5px]">
        <p>Time execution: {data.elapsed} ms</p>
        <p>Visited nodes: {data.nodes.length}</p>
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
