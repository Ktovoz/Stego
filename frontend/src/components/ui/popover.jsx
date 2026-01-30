import * as React from "react"
import { cn } from "../../lib/utils"

const PopoverContext = React.createContext({
  isOpen: false,
  setIsOpen: () => {},
})

const Popover = React.forwardRef(({ className, children, ...props }, ref) => {
  const [isOpen, setIsOpen] = React.useState(false)

  return (
    <PopoverContext.Provider value={{ isOpen, setIsOpen }}>
      <div
        ref={ref}
        className={cn("relative", className)}
        {...props}
      >
        {children}
      </div>
    </PopoverContext.Provider>
  )
})
Popover.displayName = "Popover"

const PopoverTrigger = React.forwardRef(({ className, children, ...props }, ref) => {
  const { isOpen, setIsOpen } = React.useContext(PopoverContext)

  return (
    <button
      ref={ref}
      type="button"
      onClick={() => setIsOpen(!isOpen)}
      className={cn("w-full", className)}
      {...props}
    >
      {children}
    </button>
  )
})
PopoverTrigger.displayName = "PopoverTrigger"

const PopoverContent = React.forwardRef(({ className, children, ...props }, ref) => {
  const { isOpen, setIsOpen } = React.useContext(PopoverContext)
  const contentRef = React.useRef(null)

  React.useEffect(() => {
    if (!isOpen) return

    const handleClick = (e) => {
      if (contentRef.current && !contentRef.current.contains(e.target)) {
        setIsOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClick)
    return () => document.removeEventListener('mousedown', handleClick)
  }, [isOpen, setIsOpen])

  if (!isOpen) return null

  return (
    <div
      ref={(node) => {
        contentRef.current = node
        if (typeof ref === 'function') ref(node)
        else if (ref) ref.current = node
      }}
      className={cn(
        "absolute z-50 w-72 rounded-2xl liquid-panel p-4 text-popover-foreground outline-none transition-gentle",
        "top-full left-0 mt-1",
        className
      )}
      {...props}
    >
      {children}
    </div>
  )
})
PopoverContent.displayName = "PopoverContent"

export { Popover, PopoverTrigger, PopoverContent }
