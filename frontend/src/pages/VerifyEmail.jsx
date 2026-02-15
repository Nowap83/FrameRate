import React, { useEffect, useState, useRef } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { motion, AnimatePresence } from "framer-motion";
import { Mail, CheckCircle, XCircle, Loader2 } from "lucide-react";
import { authService } from "../api/auth";
import Button from "../components/Button";

const VerifyEmail = () => {
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();
    const token = searchParams.get("token");
    const hasRun = useRef(false);

    const [status, setStatus] = useState("verifying"); // verifying, success, error
    const [message, setMessage] = useState("");
    const [countdown, setCountdown] = useState(5);

    useEffect(() => {
        const verify = async () => {
            if (!token) {
                setStatus("error");
                setMessage("Missing verification token.");
                return;
            }

            // Empêche le double appel en développement (StrictMode)
            if (hasRun.current) return;
            hasRun.current = true;

            try {
                const response = await authService.verifyEmail(token);
                setStatus("success");
                setMessage(response.message || "Your email has been successfully verified!");

                // Si le backend renvoie un token, on connecte l'utilisateur directement
                if (response.token) {
                    localStorage.setItem("token", response.token);
                }
            } catch (error) {
                setStatus("error");
                setMessage(
                    error.response?.data?.error ||
                    error.response?.data?.message ||
                    "Verification failed. The link may be expired or invalid."
                );
            }
        };

        verify();
    }, [token]);

    // Timer de redirection automatique en cas de succès
    useEffect(() => {
        if (status === "success" && countdown > 0) {
            const timer = setTimeout(() => setCountdown(countdown - 1), 1000);
            return () => clearTimeout(timer);
        } else if (status === "success" && countdown === 0) {
            navigate("/"); // Redirection vers la Home au lieu du Login car il est déjà co !
        }
    }, [status, countdown, navigate]);

    return (
        <div className="min-h-screen bg-[#0A0F0D] flex items-center justify-center p-4">
            <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                className="w-full max-w-md bg-header-bg rounded-3xl p-10 shadow-2xl text-center"
            >
                <div className="flex justify-center mb-8">
                    <div className="text-white font-bold text-2xl flex items-center gap-2">
                        <div className="text-mint">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M18 4l2 4h-3l-2-4h-2l2 4h-3l-2-4H8l2 4H7L5 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V4h-4z" />
                            </svg>
                        </div>
                        FrameRate
                    </div>
                </div>

                <AnimatePresence mode="wait">
                    {status === "verifying" && (
                        <motion.div
                            key="verifying"
                            initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}
                            className="space-y-4"
                        >
                            <Loader2 className="h-16 w-16 text-mint animate-spin mx-auto mb-6" />
                            <h2 className="text-2xl font-bold text-white uppercase italic">Verifying...</h2>
                            <p className="text-gray-400 text-sm italic">Please wait while we confirm your account.</p>
                        </motion.div>
                    )}

                    {status === "success" && (
                        <motion.div
                            key="success"
                            initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}
                            className="space-y-4"
                        >
                            <CheckCircle className="h-16 w-16 text-mint mx-auto mb-6" />
                            <h2 className="text-2xl font-bold text-white uppercase italic">Verified!</h2>
                            <p className="text-gray-400">{message}</p>
                            <div className="pt-6">
                                <p className="text-xs text-gray-500 mb-4 italic uppercase">Redirecting to Home in {countdown}s...</p>
                                <Button onClick={() => navigate("/")} className="w-full">
                                    Go to Home now
                                </Button>
                            </div>
                        </motion.div>
                    )}

                    {status === "error" && (
                        <motion.div
                            key="error"
                            initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}
                            className="space-y-4"
                        >
                            <XCircle className="h-16 w-16 text-red-500 mx-auto mb-6" />
                            <h2 className="text-2xl font-bold text-white uppercase italic">Failed</h2>
                            <p className="text-red-400 text-sm px-4 py-2 bg-red-500/10 border border-red-500/20 rounded-lg inline-block italic">
                                {message}
                            </p>
                            <div className="pt-6">
                                <Button onClick={() => navigate("/register")} className="w-full">
                                    Return to Sign Up
                                </Button>
                                <button
                                    onClick={() => navigate("/login")}
                                    className="text-sm text-gray-400 hover:text-white mt-6 transition-colors uppercase tracking-widest font-bold"
                                >
                                    Back to Login
                                </button>
                            </div>
                        </motion.div>
                    )}
                </AnimatePresence>
            </motion.div>
        </div>
    );
};

export default VerifyEmail;
