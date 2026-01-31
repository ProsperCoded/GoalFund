import { motion, useInView } from "framer-motion"
import { useRef } from "react"

interface CountUpProps {
  end: number
  duration?: number
  prefix?: string
  suffix?: string
  className?: string
}

export function CountUp({ 
  end, 
  duration = 2, 
  prefix = "", 
  suffix = "", 
  className = "" 
}: CountUpProps) {
  const ref = useRef(null)
  const isInView = useInView(ref, { once: true })

  return (
    <span ref={ref} className={className}>
      {prefix}
      <motion.span
        initial={{ opacity: 0 }}
        animate={isInView ? { opacity: 1 } : {}}
      >
        {isInView && (
          <motion.span
            initial={0}
            animate={end}
            transition={{ duration, ease: "easeOut" }}
          >
            {Math.round(end).toLocaleString()}
          </motion.span>
        )}
        {!isInView && "0"}
      </motion.span>
      {suffix}
    </span>
  )
}
