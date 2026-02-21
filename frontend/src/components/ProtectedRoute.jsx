import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

const ProtectedRoute = ({ children, requireAdmin = false }) => {
    const { user, loading } = useAuth();
    const location = useLocation();

    if (loading) {
        return (
            <div className="min-h-screen bg-[var(--color-body-bg)] flex items-center justify-center text-white">
                <div className="loader"></div>
            </div>
        );
    }

    if (!user) {
        // Rediriger vers login, mais garder en mémoire d'où l'utilisateur vient
        return <Navigate to="/login" state={{ from: location }} replace />;
    }

    if (requireAdmin && !user.is_admin) {
        // Rediriger vers l'accueil si ce n'est pas un admin
        return <Navigate to="/" replace />;
    }

    // Autoriser l'accès
    return children;
};

export default ProtectedRoute;
