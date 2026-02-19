import { useEffect } from "react";
import { Route, Routes, Navigate, useLocation } from "react-router-dom"
import Dashboard from "./pages/Dashboard"
import LandingPage from "./pages/LandingPage"
import AuthPage from "./pages/AuthPage"
import VerifyEmail from "./pages/VerifyEmail"
import AppLayout from "./layouts/AppLayout"
import MovieDetails from "./pages/MovieDetails"
import Profile from "./pages/Profile"
import Settings from "./pages/Settings"
import { useAuth } from "./context/AuthContext"

function App() {
  const { user, loading } = useAuth();
  const { pathname } = useLocation();

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [pathname]);

  if (loading) {
    return <div className="min-h-screen bg-[var(--color-body-bg)] flex items-center justify-center text-white">Loading...</div>;
  }

  return (
    <Routes>
      <Route path="/" element={
        <AppLayout>
          {user ? <Dashboard /> : <LandingPage />}
        </AppLayout>
      } />
      <Route path="/login" element={user ? <Navigate to="/" /> : <AuthPage />} />
      <Route path="/register" element={user ? <Navigate to="/" /> : <AuthPage />} />
      <Route path="/movie/:id" element={
        <AppLayout>
          <MovieDetails />
        </AppLayout>
      } />
      <Route path="/profile" element={
        <AppLayout>
          <Profile />
        </AppLayout>
      } />
      <Route path="/verify-email" element={
        <AppLayout>
          <VerifyEmail />
        </AppLayout>
      } />
      <Route path="/settings" element={
        <AppLayout>
          <Settings />
        </AppLayout>
      } />
    </Routes>
  );
};

export default App
