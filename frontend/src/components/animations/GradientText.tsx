import { motion } from "framer-motion"

interface GradientTextProps {
  text: string
  className?: string
  from?: string
  via?: string
  to?: string
  animate?: boolean
}

export function GradientText({
  text,
  className = "",
  from = "#f5d78e",
  via = "#d4a853",
  to = "#b8942e",
  animate = false,
}: GradientTextProps) {
  return (
    <motion.span
      className={`inline-block ${className}`}
      style={{
        background: `linear-gradient(135deg, ${from} 0%, ${via} 50%, ${to} 100%)`,
        backgroundSize: animate ? "200% auto" : "100% auto",
        WebkitBackgroundClip: "text",
        WebkitTextFillColor: "transparent",
        backgroundClip: "text",
      }}
      animate={animate ? {
        backgroundPosition: ["0% center", "200% center", "0% center"],
      } : undefined}
      transition={animate ? {
        duration: 4,
        ease: "linear",
        repeat: Infinity,
      } : undefined}
    >
      {text}
    </motion.span>
  )
}
