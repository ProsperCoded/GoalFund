import { motion } from "framer-motion"
import {
  Target,
  Shield,
  Wallet,
  Users,
  TrendingUp,
  CheckCircle2,
  ArrowRight,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  SplitText,
  GradientText,
  FadeIn,
  StaggerContainer,
  StaggerItem,
  TiltCard,
} from "@/components/animations"

export function AboutSection() {
  const features = [
    {
      icon: Shield,
      title: "Bank-Level Security",
      description:
        "All transactions are encrypted and secured with industry-standard protocols. Your money is safe with us.",
      color: "text-primary",
      bg: "bg-primary/10",
    },
    {
      icon: Wallet,
      title: "Transparent Ledger",
      description:
        "Every naira is tracked with immutable ledger entries. No hidden fees, no surprises - just pure transparency.",
      color: "text-accent",
      bg: "bg-accent/10",
    },
    {
      icon: Users,
      title: "Community Verified",
      description:
        "Contributors can vote on proof of accomplishment, ensuring funds are used as intended.",
      color: "text-primary",
      bg: "bg-primary/10",
    },
    {
      icon: TrendingUp,
      title: "Continuous Funding",
      description:
        "Goals aren't capped. Continue receiving contributions beyond your target and withdraw anytime.",
      color: "text-accent",
      bg: "bg-accent/10",
    },
  ]

  const benefits = [
    "No minimum contribution amount",
    "Instant payment verification",
    "Multiple withdrawal support",
    "Milestone-based tracking",
    "Real-time notifications",
    "Proof of accomplishment",
  ]

  return (
    <section className="py-24 relative overflow-hidden">
      {/* Background Pattern */}
      <div className="absolute inset-0 opacity-30">
        <div
          className="absolute inset-0"
          style={{
            backgroundImage: `radial-gradient(circle at 2px 2px, rgba(16, 185, 129, 0.15) 1px, transparent 0)`,
            backgroundSize: "40px 40px",
          }}
        />
      </div>

      <div className="relative z-10 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="text-center mb-16">
          <FadeIn>
            <motion.div
              whileHover={{ scale: 1.05 }}
              className="inline-flex items-center gap-2 px-4 py-2 rounded-full glass mb-6"
            >
              <Target className="w-4 h-4 text-primary" />
              <span className="text-sm text-muted-foreground">
                Why Choose GoFund?
              </span>
            </motion.div>
          </FadeIn>

          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-bold mb-6">
            <SplitText
              text="Built for Trust,"
              className="mb-2"
            />
            <span className="block">
              <GradientText text="Designed for Impact" className="text-3xl sm:text-4xl lg:text-5xl font-bold" />
            </span>
          </h2>

          <FadeIn delay={0.2}>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              GoFund combines the power of community with the precision of fintech.
              Every contribution is tracked, every withdrawal is auditable, and
              every goal is achievable.
            </p>
          </FadeIn>
        </div>

        {/* Features Grid */}
        <StaggerContainer className="grid md:grid-cols-2 gap-6 mb-16" staggerDelay={0.1}>
          {features.map((feature) => {
            const Icon = feature.icon
            return (
              <StaggerItem key={feature.title}>
                <TiltCard>
                  <motion.div
                    whileHover={{ borderColor: "rgba(16, 185, 129, 0.5)" }}
                    className="h-full glass rounded-2xl p-6 border border-transparent transition-colors"
                  >
                    <div
                      className={`w-14 h-14 rounded-xl ${feature.bg} flex items-center justify-center mb-4`}
                    >
                      <Icon className={`w-7 h-7 ${feature.color}`} />
                    </div>
                    <h3 className="text-xl font-semibold mb-2">{feature.title}</h3>
                    <p className="text-muted-foreground">{feature.description}</p>
                  </motion.div>
                </TiltCard>
              </StaggerItem>
            )
          })}
        </StaggerContainer>

        {/* Benefits & CTA */}
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          {/* Benefits List */}
          <FadeIn direction="right">
            <div className="glass rounded-2xl p-8">
              <h3 className="text-2xl font-bold mb-6">
                Everything you need to{" "}
                <GradientText text="succeed" />
              </h3>
              <ul className="space-y-4">
                {benefits.map((benefit, index) => (
                  <motion.li
                    key={benefit}
                    initial={{ opacity: 0, x: -20 }}
                    whileInView={{ opacity: 1, x: 0 }}
                    viewport={{ once: true }}
                    transition={{ delay: index * 0.1 }}
                    className="flex items-center gap-3"
                  >
                    <CheckCircle2 className="w-5 h-5 text-primary flex-shrink-0" />
                    <span>{benefit}</span>
                  </motion.li>
                ))}
              </ul>
            </div>
          </FadeIn>

          {/* Visual / Stats */}
          <FadeIn direction="left">
            <div className="relative">
              {/* Main Card */}
              <motion.div
                whileHover={{ scale: 1.02 }}
                className="glass rounded-2xl p-8 text-center"
              >
                <div className="mb-6">
                  <div className="text-6xl font-bold mb-2">
                    <GradientText text="99.9%" />
                  </div>
                  <p className="text-muted-foreground">Transaction Success Rate</p>
                </div>

                <div className="grid grid-cols-2 gap-4 mb-6">
                  <div className="p-4 rounded-xl bg-muted/50">
                    <div className="text-2xl font-bold text-primary">2 mins</div>
                    <div className="text-sm text-muted-foreground">Avg. Verification</div>
                  </div>
                  <div className="p-4 rounded-xl bg-muted/50">
                    <div className="text-2xl font-bold text-accent">24/7</div>
                    <div className="text-sm text-muted-foreground">Support Available</div>
                  </div>
                </div>

                <Button variant="gradient" size="lg" className="w-full gap-2">
                  Start Your Journey
                  <ArrowRight className="w-5 h-5" />
                </Button>
              </motion.div>

              {/* Decorative Elements */}
              <motion.div
                className="absolute -top-4 -right-4 w-24 h-24 bg-primary/20 rounded-full blur-2xl"
                animate={{ scale: [1, 1.2, 1] }}
                transition={{ duration: 4, repeat: Infinity }}
              />
              <motion.div
                className="absolute -bottom-4 -left-4 w-32 h-32 bg-accent/20 rounded-full blur-2xl"
                animate={{ scale: [1.2, 1, 1.2] }}
                transition={{ duration: 4, repeat: Infinity }}
              />
            </div>
          </FadeIn>
        </div>
      </div>
    </section>
  )
}
