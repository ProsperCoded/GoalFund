import api from "./config"

// Types
export interface InitializePaymentRequest {
  user_id: string
  goal_id: string
  amount: number // Amount in kobo (100 kobo = 1 NGN)
  currency: string
  email: string
  callback_url?: string
  metadata?: Record<string, unknown>
}

export interface InitializePaymentResponse {
  payment_id: string
  authorization_url: string
  access_code: string
  reference: string
}

export interface PaymentStatus {
  payment_id: string
  reference: string
  status: "INITIATED" | "PENDING" | "VERIFIED" | "FAILED"
  amount: number
  currency: string
  paid_at?: string
  channel?: string
  created_at: string
  updated_at: string
  metadata?: Record<string, unknown>
}

export interface Bank {
  id: number
  name: string
  code: string
}

export interface ResolvedAccount {
  account_number: string
  account_name: string
  bank_code: string
  bank_name: string
}

// API Functions
export const paymentsApi = {
  /**
   * Initialize a payment for a goal contribution
   * Returns Paystack checkout URL for user to complete payment
   */
  initialize: async (data: InitializePaymentRequest) => {
    const response = await api.post<{ status: string; data: InitializePaymentResponse }>(
      "/payments/initialize",
      data
    )
    return response.data.data
  },

  /**
   * Verify a payment after user completes Paystack checkout
   * Call this after redirect from Paystack or to check payment status
   */
  verify: async (reference: string) => {
    const response = await api.get<{ status: string; data: PaymentStatus }>(
      `/payments/verify/${reference}`
    )
    return response.data.data
  },

  /**
   * Get payment status by payment ID
   */
  getStatus: async (paymentId: string) => {
    const response = await api.get<{ status: string; data: PaymentStatus }>(
      `/payments/${paymentId}/status`
    )
    return response.data.data
  },

  /**
   * Get list of supported banks for account resolution
   */
  getBanks: async (country: string = "nigeria") => {
    const response = await api.get<{ status: string; data: Bank[] }>(
      `/payments/banks`,
      { params: { country } }
    )
    return response.data.data
  },

  /**
   * Resolve bank account to get account holder name
   * Used to verify bank details before saving
   */
  resolveAccount: async (accountNumber: string, bankCode: string) => {
    const response = await api.get<{ status: string; data: ResolvedAccount }>(
      `/payments/resolve-account`,
      { params: { account_number: accountNumber, bank_code: bankCode } }
    )
    return response.data.data
  },
}

// Helper function to convert NGN to kobo
export function nairaToKobo(naira: number): number {
  return Math.round(naira * 100)
}

// Helper function to convert kobo to NGN
export function koboToNaira(kobo: number): number {
  return kobo / 100
}

// Helper to format payment amount for display
export function formatPaymentAmount(amountInKobo: number): string {
  const naira = koboToNaira(amountInKobo)
  return new Intl.NumberFormat("en-NG", {
    style: "currency",
    currency: "NGN",
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(naira)
}
