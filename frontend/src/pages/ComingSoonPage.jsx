import { motion } from "framer-motion";
import { Clock, ArrowLeft } from "lucide-react";
import Button from "../components/Button";
import { Link } from "react-router-dom";
import useDocumentTitle from "../hooks/useDocumentTitle";

const ComingSoonPage = ({ title = "Under Construction" }) => {
    useDocumentTitle(title);

    return (
        <div className="min-h-[80vh] flex flex-col items-center justify-center text-center px-4">
            <motion.div
                initial={{ opacity: 0, scale: 0.9 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ duration: 0.5 }}
                className="bg-white/5 border border-white/10 rounded-3xl p-10 md:p-16 max-w-lg w-full shadow-2xl backdrop-blur-md"
            >
                <motion.div
                    initial={{ y: -20, opacity: 0 }}
                    animate={{ y: 0, opacity: 1 }}
                    transition={{ delay: 0.2 }}
                    className="flex justify-center mb-6"
                >
                    <div className="w-20 h-20 bg-mint/20 rounded-full flex items-center justify-center text-mint">
                        <Clock size={40} />
                    </div>
                </motion.div>

                <motion.h1
                    initial={{ y: 20, opacity: 0 }}
                    animate={{ y: 0, opacity: 1 }}
                    transition={{ delay: 0.3 }}
                    className="text-3xl md:text-4xl font-bold font-display text-white mb-4"
                >
                    {title}
                </motion.h1>

                <motion.p
                    initial={{ y: 20, opacity: 0 }}
                    animate={{ y: 0, opacity: 1 }}
                    transition={{ delay: 0.4 }}
                    className="text-gray-400 mb-8"
                >
                    We're working hard to bring this feature to life. Check back soon for updates as our community grows!
                </motion.p>

                <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ delay: 0.5 }}
                >
                    <Link to="/">
                        <Button className="px-8 py-3 rounded-full font-semibold flex items-center gap-2 mx-auto">
                            <ArrowLeft size={18} /> Back to Home
                        </Button>
                    </Link>
                </motion.div>
            </motion.div>
        </div>
    );
};

export default ComingSoonPage;
