import { useState, useEffect } from "react"
import { Link } from "react-router-dom"
import { motion } from "framer-motion"
import {
  Search,
  Compass,
  Loader2,
  SlidersHorizontal,
  TrendingUp,
  Clock,
  Target,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { goalsApi, type Goal } from "@/lib/api"
import { cn, formatCurrency } from "@/lib/utils"

type SortOption = "trending" | "recent" | "ending_soon" | "most_funded"

export function ExplorePage() {
  const [goals, setGoals] = useState<Goal[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState("")
  const [sortBy, setSortBy] = useState<SortOption>("trending")

  useEffect(() => {
    fetchPublicGoals()
  }, [])

  const fetchPublicGoals = async () => {
    setIsLoading(true)
    try {
      const response = await goalsApi.getPublicGoals()
      setGoals(response.goals || [])
    } catch (error) {
      console.error("Failed to fetch public goals:", error)
      // Use mock data
      setGoals(getMockPublicGoals())
    } finally {
      setIsLoading(false)
    }
  }

  const filteredGoals = goals.filter((goal) =>
    goal.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
    goal.description.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const sortedGoals = [...filteredGoals].sort((a, b) => {
    switch (sortBy) {
      case "recent":
        return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      case "ending_soon":
        if (!a.deadline) return 1
        if (!b.deadline) return -1
        return new Date(a.deadline).getTime() - new Date(b.deadline).getTime()
      case "most_funded":
        return b.current_amount - a.current_amount
      case "trending":
      default:
        return b.contributor_count - a.contributor_count
    }
  })

  const sortOptions: { value: SortOption; label: string; icon: typeof TrendingUp }[] = [
    { value: "trending", label: "Trending", icon: TrendingUp },
    { value: "recent", label: "Most Recent", icon: Clock },
    { value: "ending_soon", label: "Ending Soon", icon: Clock },
    { value: "most_funded", label: "Most Funded", icon: Target },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold">Explore Goals</h1>
        <p className="text-muted-foreground mt-1">
          Discover and support public goals from the community
        </p>
      </div>

      {/* Search and Filters */}
      <div className="flex flex-col lg:flex-row gap-4">
        {/* Search */}
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <input
            type="text"
            placeholder="Search public goals..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-10 pr-4 py-2.5 bg-background border border-input rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
          />
        </div>

        {/* Sort Options */}
        <div className="flex items-center gap-2 overflow-x-auto pb-2 lg:pb-0">
          <SlidersHorizontal className="w-4 h-4 text-muted-foreground flex-shrink-0" />
          {sortOptions.map((option) => (
            <Button
              key={option.value}
              variant={sortBy === option.value ? "default" : "outline"}
              size="sm"
              onClick={() => setSortBy(option.value)}
              className="gap-1.5 flex-shrink-0"
            >
              <option.icon className="w-3.5 h-3.5" />
              {option.label}
            </Button>
          ))}
        </div>
      </div>

      {/* Goals Grid */}
      {isLoading ? (
        <div className="flex items-center justify-center min-h-[400px]">
          <Loader2 className="w-8 h-8 animate-spin text-primary" />
        </div>
      ) : sortedGoals.length === 0 ? (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="flex flex-col items-center justify-center min-h-[400px] text-center"
        >
          <Compass className="w-12 h-12 text-muted-foreground mb-4" />
          <h3 className="font-semibold text-lg mb-2">
            {searchQuery ? "No goals found" : "No public goals yet"}
          </h3>
          <p className="text-muted-foreground mb-4 max-w-md">
            {searchQuery
              ? "Try adjusting your search terms"
              : "Be the first to create a public goal and share it with the community!"}
          </p>
          {!searchQuery && (
            <Button asChild>
              <Link to="/dashboard/goals/create">Create a public goal</Link>
            </Button>
          )}
        </motion.div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {sortedGoals.map((goal, index) => (
            <PublicGoalCard key={goal.id} goal={goal} index={index} />
          ))}
        </div>
      )}
    </div>
  )
}

function PublicGoalCard({ goal, index }: { goal: Goal; index: number }) {
  const progress = Math.min((goal.current_amount / goal.target_amount) * 100, 100)
  const isOverfunded = goal.current_amount > goal.target_amount
  const daysLeft = goal.deadline
    ? Math.max(0, Math.ceil((new Date(goal.deadline).getTime() - Date.now()) / (1000 * 60 * 60 * 24)))
    : null

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.05 }}
      className="group bg-card border border-border rounded-lg overflow-hidden hover:border-primary/30 transition-all"
    >
      {/* Placeholder image area */}
      <div className="h-40 bg-gradient-to-br from-primary/10 to-primary/5 flex items-center justify-center">
        <Target className="w-12 h-12 text-primary/30" />
      </div>

      <div className="p-4">
        {/* Title and description */}
        <Link to={`/dashboard/goals/${goal.id}`}>
          <h3 className="font-semibold text-lg group-hover:text-primary transition-colors line-clamp-1">
            {goal.title}
          </h3>
        </Link>
        <p className="text-sm text-muted-foreground mt-1 line-clamp-2">
          {goal.description}
        </p>

        {/* Progress */}
        <div className="mt-4">
          <div className="flex items-center justify-between text-sm mb-1.5">
            <span className="font-medium">
              {formatCurrency(goal.current_amount)}
            </span>
            <span className="text-muted-foreground">
              {progress.toFixed(0)}%
            </span>
          </div>
          <div className="h-2 bg-muted rounded-full overflow-hidden">
            <div
              className={cn(
                "h-full rounded-full transition-all",
                isOverfunded ? "bg-green-500" : "bg-primary"
              )}
              style={{ width: `${Math.min(progress, 100)}%` }}
            />
          </div>
          <div className="flex items-center justify-between text-xs text-muted-foreground mt-1.5">
            <span>Goal: {formatCurrency(goal.target_amount)}</span>
            {daysLeft !== null && (
              <span className={daysLeft <= 7 ? "text-orange-500" : ""}>
                {daysLeft === 0 ? "Ends today" : `${daysLeft} days left`}
              </span>
            )}
          </div>
        </div>

        {/* Footer */}
        <div className="flex items-center justify-between mt-4 pt-4 border-t border-border">
          <div className="flex items-center gap-1 text-sm text-muted-foreground">
            <span className="font-medium text-foreground">{goal.contributor_count}</span>
            contributors
          </div>
          <Button size="sm" asChild>
            <Link to={`/dashboard/goals/${goal.id}`}>Contribute</Link>
          </Button>
        </div>
      </div>
    </motion.div>
  )
}

// Mock public goals
function getMockPublicGoals(): Goal[] {
  return [
    {
      id: "p1",
      user_id: "user2",
      title: "Build a Library for Rural Kids",
      description: "Help us build a community library to provide educational resources for children in rural areas who lack access to books and learning materials.",
      target_amount: 8000000,
      current_amount: 5600000,
      currency: "NGN",
      status: "open",
      is_public: true,
      contributor_count: 127,
      deadline: new Date(Date.now() + 45 * 24 * 60 * 60 * 1000).toISOString(),
      created_at: new Date(Date.now() - 15 * 24 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date().toISOString(),
    },
    {
      id: "p2",
      user_id: "user3",
      title: "Support Local Farmers Initiative",
      description: "Providing agricultural equipment and training to local farmers to improve crop yields and food security in our community.",
      target_amount: 3500000,
      current_amount: 2100000,
      currency: "NGN",
      status: "open",
      is_public: true,
      contributor_count: 89,
      deadline: new Date(Date.now() + 21 * 24 * 60 * 60 * 1000).toISOString(),
      created_at: new Date(Date.now() - 10 * 24 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date().toISOString(),
    },
    {
      id: "p3",
      user_id: "user4",
      title: "Medical Fund for Mama Ada",
      description: "Raising funds for Mama Ada's heart surgery. She has been a pillar of our community for over 40 years and now needs our help.",
      target_amount: 4500000,
      current_amount: 4200000,
      currency: "NGN",
      status: "open",
      is_public: true,
      contributor_count: 203,
      deadline: new Date(Date.now() + 5 * 24 * 60 * 60 * 1000).toISOString(),
      created_at: new Date(Date.now() - 25 * 24 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date().toISOString(),
    },
    {
      id: "p4",
      user_id: "user5",
      title: "Clean Water for Village Schools",
      description: "Installing water purification systems in 5 village schools to provide clean drinking water for over 2000 students.",
      target_amount: 6000000,
      current_amount: 1500000,
      currency: "NGN",
      status: "open",
      is_public: true,
      contributor_count: 42,
      deadline: new Date(Date.now() + 60 * 24 * 60 * 60 * 1000).toISOString(),
      created_at: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date().toISOString(),
    },
    {
      id: "p5",
      user_id: "user6",
      title: "Youth Skills Training Center",
      description: "Building a skills acquisition center to train young people in vocational skills like tailoring, carpentry, and digital skills.",
      target_amount: 10000000,
      current_amount: 3500000,
      currency: "NGN",
      status: "open",
      is_public: true,
      contributor_count: 78,
      created_at: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date().toISOString(),
    },
    {
      id: "p6",
      user_id: "user7",
      title: "Solar Power for Health Clinic",
      description: "Installing solar panels at the community health clinic to ensure 24/7 power supply for medical equipment and lighting.",
      target_amount: 2500000,
      current_amount: 2600000,
      currency: "NGN",
      status: "open",
      is_public: true,
      contributor_count: 156,
      deadline: new Date(Date.now() + 10 * 24 * 60 * 60 * 1000).toISOString(),
      created_at: new Date(Date.now() - 20 * 24 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date().toISOString(),
    },
  ]
}
