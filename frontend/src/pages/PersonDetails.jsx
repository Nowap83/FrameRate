import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import apiClient from "../api/apiClient";
import { User, Calendar, MapPin } from "lucide-react";

const PersonDetails = () => {
    const { id } = useParams();
    const [person, setPerson] = useState(null);
    const [credits, setCredits] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchPersonData = async () => {
            setLoading(true);
            try {
                const [detailsRes, creditsRes] = await Promise.all([
                    apiClient.get(`/tmdb/person/${id}`),
                    apiClient.get(`/tmdb/person/${id}/credits`)
                ]);

                if (detailsRes.data.success) {
                    setPerson(detailsRes.data.data);
                }
                if (creditsRes.data.success) {
                    setCredits(creditsRes.data.data);
                }
            } catch (error) {
                console.error("Failed to fetch person data:", error);
            } finally {
                setLoading(false);
            }
        };

        if (id) {
            fetchPersonData();
        }
    }, [id]);

    if (loading) {
        return (
            <div className="min-h-screen bg-[var(--color-body-bg)] flex items-center justify-center text-white">
                <div className="loader"></div>
            </div>
        );
    }

    if (!person) {
        return (
            <div className="min-h-screen bg-[var(--color-body-bg)] flex items-center justify-center text-white">
                <div className="text-center">
                    <h2 className="text-2xl font-bold mb-2">Person Not Found</h2>
                    <p className="text-gray-400">The person you are looking for does not exist or we couldn't load their data.</p>
                </div>
            </div>
        );
    }

    const isActor = person.known_for_department === 'Acting' || (credits?.cast?.length > credits?.crew?.length);
    const filmography = isActor ? credits?.cast || [] : credits?.crew || [];

    const uniqueFilmography = Array.from(new Map(filmography.map(movie => [movie.id, movie])).values())
        .sort((a, b) => {
            const dateA = a.release_date ? new Date(a.release_date) : new Date('1900-01-01');
            const dateB = b.release_date ? new Date(b.release_date) : new Date('1900-01-01');
            return dateB - dateA;
        });


    return (
        <div className="min-h-screen bg-[var(--color-body-bg)]">
            <div className="container mx-auto px-4 py-8 md:py-12 mt-16 md:mt-24">
                <div className="flex flex-col md:flex-row gap-8 lg:gap-12">

                    {/* left column: photo and personal info */}
                    <div className="md:w-1/3 lg:w-1/4 flex-shrink-0">
                        <div className="bg-[#1a1a1a] rounded-2xl p-2 w-full max-w-sm mx-auto shadow-2xl overflow-hidden border border-white/5 relative group">
                            <div className="aspect-[2/3] rounded-xl overflow-hidden relative">
                                {person.profile_path ? (
                                    <img
                                        src={`https://image.tmdb.org/t/p/w500${person.profile_path}`}
                                        alt={person.name}
                                        className="w-full h-full object-cover"
                                    />
                                ) : (
                                    <div className="w-full h-full bg-white/5 flex items-center justify-center text-gray-500">
                                        <User size={64} className="opacity-20" />
                                    </div>
                                )}
                            </div>
                        </div>

                        {/* personal info box */}
                        <div className="mt-6 bg-[#12201B] rounded-2xl p-6 border border-white/5 space-y-4">
                            <h3 className="text-lg font-bold text-white border-b border-white/10 pb-2 mb-4">Personal Info</h3>

                            {person.known_for_department && (
                                <div>
                                    <p className="text-xs text-gray-400 uppercase tracking-wider mb-1">Known For</p>
                                    <p className="text-sm text-gray-200">{person.known_for_department}</p>
                                </div>
                            )}

                            {person.gender > 0 && (
                                <div>
                                    <p className="text-xs text-gray-400 uppercase tracking-wider mb-1">Gender</p>
                                    <p className="text-sm text-gray-200">{person.gender === 1 ? 'Female' : person.gender === 2 ? 'Male' : 'Other'}</p>
                                </div>
                            )}

                            {person.birthday && (
                                <div>
                                    <p className="text-xs text-gray-400 uppercase tracking-wider mb-1 flex items-center gap-1">
                                        <Calendar size={12} /> Born
                                    </p>
                                    <p className="text-sm text-gray-200">
                                        {new Date(person.birthday).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}
                                        {!person.deathday && (
                                            <span className="text-gray-500 ml-1">
                                                ({new Date().getFullYear() - new Date(person.birthday).getFullYear()} years old)
                                            </span>
                                        )}
                                    </p>
                                </div>
                            )}

                            {person.deathday && (
                                <div>
                                    <p className="text-xs text-gray-400 uppercase tracking-wider mb-1">Died</p>
                                    <p className="text-sm text-gray-200">
                                        {new Date(person.deathday).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}
                                    </p>
                                </div>
                            )}

                            {person.place_of_birth && (
                                <div>
                                    <p className="text-xs text-gray-400 uppercase tracking-wider mb-1 flex items-center gap-1">
                                        <MapPin size={12} /> Place of Birth
                                    </p>
                                    <p className="text-sm text-gray-200">{person.place_of_birth}</p>
                                </div>
                            )}
                        </div>
                    </div>

                    {/* right column: bio and filmography */}
                    <div className="md:w-2/3 lg:w-3/4">
                        <div className="mb-10">
                            <h1 className="text-4xl md:text-5xl font-black text-white mb-6 tracking-tight">{person.name}</h1>

                            {person.biography && (
                                <div>
                                    <h3 className="text-xl font-bold text-white mb-3 flex items-center gap-2">
                                        Biography
                                    </h3>
                                    <div className="text-gray-300 text-sm md:text-base leading-relaxed space-y-4 opacity-90">
                                        {person.biography.split('\n\n').map((paragraph, index) => (
                                            <p key={index}>{paragraph}</p>
                                        ))}
                                    </div>
                                </div>
                            )}
                        </div>

                        {/* filmography grid */}
                        <div className="mt-12">
                            <div className="flex items-center justify-between border-b border-white/10 pb-4 mb-6">
                                <h2 className="text-2xl font-bold text-white">Filmography</h2>
                                <span className="text-sm text-gray-500 bg-white/5 px-3 py-1 rounded-full">{uniqueFilmography.length} movies</span>
                            </div>

                            {uniqueFilmography.length > 0 ? (
                                <div className="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 lg:grid-cols-6 gap-4">
                                    {uniqueFilmography.map((movie) => (
                                        <Link
                                            key={movie.id}
                                            to={`/movie/${movie.id}`}
                                            className="group relative block rounded-xl overflow-hidden shadow-md hover:shadow-[0_0_15px_rgba(33,184,136,0.3)] transition-all hover:-translate-y-1"
                                            title={`${movie.title} (${movie.release_date?.substring(0, 4) || 'Unknown'})`}
                                        >
                                            <div className="aspect-[2/3] bg-[#1a1a1a]">
                                                {movie.poster_path ? (
                                                    <img
                                                        src={`https://image.tmdb.org/t/p/w342${movie.poster_path}`}
                                                        alt={movie.title}
                                                        className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
                                                        loading="lazy"
                                                    />
                                                ) : (
                                                    <div className="w-full h-full flex flex-col items-center justify-center p-2 text-center text-gray-600 border border-white/5">
                                                        <span className="text-xs font-medium line-clamp-3">{movie.title}</span>
                                                    </div>
                                                )}

                                                {/* hover overlay with year */}
                                                <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex flex-col justify-end p-2 md:p-3">
                                                    <span className="text-white text-xs font-bold truncate">{movie.title}</span>
                                                    {movie.release_date && (
                                                        <span className="text-[var(--color-primary)] text-[10px] font-medium">{movie.release_date.substring(0, 4)}</span>
                                                    )}
                                                </div>
                                            </div>
                                        </Link>
                                    ))}
                                </div>
                            ) : (
                                <div className="text-center py-12 bg-white/5 rounded-2xl border border-white/5">
                                    <p className="text-gray-400">No filmography data available.</p>
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default PersonDetails;
