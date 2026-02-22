import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { X, Plus, Trash2, ChevronLeft, Save, Lock, User, Upload } from 'lucide-react';
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import apiClient from '../api/apiClient';
import { useAuth } from '../context/AuthContext';
import { getAvatarUrl } from '../utils/image';
import useDocumentTitle from '../hooks/useDocumentTitle';

import MovieSearch from '../components/MovieSearch';
import Input from '../components/Input';
import { changePasswordSchema } from '../validators/auth';
const Settings = () => {
    const navigate = useNavigate();
    const { user: authUser, loading: authLoading, setUser } = useAuth();
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [activeTab, setActiveTab] = useState('profile'); // 'profile' | 'auth' | 'avatar'
    const [passwordError, setPasswordError] = useState('');
    const [passwordSuccess, setPasswordSuccess] = useState('');
    const [avatarError, setAvatarError] = useState('');
    const [avatarSuccess, setAvatarSuccess] = useState('');
    const [uploadingAvatar, setUploadingAvatar] = useState(false);
    const [checkingUsername, setCheckingUsername] = useState(false);
    const [usernameAvailable, setUsernameAvailable] = useState(null); // null, true, 'taken', 'current'

    useDocumentTitle("Settings");

    // password form
    const {
        register,
        handleSubmit,
        formState: { errors },
        reset
    } = useForm({
        resolver: zodResolver(changePasswordSchema),
        defaultValues: {
            current_password: '',
            new_password: '',
            confirm_new_password: ''
        }
    });

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
        if (name === 'username') {
            setUsernameAvailable(null);
        }
    };

    const checkUsernameAvailability = async () => {
        if (!formData.username || formData.username === authUser.username) {
            setUsernameAvailable('current');
            return;
        }

        setCheckingUsername(true);
        try {
            const response = await apiClient.get(`/users/check-username?username=${encodeURIComponent(formData.username)}`);
            setUsernameAvailable(response.data.available);
        } catch (error) {
            console.error("Failed to check username", error);
        } finally {
            setCheckingUsername(false);
        }
    };

    const selectMovie = (movie) => {
        if (activeSlot !== null) {
            const newFavs = [...favorites];
            newFavs[activeSlot] = movie;
            setFavorites(newFavs);
            setActiveSlot(null);
        }
    };

    const removeMovie = (index) => {
        const newFavs = [...favorites];
        newFavs[index] = null;
        setFavorites(newFavs);
    };

    const handleProfileSubmit = async (e) => {
        e.preventDefault();
        setSaving(true);
        try {
            // Start with a copy of form data
            const cleanFormData = { ...formData };

            // We do NOT update profile_picture_url via this form anymore (use Avatar tab)
            delete cleanFormData.profile_picture_url;

            // Remove empty strings for optional fields to avoid validation errors
            if (!cleanFormData.website) delete cleanFormData.website;
            if (!cleanFormData.location) delete cleanFormData.location;
            if (!cleanFormData.given_name) delete cleanFormData.given_name;
            if (!cleanFormData.family_name) delete cleanFormData.family_name;

            const payload = {
                ...cleanFormData,
                favorite_films: favorites
                    .filter(m => m !== null)
                    .map(m => ({
                        tmdb_id: m.tmdb_id || m.id,
                        title: m.title,
                        poster_url: m.poster_path || m.poster_url,
                        release_year: m.release_date ? parseInt(m.release_date.split('-')[0]) : 0,
                    }))
            };

            await apiClient.put('/users/me', payload);
            navigate('/profile');
        } catch (error) {
            console.error("Update failed", error);
            setSaving(false);
        }
    };

    const onSubmitPassword = async (data) => {
        setSaving(true);
        setPasswordError('');
        setPasswordSuccess('');
        try {
            await apiClient.put('/users/me/password', {
                current_password: data.current_password,
                new_password: data.new_password
            });
            setPasswordSuccess('Password updated successfully');
            reset();
        } catch (error) {
            console.error("Password update failed", error);
            setPasswordError(error.response?.data?.error || 'Failed to update password');
        } finally {
            setSaving(false);
        }
    };

    const handleAvatarUpload = async (e) => {
        const file = e.target.files[0];
        if (!file) return;

        setUploadingAvatar(true);
        setAvatarError('');
        setAvatarSuccess('');

        const formData = new FormData();
        formData.append('avatar', file);

        try {
            const response = await apiClient.post('/users/me/avatar', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });

            if (authUser) {
                setUser({ ...authUser, profile_picture_url: response.data.url });
            }
            setFormData(prev => ({ ...prev, profile_picture_url: response.data.url }));

            setAvatarSuccess('Avatar updated successfully!');
        } catch (error) {
            console.error("Avatar upload failed", error);
            setAvatarError(error.response?.data?.error || 'Failed to upload avatar');
        } finally {
            setUploadingAvatar(false);
        }
    };

    if (loading) return <div className="min-h-screen bg-[var(--color-body-bg)] flex items-center justify-center text-white font-mono text-sm"><div className="loader"></div></div>;

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
                    {activeTab === 'profile' && (
                        <button
                            type="submit"
                            form="settings-form"
                            disabled={saving}
                            className="flex items-center gap-2 px-6 py-2.5 bg-[var(--color-primary)] text-black font-bold rounded-full hover:bg-[#a6e6ce] transition-all hover:scale-105 shadow-[0_0_20px_rgba(189,244,223,0.3)] disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            <Save size={18} />
                            {saving ? 'Saving...' : 'Save Changes'}
                        </button>
                    )}
                </div>

                {/* Tabs */}
                <div className="flex gap-4 mb-8 border-b border-white/5 pb-1">
                    <button
                        onClick={() => setActiveTab('profile')}
                        className={`flex items-center gap-2 px-4 py-2 text-sm font-bold uppercase tracking-wider transition-colors relative ${activeTab === 'profile' ? 'text-[var(--color-primary)]' : 'text-gray-500 hover:text-white'}`}
                    >
                        <User size={16} />
                        Profile
                        {activeTab === 'profile' && <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-[var(--color-primary)] shadow-[0_0_10px_var(--color-primary)]" />}
                    </button>
                    <button
                        onClick={() => setActiveTab('auth')}
                        className={`flex items-center gap-2 px-4 py-2 text-sm font-bold uppercase tracking-wider transition-colors relative ${activeTab === 'auth' ? 'text-[var(--color-primary)]' : 'text-gray-500 hover:text-white'}`}
                    >
                        <Lock size={16} />
                        Auth
                        {activeTab === 'auth' && <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-[var(--color-primary)] shadow-[0_0_10px_var(--color-primary)]" />}
                    </button>
                    <button
                        onClick={() => setActiveTab('avatar')}
                        className={`flex items-center gap-2 px-4 py-2 text-sm font-bold uppercase tracking-wider transition-colors relative ${activeTab === 'avatar' ? 'text-[var(--color-primary)]' : 'text-gray-500 hover:text-white'}`}
                    >
                        <User size={16} />
                        Avatar
                        {activeTab === 'avatar' && <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-[var(--color-primary)] shadow-[0_0_10px_var(--color-primary)]" />}
                    </button>
                </div>

                <div className="bg-[#121212] rounded-2xl border border-white/10">

                    {/* PROFILE TAB */}
                    {activeTab === 'profile' && (
                        <form id="settings-form" onSubmit={handleProfileSubmit} className="p-6 md:p-8 space-y-10">

                            {/* 1. Profile Information */}
                            <div className="space-y-6">
                                <h3 className="text-sm font-bold text-[var(--color-primary)] uppercase tracking-wider border-b border-white/5 pb-2">Profile</h3>

                                <div className="space-y-1.5">
                                    <label className="text-xs text-gray-400 font-bold uppercase">Username</label>
                                    <div className="flex gap-2">
                                        <input
                                            type="text"
                                            name="username"
                                            value={formData.username}
                                            onChange={handleInputChange}
                                            className={`w-full bg-white/5 border rounded-lg p-3 text-white focus:outline-none transition-colors ${usernameAvailable === true ? 'border-green-500/50 focus:border-green-500' :
                                                usernameAvailable === false ? 'border-red-500/50 focus:border-red-500' :
                                                    'border-white/10 focus:border-[var(--color-primary)]'
                                                }`}
                                        />
                                        <button
                                            type="button"
                                            onClick={checkUsernameAvailability}
                                            disabled={checkingUsername || !formData.username || formData.username === authUser.username}
                                            className="px-4 py-2 bg-white/5 hover:bg-white/10 text-white rounded-lg text-sm font-bold uppercase disabled:opacity-50 transition-colors"
                                        >
                                            {checkingUsername ? '...' : 'Check'}
                                        </button>
                                    </div>
                                    {usernameAvailable === true && <p className="text-xs text-green-500">Username is available!</p>}
                                    {usernameAvailable === false && <p className="text-xs text-red-500">Username is already taken.</p>}
                                    {usernameAvailable === 'current' && <p className="text-xs text-gray-500">This is your current username.</p>}
                                </div>

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
                                    <div className="mt-4">
                                        <MovieSearch
                                            onSelect={(movie) => selectMovie({
                                                tmdb_id: movie.id,
                                                title: movie.title,
                                                poster_path: movie.poster_path, // Ensure consistency with API response
                                                poster_url: movie.poster_path,   // For backward compatibility if needed
                                                release_date: movie.release_date
                                            })}
                                            onCancel={() => setActiveSlot(null)}
                                        />
                                    </div>
                                )}
                            </div>

                        </form>
                    )}

                    {/* AUTH TAB */}
                    {activeTab === 'auth' && (
                        <div className="p-6 md:p-8 space-y-10">
                            <div className="space-y-6">
                                <h3 className="text-sm font-bold text-[var(--color-primary)] uppercase tracking-wider border-b border-white/5 pb-2">Change Password</h3>

                                {passwordError && (
                                    <div className="p-4 rounded-lg text-sm border bg-red-500/10 border-red-500/20 text-red-400">
                                        {passwordError}
                                    </div>
                                )}
                                {passwordSuccess && (
                                    <div className="p-4 rounded-lg text-sm border bg-green-500/10 border-green-500/20 text-green-400">
                                        {passwordSuccess}
                                    </div>
                                )}

                                <form onSubmit={handleSubmit(onSubmitPassword)} className="space-y-6 max-w-md">
                                    <Input
                                        label="Current Password"
                                        type="password"
                                        placeholder="••••••••"
                                        {...register("current_password")}
                                        error={errors.current_password?.message}
                                        disabled={saving}
                                    />
                                    <Input
                                        label="New Password"
                                        type="password"
                                        placeholder="••••••••"
                                        {...register("new_password")}
                                        error={errors.new_password?.message}
                                        disabled={saving}
                                    />
                                    <Input
                                        label="Confirm New Password"
                                        type="password"
                                        placeholder="••••••••"
                                        {...register("confirm_new_password")}
                                        error={errors.confirm_new_password?.message}
                                        disabled={saving}
                                    />

                                    <button
                                        type="submit"
                                        disabled={saving}
                                        className="flex items-center gap-2 px-6 py-2.5 bg-[var(--color-primary)] text-black font-bold rounded-full hover:bg-[#a6e6ce] transition-all hover:scale-105 shadow-[0_0_20px_rgba(189,244,223,0.3)] disabled:opacity-50 disabled:cursor-not-allowed"
                                    >
                                        {saving ? 'Updating...' : 'Update Password'}
                                    </button>
                                </form>
                            </div>
                        </div>
                    )}

                    {/* AVATAR TAB */}
                    {activeTab === 'avatar' && (
                        <div className="p-6 md:p-8 space-y-10">
                            <div className="space-y-6">
                                <h3 className="text-sm font-bold text-[var(--color-primary)] uppercase tracking-wider border-b border-white/5 pb-2">Change Avatar</h3>

                                {avatarError && (
                                    <div className="p-4 rounded-lg text-sm border bg-red-500/10 border-red-500/20 text-red-400">
                                        {avatarError}
                                    </div>
                                )}
                                {avatarSuccess && (
                                    <div className="p-4 rounded-lg text-sm border bg-green-500/10 border-green-500/20 text-green-400">
                                        {avatarSuccess}
                                    </div>
                                )}

                                <div className="flex flex-col items-center gap-8 py-8">
                                    <div className="relative group w-40 h-40 rounded-full overflow-hidden border-2 border-dashed border-white/20 hover:border-[var(--color-primary)] transition-colors bg-white/5 flex items-center justify-center">
                                        {authUser?.profile_picture_url ? (
                                            <img
                                                src={getAvatarUrl(authUser.profile_picture_url)}
                                                alt="Avatar"
                                                className="w-full h-full object-cover"
                                            />
                                        ) : (
                                            <User size={48} className="text-gray-500" />
                                        )}

                                        <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col items-center justify-center cursor-pointer">
                                            <Upload size={24} className="text-white mb-2" />
                                            <span className="text-xs text-white font-bold uppercase">Upload</span>
                                        </div>

                                        <input
                                            type="file"
                                            accept="image/*"
                                            onChange={handleAvatarUpload}
                                            disabled={uploadingAvatar}
                                            className="absolute inset-0 w-full h-full opacity-0 cursor-pointer disabled:cursor-not-allowed"
                                        />
                                    </div>

                                    <div className="text-center space-y-2">
                                        <p className="text-sm text-gray-400">
                                            Click on the image to upload a new avatar.
                                        </p>
                                        <p className="text-xs text-gray-600 uppercase font-bold tracking-wider">
                                            Max size 5MB. Formats: JPEG, PNG, WEBP.
                                        </p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    )}


                </div>
            </div>
        </div>
    );
};

export default Settings;
