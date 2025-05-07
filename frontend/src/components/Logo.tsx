import React from 'react';
import { Josefin_Sans } from "next/font/google";

const josefinSans = Josefin_Sans({
    subsets: ["latin"],
    weight: ["700"],
})

const Logo = () => {
    return (
        <div className="flex h-[93px] w-[600px] items-center">
            <div className="w-[100px] h-[100px] bg-[#d6bd98]" style={{ maskImage: 'url("/little-alchemy-2-icon.png")', WebkitMaskImage: 'url("/little-alchemy-2-icon.png")' }}></div>
            <h1
                className={`text-[#f3f3f3] ${josefinSans.className}`}
                style={{
                    textShadow: '-1px -1px 0 #000, 1px -1px 0 #000, -1px 1px 0 #000, 1px 1px 0 #000',
                }}
            >
                Home
            </h1>
        </div>
    );
};

export default Logo;