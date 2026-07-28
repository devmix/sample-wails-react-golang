import { X } from 'lucide-react';
import type { TodoFilter } from '@/types/models';
import { Button } from '@/components/ui/button';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';

/** Props for the TodoFilters component. */
interface TodoFiltersProps {
  /** Current filter state containing status, priority and tag selections. */
  filter: TodoFilter;
  /** Callback to apply a new filter configuration. */
  onFilterChange: (f: TodoFilter) => void;
  /** Tags available for filtering display as toggle buttons. */
  availableTags: { id: number; name: string }[];
}

/**
 * Filter bar with dropdowns for status and priority, tag toggle buttons,
 * and a "Clear filters" button that appears when any filter is active.
 */
export default function TodoFilters({ filter, onFilterChange, availableTags }: TodoFiltersProps) {
  return (
    <div className="flex flex-wrap items-center gap-3">
      <Select value={filter.status || ''} onValueChange={v => onFilterChange({ ...filter, status: v || undefined })}>
        <SelectTrigger className="w-[160px]">
          <SelectValue placeholder="All Status" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="">All Status</SelectItem>
          <SelectItem value="pending">Pending</SelectItem>
          <SelectItem value="in_progress">In Progress</SelectItem>
          <SelectItem value="completed">Completed</SelectItem>
        </SelectContent>
      </Select>

      <Select value={filter.priority || ''} onValueChange={v => onFilterChange({ ...filter, priority: v || undefined })}>
        <SelectTrigger className="w-[160px]">
          <SelectValue placeholder="All Priority" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="">All Priority</SelectItem>
          <SelectItem value="low">Low</SelectItem>
          <SelectItem value="medium">Medium</SelectItem>
          <SelectItem value="high">High</SelectItem>
          <SelectItem value="urgent">Urgent</SelectItem>
        </SelectContent>
      </Select>

      {availableTags.length > 0 && (
        <div className="flex flex-wrap gap-1">
          {availableTags.map(tag => (
            <Button
              key={tag.id}
              variant={filter.tagName === tag.name ? 'default' : 'secondary'}
              size="sm"
              className="rounded-full px-3 text-xs font-normal"
              onClick={() => onFilterChange({ ...filter, tagName: filter.tagName === tag.name ? undefined : tag.name })}
            >
              {tag.name}
            </Button>
          ))}
        </div>
      )}

      {(filter.status || filter.priority || filter.tagName) && (
        <Button variant="ghost" size="sm" className="text-xs font-normal h-auto px-2" onClick={() => onFilterChange({})}>
          <X className="h-3 w-3 mr-1" />
          Clear filters
        </Button>
      )}
    </div>
  );
}
