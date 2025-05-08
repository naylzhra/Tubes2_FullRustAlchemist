import { Josefin_Sans } from "next/font/google";

const josefinSans = Josefin_Sans({
  subsets: ["latin"],
  weight: ["700"],
});

const Footer = () => {
  return (
    <div className="w-full border-t border-[#b3b3b3] bg-[var(--background)] px-10 md:px-10 py-[15px] flex items-center">
      <p className={`ml-8 mr-6 font-semibold text-white ${josefinSans.className}`}>Contributors</p>
      
      <div className="flex space-x-10 text-sm text-white ml-8">
        <p>Ranashahira Reztaputri</p>
        <p>Syahrizal Bani Khairan</p>
        <p>Nayla Zahira</p>
      </div>
    </div>
  );
};

export default Footer;
