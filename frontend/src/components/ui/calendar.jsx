import * as React from "react"
import { cn } from "../../lib/utils"

function Calendar({ className, ...props }) {
  return (
    <div
      className={cn("p-3 rounded-2xl bg-glass/60 backdrop-blur-xl border border-glass-border/40 shadow-glass", className)}
      {...props}
    />
  )
}

function CalendarHeader({ className, ...props }) {
  return (
    <div
      className={cn("flex items-center justify-between mb-4", className)}
      {...props}
    />
  )
}

function CalendarTitle({ className, ...props }) {
  return (
    <div
      className={cn("text-sm font-medium", className)}
      {...props}
    />
  )
}

function CalendarPrevButton({ className, ...props }) {
  return (
    <button
      className={cn(
        "h-7 w-7 flex items-center justify-center rounded-lg hover:bg-glass/60 transition-gentle",
        className
      )}
      {...props}
    >
      ‹
    </button>
  )
}

function CalendarNextButton({ className, ...props }) {
  return (
    <button
      className={cn(
        "h-7 w-7 flex items-center justify-center rounded-lg hover:bg-glass/60 transition-gentle",
        className
      )}
      {...props}
    >
      ›
    </button>
  )
}

function CalendarGrid({ className, ...props }) {
  return (
    <div
      className={cn("grid grid-cols-7 gap-1", className)}
      {...props}
    />
  )
}

function CalendarDayHeader({ className, ...props }) {
  return (
    <div
      className={cn("text-xs text-muted-foreground text-center py-1", className)}
      {...props}
    />
  )
}

function CalendarDay({ className, ...props }) {
  return (
    <button
      className={cn(
        "h-8 w-8 text-xs rounded-lg hover:bg-glass/60 transition-gentle",
        "focus:outline-none focus:ring-2 focus:ring-ring focus:ring-primary/30",
        className
      )}
      {...props}
    />
  )
}

export {
  Calendar,
  CalendarHeader,
  CalendarTitle,
  CalendarPrevButton,
  CalendarNextButton,
  CalendarGrid,
  CalendarDayHeader,
  CalendarDay,
}
