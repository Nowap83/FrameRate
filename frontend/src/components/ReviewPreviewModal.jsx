import React, { useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { X } from 'lucide-react';
import RatingStars from './RatingStars';

const ReviewPreviewModal = ({ isOpen, onClose, review }) => {
    useEffect(() => {
        if (isOpen) {
            document.body.style.overflow = 'hidden';
        } else {
            document.body.style.overflow = 'unset';
        }
        return () => {
            document.body.style.overflow = 'unset';
        };
    }, [isOpen]);

    if (!isOpen || !review) return null;

    const {
        title,
        release_year,
        rating,
        content,
        watched_date,
    } = review;

    const watchedDateFormatted = watched_date
        ? new Date(watched_date).toLocaleDateString("en-GB", { day: '2-digit', month: 'short', year: 'numeric' })
        : null;

    return (
        <AnimatePresence>
            {isOpen && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        onClick={onClose}
                        className="absolute inset-0 bg-black/80 backdrop-blur-sm"
                    />
                    
                    <motion.div
                        initial={{ opacity: 0, scale: 0.95, y: 20 }}
                        animate={{ opacity: 1, scale: 1, y: 0 }}
                        exit={{ opacity: 0, scale: 0.95, y: 20 }}
                        className="relative w-full max-w-lg bg-[#1A2C24] border border-white/10 rounded-2xl shadow-2xl overflow-hidden z-10"
                    >
                        {/* Header */}
                        <div className="flex items-center justify-between p-4 md:p-6 border-b border-white/5 bg-black/20">
                            <div>
                                <h2 className="text-xl md:text-2xl font-display font-bold text-white flex items-baseline gap-2">
                                    {title}
                                    {release_year > 0 && <span className="text-sm font-normal text-gray-400">{release_year}</span>}
                                </h2>
                                <div className="flex items-center gap-3 mt-2 text-sm">
                                    {rating > 0 && <RatingStars rating={rating} size={14} readOnly />}
                                    {watchedDateFormatted && <span className="text-gray-400">Watched {watchedDateFormatted}</span>}
                                </div>
                            </div>
                            <button
                                onClick={onClose}
                                className="p-2 text-gray-500 hover:text-white hover:bg-white/10 rounded-full transition-colors self-start"
                            >
                                <X size={20} />
                            </button>
                        </div>

                        {/* Content */}
                        <div className="p-4 md:p-6 max-h-[60vh] overflow-y-auto custom-scrollbar">
                            <p className="text-gray-300 text-sm md:text-base leading-relaxed whitespace-pre-wrap">
                                {content || <span className="italic text-gray-500">No review text provided.</span>}
                            </p>
                        </div>
                    </motion.div>
                </div>
            )}
        </AnimatePresence>
    );
};

export default ReviewPreviewModal;
