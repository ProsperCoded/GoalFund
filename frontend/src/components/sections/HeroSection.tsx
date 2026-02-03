import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import {
  ArrowRight,
  BookOpen,
  ShieldCheck,
  Eye,
  Users,
  ChevronRight,
} from "lucide-react";
import { useAuth } from "@/contexts";
import { Button } from "@/components/ui/button";
import {
  SplitText,
  BlurText,
  GradientText,
  FadeIn,
  MagneticButton,
  Squares,
} from "@/components/animations";

export function HeroSection() {
  const { isAuthenticated } = useAuth();

  return (
    <section className="relative min-h-[calc(100vh-4rem)] flex items-center overflow-hidden">
      {/* Squares Background */}
      <div className="absolute inset-0">
        <Squares
          speed={0.3}
          squareSize={50}
          direction="diagonal"
          borderColor="rgba(212, 168, 83, 0.15)"
          hoverFillColor="rgba(212, 168, 83, 0.08)"
        />
      </div>

      {/* Accent elements */}
      <div className="absolute top-40 left-10 w-px h-32 bg-gradient-to-b from-transparent via-primary/40 to-transparent hidden lg:block" />
      <div className="absolute bottom-40 right-10 w-px h-32 bg-gradient-to-b from-transparent via-primary/40 to-transparent hidden lg:block" />

      {/* Decorative circles */}
      <motion.div
        className="absolute top-32 right-1/4 w-2 h-2 rounded-full bg-primary/60"
        animate={{ scale: [1, 1.5, 1], opacity: [0.6, 1, 0.6] }}
        transition={{ duration: 3, repeat: Infinity }}
      />
      <motion.div
        className="absolute bottom-40 left-1/4 w-1.5 h-1.5 rounded-full bg-primary/40"
        animate={{ scale: [1, 1.3, 1], opacity: [0.4, 0.8, 0.4] }}
        transition={{ duration: 4, repeat: Infinity, delay: 1 }}
      />

      <div className="relative z-10 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <div className="grid lg:grid-cols-12 gap-12 items-center">
          {/* Left Content - 7 columns */}
          <div className="lg:col-span-7 text-center lg:text-left">
            {/* Small label */}
            <FadeIn delay={0.1}>
              <div className="inline-flex items-center gap-2 mb-8">
                <div className="w-8 h-px bg-primary/60" />
                <span className="text-xs uppercase tracking-[0.2em] text-primary font-medium">
                  Accountable Group Funding
                </span>
              </div>
            </FadeIn>

            {/* Headline */}
            <div className="mb-8">
              <h1 className="text-4xl sm:text-5xl lg:text-6xl xl:text-7xl font-bold leading-[1.1] tracking-tight">
                <SplitText
                  text="Fund Big Goals"
                  className="mb-3"
                  textAlign="left"
                />
                <span className="block mt-2">
                  With Full{" "}
                  <span className="relative inline-block">
                    <GradientText
                      text="Transparency"
                      className="text-4xl sm:text-5xl lg:text-6xl xl:text-7xl font-bold"
                    />
                    <motion.span
                      className="absolute -bottom-2 left-0 w-full h-0.5 bg-primary/30"
                      initial={{ scaleX: 0 }}
                      animate={{ scaleX: 1 }}
                      transition={{ delay: 1, duration: 0.8 }}
                    />
                  </span>
                </span>
              </h1>
            </div>

            {/* Description */}
            <FadeIn delay={0.3}>
              <p className="text-lg sm:text-xl text-muted-foreground mb-10 max-w-2xl mx-auto lg:mx-0 leading-relaxed">
                The platform for community projects, group contributions, and
                large-scale funding. Every naira tracked. Every transaction
                verified. Complete accountability.
              </p>
            </FadeIn>

            {/* CTA Buttons */}
            <FadeIn delay={0.4}>
              <div className="flex flex-col sm:flex-row gap-4 justify-center lg:justify-start mb-12">
                <MagneticButton>
                  <Link to={isAuthenticated ? "/dashboard" : "/register"}>
                    <Button
                      size="xl"
                      className="w-full sm:w-auto gap-2 bg-primary text-primary-foreground hover:bg-primary/90 px-8"
                    >
                      {isAuthenticated ? "Go to Dashboard" : "Create a Goal"}
                      <ArrowRight className="w-5 h-5" />
                    </Button>
                  </Link>
                </MagneticButton>
                <MagneticButton>
                  <Link to="/dashboard/explore">
                    <Button
                      size="xl"
                      variant="outline"
                      className="w-full sm:w-auto gap-2 border-primary/30 hover:bg-primary/5 hover:border-primary/50"
                    >
                      {isAuthenticated ? "Contribute Publicly" : "Explore Goals"}
                      <ChevronRight className="w-5 h-5" />
                    </Button>
                  </Link>
                </MagneticButton>
              </div>
            </FadeIn>

            {/* Trust Pillars */}
            <FadeIn delay={0.5}>
              <div className="grid grid-cols-3 gap-6 max-w-lg mx-auto lg:mx-0">
                <div className="text-center lg:text-left">
                  <div className="flex items-center justify-center lg:justify-start gap-2 mb-1">
                    <BookOpen className="w-4 h-4 text-primary" />
                    <span className="text-sm font-medium">Ledger-backed</span>
                  </div>
                  <p className="text-xs text-muted-foreground">
                    Immutable records
                  </p>
                </div>
                <div className="text-center lg:text-left">
                  <div className="flex items-center justify-center lg:justify-start gap-2 mb-1">
                    <Eye className="w-4 h-4 text-primary" />
                    <span className="text-sm font-medium">Transparent</span>
                  </div>
                  <p className="text-xs text-muted-foreground">
                    Public tracking
                  </p>
                </div>
                <div className="text-center lg:text-left">
                  <div className="flex items-center justify-center lg:justify-start gap-2 mb-1">
                    <ShieldCheck className="w-4 h-4 text-primary" />
                    <span className="text-sm font-medium">Verified</span>
                  </div>
                  <p className="text-xs text-muted-foreground">
                    Payment proofs
                  </p>
                </div>
              </div>
            </FadeIn>
          </div>

          {/* Right Content - 5 columns */}
          <div className="lg:col-span-5 relative">
            {/* Main Card */}
            <FadeIn direction="left" delay={0.3}>
              <div className="relative">
                {/* Card border glow */}
                <div className="absolute -inset-0.5 bg-gradient-to-b from-primary/20 to-transparent rounded-2xl blur-sm" />

                <motion.div
                  className="relative bg-card rounded-2xl p-6 border border-primary/10"
                  whileHover={{ y: -3 }}
                  transition={{ type: "spring", stiffness: 300 }}
                >
                  {/* Card header */}
                  <div className="flex items-center justify-between mb-6">
                    <span className="text-xs uppercase tracking-wider text-primary/70">
                      Live Project
                    </span>
                    <div className="flex items-center gap-1.5">
                      <span className="relative flex h-2 w-2">
                        <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
                        <span className="relative inline-flex rounded-full h-2 w-2 bg-primary"></span>
                      </span>
                      <span className="text-xs text-muted-foreground">
                        Active
                      </span>
                    </div>
                  </div>

                  {/* Progress visualization */}
                  <div className="flex items-center justify-center mb-6">
                    <div className="relative w-36 h-36">
                      <svg className="w-full h-full transform -rotate-90">
                        <circle
                          cx="72"
                          cy="72"
                          r="64"
                          stroke="currentColor"
                          strokeWidth="6"
                          fill="none"
                          className="text-muted/30"
                        />
                        <motion.circle
                          cx="72"
                          cy="72"
                          r="64"
                          stroke="url(#gold-gradient)"
                          strokeWidth="6"
                          fill="none"
                          strokeLinecap="round"
                          initial={{ pathLength: 0 }}
                          whileInView={{ pathLength: 0.65 }}
                          viewport={{ once: true }}
                          transition={{ duration: 2, ease: "easeOut" }}
                        />
                        <defs>
                          <linearGradient
                            id="gold-gradient"
                            x1="0%"
                            y1="0%"
                            x2="100%"
                            y2="0%"
                          >
                            <stop offset="0%" stopColor="#f5d78e" />
                            <stop offset="100%" stopColor="#b8942e" />
                          </linearGradient>
                        </defs>
                      </svg>
                      <div className="absolute inset-0 flex flex-col items-center justify-center">
                        <BlurText text="65%" className="text-3xl font-bold" />
                        <span className="text-xs text-muted-foreground">
                          Progress
                        </span>
                      </div>
                    </div>
                  </div>

                  {/* Project Info */}
                  <div className="text-center mb-5">
                    <h3 className="font-semibold text-lg mb-1">
                      Community Borehole Project
                    </h3>
                    <p className="text-muted-foreground text-sm">
                      Clean water access for Ikeja community
                    </p>
                  </div>

                  {/* Amount display */}
                  <div className="flex items-end justify-between px-2 mb-4">
                    <div>
                      <p className="text-xs text-muted-foreground mb-1">
                        Raised
                      </p>
                      <p className="text-2xl font-bold text-primary">₦6.5M</p>
                    </div>
                    <div className="text-right">
                      <p className="text-xs text-muted-foreground mb-1">
                        Target
                      </p>
                      <p className="text-2xl font-bold text-foreground/80">
                        ₦10M
                      </p>
                    </div>
                  </div>

                  {/* Progress bar */}
                  <div className="h-1.5 bg-muted rounded-full overflow-hidden mb-5">
                    <motion.div
                      className="h-full bg-gradient-to-r from-primary to-accent rounded-full"
                      initial={{ width: 0 }}
                      whileInView={{ width: "65%" }}
                      viewport={{ once: true }}
                      transition={{ duration: 1.5, ease: "easeOut" }}
                    />
                  </div>

                  {/* Contributors */}
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <Users className="w-4 h-4 text-muted-foreground" />
                      <span className="text-sm text-muted-foreground">
                        47 contributors
                      </span>
                    </div>
                    <span className="text-xs text-primary/70">
                      View details →
                    </span>
                  </div>
                </motion.div>
              </div>
            </FadeIn>

            {/* Floating notification - top left */}
            <motion.div
              className="absolute -top-6 -left-6 bg-card rounded-lg p-3 border border-primary/10 shadow-lg hidden lg:block"
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 1 }}
            >
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center">
                  <ShieldCheck className="w-4 h-4 text-primary" />
                </div>
                <div>
                  <div className="text-xs font-medium">Payment Verified</div>
                  <div className="text-[10px] text-muted-foreground">
                    ₦250,000 • 2 min ago
                  </div>
                </div>
              </div>
            </motion.div>

            {/* Floating stat - bottom right */}
            <motion.div
              className="absolute -bottom-4 -right-4 bg-card rounded-lg px-4 py-2 border border-primary/10 shadow-lg hidden lg:block"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 1.2 }}
            >
              <div className="flex items-center gap-2">
                <span className="text-xs text-muted-foreground">
                  3 milestones completed
                </span>
                <span className="text-primary">✓</span>
              </div>
            </motion.div>
          </div>
        </div>
      </div>
    </section>
  );
}
