import { Route, Routes } from "react-router-dom"
import Home from "./pages/Home"
import AuthPage from "./pages/AuthPage"

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/login" element={<AuthPage />} />
      <Route path="/register" element={<AuthPage />} />
    </Routes>
  );
};

export default App
