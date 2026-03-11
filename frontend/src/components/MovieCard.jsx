import { useState } from "react";
import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { Star, Plus, AlignLeft } from "lucide-react";
import ReviewPreviewModal from "./ReviewPreviewModal";
import { getMovieInteraction } from "../api/tmdb";


const MovieCard = ({ movie }) => {
    // TMDB API format et Backend DB format
    const id = movie.tmdb_id || movie.id;
    const posterPath = movie.poster_url || movie.poster_path;
    const title = movie.title;
    
    const [isReviewModalOpen, setIsReviewModalOpen] = useState(false);
    const [reviewData, setReviewData] = useState(null);

    // rating logic: user_rating si available (Profile), sinon vote_average (TMDB)
    let rating = null;
    if (movie.user_rating != null && movie.user_rating > 0) {
        rating = movie.user_rating.toFixed(1);
    } else if (movie.vote_average != null && movie.vote_average > 0) {
        rating = movie.vote_average.toFixed(1);
    }

    // year logic
    const releaseDate = movie.release_date || (movie.release_year ? `${movie.release_year}-01-01` : null);
    const year = releaseDate ? new Date(releaseDate).getFullYear() : 'Unknown';

    if (!posterPath) return null; // pas de render si pas d'image

    const handleReviewClick = async (e) => {
        e.preventDefault();
        e.stopPropagation();
        
        try {
            const data = await getMovieInteraction(id);
            if (data && data.user_review) {
                setReviewData({
                    title: movie.title,
                    release_year: year !== 'Unknown' ? year : 0,
                    rating: data.user_rating || 0,
                    content: data.user_review.content,
                    watched_date: data.watched_date
                });
                setIsReviewModalOpen(true);
            }
        } catch (error) {
            console.error("Failed to fetch review", error);
        }
    };

    return (
        <div className="relative aspect-[2/3] rounded-xl overflow-hidden group shadow-lg">
            <Link to={`/movie/${id}`} className="block w-full h-full">
                <img
                    src={`https://image.tmdb.org/t/p/w500${posterPath}`}
                    alt={title}
                    className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
                    loading="lazy"
                />

                {/* 
                  overlay + quick actions
                */}
                <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col items-center justify-center gap-4">
                    <div className="flex items-center gap-4 transform translate-y-4 group-hover:translate-y-0 transition-transform duration-300">
                        <button
                            className="p-3 bg-mint rounded-full text-black hover:scale-110 transition-transform shadow-lg"
                            onClick={(e) => e.preventDefault()}
                            title="Rate Movie"
                        >
                            <Star size={20} />
                        </button>
                        <button
                            className="p-3 bg-white/20 rounded-full text-white hover:scale-110 transition-transform shadow-lg"
                            onClick={(e) => e.preventDefault()}
                            title="Add to Watchlist"
                        >
                            <Plus size={20} />
                        </button>
                    </div>

                    <div className="absolute bottom-0 left-0 w-full p-4 bg-gradient-to-t from-black/90 to-transparent flex flex-col items-center">
                        <span className="text-white font-bold text-sm text-center line-clamp-2 leading-tight">{title}</span>
                        <span className="text-mint text-xs mt-1 font-medium">{year}</span>
                    </div>
                </div>

                {/* review badge */}
                {movie.has_review && (
                    <div 
                        className="absolute top-2 left-2 bg-black/70 backdrop-blur-md p-1.5 rounded-md flex items-center ring-1 ring-white/10 z-20 transition-transform group-hover:scale-110 cursor-pointer"
                        onClick={handleReviewClick}
                    >
                        <AlignLeft size={14} className="text-[var(--color-primary)]" />
                    </div>
                )}

                {/* rating badge */}
                {rating && (
                    <div className="absolute top-2 right-2 bg-black/70 backdrop-blur-md px-2 py-1 rounded-md flex items-center gap-1 ring-1 ring-white/10 z-10">
                        <Star size={12} className="text-yellow-400 fill-yellow-400" />
                        <span className="text-white text-xs font-bold">{rating}</span>
                    </div>
                )}
            </Link>
            
            <ReviewPreviewModal 
                isOpen={isReviewModalOpen} 
                onClose={() => setIsReviewModalOpen(false)} 
                review={reviewData} 
            />
        </div>
    );
};

export default MovieCard;
