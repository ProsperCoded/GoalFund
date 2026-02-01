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

export function DashboardPage() {
  const { user } = useAuth()
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
        goalsApi.getMyGoals().catch(() => ({ goals: getMockGoals() })),
        contributionsApi.getMyContributions().catch(() => ({ contributions: getMockContributions() })),
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
      // Use mock data for demo
      const mockGoals = getMockGoals()
      const mockContributions = getMockContributions()
      setGoals(mockGoals)
      setContributions(mockContributions)
      setStats({
        totalGoals: mockGoals.length,
        totalRaised: mockGoals.reduce((sum, g) => sum + g.current_amount, 0),
        totalContributed: mockContributions.reduce((sum, c) => sum + c.amount, 0),
        activeGoals: mockGoals.filter((g) => g.status === "open").length,
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

// Mock data for demo
function getMockGoals(): Goal[] {
  return [
    {
      id: "1",
      user_id: "user1",
      title: "Community Borehole Project",
      description: "Help us provide clean water to the community",
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
      description: "Raising funds for my university tuition",
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
      description: "Help with medical expenses",
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
  ]
}

function getMockContributions(): Contribution[] {
  return [
    {
      id: "1",
      goal_id: "g1",
      goal_title: "Build a Library for Rural Kids",
      user_id: "user1",
      amount: 50000,
      status: "confirmed",
      is_anonymous: false,
      created_at: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date().toISOString(),
    },
    {
      id: "2",
      goal_id: "g2",
      goal_title: "Support Local Farmers Initiative",
      user_id: "user1",
      amount: 25000,
      status: "confirmed",
      is_anonymous: false,
      message: "Keep up the great work!",
      created_at: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date().toISOString(),
    },
    {
      id: "3",
      goal_id: "g3",
      goal_title: "Medical Fund for Mama Ada",
      user_id: "user1",
      amount: 100000,
      status: "pending",
      is_anonymous: true,
      created_at: new Date(Date.now() - 48 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date().toISOString(),
    },
  ]
}
