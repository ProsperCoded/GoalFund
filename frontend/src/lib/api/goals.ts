import apiClient from "./config"
import type {
  Goal,
  GoalCreateRequest,
  GoalListResponse,
  Milestone,
  MilestoneCreateRequest,
  Contribution,
  ContributionListResponse,
  WithdrawalRequest,
  ProofSubmitRequest,
  RefundRequest,
  GoalDetailResponse,
} from "./types"

export const goalsApi = {
  // Create a new goal
  create: async (data: GoalCreateRequest): Promise<Goal> => {
    // Map frontend payload to backend expected fields (no json tags in DTOs)
    const payload: any = {
      Title: data.title,
      Description: data.description,
      TargetAmount: data.target_amount,
      FixedContributionAmount: data.fixed_contribution_amount || null,
      Currency: data.currency,
      Deadline: data.deadline,
      BankName: data.bank_name,
      AccountNumber: data.bank_account_number,
      AccountName: data.bank_account_name,
      IsPublic: data.is_public,
      // Milestones mapping can be added when UI supports it
    }
    const response = await apiClient.post<Goal>("/goals", payload)
    return response.data
  },

  // Get user's created goals
  getMyGoals: async (params?: { status?: string; page?: number; limit?: number }): Promise<GoalListResponse> => {
    const response = await apiClient.get<any>("/goals/my", { params })
    const { goals, total, page, limit } = response.data
    return {
      goals,
      total,
      page,
      limit,
      has_more: total > page * limit,
    }
  },

  // Get public goals (for explore page)
  getPublicGoals: async (params?: { 
    search?: string
    page?: number
    limit?: number 
    sort?: string
  }): Promise<GoalListResponse> => {
    const response = await apiClient.get<any>("/goals/list", { params })
    const { data, total, page, size } = response.data
    return {
      goals: data,
      total,
      page,
      limit: size,
      has_more: total > page * size,
    }
  },

  // Get single goal details (public view)
  getGoal: async (id: string): Promise<GoalDetailResponse> => {
    const response = await apiClient.get<GoalDetailResponse>(`/goals/view/${id}`)
    return response.data
  },

  // Get single goal by ID (alias for convenience)
  getGoalById: async (id: string): Promise<GoalDetailResponse> => {
    const response = await apiClient.get<GoalDetailResponse>(`/goals/view/${id}`)
    return response.data
  },

  // Get single goal details (owner view - includes more data)
  getGoalAsOwner: async (id: string): Promise<GoalDetailResponse> => {
    const response = await apiClient.get<GoalDetailResponse>(`/goals/${id}`)
    return response.data
  },

  // Update goal
  update: async (id: string, data: Partial<GoalCreateRequest>): Promise<Goal> => {
    const payload: any = {
      Title: data.title,
      Description: data.description,
      TargetAmount: data.target_amount,
      Deadline: data.deadline,
      BankName: data.bank_name,
      AccountNumber: data.bank_account_number,
      AccountName: data.bank_account_name,
      IsPublic: data.is_public,
    }
    const response = await apiClient.patch<Goal>(`/goals/${id}`, payload)
    return response.data
  },

  // Update goal (alias for edit page)
  updateGoal: async (id: string, data: {
    title?: string
    description?: string
    target_amount?: number
    fixed_contribution_amount?: number | null
    deadline?: string
    is_public?: boolean
    deposit_bank_name?: string
    deposit_account_number?: string
    deposit_account_name?: string
  }): Promise<Goal> => {
    const payload: any = {
      Title: data.title,
      Description: data.description,
      TargetAmount: data.target_amount,
      FixedContributionAmount: data.fixed_contribution_amount,
      Deadline: data.deadline,
      IsPublic: data.is_public,
      BankName: data.deposit_bank_name,
      AccountNumber: data.deposit_account_number,
      AccountName: data.deposit_account_name,
    }
    // Remove undefined values
    Object.keys(payload).forEach(key => payload[key] === undefined && delete payload[key])
    const response = await apiClient.patch<Goal>(`/goals/${id}`, payload)
    return response.data
  },

  // Close goal
  close: async (id: string): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>(`/goals/${id}/close`)
    return response.data
  },

  // Cancel goal
  cancel: async (id: string, reason?: string): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>(`/goals/${id}/cancel`, { reason })
    return response.data
  },

  // Request withdrawal
  withdraw: async (id: string, data: WithdrawalRequest): Promise<any> => {
    const payload: any = {
      GoalID: id,
      Amount: data.amount,
      MilestoneID: data.milestone_id,
      BankName: data.bank_name,
      AccountNumber: data.bank_account_number,
      AccountName: data.bank_account_name,
    }
    const response = await apiClient.post(`/goals/withdraw`, payload)
    return response.data
  },

  // Submit proof
  submitProof: async (id: string, data: ProofSubmitRequest & { title?: string }): Promise<any> => {
    const payload: any = {
      GoalID: id,
      MilestoneID: data.milestone_id,
      Title: data["title"] || "",
      Description: data.description,
      MediaURLs: data.attachments,
    }
    const response = await apiClient.post(`/goals/proofs`, payload)
    return response.data
  },

  // Initiate refund
  refund: async (id: string, data: RefundRequest): Promise<any> => {
    const payload = {
      goal_id: id,
      refund_percentage: data.percentage,
      reason: data.reason,
    }
    const response = await apiClient.post(`/goals/refunds`, payload)
    return response.data
  },

  // Add milestone
  addMilestone: async (goalId: string, data: MilestoneCreateRequest): Promise<Milestone> => {
    const response = await apiClient.post<Milestone>(`/goals/${goalId}/milestones`, data)
    return response.data
  },

  // Get goal milestones
  getMilestones: async (goalId: string): Promise<Milestone[]> => {
    const response = await apiClient.get<{ milestones: Milestone[] }>(`/goals/${goalId}/milestones`)
    return response.data.milestones
  },

  // Vote on proof
  voteOnProof: async (_goalId: string, proofId: string, vote: boolean, comment?: string): Promise<any> => {
    const payload = {
      ProofID: proofId,
      IsSatisfied: vote,
      Comment: comment || "",
    }
    const response = await apiClient.post(`/goals/votes`, payload)
    return response.data
  },
}

export const contributionsApi = {
  // Get user's contributions
  getMyContributions: async (params?: { 
    status?: string
    page?: number
    limit?: number 
  }): Promise<ContributionListResponse> => {
    const response = await apiClient.get<any>("/contributions/my", { params })
    const { contributions, total } = response.data
    const page = params?.page ?? 1
    const limit = params?.limit ?? total
    return {
      contributions,
      total,
      page,
      limit,
      has_more: total > page * limit,
    }
  },

  // Make a contribution (initialize payment)
  contribute: async (goalId: string, amount: number, milestoneId?: string): Promise<Contribution> => {
    const payload: any = {
      GoalID: goalId,
      Amount: amount,
      MilestoneID: milestoneId,
    }
    const response = await apiClient.post<Contribution>("/contributions", payload)
    return response.data
  },

  // Get contribution details
  getContribution: async (id: string): Promise<Contribution> => {
    const response = await apiClient.get<Contribution>(`/contributions/${id}`)
    return response.data
  },
}
