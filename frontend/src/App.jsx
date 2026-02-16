import { Route, Routes, Navigate } from "react-router-dom"
import Dashboard from "./pages/Dashboard"
import LandingPage from "./pages/LandingPage"
import AuthPage from "./pages/AuthPage"
import VerifyEmail from "./pages/VerifyEmail"
import AppLayout from "./layouts/AppLayout"
import { useAuth } from "./context/AuthContext"

function App() {
  const { user, loading } = useAuth();

  if (loading) {
    return <div className="min-h-screen bg-[#12201B] flex items-center justify-center text-white">Loading...</div>;
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
      <Route path="/verify-email" element={
        <AppLayout>
          <VerifyEmail />
        </AppLayout>
      } />
    </Routes>
  );
};

export default App
