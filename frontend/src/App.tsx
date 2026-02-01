import { BrowserRouter as Router, Routes, Route } from "react-router-dom"
import { Layout } from "@/components/layout"
import { HomePage, LoginPage, RegisterPage, ForgotPasswordPage, OnboardingPage } from "@/pages"
import {
  DashboardPage,
  MyGoalsPage,
  MyContributionsPage,
  ExplorePage,
  CreateGoalPage,
  GoalDetailPage,
  SettingsPage,
} from "@/pages/dashboard"
import { DashboardLayout } from "@/components/dashboard"
import { ProtectedRoute } from "@/components/auth/ProtectedRoute"
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

          {/* Dashboard Routes - Protected */}
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <DashboardLayout />
              </ProtectedRoute>
            }
          >
            <Route index element={<DashboardPage />} />
            <Route path="goals" element={<MyGoalsPage />} />
            <Route path="goals/create" element={<CreateGoalPage />} />
            <Route path="goals/:goalId" element={<GoalDetailPage />} />
            <Route path="contributions" element={<MyContributionsPage />} />
            <Route path="explore" element={<ExplorePage />} />
            <Route path="settings" element={<SettingsPage />} />
          </Route>

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
