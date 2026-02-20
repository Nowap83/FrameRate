import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { getMovieDetails, getMovieVideos } from "../api/tmdb";
import Button from "../components/Button";
import { Star, Eye, Plus, List, Play, Heart } from "lucide-react";

const IMAGE_BASE_URL = "https://image.tmdb.org/t/p/original";
const POSTER_BASE_URL = "https://image.tmdb.org/t/p/w500";

const MovieDetails = () => {
    const { id } = useParams();
    const [movie, setMovie] = useState(null);
    const [loading, setLoading] = useState(true);
    const [activeTab, setActiveTab] = useState("CAST");
    const [trailerKey, setTrailerKey] = useState(null);

    useEffect(() => {
        const fetchDetails = async () => {
            try {
                const data = await getMovieDetails(id);
                // console.log("Movie Data Received:", data);
                if (data.data) {
                    setMovie(data.data);
                } else {
                    console.error("No data property in response", data);
                }
            } catch (error) {
                console.error("Failed to fetch movie details", error);
            } finally {
                setLoading(false);
            }
        };

        const fetchTrailer = async () => {
            try {
                const data = await getMovieVideos(id);
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

        if (id) {
            fetchDetails();
            fetchTrailer();
        }
    }, [id]);

    if (loading) {
        return <div className="min-h-screen flex items-center justify-center text-white"><div className="loader"></div></div>;
    }

    if (!movie) {
        return <div className="min-h-screen flex items-center justify-center text-white">Movie not found</div>;
    }

    const director = movie.credits?.crew?.find(c => c.job === "Director")?.name;
    const writers = movie.credits?.crew?.filter(c => c.department === "Writing").slice(0, 2).map(c => c.name).join(", ");

    // format durée
    const hours = Math.floor(movie.runtime / 60);
    const minutes = movie.runtime % 60;
    const formattedRuntime = `${hours}h ${minutes}m`;

    const formatCurrency = (value) => {
        if (!value) return "Unknown";
        return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', maximumFractionDigits: 0 }).format(value);
    };

    const handleWatchTrailer = () => {
        if (trailerKey) {
            window.open(`https://www.youtube.com/watch?v=${trailerKey}`, "_blank");
        }
    };

    const renderTabContent = () => {
        switch (activeTab) {
            case "CAST":
                return (
                    <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
                        {movie.credits?.cast?.slice(0, 10).map(actor => (
                            <Link key={actor.id} to={`/person/${actor.id}`} className="group block">
                                <div className="aspect-[2/3] overflow-hidden rounded-lg bg-gray-800 mb-2">
                                    {actor.profile_path ? (
                                        <img
                                            src={`${POSTER_BASE_URL}${actor.profile_path}`}
                                            alt={actor.name}
                                            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                                        />
                                    ) : (
                                        <div className="w-full h-full flex items-center justify-center text-gray-600 text-xs text-center border border-white/5 p-2">No Image</div>
                                    )}
                                </div>
                                <h4 className="font-bold text-sm truncate group-hover:text-mint transition-colors">{actor.name}</h4>
                                <p className="text-xs text-gray-500 truncate">{actor.character}</p>
                            </Link>
                        ))}
                    </div>
                );
            case "CREW": {
                const importantCrew = movie.credits?.crew?.filter(c => ["Directing", "Writing", "Production"].includes(c.department)) || [];
                return (
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
                        {importantCrew.slice(0, 12).map(member => (
                            <Link key={`${member.id}-${member.job}`} to={`/person/${member.id}`} className="bg-white/5 p-4 rounded-lg block hover:bg-white/10 transition-colors">
                                <h4 className="font-bold text-sm">{member.name}</h4>
                                <p className="text-xs text-mint">{member.job}</p>
                            </Link>
                        ))}
                    </div>
                );
            }
            case "DETAILS":
                return (
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-8 text-sm">
                        <div>
                            <h4 className="font-bold text-gray-400 uppercase tracking-widest mb-2">Original Title</h4>
                            <p className="text-xl mb-6">{movie.original_title}</p>

                            <h4 className="font-bold text-gray-400 uppercase tracking-widest mb-2">Original Language</h4>
                            <p className="text-xl mb-6 uppercase">{movie.original_language}</p>
                        </div>
                        <div>
                            <h4 className="font-bold text-gray-400 uppercase tracking-widest mb-2">Budget</h4>
                            <p className="text-xl mb-6">{formatCurrency(movie.budget)}</p>

                            <h4 className="font-bold text-gray-400 uppercase tracking-widest mb-2">Revenue</h4>
                            <p className="text-xl mb-6">{formatCurrency(movie.revenue)}</p>
                        </div>
                    </div>
                );
            case "GENRE":
                return (
                    <div className="flex flex-wrap gap-3">
                        {movie.genres?.map(genre => (
                            <span key={genre.id} className="px-4 py-2 bg-white/10 rounded-lg text-sm font-medium hover:bg-white/20 transition-colors cursor-default">
                                {genre.name}
                            </span>
                        ))}
                    </div>
                );
            case "RELEASE": {
                const releaseDate = movie.release_date && !isNaN(new Date(movie.release_date).getTime())
                    ? new Date(movie.release_date).toLocaleDateString(undefined, { dateStyle: "long" })
                    : "Unknown Release Date";

                return (
                    <div className="bg-white/5 p-6 rounded-xl inline-block">
                        <h4 className="font-bold text-gray-400 uppercase tracking-widest mb-2">Release Date</h4>
                        <p className="text-3xl font-display">{releaseDate}</p>
                    </div>
                );
            }
            default:
                return null;
        }
    };

    return (
        <div className="min-h-screen bg-[#12201B] text-white font-sans pb-20">
            {/* hero section */}
            <div className="relative w-full h-[500px] overflow-hidden">
                <div className="absolute inset-0 bg-gradient-to-b from-transparent to-[#12201B] z-10" />
                {movie.backdrop_path && (
                    <img
                        src={`${IMAGE_BASE_URL}${movie.backdrop_path}`}
                        alt={movie.title}
                        className="w-full h-full object-cover opacity-60"
                    />
                )}
            </div>

            <div className="max-w-7xl mx-auto px-6 -mt-80 relative z-20">
                <div className="flex flex-col md:flex-row gap-10">
                    {/* poster */}
                    <div className="shrink-0 w-72">
                        {movie.poster_path ? (
                            <img
                                src={`${POSTER_BASE_URL}${movie.poster_path}`}
                                alt={movie.title}
                                className="w-full rounded-xl shadow-2xl border border-gray-700/50"
                            />
                        ) : (
                            <div className="w-full h-[430px] bg-gray-800 rounded-xl flex items-center justify-center">No Image</div>
                        )}

                        <button
                            className={`w-full mt-4 flex items-center justify-center gap-2 bg-white/10 hover:bg-white/20 backdrop-blur-md py-3 rounded-xl font-semibold transition-all ${!trailerKey ? 'opacity-50 cursor-not-allowed' : ''}`}
                            onClick={handleWatchTrailer}
                            disabled={!trailerKey}
                        >
                            <Play size={20} fill="currentColor" />
                            WATCH TRAILER
                        </button>
                    </div>

                    {/* content */}
                    <div type="info" className="flex-1 pt-10 md:pt-32">
                        <h1 className="text-6xl font-display font-bold mb-2">{movie.title}</h1>

                        <div className="flex items-center gap-4 text-gray-300 mb-6 text-sm">
                            <span>{movie.release_date?.split("-")[0]}</span>
                            <span>•</span>
                            <span>{formattedRuntime}</span>
                            <span>•</span>
                            <span>{director}</span>
                        </div>

                        {/* genres */}
                        <div className="flex gap-2 mb-8">
                            {movie.genres?.map(genre => (
                                <span key={genre.id} className="px-3 py-1 bg-white/10 rounded-full text-xs font-medium uppercase tracking-wider backdrop-blur-sm">
                                    {genre.name}
                                </span>
                            ))}
                        </div>

                        {/* Synopsis */}
                        <div className="mb-8 max-w-2xl">
                            <h3 className="text-xs font-bold text-gray-400 uppercase tracking-widest mb-2">Synopsis</h3>
                            <p className="text-gray-300 leading-relaxed text-lg">
                                {movie.overview}
                            </p>
                        </div>

                        {/* actions card */}
                        <div className="flex flex-wrap items-center gap-6 mb-12 bg-[#1A2C24] p-6 rounded-2xl border border-white/5 max-w-3xl">
                            {/* rating */}
                            <div className="flex flex-col pr-8 border-r border-white/10">
                                <span className="text-xs font-bold text-gray-400 uppercase tracking-widest mb-1">Avg Rating</span>
                                <div className="flex items-end gap-2">
                                    <span className="text-4xl font-bold">{movie.vote_average ? movie.vote_average.toFixed(1) : "N/A"}</span>
                                    <div className="flex pb-1">
                                        {[...Array(5)].map((_, i) => (
                                            <Star
                                                key={i}
                                                size={16}
                                                className={i < Math.round((movie.vote_average || 0) / 2) ? "fill-white text-white" : "text-gray-600"}
                                            />
                                        ))}
                                    </div>
                                </div>
                                <span className="text-xs text-gray-500 mt-1">{movie.vote_count} Reviews</span>
                            </div>

                            {/* user actions */}
                            <div className="flex items-center gap-4 flex-1 justify-end">
                                <div className="flex flex-col items-center gap-2">
                                    <button className="p-3 bg-white/5 hover:bg-white/10 rounded-full transition-colors">
                                        <Eye size={24} />
                                    </button>
                                </div>

                                <button className="flex items-center gap-2 px-6 py-3 bg-mint text-emerald-950 font-bold rounded-xl hover:brightness-110 transition-all">
                                    <Plus size={20} className="stroke-[3]" />
                                    ADD TO WATCHLIST
                                </button>

                                <button className="flex items-center gap-2 px-6 py-3 bg-white/5 text-white font-bold rounded-xl hover:bg-white/10 transition-all">
                                    <List size={20} />
                                    ADD TO LIST
                                </button>

                                <button className="p-3 bg-white/5 hover:bg-white/10 rounded-full transition-colors">
                                    <Heart size={24} />
                                </button>
                            </div>
                        </div>

                        {/* tabs */}
                        <div className="flex gap-8 border-b border-white/10 mb-8">
                            {["CAST", "CREW", "DETAILS", "GENRE", "RELEASE"].map((tab) => (
                                <button
                                    key={tab}
                                    onClick={() => setActiveTab(tab)}
                                    className={`pb-4 text-sm font-bold tracking-widest transition-colors ${activeTab === tab ? "text-mint border-b-2 border-mint" : "text-gray-500 hover:text-white"}`}
                                >
                                    {tab}
                                </button>
                            ))}
                        </div>

                        {/* tab content */}
                        <div className="min-h-[300px] animate-in fade-in duration-500">
                            {renderTabContent()}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default MovieDetails;
