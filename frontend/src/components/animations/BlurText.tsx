import { motion } from "framer-motion"
import { useMemo } from "react"

interface BlurTextProps {
  text: string
  className?: string
  delay?: number
  animateByCharacter?: boolean
}

export function BlurText({
  text,
  className = "",
  delay = 0.03,
  animateByCharacter = false,
}: BlurTextProps) {
  const elements = useMemo(() => {
    if (animateByCharacter) {
      return text.split("").map((char, i) => (char === " " ? { char: "\u00A0", key: i } : { char, key: i }))
    }
    return text.split(" ").map((word, i) => ({ char: word, key: i }))
  }, [text, animateByCharacter])

  const containerVariants = {
    hidden: {},
    visible: {
      transition: {
        staggerChildren: delay,
      },
    },
  }

  const charVariants = {
    hidden: {
      opacity: 0,
      filter: "blur(10px)",
      y: 20,
    },
    visible: {
      opacity: 1,
      filter: "blur(0px)",
      y: 0,
      transition: {
        type: "spring" as const,
        damping: 15,
        stiffness: 150,
      },
    },
  }

  return (
    <motion.div
      className={`flex flex-wrap justify-center gap-x-2 ${className}`}
      variants={containerVariants}
      initial="hidden"
      whileInView="visible"
      viewport={{ once: true, amount: 0.5 }}
    >
      {elements.map(({ char, key }) => (
        <motion.span
          key={key}
          variants={charVariants}
          className="inline-block"
        >
          {char}
        </motion.span>
      ))}
    </motion.div>
  )
}
