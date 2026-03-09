import { useEffect, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Link } from "react-router-dom";
import apiClient from "../api/apiClient";
import { getMovieVideos } from "../api/tmdb";
import { Play, Star, Plus, ChevronLeft, ChevronRight } from "lucide-react";
import useDocumentTitle from "../hooks/useDocumentTitle";
import MovieCard from "../components/MovieCard";

const HomePage = () => {
    const [trendingMovies, setTrendingMovies] = useState([]);
    const [loading, setLoading] = useState(true);
    const [currentHeroIndex, setCurrentHeroIndex] = useState(0);
    const [trailerKey, setTrailerKey] = useState(null);

    useDocumentTitle("Home");

    useEffect(() => {
        const fetchTrending = async () => {
            try {
                const response = await apiClient.get("/tmdb/trending?timeWindow=week");
                if (response.data.success) {
                    setTrendingMovies(response.data.data.results);
                }
            } catch (error) {
                console.error("Failed to fetch trending movies", error);
            } finally {
                setLoading(false);
            }
        };

        fetchTrending();
    }, []);

    const currentMovie = trendingMovies[currentHeroIndex];

    // recupere le trailer pour le film hero actuel
    useEffect(() => {
        if (!currentMovie) return;

        const fetchTrailer = async () => {
            try {
                const data = await getMovieVideos(currentMovie.id);
                if (data.success && data.data.results) {
                    const trailer = data.data.results.find(
                        (video) => video.type === "Trailer" && video.site === "YouTube"
                    );
                    setTrailerKey(trailer ? trailer.key : null);
                }
            } catch (error) {
                console.error("Failed to fetch trailer", error);
                setTrailerKey(null);
            }
        };

        fetchTrailer();
    }, [currentMovie]);

    // Auto-rotate carousel
    useEffect(() => {
        if (trendingMovies.length === 0) return;

        const interval = setInterval(() => {
            setCurrentHeroIndex((prev) => (prev + 1) % Math.min(5, trendingMovies.length)); // Rotate through top 5
        }, 8000);

        return () => clearInterval(interval);
    }, [trendingMovies, currentHeroIndex]);

    const handleWatchTrailer = (e) => {
        e.preventDefault();
        e.stopPropagation();
        if (trailerKey) {
            window.open(`https://www.youtube.com/watch?v=${trailerKey}`, "_blank");
        }
    };

    const handlePrev = (e) => {
        e.preventDefault();
        e.stopPropagation();
        setCurrentHeroIndex((prev) => (prev - 1 + Math.min(5, trendingMovies.length)) % Math.min(5, trendingMovies.length));
    };

    const handleNext = (e) => {
        e.preventDefault();
        e.stopPropagation();
        setCurrentHeroIndex((prev) => (prev + 1) % Math.min(5, trendingMovies.length));
    };

    return (
        <div className="min-h-screen bg-[var(--color-bg-primary)] text-[var(--color-text-primary)] p-8">
            <header className="mb-10">
                <h1 className="text-3xl font-bold font-display mb-2">Trending on FrameRate</h1>
                <p className="text-gray-400">Discover what the community is watching this week.</p>
            </header>

            {loading ? (
                <div className="flex justify-center py-20"><div className="loader"></div></div>
            ) : (
                <div className="relative">
                    {/* hero carousel */}
                    {trendingMovies.length > 0 && currentMovie && (
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
                                    <Link to={`/movie/${currentMovie.id}`} className="block w-full h-full">
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
                                                className="text-4xl md:text-6xl font-black font-display mb-4 leading-tight hover:text-mint transition-colors"
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
                                                className="flex flex-wrap items-center gap-3 md:gap-4"
                                            >
                                                <button
                                                    className={`bg-white text-black px-4 md:px-6 py-2 md:py-3 text-sm md:text-base rounded-full font-bold flex items-center gap-2 hover:bg-gray-200 transition-colors z-20 ${!trailerKey ? 'opacity-50 cursor-not-allowed' : ''}`}
                                                    onClick={handleWatchTrailer}
                                                    disabled={!trailerKey}
                                                >
                                                    <Play size={20} fill="black" /> Watch Trailer
                                                </button>
                                                <button className="bg-white/10 backdrop-blur-md px-4 md:px-6 py-2 md:py-3 text-sm md:text-base rounded-full font-bold flex items-center gap-2 hover:bg-white/20 transition-colors z-20">
                                                    <Plus size={20} /> Add to Watchlist
                                                </button>
                                            </motion.div>
                                        </div>
                                    </Link>
                                </motion.div>
                            </AnimatePresence>

                            {/* Carousel Controls */}
                            <button
                                onClick={handlePrev}
                                className="absolute left-4 top-1/2 -translate-y-1/2 p-2 rounded-full bg-black/30 text-white backdrop-blur-sm opacity-0 group-hover:opacity-100 transition-opacity hover:bg-[var(--color-primary)] hover:text-black z-20"
                            >
                                <ChevronLeft size={32} />
                            </button>
                            <button
                                onClick={handleNext}
                                className="absolute right-4 top-1/2 -translate-y-1/2 p-2 rounded-full bg-black/30 text-white backdrop-blur-sm opacity-0 group-hover:opacity-100 transition-opacity hover:bg-[var(--color-primary)] hover:text-black z-20"
                            >
                                <ChevronRight size={32} />
                            </button>

                            {/* Carousel Indicators */}
                            <div className="absolute bottom-8 right-8 flex gap-2 z-20">
                                {[...Array(Math.min(5, trendingMovies.length))].map((_, idx) => (
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
                        {trendingMovies.slice(1).map((movie) => (
                            <MovieCard key={movie.id} movie={movie} />
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
};

export default HomePage;