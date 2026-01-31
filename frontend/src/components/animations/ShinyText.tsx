import { motion } from "framer-motion"

interface ShinyTextProps {
  text: string
  className?: string
  shimmerWidth?: number
}

export function ShinyText({
  text,
  className = "",
  shimmerWidth = 100,
}: ShinyTextProps) {
  return (
    <motion.span
      className={`relative inline-block ${className}`}
      style={{
        background: `linear-gradient(
          120deg,
          rgba(255, 255, 255, 0) 40%,
          rgba(255, 255, 255, 0.8) 50%,
          rgba(255, 255, 255, 0) 60%
        ) var(--color-foreground)`,
        backgroundSize: `${shimmerWidth}% 100%`,
        WebkitBackgroundClip: "text",
        WebkitTextFillColor: "transparent",
        backgroundClip: "text",
      }}
      animate={{
        backgroundPosition: ["200% center", "-200% center"],
      }}
      transition={{
        duration: 3,
        ease: "linear",
        repeat: Infinity,
        repeatDelay: 1,
      }}
    >
      {text}
    </motion.span>
  )
}
