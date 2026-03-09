import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import Button from "./Button";
import { LogOut, User, Menu, Search, ShieldAlert, X } from "lucide-react";
import { getAvatarUrl } from "../utils/image";
import { useState } from "react";

const Header = () => {
    const { user, logout } = useAuth();
    const navigate = useNavigate();
    const [searchQuery, setSearchQuery] = useState("");
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

    const handleLogout = () => {
        logout();
        navigate("/");
    };

    const handleSearchSubmit = (e) => {
        e.preventDefault();
        if (searchQuery.trim()) {
            navigate(`/search?q=${encodeURIComponent(searchQuery.trim())}`);
            setSearchQuery("");
        }
    };

    return (
        <header className="bg-[#12201B]/90 backdrop-blur-md sticky top-0 z-50 py-4 px-6 border-b border-white/5">
            <div className="container mx-auto flex items-center justify-between">

                <div className="flex items-center gap-8">
                    {/* logo */}
                    <Link to="/" className="flex items-center gap-2 font-bold text-xl tracking-tight">
                        <div className="w-8 h-8 bg-mint rounded-lg flex items-center justify-center text-[#12201B]">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M18 4l2 4h-3l-2-4h-2l2 4h-3l-2-4H8l2 4H7L5 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V4h-4z" />
                                <rect x="5" y="14" width="2" height="2" fill="currentColor" fillOpacity="0.3" />
                                <rect x="5" y="10" width="2" height="2" fill="currentColor" fillOpacity="0.3" />
                                <rect x="17" y="14" width="2" height="2" fill="currentColor" fillOpacity="0.3" />
                                <rect x="17" y="10" width="2" height="2" fill="currentColor" fillOpacity="0.3" />
                            </svg>
                        </div>
                        <span className="font-display text-white">FrameRate</span>
                    </Link>

                    {/* nav */}
                    <nav className="hidden md:flex items-center gap-6">
                        <Link to="/movies" className="text-gray-400 hover:text-white transition-colors text-sm font-medium">Films</Link>
                        <Link to="/lists" className="text-gray-400 hover:text-white transition-colors text-sm font-medium">Lists</Link>
                        <Link to="/community" className="text-gray-400 hover:text-white transition-colors text-sm font-medium">Community</Link>
                    </nav>
                </div>

                {/* search bar */}
                <div className="hidden md:block flex-1 max-w-xl px-8">
                    <form onSubmit={handleSearchSubmit} className="relative w-full">
                        <div className="relative flex items-center">
                            <span className="absolute left-3 text-gray-400">
                                <Search size={16} />
                            </span>
                            <input
                                type="text"
                                placeholder="Search movies..."
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                className="w-full bg-white/5 text-sm text-gray-300 py-2 pl-10 pr-4 rounded-full border border-white/10 focus:outline-none focus:border-[var(--color-primary)]/50 transition-colors"
                            />
                        </div>
                    </form>
                </div>

                {/* right actions */}
                <div className="flex items-center gap-4">
                    {user ? (
                        <>

                            <div className="hidden md:flex items-center gap-3">
                                {/* Admin Menu */}
                                {user.is_admin && (
                                    <Link to="/admin" className="p-2 text-gray-400 hover:text-[var(--color-mint)] hover:bg-[var(--color-mint)]/10 rounded-full transition-all" title="Admin Dashboard">
                                        <ShieldAlert size={20} />
                                    </Link>
                                )}

                                {/* User Menu */}
                                <Link to="/profile">
                                    <div className="flex items-center gap-2 cursor-pointer hover:bg-white/5 px-2 py-1 rounded-full transition-colors">
                                        <div className="w-8 h-8 rounded-full bg-gradient-to-br from-mint to-emerald-600 flex items-center justify-center text-[#12201B] font-bold text-xs ring-2 ring-[#12201B] overflow-hidden">
                                            {user.profile_picture_url ? (
                                                <img
                                                    src={getAvatarUrl(user.profile_picture_url)}
                                                    alt={user.username}
                                                    className="w-full h-full object-cover"
                                                />
                                            ) : (
                                                user.username?.charAt(0).toUpperCase()
                                            )}
                                        </div>
                                        <span className="text-sm font-medium hidden lg:block">{user.username}</span>
                                    </div>
                                </Link>


                                <button
                                    onClick={handleLogout}
                                    className="p-2 text-gray-400 hover:text-white hover:bg-red-500/10 hover:text-red-400 rounded-full transition-all"
                                    title="Logout"
                                >
                                    <LogOut size={20} />
                                </button>
                            </div>
                        </>
                    ) : (
                        <div className="flex items-center gap-3">
                            <Link to="/login" className="text-gray-300 hover:text-white font-medium text-xs sm:text-sm px-2 sm:px-3 py-2">
                                Sign In
                            </Link>
                            <Link to="/register">
                                <Button className="px-3 sm:px-5 py-1.5 sm:py-2 text-xs sm:text-sm rounded-full whitespace-nowrap">Get Started</Button>
                            </Link>
                        </div>
                    )}

                    {/* mobile menu button fixed layout */}
                    <div className="md:hidden flex items-center ml-4">
                        <button
                            className="p-2 text-gray-300 hover:text-white transition-colors"
                            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                        >
                            {isMobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
                        </button>
                    </div>
                </div>
            </div>

            {/* menu burger mobile */}
            {isMobileMenuOpen && (
                <div className="md:hidden absolute top-full left-0 w-full bg-[#12201B]/95 backdrop-blur-xl border-b border-white/10 py-6 px-6 flex flex-col gap-6 shadow-2xl animate-in slide-in-from-top-2">
                    <nav className="flex flex-col gap-4">
                        <Link to="/movies" className="text-gray-300 hover:text-white text-lg font-medium" onClick={() => setIsMobileMenuOpen(false)}>Films</Link>
                        <Link to="/lists" className="text-gray-300 hover:text-white text-lg font-medium" onClick={() => setIsMobileMenuOpen(false)}>Lists</Link>
                        <Link to="/community" className="text-gray-300 hover:text-white text-lg font-medium" onClick={() => setIsMobileMenuOpen(false)}>Community</Link>
                    </nav>

                    {user ? (
                        <div className="flex flex-col gap-4 border-t border-white/10 pt-6">
                            <Link to="/profile" className="flex items-center gap-3 text-white" onClick={() => setIsMobileMenuOpen(false)}>
                                <div className="w-10 h-10 rounded-full bg-gradient-to-br from-mint to-emerald-600 flex items-center justify-center text-[#12201B] font-bold text-sm overflow-hidden">
                                    {user.profile_picture_url ? (
                                        <img src={getAvatarUrl(user.profile_picture_url)} alt={user.username} className="w-full h-full object-cover" />
                                    ) : (
                                        user.username?.charAt(0).toUpperCase()
                                    )}
                                </div>
                                <span className="font-bold">{user.username}</span>
                            </Link>

                            {user.is_admin && (
                                <Link to="/admin" className="flex items-center gap-3 text-gray-300 hover:text-mint transition-colors" onClick={() => setIsMobileMenuOpen(false)}>
                                    <ShieldAlert size={20} />
                                    <span>Admin Dashboard</span>
                                </Link>
                            )}

                            <button
                                onClick={() => {
                                    handleLogout();
                                    setIsMobileMenuOpen(false);
                                }}
                                className="flex items-center gap-3 text-red-400 hover:text-red-300 transition-colors w-full text-left"
                            >
                                <LogOut size={20} />
                                <span>Logout</span>
                            </button>
                        </div>
                    ) : (
                        <div className="flex flex-col gap-4 border-t border-white/10 pt-6">
                            <Link to="/login" onClick={() => setIsMobileMenuOpen(false)}>
                                <Button className="w-full bg-white/10 hover:bg-white/20 text-white">Sign In</Button>
                            </Link>
                            <Link to="/register" onClick={() => setIsMobileMenuOpen(false)}>
                                <Button className="w-full">Get Started</Button>
                            </Link>
                        </div>
                    )}
                </div>
            )}
        </header>
    );
};

export default Header;