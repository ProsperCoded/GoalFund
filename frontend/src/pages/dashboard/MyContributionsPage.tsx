import { useState, useEffect } from "react"
import { motion } from "framer-motion"
import { HandCoins, Search, Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { ContributionCard, ContributionRow } from "@/components/dashboard"
import { contributionsApi, type Contribution } from "@/lib/api"
import { formatCurrency } from "@/lib/utils"

type ViewMode = "cards" | "table"
type FilterStatus = "all" | "confirmed" | "pending" | "refunded"

export function MyContributionsPage() {
  const [contributions, setContributions] = useState<Contribution[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState("")
  const [statusFilter, setStatusFilter] = useState<FilterStatus>("all")
  const [viewMode] = useState<ViewMode>("cards")

  useEffect(() => {
    fetchContributions()
  }, [])

  const fetchContributions = async () => {
    setIsLoading(true)
    try {
      const response = await contributionsApi.getMyContributions()
      setContributions(response.contributions || [])
    } catch (error) {
      console.error("Failed to fetch contributions:", error)
      setContributions([])
    } finally {
      setIsLoading(false)
    }
  }

  const filteredContributions = contributions.filter((c) => {
    const matchesSearch = (c.goal_title || "").toLowerCase().includes(searchQuery.toLowerCase())
    const matchesStatus = statusFilter === "all" || c.status === statusFilter
    return matchesSearch && matchesStatus
  })

  const totalConfirmed = contributions
    .filter((c) => c.status === "confirmed")
    .reduce((sum, c) => sum + c.amount, 0)

  const totalPending = contributions
    .filter((c) => c.status === "pending")
    .reduce((sum, c) => sum + c.amount, 0)

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold">My Contributions</h1>
          <p className="text-muted-foreground mt-1">
            Track all your contributions to various goals
          </p>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="bg-card border border-border rounded-lg p-4">
          <p className="text-sm text-muted-foreground">Total Contributed</p>
          <p className="text-2xl font-bold mt-1">{formatCurrency(totalConfirmed)}</p>
        </div>
        <div className="bg-card border border-border rounded-lg p-4">
          <p className="text-sm text-muted-foreground">Pending</p>
          <p className="text-2xl font-bold mt-1 text-yellow-500">
            {formatCurrency(totalPending)}
          </p>
        </div>
        <div className="bg-card border border-border rounded-lg p-4">
          <p className="text-sm text-muted-foreground">Goals Supported</p>
          <p className="text-2xl font-bold mt-1">
            {new Set(contributions.map((c) => c.goal_id)).size}
          </p>
        </div>
      </div>

      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-4">
        {/* Search */}
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <input
            type="text"
            placeholder="Search by goal name..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-10 pr-4 py-2 bg-background border border-input rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
          />
        </div>

        {/* Status Filter */}
        <div className="flex items-center gap-2 flex-wrap">
          {(["all", "confirmed", "pending", "refunded"] as FilterStatus[]).map((status) => (
            <Button
              key={status}
              variant={statusFilter === status ? "default" : "outline"}
              size="sm"
              onClick={() => setStatusFilter(status)}
              className="capitalize"
            >
              {status}
            </Button>
          ))}
        </div>
      </div>

      {/* Content */}
      {isLoading ? (
        <div className="flex items-center justify-center min-h-[300px]">
          <Loader2 className="w-8 h-8 animate-spin text-primary" />
        </div>
      ) : filteredContributions.length === 0 ? (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="flex flex-col items-center justify-center min-h-[300px] text-center"
        >
          <HandCoins className="w-12 h-12 text-muted-foreground mb-4" />
          <h3 className="font-semibold text-lg mb-2">
            {searchQuery || statusFilter !== "all"
              ? "No contributions found"
              : "No contributions yet"}
          </h3>
          <p className="text-muted-foreground mb-4 max-w-md">
            {searchQuery || statusFilter !== "all"
              ? "Try adjusting your search or filters"
              : "Explore public goals and make your first contribution!"}
          </p>
        </motion.div>
      ) : viewMode === "table" ? (
        <div className="bg-card border border-border rounded-lg overflow-hidden">
          <table className="w-full">
            <thead>
              <tr className="border-b border-border bg-muted/50">
                <th className="text-left py-3 px-4 text-sm font-medium">Goal</th>
                <th className="text-right py-3 px-4 text-sm font-medium">Amount</th>
                <th className="text-left py-3 px-4 text-sm font-medium">Status</th>
                <th className="text-left py-3 px-4 text-sm font-medium">Date</th>
              </tr>
            </thead>
            <tbody>
              {filteredContributions.map((contribution) => (
                <ContributionRow key={contribution.id} contribution={contribution} />
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {filteredContributions.map((contribution, index) => (
            <motion.div
              key={contribution.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.05 }}
            >
              <ContributionCard contribution={contribution} />
            </motion.div>
          ))}
        </div>
      )}
    </div>
  )
}


