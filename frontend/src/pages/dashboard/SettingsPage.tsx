import { useState } from "react"
import { motion } from "framer-motion"
import { User, Bell, Shield, CreditCard, Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useAuth } from "@/contexts"
import { useToast } from "@/hooks/use-toast"
import { cn } from "@/lib/utils"

type TabKey = "profile" | "notifications" | "security" | "payment"

const tabs = [
  { key: "profile" as TabKey, label: "Profile", icon: User },
  { key: "notifications" as TabKey, label: "Notifications", icon: Bell },
  { key: "security" as TabKey, label: "Security", icon: Shield },
  { key: "payment" as TabKey, label: "Payment", icon: CreditCard },
]

export function SettingsPage() {
  const [activeTab, setActiveTab] = useState<TabKey>("profile")
  const { user } = useAuth()

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-8">
        <h1 className="text-2xl font-bold">Settings</h1>
        <p className="text-muted-foreground mt-1">
          Manage your account settings and preferences
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        {/* Tabs */}
        <nav className="space-y-1">
          {tabs.map((tab) => (
            <button
              key={tab.key}
              onClick={() => setActiveTab(tab.key)}
              className={cn(
                "w-full flex items-center gap-3 px-4 py-2.5 rounded-lg text-sm font-medium transition-colors text-left",
                activeTab === tab.key
                  ? "bg-primary/10 text-primary"
                  : "text-muted-foreground hover:bg-muted hover:text-foreground"
              )}
            >
              <tab.icon className="w-4 h-4" />
              {tab.label}
            </button>
          ))}
        </nav>

        {/* Content */}
        <div className="md:col-span-3">
          <motion.div
            key={activeTab}
            initial={{ opacity: 0, x: 10 }}
            animate={{ opacity: 1, x: 0 }}
            className="bg-card border border-border rounded-lg p-6"
          >
            {activeTab === "profile" && <ProfileSettings user={user} />}
            {activeTab === "notifications" && <NotificationSettings />}
            {activeTab === "security" && <SecuritySettings />}
            {activeTab === "payment" && <PaymentSettings />}
          </motion.div>
        </div>
      </div>
    </div>
  )
}

function ProfileSettings({ user }: { user: any }) {
  const [isLoading, setIsLoading] = useState(false)
  const { toast } = useToast()

  const handleSave = async () => {
    setIsLoading(true)
    await new Promise((r) => setTimeout(r, 1000))
    toast({ title: "Profile updated!", description: "Your changes have been saved." })
    setIsLoading(false)
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-lg font-semibold">Profile Information</h2>
        <p className="text-sm text-muted-foreground">
          Update your personal information
        </p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div className="space-y-2">
          <label className="text-sm font-medium">First Name</label>
          <input
            type="text"
            defaultValue={user?.first_name || ""}
            className="w-full px-4 py-2 bg-background border border-input rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
          />
        </div>
        <div className="space-y-2">
          <label className="text-sm font-medium">Last Name</label>
          <input
            type="text"
            defaultValue={user?.last_name || ""}
            className="w-full px-4 py-2 bg-background border border-input rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
          />
        </div>
      </div>

      <div className="space-y-2">
        <label className="text-sm font-medium">Email</label>
        <input
          type="email"
          defaultValue={user?.email || ""}
          disabled
          className="w-full px-4 py-2 bg-muted border border-input rounded-lg text-sm text-muted-foreground"
        />
        <p className="text-xs text-muted-foreground">
          Email cannot be changed
        </p>
      </div>

      <div className="space-y-2">
        <label className="text-sm font-medium">Phone Number</label>
        <input
          type="tel"
          defaultValue={user?.phone_number || ""}
          placeholder="+234 XXX XXX XXXX"
          className="w-full px-4 py-2 bg-background border border-input rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
        />
      </div>

      <Button onClick={handleSave} disabled={isLoading}>
        {isLoading && <Loader2 className="w-4 h-4 mr-2 animate-spin" />}
        Save Changes
      </Button>
    </div>
  )
}

function NotificationSettings() {
  const [settings, setSettings] = useState({
    email_contributions: true,
    email_withdrawals: true,
    email_updates: false,
    push_contributions: true,
    push_milestones: true,
  })

  const toggleSetting = (key: keyof typeof settings) => {
    setSettings((prev) => ({ ...prev, [key]: !prev[key] }))
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-lg font-semibold">Notification Preferences</h2>
        <p className="text-sm text-muted-foreground">
          Choose how you want to be notified
        </p>
      </div>

      <div className="space-y-4">
        <h3 className="text-sm font-medium">Email Notifications</h3>
        <div className="space-y-3">
          <ToggleRow
            label="Contribution received"
            description="When someone contributes to your goal"
            checked={settings.email_contributions}
            onChange={() => toggleSetting("email_contributions")}
          />
          <ToggleRow
            label="Withdrawal completed"
            description="When a withdrawal is processed"
            checked={settings.email_withdrawals}
            onChange={() => toggleSetting("email_withdrawals")}
          />
          <ToggleRow
            label="Product updates"
            description="News and feature announcements"
            checked={settings.email_updates}
            onChange={() => toggleSetting("email_updates")}
          />
        </div>
      </div>

      <div className="space-y-4">
        <h3 className="text-sm font-medium">Push Notifications</h3>
        <div className="space-y-3">
          <ToggleRow
            label="Real-time contributions"
            description="Instant notification for new contributions"
            checked={settings.push_contributions}
            onChange={() => toggleSetting("push_contributions")}
          />
          <ToggleRow
            label="Milestone alerts"
            description="When your goal reaches a milestone"
            checked={settings.push_milestones}
            onChange={() => toggleSetting("push_milestones")}
          />
        </div>
      </div>
    </div>
  )
}

function ToggleRow({
  label,
  description,
  checked,
  onChange,
}: {
  label: string
  description: string
  checked: boolean
  onChange: () => void
}) {
  return (
    <div className="flex items-center justify-between py-2">
      <div>
        <p className="text-sm font-medium">{label}</p>
        <p className="text-xs text-muted-foreground">{description}</p>
      </div>
      <button
        type="button"
        role="switch"
        aria-checked={checked}
        onClick={onChange}
        className={cn(
          "relative inline-flex h-6 w-11 items-center rounded-full transition-colors",
          checked ? "bg-primary" : "bg-muted"
        )}
      >
        <span
          className={cn(
            "inline-block h-4 w-4 transform rounded-full bg-white transition-transform",
            checked ? "translate-x-6" : "translate-x-1"
          )}
        />
      </button>
    </div>
  )
}

function SecuritySettings() {
  const [isLoading, setIsLoading] = useState(false)
  const { toast } = useToast()

  const handleChangePassword = async () => {
    setIsLoading(true)
    await new Promise((r) => setTimeout(r, 1000))
    toast({ title: "Password updated!", description: "Your password has been changed." })
    setIsLoading(false)
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-lg font-semibold">Security Settings</h2>
        <p className="text-sm text-muted-foreground">
          Manage your password and account security
        </p>
      </div>

      <div className="space-y-4">
        <div className="space-y-2">
          <label className="text-sm font-medium">Current Password</label>
          <input
            type="password"
            placeholder="••••••••"
            className="w-full px-4 py-2 bg-background border border-input rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
          />
        </div>
        <div className="space-y-2">
          <label className="text-sm font-medium">New Password</label>
          <input
            type="password"
            placeholder="••••••••"
            className="w-full px-4 py-2 bg-background border border-input rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
          />
        </div>
        <div className="space-y-2">
          <label className="text-sm font-medium">Confirm New Password</label>
          <input
            type="password"
            placeholder="••••••••"
            className="w-full px-4 py-2 bg-background border border-input rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
          />
        </div>
      </div>

      <Button onClick={handleChangePassword} disabled={isLoading}>
        {isLoading && <Loader2 className="w-4 h-4 mr-2 animate-spin" />}
        Update Password
      </Button>
    </div>
  )
}

function PaymentSettings() {
  const { user } = useAuth()
  const hasSettlementAccount = user?.settlement_account_status === "verified"

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-lg font-semibold">Payment Settings</h2>
        <p className="text-sm text-muted-foreground">
          Manage your payment methods and settlement account
        </p>
      </div>

      <div className="p-4 bg-muted/50 rounded-lg border border-border">
        <div className="flex items-center justify-between">
          <div>
            <p className="font-medium">Settlement Account</p>
            <p className="text-sm text-muted-foreground">
              {hasSettlementAccount
                ? "Your settlement account is set up and verified"
                : "Set up your bank account to receive withdrawals"}
            </p>
          </div>
          <span
            className={cn(
              "px-2 py-1 text-xs font-medium rounded-full",
              hasSettlementAccount
                ? "bg-green-500/10 text-green-500"
                : "bg-yellow-500/10 text-yellow-500"
            )}
          >
            {hasSettlementAccount ? "Verified" : "Not Set Up"}
          </span>
        </div>
        {!hasSettlementAccount && (
          <Button variant="outline" size="sm" className="mt-4">
            Set Up Settlement Account
          </Button>
        )}
      </div>

      <div className="p-4 bg-muted/50 rounded-lg border border-border">
        <div className="flex items-center justify-between">
          <div>
            <p className="font-medium">KYC Verification</p>
            <p className="text-sm text-muted-foreground">
              {user?.kyc_status === "verified"
                ? "Your identity has been verified"
                : "Complete KYC to increase your withdrawal limits"}
            </p>
          </div>
          <span
            className={cn(
              "px-2 py-1 text-xs font-medium rounded-full",
              user?.kyc_status === "verified"
                ? "bg-green-500/10 text-green-500"
                : "bg-yellow-500/10 text-yellow-500"
            )}
          >
            {user?.kyc_status === "verified" ? "Verified" : "Pending"}
          </span>
        </div>
        {user?.kyc_status !== "verified" && (
          <Button variant="outline" size="sm" className="mt-4">
            Complete KYC
          </Button>
        )}
      </div>
    </div>
  )
}
