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
    const response = await apiClient.post<Goal>("/goals", data)
    return response.data
  },

  // Get user's created goals
  getMyGoals: async (params?: { status?: string; page?: number; limit?: number }): Promise<GoalListResponse> => {
    const response = await apiClient.get<GoalListResponse>("/goals/my", { params })
    return response.data
  },

  // Get public goals (for explore page)
  getPublicGoals: async (params?: { 
    search?: string
    page?: number
    limit?: number 
    sort?: string
  }): Promise<GoalListResponse> => {
    const response = await apiClient.get<GoalListResponse>("/goals/list", { params })
    return response.data
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
    const response = await apiClient.put<Goal>(`/goals/${id}`, data)
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
  withdraw: async (id: string, data: WithdrawalRequest): Promise<{ message: string; withdrawal_id: string }> => {
    const response = await apiClient.post<{ message: string; withdrawal_id: string }>(`/goals/${id}/withdraw`, data)
    return response.data
  },

  // Submit proof
  submitProof: async (id: string, data: ProofSubmitRequest): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>(`/goals/${id}/proof`, data)
    return response.data
  },

  // Initiate refund
  refund: async (id: string, data: RefundRequest): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>(`/goals/${id}/refund`, data)
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
  voteOnProof: async (goalId: string, proofId: string, vote: boolean): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>(`/goals/${goalId}/proofs/${proofId}/vote`, { 
      satisfied: vote 
    })
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
    const response = await apiClient.get<ContributionListResponse>("/contributions/my", { params })
    return response.data
  },

  // Make a contribution (initialize payment)
  contribute: async (goalId: string, amount: number, milestoneId?: string): Promise<{
    contribution_id: string
    payment_url: string
    reference: string
  }> => {
    const response = await apiClient.post<{
      contribution_id: string
      payment_url: string
      reference: string
    }>("/contributions", {
      goal_id: goalId,
      amount,
      milestone_id: milestoneId,
    })
    return response.data
  },

  // Get contribution details
  getContribution: async (id: string): Promise<Contribution> => {
    const response = await apiClient.get<Contribution>(`/contributions/${id}`)
    return response.data
  },
}
