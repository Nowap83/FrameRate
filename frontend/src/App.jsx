import { Route, Routes } from "react-router-dom"
import Home from "./pages/Home"
import AuthPage from "./pages/AuthPage"
import VerifyEmail from "./pages/VerifyEmail"

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/login" element={<AuthPage />} />
      <Route path="/register" element={<AuthPage />} />
      <Route path="/verify-email" element={<VerifyEmail />} />
    </Routes>
  );
};

export default App
