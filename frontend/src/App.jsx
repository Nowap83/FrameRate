import { useEffect } from "react";
import { Route, Routes, Navigate, useLocation } from "react-router-dom"
import HomePage from "./pages/HomePage"
import LandingPage from "./pages/LandingPage"
import AuthPage from "./pages/AuthPage"
import VerifyEmail from "./pages/VerifyEmail"
import AppLayout from "./layouts/AppLayout"
import MovieDetails from "./pages/MovieDetails"
import Profile from "./pages/Profile"
import Settings from "./pages/Settings"
import SearchPage from "./pages/SearchPage"
import PersonDetails from "./pages/PersonDetails"
import AdminDashboard from "./pages/AdminDashboard"
import { useAuth } from "./context/AuthContext"
import ProtectedRoute from "./components/ProtectedRoute"
import ErrorBoundary from "./components/ErrorBoundary"

function App() {
  const { user, loading } = useAuth();
  const { pathname } = useLocation();

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [pathname]);

  if (loading) {
    return <div className="min-h-screen bg-[var(--color-body-bg)] flex items-center justify-center text-white"><div className="loader"></div></div>;
  }

  const withLayout = (Component) => (
    <ErrorBoundary>
      <AppLayout>
        {Component}
      </AppLayout>
    </ErrorBoundary>
  );

  return (
    <Routes>
      <Route path="/" element={withLayout(user ? <HomePage /> : <LandingPage />)} />
      <Route path="/login" element={user ? <Navigate to="/" /> : <AuthPage />} />
      <Route path="/register" element={user ? <Navigate to="/" /> : <AuthPage />} />
      <Route path="/movie/:id" element={withLayout(<MovieDetails />)} />
      <Route path="/search" element={withLayout(<SearchPage />)} />
      <Route path="/person/:id" element={withLayout(<PersonDetails />)} />
      <Route path="/verify-email" element={withLayout(<VerifyEmail />)} />

      {/* Routes Protégées (Utilisateurs uniquement) */}
      <Route path="/profile" element={
        <ProtectedRoute>
          {withLayout(<Profile />)}
        </ProtectedRoute>
      } />
      <Route path="/settings" element={
        <ProtectedRoute>
          {withLayout(<Settings />)}
        </ProtectedRoute>
      } />

      {/* Routes Admin (Admin uniquement) */}
      <Route path="/admin" element={
        <ProtectedRoute requireAdmin={true}>
          {withLayout(<AdminDashboard />)}
        </ProtectedRoute>
      } />
    </Routes>
  );
};

export default App
