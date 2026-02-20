import { useState, useEffect, useRef } from 'react';
import { Search } from 'lucide-react';
import apiClient from '../api/apiClient';

const MovieSearch = ({ onSelect, onCancel, placeholder = "Search for a movie..." }) => {
    const [query, setQuery] = useState('');
    const [results, setResults] = useState([]);
    const [loading, setLoading] = useState(false);
    const searchContainerRef = useRef(null);

    useEffect(() => {
        const searchMovies = async () => {
            if (query.length > 2) {
                setLoading(true);
                try {
                    const response = await apiClient.get(`/tmdb/search?q=${encodeURIComponent(query)}`);
                    setResults(response.data.data.results || []);
                } catch (error) {
                    console.error("Search failed", error);
                    setResults([]);
                } finally {
                    setLoading(false);
                }
            } else {
                setResults([]);
            }
        };

        const timeoutId = setTimeout(searchMovies, 300); // debounce
        return () => clearTimeout(timeoutId);
    }, [query]);

    useEffect(() => {
        const handleClickOutside = (event) => {
            if (searchContainerRef.current && !searchContainerRef.current.contains(event.target)) {
                setQuery('');
                if (onCancel) onCancel();
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, [onCancel]);


    return (
        <div ref={searchContainerRef} className="relative w-full">
            <div className="relative flex items-center">
                <span className="absolute left-3 text-gray-400">
                    <Search size={16} />
                </span>
                <input
                    type="text"
                    placeholder={placeholder}
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    className="w-full bg-white/5 text-sm text-gray-300 py-2 pl-10 pr-4 rounded-full border border-white/10 focus:outline-none focus:border-[var(--color-primary)]/50 transition-colors"
                />
                {query && (
                    <button
                        onClick={() => { setQuery(''); onCancel && onCancel(); }}
                        className="absolute right-3 text-gray-500 hover:text-white"
                    >
                        <Search size={14} className="rotate-45" /> {/* Close icon visual using Search rotated or just X */}
                    </button>
                )}
            </div>

            {/* Results Dropdown */}
            {(results.length > 0 || (query.length > 2 && loading)) && (
                <div className="absolute top-full left-0 right-0 mt-2 bg-[#1a1a1a] rounded-xl border border-white/10 shadow-2xl overflow-hidden z-50 animate-in fade-in slide-in-from-top-2">
                    <div className="max-h-80 overflow-y-auto custom-scrollbar">
                        {results.map(movie => (
                            <button
                                key={movie.id}
                                type="button"
                                onClick={() => {
                                    onSelect && onSelect(movie);
                                    setQuery(''); // Clear on select
                                    setResults([]);
                                }}
                                className="w-full flex items-center gap-3 p-3 hover:bg-white/5 transition-colors text-left border-b border-white/5 last:border-0"
                            >
                                <div className="w-10 h-14 bg-gray-800 rounded overflow-hidden flex-shrink-0 shadow-sm">
                                    {movie.poster_path ? (
                                        <img src={`https://image.tmdb.org/t/p/w92${movie.poster_path}`} alt={movie.title} className="w-full h-full object-cover" />
                                    ) : (
                                        <div className="w-full h-full bg-white/10" />
                                    )}
                                </div>
                                <div>
                                    <p className="font-bold text-white text-sm line-clamp-1">{movie.title}</p>
                                    <p className="text-xs text-gray-400">
                                        {movie.release_date ? movie.release_date.split('-')[0] : 'Unknown'}
                                    </p>
                                </div>
                            </button>
                        ))}

                        {loading && (
                            <div className="p-4 text-center text-gray-500 text-sm">
                                Searching...
                            </div>
                        )}

                        {!loading && results.length === 0 && query.length > 2 && (
                            <div className="p-4 text-center text-gray-500 text-sm">
                                No results found.
                            </div>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
};

export default MovieSearch;
