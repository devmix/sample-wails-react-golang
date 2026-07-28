import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';

const priorityStyles: Record<string, string> = {
  low: 'bg-priority-low/15 text-priority-low border-priority-low/30',
  medium: 'bg-priority-medium/15 text-priority-medium border-priority-medium/30',
  high: 'bg-priority-high/15 text-priority-high border-priority-high/30',
  urgent: 'bg-priority-urgent/15 text-priority-urgent border-priority-urgent/30',
};

/** Props for the PriorityBadge component. */
interface PriorityBadgeProps {
  /** Priority level key that determines badge color (low, medium, high, urgent). */
  priority: string;
}

/**
 * Displays a colored badge representing a task's priority level.
 * Color is mapped from the `priority` value via an internal style map.
 */
export default function PriorityBadge({ priority }: PriorityBadgeProps) {
  return (
    <Badge variant="outline" className={cn('rounded-full font-normal capitalize', priorityStyles[priority] || '')}>
      {priority}
    </Badge>
  );
}
