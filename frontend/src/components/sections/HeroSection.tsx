import { Link } from "react-router-dom"
import { motion } from "framer-motion"
import { ArrowRight, Play, Shield, Zap, Users } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  SplitText,
  BlurText,
  GradientText,
  FadeIn,
  AuroraBackground,
  MagneticButton,
} from "@/components/animations"

export function HeroSection() {
  const stats = [
    { value: "₦50M+", label: "Funds Raised" },
    { value: "10K+", label: "Goals Funded" },
    { value: "50K+", label: "Active Users" },
  ]

  return (
    <section className="relative min-h-[calc(100vh-4rem)] flex items-center overflow-hidden">
      {/* Background */}
      <AuroraBackground />
      
      {/* Gradient Orbs */}
      <motion.div
        className="absolute top-20 -left-32 w-96 h-96 bg-primary/20 rounded-full blur-3xl"
        animate={{
          scale: [1, 1.2, 1],
          opacity: [0.3, 0.5, 0.3],
        }}
        transition={{ duration: 8, repeat: Infinity }}
      />
      <motion.div
        className="absolute bottom-20 -right-32 w-96 h-96 bg-accent/20 rounded-full blur-3xl"
        animate={{
          scale: [1.2, 1, 1.2],
          opacity: [0.3, 0.5, 0.3],
        }}
        transition={{ duration: 8, repeat: Infinity }}
      />

      <div className="relative z-10 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          {/* Left Content */}
          <div className="text-center lg:text-left">
            {/* Badge */}
            <FadeIn delay={0.1}>
              <motion.div
                whileHover={{ scale: 1.05 }}
                className="inline-flex items-center gap-2 px-4 py-2 rounded-full glass mb-6"
              >
                <span className="relative flex h-2 w-2">
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-2 w-2 bg-primary"></span>
                </span>
                <span className="text-sm text-muted-foreground">
                  Trusted by 50,000+ community members
                </span>
              </motion.div>
            </FadeIn>

            {/* Headline */}
            <div className="mb-6">
              <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold leading-tight">
                <SplitText
                  text="Fund Your Goals"
                  className="mb-2"
                  textAlign="left"
                />
                <span className="block">
                  <GradientText text="Together" className="text-4xl sm:text-5xl lg:text-6xl font-bold" />
                </span>
              </h1>
            </div>

            {/* Description */}
            <FadeIn delay={0.3}>
              <p className="text-lg text-muted-foreground mb-8 max-w-xl mx-auto lg:mx-0">
                The most transparent and accountable group funding platform. 
                Create goals, track contributions, and achieve milestones 
                with full financial transparency.
              </p>
            </FadeIn>

            {/* CTA Buttons */}
            <FadeIn delay={0.4}>
              <div className="flex flex-col sm:flex-row gap-4 justify-center lg:justify-start">
                <MagneticButton>
                  <Link to="/register">
                    <Button size="xl" variant="gradient" className="w-full sm:w-auto gap-2">
                      Start Your Goal
                      <ArrowRight className="w-5 h-5" />
                    </Button>
                  </Link>
                </MagneticButton>
                <MagneticButton>
                  <Button size="xl" variant="outline" className="w-full sm:w-auto gap-2">
                    <Play className="w-5 h-5" />
                    Watch Demo
                  </Button>
                </MagneticButton>
              </div>
            </FadeIn>

            {/* Trust Indicators */}
            <FadeIn delay={0.5}>
              <div className="flex flex-wrap items-center gap-6 mt-8 justify-center lg:justify-start">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Shield className="w-4 h-4 text-primary" />
                  Bank-level Security
                </div>
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Zap className="w-4 h-4 text-accent" />
                  Instant Withdrawals
                </div>
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Users className="w-4 h-4 text-primary" />
                  Community Verified
                </div>
              </div>
            </FadeIn>
          </div>

          {/* Right Content - Stats & Visual */}
          <div className="relative">
            {/* Main Card */}
            <FadeIn direction="left" delay={0.3}>
              <motion.div
                className="relative glass rounded-2xl p-6 max-w-md mx-auto"
                whileHover={{ y: -5 }}
                transition={{ type: "spring", stiffness: 300 }}
              >
                {/* Progress Circle */}
                <div className="flex items-center justify-center mb-6">
                  <div className="relative w-40 h-40">
                    <svg className="w-full h-full transform -rotate-90">
                      <circle
                        cx="80"
                        cy="80"
                        r="70"
                        stroke="currentColor"
                        strokeWidth="8"
                        fill="none"
                        className="text-muted"
                      />
                      <motion.circle
                        cx="80"
                        cy="80"
                        r="70"
                        stroke="url(#gradient)"
                        strokeWidth="8"
                        fill="none"
                        strokeLinecap="round"
                        initial={{ pathLength: 0 }}
                        whileInView={{ pathLength: 0.78 }}
                        viewport={{ once: true }}
                        transition={{ duration: 2, ease: "easeOut" }}
                        style={{
                          strokeDasharray: "440",
                        }}
                      />
                      <defs>
                        <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="0%">
                          <stop offset="0%" stopColor="#10b981" />
                          <stop offset="100%" stopColor="#f59e0b" />
                        </linearGradient>
                      </defs>
                    </svg>
                    <div className="absolute inset-0 flex flex-col items-center justify-center">
                      <BlurText text="78%" className="text-3xl font-bold" />
                      <span className="text-sm text-muted-foreground">Funded</span>
                    </div>
                  </div>
                </div>

                {/* Goal Info */}
                <div className="text-center mb-4">
                  <h3 className="font-semibold text-lg mb-1">Community School Project</h3>
                  <p className="text-muted-foreground text-sm">
                    Building a better future for our children
                  </p>
                </div>

                {/* Stats */}
                <div className="grid grid-cols-2 gap-4">
                  <div className="text-center p-3 rounded-lg bg-muted/50">
                    <div className="text-xl font-bold text-primary">₦7.8M</div>
                    <div className="text-xs text-muted-foreground">Raised</div>
                  </div>
                  <div className="text-center p-3 rounded-lg bg-muted/50">
                    <div className="text-xl font-bold text-accent">₦10M</div>
                    <div className="text-xs text-muted-foreground">Target</div>
                  </div>
                </div>

                {/* Contributors */}
                <div className="mt-4 flex items-center justify-between">
                  <div className="flex -space-x-2">
                    {[1, 2, 3, 4, 5].map((i) => (
                      <motion.div
                        key={i}
                        initial={{ scale: 0, x: -20 }}
                        whileInView={{ scale: 1, x: 0 }}
                        viewport={{ once: true }}
                        transition={{ delay: i * 0.1 }}
                        className="w-8 h-8 rounded-full bg-gradient-to-br from-primary to-accent border-2 border-card"
                      />
                    ))}
                    <motion.div
                      initial={{ scale: 0 }}
                      whileInView={{ scale: 1 }}
                      viewport={{ once: true }}
                      transition={{ delay: 0.6 }}
                      className="w-8 h-8 rounded-full bg-muted border-2 border-card flex items-center justify-center text-xs font-medium"
                    >
                      +99
                    </motion.div>
                  </div>
                  <span className="text-sm text-muted-foreground">104 contributors</span>
                </div>
              </motion.div>
            </FadeIn>

            {/* Floating Cards */}
            <motion.div
              className="absolute -top-4 -left-4 glass rounded-xl p-4 hidden lg:block"
              animate={{ y: [0, -10, 0] }}
              transition={{ duration: 4, repeat: Infinity }}
            >
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-primary/20 flex items-center justify-center">
                  <Shield className="w-5 h-5 text-primary" />
                </div>
                <div>
                  <div className="text-sm font-medium">Verified Payment</div>
                  <div className="text-xs text-muted-foreground">₦50,000 • Just now</div>
                </div>
              </div>
            </motion.div>

            <motion.div
              className="absolute -bottom-4 -right-4 glass rounded-xl p-4 hidden lg:block"
              animate={{ y: [0, 10, 0] }}
              transition={{ duration: 4, repeat: Infinity, delay: 2 }}
            >
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-accent/20 flex items-center justify-center">
                  <Zap className="w-5 h-5 text-accent" />
                </div>
                <div>
                  <div className="text-sm font-medium">Goal Achieved! 🎉</div>
                  <div className="text-xs text-muted-foreground">Wedding Fund</div>
                </div>
              </div>
            </motion.div>
          </div>
        </div>

        {/* Bottom Stats */}
        <FadeIn delay={0.6}>
          <div className="mt-16 grid grid-cols-3 gap-8 max-w-2xl mx-auto">
            {stats.map((stat, index) => (
              <motion.div
                key={stat.label}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: 0.7 + index * 0.1 }}
                className="text-center"
              >
                <div className="text-3xl sm:text-4xl font-bold">
                  <GradientText text={stat.value} />
                </div>
                <div className="text-sm text-muted-foreground mt-1">{stat.label}</div>
              </motion.div>
            ))}
          </div>
        </FadeIn>
      </div>
    </section>
  )
}
