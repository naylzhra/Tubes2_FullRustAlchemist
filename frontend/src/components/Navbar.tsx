import React from 'react';
import Logo from "./Logo";
import ModeToggleButton from "./ModeToggleButton";

export const Navbar = () => {
return (
    <>
    <div className="navbar flex p-[10px] justify-between items-center border-b-[0.2px] border-[#b3b3b3]-300 h-[93px]">
        <Logo />
        <ModeToggleButton />
    </div>
    </>
);
};

export default Navbar;