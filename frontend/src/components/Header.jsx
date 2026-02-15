import { Link } from "react-router-dom";
import Button from "./Button";
import SearchBar from "./SearchBar";

const Header = () => {
    return (
        <header className="bg-header-bg py-4 px-6 flex items-center justify-between shadow-lg">
            <div className="flex items-center gap-8">
                {/* logo mock */}
                <Link to="/" className="flex items-center gap-2 text-white font-bold text-xl">
                    <div className="text-mint">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8" viewBox="0 0 24 24" fill="currentColor">
                            <path d="M18 4l2 4h-3l-2-4h-2l2 4h-3l-2-4H8l2 4H7L5 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V4h-4zM7 18H5v-2h2v2zm0-4H5v-2h2v2zm0-4H5V8h2v2zm4 8H9v-2h2v2zm0-4H9v-2h2v2zm0-4H9V8h2v2zm4 8h-2v-2h2v2zm0-4h-2v-2h2v2zm0-4h-2V8h2v2zm4 8h-2v-2h2v2zm0-4h-2v-2h2v2zm0-4h-2V8h2v2z" />
                        </svg>
                    </div>
                    FrameRate
                </Link>

                {/* nav */}
                <nav className="hidden md:flex items-center gap-6">
                    <Link to="/movies" className="text-gray-400 hover:text-white transition-colors text-sm font-medium">Films</Link>
                    <Link to="/lists" className="text-gray-400 hover:text-white transition-colors text-sm font-medium">Lists</Link>
                    <Link to="/community" className="text-gray-400 hover:text-white transition-colors text-sm font-medium">Community</Link>
                </nav>
            </div>

            <div className="flex items-center gap-6 flex-1 max-w-2xl ml-8">
                <SearchBar className="flex-1" />
            </div>

            <div className="flex items-center gap-3">
                <Button className="bg-opacity-80 px-4 py-2 text-sm">Sign In</Button>
                <Button className="px-4 py-2 text-sm">Register</Button>
            </div>
        </header>
    );
};

export default Header;