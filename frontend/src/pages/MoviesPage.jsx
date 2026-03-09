import { useEffect, useState, useRef } from "react";
import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { ChevronLeft, ChevronRight } from "lucide-react";
import apiClient from "../api/apiClient";
import useDocumentTitle from "../hooks/useDocumentTitle";
import MovieCard from "../components/MovieCard";

const MoviesPage = () => {
    const [categories, setCategories] = useState({
        popular: [],
        topRated: [],
        upcoming: [],
    });
    const [loading, setLoading] = useState(true);

    useDocumentTitle("Explore Movies");

    useEffect(() => {
        const fetchAllMovies = async () => {
            setLoading(true);
            try {
                const [popularRes, topRatedRes, upcomingRes] = await Promise.all([
                    apiClient.get("/tmdb/popular"),
                    apiClient.get("/tmdb/top-rated"),
                    apiClient.get("/tmdb/upcoming"),
                ]);

                setCategories({
                    popular: popularRes.data.success ? popularRes.data.data.results : [],
                    topRated: topRatedRes.data.success ? topRatedRes.data.data.results : [],
                    upcoming: upcomingRes.data.success ? upcomingRes.data.data.results : [],
                });
            } catch (error) {
                console.error("Failed to fetch movies collections:", error);
            } finally {
                setLoading(false);
            }
        };

        fetchAllMovies();
    }, []);

    const MovieCarousel = ({ title, movies }) => {
        const carouselRef = useRef(null);

        if (!movies || movies.length === 0) return null;

        const scroll = (direction) => {
            if (carouselRef.current) {
                const { scrollLeft, clientWidth } = carouselRef.current;
                const scrollAmount = clientWidth * 0.8; // Scroll by 80% of container width

                carouselRef.current.scrollTo({
                    left: direction === "left" ? scrollLeft - scrollAmount : scrollLeft + scrollAmount,
                    behavior: "smooth"
                });
            }
        };

        return (
            <div className="mb-12 relative group/carousel">
                <h3 className="text-2xl font-bold mb-6 font-display flex items-center gap-2">
                    <div className="w-1 h-6 bg-mint rounded-full"></div>
                    {title}
                </h3>
                <div
                    ref={carouselRef}
                    className="flex overflow-x-auto gap-6 pb-6 px-2 -mx-2 hide-scrollbar snap-x relative z-10"
                >
                    {movies.map((movie) => (
                        <div key={movie.id} className="min-w-[160px] md:min-w-[200px] snap-start">
                            <motion.div
                                initial={{ opacity: 0, y: 10 }}
                                whileInView={{ opacity: 1, y: 0 }}
                                viewport={{ once: true, margin: "-50px" }}
                                className="h-full"
                            >
                                <MovieCard movie={movie} />
                            </motion.div>
                        </div>
                    ))}
                </div>

                {/* Left Arrow */}
                <button
                    onClick={() => scroll("left")}
                    className="absolute left-0 top-[60%] -translate-y-1/2 -ml-4 md:-ml-8 bg-black/50 hover:bg-black/80 text-white p-2 rounded-full backdrop-blur-md opacity-0 group-hover/carousel:opacity-100 transition-opacity z-20 shadow-lg border border-white/10 hidden md:flex"
                >
                    <ChevronLeft size={24} />
                </button>

                {/* Right Arrow */}
                <button
                    onClick={() => scroll("right")}
                    className="absolute right-0 top-[60%] -translate-y-1/2 -mr-4 md:-mr-8 bg-black/50 hover:bg-black/80 text-white p-2 rounded-full backdrop-blur-md opacity-0 group-hover/carousel:opacity-100 transition-opacity z-20 shadow-lg border border-white/10 hidden md:flex"
                >
                    <ChevronRight size={24} />
                </button>
            </div>
        );
    };

    return (
        <div className="min-h-screen bg-[#12201B] text-white p-6 md:p-8 lg:px-12 pt-24">
            <header className="mb-12 max-w-4xl">
                <motion.h1
                    initial={{ opacity: 0, y: -20 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="text-4xl md:text-5xl font-black font-display mb-4"
                >
                    Explore <span className="text-mint">Films</span>
                </motion.h1>
                <motion.p
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ delay: 0.2 }}
                    className="text-gray-400 text-lg"
                >
                    Discover the most popular releases, highest rated classics, and upcoming theatrical hits.
                </motion.p>
            </header>

            {loading ? (
                <div className="flex justify-center py-32"><div className="loader"></div></div>
            ) : (
                <div className="animate-fade-in relative">
                    <MovieCarousel title="Now Popular" movies={categories.popular} />
                    <MovieCarousel title="Highest Rated" movies={categories.topRated} />
                    <MovieCarousel title="Upcoming Releases" movies={categories.upcoming} />
                </div>
            )}

            {/* Minimalist styling for hiding horizontal scrollbar but keeping functionality */}
            <style dangerouslySetInnerHTML={{
                __html: `
                .hide-scrollbar::-webkit-scrollbar {
                    display: none;
                }
                .hide-scrollbar {
                    -ms-overflow-style: none;
                    scrollbar-width: none;
                }
            `}} />
        </div>
    );
};

export default MoviesPage;
