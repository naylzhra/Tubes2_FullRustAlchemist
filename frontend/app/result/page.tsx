import React from "react";
import RecipeResult from "../_components/RecipeResult";

const Result = () => {
    return (
        <div>
            <div className="flex flex-col items-center justify-center p-[2%]">
                <p className="w-[510px] h-[58px] m-[5px] p-[10px] border border-[var(--foreground)] bg-[var(--foreground)] rounded-[12px] text-left align-middle">Acid Rain</p>
                <div className="flex justify-between w-[510px] text-[#b3b3b3] m-[5px]">
                    <p>Time execution: 1.35s</p>
                    <p>Visited nodes: 100</p>
                </div>
                <RecipeResult />
                <button className="m-[10px] p-[10px] w-[199px] h-[44px] border border-[#d6bd98] rounded-[12px] bg-[#d6bd98] text-[#000000] text-[20px] text-center align-middle">Back</button>
            </div>
        </div>
    );
};


export default Result;