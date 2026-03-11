import { Link } from "react-router-dom";
import RatingStars from "./RatingStars";

const POSTER_BASE_URL = "https://image.tmdb.org/t/p/w500";

const ReviewCard = ({ review }) => {
    const {
        tmdb_id,
        title,
        release_year,
        poster_url,
        rating,
        content,
        watched_date,
    } = review;

    // Use watched_date if available, otherwise just show year
    const watchedDateFormatted = watched_date 
        ? new Date(watched_date).toLocaleDateString("en-GB", { day: '2-digit', month: 'short', year: 'numeric' })
        : null;

    return (
        <div className="flex gap-4 md:gap-6 py-6 border-b border-white/5 last:border-0 px-4 rounded-xl">
            {/* Poster */}
            <Link to={`/movie/${tmdb_id}`} className="shrink-0 w-24 md:w-32 aspect-[2/3] block rounded-md overflow-hidden bg-gray-800 shadow-md">
                {poster_url ? (
                    <img
                        src={`${POSTER_BASE_URL}${poster_url}`}
                        alt={title}
                        className="w-full h-full object-cover transition-transform hover:scale-105"
                        loading="lazy"
                    />
                ) : (
                    <div className="w-full h-full flex items-center justify-center text-xs text-gray-500">No Img</div>
                )}
            </Link>

            {/* Content */}
            <div className="flex-1 flex flex-col">
                <div className="flex flex-wrap items-baseline gap-2 mb-1">
                    <Link to={`/movie/${tmdb_id}`} className="text-xl md:text-2xl font-bold font-display hover:text-[var(--color-primary)] transition-colors">
                        {title}
                    </Link>
                    {release_year > 0 && <span className="text-gray-400 font-normal text-sm md:text-base">{release_year}</span>}
                </div>

                <div className="flex flex-wrap items-center gap-2 md:gap-3 mb-4">
                    {rating > 0 && (
                        <div className="flex">
                            <RatingStars rating={rating} size={14} readOnly />
                        </div>
                    )}
                    {watchedDateFormatted && (
                       <span className="text-xs md:text-sm text-gray-400 font-medium tracking-wide">Watched {watchedDateFormatted}</span> 
                    )}
                </div>

                <div className="text-gray-300 text-sm md:text-base leading-relaxed whitespace-pre-wrap mb-4">
                    {content}
                </div>
            </div>
        </div>
    );
};

export default ReviewCard;
