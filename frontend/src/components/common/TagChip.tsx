import { Badge } from '@/components/ui/badge';
import type { Tag } from '@/types/models';

/** Props for the TagChip component. */
interface TagChipProps {
  /** Display label of the tag. */
  name: string;
  /** Optional callback to remove this tag from its parent collection. */
  onRemove?: () => void;
}

/**
 * Renders a single removable tag badge.
 * Shows an "×" button only when `onRemove` is provided.
 */
export function TagChip({ name, onRemove }: TagChipProps) {
  return (
    <Badge variant="outline" className="rounded-full bg-tag-bg text-tag-fg border-tag-bg/30 font-normal gap-1 pr-1">
      {name}
      {onRemove && (
        <button
          onClick={onRemove}
          className="ml-0.5 rounded-full hover:bg-tag-fg/20 w-4 h-4 flex items-center justify-center"
          type="button"
        >
          ×
        </button>
      )}
    </Badge>
  );
}

/** Props for the TagList component. */
interface TagListProps {
  /** Array of tag objects to render as chips. */
  tags: Tag[];
}

/**
 * Renders a flex-wrapped row of `TagChip` components.
 * Returns null when the tags array is empty or undefined.
 */
export function TagList({ tags }: TagListProps) {
  if (!tags || tags.length === 0) return null;
  return (
    <div className="flex flex-wrap gap-1">
      {tags.map(tag => (
        <TagChip key={tag.id} name={tag.name} />
      ))}
    </div>
  );
}
