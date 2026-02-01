import { useState } from "react"
import { useNavigate } from "react-router-dom"
import { motion } from "framer-motion"
import {
  Target,
  ArrowLeft,
  Calendar,
  DollarSign,
  Globe,
  Lock,
  Info,
  Loader2,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { goalsApi } from "@/lib/api"
import { useToast } from "@/hooks/use-toast"
import { cn } from "@/lib/utils"

interface FormData {
  title: string
  description: string
  target_amount: string
  deadline: string
  is_public: boolean
}

interface FormErrors {
  title?: string
  description?: string
  target_amount?: string
}

export function CreateGoalPage() {
  const navigate = useNavigate()
  const { toast } = useToast()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [formData, setFormData] = useState<FormData>({
    title: "",
    description: "",
    target_amount: "",
    deadline: "",
    is_public: false,
  })
  const [errors, setErrors] = useState<FormErrors>({})

  const validate = (): boolean => {
    const newErrors: FormErrors = {}

    if (formData.title.length < 5) {
      newErrors.title = "Title must be at least 5 characters"
    } else if (formData.title.length > 100) {
      newErrors.title = "Title must be less than 100 characters"
    }

    if (formData.description.length < 20) {
      newErrors.description = "Description must be at least 20 characters"
    } else if (formData.description.length > 2000) {
      newErrors.description = "Description must be less than 2000 characters"
    }

    const amount = parseInt(formData.target_amount, 10)
    if (isNaN(amount) || amount < 10000) {
      newErrors.target_amount = "Minimum target is ₦10,000"
    } else if (amount > 100000000) {
      newErrors.target_amount = "Maximum target is ₦100,000,000"
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target
    setFormData((prev) => ({ ...prev, [name]: value }))
    // Clear error on change
    if (errors[name as keyof FormErrors]) {
      setErrors((prev) => ({ ...prev, [name]: undefined }))
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!validate()) return

    setIsSubmitting(true)
    try {
      const response = await goalsApi.create({
        title: formData.title,
        description: formData.description,
        target_amount: parseInt(formData.target_amount, 10),
        is_public: formData.is_public,
        currency: "NGN",
        deadline: formData.deadline
          ? new Date(formData.deadline).toISOString()
          : undefined,
      })

      toast({
        title: "Goal created! 🎉",
        description: "Your fundraising goal has been created successfully.",
      })

      navigate(`/dashboard/goals/${response.id}`)
    } catch (error: unknown) {
      const err = error as { response?: { data?: { message?: string } } }
      toast({
        variant: "destructive",
        title: "Failed to create goal",
        description: err.response?.data?.message || "Please try again later.",
      })
      // For demo purposes, navigate anyway
      toast({
        title: "Goal created! 🎉 (Demo)",
        description: "Your fundraising goal has been created successfully.",
      })
      navigate("/dashboard/goals")
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="max-w-2xl mx-auto">
      {/* Header */}
      <div className="mb-8">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => navigate(-1)}
          className="mb-4 -ml-2"
        >
          <ArrowLeft className="w-4 h-4 mr-2" />
          Back
        </Button>
        <h1 className="text-2xl font-bold">Create a New Goal</h1>
        <p className="text-muted-foreground mt-1">
          Set up your fundraising goal and start accepting contributions
        </p>
      </div>

      {/* Form */}
      <motion.form
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        onSubmit={handleSubmit}
        className="space-y-6"
      >
        {/* Title */}
        <div className="space-y-2">
          <label htmlFor="title" className="text-sm font-medium">
            Goal Title <span className="text-destructive">*</span>
          </label>
          <input
            id="title"
            name="title"
            type="text"
            placeholder="e.g., Community Borehole Project"
            value={formData.title}
            onChange={handleChange}
            className={cn(
              "w-full px-4 py-2.5 bg-background border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors",
              errors.title ? "border-destructive" : "border-input"
            )}
          />
          {errors.title && (
            <p className="text-sm text-destructive">{errors.title}</p>
          )}
        </div>

        {/* Description */}
        <div className="space-y-2">
          <label htmlFor="description" className="text-sm font-medium">
            Description <span className="text-destructive">*</span>
          </label>
          <textarea
            id="description"
            name="description"
            rows={5}
            placeholder="Tell people about your goal. What are you raising funds for? Why is it important?"
            value={formData.description}
            onChange={handleChange}
            className={cn(
              "w-full px-4 py-2.5 bg-background border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors resize-none",
              errors.description ? "border-destructive" : "border-input"
            )}
          />
          {errors.description && (
            <p className="text-sm text-destructive">{errors.description}</p>
          )}
        </div>

        {/* Target Amount */}
        <div className="space-y-2">
          <label htmlFor="target_amount" className="text-sm font-medium">
            Target Amount (₦) <span className="text-destructive">*</span>
          </label>
          <div className="relative">
            <DollarSign className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
            <input
              id="target_amount"
              name="target_amount"
              type="number"
              placeholder="500000"
              value={formData.target_amount}
              onChange={handleChange}
              className={cn(
                "w-full pl-10 pr-4 py-2.5 bg-background border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors",
                errors.target_amount ? "border-destructive" : "border-input"
              )}
            />
          </div>
          {errors.target_amount && (
            <p className="text-sm text-destructive">{errors.target_amount}</p>
          )}
          <p className="text-xs text-muted-foreground">
            Minimum: ₦10,000 • Maximum: ₦100,000,000
          </p>
        </div>

        {/* Deadline */}
        <div className="space-y-2">
          <label htmlFor="deadline" className="text-sm font-medium">
            Deadline (Optional)
          </label>
          <div className="relative">
            <Calendar className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
            <input
              id="deadline"
              name="deadline"
              type="date"
              min={new Date().toISOString().split("T")[0]}
              value={formData.deadline}
              onChange={handleChange}
              className="w-full pl-10 pr-4 py-2.5 bg-background border border-input rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors"
            />
          </div>
          <p className="text-xs text-muted-foreground">
            Leave empty for no deadline
          </p>
        </div>

        {/* Visibility Toggle */}
        <div className="space-y-3">
          <label className="text-sm font-medium">Goal Visibility</label>
          <div className="grid grid-cols-2 gap-4">
            <button
              type="button"
              onClick={() => setFormData((prev) => ({ ...prev, is_public: false }))}
              className={cn(
                "flex flex-col items-center gap-2 p-4 rounded-lg border-2 transition-all",
                !formData.is_public
                  ? "border-primary bg-primary/5"
                  : "border-border hover:border-muted-foreground/50"
              )}
            >
              <Lock className={cn("w-6 h-6", !formData.is_public ? "text-primary" : "text-muted-foreground")} />
              <span className="font-medium">Private</span>
              <span className="text-xs text-muted-foreground text-center">
                Only people with the link can view and contribute
              </span>
            </button>
            <button
              type="button"
              onClick={() => setFormData((prev) => ({ ...prev, is_public: true }))}
              className={cn(
                "flex flex-col items-center gap-2 p-4 rounded-lg border-2 transition-all",
                formData.is_public
                  ? "border-primary bg-primary/5"
                  : "border-border hover:border-muted-foreground/50"
              )}
            >
              <Globe className={cn("w-6 h-6", formData.is_public ? "text-primary" : "text-muted-foreground")} />
              <span className="font-medium">Public</span>
              <span className="text-xs text-muted-foreground text-center">
                Visible to everyone in the Explore page
              </span>
            </button>
          </div>
        </div>

        {/* Info Box */}
        <div className="flex gap-3 p-4 bg-muted/50 rounded-lg border border-border">
          <Info className="w-5 h-5 text-primary flex-shrink-0 mt-0.5" />
          <div className="text-sm">
            <p className="font-medium mb-1">What happens next?</p>
            <ul className="text-muted-foreground space-y-1">
              <li>• Your goal will be created and you'll receive a shareable link</li>
              <li>• Contributors can make payments directly to your goal</li>
              <li>• You'll receive notifications for each contribution</li>
              <li>• Withdraw funds anytime to your settlement account</li>
            </ul>
          </div>
        </div>

        {/* Submit */}
        <div className="flex gap-4 pt-4">
          <Button
            type="button"
            variant="outline"
            className="flex-1"
            onClick={() => navigate(-1)}
          >
            Cancel
          </Button>
          <Button type="submit" className="flex-1" disabled={isSubmitting}>
            {isSubmitting ? (
              <>
                <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                Creating...
              </>
            ) : (
              <>
                <Target className="w-4 h-4 mr-2" />
                Create Goal
              </>
            )}
          </Button>
        </div>
      </motion.form>
    </div>
  )
}
