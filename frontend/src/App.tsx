import { BrowserRouter as Router, Routes, Route } from "react-router-dom"
import { Layout } from "@/components/layout"
import { HomePage, LoginPage, RegisterPage, ForgotPasswordPage } from "@/pages"

function App() {
  return (
    <Router>
      <Routes>
        {/* Auth Routes - No Layout */}
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/forgot-password" element={<ForgotPasswordPage />} />

        {/* Main Routes - With Layout */}
        <Route element={<Layout />}>
          <Route path="/" element={<HomePage />} />
          {/* Add more routes here */}
        </Route>
      </Routes>
    </Router>
  )
}

export default App
