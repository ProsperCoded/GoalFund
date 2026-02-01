import { useState, useEffect } from "react"
import { Link } from "react-router-dom"
import { motion } from "framer-motion"
import { Plus, Search, Target, Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { GoalCard } from "@/components/dashboard"
import { goalsApi, type Goal } from "@/lib/api"

type FilterStatus = "all" | "open" | "closed"

export function MyGoalsPage() {
  const [goals, setGoals] = useState<Goal[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState("")
  const [statusFilter, setStatusFilter] = useState<FilterStatus>("all")

  useEffect(() => {
    fetchGoals()
  }, [])

  const fetchGoals = async () => {
    setIsLoading(true)
    try {
      const response = await goalsApi.getMyGoals()
      setGoals(response.goals || [])
    } catch (error) {
      console.error("Failed to fetch goals:", error)
      // Use mock data for demo
      setGoals(getMockGoals())
    } finally {
      setIsLoading(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm("Are you sure you want to delete this goal?")) return
    try {
      // await goalsApi.deleteGoal(id)
      setGoals((prev) => prev.filter((g) => g.id !== id))
    } catch (error) {
      console.error("Failed to delete goal:", error)
    }
  }

  const filteredGoals = goals.filter((goal) => {
    const matchesSearch = goal.title.toLowerCase().includes(searchQuery.toLowerCase())
    const matchesStatus = statusFilter === "all" || goal.status === statusFilter
    return matchesSearch && matchesStatus
  })

  const openGoals = goals.filter((g) => g.status === "open").length
  const closedGoals = goals.filter((g) => g.status === "closed").length

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold">My Goals</h1>
          <p className="text-muted-foreground mt-1">
            Manage and track all your fundraising goals
          </p>
        </div>
        <Button asChild>
          <Link to="/dashboard/goals/create" className="gap-2">
            <Plus className="w-4 h-4" />
            Create New Goal
          </Link>
        </Button>
      </div>

      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-4">
        {/* Search */}
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <input
            type="text"
            placeholder="Search goals..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-10 pr-4 py-2 bg-background border border-input rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
          />
        </div>

        {/* Status Filter */}
        <div className="flex items-center gap-2">
          {(["all", "open", "closed"] as FilterStatus[]).map((status) => (
            <Button
              key={status}
              variant={statusFilter === status ? "default" : "outline"}
              size="sm"
              onClick={() => setStatusFilter(status)}
              className="capitalize"
            >
              {status}
              {status === "all" && ` (${goals.length})`}
              {status === "open" && ` (${openGoals})`}
              {status === "closed" && ` (${closedGoals})`}
            </Button>
          ))}
        </div>
      </div>

      {/* Goals Grid */}
      {isLoading ? (
        <div className="flex items-center justify-center min-h-[300px]">
          <Loader2 className="w-8 h-8 animate-spin text-primary" />
        </div>
      ) : filteredGoals.length === 0 ? (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="flex flex-col items-center justify-center min-h-[300px] text-center"
        >
          <Target className="w-12 h-12 text-muted-foreground mb-4" />
          <h3 className="font-semibold text-lg mb-2">
            {searchQuery || statusFilter !== "all"
              ? "No goals found"
              : "No goals yet"}
          </h3>
          <p className="text-muted-foreground mb-4 max-w-md">
            {searchQuery || statusFilter !== "all"
              ? "Try adjusting your search or filters"
              : "Create your first fundraising goal and start accepting contributions!"}
          </p>
          {!searchQuery && statusFilter === "all" && (
            <Button asChild>
              <Link to="/dashboard/goals/create">Create your first goal</Link>
            </Button>
          )}
        </motion.div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {filteredGoals.map((goal, index) => (
            <motion.div
              key={goal.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.05 }}
            >
              <GoalCard goal={goal} showActions onDelete={handleDelete} />
            </motion.div>
          ))}
        </div>
      )}
    </div>
  )
}

// Mock data
function getMockGoals(): Goal[] {
  return [
    {
      id: "1",
      user_id: "user1",
      title: "Community Borehole Project",
      description: "Help us provide clean water to our community. This project will serve over 500 families in the area.",
      target_amount: 5000000,
      current_amount: 3250000,
      currency: "NGN",
      status: "open",
      is_public: true,
      contributor_count: 45,
      deadline: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString(),
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    },
    {
      id: "2",
      user_id: "user1",
      title: "School Fees Fund",
      description: "Raising funds for my final year university tuition. Your support means the world to me!",
      target_amount: 500000,
      current_amount: 500000,
      currency: "NGN",
      status: "open",
      is_public: false,
      contributor_count: 12,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    },
    {
      id: "3",
      user_id: "user1",
      title: "Medical Emergency Fund",
      description: "Help with medical expenses for a family member going through a difficult time.",
      target_amount: 2000000,
      current_amount: 750000,
      currency: "NGN",
      status: "open",
      is_public: true,
      contributor_count: 28,
      deadline: new Date(Date.now() + 14 * 24 * 60 * 60 * 1000).toISOString(),
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    },
    {
      id: "4",
      user_id: "user1",
      title: "Wedding Contribution",
      description: "Collecting contributions for our upcoming wedding celebration.",
      target_amount: 3000000,
      current_amount: 3500000,
      currency: "NGN",
      status: "closed",
      is_public: false,
      contributor_count: 67,
      created_at: new Date(Date.now() - 60 * 24 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date().toISOString(),
    },
  ]
}
