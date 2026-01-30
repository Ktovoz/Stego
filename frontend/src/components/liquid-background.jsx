import { cn } from '@/lib/utils';

export function LiquidBackground({ className }) {
  return (
    <div className={cn('liquid-bg', className)} aria-hidden="true">
      <div className="liquid-blob liquid-blob-1" />
      <div className="liquid-blob liquid-blob-2" />
      <div className="liquid-blob liquid-blob-3" />
      <div className="liquid-bg-noise" />
    </div>
  );
}
