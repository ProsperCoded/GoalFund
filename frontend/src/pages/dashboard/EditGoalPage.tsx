import { useState, useEffect } from "react"
import { useParams, useNavigate } from "react-router-dom"
import { motion } from "framer-motion"
import {
  ArrowLeft,
  Loader2,
  Save,
  Target,
  Globe,
  Lock,
  AlertCircle,
  Users,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { goalsApi, type Goal } from "@/lib/api"
import { useAuth } from "@/contexts"
import { useToast } from "@/hooks/use-toast"
import { cn } from "@/lib/utils"

export function EditGoalPage() {
  const { goalId } = useParams<{ goalId: string }>()
  const navigate = useNavigate()
  const { user } = useAuth()
  const { toast } = useToast()

  const [isLoading, setIsLoading] = useState(true)
  const [isSaving, setIsSaving] = useState(false)
  const [goal, setGoal] = useState<Goal | null>(null)
  const [errors, setErrors] = useState<Record<string, string>>({})

  // Form state
  const [title, setTitle] = useState("")
  const [description, setDescription] = useState("")
  const [targetAmount, setTargetAmount] = useState("")
  const [deadline, setDeadline] = useState("")
  const [isPublic, setIsPublic] = useState(true)
  const [depositBankName, setDepositBankName] = useState("")
  const [depositAccountNumber, setDepositAccountNumber] = useState("")
  const [depositAccountName, setDepositAccountName] = useState("")
  const [fixedContributionAmount, setFixedContributionAmount] = useState("")
  const [hasFixedAmount, setHasFixedAmount] = useState(false)

  useEffect(() => {
    if (goalId) {
      fetchGoal()
    }
  }, [goalId])

  const fetchGoal = async () => {
    setIsLoading(true)
    try {
      const response = await goalsApi.getGoalById(goalId!)
      const goalData = response.goal || response

      // Check ownership
      const ownerId = goalData.owner_id || goalData.user_id
      if (ownerId !== user?.id) {
        toast({
          variant: "destructive",
          title: "Access denied",
          description: "You can only edit your own goals.",
        })
        navigate("/dashboard/goals")
        return
      }

      setGoal(goalData)
      
      // Populate form
      setTitle(goalData.title || "")
      setDescription(goalData.description || "")
      setTargetAmount(goalData.target_amount?.toString() || "")
      setDeadline(goalData.deadline ? goalData.deadline.split("T")[0] : "")
      setIsPublic(goalData.is_public ?? true)
      setDepositBankName(goalData.deposit_bank_name || "")
      setDepositAccountNumber(goalData.deposit_account_number || "")
      setDepositAccountName(goalData.deposit_account_name || "")
      
      // Fixed contribution amount
      if (goalData.fixed_contribution_amount) {
        setHasFixedAmount(true)
        setFixedContributionAmount(goalData.fixed_contribution_amount.toString())
      } else {
        setHasFixedAmount(false)
        setFixedContributionAmount("")
      }
    } catch (error) {
      console.error("Failed to fetch goal:", error)
      toast({
        variant: "destructive",
        title: "Failed to load goal",
        description: "Please try again.",
      })
      navigate("/dashboard/goals")
    } finally {
      setIsLoading(false)
    }
  }

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!title.trim()) {
      newErrors.title = "Title is required"
    }

    if (!description.trim()) {
      newErrors.description = "Description is required"
    }

    const amount = parseInt(targetAmount.replace(/[^0-9]/g, ""), 10)
    if (isNaN(amount) || amount < 1000) {
      newErrors.targetAmount = "Target amount must be at least ₦1,000"
    }

    // Validate fixed contribution amount if enabled
    if (hasFixedAmount) {
      const fixedAmount = parseInt(fixedContributionAmount.replace(/[^0-9]/g, ""), 10)
      if (isNaN(fixedAmount) || fixedAmount < 100) {
        newErrors.fixedContributionAmount = "Fixed contribution amount must be at least ₦100"
      }
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateForm()) {
      return
    }

    setIsSaving(true)

    // Calculate fixed amount: null to remove, number to set
    let fixedAmount: number | null = null
    if (hasFixedAmount && fixedContributionAmount) {
      fixedAmount = parseInt(fixedContributionAmount.replace(/[^0-9]/g, ""), 10)
    } else if (!hasFixedAmount && goal?.fixed_contribution_amount) {
      // User disabled fixed amount - send 0 to remove it
      fixedAmount = 0
    }

    try {
      await goalsApi.updateGoal(goalId!, {
        title,
        description,
        target_amount: parseInt(targetAmount.replace(/[^0-9]/g, ""), 10),
        deadline: deadline || undefined,
        is_public: isPublic,
        deposit_bank_name: depositBankName || undefined,
        deposit_account_number: depositAccountNumber || undefined,
        deposit_account_name: depositAccountName || undefined,
        fixed_contribution_amount: fixedAmount,
      })

      toast({
        title: "Goal updated!",
        description: "Your changes have been saved.",
      })

      navigate(`/dashboard/goals/${goalId}`)
    } catch (error) {
      console.error("Failed to update goal:", error)
      toast({
        variant: "destructive",
        title: "Failed to update goal",
        description: "Please try again.",
      })
    } finally {
      setIsSaving(false)
    }
  }

  const handleAmountChange = (value: string) => {
    const numericValue = value.replace(/[^0-9]/g, "")
    setTargetAmount(numericValue)
  }

  const handleFixedAmountChange = (value: string) => {
    const numericValue = value.replace(/[^0-9]/g, "")
    setFixedContributionAmount(numericValue)
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

  return (
    <div className="max-w-2xl mx-auto">
      {/* Header */}
      <div className="mb-6">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => navigate(-1)}
          className="-ml-2 mb-4"
        >
          <ArrowLeft className="w-4 h-4 mr-2" />
          Back
        </Button>
        <h1 className="text-2xl font-bold">Edit Goal</h1>
        <p className="text-muted-foreground">Update your goal details</p>
      </div>

      {/* Warning if goal has contributions */}
      {goal.current_amount > 0 && (
        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-6 p-4 bg-yellow-500/10 border border-yellow-500/20 rounded-lg flex items-start gap-3"
        >
          <AlertCircle className="w-5 h-5 text-yellow-500 flex-shrink-0 mt-0.5" />
          <div className="text-sm">
            <p className="font-medium text-yellow-500">Goal has contributions</p>
            <p className="text-muted-foreground">
              Some fields like target amount cannot be reduced below the current amount raised.
            </p>
          </div>
        </motion.div>
      )}

      {/* Form */}
      <motion.form
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        onSubmit={handleSubmit}
        className="bg-card border border-border rounded-lg p-6 space-y-6"
      >
        {/* Title */}
        <div className="space-y-2">
          <label className="text-sm font-medium">Goal Title</label>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="e.g., New Laptop Fund"
            className={cn(
              "w-full px-4 py-2 bg-background border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary",
              errors.title ? "border-destructive" : "border-input"
            )}
          />
          {errors.title && (
            <p className="text-xs text-destructive">{errors.title}</p>
          )}
        </div>

        {/* Description */}
        <div className="space-y-2">
          <label className="text-sm font-medium">Description</label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Tell people about your goal..."
            rows={4}
            className={cn(
              "w-full px-4 py-2 bg-background border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary resize-none",
              errors.description ? "border-destructive" : "border-input"
            )}
          />
          {errors.description && (
            <p className="text-xs text-destructive">{errors.description}</p>
          )}
        </div>

        {/* Target Amount */}
        <div className="space-y-2">
          <label className="text-sm font-medium">Target Amount (₦)</label>
          <div className="relative">
            <span className="absolute left-4 top-1/2 -translate-y-1/2 text-muted-foreground">
              ₦
            </span>
            <input
              type="text"
              inputMode="numeric"
              value={targetAmount}
              onChange={(e) => handleAmountChange(e.target.value)}
              placeholder="100000"
              className={cn(
                "w-full pl-8 pr-4 py-2 bg-background border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary",
                errors.targetAmount ? "border-destructive" : "border-input"
              )}
            />
          </div>
          {errors.targetAmount && (
            <p className="text-xs text-destructive">{errors.targetAmount}</p>
          )}
          {goal.current_amount > 0 && (
            <p className="text-xs text-muted-foreground">
              Current amount raised: ₦{goal.current_amount.toLocaleString()}
            </p>
          )}
        </div>

        {/* Deadline */}
        <div className="space-y-2">
          <label className="text-sm font-medium">Deadline (Optional)</label>
          <input
            type="date"
            value={deadline}
            onChange={(e) => setDeadline(e.target.value)}
            min={new Date().toISOString().split("T")[0]}
            className="w-full px-4 py-2 bg-background border border-input rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
          />
        </div>

        {/* Fixed Contribution Amount (Group Contribution) */}
        <div className="space-y-3">
          <div className="flex items-center gap-3">
            <input
              type="checkbox"
              id="hasFixedAmount"
              checked={hasFixedAmount}
              onChange={(e) => {
                setHasFixedAmount(e.target.checked)
                if (!e.target.checked) {
                  setFixedContributionAmount("")
                }
              }}
              className="w-4 h-4 rounded border-input"
            />
            <label htmlFor="hasFixedAmount" className="flex items-center gap-2 cursor-pointer">
              <Users className="w-4 h-4 text-primary" />
              <span className="text-sm font-medium">Fixed Contribution Amount (Group Contribution)</span>
            </label>
          </div>
          
          {hasFixedAmount && (
            <div className="ml-7 space-y-2">
              <p className="text-xs text-muted-foreground">
                Set a fixed amount that every contributor must pay. Useful for group contributions.
              </p>
              <div className="relative">
                <span className="absolute left-4 top-1/2 -translate-y-1/2 text-muted-foreground">
                  ₦
                </span>
                <input
                  type="text"
                  inputMode="numeric"
                  value={fixedContributionAmount}
                  onChange={(e) => handleFixedAmountChange(e.target.value)}
                  placeholder="e.g., 5000"
                  className={cn(
                    "w-full pl-8 pr-4 py-2 bg-background border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary",
                    errors.fixedContributionAmount ? "border-destructive" : "border-input"
                  )}
                />
              </div>
              {errors.fixedContributionAmount && (
                <p className="text-xs text-destructive">{errors.fixedContributionAmount}</p>
              )}
            </div>
          )}
        </div>

        {/* Visibility */}
        <div className="space-y-3">
          <label className="text-sm font-medium">Visibility</label>
          <div className="grid grid-cols-2 gap-3">
            <button
              type="button"
              onClick={() => setIsPublic(true)}
              className={cn(
                "flex items-center gap-3 p-4 rounded-lg border transition-all",
                isPublic
                  ? "border-primary bg-primary/5"
                  : "border-border hover:border-primary/50"
              )}
            >
              <Globe className={cn("w-5 h-5", isPublic && "text-primary")} />
              <div className="text-left">
                <p className="font-medium text-sm">Public</p>
                <p className="text-xs text-muted-foreground">Anyone can contribute</p>
              </div>
            </button>
            <button
              type="button"
              onClick={() => setIsPublic(false)}
              className={cn(
                "flex items-center gap-3 p-4 rounded-lg border transition-all",
                !isPublic
                  ? "border-primary bg-primary/5"
                  : "border-border hover:border-primary/50"
              )}
            >
              <Lock className={cn("w-5 h-5", !isPublic && "text-primary")} />
              <div className="text-left">
                <p className="font-medium text-sm">Private</p>
                <p className="text-xs text-muted-foreground">Only you can see</p>
              </div>
            </button>
          </div>
        </div>

        {/* Bank Details */}
        <div className="space-y-4 pt-4 border-t border-border">
          <h3 className="font-medium">Withdrawal Bank Details (Optional)</h3>
          <p className="text-sm text-muted-foreground">
            Add bank details for when you want to withdraw funds
          </p>

          <div className="space-y-2">
            <label className="text-sm font-medium">Bank Name</label>
            <input
              type="text"
              value={depositBankName}
              onChange={(e) => setDepositBankName(e.target.value)}
              placeholder="e.g., GTBank"
              className="w-full px-4 py-2 bg-background border border-input rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Account Number</label>
            <input
              type="text"
              value={depositAccountNumber}
              onChange={(e) => setDepositAccountNumber(e.target.value)}
              placeholder="0123456789"
              maxLength={10}
              className="w-full px-4 py-2 bg-background border border-input rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Account Name</label>
            <input
              type="text"
              value={depositAccountName}
              onChange={(e) => setDepositAccountName(e.target.value)}
              placeholder="Account holder name"
              className="w-full px-4 py-2 bg-background border border-input rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
            />
          </div>
        </div>

        {/* Submit */}
        <div className="flex items-center gap-3 pt-4">
          <Button
            type="button"
            variant="outline"
            onClick={() => navigate(-1)}
            className="flex-1"
          >
            Cancel
          </Button>
          <Button type="submit" disabled={isSaving} className="flex-1 gap-2">
            {isSaving ? (
              <Loader2 className="w-4 h-4 animate-spin" />
            ) : (
              <Save className="w-4 h-4" />
            )}
            Save Changes
          </Button>
        </div>
      </motion.form>
    </div>
  )
}
