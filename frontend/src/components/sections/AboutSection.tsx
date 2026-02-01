import { motion } from "framer-motion";
import {
  BookOpen,
  Eye,
  ShieldCheck,
  Banknote,
  Target,
  Users,
  ArrowRight,
  ChevronRight,
  CircleDot,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";
import {
  GradientText,
  FadeIn,
  StaggerContainer,
  StaggerItem,
  TiltCard,
} from "@/components/animations";

export function AboutSection() {
  const coreFeatures = [
    {
      icon: BookOpen,
      title: "Ledger-Backed Accounting",
      description:
        "Every contribution and withdrawal is recorded as immutable double-entry ledger entries. Balances are computed, never stored—ensuring financial accuracy.",
    },
    {
      icon: Eye,
      title: "Complete Transparency",
      description:
        "Contributors see exactly where their money goes. Track every transaction, view withdrawal history, and monitor fund utilization in real-time.",
    },
    {
      icon: ShieldCheck,
      title: "Verified Payments",
      description:
        "All payments go through Paystack with webhook verification. Idempotent processing ensures no duplicate charges or missed contributions.",
    },
    {
      icon: Banknote,
      title: "Flexible Withdrawals",
      description:
        "Goal owners can withdraw available funds anytime. Support for multiple withdrawals with complete audit trail. Bank details securely stored.",
    },
    {
      icon: Target,
      title: "Milestone Tracking",
      description:
        "Break down large goals into milestones. Support for recurring contributions—weekly, monthly, or yearly. Track progress at each stage.",
    },
    {
      icon: Users,
      title: "Community Accountability",
      description:
        "Contributors can vote on proof of accomplishment. Build trust through transparent fund usage and community feedback.",
    },
  ];

  const howItWorks = [
    {
      step: "01",
      title: "Create a Goal",
      description:
        "Set your funding target, add milestones, and describe how funds will be used.",
    },
    {
      step: "02",
      title: "Share & Collect",
      description:
        "Share your goal link. Contributors can fund with just an email—no signup required.",
    },
    {
      step: "03",
      title: "Track Progress",
      description:
        "Watch contributions come in. Every payment is verified and recorded in real-time.",
    },
    {
      step: "04",
      title: "Withdraw & Report",
      description:
        "Withdraw funds when ready. Submit proof of accomplishment for accountability.",
    },
  ];

  return (
    <section className="py-24 relative overflow-hidden">
      {/* Decorative line */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-px h-24 bg-linear-to-b from-transparent to-primary/30" />

      <div className="relative z-10 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section Header */}
        <div className="text-center mb-20">
          <FadeIn>
            <div className="inline-flex items-center gap-2 mb-6">
              <div className="w-8 h-px bg-primary/60" />
              <span className="text-xs uppercase tracking-[0.2em] text-primary font-medium">
                Built for Accountability
              </span>
              <div className="w-8 h-px bg-primary/60" />
            </div>
          </FadeIn>

          <FadeIn delay={0.1}>
            <h2 className="text-3xl sm:text-4xl lg:text-5xl font-bold mb-6 leading-tight">
              Financial Transparency{" "}
              <span className="block sm:inline">
                for{" "}
                <GradientText
                  text="Large Goals"
                  className="text-3xl sm:text-4xl lg:text-5xl font-bold"
                />
              </span>
            </h2>
          </FadeIn>

          <FadeIn delay={0.2}>
            <p className="text-lg text-muted-foreground max-w-3xl mx-auto leading-relaxed">
              GoalFund is built for community projects, group contributions, and
              large-scale funding where trust and accountability matter. Every
              naira is tracked with fintech-grade precision.
            </p>
          </FadeIn>
        </div>

        {/* Core Features - Bento Grid Style */}
        <div className="mb-24">
          <FadeIn delay={0.1}>
            <h3 className="text-sm uppercase tracking-wider text-primary/70 mb-8 text-center">
              Core Features
            </h3>
          </FadeIn>

          <StaggerContainer
            className="grid md:grid-cols-2 lg:grid-cols-3 gap-5"
            staggerDelay={0.08}
          >
            {coreFeatures.map((feature) => {
              const Icon = feature.icon;
              return (
                <StaggerItem key={feature.title}>
                  <TiltCard>
                    <motion.div
                      className="h-full bg-card rounded-xl p-6 border border-border/50 card-hover group"
                      whileHover={{ borderColor: "rgba(212, 168, 83, 0.3)" }}
                    >
                      <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center mb-4 group-hover:bg-primary/20 transition-colors">
                        <Icon className="w-5 h-5 text-primary" />
                      </div>
                      <h4 className="text-lg font-semibold mb-2">
                        {feature.title}
                      </h4>
                      <p className="text-sm text-muted-foreground leading-relaxed">
                        {feature.description}
                      </p>
                    </motion.div>
                  </TiltCard>
                </StaggerItem>
              );
            })}
          </StaggerContainer>
        </div>

        {/* How It Works */}
        <div className="mb-24">
          <FadeIn>
            <div className="text-center mb-12">
              <h3 className="text-2xl sm:text-3xl font-bold mb-4">
                How <GradientText text="GoalFund" /> Works
              </h3>
              <p className="text-muted-foreground max-w-xl mx-auto">
                Simple process, powerful accountability. Get your project funded
                in four steps.
              </p>
            </div>
          </FadeIn>

          <div className="grid md:grid-cols-4 gap-6 relative">
            {/* Connection line - desktop only */}
            <div className="absolute top-8 left-[12.5%] right-[12.5%] h-px bg-border hidden md:block">
              <motion.div
                className="h-full bg-linear-to-r from-primary/50 via-primary to-primary/50"
                initial={{ scaleX: 0 }}
                whileInView={{ scaleX: 1 }}
                viewport={{ once: true }}
                transition={{ duration: 1.5, delay: 0.5 }}
              />
            </div>

            {howItWorks.map((item, index) => (
              <FadeIn key={item.step} delay={0.1 + index * 0.15}>
                <div className="relative text-center">
                  {/* Step indicator */}
                  <div className="relative z-10 w-16 h-16 mx-auto mb-5 rounded-full bg-card border-2 border-primary/30 flex items-center justify-center">
                    <span className="text-sm font-bold text-primary">
                      {item.step}
                    </span>
                  </div>

                  <h4 className="font-semibold mb-2">{item.title}</h4>
                  <p className="text-sm text-muted-foreground">
                    {item.description}
                  </p>
                </div>
              </FadeIn>
            ))}
          </div>
        </div>

        {/* Use Cases */}
        <FadeIn>
          <div className="bg-card rounded-2xl border border-border/50 p-8 md:p-12">
            <div className="grid lg:grid-cols-2 gap-12 items-center">
              {/* Left - Text */}
              <div>
                <h3 className="text-2xl sm:text-3xl font-bold mb-4">
                  Perfect for <GradientText text="Large-Scale" /> Projects
                </h3>
                <p className="text-muted-foreground mb-6">
                  GoalFund is designed for projects where accountability and
                  transparency matter most. From community infrastructure to
                  educational initiatives.
                </p>

                <ul className="space-y-3 mb-8">
                  {[
                    "Community borehole and infrastructure projects",
                    "School fees and educational funding",
                    "Religious and community center buildings",
                    "Medical emergency and health funds",
                    "Group investment and business funding",
                    "Event and celebration contributions",
                  ].map((item, i) => (
                    <motion.li
                      key={i}
                      initial={{ opacity: 0, x: -10 }}
                      whileInView={{ opacity: 1, x: 0 }}
                      viewport={{ once: true }}
                      transition={{ delay: i * 0.05 }}
                      className="flex items-start gap-3 text-sm"
                    >
                      <CircleDot className="w-4 h-4 text-primary mt-0.5 flex-shrink-0" />
                      <span>{item}</span>
                    </motion.li>
                  ))}
                </ul>

                <div className="flex flex-col sm:flex-row gap-3">
                  <Link to="/register">
                    <Button
                      size="lg"
                      className="bg-primary text-primary-foreground hover:bg-primary/90 gap-2"
                    >
                      Start a Goal
                      <ArrowRight className="w-4 h-4" />
                    </Button>
                  </Link>
                  <Button
                    variant="outline"
                    size="lg"
                    className="gap-2 border-primary/30"
                  >
                    Learn More
                    <ChevronRight className="w-4 h-4" />
                  </Button>
                </div>
              </div>

              {/* Right - Stats */}
              <div className="relative">
                <div className="relative flex items-center justify-center">
                  <div className="absolute inset-0 bg-primary/10 blur-3xl rounded-full opacity-30" />
                  <img
                    src="/assets/secure-book.png"
                    alt="Secure Ledger"
                    className="relative w-full max-w-[400px] object-contain drop-shadow-xl hover:scale-105 transition-transform duration-500"
                  />
                </div>
              </div>

              {/* Accent */}
              <div className="absolute -top-2 -right-2 w-4 h-4 border-t-2 border-r-2 border-primary/30" />
              <div className="absolute -bottom-2 -left-2 w-4 h-4 border-b-2 border-l-2 border-primary/30" />
            </div>
          </div>
        </FadeIn>
      </div>
    </section>
  );
}
