import apiClient from "./config"
import type {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  RefreshTokenRequest,
  ForgotPasswordRequest,
  ResetPasswordRequest,
  User,
  UpdateProfileRequest,
  UpdateSettlementAccountRequest,
  SubmitNINRequest,
  KYCStatusResponse,
} from "./types"

export const authApi = {
  // Login user
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>("/auth/login", data)
    return response.data
  },

  // Register new user
  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>("/auth/register", data)
    return response.data
  },

  // Refresh access token
  refresh: async (data: RefreshTokenRequest): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>("/auth/refresh", data)
    return response.data
  },

  // Logout user
  logout: async (refreshToken: string): Promise<void> => {
    await apiClient.post("/auth/logout", { refresh_token: refreshToken })
  },

  // Forgot password
  forgotPassword: async (data: ForgotPasswordRequest): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>("/auth/forgot-password", data)
    return response.data
  },

  // Reset password
  resetPassword: async (data: ResetPasswordRequest): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>("/auth/reset-password", data)
    return response.data
  },

  // Get user profile
  getProfile: async (): Promise<{ user: User }> => {
    const response = await apiClient.get<{ user: User }>("/users/profile")
    return response.data
  },

  // Update user profile
  updateProfile: async (data: UpdateProfileRequest): Promise<{ user: User; message: string }> => {
    const response = await apiClient.put<{ user: User; message: string }>("/users/profile", data)
    return response.data
  },

  // Update settlement account (optional)
  updateSettlementAccount: async (data: UpdateSettlementAccountRequest): Promise<{ message: string }> => {
    const response = await apiClient.put<{ message: string }>("/users/settlement-account", data)
    return response.data
  },

  // Submit NIN for KYC (optional)
  submitNIN: async (data: SubmitNINRequest): Promise<KYCStatusResponse> => {
    const response = await apiClient.post<KYCStatusResponse>("/users/kyc/submit-nin", data)
    return response.data
  },

  // Get KYC status
  getKYCStatus: async (): Promise<KYCStatusResponse> => {
    const response = await apiClient.get<KYCStatusResponse>("/users/kyc/status")
    return response.data
  },
}
