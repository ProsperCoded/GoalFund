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
  kyc_status?: "pending" | "submitted" | "verified" | "rejected"
  settlement_account_status?: "pending" | "verified" | "rejected"
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

// Goal Types
export type GoalStatus = "open" | "closed" | "cancelled" | "OPEN" | "FUNDED" | "WITHDRAWN" | "CLOSED" | "CANCELLED"
export type MilestoneRecurrence = "one_time" | "weekly" | "monthly" | "semester" | "yearly"

export interface Goal {
  id: string
  user_id?: string
  owner_id?: string
  owner_name?: string
  title: string
  description: string
  target_amount: number
  fixed_contribution_amount?: number | null // Fixed amount each contributor must pay (null = any amount allowed)
  current_amount: number
  currency?: string
  contributor_count: number
  status: GoalStatus
  is_public: boolean
  deadline?: string
  created_at: string
  updated_at: string
  bank_name?: string
  bank_account_number?: string
  bank_account_name?: string
  deposit_bank_name?: string
  deposit_account_number?: string
  deposit_account_name?: string
  total_withdrawn?: number
  milestones?: Milestone[]
  proofs?: Proof[]
}

export interface GoalCreateRequest {
  title: string
  description: string
  target_amount: number
  fixed_contribution_amount?: number | null // Fixed amount each contributor must pay (null = any amount allowed)
  currency?: string
  is_public: boolean
  deadline?: string
  bank_name?: string
  bank_account_number?: string
  bank_account_name?: string
}

export interface GoalListResponse {
  goals: Goal[]
  total: number
  page: number
  limit: number
  has_more: boolean
}

export interface GoalDetailResponse {
  goal: Goal
  milestones: Milestone[]
  contributions: Contribution[]
  proofs: Proof[]
  withdrawals: Withdrawal[]
  is_owner: boolean
  user_contribution?: number
}

export interface Milestone {
  id: string
  goal_id: string
  title: string
  description: string
  target_amount: number
  current_amount: number
  recurrence: MilestoneRecurrence
  order: number
  completed: boolean
  completed_at?: string
  created_at: string
}

export interface MilestoneCreateRequest {
  title: string
  description: string
  target_amount: number
  recurrence: MilestoneRecurrence
}

// Contribution Types
export type ContributionStatus = "pending" | "verified" | "confirmed" | "failed" | "refunded"

export interface Contribution {
  id: string
  goal_id: string
  goal_title?: string
  user_id: string
  user_email?: string
  user_name?: string
  amount: number
  status: ContributionStatus
  message?: string
  is_anonymous?: boolean
  milestone_id?: string
  payment_reference?: string
  created_at: string
  updated_at?: string
  verified_at?: string
}

export interface ContributionListResponse {
  contributions: Contribution[]
  total: number
  page: number
  limit: number
  has_more: boolean
}

// Withdrawal Types
export type WithdrawalStatus = "pending" | "processing" | "completed" | "failed"

export interface Withdrawal {
  id: string
  goal_id: string
  amount: number
  status: WithdrawalStatus
  bank_name: string
  bank_account_number: string
  bank_account_name: string
  milestone_id?: string
  created_at: string
  completed_at?: string
}

export interface WithdrawalRequest {
  amount: number
  milestone_id?: string
  bank_name?: string
  bank_account_number?: string
  bank_account_name?: string
}

// Proof Types
export interface Proof {
  id: string
  goal_id: string
  milestone_id?: string
  description: string
  attachments: string[]
  votes_satisfied: number
  votes_unsatisfied: number
  user_vote?: boolean
  created_at: string
}

export interface ProofSubmitRequest {
  description: string
  attachments?: string[]
  milestone_id?: string
}

// Refund Types
export interface RefundRequest {
  percentage: number
  reason: string
}

// Notification Types
export type NotificationType = 
  | "contribution_received"
  | "contribution_verified"
  | "withdrawal_completed"
  | "goal_funded"
  | "proof_submitted"
  | "refund_initiated"
  | "refund_completed"
  | "goal_closed"
  | "system"

export interface Notification {
  id: string
  user_id: string
  type: NotificationType
  title: string
  message: string
  data?: Record<string, unknown>
  read: boolean
  created_at: string
}

export interface NotificationListResponse {
  notifications: Notification[]
  total: number
  page: number
  limit: number
  unread_count: number
}

// Dashboard Stats
export interface DashboardStats {
  total_contributed: number
  total_raised: number
  active_goals: number
  total_contributors: number
  pending_withdrawals: number
  recent_activity: ActivityItem[]
}

export interface ActivityItem {
  id: string
  type: "contribution" | "withdrawal" | "goal_created" | "proof_submitted"
  description: string
  amount?: number
  goal_title?: string
  created_at: string
}
