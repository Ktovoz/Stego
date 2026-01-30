import * as React from "react"
import { cn } from "../../lib/utils"

function Badge({
  className,
  variant = "default",
  ...props
}) {
  const variants = {
    default: "bg-primary/80 text-primary-foreground backdrop-blur-sm",
    secondary: "bg-glass/70 text-secondary-foreground backdrop-blur-sm border border-glass-border/40",
    destructive: "bg-destructive/80 text-destructive-foreground backdrop-blur-sm",
    outline: "text-foreground border border-input bg-glass/40 backdrop-blur-sm",
  }

  return (
    <div
      className={cn(
        "inline-flex items-center rounded-lg border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
        variants[variant],
        className
      )}
      {...props}
    />
  )
}

export { Badge }
