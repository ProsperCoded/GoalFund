import { motion } from "framer-motion"
import type { ReactNode } from "react"

interface MagneticButtonProps {
  children: ReactNode
  className?: string
}

export function MagneticButton({
  children,
  className = "",
}: MagneticButtonProps) {
  return (
    <motion.div
      className={`inline-block ${className}`}
      whileHover={{ scale: 1.05 }}
      whileTap={{ scale: 0.95 }}
      transition={{
        type: "spring",
        stiffness: 400,
        damping: 17,
      }}
    >
      {children}
    </motion.div>
  )
}

interface TiltCardProps {
  children: ReactNode
  className?: string
  tiltAmount?: number
}

export function TiltCard({
  children,
  className = "",
  tiltAmount = 10,
}: TiltCardProps) {
  return (
    <motion.div
      className={className}
      whileHover={{
        rotateX: -tiltAmount,
        rotateY: tiltAmount,
        scale: 1.02,
      }}
      transition={{
        type: "spring",
        stiffness: 300,
        damping: 20,
      }}
      style={{
        transformStyle: "preserve-3d",
        perspective: "1000px",
      }}
    >
      {children}
    </motion.div>
  )
}

interface RotatingBorderProps {
  children: ReactNode
  className?: string
  borderWidth?: number
  duration?: number
}

export function RotatingBorder({
  children,
  className = "",
  borderWidth = 2,
  duration = 3,
}: RotatingBorderProps) {
  return (
    <div className={`relative p-[${borderWidth}px] ${className}`}>
      <motion.div
        className="absolute inset-0 rounded-xl"
        style={{
          background: "linear-gradient(90deg, #10b981, #f59e0b, #10b981)",
          backgroundSize: "200% 100%",
        }}
        animate={{
          backgroundPosition: ["0% 0%", "200% 0%"],
        }}
        transition={{
          duration,
          repeat: Infinity,
          ease: "linear",
        }}
      />
      <div className="relative bg-card rounded-xl">{children}</div>
    </div>
  )
}
