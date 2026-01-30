import * as React from "react"
import { cn } from "../../lib/utils"
import { Button } from "./button"

function DatePicker({
  value,
  onChange,
  className,
  placeholder = "选择日期",
  ...props
}) {
  const [isOpen, setIsOpen] = React.useState(false)
  const [currentMonth, setCurrentMonth] = React.useState(new Date())
  const containerRef = React.useRef(null)

  const selectedDate = value ? new Date(value) : null

  const daysInMonth = new Date(
    currentMonth.getFullYear(),
    currentMonth.getMonth() + 1,
    0
  ).getDate()

  const firstDayOfMonth = new Date(
    currentMonth.getFullYear(),
    currentMonth.getMonth(),
    1
  ).getDay()

  const weekDays = ['日', '一', '二', '三', '四', '五', '六']

  const handlePrevMonth = () => {
    setCurrentMonth(new Date(currentMonth.getFullYear(), currentMonth.getMonth() - 1))
  }

  const handleNextMonth = () => {
    setCurrentMonth(new Date(currentMonth.getFullYear(), currentMonth.getMonth() + 1))
  }

  const handleDateClick = (day) => {
    const year = currentMonth.getFullYear()
    const month = currentMonth.getMonth()

    const yyyy = year
    const mm = String(month + 1).padStart(2, '0')
    const dd = String(day).padStart(2, '0')
    const isoString = `${yyyy}-${mm}-${dd}`

    onChange(isoString)
    setIsOpen(false)
  }

  const formatDate = (date) => {
    if (!date) return ''
    const d = new Date(date)
    return d.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    })
  }

  const isSelected = (day) => {
    if (!selectedDate) return false
    return (
      day === selectedDate.getDate() &&
      currentMonth.getMonth() === selectedDate.getMonth() &&
      currentMonth.getFullYear() === selectedDate.getFullYear()
    )
  }

  const isToday = (day) => {
    const today = new Date()
    return (
      day === today.getDate() &&
      currentMonth.getMonth() === today.getMonth() &&
      currentMonth.getFullYear() === today.getFullYear()
    )
  }

  React.useEffect(() => {
    if (!isOpen) return

    const handleClick = (e) => {
      if (containerRef.current && !containerRef.current.contains(e.target)) {
        setIsOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClick)
    return () => document.removeEventListener('mousedown', handleClick)
  }, [isOpen])

  return (
    <div ref={containerRef} className={cn("relative", className)} {...props}>
      <Button
        variant="outline"
        className={cn(
          "w-full justify-start text-left font-normal h-9 text-sm",
          !value && "text-muted-foreground"
        )}
        onClick={() => setIsOpen(!isOpen)}
      >
        {value ? formatDate(value) : placeholder}
      </Button>

      {isOpen && (
        <div className="absolute z-50 w-72 rounded-md border bg-popover p-4 text-popover-foreground shadow-md outline-none top-full left-0 mt-1">
          <div className="flex items-center justify-between mb-4">
            <button
              onClick={handlePrevMonth}
              className="h-7 w-7 flex items-center justify-center rounded-md hover:bg-accent"
            >
              ‹
            </button>
            <div className="text-sm font-medium">
              {currentMonth.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long' })}
            </div>
            <button
              onClick={handleNextMonth}
              className="h-7 w-7 flex items-center justify-center rounded-md hover:bg-accent"
            >
              ›
            </button>
          </div>

          <div className="grid grid-cols-7 gap-1">
            {weekDays.map((day) => (
              <div key={day} className="text-xs text-muted-foreground text-center py-1">
                {day}
              </div>
            ))}
            {Array.from({ length: firstDayOfMonth }).map((_, i) => (
              <div key={`empty-${i}`} />
            ))}
            {Array.from({ length: daysInMonth }).map((_, i) => {
              const day = i + 1
              return (
                <button
                  key={day}
                  onClick={() => handleDateClick(day)}
                  className={cn(
                    "h-8 w-8 text-xs rounded-md hover:bg-accent focus:outline-none focus:ring-2 focus:ring-ring",
                    isSelected(day) && "bg-primary text-primary-foreground",
                    isToday(day) && !isSelected(day) && "border border-primary"
                  )}
                >
                  {day}
                </button>
              )
            })}
          </div>
        </div>
      )}
    </div>
  )
}

export { DatePicker }
