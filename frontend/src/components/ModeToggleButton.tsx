import React, { useState } from "react";


const ModeToggleButton = () => {
    // 0: Unselected, 1: Shortest recipe, 2: Multiple recipe
    let [choice, setChoice] = useState<number>(1);

    const handleChoiceChange = (newChoice: number) => {
        // Should be static on result page
        setChoice(newChoice);
    }

    const unselectedBackground = "grid grid-cols-2 items-center justify-center mx-[10px] p-0 h-[60%] w-[clamp(470px,20%,570px)] rounded-[200px] border-2 border-[var(--foreground)] bg-[var(--foreground)]";
    const background = choice === 0
        ? unselectedBackground
        : `${unselectedBackground} rounded-[200px] border-2 border-[var(--foreground)] bg-[var(--foreground)]`;
    const highlight = "h-[85%] w-[95%] bg-[#FFFFFF] mx-[5px] text-[var(--foreground)] rounded-[200px] flex justify-center items-center";
    const highlightShortest = choice === 1 ? highlight : "flex justify-center align-center";
    const highlightMultiple = choice === 2 ? highlight : "flex justify-center align-center";

    return (
        <div className={background}>
            <p className={highlightShortest} onClick={()=> handleChoiceChange(1)}>Shortest recipe</p>
            <p className={highlightMultiple} onClick={()=> handleChoiceChange(2)}>Multiple recipe</p>
        </div>
    );
};

export default ModeToggleButton;