import { motion } from "framer-motion";
import { Search, PlayCircle } from "lucide-react";
import Button from "../components/Button";
import { Link } from "react-router-dom";

const LandingPage = () => {
    return (
        <div className="min-h-screen text-white selection:bg-mint selection:text-black">
            {/* hero section */}
            <section className="relative h-screen flex items-center justify-center overflow-hidden">
                {/* background gradient/mesh */}
                <div className="absolute inset-0 bg-gradient-to-b from-transparent to-[#12201B] z-10" />
                <div className="absolute inset-0 opacity-20 bg-[url('https://image.tmdb.org/t/p/original/uDgy6hyPd82kOHh6I95FLtLnj6p.jpg')] bg-cover bg-center animate-pulse-slow" />

                <div className="absolute inset-0 bg-[#0A1410]/60 z-0" />

                <div className="relative z-20 container mx-auto px-4 text-center">
                    <motion.h1
                        initial={{ opacity: 0, y: 30 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.8 }}
                        className="text-5xl md:text-7xl font-bold mb-6 font-display"
                    >
                        Track your <span className="text-mint">cinematic</span> <br /> journey.
                    </motion.h1>

                    <motion.p
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.8, delay: 0.2 }}
                        className="text-gray-300 text-lg md:text-xl max-w-2xl mx-auto mb-10"
                    >
                        The world's most beautiful social network for movie lovers. Rate, review, and discover your next favorite film.
                    </motion.p>

                    <motion.div
                        initial={{ opacity: 0, scale: 0.9 }}
                        animate={{ opacity: 1, scale: 1 }}
                        transition={{ duration: 0.5, delay: 0.4 }}
                        className="flex flex-col md:flex-row items-center justify-center gap-4"
                    >
                        <div className="relative w-full md:w-96">
                            <div className="absolute inset-y-0 left-3 flex items-center pointer-events-none">
                                <Search className="h-5 w-5 text-gray-400" />
                            </div>
                            <input
                                type="text"
                                placeholder="Search for a movie, actor, or director..."
                                className="w-full pl-10 pr-4 py-3 bg-white/10 backdrop-blur-md border border-white/10 rounded-full text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-mint focus:border-transparent transition-all"
                            />
                        </div>
                        <Button className="w-full md:w-auto px-8 py-3 rounded-full font-semibold">
                            Search
                        </Button>
                    </motion.div>

                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        transition={{ duration: 1, delay: 0.8 }}
                        className="mt-8 text-xs text-mint/60 uppercase tracking-widest"
                    >
                        Trending Searches: Dune 2   Furiosa   Civil War
                    </motion.div>
                </div>
            </section>

            {/* stats section */}
            <section className="py-20 bg-[#0A1410] border-y border-white/5">
                <div className="container mx-auto px-4">
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-8 text-center">
                        {[
                            { number: "42M+", label: "Films Watched" },
                            { number: "12M+", label: "Reviews Shared" },
                            { number: "8M+", label: "Cinephiles Joined" },
                        ].map((stat, index) => (
                            <motion.div
                                key={index}
                                initial={{ opacity: 0, y: 20 }}
                                whileInView={{ opacity: 1, y: 0 }}
                                viewport={{ once: true }}
                                transition={{ delay: index * 0.1 }}
                            >
                                <h3 className="text-4xl font-bold text-white mb-2 font-display">{stat.number}</h3>
                                <p className="text-gray-400 text-sm tracking-widest uppercase">{stat.label}</p>
                            </motion.div>
                        ))}
                    </div>
                </div>
            </section>

            {/* CTA / Community Section */}
            <section className="py-32 relative overflow-hidden">
                <div className="absolute inset-0 bg-gradient-to-t from-[#12201B] to-transparent z-10 pointer-events-none" />

                <div className="container mx-auto px-4 relative z-20 text-center">
                    <motion.div
                        initial={{ opacity: 0, scale: 0.95 }}
                        whileInView={{ opacity: 1, scale: 1 }}
                        viewport={{ once: true }}
                        className="bg-gradient-to-br from-white/5 to-white/0 backdrop-blur-3xl border border-white/10 rounded-3xl p-12 md:p-20 max-w-4xl mx-auto shadow-2xl"
                    >
                        <h2 className="text-3xl md:text-5xl font-bold mb-6 font-display">
                            Join the world's most <br /> passionate community of <br /> film lovers.
                        </h2>
                        <p className="text-gray-300 mb-10 max-w-xl mx-auto">
                            Log movies as you watch them. Save those you want to see. Reviews, lists, and recommendations from your friends.
                        </p>
                        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
                            <Link to="/register">
                                <Button className="px-8 py-3 rounded-full text-lg w-full sm:w-auto">
                                    Get Started - It's Free
                                </Button>
                            </Link>
                            <Link to="/features">
                                <button className="px-8 py-3 rounded-full text-lg border border-white/20 hover:bg-white/5 transition-colors w-full sm:w-auto">
                                    Explore Features
                                </button>
                            </Link>
                        </div>
                        <p className="mt-8 text-xs text-gray-500 uppercase tracking-widest">
                            Verified by 8 million members
                        </p>
                    </motion.div>
                </div>
            </section>
        </div>
    );
};

export default LandingPage;
