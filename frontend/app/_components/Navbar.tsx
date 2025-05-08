"use client";

import React, { useState } from "react";
import { Josefin_Sans } from "next/font/google";

const josefinSans = Josefin_Sans({
  subsets: ["latin"],
  weight: ["700"],
});

const Navbar = () => {
  // 0 = none, 1 = Shortest, 2 = Multiple
  const [choice, setChoice] = useState<number>(1);

  /* ----- styling helpers for the pill ----- */
  const wrapper =
    "grid grid-cols-2 items-center justify-center mx-[10px] p-0 h-[60%] " +
    "w-[clamp(470px,20%,570px)] rounded-[200px] border-2 border-[var(--foreground)] " +
    "bg-[var(--foreground)]";
  const highlight =
    "h-[85%] w-[95%] bg-[#FFFFFF] mx-[5px] text-[var(--foreground)] " +
    "rounded-[200px] flex justify-center items-center";
  const maybeHighlight = (n: number) =>
    choice === n ? highlight : "flex justify-center items-center";

  return (
    <div className="sticky top-0 z-50 h-auto w-full border-b border-[#b3b3b3] bg-[var(--background)] px-10 md:px-10 flex items-center justify-between">
        {/* ---- logo area ---- */}
        <div className="flex items-center gap-x-1 px-11">
            <div
                className="w-[50px] h-[50px] bg-[#d6bd98]"
                style={{
                    maskImage: 'url("/little-alchemy-2-icon.png")',
                    WebkitMaskImage: 'url("/little-alchemy-2-icon.png")',
                    maskSize: "50px 50px",
                    WebkitMaskSize: "50px 50px",
                    maskRepeat: "no-repeat",
                    WebkitMaskRepeat: "no-repeat",
                    maskPosition: "center",
                    WebkitMaskPosition: "center",
                }}
            />
            <text
                className={`text-[#f3f3f3] text-lg leading-none ${josefinSans.className}`}
                // style={{
                //     textShadow: "0 0 5px #000000, 0 0 10px #000000",
                // }}
            >
                Home
            </text>
        </div>
        
        {/* ---- toggle pill ---- */}
        <div className={wrapper}>
            <p className={maybeHighlight(1)} onClick={() => setChoice(1)}>
                Shortest recipe
            </p>
            <p className={maybeHighlight(2)} onClick={() => setChoice(2)}>
                Multiple recipe
            </p>
        </div>
    </div>
  );
};

export default Navbar;