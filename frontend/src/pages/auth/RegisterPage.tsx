import { useState } from "react"
import { Link, useNavigate } from "react-router-dom"
import { motion } from "framer-motion"
import {
  Mail,
  Lock,
  Eye,
  EyeOff,
  ArrowRight,
  Github,
  Loader2,
  User,
  CheckCircle2,
  Target,
  Users,
  Banknote,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Checkbox } from "@/components/ui/checkbox"
import { Separator } from "@/components/ui/separator"
import {
  GradientText,
  FadeIn,
  BlurText,
  StaggerContainer,
  StaggerItem,
} from "@/components/animations"

export function RegisterPage() {
  const navigate = useNavigate()
  const [showPassword, setShowPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    password: "",
    acceptTerms: false,
  })

  const passwordRequirements = [
    { label: "At least 8 characters", met: formData.password.length >= 8 },
    { label: "One uppercase letter", met: /[A-Z]/.test(formData.password) },
    { label: "One lowercase letter", met: /[a-z]/.test(formData.password) },
    { label: "One number", met: /\d/.test(formData.password) },
  ]

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)

    // Simulate API call
    setTimeout(() => {
      setIsLoading(false)
      navigate("/dashboard")
    }, 1500)
  }

  const features = [
    { 
      icon: Target, 
      title: "Create Goals",
      description: "Set up funding goals with milestones and track progress" 
    },
    { 
      icon: Users, 
      title: "Collect Contributions",
      description: "Contributors can fund with just an email—no signup needed" 
    },
    { 
      icon: Banknote, 
      title: "Withdraw Anytime",
      description: "Access funds when you need them with full audit trail" 
    },
  ]

  return (
    <div className="min-h-screen flex">
      {/* Left Side - Visual */}
      <div className="hidden lg:flex flex-1 relative overflow-hidden bg-card border-r border-border/50">
        {/* Subtle grid */}
        <div className="absolute inset-0 opacity-[0.03]">
          <div className="absolute inset-0" style={{
            backgroundImage: `linear-gradient(rgba(212, 168, 83, 0.5) 1px, transparent 1px),
                             linear-gradient(90deg, rgba(212, 168, 83, 0.5) 1px, transparent 1px)`,
            backgroundSize: '40px 40px'
          }} />
        </div>
        
        {/* Corner accents */}
        <div className="absolute top-8 right-8 w-16 h-16 border-t border-r border-primary/20" />
        <div className="absolute bottom-8 left-8 w-16 h-16 border-b border-l border-primary/20" />
        
        <div className="relative z-10 flex items-center justify-center w-full p-12">
          <FadeIn direction="right">
            <div className="max-w-md">
              <div className="mb-10">
                <h2 className="text-2xl font-bold mb-4">
                  Start Funding <GradientText text="Big Goals" />
                </h2>
                <p className="text-muted-foreground text-sm leading-relaxed">
                  Create goals for community projects, group contributions, and large-scale 
                  funding with complete accountability.
                </p>
              </div>

              {/* Features */}
              <StaggerContainer className="space-y-4" staggerDelay={0.1}>
                {features.map((feature) => {
                  const Icon = feature.icon
                  return (
                    <StaggerItem key={feature.title}>
                      <div className="flex items-start gap-4 p-4 bg-muted/30 rounded-lg border border-border/50">
                        <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0">
                          <Icon className="w-5 h-5 text-primary" />
                        </div>
                        <div>
                          <h3 className="font-medium text-sm mb-1">{feature.title}</h3>
                          <p className="text-xs text-muted-foreground">{feature.description}</p>
                        </div>
                      </div>
                    </StaggerItem>
                  )
                })}
              </StaggerContainer>

              {/* Trust line */}
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ delay: 0.8 }}
                className="mt-10 pt-6 border-t border-border/50"
              >
                <p className="text-xs text-muted-foreground text-center">
                  Ledger-backed accounting • Verified payments • Complete transparency
                </p>
              </motion.div>
            </div>
          </FadeIn>
        </div>
      </div>

      {/* Right Side - Form */}
      <div className="flex-1 flex items-center justify-center p-8">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          className="w-full max-w-md"
        >
          {/* Logo */}
          <Link to="/" className="flex items-center gap-2.5 mb-10">
            <div className="w-8 h-8 rounded-lg bg-primary/10 border border-primary/20 flex items-center justify-center">
              <span className="text-primary font-bold text-sm">G</span>
            </div>
            <span className="text-xl font-bold">
              <GradientText text="GoFund" />
            </span>
          </Link>

          {/* Header */}
          <div className="mb-8">
            <h1 className="text-3xl font-bold mb-2">
              <BlurText text="Create your account" />
            </h1>
            <p className="text-muted-foreground">
              Get started with transparent group funding
            </p>
          </div>

          {/* Social Login */}
          <div className="space-y-3 mb-6">
            <Button variant="outline" className="w-full gap-2 border-border/50 hover:bg-muted/50" type="button">
              <svg className="w-5 h-5" viewBox="0 0 24 24">
                <path
                  fill="currentColor"
                  d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                />
                <path
                  fill="currentColor"
                  d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                />
                <path
                  fill="currentColor"
                  d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                />
                <path
                  fill="currentColor"
                  d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                />
              </svg>
              Continue with Google
            </Button>
            <Button variant="outline" className="w-full gap-2 border-border/50 hover:bg-muted/50" type="button">
              <Github className="w-5 h-5" />
              Continue with GitHub
            </Button>
          </div>

          <div className="relative mb-6">
            <Separator className="bg-border/50" />
            <span className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 bg-background px-2 text-xs text-muted-foreground uppercase tracking-wider">
              or
            </span>
          </div>

          {/* Register Form */}
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name" className="text-sm">Full Name</Label>
              <Input
                id="name"
                type="text"
                placeholder="John Doe"
                icon={<User className="w-4 h-4" />}
                value={formData.name}
                onChange={(e) =>
                  setFormData({ ...formData, name: e.target.value })
                }
                required
                className="bg-muted/30 border-border/50 focus:border-primary/50"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="email" className="text-sm">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="you@example.com"
                icon={<Mail className="w-4 h-4" />}
                value={formData.email}
                onChange={(e) =>
                  setFormData({ ...formData, email: e.target.value })
                }
                required
                className="bg-muted/30 border-border/50 focus:border-primary/50"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="password" className="text-sm">Password</Label>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  placeholder="••••••••"
                  icon={<Lock className="w-4 h-4" />}
                  value={formData.password}
                  onChange={(e) =>
                    setFormData({ ...formData, password: e.target.value })
                  }
                  required
                  className="bg-muted/30 border-border/50 focus:border-primary/50"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                >
                  {showPassword ? (
                    <EyeOff className="w-4 h-4" />
                  ) : (
                    <Eye className="w-4 h-4" />
                  )}
                </button>
              </div>

              {/* Password Requirements */}
              {formData.password && (
                <motion.div
                  initial={{ opacity: 0, height: 0 }}
                  animate={{ opacity: 1, height: "auto" }}
                  className="grid grid-cols-2 gap-2 mt-2"
                >
                  {passwordRequirements.map((req) => (
                    <div
                      key={req.label}
                      className={`flex items-center gap-1 text-xs ${
                        req.met ? "text-primary" : "text-muted-foreground"
                      }`}
                    >
                      <CheckCircle2
                        className={`w-3 h-3 ${
                          req.met ? "text-primary" : "text-muted-foreground"
                        }`}
                      />
                      {req.label}
                    </div>
                  ))}
                </motion.div>
              )}
            </div>

            <div className="flex items-start gap-2">
              <Checkbox
                id="terms"
                checked={formData.acceptTerms}
                onCheckedChange={(checked) =>
                  setFormData({ ...formData, acceptTerms: checked as boolean })
                }
                className="mt-0.5"
              />
              <Label htmlFor="terms" className="text-sm cursor-pointer leading-tight text-muted-foreground">
                I agree to the{" "}
                <Link to="/terms" className="text-primary hover:underline">
                  Terms of Service
                </Link>{" "}
                and{" "}
                <Link to="/privacy" className="text-primary hover:underline">
                  Privacy Policy
                </Link>
              </Label>
            </div>

            <Button
              type="submit"
              size="lg"
              className="w-full gap-2 bg-primary text-primary-foreground hover:bg-primary/90"
              disabled={isLoading || !formData.acceptTerms}
            >
              {isLoading ? (
                <>
                  <Loader2 className="w-4 h-4 animate-spin" />
                  Creating account...
                </>
              ) : (
                <>
                  Create Account
                  <ArrowRight className="w-4 h-4" />
                </>
              )}
            </Button>
          </form>

          {/* Sign In Link */}
          <p className="mt-6 text-center text-sm text-muted-foreground">
            Already have an account?{" "}
            <Link to="/login" className="text-primary hover:underline font-medium">
              Sign in
            </Link>
          </p>
        </motion.div>
      </div>
    </div>
  )
}
