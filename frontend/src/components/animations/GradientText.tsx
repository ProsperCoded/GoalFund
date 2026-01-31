import { motion } from "framer-motion"

interface GradientTextProps {
  text: string
  className?: string
  from?: string
  via?: string
  to?: string
}

export function GradientText({
  text,
  className = "",
  from = "#10b981",
  via = "#f59e0b",
  to = "#10b981",
}: GradientTextProps) {
  return (
    <motion.span
      className={`inline-block ${className}`}
      style={{
        background: `linear-gradient(135deg, ${from} 0%, ${via} 50%, ${to} 100%)`,
        backgroundSize: "200% auto",
        WebkitBackgroundClip: "text",
        WebkitTextFillColor: "transparent",
        backgroundClip: "text",
      }}
      animate={{
        backgroundPosition: ["0% center", "200% center", "0% center"],
      }}
      transition={{
        duration: 4,
        ease: "linear",
        repeat: Infinity,
      }}
    >
      {text}
    </motion.span>
  )
}
