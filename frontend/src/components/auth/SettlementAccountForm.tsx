import { useState } from "react"
import { motion } from "framer-motion"
import { Wallet, Loader2, CheckCircle2, AlertCircle } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { authApi } from "@/lib/api"
import { useToast } from "@/hooks/use-toast"

interface SettlementAccountFormProps {
  onComplete?: () => void
  onSkip?: () => void
  canSkip?: boolean
}

export function SettlementAccountForm({ onComplete, onSkip, canSkip = true }: SettlementAccountFormProps) {
  const [isLoading, setIsLoading] = useState(false)
  const [isComplete, setIsComplete] = useState(false)
  const [formData, setFormData] = useState({
    bankName: "",
    accountNumber: "",
    accountName: "",
  })
  const { toast } = useToast()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)

    try {
      await authApi.updateSettlementAccount({
        bank_name: formData.bankName,
        account_number: formData.accountNumber,
        account_name: formData.accountName,
      })
      
      setIsComplete(true)

      toast({
        title: "Settlement Account Added!",
        description: "Your bank account has been successfully linked",
      })

      // Call onComplete after a short delay to show success state
      setTimeout(() => {
        onComplete?.()
      }, 1500)
    } catch (error: any) {
      const errorMessage = error.response?.data?.error || "Failed to add settlement account"
      toast({
        title: "Failed",
        description: errorMessage,
        variant: "destructive",
      })
    } finally {
      setIsLoading(false)
    }
  }

  if (isComplete) {
    return (
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        className="text-center py-8"
      >
        <div className="w-16 h-16 bg-green-500/10 rounded-full flex items-center justify-center mx-auto mb-4">
          <CheckCircle2 className="w-8 h-8 text-green-500" />
        </div>
        <h3 className="text-xl font-semibold mb-2">Account Linked!</h3>
        <p className="text-muted-foreground">Your settlement account has been added</p>
      </motion.div>
    )
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="space-y-6"
    >
      <div className="text-center">
        <div className="w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-4">
          <Wallet className="w-6 h-6 text-primary" />
        </div>
        <h3 className="text-xl font-semibold mb-2">Add Settlement Account</h3>
        <p className="text-sm text-muted-foreground">
          Link your bank account to receive withdrawals
        </p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="bankName" className="text-sm">Bank Name</Label>
          <Input
            id="bankName"
            type="text"
            placeholder="e.g., GTBank, Access Bank"
            value={formData.bankName}
            onChange={(e) => setFormData({ ...formData, bankName: e.target.value })}
            required
            className="bg-muted/30 border-border/50 focus:border-primary/50"
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="accountNumber" className="text-sm">Account Number</Label>
          <Input
            id="accountNumber"
            type="text"
            placeholder="0123456789"
            value={formData.accountNumber}
            onChange={(e) => setFormData({ ...formData, accountNumber: e.target.value.replace(/\D/g, "").slice(0, 10) })}
            maxLength={10}
            required
            className="bg-muted/30 border-border/50 focus:border-primary/50"
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="accountName" className="text-sm">Account Name</Label>
          <Input
            id="accountName"
            type="text"
            placeholder="John Doe"
            value={formData.accountName}
            onChange={(e) => setFormData({ ...formData, accountName: e.target.value })}
            required
            className="bg-muted/30 border-border/50 focus:border-primary/50"
          />
        </div>

        <div className="bg-muted/30 rounded-lg p-4 border border-border/50">
          <div className="flex items-start gap-2">
            <AlertCircle className="w-4 h-4 text-primary mt-0.5 flex-shrink-0" />
            <p className="text-xs text-muted-foreground">
              This account will be used for withdrawing funds from your goals. You can update this later.
            </p>
          </div>
        </div>

        <div className="flex gap-3">
          <Button
            type="submit"
            className="flex-1 gap-2"
            disabled={isLoading}
          >
            {isLoading ? (
              <>
                <Loader2 className="w-4 h-4 animate-spin" />
                Saving...
              </>
            ) : (
              "Save Account"
            )}
          </Button>
          
          {canSkip && onSkip && (
            <Button
              type="button"
              variant="outline"
              onClick={onSkip}
              disabled={isLoading}
            >
              Skip for now
            </Button>
          )}
        </div>
      </form>
    </motion.div>
  )
}
