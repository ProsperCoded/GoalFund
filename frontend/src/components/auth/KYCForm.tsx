import { useState } from "react"
import { motion } from "framer-motion"
import { Shield, Loader2, CheckCircle2, AlertCircle } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { authApi } from "@/lib/api"
import { useToast } from "@/hooks/use-toast"
import { useAuth } from "@/contexts"

interface KYCFormProps {
  onComplete?: () => void
  onSkip?: () => void
  canSkip?: boolean
}

export function KYCForm({ onComplete, onSkip, canSkip = true }: KYCFormProps) {
  const [isLoading, setIsLoading] = useState(false)
  const [isVerified, setIsVerified] = useState(false)
  const [nin, setNin] = useState("")
  const { toast } = useToast()
  const { user, updateUser } = useAuth()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (nin.length !== 11) {
      toast({
        title: "Invalid NIN",
        description: "NIN must be exactly 11 digits",
        variant: "destructive",
      })
      return
    }

    setIsLoading(true)

    try {
      const response = await authApi.submitNIN({ nin })
      
      setIsVerified(true)
      
      // Update user with KYC status
      if (user) {
        updateUser({
          ...user,
          kyc_verified: response.kyc_verified,
          kyc_verified_at: response.kyc_verified_at,
        })
      }

      toast({
        title: "KYC Verified!",
        description: "Your identity has been successfully verified",
      })

      // Call onComplete after a short delay to show success state
      setTimeout(() => {
        onComplete?.()
      }, 1500)
    } catch (error: any) {
      const errorMessage = error.response?.data?.error || "KYC verification failed"
      toast({
        title: "Verification Failed",
        description: errorMessage,
        variant: "destructive",
      })
    } finally {
      setIsLoading(false)
    }
  }

  if (isVerified) {
    return (
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        className="text-center py-8"
      >
        <div className="w-16 h-16 bg-green-500/10 rounded-full flex items-center justify-center mx-auto mb-4">
          <CheckCircle2 className="w-8 h-8 text-green-500" />
        </div>
        <h3 className="text-xl font-semibold mb-2">KYC Verified!</h3>
        <p className="text-muted-foreground">Your identity has been successfully verified</p>
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
          <Shield className="w-6 h-6 text-primary" />
        </div>
        <h3 className="text-xl font-semibold mb-2">Verify Your Identity</h3>
        <p className="text-sm text-muted-foreground">
          Complete KYC verification to unlock full platform features
        </p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="nin" className="text-sm">National Identification Number (NIN)</Label>
          <Input
            id="nin"
            type="text"
            placeholder="12345678901"
            value={nin}
            onChange={(e) => setNin(e.target.value.replace(/\D/g, "").slice(0, 11))}
            maxLength={11}
            required
            className="bg-muted/30 border-border/50 focus:border-primary/50"
          />
          <p className="text-xs text-muted-foreground">
            Your NIN will be used to verify your identity. We keep this information secure.
          </p>
        </div>

        <div className="bg-muted/30 rounded-lg p-4 border border-border/50">
          <div className="flex items-start gap-2">
            <AlertCircle className="w-4 h-4 text-primary mt-0.5 flex-shrink-0" />
            <p className="text-xs text-muted-foreground">
              KYC verification is optional but recommended for goal creators who want to withdraw funds.
            </p>
          </div>
        </div>

        <div className="flex gap-3">
          <Button
            type="submit"
            className="flex-1 gap-2"
            disabled={isLoading || nin.length !== 11}
          >
            {isLoading ? (
              <>
                <Loader2 className="w-4 h-4 animate-spin" />
                Verifying...
              </>
            ) : (
              "Verify Identity"
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
