import React, { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import Button from "../components/Button";
import Input from "../components/Input";
import { AUTH_MOVIES } from "../data/authMovies";
import { loginSchema, registerSchema } from "../validators/auth";

const AuthPage = () => {
    const location = useLocation();
    const navigate = useNavigate();

    const isLogin = location.pathname === "/login";
    const [currentMovieIndex, setCurrentMovieIndex] = useState(0);

    // init form avec react-hook-form
    const {
        register,
        handleSubmit,
        formState: { errors },
        reset
    } = useForm({
        resolver: zodResolver(isLogin ? loginSchema : registerSchema),
        mode: "onBlur"
    });

    // reset form quand on switch
    useEffect(() => {
        reset();
    }, [isLogin, reset]);

    const onSubmit = (data) => {
        console.log("Form Submitted:", data);
        // call api
    };

    const pickRandomMovie = () => {
        let nextIndex;
        do {
            nextIndex = Math.floor(Math.random() * AUTH_MOVIES.length);
        } while (nextIndex === currentMovieIndex && AUTH_MOVIES.length > 1);
        setCurrentMovieIndex(nextIndex);
    };

    const toggleAuthMode = () => {
        navigate(isLogin ? "/register" : "/login");
        pickRandomMovie();
    };

    const movie = AUTH_MOVIES[currentMovieIndex];

    return (
        <div className="min-h-screen bg-[#0A0F0D] flex items-center justify-center p-4 overflow-hidden">
            <div className="w-full max-w-6xl h-[750px] bg-header-bg rounded-3xl shadow-2xl overflow-hidden flex relative">

                {/* image */}
                <motion.div
                    className="absolute top-0 bottom-0 w-[60%] overflow-hidden z-10 hidden md:block"
                    animate={{ x: isLogin ? 0 : "66.666%" }}
                    transition={{ type: "spring", stiffness: 50, damping: 15 }}
                >
                    <AnimatePresence mode="wait">
                        <motion.div
                            key={movie.id}
                            initial={{ opacity: 0, scale: 1.1 }}
                            animate={{ opacity: 1, scale: 1 }}
                            exit={{ opacity: 0, scale: 0.95 }}
                            transition={{ duration: 0.6 }}
                            className="relative w-full h-full"
                        >
                            <img
                                src={movie.image}
                                className="w-full h-full object-cover brightness-50"
                                alt={movie.title}
                            />
                            <div className="absolute inset-0 bg-gradient-to-t from-header-bg via-transparent to-transparent opacity-60" />
                            <div className="absolute inset-0 flex flex-col items-center justify-center p-12 text-center text-white">
                                <motion.h2
                                    initial={{ y: 20, opacity: 0 }}
                                    animate={{ y: 0, opacity: 1 }}
                                    key={`title-${movie.id}`}
                                    className="text-6xl font-bold mb-6 italic"
                                >
                                    {isLogin ? "Welcome Back!" : "Join us!"}
                                </motion.h2>
                                <motion.p
                                    initial={{ y: 20, opacity: 0 }}
                                    animate={{ y: 0, opacity: 1 }}
                                    transition={{ delay: 0.1 }}
                                    key={`quote-${movie.id}`}
                                    className="text-xl text-gray-200 leading-relaxed max-w-md"
                                >
                                    {movie.quote}
                                </motion.p>
                                <div className="absolute bottom-10 text-xs text-gray-400">
                                    {movie.title} - {movie.year}
                                </div>
                            </div>
                        </motion.div>
                    </AnimatePresence>
                </motion.div>

                {/* form */}
                <motion.div
                    className="w-full md:w-[40%] p-8 md:p-14 flex flex-col justify-center bg-[#162520] z-0"
                    animate={{ x: isLogin ? "150%" : 0 }}
                    transition={{ type: "spring", stiffness: 50, damping: 15 }}
                >
                    <div className="flex items-center gap-2 text-white font-bold text-lg mb-12 justify-center md:justify-start">
                        <div className="text-mint">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M18 4l2 4h-3l-2-4h-2l2 4h-3l-2-4H8l2 4H7L5 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V4h-4z" />
                            </svg>
                        </div>
                        FrameRate
                    </div>

                    <AnimatePresence mode="wait">
                        {isLogin ? (
                            <motion.div
                                key="login"
                                initial={{ opacity: 0, x: 20 }}
                                animate={{ opacity: 1, x: 0 }}
                                exit={{ opacity: 0, x: -20 }}
                                className="w-full"
                            >
                                <h3 className="text-2xl font-bold text-white mb-2">Welcome Back</h3>
                                <p className="text-gray-400 text-sm mb-8">Sign in to track your watch list and rate your favorites.</p>

                                <form className="space-y-5" onSubmit={handleSubmit(onSubmit)}>
                                    <Input
                                        label="Email or Username"
                                        placeholder="Enter your email..."
                                        {...register("email")}
                                        error={errors.email?.message}
                                    />
                                    <Input
                                        label="Password"
                                        type="password"
                                        placeholder="••••••••"
                                        {...register("password")}
                                        error={errors.password?.message}
                                    />

                                    <div className="flex items-center justify-between text-xs">
                                        <label className="flex items-center gap-2 text-gray-400 cursor-pointer">
                                            <input type="checkbox" className="rounded border-white/10 bg-white/5 text-mint focus:ring-mint" />
                                            Remember me
                                        </label>
                                        <a href="#" className="text-gray-400 hover:text-white transition-colors">Forgot Password?</a>
                                    </div>

                                    <Button type="submit" className="w-full py-3 mt-4">Login</Button>
                                </form>

                                <p className="text-center text-sm text-gray-400 mt-8">
                                    Don't have an account?{" "}
                                    <button onClick={toggleAuthMode} className="text-white font-semibold hover:text-mint transition-colors">
                                        Sign Up
                                    </button>
                                </p>
                            </motion.div>
                        ) : (
                            <motion.div
                                key="register"
                                initial={{ opacity: 0, x: -20 }}
                                animate={{ opacity: 1, x: 0 }}
                                exit={{ opacity: 0, x: 20 }}
                                className="w-full"
                            >
                                <h3 className="text-2xl font-bold text-white mb-2">Create Account</h3>
                                <p className="text-gray-400 text-sm mb-8">Join the community of cinema lovers today.</p>

                                <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
                                    <Input
                                        label="Email Address"
                                        type="email"
                                        placeholder="name@example.com"
                                        {...register("email")}
                                        error={errors.email?.message}
                                    />
                                    <Input
                                        label="Username"
                                        placeholder="Choose a username"
                                        {...register("username")}
                                        error={errors.username?.message}
                                    />
                                    <Input
                                        label="Password"
                                        type="password"
                                        placeholder="••••••••"
                                        {...register("password")}
                                        error={errors.password?.message}
                                    />
                                    <Input
                                        label="Confirm Password"
                                        type="password"
                                        placeholder="••••••••"
                                        {...register("confirmPassword")}
                                        error={errors.confirmPassword?.message}
                                    />

                                    <Button type="submit" className="w-full py-3 mt-4">Register</Button>
                                </form>

                                <p className="text-center text-sm text-gray-400 mt-8">
                                    Already have an account?{" "}
                                    <button onClick={toggleAuthMode} className="text-white font-semibold hover:text-mint transition-colors">
                                        Log In
                                    </button>
                                </p>
                            </motion.div>
                        )}
                    </AnimatePresence>
                </motion.div>

            </div>
        </div>
    );
};

export default AuthPage;
