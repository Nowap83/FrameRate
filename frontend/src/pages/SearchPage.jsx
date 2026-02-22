import { useEffect, useState } from "react";
import { useSearchParams, Link } from "react-router-dom";
import apiClient from "../api/apiClient";
import useDocumentTitle from "../hooks/useDocumentTitle";

const SearchPage = () => {
    const [searchParams] = useSearchParams();
    const query = searchParams.get("q");
    const [results, setResults] = useState([]);
    const [loading, setLoading] = useState(false);

    useDocumentTitle(query ? `Search: ${query}` : "Search");

    useEffect(() => {
        const fetchResults = async () => {
            if (!query) return;
            setLoading(true);
            try {
                const response = await apiClient.get(`/tmdb/search?q=${encodeURIComponent(query)}`);
                setResults(response.data.data.results || []);
            } catch (error) {
                console.error("Search failed:", error);
                setResults([]);
            } finally {
                setLoading(false);
            }
        };

        fetchResults();
    }, [query]);

    return (
        <div className="min-h-screen bg-[var(--color-body-bg)] pt-24 pb-20 px-4 md:px-8">
            <div className="max-w-7xl mx-auto">
                <h1 className="text-3xl font-bold text-white mb-8 border-b border-white/10 pb-4">
                    Search Results for <span className="text-[var(--color-primary)]">"{query}"</span>
                </h1>

                {loading ? (
                    <div className="text-white text-center py-20 animate-pulse">Loading results...</div>
                ) : results.length > 0 ? (
                    <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-6">
                        {results.map((movie) => (
                            <Link key={movie.id} to={`/movie/${movie.id}`} className="group relative block bg-[#1a1a1a] rounded-xl overflow-hidden shadow-lg hover:shadow-[0_0_20px_rgba(0,0,0,0.5)] transition-all hover:-translate-y-1">
                                <div className="aspect-[2/3] overflow-hidden">
                                    {movie.poster_path ? (
                                        <img
                                            src={`https://image.tmdb.org/t/p/w500${movie.poster_path}`}
                                            alt={movie.title}
                                            className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
                                        />
                                    ) : (
                                        <div className="w-full h-full bg-white/5 flex items-center justify-center text-gray-500">
                                            No Image
                                        </div>
                                    )}
                                </div>
                                <div className="p-4">
                                    <h3 className="font-bold text-white text-sm line-clamp-2 mb-1 group-hover:text-[var(--color-primary)] transition-colors">
                                        {movie.title}
                                    </h3>
                                    <p className="text-xs text-gray-400">
                                        {movie.release_date ? movie.release_date.split('-')[0] : 'Unknown'}
                                    </p>
                                </div>
                            </Link>
                        ))}
                    </div>
                ) : (
                    <div className="text-center py-20">
                        <p className="text-xl text-gray-400">No results found for "{query}".</p>
                        <p className="text-sm text-gray-500 mt-2">Try checking for typos or using broader keywords.</p>
                    </div>
                )}
            </div>
        </div>
    );
};

export default SearchPage;
