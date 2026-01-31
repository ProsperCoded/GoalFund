import { motion } from "framer-motion"
import { useMemo } from "react"

interface SplitTextProps {
  text: string
  className?: string
  delay?: number
  animationFrom?: { opacity: number; y: number }
  animationTo?: { opacity: number; y: number }
  threshold?: number
  rootMargin?: string
  textAlign?: "left" | "center" | "right"
}

export function SplitText({
  text,
  className = "",
  delay = 0.05,
  animationFrom = { opacity: 0, y: 40 },
  animationTo = { opacity: 1, y: 0 },
  textAlign = "center",
}: SplitTextProps) {
  const words = useMemo(() => text.split(" "), [text])

  const containerVariants = {
    hidden: {},
    visible: {
      transition: {
        staggerChildren: delay,
      },
    },
  }

  const wordVariants = {
    hidden: animationFrom,
    visible: {
      ...animationTo,
      transition: {
        type: "spring" as const,
        damping: 12,
        stiffness: 100,
      },
    },
  }

  return (
    <motion.div
      className={`flex flex-wrap gap-x-2 gap-y-1 ${
        textAlign === "center" ? "justify-center" : textAlign === "right" ? "justify-end" : "justify-start"
      } ${className}`}
      variants={containerVariants}
      initial="hidden"
      whileInView="visible"
      viewport={{ once: true, amount: 0.3 }}
    >
      {words.map((word, index) => (
        <motion.span
          key={index}
          variants={wordVariants}
          className="inline-block"
        >
          {word}
        </motion.span>
      ))}
    </motion.div>
  )
}
