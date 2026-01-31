import { useState } from "react"
import { useNavigate } from "react-router-dom"
import { motion, AnimatePresence } from "framer-motion"
import { X } from "lucide-react"
import { Button } from "@/components/ui/button"
import { KYCForm } from "@/components/auth/KYCForm"
import { SettlementAccountForm } from "@/components/auth/SettlementAccountForm"
import { GradientText } from "@/components/animations"

type OnboardingStep = "welcome" | "kyc" | "settlement" | "complete"

export function OnboardingPage() {
  const navigate = useNavigate()
  const [currentStep, setCurrentStep] = useState<OnboardingStep>("welcome")

  const handleSkipAll = () => {
    navigate("/")
  }

  const handleKYCComplete = () => {
    setCurrentStep("settlement")
  }

  const handleKYCSkip = () => {
    setCurrentStep("settlement")
  }

  const handleSettlementComplete = () => {
    setCurrentStep("complete")
  }

  const handleSettlementSkip = () => {
    navigate("/")
  }

  const handleFinish = () => {
    navigate("/")
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4 bg-background">
      <div className="w-full max-w-lg">
        <AnimatePresence mode="wait">
          {currentStep === "welcome" && (
            <motion.div
              key="welcome"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              className="bg-card border border-border rounded-lg p-8 relative"
            >
              <Button
                variant="ghost"
                size="icon"
                className="absolute top-4 right-4"
                onClick={handleSkipAll}
              >
                <X className="w-4 h-4" />
              </Button>

              <div className="text-center space-y-6">
                <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center mx-auto">
                  <span className="text-2xl">🎉</span>
                </div>

                <div>
                  <h1 className="text-2xl font-bold mb-2">
                    Welcome to <GradientText text="GoFund!" />
                  </h1>
                  <p className="text-muted-foreground">
                    Let's set up your account to unlock all features
                  </p>
                </div>

                <div className="space-y-3 text-left bg-muted/30 rounded-lg p-4 border border-border/50">
                  <h3 className="font-semibold text-sm">What's next?</h3>
                  <ul className="space-y-2 text-sm text-muted-foreground">
                    <li className="flex items-start gap-2">
                      <span className="text-primary">•</span>
                      <span>Verify your identity (KYC) - Optional</span>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary">•</span>
                      <span>Add settlement account for withdrawals - Optional</span>
                    </li>
                  </ul>
                </div>

                <div className="flex flex-col gap-3">
                  <Button
                    onClick={() => setCurrentStep("kyc")}
                    className="w-full"
                  >
                    Get Started
                  </Button>
                  <Button
                    variant="outline"
                    onClick={handleSkipAll}
                  >
                    Skip for now
                  </Button>
                </div>
              </div>
            </motion.div>
          )}

          {currentStep === "kyc" && (
            <motion.div
              key="kyc"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -20 }}
              className="bg-card border border-border rounded-lg p-8 relative"
            >
              <Button
                variant="ghost"
                size="icon"
                className="absolute top-4 right-4"
                onClick={handleSkipAll}
              >
                <X className="w-4 h-4" />
              </Button>

              <div className="mb-4">
                <div className="flex gap-2 mb-6">
                  <div className="h-1 flex-1 bg-primary rounded" />
                  <div className="h-1 flex-1 bg-muted rounded" />
                </div>
              </div>

              <KYCForm
                onComplete={handleKYCComplete}
                onSkip={handleKYCSkip}
                canSkip={true}
              />
            </motion.div>
          )}

          {currentStep === "settlement" && (
            <motion.div
              key="settlement"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -20 }}
              className="bg-card border border-border rounded-lg p-8 relative"
            >
              <Button
                variant="ghost"
                size="icon"
                className="absolute top-4 right-4"
                onClick={handleSkipAll}
              >
                <X className="w-4 h-4" />
              </Button>

              <div className="mb-4">
                <div className="flex gap-2 mb-6">
                  <div className="h-1 flex-1 bg-primary rounded" />
                  <div className="h-1 flex-1 bg-primary rounded" />
                </div>
              </div>

              <SettlementAccountForm
                onComplete={handleSettlementComplete}
                onSkip={handleSettlementSkip}
                canSkip={true}
              />
            </motion.div>
          )}

          {currentStep === "complete" && (
            <motion.div
              key="complete"
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.95 }}
              className="bg-card border border-border rounded-lg p-8 text-center"
            >
              <div className="space-y-6">
                <div className="w-16 h-16 bg-green-500/10 rounded-full flex items-center justify-center mx-auto">
                  <span className="text-2xl">✓</span>
                </div>

                <div>
                  <h2 className="text-2xl font-bold mb-2">All Set!</h2>
                  <p className="text-muted-foreground">
                    Your account is ready. Start creating goals or contributing to existing ones.
                  </p>
                </div>

                <Button onClick={handleFinish} className="w-full">
                  Go to Dashboard
                </Button>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </div>
  )
}
