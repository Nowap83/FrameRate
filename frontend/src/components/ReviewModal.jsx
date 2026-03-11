import React, { useState, useEffect } from "react";
import { X } from "lucide-react";
import RatingStars from "./RatingStars";

const POSTER_BASE_URL = "https://image.tmdb.org/t/p/w500";

const ReviewModal = ({ isOpen, onClose, movie, onSave, initialRating = 0, initialWatched = false }) => {
    const [rating, setRating] = useState(initialRating);
    const [reviewText, setReviewText] = useState("");
    const [isWatched, setIsWatched] = useState(initialWatched);
    const [watchedDate, setWatchedDate] = useState(() => {
        const today = new Date();
        return today.toISOString().split("T")[0]; // YYYY-MM-DD
    });
    const [isSubmitting, setIsSubmitting] = useState(false);

    useEffect(() => {
        if (isOpen) {
            setRating(initialRating);
            setIsWatched(initialWatched);
            setReviewText("");
        }
    }, [isOpen, initialRating, initialWatched]);

    if (!isOpen || !movie) return null;

    // Auto-check "Watched" if typing a review or rating > 0
    const handleRatingChange = (newRating) => {
        setRating(newRating);
        if (!isWatched && newRating > 0) {
            setIsWatched(true);
        }
    };

    const handleReviewChange = (e) => {
        setReviewText(e.target.value);
        if (!isWatched && e.target.value.length > 0) {
            setIsWatched(true);
        }
    };

    const handleSubmit = async () => {
        setIsSubmitting(true);
        try {
            await onSave({
                rating: rating > 0 ? rating : null,
                review_text: reviewText.trim() || null,
                is_spoiler: false,
                watched_date: isWatched ? new Date(watchedDate).toISOString() : null,
            });
            onClose();
        } catch (error) {
            console.error("Failed to save log:", error);
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            {/* Backdrop */}
            <div 
                className="absolute inset-0 bg-black/80 backdrop-blur-sm"
                onClick={onClose}
            ></div>

            {/* Modal Content */}
            <div className="relative bg-[#445566] text-white w-full max-w-2xl rounded-lg shadow-2xl overflow-hidden animate-in fade-in zoom-in-95 duration-200">
                {/* Header */}
                <div className="flex items-center justify-between px-6 py-4 border-b border-white/10 bg-[#3B4A59]">
                    <h2 className="text-xl font-bold font-display">I watched...</h2>
                    <button 
                        onClick={onClose}
                        className="p-1 hover:bg-white/10 rounded-full transition-colors"
                    >
                        <X size={24} className="text-gray-300" />
                    </button>
                </div>

                <div className="p-6">
                    <div className="flex gap-6">
                        {/* Poster */}
                        <div className="shrink-0 w-32 md:w-40">
                            {movie.poster_path ? (
                                <img
                                    src={`${POSTER_BASE_URL}${movie.poster_path}`}
                                    alt={movie.title}
                                    className="w-full rounded-md shadow-lg border border-white/5"
                                />
                            ) : (
                                <div className="w-full aspect-[2/3] bg-gray-800 rounded-md flex items-center justify-center text-sm">No Image</div>
                            )}
                        </div>

                        {/* Form Area */}
                        <div className="flex-1 flex flex-col">
                            <h3 className="text-2xl font-bold font-display mb-1">
                                {movie.title} <span className="text-gray-400 font-normal text-lg">{movie.release_date?.split("-")[0]}</span>
                            </h3>

                            <div className="flex flex-wrap gap-6 mt-4 mb-4">
                                <label className="flex items-center gap-2 cursor-pointer group">
                                    <input 
                                        type="checkbox" 
                                        checked={isWatched}
                                        onChange={(e) => setIsWatched(e.target.checked)}
                                        className="w-5 h-5 rounded border-gray-500 bg-black/20 text-mint focus:ring-mint focus:ring-offset-0"
                                    />
                                    <span className="text-sm font-medium text-gray-200 group-hover:text-white transition-colors">Watched on</span>
                                    {isWatched && (
                                        <input 
                                            type="date"
                                            value={watchedDate}
                                            onChange={(e) => setWatchedDate(e.target.value)}
                                            className="bg-black/20 border border-white/10 rounded px-2 py-1 text-sm focus:outline-none focus:border-mint"
                                        />
                                    )}
                                </label>
                            </div>

                            <textarea
                                placeholder="Add a review..."
                                value={reviewText}
                                onChange={handleReviewChange}
                                className="w-full flex-1 min-h-[120px] bg-[#EEF2F5] text-gray-900 placeholder:text-gray-500 p-4 rounded-md resize-none focus:outline-none focus:ring-2 focus:ring-mint"
                            ></textarea>

                            {/* Tags and Rating */}
                            <div className="flex items-center justify-between mt-4">
                                <div className="flex gap-4 items-center">
                                    <div className="flex flex-col">
                                        <span className="text-xs font-bold text-gray-300 mb-1">Rating</span>
                                        <RatingStars 
                                            rating={rating} 
                                            onChange={handleRatingChange} 
                                            size={24} 
                                        />
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Footer Actions */}
                <div className="flex justify-end gap-3 px-6 py-4 border-t border-white/10 bg-[#3B4A59]">
                    <button 
                        className="px-6 py-2 bg-[#00E054] hover:bg-[#00c248] text-white font-bold rounded transition-colors disabled:opacity-50"
                        onClick={handleSubmit}
                        disabled={isSubmitting}
                    >
                        {isSubmitting ? "SAVING..." : "SAVE"}
                    </button>
                </div>
            </div>
        </div>
    );
};

export default ReviewModal;
