import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { X, Search, Plus, Trash2, ChevronLeft, Save } from 'lucide-react';
import apiClient from '../api/apiClient';
import { useAuth } from '../context/AuthContext';

const Settings = () => {
    const navigate = useNavigate();
    const { user: authUser, loading: authLoading } = useAuth();
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);

    const [formData, setFormData] = useState({
        username: '',
        bio: '',
        given_name: '',
        family_name: '',
        location: '',
        website: '',
        profile_picture_url: '',
    });

    const [favorites, setFavorites] = useState(Array(4).fill(null));
    const [searchQuery, setSearchQuery] = useState('');
    const [searchResults, setSearchResults] = useState([]);
    const [activeSlot, setActiveSlot] = useState(null);

    useEffect(() => {
        const fetchProfile = async () => {
            try {
                const response = await apiClient.get('/users/me');
                const userData = response.data.user;
                const userFavs = response.data.favorites;

                setFormData({
                    username: userData.username || '',
                    bio: userData.bio || '',
                    given_name: userData.given_name || '',
                    family_name: userData.family_name || '',
                    location: userData.location || '',
                    website: userData.website || '',
                    profile_picture_url: userData.profile_picture_url || '',
                });

                if (userFavs && userFavs.length > 0) {
                    const initialFavs = Array(4).fill(null);
                    userFavs.slice(0, 4).forEach((movie, index) => {
                        initialFavs[index] = movie;
                    });
                    setFavorites(initialFavs);
                }
            } catch (error) {
                console.error("Failed to fetch profile:", error);
            } finally {
                setLoading(false);
            }
        };

        if (!authLoading && authUser) {
            fetchProfile();
        } else if (!authLoading && !authUser) {
            navigate('/login');
        }
    }, [authLoading, authUser, navigate]);

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));
    };

    const handleSearch = async (query) => {
        setSearchQuery(query);
        if (query.length > 2) {
            try {
                const response = await apiClient.get(`/tmdb/search/movie?query=${encodeURIComponent(query)}`);
                setSearchResults(response.data.results || []);
            } catch (error) {
                console.error("Search failed", error);
            }
        } else {
            setSearchResults([]);
        }
    };

    const selectMovie = (movie) => {
        if (activeSlot !== null) {
            const newFavs = [...favorites];
            newFavs[activeSlot] = movie;
            setFavorites(newFavs);
            setActiveSlot(null);
            setSearchQuery('');
            setSearchResults([]);
        }
    };

    const removeMovie = (index) => {
        const newFavs = [...favorites];
        newFavs[index] = null;
        setFavorites(newFavs);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setSaving(true);
        try {
            const payload = {
                ...formData,
                favorite_films: favorites
                    .filter(m => m !== null)
                    .map(m => ({
                        tmdb_id: m.tmdb_id || m.id,
                        title: m.title,
                        poster_url: m.poster_path || m.poster_url,
                        release_year: m.release_date ? parseInt(m.release_date.split('-')[0]) : 0,
                        // Add defaults for fields we might miss from search but are needed for model
                        // Note: Backend might need more, but UpsertMovie handles updates.
                        // Ideally we should try to get these from search result if available.
                    }))
            };

            await apiClient.put('/users/me', payload);
            navigate('/profile');
        } catch (error) {
            console.error("Update failed", error);
            setSaving(false);
        }
    };

    if (loading) return <div className="min-h-screen bg-[var(--color-body-bg)] flex items-center justify-center text-white font-mono text-sm animate-pulse">LOADING SETTINGS...</div>;

    return (
        <div className="min-h-screen bg-[var(--color-body-bg)] pb-20 pt-20">
            <div className="max-w-4xl mx-auto px-4 md:px-8">

                {/* Header */}
                <div className="flex items-center justify-between mb-8">
                    <div className="flex items-center gap-4">
                        <button onClick={() => navigate('/profile')} className="p-2 hover:bg-white/5 rounded-full text-gray-400 hover:text-white transition-colors">
                            <ChevronLeft size={24} />
                        </button>
                        <div>
                            <h1 className="text-3xl font-bold text-white font-display">Settings</h1>
                            <p className="text-gray-400 text-sm">Manage your profile and preferences.</p>
                        </div>
                    </div>
                    <button
                        type="submit"
                        form="settings-form"
                        disabled={saving}
                        className="flex items-center gap-2 px-6 py-2.5 bg-[var(--color-primary)] text-black font-bold rounded-full hover:bg-[#a6e6ce] transition-all hover:scale-105 shadow-[0_0_20px_rgba(189,244,223,0.3)] disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        <Save size={18} />
                        {saving ? 'Saving...' : 'Save Changes'}
                    </button>
                </div>

                <div className="bg-[#121212] rounded-2xl border border-white/10 overflow-hidden">
                    <form id="settings-form" onSubmit={handleSubmit} className="p-6 md:p-8 space-y-10">

                        {/* 1. Profile Information */}
                        <div className="space-y-6">
                            <h3 className="text-sm font-bold text-[var(--color-primary)] uppercase tracking-wider border-b border-white/5 pb-2">Profile Information</h3>

                            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                <div className="space-y-1.5">
                                    <label className="text-xs text-gray-400 font-bold uppercase">Given Name</label>
                                    <input
                                        type="text"
                                        name="given_name"
                                        value={formData.given_name}
                                        onChange={handleInputChange}
                                        className="w-full bg-white/5 border border-white/10 rounded-lg p-3 text-white focus:border-[var(--color-primary)] focus:outline-none transition-colors"
                                    />
                                </div>
                                <div className="space-y-1.5">
                                    <label className="text-xs text-gray-400 font-bold uppercase">Family Name</label>
                                    <input
                                        type="text"
                                        name="family_name"
                                        value={formData.family_name}
                                        onChange={handleInputChange}
                                        className="w-full bg-white/5 border border-white/10 rounded-lg p-3 text-white focus:border-[var(--color-primary)] focus:outline-none transition-colors"
                                    />
                                </div>
                            </div>

                            <div className="space-y-1.5">
                                <label className="text-xs text-gray-400 font-bold uppercase">Location</label>
                                <input
                                    type="text"
                                    name="location"
                                    value={formData.location}
                                    onChange={handleInputChange}
                                    placeholder="e.g. Paris, France"
                                    className="w-full bg-white/5 border border-white/10 rounded-lg p-3 text-white focus:border-[var(--color-primary)] focus:outline-none transition-colors"
                                />
                            </div>

                            <div className="space-y-1.5">
                                <label className="text-xs text-gray-400 font-bold uppercase">Website</label>
                                <input
                                    type="text"
                                    name="website"
                                    value={formData.website}
                                    onChange={handleInputChange}
                                    placeholder="https://"
                                    className="w-full bg-white/5 border border-white/10 rounded-lg p-3 text-white focus:border-[var(--color-primary)] focus:outline-none transition-colors"
                                />
                            </div>

                            <div className="space-y-1.5">
                                <label className="text-xs text-gray-400 font-bold uppercase">Profile Picture URL</label>
                                <input
                                    type="text"
                                    name="profile_picture_url"
                                    value={formData.profile_picture_url}
                                    onChange={handleInputChange}
                                    placeholder="https://"
                                    className="w-full bg-white/5 border border-white/10 rounded-lg p-3 text-white focus:border-[var(--color-primary)] focus:outline-none transition-colors"
                                />
                                <p className="text-[10px] text-gray-500">Leave empty to use generated avatar.</p>
                            </div>

                            <div className="space-y-1.5">
                                <label className="text-xs text-gray-400 font-bold uppercase">Bio</label>
                                <textarea
                                    name="bio"
                                    value={formData.bio}
                                    onChange={handleInputChange}
                                    rows={4}
                                    className="w-full bg-white/5 border border-white/10 rounded-lg p-3 text-white focus:border-[var(--color-primary)] focus:outline-none transition-colors resize-none"
                                />
                            </div>
                        </div>

                        {/* 2. Favorite Films */}
                        <div className="space-y-6">
                            <div className="flex justify-between items-end border-b border-white/5 pb-2">
                                <h3 className="text-sm font-bold text-[var(--color-primary)] uppercase tracking-wider">Favorite Films</h3>
                                <p className="text-xs text-gray-500">Select your top 4 favorites</p>
                            </div>

                            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                                {favorites.map((movie, index) => (
                                    <div key={index} className="relative aspect-[2/3] bg-white/5 rounded-xl border border-white/10 overflow-hidden group hover:border-[var(--color-primary)]/50 transition-colors">
                                        {movie ? (
                                            <>
                                                <img
                                                    src={`https://image.tmdb.org/t/p/w300${movie.poster_path || movie.poster_url}`}
                                                    alt={movie.title}
                                                    className="w-full h-full object-cover"
                                                />
                                                <button
                                                    type="button"
                                                    onClick={() => removeMovie(index)}
                                                    className="absolute top-2 right-2 p-1.5 bg-black/60 rounded-full text-white hover:bg-red-500/80 transition-colors opacity-0 group-hover:opacity-100"
                                                >
                                                    <Trash2 size={14} />
                                                </button>
                                                <button
                                                    type="button"
                                                    onClick={() => setActiveSlot(index)} // Allow replacing
                                                    className="absolute inset-0 bg-black/20 opacity-0 group-hover:opacity-100 transition-opacity"
                                                />
                                            </>
                                        ) : (
                                            <button
                                                type="button"
                                                onClick={() => setActiveSlot(index)}
                                                className="w-full h-full flex flex-col items-center justify-center text-gray-500 hover:text-[var(--color-primary)] transition-colors gap-2"
                                            >
                                                <Plus size={32} />
                                                <span className="text-xs font-bold uppercase">Add Film</span>
                                            </button>
                                        )}
                                    </div>
                                ))}
                            </div>

                            {/* Movie Search Overlay - Inline/Modal based on slot active */}
                            {activeSlot !== null && (
                                <div className="mt-4 p-6 bg-[#1a1a1a] rounded-xl border border-white/10 animate-in fade-in slide-in-from-top-4">
                                    <div className="flex items-center gap-3 bg-white/5 p-3 rounded-lg border border-white/10 mb-4">
                                        <Search size={20} className="text-gray-400" />
                                        <input
                                            type="text"
                                            autoFocus
                                            placeholder="Search for a movie..."
                                            value={searchQuery}
                                            onChange={(e) => handleSearch(e.target.value)}
                                            className="bg-transparent border-none focus:outline-none text-white w-full"
                                        />
                                        <button type="button" onClick={() => setActiveSlot(null)} className="text-xs text-gray-500 hover:text-white">CANCEL</button>
                                    </div>

                                    <div className="space-y-2 max-h-60 overflow-y-auto custom-scrollbar">
                                        {searchResults.map(result => (
                                            <button
                                                key={result.id}
                                                type="button"
                                                onClick={() => selectMovie({
                                                    tmdb_id: result.id,
                                                    title: result.title,
                                                    poster_path: result.poster_path, // Keep this for UI
                                                    poster_url: result.poster_path,  // For backend consistency
                                                    release_date: result.release_date
                                                })}
                                                className="w-full flex items-center gap-3 p-2 hover:bg-white/10 rounded-lg transition-colors text-left"
                                            >
                                                <div className="w-10 h-14 bg-gray-800 rounded overflow-hidden flex-shrink-0">
                                                    {result.poster_path && (
                                                        <img src={`https://image.tmdb.org/t/p/w92${result.poster_path}`} alt="" className="w-full h-full object-cover" />
                                                    )}
                                                </div>
                                                <div>
                                                    <p className="font-bold text-white text-sm">{result.title}</p>
                                                    <p className="text-xs text-gray-400">
                                                        {result.release_date ? result.release_date.split('-')[0] : 'Unknown'}
                                                    </p>
                                                </div>
                                            </button>
                                        ))}
                                        {searchQuery.length > 2 && searchResults.length === 0 && (
                                            <p className="text-center text-gray-500 text-sm py-4">No movies found.</p>
                                        )}
                                    </div>
                                </div>
                            )}
                        </div>

                    </form>
                </div>
            </div>
        </div>
    );
};

export default Settings;
