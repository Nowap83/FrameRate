import React from "react";

const SearchBar = ({ className = "" }) => {
    return (
        <div className={`relative flex items-center ${className}`}>
            <span className="absolute left-3 text-gray-400">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
            </span>
            <input
                type="text"
                placeholder="Search movies, actors, directors..."
                className="w-full bg-white/5 text-sm text-gray-300 py-1.5 pl-10 pr-4 rounded-md border border-white/10 focus:outline-none focus:border-mint/50 transition-colors"
            />
        </div>
    );
};

export default SearchBar;
