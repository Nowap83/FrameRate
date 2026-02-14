import { Link } from "react-router-dom";

const Footer = () => {
    const currentYear = new Date().getFullYear();

    return (
        <footer className="bg-header-bg py-8 px-6 mt-auto border-t border-white/5">
            <div className="max-w-7xl mx-auto flex flex-col md:flex-row justify-between items-center gap-4 text-xs text-gray-400">
                <div className="text-center md:text-left">
                    <p>© {currentYear} FrameRate. All rights reserved. Cinematic metadata courtesy of TMDB.</p>
                </div>

                <nav className="flex items-center gap-6">
                    <Link to="/about" className="hover:text-white transition-colors">About</Link>
                    <Link to="/terms" className="hover:text-white transition-colors">Terms</Link>
                    <Link to="/privacy" className="hover:text-white transition-colors">Privacy</Link>
                    <Link to="/support" className="hover:text-white transition-colors">Support</Link>
                </nav>
            </div>
        </footer>
    );
};

export default Footer;
