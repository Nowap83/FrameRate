const Button = ({ children, className = "", ...rest }) => {
    return (
        <button
            className={`px-6 py-2.5 bg-mint text-gray-900 font-semibold rounded-lg shadow-sm hover:brightness-95 transition-all active:scale-95 ${className}`}
            {...rest}
        >
            {children}
        </button>
    );
};

export default Button;