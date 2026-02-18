import { useEffect, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import apiClient from "../api/apiClient";
import { Play, Star, Plus } from "lucide-react";

const Dashboard = () => {
    const [popularMovies, setPopularMovies] = useState([]);
    const [loading, setLoading] = useState(true);
    const [currentHeroIndex, setCurrentHeroIndex] = useState(0);

    useEffect(() => {
        const fetchPopular = async () => {
            try {
                const response = await apiClient.get("/tmdb/popular");
                if (response.data.success) {
                    setPopularMovies(response.data.data.results);
                }
            } catch (error) {
                console.error("Failed to fetch popular movies", error);
            } finally {
                setLoading(false);
            }
        };

        fetchPopular();
    }, []);

    // Auto-rotate carousel
    useEffect(() => {
        if (popularMovies.length === 0) return;

        const interval = setInterval(() => {
            setCurrentHeroIndex((prev) => (prev + 1) % Math.min(5, popularMovies.length)); // Rotate through top 5
        }, 8000);

        return () => clearInterval(interval);
    }, [popularMovies, currentHeroIndex]);

    const currentMovie = popularMovies[currentHeroIndex];

    return (
        <div className="min-h-screen bg-[#12201B] text-white p-8">
            <header className="mb-10">
                <h1 className="text-3xl font-bold font-display mb-2">Popular on FrameRate</h1>
                <p className="text-gray-400">Discover what the community is watching.</p>
            </header>

            {loading ? (
                <div className="text-center py-20">Loading...</div>
            ) : (
                <div className="relative">
                    {/* hero carousel */}
                    {popularMovies.length > 0 && currentMovie && (
                        <div className="relative h-[500px] w-full rounded-2xl overflow-hidden mb-12 shadow-2xl group cursor-pointer">
                            <AnimatePresence mode="wait">
                                <motion.div
                                    key={currentMovie.id}
                                    initial={{ opacity: 0 }}
                                    animate={{ opacity: 1 }}
                                    exit={{ opacity: 0 }}
                                    transition={{ duration: 1 }}
                                    className="absolute inset-0 w-full h-full"
                                >
                                    <img
                                        src={`https://image.tmdb.org/t/p/original${currentMovie.backdrop_path}`}
                                        alt={currentMovie.title}
                                        className="w-full h-full object-cover transition-transform duration-[10s] ease-in-out transform scale-100 hover:scale-105"
                                        style={{ animation: 'slowZoom 10s linear infinite alternate' }}
                                    />
                                    <div className="absolute inset-0 bg-gradient-to-t from-[#0A1410] via-black/40 to-transparent" />

                                    <div className="absolute bottom-0 left-0 p-8 md:p-12 w-full md:w-2/3 z-10">
                                        <motion.h2
                                            initial={{ y: 20, opacity: 0 }}
                                            animate={{ y: 0, opacity: 1 }}
                                            transition={{ delay: 0.3 }}
                                            className="text-4xl md:text-6xl font-black font-display mb-4 leading-tight"
                                        >
                                            {currentMovie.title}
                                        </motion.h2>
                                        <motion.p
                                            initial={{ y: 20, opacity: 0 }}
                                            animate={{ y: 0, opacity: 1 }}
                                            transition={{ delay: 0.4 }}
                                            className="text-gray-200 line-clamp-2 text-lg mb-6"
                                        >
                                            {currentMovie.overview}
                                        </motion.p>

                                        <motion.div
                                            initial={{ y: 20, opacity: 0 }}
                                            animate={{ y: 0, opacity: 1 }}
                                            transition={{ delay: 0.5 }}
                                            className="flex items-center gap-4"
                                        >
                                            <button className="bg-white text-black px-6 py-3 rounded-full font-bold flex items-center gap-2 hover:bg-gray-200 transition-colors">
                                                <Play size={20} fill="black" /> Watch Trailer
                                            </button>
                                            <button className="bg-white/10 backdrop-blur-md px-6 py-3 rounded-full font-bold flex items-center gap-2 hover:bg-white/20 transition-colors">
                                                <Plus size={20} /> Add to Watchlist
                                            </button>
                                        </motion.div>
                                    </div>
                                </motion.div>
                            </AnimatePresence>

                            {/* Carousel Indicators */}
                            <div className="absolute bottom-8 right-8 flex gap-2 z-20">
                                {[...Array(Math.min(5, popularMovies.length))].map((_, idx) => (
                                    <button
                                        key={idx}
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            setCurrentHeroIndex(idx);
                                        }}
                                        className={`h-1 rounded-full transition-all duration-300 ${idx === currentHeroIndex ? "w-8 bg-mint" : "w-4 bg-white/30 hover:bg-white/50"
                                            }`}
                                    />
                                ))}
                            </div>
                        </div>
                    )}

                    {/* movie carousel */}
                    <h3 className="text-xl font-bold mb-6 font-display flex items-center gap-2">
                        <div className="w-1 h-6 bg-mint rounded-full"></div>
                        Trending Now
                    </h3>

                    <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-6">
                        {popularMovies.slice(1).map((movie) => (
                            <motion.div
                                key={movie.id}
                                initial={{ opacity: 0, y: 10 }}
                                whileInView={{ opacity: 1, y: 0 }}
                                viewport={{ once: true }}
                                className="group relative"
                            >
                                <div className="aspect-[2/3] rounded-xl overflow-hidden mb-3 shadow-lg relative">
                                    <img
                                        src={`https://image.tmdb.org/t/p/w500${movie.poster_path}`}
                                        alt={movie.title}
                                        className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
                                    />
                                    <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-4 backdrop-blur-sm">
                                        <button className="p-3 bg-mint rounded-full text-black hover:scale-110 transition-transform">
                                            <Star size={20} />
                                        </button>
                                        <button className="p-3 bg-white/20 rounded-full text-white hover:scale-110 transition-transform">
                                            <Plus size={20} />
                                        </button>
                                    </div>
                                    <div className="absolute top-2 right-2 bg-black/60 backdrop-blur-md px-2 py-1 rounded-md flex items-center gap-1 text-xs font-bold ring-1 ring-white/10">
                                        <Star size={12} className="text-yellow-400 fill-yellow-400" />
                                        {movie.vote_average.toFixed(1)}
                                    </div>
                                </div>
                                <h4 className="font-bold truncate group-hover:text-mint transition-colors">{movie.title}</h4>
                                <p className="text-xs text-gray-500">{new Date(movie.release_date).getFullYear()}</p>
                            </motion.div>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
};

export default Dashboard;