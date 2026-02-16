import React, { useState } from "react";
import { Eye, EyeOff } from "lucide-react";

const Input = React.forwardRef(({ label, error, type = "text", className = "", ...props }, ref) => {
    const [showPassword, setShowPassword] = useState(false);
    const isPassword = type === "password";

    const togglePasswordVisibility = () => {
        setShowPassword(!showPassword);
    };

    const inputType = isPassword ? (showPassword ? "text" : "password") : type;

    return (
        <div className={`flex flex-col gap-1.5 ${className}`}>
            {label && (
                <label className="text-xs font-medium text-gray-400 ml-1 uppercase tracking-wider">
                    {label}
                </label>
            )}
            <div className="relative group">
                <input
                    {...props}
                    ref={ref}
                    type={inputType}
                    className={`w-full bg-white/5 border ${error ? "border-red-500/50" : "border-white/10"
                        } text-white text-sm rounded-lg p-3 outline-none 
          group-hover:border-white/20 focus:border-mint/50 focus:ring-1 focus:ring-mint/50 
          transition-all placeholder:text-gray-500`}
                />
                {isPassword && (
                    <button
                        type="button"
                        onClick={togglePasswordVisibility}
                        className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-mint transition-colors"
                    >
                        {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                    </button>
                )}
            </div>
            {error && <span className="text-xs text-red-500 ml-1">{error}</span>}
        </div>
    );
});

export default Input;
