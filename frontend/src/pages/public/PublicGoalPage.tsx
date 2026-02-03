import { useState, useEffect } from "react"
import { useParams, useNavigate, useSearchParams, Link } from "react-router-dom"
import { motion } from "framer-motion"
import {
  Target,
  Users,
  Calendar,
  Share2,
  Loader2,
  Heart,
  Clock,
  CheckCircle,
  CreditCard,
  ExternalLink,
  AlertCircle,
  ArrowRight,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { goalsApi, paymentsApi, nairaToKobo, type Goal, type Contribution } from "@/lib/api"
import { useAuth } from "@/contexts"
import { useToast } from "@/hooks/use-toast"
import { cn, formatCurrency, formatDate, formatDistanceToNow } from "@/lib/utils"

type ContributionStep = "form" | "processing" | "redirect" | "verifying" | "success" | "error"

export function PublicGoalPage() {
  const { goalId } = useParams<{ goalId: string }>()
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const { user, isAuthenticated } = useAuth()
  const { toast } = useToast()

  const [goal, setGoal] = useState<Goal | null>(null)
  const [contributions, setContributions] = useState<Contribution[]>([])
  const [isLoading, setIsLoading] = useState(true)
  
  // Contribution form state
  const [step, setStep] = useState<ContributionStep>("form")
  const [amount, setAmount] = useState("")
  const [email, setEmail] = useState("")
  const [name, setName] = useState("")
  const [message, setMessage] = useState("")
  const [isAnonymous, setIsAnonymous] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const quickAmounts = [5000, 10000, 25000, 50000, 100000]

  useEffect(() => {
    if (goalId) {
      fetchGoalDetails()
    }
  }, [goalId])

  // Pre-fill user info if authenticated
  useEffect(() => {
    if (user) {
      setEmail(user.email || "")
      setName(`${user.first_name || ""} ${user.last_name || ""}`.trim())
    }
  }, [user])

  // Handle payment verification on return from Paystack
  useEffect(() => {
    const reference = searchParams.get("reference") || searchParams.get("trxref")

    if (reference) {
      setStep("verifying")
      verifyPayment(reference)
      // Clean URL params
      navigate(`/g/${goalId}`, { replace: true })
    }
  }, [searchParams, goalId])

  const fetchGoalDetails = async () => {
    setIsLoading(true)
    try {
      // Use public endpoint to fetch goal
      const response = await goalsApi.getGoalById(goalId!)
      const goalData = response.goal || response
      
      // Check if goal is public
      if (!goalData.is_public) {
        setGoal(null)
        return
      }
      
      setGoal(goalData)
      setContributions(response.contributions || [])
    } catch (error) {
      console.error("Failed to fetch goal:", error)
      setGoal(null)
    } finally {
      setIsLoading(false)
    }
  }

  const verifyPayment = async (reference: string) => {
    try {
      const result = await paymentsApi.verify(reference)

      if (result.status === "VERIFIED") {
        setStep("success")
        toast({
          title: "Contribution successful! 🎉",
          description: "Thank you for your contribution!",
        })
        fetchGoalDetails()
      } else if (result.status === "FAILED") {
        setError("Payment was not successful. Please try again.")
        setStep("error")
      } else {
        setError("Payment is still being processed. Please check back later.")
        setStep("error")
      }
    } catch (err: unknown) {
      const error = err as { response?: { data?: { message?: string } } }
      setError(error.response?.data?.message || "Failed to verify payment")
      setStep("error")
    }
  }

  const handleAmountChange = (value: string) => {
    const numericValue = value.replace(/[^0-9]/g, "")
    setAmount(numericValue)
  }

  const handleQuickAmount = (value: number) => {
    setAmount(value.toString())
  }

  const initializePayment = async () => {
    // Use fixed amount if available, otherwise use user input
    const contributionAmount = goal?.fixed_contribution_amount 
      ? goal.fixed_contribution_amount 
      : parseInt(amount, 10)
    
    if (isNaN(contributionAmount) || contributionAmount < 100) {
      toast({
        variant: "destructive",
        title: "Invalid amount",
        description: "Minimum contribution is ₦100",
      })
      return
    }

    if (!email) {
      toast({
        variant: "destructive",
        title: "Email required",
        description: "Please enter your email address",
      })
      return
    }

    setStep("processing")
    setError(null)

    try {
      const response = await paymentsApi.initialize({
        user_id: user?.id || "", // Empty for guest contributions
        goal_id: goal!.id,
        amount: nairaToKobo(contributionAmount),
        currency: "NGN",
        email: email,
        callback_url: `${window.location.origin}/g/${goal!.id}`,
        metadata: {
          goal_title: goal!.title,
          contributor_name: isAnonymous ? "Anonymous" : name || "Anonymous",
          contributor_email: email,
          message: message || undefined,
          is_anonymous: isAnonymous,
          is_guest: !isAuthenticated,
        },
      })

      setStep("redirect")
      // Redirect to Paystack checkout
      window.location.href = response.authorization_url
    } catch (err: unknown) {
      const error = err as { response?: { data?: { message?: string } } }
      setError(error.response?.data?.message || "Failed to initialize payment")
      setStep("error")
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
        // User cancelled
      }
    } else {
      await navigator.clipboard.writeText(url)
      toast({
        title: "Link copied!",
        description: "Share this link with friends to support this goal.",
      })
    }
  }

  // Calculate contribution amount - use fixed amount if available
  const fixedAmount = goal?.fixed_contribution_amount
  const numericAmount = fixedAmount || parseInt(amount, 10) || 0
  const isValidAmount = numericAmount >= 100

  if (isLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    )
  }

  if (!goal) {
    return (
      <div className="min-h-screen bg-background flex flex-col items-center justify-center text-center p-4">
        <Target className="w-16 h-16 text-muted-foreground mb-4" />
        <h1 className="text-2xl font-bold mb-2">Goal Not Found</h1>
        <p className="text-muted-foreground mb-6">
          This goal doesn't exist or is not available for public contributions.
        </p>
        <Link to="/">
          <Button>Go to Homepage</Button>
        </Link>
      </div>
    )
  }

  const progress = Math.min(((goal.current_amount || 0) / (goal.target_amount || 1)) * 100, 100)
  const isOverfunded = (goal.current_amount || 0) > (goal.target_amount || 1)
  const statusLower = goal.status?.toLowerCase() || "open"
  const daysLeft = goal.deadline
    ? Math.max(
        0,
        Math.ceil(
          (new Date(goal.deadline).getTime() - Date.now()) / (1000 * 60 * 60 * 24)
        )
      )
    : null

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="sticky top-0 z-50 bg-background/95 backdrop-blur border-b border-border">
        <div className="max-w-5xl mx-auto px-4 py-4 flex items-center justify-between">
          <Link to="/" className="flex items-center gap-2">
            <Target className="w-6 h-6 text-primary" />
            <span className="font-bold text-xl">GoalFund</span>
          </Link>
          {!isAuthenticated && (
            <div className="flex items-center gap-2">
              <Link to="/login">
                <Button variant="ghost" size="sm">Login</Button>
              </Link>
              <Link to="/register">
                <Button size="sm">Sign Up</Button>
              </Link>
            </div>
          )}
          {isAuthenticated && (
            <Link to="/dashboard">
              <Button variant="outline" size="sm">
                Dashboard <ArrowRight className="w-4 h-4 ml-2" />
              </Button>
            </Link>
          )}
        </div>
      </header>

      <main className="max-w-5xl mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-5 gap-8">
          {/* Goal Info - Left Column */}
          <div className="lg:col-span-3 space-y-6">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="bg-card border border-border rounded-lg p-6"
            >
              <div className="flex items-center gap-2 mb-4">
                <span
                  className={cn(
                    "px-2 py-0.5 text-xs font-medium rounded-full border capitalize",
                    statusLower === "open"
                      ? "bg-green-500/10 text-green-500 border-green-500/20"
                      : "bg-gray-500/10 text-gray-500 border-gray-500/20"
                  )}
                >
                  {statusLower}
                </span>
              </div>

              <h1 className="text-2xl md:text-3xl font-bold mb-4">{goal.title}</h1>
              <p className="text-muted-foreground whitespace-pre-wrap">{goal.description}</p>

              {/* Progress */}
              <div className="mt-6 pt-6 border-t border-border">
                <div className="flex items-end justify-between mb-3">
                  <div>
                    <p className="text-3xl font-bold">
                      {formatCurrency(goal.current_amount || 0)}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      raised of {formatCurrency(goal.target_amount || 0)}
                    </p>
                  </div>
                  <p className="text-lg font-semibold text-primary">
                    {isNaN(progress) ? 0 : progress.toFixed(0)}%
                  </p>
                </div>
                <div className="h-3 bg-muted rounded-full overflow-hidden">
                  <div
                    className={cn(
                      "h-full rounded-full transition-all",
                      isOverfunded ? "bg-green-500" : "bg-primary"
                    )}
                    style={{ width: `${Math.min(isNaN(progress) ? 0 : progress, 100)}%` }}
                  />
                </div>
              </div>

              {/* Meta */}
              <div className="flex flex-wrap items-center gap-4 mt-6 pt-4 border-t border-border text-sm text-muted-foreground">
                <div className="flex items-center gap-1.5">
                  <Users className="w-4 h-4" />
                  <span>{goal.contributor_count || 0} contributors</span>
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

            {/* Recent Contributions */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 }}
              className="bg-card border border-border rounded-lg"
            >
              <div className="p-4 border-b border-border">
                <h2 className="font-semibold">Recent Supporters</h2>
              </div>
              <div className="divide-y divide-border">
                {contributions.length === 0 ? (
                  <div className="p-8 text-center">
                    <Heart className="w-10 h-10 text-muted-foreground mx-auto mb-2" />
                    <p className="text-sm text-muted-foreground">
                      Be the first to support this goal!
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

          {/* Contribution Form - Right Column */}
          <div className="lg:col-span-2">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.2 }}
              className="bg-card border border-border rounded-lg p-6 sticky top-24"
            >
              {/* Step: Form */}
              {step === "form" && (
                <>
                  <h2 className="font-semibold text-lg mb-4">Contribute to this Goal</h2>

                  {statusLower !== "open" ? (
                    <div className="text-center py-8">
                      <AlertCircle className="w-10 h-10 text-muted-foreground mx-auto mb-2" />
                      <p className="text-muted-foreground">
                        This goal is no longer accepting contributions.
                      </p>
                    </div>
                  ) : (
                    <div className="space-y-4">
                      {/* Fixed Contribution Amount Notice */}
                      {goal.fixed_contribution_amount && (
                        <div className="bg-primary/5 border border-primary/20 rounded-lg p-4">
                          <div className="flex items-center gap-2 mb-2">
                            <Users className="w-4 h-4 text-primary" />
                            <span className="text-sm font-medium text-primary">Group Contribution</span>
                          </div>
                          <p className="text-sm text-muted-foreground">
                            This goal requires a fixed contribution of{" "}
                            <span className="font-semibold text-foreground">
                              {formatCurrency(goal.fixed_contribution_amount)}
                            </span>
                          </p>
                        </div>
                      )}

                      {/* Amount Input - Only show if no fixed amount */}
                      {!goal.fixed_contribution_amount ? (
                        <>
                          <div className="space-y-2">
                            <label className="text-sm font-medium">Amount</label>
                            <div className="relative">
                              <span className="absolute left-4 top-1/2 -translate-y-1/2 text-muted-foreground">
                                ₦
                              </span>
                              <input
                                type="text"
                                inputMode="numeric"
                                value={amount}
                                onChange={(e) => handleAmountChange(e.target.value)}
                                placeholder="0"
                                className="w-full pl-8 pr-4 py-3 text-xl font-semibold bg-background border border-input rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
                              />
                            </div>
                            {amount && !isValidAmount && (
                              <p className="text-xs text-destructive">
                                Minimum contribution is ₦100
                              </p>
                            )}
                          </div>

                          {/* Quick Amounts */}
                          <div className="flex flex-wrap gap-2">
                            {quickAmounts.map((value) => (
                              <button
                                key={value}
                                type="button"
                                onClick={() => handleQuickAmount(value)}
                                className={cn(
                                  "px-3 py-1 text-xs rounded-full border transition-colors",
                                  amount === value.toString()
                                    ? "border-primary bg-primary/10 text-primary"
                                    : "border-border hover:border-primary/50"
                                )}
                              >
                                ₦{value.toLocaleString()}
                              </button>
                            ))}
                          </div>
                        </>
                      ) : (
                        /* Fixed Amount Display */
                        <div className="space-y-2">
                          <label className="text-sm font-medium">Contribution Amount</label>
                          <div className="w-full px-4 py-3 text-xl font-semibold bg-muted/50 border border-input rounded-lg text-center">
                            {formatCurrency(goal.fixed_contribution_amount)}
                          </div>
                        </div>
                      )}

                      {/* Email */}
                      <div className="space-y-2">
                        <label className="text-sm font-medium">Email</label>
                        <input
                          type="email"
                          value={email}
                          onChange={(e) => setEmail(e.target.value)}
                          placeholder="your@email.com"
                          className="w-full px-4 py-2 bg-background border border-input rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
                        />
                      </div>

                      {/* Name */}
                      {!isAnonymous && (
                        <div className="space-y-2">
                          <label className="text-sm font-medium">
                            Name (Optional)
                          </label>
                          <input
                            type="text"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            placeholder="Your name"
                            className="w-full px-4 py-2 bg-background border border-input rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
                          />
                        </div>
                      )}

                      {/* Message */}
                      <div className="space-y-2">
                        <label className="text-sm font-medium">
                          Message (Optional)
                        </label>
                        <textarea
                          value={message}
                          onChange={(e) => setMessage(e.target.value)}
                          placeholder="Add an encouraging message..."
                          rows={2}
                          maxLength={200}
                          className="w-full px-4 py-2 bg-background border border-input rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary resize-none"
                        />
                      </div>

                      {/* Anonymous */}
                      <label className="flex items-center gap-3 cursor-pointer">
                        <input
                          type="checkbox"
                          checked={isAnonymous}
                          onChange={(e) => setIsAnonymous(e.target.checked)}
                          className="w-4 h-4 rounded border-input"
                        />
                        <span className="text-sm">Contribute anonymously</span>
                      </label>

                      {/* Submit */}
                      <Button
                        onClick={initializePayment}
                        disabled={!isValidAmount || !email}
                        className="w-full gap-2"
                        size="lg"
                      >
                        <CreditCard className="w-4 h-4" />
                        Contribute {isValidAmount ? formatCurrency(numericAmount) : ""}
                      </Button>

                      <p className="text-xs text-center text-muted-foreground">
                        Secured by Paystack
                      </p>
                    </div>
                  )}

                  {/* Share Button */}
                  <div className="mt-6 pt-4 border-t border-border">
                    <Button
                      variant="outline"
                      onClick={handleShare}
                      className="w-full gap-2"
                    >
                      <Share2 className="w-4 h-4" />
                      Share this Goal
                    </Button>
                  </div>
                </>
              )}

              {/* Step: Processing */}
              {step === "processing" && (
                <div className="flex flex-col items-center justify-center py-12 text-center">
                  <Loader2 className="w-12 h-12 text-primary animate-spin mb-4" />
                  <h3 className="font-semibold text-lg mb-2">
                    Initializing Payment
                  </h3>
                  <p className="text-sm text-muted-foreground">
                    Please wait...
                  </p>
                </div>
              )}

              {/* Step: Redirect */}
              {step === "redirect" && (
                <div className="flex flex-col items-center justify-center py-12 text-center">
                  <ExternalLink className="w-12 h-12 text-primary mb-4" />
                  <h3 className="font-semibold text-lg mb-2">
                    Redirecting to Paystack
                  </h3>
                  <p className="text-sm text-muted-foreground mb-4">
                    You'll be redirected to complete your payment securely.
                  </p>
                  <Loader2 className="w-6 h-6 animate-spin text-muted-foreground" />
                </div>
              )}

              {/* Step: Verifying */}
              {step === "verifying" && (
                <div className="flex flex-col items-center justify-center py-12 text-center">
                  <Loader2 className="w-12 h-12 text-primary animate-spin mb-4" />
                  <h3 className="font-semibold text-lg mb-2">
                    Verifying Payment
                  </h3>
                  <p className="text-sm text-muted-foreground">
                    Please wait while we confirm your payment...
                  </p>
                </div>
              )}

              {/* Step: Success */}
              {step === "success" && (
                <div className="flex flex-col items-center justify-center py-12 text-center">
                  <div className="w-16 h-16 rounded-full bg-green-500/10 flex items-center justify-center mb-4">
                    <CheckCircle className="w-10 h-10 text-green-500" />
                  </div>
                  <h3 className="font-semibold text-lg mb-2">
                    Thank You! 🎉
                  </h3>
                  <p className="text-sm text-muted-foreground mb-4">
                    Your contribution of {formatCurrency(numericAmount)} has been received.
                  </p>
                  <Heart className="w-6 h-6 text-primary fill-current" />
                  <div className="mt-6 space-y-2 w-full">
                    <Button onClick={handleShare} variant="outline" className="w-full gap-2">
                      <Share2 className="w-4 h-4" />
                      Share with Friends
                    </Button>
                    <Button onClick={() => setStep("form")} variant="ghost" className="w-full">
                      Contribute Again
                    </Button>
                  </div>
                </div>
              )}

              {/* Step: Error */}
              {step === "error" && (
                <div className="flex flex-col items-center justify-center py-12 text-center">
                  <div className="w-16 h-16 rounded-full bg-destructive/10 flex items-center justify-center mb-4">
                    <AlertCircle className="w-10 h-10 text-destructive" />
                  </div>
                  <h3 className="font-semibold text-lg mb-2">
                    Payment Failed
                  </h3>
                  <p className="text-sm text-muted-foreground mb-6">
                    {error || "Something went wrong. Please try again."}
                  </p>
                  <Button onClick={() => setStep("form")} className="w-full">
                    Try Again
                  </Button>
                </div>
              )}
            </motion.div>
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t border-border mt-12 py-8">
        <div className="max-w-5xl mx-auto px-4 text-center text-sm text-muted-foreground">
          <p>Powered by <Link to="/" className="text-primary hover:underline">GoalFund</Link></p>
        </div>
      </footer>
    </div>
  )
}
