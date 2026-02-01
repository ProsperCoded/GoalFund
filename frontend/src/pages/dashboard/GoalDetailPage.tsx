import { useState, useEffect } from "react"
import { useParams, useNavigate, useSearchParams } from "react-router-dom"
import { motion } from "framer-motion"
import {
  ArrowLeft,
  Users,
  Calendar,
  Share2,
  Target,
  Loader2,
  Heart,
  Globe,
  Lock,
  CheckCircle,
  Clock,
  MoreVertical,
  Edit,
  XCircle,
  Banknote,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { ContributeModal } from "@/components/dashboard/ContributeModal"
import { goalsApi, paymentsApi, type Goal, type Contribution } from "@/lib/api"
import { useAuth } from "@/contexts"
import { useToast } from "@/hooks/use-toast"
import { cn, formatCurrency, formatDate, formatDistanceToNow } from "@/lib/utils"

export function GoalDetailPage() {
  const { goalId } = useParams<{ goalId: string }>()
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const { user } = useAuth()
  const { toast } = useToast()

  const [goal, setGoal] = useState<Goal | null>(null)
  const [contributions, setContributions] = useState<Contribution[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isContributeOpen, setIsContributeOpen] = useState(false)
  const [isVerifyingPayment, setIsVerifyingPayment] = useState(false)

  const isOwner = goal?.user_id === user?.id || goal?.owner_id === user?.id

  useEffect(() => {
    if (goalId) {
      fetchGoalDetails()
    }
  }, [goalId])

  // Handle payment verification on return from Paystack
  useEffect(() => {
    const reference = searchParams.get("reference") || searchParams.get("trxref")
    const isPaymentReturn = searchParams.get("payment") === "true"

    if (reference && isPaymentReturn) {
      verifyPayment(reference)
      // Clean URL params
      navigate(`/dashboard/goals/${goalId}`, { replace: true })
    }
  }, [searchParams, goalId])

  const fetchGoalDetails = async () => {
    setIsLoading(true)
    try {
      const response = await goalsApi.getGoalById(goalId!)
      setGoal(response.goal || response)
      setContributions(response.contributions || [])
    } catch (error) {
      console.error("Failed to fetch goal:", error)
      setGoal(null)
      setContributions([])
    } finally {
      setIsLoading(false)
    }
  }

  const verifyPayment = async (reference: string) => {
    setIsVerifyingPayment(true)
    try {
      const result = await paymentsApi.verify(reference)

      if (result.status === "VERIFIED") {
        toast({
          title: "Contribution successful! 🎉",
          description: "Thank you for your contribution!",
        })
        // Refresh goal data
        fetchGoalDetails()
      } else if (result.status === "FAILED") {
        toast({
          variant: "destructive",
          title: "Payment failed",
          description: "Your payment was not successful. Please try again.",
        })
      }
    } catch (error) {
      console.error("Payment verification error:", error)
    } finally {
      setIsVerifyingPayment(false)
    }
  }

  const handleShare = async () => {
    const url = window.location.href
    if (navigator.share) {
      try {
        await navigator.share({
          title: goal?.title,
          text: goal?.description,
          url,
        })
      } catch (error) {
        // User cancelled or error
      }
    } else {
      await navigator.clipboard.writeText(url)
      toast({
        title: "Link copied!",
        description: "Goal link has been copied to clipboard.",
      })
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!goal) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[400px] text-center">
        <Target className="w-12 h-12 text-muted-foreground mb-4" />
        <h2 className="text-xl font-semibold mb-2">Goal not found</h2>
        <p className="text-muted-foreground mb-4">
          This goal may have been removed or doesn't exist.
        </p>
        <Button onClick={() => navigate("/dashboard/goals")}>
          Back to Goals
        </Button>
      </div>
    )
  }

  const progress = Math.min((goal.current_amount / goal.target_amount) * 100, 100)
  const isOverfunded = goal.current_amount > goal.target_amount
  const daysLeft = goal.deadline
    ? Math.max(
        0,
        Math.ceil(
          (new Date(goal.deadline).getTime() - Date.now()) / (1000 * 60 * 60 * 24)
        )
      )
    : null

  return (
    <div className="max-w-4xl mx-auto">
      {/* Back Button */}
      <Button
        variant="ghost"
        size="sm"
        onClick={() => navigate(-1)}
        className="mb-6 -ml-2"
      >
        <ArrowLeft className="w-4 h-4 mr-2" />
        Back
      </Button>

      {/* Payment Verification Overlay */}
      {isVerifyingPayment && (
        <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center">
          <div className="bg-card p-6 rounded-lg text-center">
            <Loader2 className="w-8 h-8 animate-spin text-primary mx-auto mb-4" />
            <p className="font-medium">Verifying your payment...</p>
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Goal Header */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="bg-card border border-border rounded-lg p-6"
          >
            <div className="flex items-start justify-between gap-4 mb-4">
              <div className="flex items-center gap-2">
                <span
                  className={cn(
                    "px-2 py-0.5 text-xs font-medium rounded-full border capitalize",
                    goal.status === "open"
                      ? "bg-green-500/10 text-green-500 border-green-500/20"
                      : "bg-gray-500/10 text-gray-500 border-gray-500/20"
                  )}
                >
                  {goal.status}
                </span>
                {goal.is_public ? (
                  <span className="flex items-center gap-1 text-xs text-muted-foreground">
                    <Globe className="w-3 h-3" />
                    Public
                  </span>
                ) : (
                  <span className="flex items-center gap-1 text-xs text-muted-foreground">
                    <Lock className="w-3 h-3" />
                    Private
                  </span>
                )}
              </div>

              {isOwner && (
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon">
                      <MoreVertical className="w-4 h-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem>
                      <Edit className="w-4 h-4 mr-2" />
                      Edit Goal
                    </DropdownMenuItem>
                    <DropdownMenuItem>
                      <Banknote className="w-4 h-4 mr-2" />
                      Withdraw Funds
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem className="text-destructive">
                      <XCircle className="w-4 h-4 mr-2" />
                      Close Goal
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              )}
            </div>

            <h1 className="text-2xl font-bold mb-3">{goal.title}</h1>
            <p className="text-muted-foreground">{goal.description}</p>

            {/* Meta Info */}
            <div className="flex flex-wrap items-center gap-4 mt-6 pt-4 border-t border-border text-sm text-muted-foreground">
              <div className="flex items-center gap-1.5">
                <Users className="w-4 h-4" />
                <span>{goal.contributor_count} contributors</span>
              </div>
              {goal.deadline && (
                <div className="flex items-center gap-1.5">
                  <Calendar className="w-4 h-4" />
                  <span>
                    {daysLeft === 0
                      ? "Ends today"
                      : daysLeft === 1
                      ? "1 day left"
                      : `${daysLeft} days left`}
                  </span>
                </div>
              )}
              <div className="flex items-center gap-1.5">
                <Clock className="w-4 h-4" />
                <span>Created {formatDate(goal.created_at)}</span>
              </div>
            </div>
          </motion.div>

          {/* Recent Contributors */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 }}
            className="bg-card border border-border rounded-lg"
          >
            <div className="p-4 border-b border-border">
              <h2 className="font-semibold">Recent Contributions</h2>
            </div>
            <div className="divide-y divide-border">
              {contributions.length === 0 ? (
                <div className="p-8 text-center">
                  <Heart className="w-10 h-10 text-muted-foreground mx-auto mb-2" />
                  <p className="text-sm text-muted-foreground">
                    No contributions yet. Be the first to contribute!
                  </p>
                </div>
              ) : (
                contributions.slice(0, 10).map((contribution) => (
                  <div
                    key={contribution.id}
                    className="flex items-center justify-between p-4"
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
                        <span className="text-primary font-semibold text-sm">
                          {contribution.user_name?.[0]?.toUpperCase() || "A"}
                        </span>
                      </div>
                      <div>
                        <p className="font-medium text-sm">
                          {contribution.user_name || "Anonymous"}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          {formatDistanceToNow(contribution.created_at)}
                        </p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="font-semibold">
                        {formatCurrency(contribution.amount)}
                      </p>
                      {contribution.status === "verified" && (
                        <CheckCircle className="w-4 h-4 text-green-500 ml-auto" />
                      )}
                    </div>
                  </div>
                ))
              )}
            </div>
          </motion.div>
        </div>

        {/* Sidebar - Progress & Actions */}
        <div className="space-y-6">
          {/* Progress Card */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="bg-card border border-border rounded-lg p-6 sticky top-24"
          >
            {/* Amount Raised */}
            <div className="mb-4">
              <p className="text-3xl font-bold">
                {formatCurrency(goal.current_amount)}
                {isOverfunded && (
                  <span className="text-green-500 text-lg ml-1">
                    (+{formatCurrency(goal.current_amount - goal.target_amount)})
                  </span>
                )}
              </p>
              <p className="text-sm text-muted-foreground">
                raised of {formatCurrency(goal.target_amount)} goal
              </p>
            </div>

            {/* Progress Bar */}
            <div className="mb-4">
              <div className="h-3 bg-muted rounded-full overflow-hidden">
                <div
                  className={cn(
                    "h-full rounded-full transition-all",
                    isOverfunded ? "bg-green-500" : "bg-primary"
                  )}
                  style={{ width: `${Math.min(progress, 100)}%` }}
                />
              </div>
              <p className="text-sm text-muted-foreground mt-1">
                {progress.toFixed(0)}% funded
              </p>
            </div>

            {/* Action Buttons */}
            <div className="space-y-3">
              {goal.status === "open" && (
                <Button
                  onClick={() => setIsContributeOpen(true)}
                  className="w-full gap-2"
                  size="lg"
                >
                  <Heart className="w-4 h-4" />
                  Contribute Now
                </Button>
              )}

              <Button
                variant="outline"
                onClick={handleShare}
                className="w-full gap-2"
              >
                <Share2 className="w-4 h-4" />
                Share Goal
              </Button>

              {isOwner && goal.current_amount > 0 && (
                <Button variant="outline" className="w-full gap-2">
                  <Banknote className="w-4 h-4" />
                  Withdraw Funds
                </Button>
              )}
            </div>

            {/* Quick Stats */}
            <div className="grid grid-cols-2 gap-4 mt-6 pt-4 border-t border-border">
              <div>
                <p className="text-2xl font-bold">{goal.contributor_count}</p>
                <p className="text-xs text-muted-foreground">Contributors</p>
              </div>
              <div>
                <p className="text-2xl font-bold">
                  {daysLeft !== null ? daysLeft : "∞"}
                </p>
                <p className="text-xs text-muted-foreground">Days Left</p>
              </div>
            </div>
          </motion.div>
        </div>
      </div>

      {/* Contribute Modal */}
      <ContributeModal
        isOpen={isContributeOpen}
        onClose={() => setIsContributeOpen(false)}
        goal={goal}
        onSuccess={() => {
          setIsContributeOpen(false)
          fetchGoalDetails()
        }}
      />
    </div>
  )
}


