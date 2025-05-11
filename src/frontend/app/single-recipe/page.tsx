"use client";
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Josefin_Sans } from 'next/font/google';
import ScrollingElements from '../_components/ScrollingElements';

const josefinSans = Josefin_Sans({
  subsets: ['latin'],
  weight: ['700'],
});

const SingleRecipePage = () => {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedAlgo, setSelectedAlgo] = useState(1); // 0 = none, 1 = BFS, 2 = DFS
  const [errorMessage, setErrorMessage] = useState('');
  // Mock data for element grid
  const elements = Array(16).fill(null);

  const handleSearch = async () => {
    if (!searchQuery.trim()) {
      return;
    }
    setErrorMessage('');
    console.log(`Searching for: ${searchQuery} using ${selectedAlgo === 1 ? 'BFS' : 'DFS'}`);
    const algo = selectedAlgo === 1 ? 'BFS' : 'DFS';
    try {
      const response = await fetch(`/api/recipe?element=${encodeURIComponent(searchQuery.trim())}&algo=${algo}`);
      const data = await response.json();
      console.log(data);

      if (data.error) {
        setErrorMessage(data.message || `Error: ${data.type}`);
        return;
      }

      router.push(
      `/single-recipe/result?element=${encodeURIComponent(searchQuery.trim())}` +
      `&algo=${algo}`
      );
    } catch (error) {
      console.error("Error fetching recipe:", error);
    }
  };

  const handleInputChange = (e) => {
    setSearchQuery(e.target.value);
    setErrorMessage('');
  };

return (
  <div className="min-h-screen text-white p-8 overflow-y-hidden">
    {/* Title */}
    <div className="mt-4 text-center items-center">
      <h1 className={`text-6xl font-bold text-white ${josefinSans.className}`}>
        Little <span className="bg-gradient-to-br from-purple-[#798772] to-[#D6BD98] bg-clip-text text-transparent">Alchemy</span> Recipe
      </h1>
    </div>
    <div className="mt-10 flex flex-col items-center mb-10">
      <div className="flex justify-center h-10 space-x-3">
        <div className="flex items-center">
          <select 
            value={selectedAlgo}
            onChange={(e) => setSelectedAlgo(Number(e.target.value))}
            className="select-box flex h-full align-middle bg-[#D6BD98] text-[#1E1E1E] text-center items-center rounded-sm px-2">
            <option value="1">BFS</option>
            <option value="2">DFS</option>
          </select>
        </div>
        <div className="flex bg-[#40534C] h-full w-96 align-middle text-center text-white items-center rounded-sm">
          <input
            type="text"
            placeholder="Which element recipe are you looking for?"
            value={searchQuery}
            onChange={handleInputChange}
            className="w-full h-full bg-transparent text-white text-center placeholder-[#B3B3B3] px-4"
          />
        </div>
        <button 
          onClick={handleSearch}
          disabled={!searchQuery.trim()}
          className={`rounded-sm w-20 ${
            !searchQuery.trim() 
              ? 'bg-[#d6bd9877] text-[#1E1E1E] cursor-not-allowed' 
              : 'bg-[#D6BD98] text-[#1E1E1E] hover:text-[#40534C] hover:bg-white'
          }`}
        >
          Search
        </button>
      </div>
      {errorMessage && (
        <div className="p-1 bg-opacity-20 rounded-md text-center max-w-md">
          <p className="text-[#B3B3B3]">{errorMessage}</p>
        </div>
      )}
    </div>
    <ScrollingElements />
  </div>
  );
} 

export default SingleRecipePage;