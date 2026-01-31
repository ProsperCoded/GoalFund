import { useEffect, useState, useCallback } from "react"
import { motion } from "framer-motion"

interface TextScrambleProps {
  text: string
  className?: string
  speed?: number
  trigger?: boolean
}

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*"

export function TextScramble({ 
  text, 
  className = "", 
  speed = 50,
  trigger = true 
}: TextScrambleProps) {
  const [displayText, setDisplayText] = useState(text)
  const [isAnimating, setIsAnimating] = useState(false)

  const scramble = useCallback(() => {
    if (isAnimating) return
    setIsAnimating(true)
    
    let iteration = 0
    const interval = setInterval(() => {
      setDisplayText(
        text
          .split("")
          .map((char, index) => {
            if (char === " ") return " "
            if (index < iteration) return text[index]
            return chars[Math.floor(Math.random() * chars.length)]
          })
          .join("")
      )

      if (iteration >= text.length) {
        clearInterval(interval)
        setIsAnimating(false)
      }

      iteration += 1 / 3
    }, speed)

    return () => clearInterval(interval)
  }, [text, speed, isAnimating])

  useEffect(() => {
    if (trigger) {
      scramble()
    }
  }, [trigger, scramble])

  return (
    <motion.span
      className={`inline-block font-mono ${className}`}
      onHoverStart={scramble}
    >
      {displayText}
    </motion.span>
  )
}
