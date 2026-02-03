import { Link } from "react-router-dom"
import { motion } from "framer-motion"
import { Target, Users, Calendar, MoreVertical, Eye, Edit, Trash2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { cn, formatCurrency, formatDate } from "@/lib/utils"
import type { Goal } from "@/lib/api"

interface GoalCardProps {
  goal: Goal
  showActions?: boolean
  onDelete?: (id: string) => void
}

export function GoalCard({ goal, showActions = false, onDelete }: GoalCardProps) {
  // Handle potential NaN by defaulting to 0
  const currentAmount = goal.current_amount || 0
  const targetAmount = goal.target_amount || 1
  const contributorCount = goal.contributor_count || 0
  
  const progress = Math.min((currentAmount / targetAmount) * 100, 100)
  const isOverfunded = currentAmount > targetAmount

  // Normalize status to lowercase for comparison
  const statusLower = goal.status?.toLowerCase() || "open"
  
  const statusColors: Record<string, string> = {
    open: "bg-green-500/10 text-green-500 border-green-500/20",
    funded: "bg-blue-500/10 text-blue-500 border-blue-500/20",
    withdrawn: "bg-purple-500/10 text-purple-500 border-purple-500/20",
    closed: "bg-gray-500/10 text-gray-500 border-gray-500/20",
    cancelled: "bg-red-500/10 text-red-500 border-red-500/20",
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="bg-card border border-border rounded-lg p-5 hover:border-primary/30 transition-all group"
    >
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-2">
            <span
              className={cn(
                "px-2 py-0.5 text-xs font-medium rounded-full border capitalize",
                statusColors[statusLower] || statusColors.open
              )}
            >
              {statusLower}
            </span>
            {!goal.is_public && (
              <span className="px-2 py-0.5 text-xs font-medium rounded-full bg-muted text-muted-foreground border border-border">
                Private
              </span>
            )}
          </div>

          <Link to={`/dashboard/goals/${goal.id}`}>
            <h3 className="font-semibold text-lg group-hover:text-primary transition-colors line-clamp-1">
              {goal.title}
            </h3>
          </Link>

          <p className="text-sm text-muted-foreground mt-1 line-clamp-2">
            {goal.description}
          </p>
        </div>

        {showActions && (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="flex-shrink-0">
                <MoreVertical className="w-4 h-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem asChild>
                <Link to={`/dashboard/goals/${goal.id}`}>
                  <Eye className="w-4 h-4 mr-2" />
                  View Details
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link to={`/dashboard/goals/${goal.id}/edit`}>
                  <Edit className="w-4 h-4 mr-2" />
                  Edit Goal
                </Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                className="text-destructive"
                onClick={() => onDelete?.(goal.id)}
              >
                <Trash2 className="w-4 h-4 mr-2" />
                Delete Goal
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )}
      </div>

      {/* Progress */}
      <div className="mt-4">
        <div className="flex items-center justify-between text-sm mb-2">
          <span className="font-medium">
            {formatCurrency(currentAmount)}
            {isOverfunded && (
              <span className="text-green-500 ml-1">
                (+{formatCurrency(currentAmount - targetAmount)})
              </span>
            )}
          </span>
          <span className="text-muted-foreground">
            of {formatCurrency(targetAmount)}
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
        <p className="text-xs text-muted-foreground mt-1">
          {isNaN(progress) ? 0 : progress.toFixed(0)}% funded
        </p>
      </div>

      {/* Meta */}
      <div className="flex items-center gap-4 mt-4 pt-4 border-t border-border text-sm text-muted-foreground">
        <div className="flex items-center gap-1.5">
          <Users className="w-4 h-4" />
          <span>{contributorCount} contributors</span>
        </div>
        {goal.deadline && (
          <div className="flex items-center gap-1.5">
            <Calendar className="w-4 h-4" />
            <span>{formatDate(goal.deadline)}</span>
          </div>
        )}
      </div>
    </motion.div>
  )
}

// Compact version for lists
export function GoalCardCompact({ goal }: { goal: Goal }) {
  // Handle potential NaN by defaulting to 0
  const currentAmount = goal.current_amount || 0
  const targetAmount = goal.target_amount || 1
  
  const progress = Math.min((currentAmount / targetAmount) * 100, 100)

  return (
    <Link
      to={`/dashboard/goals/${goal.id}`}
      className="flex items-center gap-4 p-4 bg-card border border-border rounded-lg hover:border-primary/30 transition-all"
    >
      <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0">
        <Target className="w-5 h-5 text-primary" />
      </div>
      <div className="flex-1 min-w-0">
        <h4 className="font-medium text-sm truncate">{goal.title}</h4>
        <div className="flex items-center gap-2 mt-1">
          <div className="flex-1 h-1.5 bg-muted rounded-full overflow-hidden">
            <div
              className="h-full bg-primary rounded-full"
              style={{ width: `${isNaN(progress) ? 0 : progress}%` }}
            />
          </div>
          <span className="text-xs text-muted-foreground flex-shrink-0">
            {isNaN(progress) ? 0 : progress.toFixed(0)}%
          </span>
        </div>
      </div>
      <div className="text-right flex-shrink-0">
        <p className="text-sm font-medium">{formatCurrency(currentAmount)}</p>
        <p className="text-xs text-muted-foreground">
          of {formatCurrency(targetAmount)}
        </p>
      </div>
    </Link>
  )
}
