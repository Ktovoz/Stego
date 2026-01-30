import * as React from 'react';
import { cn } from '@/lib/utils';

const Switch = React.forwardRef(({ checked, onCheckedChange, ...props }, ref) => {
  const [internalChecked, setInternalChecked] = React.useState(checked || false);
  const currentValue = checked !== undefined ? checked : internalChecked;

  const handleChange = () => {
    const newValue = !currentValue;
    if (checked === undefined) {
      setInternalChecked(newValue);
    }
    onCheckedChange?.(newValue);
  };

  return (
    <button
      type="button"
      role="switch"
      aria-checked={currentValue}
      className={cn(
        'peer inline-flex h-6 w-11 shrink-0 cursor-pointer items-center rounded-full border-2 border-transparent focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:cursor-not-allowed disabled:opacity-50',
        currentValue ? 'bg-primary' : 'bg-glass/70 backdrop-blur-sm'
      )}
      ref={ref}
      onClick={handleChange}
      {...props}
    >
      <span
        className={cn(
          'pointer-events-none block h-5 w-5 rounded-full bg-background shadow-lg ring-0 transition-transform',
          currentValue ? 'translate-x-5' : 'translate-x-0'
        )}
      />
    </button>
  );
});
Switch.displayName = 'Switch';

export { Switch };
