import * as React from 'react';
import { cva } from 'class-variance-authority';
import { cn } from '@/lib/utils';

const buttonVariants = cva(
  'inline-flex items-center justify-center whitespace-nowrap rounded-xl text-base font-medium ring-offset-background transition-gentle focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 relative overflow-hidden',
  {
    variants: {
      variant: {
        default: 'liquid-button bg-primary text-primary-foreground shadow-glow',
        destructive: 'liquid-button bg-destructive text-destructive-foreground shadow-glow',
        outline: 'liquid-button border border-input bg-glass/60 backdrop-blur-2xl text-foreground shadow-glass-sm',
        secondary: 'liquid-button bg-glass/70 backdrop-blur-2xl text-secondary-foreground shadow-glass-sm',
        ghost: 'text-foreground',
        link: 'text-primary underline-offset-4 relative overflow-hidden',
      },
      size: {
        default: 'h-11 px-5 py-2.5',
        sm: 'h-10 rounded-lg px-4',
        lg: 'h-12 rounded-lg px-10',
        icon: 'h-11 w-11',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  }
);

const Button = React.forwardRef(({ className, variant, size, ...props }, ref) => {
  return (
    <button
      className={cn(buttonVariants({ variant, size, className }))}
      ref={ref}
      {...props}
    />
  );
});
Button.displayName = 'Button';

export { Button, buttonVariants };
