import React, { useState } from 'react';
import { Star } from 'lucide-react';

const RatingStars = ({ rating, onChange, onHoverChange, readOnly = false, size = 24 }) => {
    const [hoverValue, setHoverValue] = useState(0);

    const handleMouseMove = (e, index) => {
        if (readOnly) return;

        const rect = e.currentTarget.getBoundingClientRect();
        const width = rect.width;
        const clickX = e.clientX - rect.left;

        let newValue = index + 1;
        if (clickX < width / 2) {
            newValue -= 0.5;
        }

        setHoverValue(newValue);
        if (onHoverChange) {
            onHoverChange(newValue);
        }
    };

    const handleClick = (e, index) => {
        if (readOnly || !onChange) return;

        const rect = e.currentTarget.getBoundingClientRect();
        const width = rect.width;
        const clickX = e.clientX - rect.left;

        let newValue = index + 1;
        if (clickX < width / 2) {
            newValue -= 0.5;
        }

        if (rating === newValue) {
            // allow un-rating by clicking the same exact value
            onChange(0);
            setHoverValue(0);
        } else {
            onChange(newValue);
        }
    };

    const handleMouseLeave = () => {
        if (readOnly) return;
        setHoverValue(0);
        if (onHoverChange) {
            onHoverChange(0);
        }
    };

    const displayValue = hoverValue > 0 ? hoverValue : rating;

    return (
        <div
            className={`flex items-center gap-1 ${readOnly ? '' : 'cursor-pointer'}`}
            onMouseLeave={handleMouseLeave}
        >
            {[...Array(5)].map((_, index) => {
                const filledValue = displayValue - index;

                // Determine star state
                let fill = 'none';
                let strokeClass = 'text-gray-400';

                if (filledValue >= 1) {
                    fill = 'currentColor';
                    strokeClass = 'text-[var(--color-primary)]';
                } else if (filledValue === 0.5) {
                    // Half star handling
                    return (
                        <div
                            key={index}
                            className="relative group w-auto h-auto"
                            onMouseMove={(e) => handleMouseMove(e, index)}
                            onClick={(e) => handleClick(e, index)}
                        >
                            <Star size={size} className="text-gray-400" />
                            <div className="absolute top-0 left-0 overflow-hidden w-[50%]">
                                <Star size={size} className="text-[var(--color-primary)] fill-current" />
                            </div>
                        </div>
                    );
                }

                return (
                    <div
                        key={index}
                        onMouseMove={(e) => handleMouseMove(e, index)}
                        onClick={(e) => handleClick(e, index)}
                    >
                        <Star
                            size={size}
                            fill={fill}
                            className={`transition-colors ${strokeClass}`}
                        />
                    </div>
                );
            })}
        </div>
    );
};

export default RatingStars;
