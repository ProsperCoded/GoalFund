import { motion } from "framer-motion"
import { Link } from "react-router-dom"
import { ArrowRight, Sparkles } from "lucide-react"
import { Button } from "@/components/ui/button"
import { SplitText, GradientText, FadeIn } from "@/components/animations"

export function CTASection() {
  return (
    <section className="py-24 relative overflow-hidden">
      {/* Background */}
      <div className="absolute inset-0 bg-gradient-to-br from-primary/10 via-transparent to-accent/10" />
      
      {/* Animated Orbs */}
      <motion.div
        className="absolute top-10 left-10 w-64 h-64 bg-primary/20 rounded-full blur-3xl"
        animate={{
          x: [0, 50, 0],
          y: [0, 30, 0],
        }}
        transition={{ duration: 10, repeat: Infinity }}
      />
      <motion.div
        className="absolute bottom-10 right-10 w-64 h-64 bg-accent/20 rounded-full blur-3xl"
        animate={{
          x: [0, -50, 0],
          y: [0, -30, 0],
        }}
        transition={{ duration: 10, repeat: Infinity, delay: 5 }}
      />

      <div className="relative z-10 max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <FadeIn>
          <motion.div
            whileHover={{ scale: 1.05 }}
            className="inline-flex items-center gap-2 px-4 py-2 rounded-full glass mb-6"
          >
            <Sparkles className="w-4 h-4 text-accent" />
            <span className="text-sm text-muted-foreground">
              Join thousands of successful fundraisers
            </span>
          </motion.div>
        </FadeIn>

        <h2 className="text-3xl sm:text-4xl lg:text-5xl font-bold mb-6">
          <SplitText text="Ready to Fund" className="mb-2" />
          <span className="block">
            <GradientText text="Your Dreams?" className="text-3xl sm:text-4xl lg:text-5xl font-bold" />
          </span>
        </h2>

        <FadeIn delay={0.2}>
          <p className="text-lg text-muted-foreground mb-8 max-w-2xl mx-auto">
            Whether it's a community project, personal milestone, or emergency fund,
            GoFund makes it easy to raise money with complete transparency.
          </p>
        </FadeIn>

        <FadeIn delay={0.3}>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link to="/register">
              <Button size="xl" variant="gradient" className="gap-2 w-full sm:w-auto">
                Create Your Goal Free
                <ArrowRight className="w-5 h-5" />
              </Button>
            </Link>
            <Link to="/goals">
              <Button size="xl" variant="outline" className="w-full sm:w-auto">
                Explore Active Goals
              </Button>
            </Link>
          </div>
        </FadeIn>

        <FadeIn delay={0.4}>
          <p className="mt-6 text-sm text-muted-foreground">
            No credit card required • Free forever • Cancel anytime
          </p>
        </FadeIn>
      </div>
    </section>
  )
}
