import { BrowserRouter as Router, Routes, Route } from "react-router-dom"
import { Layout } from "@/components/layout"
import { HomePage, LoginPage, RegisterPage, ForgotPasswordPage, OnboardingPage } from "@/pages"
import { AuthProvider } from "@/contexts"
import { Toaster } from "@/components/ui/toaster"

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          {/* Auth Routes - No Layout */}
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/forgot-password" element={<ForgotPasswordPage />} />
          <Route path="/onboarding" element={<OnboardingPage />} />

          {/* Main Routes - With Layout */}
          <Route element={<Layout />}>
            <Route path="/" element={<HomePage />} />
            {/* Add more routes here */}
          </Route>
        </Routes>
        <Toaster />
      </Router>
    </AuthProvider>
  )
}

export default App
