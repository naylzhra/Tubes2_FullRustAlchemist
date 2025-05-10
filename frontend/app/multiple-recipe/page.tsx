"use client";
import { useState } from 'react';
import { Josefin_Sans } from 'next/font/google';

const josefinSans = Josefin_Sans({
  subsets: ['latin'],
  weight: ['700'],
});

const MultipleRecipePage = () => {
const [searchQuery, setSearchQuery] = useState('');
const [selectedAlgo, setSelectedAlgo] = useState(1); // 0 = none, 1 = BFS, 2 = DFS
const [selectedNumR, setSelectedNumR] = useState(5); //num of recipes

// Mock data for element grid
const elements = Array(16).fill(null);

const handleSearch = () => {
  if (!searchQuery.trim()) {
    return;
  }
  console.log(`Searching for: ${searchQuery} using ${selectedAlgo === 1 ? 'BFS' : 'DFS'}`);
  // Implement actual search functionality here
};

const handleInputQuery = (e) => {
  setSearchQuery(e.target.value);
};

const handleInputNumR = (e) => {
  setSelectedNumR(Number(e.target.value));
};

return (
  <div className="min-h-screen text-white p-8">
    {/* Title */}
    <div className="mt-4 text-center items-center">
      <h1 className={`text-6xl font-bold text-white ${josefinSans.className}`}>
        Little <span className="bg-gradient-to-br from-purple-[#798772] to-[#D6BD98] bg-clip-text text-transparent">Alchemy</span> Recipe
      </h1>
    </div>
    <div className="mt-3 flex flex-col items-center h-10">
      <div className="flex justify-center">
        <div className="text-white text-sm">
          Enter maximum number of recipes to show:
        </div>
        <input
            type="text"
            value={selectedNumR}
            onChange={(e) => setSelectedNumR(Number(e.target.value))}
            className="ml-3 w-10 h-4.5 bg-[#40534C] rounded text-white text-center"
        />
      </div>
    </div>
    <div className="flex flex-col items-center">
      <div className="flex justify-center h-10 space-x-3 gap-3">
        <div className="flex items-center">
          <select 
            value={selectedAlgo}
            onChange={(e) => setSelectedAlgo(Number(e.target.value))}
            className="select-box flex h-full align-middle bg-[#D6BD98] text-[#1E1E1E] text-center items-center rounded-sm px-3 py-1">
            <option value="1">BFS</option>
            <option value="2">DFS</option>
          </select>
        </div>
        <div className="flex bg-[#40534C] h-full w-96 align-middle text-center text-white items-center rounded-sm px-2">
          <input
            type="text"
            placeholder="Which element recipe are you looking for?"
            value={searchQuery}
            onChange={handleInputQuery}
            className="w-full h-full bg-transparent text-center text-white placeholder-[#B3B3B3] px-4"
          />
        </div>
        <button 
          onClick={handleSearch}
          disabled={!searchQuery.trim()}
          className={`px-4 rounded-sm w-20 ${
            !searchQuery.trim() 
              ? 'bg-[#d6bd9877] text-[#1E1E1E] cursor-not-allowed' 
              : 'bg-[#D6BD98] text-[#1E1E1E] hover:bg-amber-300'
          }`}
        >
          Search
        </button>
      </div>  
    </div>
  </div>
  );
} 

export default MultipleRecipePage;