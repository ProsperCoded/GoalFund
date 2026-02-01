import { Link } from "react-router-dom"
import { motion } from "framer-motion"
import { Target, ArrowUpRight, Clock, CheckCircle, XCircle } from "lucide-react"
import { cn, formatCurrency, formatDate, formatDistanceToNow } from "@/lib/utils"
import type { Contribution } from "@/lib/api"

interface ContributionCardProps {
  contribution: Contribution
}

export function ContributionCard({ contribution }: ContributionCardProps) {
  const statusConfig: Record<
    string,
    { icon: typeof CheckCircle; color: string; label: string }
  > = {
    pending: {
      icon: Clock,
      color: "text-yellow-500 bg-yellow-500/10 border-yellow-500/20",
      label: "Pending",
    },
    confirmed: {
      icon: CheckCircle,
      color: "text-green-500 bg-green-500/10 border-green-500/20",
      label: "Confirmed",
    },
    failed: {
      icon: XCircle,
      color: "text-red-500 bg-red-500/10 border-red-500/20",
      label: "Failed",
    },
    refunded: {
      icon: ArrowUpRight,
      color: "text-blue-500 bg-blue-500/10 border-blue-500/20",
      label: "Refunded",
    },
  }

  const status = statusConfig[contribution.status] || statusConfig.pending
  const StatusIcon = status.icon

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="bg-card border border-border rounded-lg p-4 hover:border-primary/30 transition-all"
    >
      <div className="flex items-start gap-4">
        <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0">
          <Target className="w-5 h-5 text-primary" />
        </div>

        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between gap-2">
            <Link
              to={`/dashboard/goals/${contribution.goal_id}`}
              className="font-medium hover:text-primary transition-colors truncate"
            >
              {contribution.goal_title || "Goal"}
            </Link>
            <span
              className={cn(
                "px-2 py-0.5 text-xs font-medium rounded-full border flex items-center gap-1",
                status.color
              )}
            >
              <StatusIcon className="w-3 h-3" />
              {status.label}
            </span>
          </div>

          <div className="flex items-center justify-between mt-2">
            <p className="text-lg font-semibold">
              {formatCurrency(contribution.amount)}
            </p>
            <p className="text-xs text-muted-foreground">
              {formatDistanceToNow(contribution.created_at)}
            </p>
          </div>

          {contribution.message && (
            <p className="text-sm text-muted-foreground mt-2 line-clamp-2">
              "{contribution.message}"
            </p>
          )}
        </div>
      </div>
    </motion.div>
  )
}

// Table row version for lists
export function ContributionRow({ contribution }: ContributionCardProps) {
  const statusColors: Record<string, string> = {
    pending: "bg-yellow-500",
    confirmed: "bg-green-500",
    failed: "bg-red-500",
    refunded: "bg-blue-500",
  }

  return (
    <tr className="border-b border-border last:border-0 hover:bg-muted/50 transition-colors">
      <td className="py-3 px-4">
        <Link
          to={`/dashboard/goals/${contribution.goal_id}`}
          className="font-medium hover:text-primary transition-colors"
        >
          {contribution.goal_title || "Goal"}
        </Link>
      </td>
      <td className="py-3 px-4 text-right font-medium">
        {formatCurrency(contribution.amount)}
      </td>
      <td className="py-3 px-4">
        <div className="flex items-center gap-2">
          <span
            className={cn(
              "w-2 h-2 rounded-full",
              statusColors[contribution.status]
            )}
          />
          <span className="capitalize text-sm">{contribution.status}</span>
        </div>
      </td>
      <td className="py-3 px-4 text-sm text-muted-foreground">
        {formatDate(contribution.created_at)}
      </td>
    </tr>
  )
}
