
import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { Link } from 'react-router-dom';
import { Settings, Share2, Heart, Activity, Calendar, MapPin, Film, User, Star, Globe } from 'lucide-react';
import { getAvatarUrl } from '../utils/image';
import useDocumentTitle from '../hooks/useDocumentTitle';

const ProfileHeader = ({ user }) => (
    <div className="relative mb-12 mt-24 px-4 md:px-8 max-w-7xl mx-auto flex flex-col md:flex-row items-end gap-6">
        <div className="relative group">
            <div className="w-32 h-32 md:w-40 md:h-40 rounded-full p-1 bg-[var(--color-body-bg)] shadow-2xl">
                <img
                    src={getAvatarUrl(user?.profile_picture_url) || `https://ui-avatars.com/api/?name=${user?.username}&background=random`}
                    alt={user?.username}
                    className="w-full h-full object-cover rounded-full border-4 border-[var(--color-body-bg)]"
                />
            </div>
            <div className="absolute bottom-2 right-2 bg-[var(--color-primary)] text-black p-1.5 rounded-full border-2 border-[var(--color-body-bg)] shadow-lg">
                <Star size={14} fill="currentColor" />
            </div>
        </div>

        <div className="flex-1 pb-2 text-center md:text-left">
            <h1 className="text-4xl md:text-5xl font-bold text-white mb-2 tracking-tight">
                {user?.given_name ? `${user.given_name} ${user.family_name || ''}` : user?.username}
            </h1>
            {user?.given_name && <p className="text-gray-500 font-mono text-sm mb-2">@{user?.username}</p>}

            <p className="text-gray-400 mb-4 max-w-2xl text-lg leading-relaxed line-clamp-2">
                {user?.bio}
            </p>

            <div className="flex items-center justify-center md:justify-start gap-6 text-sm text-gray-400 font-medium flex-wrap">
                {user?.location && (
                    <span className="flex items-center gap-1.5">
                        <MapPin size={16} className="text-[var(--color-primary)]" />
                        {user.location}
                    </span>
                )}
                {user?.website && (
                    <a href={user.website} target="_blank" rel="noopener noreferrer" className="flex items-center gap-1.5 hover:text-white transition-colors">
                        <Globe size={16} className="text-[var(--color-primary)]" />
                        Website
                    </a>
                )}
                <span className="flex items-center gap-1.5">
                    <Calendar size={16} className="text-[var(--color-primary)]" />
                    Joined {new Date(user?.created_at).toLocaleDateString(undefined, { month: 'long', year: 'numeric' })}
                </span>
            </div>
        </div>

        <div className="flex gap-3 pb-4 w-full md:w-auto justify-center md:justify-end">
            <Link
                to="/settings"
                className="flex items-center gap-2 px-6 py-2.5 bg-[#BDF4DF] text-black font-bold rounded-full hover:bg-[#a6e6ce] transition-all hover:scale-105 shadow-[0_0_20px_rgba(189,244,223,0.3)]"
            >
                <Settings size={18} />
                Edit Profile
            </Link>
            <button className="p-2.5 bg-white/5 backdrop-blur-md text-white rounded-full hover:bg-white/10 transition-all border border-white/10 hover:border-white/20">
                <Share2 size={20} />
            </button>
        </div>
    </div>
);

const StatsCard = ({ label, value, subtext, trend }) => (
    <div className="glass-panel p-5 rounded-2xl flex flex-col justify-between h-28 hover:bg-white/5 transition-colors group relative overflow-hidden">
        <div className="flex justify-between items-start">
            <span className="text-xs font-bold text-gray-500 uppercase tracking-widest leading-none">{label}</span>
        </div>
        <div>
            <div className="flex items-baseline gap-2">
                <span className="text-3xl font-bold text-white tracking-tight group-hover:text-[var(--color-primary)] transition-colors">{value}</span>
                {trend && <span className="text-xs text-green-400 font-medium">{trend}</span>}
            </div>
            {subtext && <span className="text-xs text-[var(--color-primary)] opacity-80 font-medium">{subtext}</span>}
        </div>
    </div>
);

const SectionHeader = ({ title, action }) => (
    <div className="flex items-center justify-between mb-6 border-b border-white/5 pb-2">
        <h3 className="text-lg font-bold text-white uppercase tracking-wider">{title}</h3>
        {action && <button className="text-xs font-bold text-gray-500 hover:text-white transition-colors">{action}</button>}
    </div>
);

const MovieGrid = ({ movies }) => {
    if (!movies || movies.length === 0) return (
        <div className="h-40 flex items-center justify-center border border-dashed border-white/10 rounded-xl">
            <p className="text-gray-500 text-sm">No movies found.</p>
        </div>
    );

    return (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {movies.map(movie => (
                <Link to={`/movie/${movie.tmdb_id}`} key={movie.id} className="relative aspect-[2/3] rounded-xl overflow-hidden group shadow-lg">
                    <img
                        src={`https://image.tmdb.org/t/p/w500${movie.poster_url || movie.poster_path}`}
                        alt={movie.title}
                        className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
                    />
                    <div className="absolute inset-0 bg-gradient-to-t from-black/90 via-black/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex flex-col justify-end p-4">
                        <span className="text-white font-bold text-sm line-clamp-2 leading-tight transform translate-y-4 group-hover:translate-y-0 transition-transform duration-300">{movie.title}</span>
                        <span className="text-[var(--color-primary)] text-xs mt-1 transform translate-y-4 group-hover:translate-y-0 transition-transform duration-300 delay-75">{new Date(movie.release_date).getFullYear()}</span>
                    </div>
                </Link>
            ))}
        </div>
    );
};

const Profile = () => {
    const { user: authUser, loading: authLoading } = useAuth();
    const [profileData, setProfileData] = useState(null);
    const [loading, setLoading] = useState(true);

    const fetchProfile = async () => {
        try {
            const response = await apiClient.get('/users/me');
            setProfileData(response.data);
        } catch (error) {
            console.error("Failed to fetch profile:", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (!authLoading) {
            if (authUser) {
                fetchProfile();
            } else {
                setLoading(false);
            }
        }
    }, [authLoading, authUser]);

    useDocumentTitle(authUser?.username || "Profile");

    if (authLoading || loading) return <div className="min-h-screen bg-[var(--color-body-bg)] flex items-center justify-center text-white font-mono text-sm"><div className="loader"></div></div>;

    if (!authUser) return (
        <div className="min-h-screen bg-[var(--color-body-bg)] flex items-center justify-center text-white flex-col gap-6">
            <User size={48} className="text-gray-600 mb-2" />
            <p className="text-xl font-light text-gray-400">Please login to view your profile.</p>
            <Link to="/login" className="px-8 py-3 bg-[var(--color-primary)] text-black font-bold rounded-full hover:bg-[#a6e6ce] transition-colors">LOGIN NOW</Link>
        </div>
    );

    if (!profileData) return <div className="min-h-screen bg-[var(--color-body-bg)] flex items-center justify-center text-white">Profile not found.</div>;

    const { user, stats, favorites, recent_activity } = profileData;

    // Calculate rating stats
    // Note: backend now sends rating_distribution in stats
    const avgRating = stats?.rating_distribution ?
        Object.entries(stats.rating_distribution).reduce((acc, [rating, count]) => acc + (parseFloat(rating) * count), 0) /
        (stats.reviews || 1) : 0; // rough approximation or use backend logic if provided

    return (
        <div className="min-h-screen bg-[var(--color-body-bg)] pb-20">
            <ProfileHeader user={user} />

            <div className="max-w-7xl mx-auto px-4 md:px-8">
                {/* Stats Grid */}
                <div className="grid grid-cols-2 md:grid-cols-5 gap-4 mb-16">
                    <StatsCard label="Total Films" value={stats?.total_films ? stats.total_films.toLocaleString() : 0} />
                    <StatsCard label="Movies This Year" value={stats?.movies_this_year || 0} trend="2026" />
                    <StatsCard label="Reviews" value={stats?.reviews || 0} />
                    <StatsCard label="Following" value={stats?.following || 0} />
                    <StatsCard label="Followers" value={stats?.followers || 0} />
                </div>

                {/* Tabs */}
                <div className="flex items-center justify-center md:justify-start gap-12 border-b border-white/5 mb-12 overflow-x-auto pb-1 scrollbar-hide">
                    {['Profile', 'Activity', 'Films', 'Reviews', 'Diary', 'Watchlist', 'Lists'].map((tab, i) => (
                        <button
                            key={tab}
                            className={`pb-4 text-xs font-bold tracking-[0.15em] uppercase transition-all relative
                                ${i === 0 ? 'text-[var(--color-primary)]' : 'text-gray-500 hover:text-white'}`}
                        >
                            {tab}
                            {i === 0 && <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-[var(--color-primary)] shadow-[0_0_10px_var(--color-primary)]" />}
                        </button>
                    ))}
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-3 gap-12">
                    {/* Left Column (Main Content) */}
                    <div className="lg:col-span-2 flex flex-col gap-12">
                        <section>
                            <SectionHeader title="Favorite Films" action="Manage Favorites" />
                            <MovieGrid movies={favorites} />
                        </section>

                        <section>
                            <SectionHeader title="Recent Activity" />
                            <MovieGrid movies={recent_activity} />
                        </section>
                    </div>

                    {/* Right Column (Sidebar) */}
                    <div className="flex flex-col gap-8">
                        <div className="glass-panel p-8 rounded-2xl">
                            <h3 className="font-bold text-white mb-6 uppercase tracking-wider text-sm">Rating Distribution</h3>
                            <div className="aspect-square bg-white/5 rounded-full flex items-center justify-center border border-white/5 relative group">
                                <div className="absolute inset-0 rounded-full border-[12px] border-[var(--color-primary)]/20 border-t-[var(--color-primary)] transform -rotate-45" />
                                <div className="text-center">
                                    <span className="text-5xl font-bold text-white tracking-tighter">
                                        {avgRating > 0 ? avgRating.toFixed(1) : '0.0'}
                                    </span>
                                    <p className="text-[10px] uppercase tracking-widest text-gray-400 mt-2 font-bold">Avg Rating</p>
                                </div>
                            </div>
                        </div>

                        <div className="glass-panel p-8 rounded-2xl">
                            <h3 className="font-bold text-white mb-6 uppercase tracking-wider text-sm">Following Activity</h3>
                            <div className="space-y-6">
                                <div className="flex gap-4">
                                    <div className="w-10 h-10 rounded-full bg-blue-500/20 flex items-center justify-center text-blue-400 flex-shrink-0">
                                        <Film size={18} />
                                    </div>
                                    <div className="text-sm leading-relaxed">
                                        <span className="font-bold text-white">Alex_Reels</span>
                                        <span className="text-gray-400"> liked your review of </span>
                                        <Link to="#" className="text-[var(--color-primary)] hover:underline">Dune</Link>
                                        <p className="text-xs text-gray-600 mt-1 font-medium">2 HOURS AGO</p>
                                    </div>
                                </div>
                                <div className="flex gap-4">
                                    <div className="w-10 h-10 rounded-full bg-purple-500/20 flex items-center justify-center text-purple-400 flex-shrink-0">
                                        <User size={18} />
                                    </div>
                                    <div className="text-sm leading-relaxed">
                                        <span className="font-bold text-white">Sara_Cine</span>
                                        <span className="text-gray-400"> started following you</span>
                                        <p className="text-xs text-gray-600 mt-1 font-medium">5 HOURS AGO</p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Profile;
