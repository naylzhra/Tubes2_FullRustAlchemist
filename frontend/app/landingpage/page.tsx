import { Josefin_Sans } from "next/font/google";
import Footer from "../_components/Footer";

const josefinSans = Josefin_Sans({
  subsets: ["latin"],
  weight: ["700"],
});

const LandingPage = () => {
  return (
    <div className="max-h-screen flex flex-col bg-[var(--background)] overflow-y-hidden">
      {/* Konten Utama */}
      <div className="flex flex-1 items-center px-32 my-[86px]">
        {/* Kiri: Teks dan Tombol */}
        <div className="max-w-xl ml-[60px]">
          <h1 className={`leading-14 text-5xl text-white font-bold ${josefinSans.className}`}>
            Little<br />
            <span className="bg-gradient-to-br from-[#798772] to-[#D6BD98] bg-clip-text text-transparent">
              Alchemy
            </span>
            <br />
            Recipe Finder
          </h1>
          <div className="mt-8">
            <p className="font-semibold text-sm text-white">Looking for a recipe in Little Alchemy?</p>
          </div>
          <div className="mt-5">
            <p className="text-xs text-white">
              Just type the element you're curious about, and we'll show you all the possible<br />
              ways to make it — including the simplest, shortest recipe.
            </p>
          </div>
          <div className="mt-10 flex">
            <button className="py-3.5 px-6 bg-[#D6BD98] rounded-md text-xs font-semibold text-[#1A3636] transition duration-300 hover:scale-105 hover:shadow-[0_0px_15px_rgba(214,189,152,0.25)]">
            Shortest Recipe
            </button>
            <button className="ml-5 py-3.5 px-6 bg-[#D6BD98] rounded-md text-xs font-semibold text-[#1A3636] transition duration-300 hover:scale-105 hover:shadow-[0_0px_15px_rgba(214,189,152,0.25)]">
            Multiple Recipe
            </button>
          </div>
        </div>

        {/* Kanan: Gambar Beaver */}
        <div className="ml-[136px]">
          <img
            src="Beaver_2.svg"
            className="w-[364px] h-auto drop-shadow-[0px_7px_5px_rgba(0,0,0,0.5)]"
            alt="Beaver"
          />
        </div>
      </div>

      {/* Footer */}
      <Footer />
    </div>
  );
};

export default LandingPage;
