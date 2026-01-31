import { motion } from "framer-motion"

interface CounterProps {
  value: number
  className?: string
  formatValue?: (value: number) => string
}

export function Counter({
  value,
  className = "",
  formatValue = (v) => v.toLocaleString(),
}: CounterProps) {
  return (
    <motion.span
      className={className}
      initial={{ opacity: 0 }}
      whileInView={{ opacity: 1 }}
      viewport={{ once: true }}
    >
      <motion.span
        initial={{ opacity: 0 }}
        whileInView={{
          opacity: 1,
        }}
        viewport={{ once: true }}
      >
        {formatValue(value)}
      </motion.span>
    </motion.span>
  )
}

interface AnimatedCounterProps {
  to: number
  className?: string
  duration?: number
  formatValue?: (value: number) => string
}

export function AnimatedCounter({
  to,
  className = "",
  duration = 2,
  formatValue = (v) => Math.round(v).toLocaleString(),
}: AnimatedCounterProps) {
  return (
    <motion.span
      className={className}
      initial={{ opacity: 0 }}
      whileInView={{ opacity: 1 }}
      viewport={{ once: true }}
    >
      <motion.span
        initial={{ 
          filter: "blur(4px)",
        }}
        whileInView={{ 
          filter: "blur(0px)",
        }}
        transition={{ duration: duration * 0.5 }}
        viewport={{ once: true }}
      >
        {formatValue(to)}
      </motion.span>
    </motion.span>
  )
}
