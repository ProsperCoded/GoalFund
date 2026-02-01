import { useState, useEffect } from "react"
import { motion, AnimatePresence } from "framer-motion"
import {
  X,
  Loader2,
  CreditCard,
  ExternalLink,
  CheckCircle,
  AlertCircle,
  Heart,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { paymentsApi, nairaToKobo, type Goal } from "@/lib/api"
import { useAuth } from "@/contexts"
import { useToast } from "@/hooks/use-toast"
import { cn, formatCurrency } from "@/lib/utils"

interface ContributeModalProps {
  isOpen: boolean
  onClose: () => void
  goal: Goal
  onSuccess?: () => void
}

type Step = "amount" | "processing" | "redirect" | "verifying" | "success" | "error"

export function ContributeModal({
  isOpen,
  onClose,
  goal,
  onSuccess,
}: ContributeModalProps) {
  const { user } = useAuth()
  const { toast } = useToast()
  const [step, setStep] = useState<Step>("amount")
  const [amount, setAmount] = useState("")
  const [message, setMessage] = useState("")
  const [isAnonymous, setIsAnonymous] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Predefined amounts
  const quickAmounts = [5000, 10000, 25000, 50000, 100000]

  // Reset state when modal opens/closes
  useEffect(() => {
    if (isOpen) {
      setStep("amount")
      setAmount("")
      setMessage("")
      setIsAnonymous(false)
      setError(null)
    }
  }, [isOpen])

  // Check for payment return from Paystack
  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search)
    const reference = urlParams.get("reference")
    const trxref = urlParams.get("trxref")

    if (reference || trxref) {
      const ref = reference || trxref
      setStep("verifying")
      verifyPayment(ref!)

      // Clean URL
      window.history.replaceState({}, "", window.location.pathname)
    }
  }, [])

  const handleAmountChange = (value: string) => {
    // Only allow numbers
    const numericValue = value.replace(/[^0-9]/g, "")
    setAmount(numericValue)
  }

  const handleQuickAmount = (value: number) => {
    setAmount(value.toString())
  }

  const initializePayment = async () => {
    if (!user) {
      toast({
        variant: "destructive",
        title: "Please log in",
        description: "You need to be logged in to contribute.",
      })
      return
    }

    const numericAmount = parseInt(amount, 10)
    if (isNaN(numericAmount) || numericAmount < 100) {
      toast({
        variant: "destructive",
        title: "Invalid amount",
        description: "Minimum contribution is ₦100",
      })
      return
    }

    setStep("processing")
    setError(null)

    try {
      const response = await paymentsApi.initialize({
        user_id: user.id,
        goal_id: goal.id,
        amount: nairaToKobo(numericAmount),
        currency: "NGN",
        email: user.email,
        callback_url: `${window.location.origin}/dashboard/goals/${goal.id}?payment=true`,
        metadata: {
          goal_title: goal.title,
          contributor_name: isAnonymous ? "Anonymous" : `${user.first_name} ${user.last_name}`,
          message: message || undefined,
          is_anonymous: isAnonymous,
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

  const verifyPayment = async (reference: string) => {
    try {
      const result = await paymentsApi.verify(reference)

      if (result.status === "VERIFIED") {
        setStep("success")
        toast({
          title: "Contribution successful! 🎉",
          description: `Thank you for contributing to "${goal.title}"`,
        })
        onSuccess?.()
      } else if (result.status === "FAILED") {
        setError("Payment was not successful. Please try again.")
        setStep("error")
      } else {
        // Still pending, might need to wait
        setError("Payment is still being processed. Please check back later.")
        setStep("error")
      }
    } catch (err: unknown) {
      const error = err as { response?: { data?: { message?: string } } }
      setError(error.response?.data?.message || "Failed to verify payment")
      setStep("error")
    }
  }

  const numericAmount = parseInt(amount, 10) || 0
  const isValidAmount = numericAmount >= 100

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 z-50 bg-black/50"
            onClick={step === "amount" ? onClose : undefined}
          />

          {/* Modal */}
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.95 }}
            className="fixed left-1/2 top-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2 bg-card border border-border rounded-lg shadow-xl overflow-hidden"
          >
            {/* Header */}
            <div className="flex items-center justify-between p-4 border-b border-border">
              <h2 className="font-semibold text-lg">
                {step === "success" ? "Thank You!" : "Contribute to Goal"}
              </h2>
              {step === "amount" && (
                <Button variant="ghost" size="icon" onClick={onClose}>
                  <X className="w-4 h-4" />
                </Button>
              )}
            </div>

            {/* Content */}
            <div className="p-6">
              {/* Step: Amount Selection */}
              {step === "amount" && (
                <div className="space-y-6">
                  {/* Goal Info */}
                  <div className="p-4 bg-muted/50 rounded-lg">
                    <h3 className="font-medium truncate">{goal.title}</h3>
                    <div className="flex items-center gap-2 mt-2 text-sm text-muted-foreground">
                      <span>{formatCurrency(goal.current_amount)} raised</span>
                      <span>•</span>
                      <span>Goal: {formatCurrency(goal.target_amount)}</span>
                    </div>
                  </div>

                  {/* Amount Input */}
                  <div className="space-y-3">
                    <label className="text-sm font-medium">Enter Amount</label>
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
                        className="w-full pl-8 pr-4 py-3 text-2xl font-semibold bg-background border border-input rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
                      />
                    </div>
                    {amount && !isValidAmount && (
                      <p className="text-sm text-destructive">
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
                          "px-3 py-1.5 text-sm rounded-full border transition-colors",
                          amount === value.toString()
                            ? "border-primary bg-primary/10 text-primary"
                            : "border-border hover:border-primary/50"
                        )}
                      >
                        ₦{value.toLocaleString()}
                      </button>
                    ))}
                  </div>

                  {/* Message (Optional) */}
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

                  {/* Anonymous Option */}
                  <label className="flex items-center gap-3 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={isAnonymous}
                      onChange={(e) => setIsAnonymous(e.target.checked)}
                      className="w-4 h-4 rounded border-input"
                    />
                    <span className="text-sm">Contribute anonymously</span>
                  </label>

                  {/* Submit Button */}
                  <Button
                    onClick={initializePayment}
                    disabled={!isValidAmount}
                    className="w-full gap-2"
                    size="lg"
                  >
                    <CreditCard className="w-4 h-4" />
                    Contribute {isValidAmount ? formatCurrency(numericAmount) : ""}
                  </Button>

                  <p className="text-xs text-center text-muted-foreground">
                    Secured by Paystack. You'll be redirected to complete payment.
                  </p>
                </div>
              )}

              {/* Step: Processing */}
              {step === "processing" && (
                <div className="flex flex-col items-center justify-center py-8 text-center">
                  <Loader2 className="w-12 h-12 text-primary animate-spin mb-4" />
                  <h3 className="font-semibold text-lg mb-2">
                    Initializing Payment
                  </h3>
                  <p className="text-sm text-muted-foreground">
                    Please wait while we set up your payment...
                  </p>
                </div>
              )}

              {/* Step: Redirect */}
              {step === "redirect" && (
                <div className="flex flex-col items-center justify-center py-8 text-center">
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
                <div className="flex flex-col items-center justify-center py-8 text-center">
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
                <div className="flex flex-col items-center justify-center py-8 text-center">
                  <div className="w-16 h-16 rounded-full bg-green-500/10 flex items-center justify-center mb-4">
                    <CheckCircle className="w-10 h-10 text-green-500" />
                  </div>
                  <h3 className="font-semibold text-lg mb-2">
                    Contribution Successful!
                  </h3>
                  <p className="text-sm text-muted-foreground mb-6">
                    Thank you for supporting "{goal.title}". Your contribution
                    makes a difference!
                  </p>
                  <div className="flex items-center gap-2 text-primary">
                    <Heart className="w-4 h-4 fill-current" />
                    <span className="font-medium">
                      {formatCurrency(numericAmount)} contributed
                    </span>
                  </div>
                  <Button onClick={onClose} className="mt-6">
                    Done
                  </Button>
                </div>
              )}

              {/* Step: Error */}
              {step === "error" && (
                <div className="flex flex-col items-center justify-center py-8 text-center">
                  <div className="w-16 h-16 rounded-full bg-destructive/10 flex items-center justify-center mb-4">
                    <AlertCircle className="w-10 h-10 text-destructive" />
                  </div>
                  <h3 className="font-semibold text-lg mb-2">
                    Payment Failed
                  </h3>
                  <p className="text-sm text-muted-foreground mb-6">
                    {error || "Something went wrong. Please try again."}
                  </p>
                  <div className="flex gap-3">
                    <Button variant="outline" onClick={onClose}>
                      Cancel
                    </Button>
                    <Button onClick={() => setStep("amount")}>Try Again</Button>
                  </div>
                </div>
              )}
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  )
}
