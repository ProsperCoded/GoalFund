import { createContext, useContext, useState, useEffect } from "react"
import type { ReactNode } from "react"
import { authApi, type User, type LoginRequest, type RegisterRequest } from "@/lib/api"
import { useToast } from "@/hooks/use-toast"

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (credentials: LoginRequest) => Promise<void>
  register: (data: RegisterRequest) => Promise<void>
  logout: () => Promise<void>
  updateUser: (user: User) => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

interface AuthProviderProps {
  children: ReactNode
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const { toast } = useToast()

  // Load user from localStorage on mount
  useEffect(() => {
    const loadUser = async () => {
      try {
        const storedUser = localStorage.getItem("user")
        const accessToken = localStorage.getItem("access_token")

        if (storedUser && accessToken) {
          setUser(JSON.parse(storedUser))
          
          // Optionally verify token is still valid by fetching profile
          try {
            const { user: freshUser } = await authApi.getProfile()
            setUser(freshUser)
            localStorage.setItem("user", JSON.stringify(freshUser))
          } catch (error) {
            // Token might be invalid, clear storage
            localStorage.removeItem("user")
            localStorage.removeItem("access_token")
            localStorage.removeItem("refresh_token")
            setUser(null)
          }
        }
      } catch (error) {
        console.error("Failed to load user:", error)
      } finally {
        setIsLoading(false)
      }
    }

    loadUser()
  }, [])

  const login = async (credentials: LoginRequest) => {
    try {
      const response = await authApi.login(credentials)
      
      // Store tokens and user
      localStorage.setItem("access_token", response.access_token)
      localStorage.setItem("refresh_token", response.refresh_token)
      localStorage.setItem("user", JSON.stringify(response.user))
      
      setUser(response.user)
      
      toast({
        title: "Welcome back!",
        description: `Logged in as ${response.user.email}`,
      })
    } catch (error: any) {
      const errorMessage = error.response?.data?.error || "Login failed. Please try again."
      toast({
        title: "Login Failed",
        description: errorMessage,
        variant: "destructive",
      })
      throw error
    }
  }

  const register = async (data: RegisterRequest) => {
    try {
      const response = await authApi.register(data)
      
      // Store tokens and user (if tokens were returned)
      if (response.access_token) {
        localStorage.setItem("access_token", response.access_token)
        localStorage.setItem("refresh_token", response.refresh_token)
      }
      localStorage.setItem("user", JSON.stringify(response.user))
      
      setUser(response.user)
      
      toast({
        title: "Account Created!",
        description: "Welcome to GoFund",
      })
    } catch (error: any) {
      const errorMessage = error.response?.data?.error || "Registration failed. Please try again."
      toast({
        title: "Registration Failed",
        description: errorMessage,
        variant: "destructive",
      })
      throw error
    }
  }

  const logout = async () => {
    try {
      const refreshToken = localStorage.getItem("refresh_token")
      if (refreshToken) {
        await authApi.logout(refreshToken)
      }
    } catch (error) {
      console.error("Logout error:", error)
    } finally {
      // Clear storage and state regardless of API call result
      localStorage.removeItem("access_token")
      localStorage.removeItem("refresh_token")
      localStorage.removeItem("user")
      setUser(null)
      
      toast({
        title: "Logged Out",
        description: "You have been successfully logged out",
      })
    }
  }

  const updateUser = (updatedUser: User) => {
    setUser(updatedUser)
    localStorage.setItem("user", JSON.stringify(updatedUser))
  }

  const value: AuthContextType = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    register,
    logout,
    updateUser,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider")
  }
  return context
}
