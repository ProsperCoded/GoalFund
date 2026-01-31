// Auth Types
export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  username?: string
  password?: string
  first_name?: string
  last_name?: string
  phone?: string
  settlement_bank_name?: string
  settlement_account_number?: string
  settlement_account_name?: string
}

export interface AuthResponse {
  user: User
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
}

export interface User {
  id: string
  email: string
  username: string
  first_name: string
  last_name: string
  phone: string
  email_verified: boolean
  phone_verified: boolean
  kyc_verified: boolean
  kyc_verified_at?: string
  role: "user" | "admin"
  created_at: string
}

export interface RefreshTokenRequest {
  refresh_token: string
}

export interface ForgotPasswordRequest {
  email: string
}

export interface ResetPasswordRequest {
  token: string
  new_password: string
}

export interface UpdateProfileRequest {
  first_name?: string
  last_name?: string
  phone?: string
}

export interface UpdateSettlementAccountRequest {
  bank_name: string
  account_number: string
  account_name: string
}

export interface SubmitNINRequest {
  nin: string
}

export interface KYCStatusResponse {
  kyc_verified: boolean
  kyc_verified_at?: string
  nin?: string
}

export interface ApiError {
  error: string
  details?: string
}
