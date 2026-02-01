import { useState, useEffect } from "react"
import { Link } from "react-router-dom"
import { motion } from "framer-motion"
import {
  Target,
  HandCoins,
  TrendingUp,
  Wallet,
  Plus,
  ArrowRight,
  Loader2,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { StatsCard, StatsGrid, GoalCardCompact, ContributionCard } from "@/components/dashboard"
import { useAuth } from "@/contexts"
import { goalsApi, contributionsApi, type Goal, type Contribution } from "@/lib/api"
import { formatCurrency } from "@/lib/utils"
import { useToast } from "@/hooks/use-toast"

export function DashboardPage() {
  const { user } = useAuth()
  const { toast } = useToast()
  const [goals, setGoals] = useState<Goal[]>([])
  const [contributions, setContributions] = useState<Contribution[]>([])
  const [isLoading, setIsLoading] = useState(true)

  // Stats
  const [stats, setStats] = useState({
    totalGoals: 0,
    totalRaised: 0,
    totalContributed: 0,
    activeGoals: 0,
  })

  useEffect(() => {
    fetchDashboardData()
  }, [])

  const fetchDashboardData = async () => {
    setIsLoading(true)
    try {
      // Fetch goals and contributions in parallel
      const [goalsResponse, contributionsResponse] = await Promise.all([
        goalsApi.getMyGoals(),
        contributionsApi.getMyContributions(),
      ])

      const fetchedGoals = goalsResponse.goals || []
      const fetchedContributions = contributionsResponse.contributions || []

      setGoals(fetchedGoals.slice(0, 5))
      setContributions(fetchedContributions.slice(0, 3))

      // Calculate stats
      const totalRaised = fetchedGoals.reduce((sum, g) => sum + g.current_amount, 0)
      const totalContributed = fetchedContributions
        .filter((c) => c.status === "confirmed")
        .reduce((sum, c) => sum + c.amount, 0)
      const activeGoals = fetchedGoals.filter((g) => g.status === "open").length

      setStats({
        totalGoals: fetchedGoals.length,
        totalRaised,
        totalContributed,
        activeGoals,
      })
    } catch (error) {
      console.error("Failed to fetch dashboard data:", error)
      toast({
        variant: "destructive",
        title: "Failed to load dashboard",
        description: "Could not load your dashboard data. Please try again.",
      })
    } finally {
      setIsLoading(false)
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    )
  }

  return (
    <div className="space-y-8">
      {/* Welcome Section */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold">
            Welcome back, {user?.first_name || "there"}! 👋
          </h1>
          <p className="text-muted-foreground mt-1">
            Here's an overview of your goals and contributions
          </p>
        </div>
        <Button asChild>
          <Link to="/dashboard/goals/create" className="gap-2">
            <Plus className="w-4 h-4" />
            Create New Goal
          </Link>
        </Button>
      </div>

      {/* Stats Grid */}
      <StatsGrid>
        <StatsCard
          title="Total Goals"
          value={stats.totalGoals}
          subtitle={`${stats.activeGoals} active`}
          icon={Target}
        />
        <StatsCard
          title="Total Raised"
          value={formatCurrency(stats.totalRaised)}
          subtitle="Across all goals"
          icon={TrendingUp}
          trend={{ value: 12.5, isPositive: true }}
        />
        <StatsCard
          title="Total Contributed"
          value={formatCurrency(stats.totalContributed)}
          subtitle="To other goals"
          icon={HandCoins}
        />
        <StatsCard
          title="Available Balance"
          value={formatCurrency(0)}
          subtitle="Ready to withdraw"
          icon={Wallet}
        />
      </StatsGrid>

      {/* Goals and Contributions */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* My Goals */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="bg-card border border-border rounded-lg"
        >
          <div className="flex items-center justify-between p-4 border-b border-border">
            <h2 className="font-semibold">My Goals</h2>
            <Button variant="ghost" size="sm" asChild>
              <Link to="/dashboard/goals" className="gap-1">
                View all
                <ArrowRight className="w-4 h-4" />
              </Link>
            </Button>
          </div>
          <div className="p-4 space-y-3">
            {goals.length === 0 ? (
              <div className="text-center py-8">
                <Target className="w-10 h-10 mx-auto text-muted-foreground mb-2" />
                <p className="text-sm text-muted-foreground">
                  You haven't created any goals yet
                </p>
                <Button variant="outline" size="sm" className="mt-3" asChild>
                  <Link to="/dashboard/goals/create">Create your first goal</Link>
                </Button>
              </div>
            ) : (
              goals.map((goal) => <GoalCardCompact key={goal.id} goal={goal} />)
            )}
          </div>
        </motion.div>

        {/* Recent Contributions */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          className="bg-card border border-border rounded-lg"
        >
          <div className="flex items-center justify-between p-4 border-b border-border">
            <h2 className="font-semibold">Recent Contributions</h2>
            <Button variant="ghost" size="sm" asChild>
              <Link to="/dashboard/contributions" className="gap-1">
                View all
                <ArrowRight className="w-4 h-4" />
              </Link>
            </Button>
          </div>
          <div className="p-4 space-y-3">
            {contributions.length === 0 ? (
              <div className="text-center py-8">
                <HandCoins className="w-10 h-10 mx-auto text-muted-foreground mb-2" />
                <p className="text-sm text-muted-foreground">
                  You haven't made any contributions yet
                </p>
                <Button variant="outline" size="sm" className="mt-3" asChild>
                  <Link to="/dashboard/explore">Explore goals</Link>
                </Button>
              </div>
            ) : (
              contributions.map((contribution) => (
                <ContributionCard key={contribution.id} contribution={contribution} />
              ))
            )}
          </div>
        </motion.div>
      </div>

      {/* Quick Actions */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.3 }}
        className="bg-gradient-to-r from-primary/10 to-primary/5 border border-primary/20 rounded-lg p-6"
      >
        <h2 className="font-semibold mb-2">Discover Public Goals</h2>
        <p className="text-sm text-muted-foreground mb-4">
          Browse and contribute to public goals from the community. Help others achieve their dreams!
        </p>
        <Button asChild>
          <Link to="/dashboard/explore" className="gap-2">
            Explore Goals
            <ArrowRight className="w-4 h-4" />
          </Link>
        </Button>
      </motion.div>
    </div>
  )
}


