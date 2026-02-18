import { Link } from "react-router-dom";
import { Film } from "lucide-react";

const Footer = () => {
    const currentYear = new Date().getFullYear();

    return (
        <footer className="bg-emerald-dark pt-12 pb-8 px-6 mt-auto border-t border-white/5 relative">
            {/* Subtle top gradient line */}
            <div className="absolute top-0 left-0 right-0 h-[1px] bg-gradient-to-r from-transparent via-mint/20 to-transparent" />

            <div className="max-w-7xl mx-auto flex flex-col md:flex-row justify-between items-center gap-6">
                <div className="flex flex-col md:flex-row items-center gap-2 md:gap-6 text-center md:text-left">
                    <div className="flex items-center gap-2 text-white/90 font-display font-medium">
                        <Film className="w-4 h-4 text-mint" />
                        <span>FrameRate</span>
                    </div>
                    <span className="hidden md:block text-white/10">|</span>
                    <p className="text-xs text-gray-500">
                        © {currentYear} FrameRate. <span className="hidden sm:inline">Cinematic metadata courtesy of TMDB.</span>
                    </p>
                </div>

                <nav className="flex items-center gap-8">
                    {["About", "Terms", "Privacy", "Support"].map((item) => (
                        <Link
                            key={item}
                            to={`/${item.toLowerCase()}`}
                            className="text-[10px] uppercase tracking-widest text-gray-400 hover:text-mint transition-all duration-300 font-medium"
                        >
                            {item}
                        </Link>
                    ))}
                </nav>
            </div>
        </footer>
    );
};

export default Footer;
