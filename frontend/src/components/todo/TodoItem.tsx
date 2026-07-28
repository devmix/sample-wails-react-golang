import { Check, Play, Pause, Trash2, Pencil } from 'lucide-react';
import type { Todo } from '@/types/models';
import PriorityBadge from '@/components/common/PriorityBadge';
import { TagList } from '@/components/common/TagChip';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

/** Props for the TodoItem component. */
interface TodoItemProps {
  /** The todo object to display. */
  todo: Todo;
  /** Callback to update a field of this todo. */
  onUpdate: (id: number, input: Partial<Todo>) => void;
  /** Callback to delete this todo. */
  onDelete: (id: number) => void;
  /** Callback to open the edit form for this todo. */
  onEdit: (todo: Todo) => void;
}

/**
 * Single todo row showing title, description, priority badge, due date and tags.
 * Action buttons (start/pause, edit, delete) appear on hover.
 * Overdue items are highlighted with a destructive border color.
 */
export default function TodoItem({ todo, onUpdate, onDelete, onEdit }: TodoItemProps) {
  const isOverdue = todo.dueDate && !todo.completedAt && new Date(todo.dueDate) < new Date(new Date().toDateString());

  return (
    <div className={cn(
      "group rounded-lg border p-4 transition-all hover:shadow-md",
      isOverdue ? 'border-destructive/50 bg-destructive/5' : 'bg-card'
    )}>
      <div className="flex items-start gap-3">
        <button
          onClick={() => onUpdate(todo.id, { status: todo.status === 'completed' ? 'pending' : 'completed' })}
          className={cn(
            "mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full border-2 transition-colors",
            todo.status === 'completed'
              ? 'border-green-500 bg-green-500 text-white'
              : todo.status === 'in_progress'
                ? 'border-yellow-400 bg-transparent'
                : 'border-input hover:border-primary'
          )}
        >
          {todo.status === 'completed' && <Check className="h-3 w-3" />}
          {todo.status === 'in_progress' && <div className="h-2 w-2 rounded-full bg-yellow-400" />}
        </button>

        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <h3 className={cn(
              "text-sm font-medium",
              todo.status === 'completed' ? 'text-muted-foreground line-through' : isOverdue ? 'text-destructive' : ''
            )}>
              {todo.title}
            </h3>
            <PriorityBadge priority={todo.priority as string} />
          </div>

          {todo.description && (
            <p className={cn(
              "mt-1 text-xs line-clamp-2",
              todo.status === 'completed' ? 'text-muted-foreground/60 line-through' : 'text-muted-foreground'
            )}>
              {todo.description}
            </p>
          )}

          <div className="mt-2 flex items-center gap-3">
            {todo.dueDate && (
              <span className={cn(
                "text-xs",
                isOverdue ? 'font-medium text-destructive' : 'text-muted-foreground'
              )}>
                Due: {todo.dueDate}
              </span>
            )}
            {todo.tags && todo.tags.length > 0 && <TagList tags={todo.tags} />}
          </div>
        </div>

        <div className="flex shrink-0 gap-1 opacity-0 transition-opacity group-hover:opacity-100">
          {todo.status !== 'in_progress' && (
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7"
              onClick={() => onUpdate(todo.id, { status: todo.status === 'pending' ? 'in_progress' : 'pending' })}
              title={todo.status === 'pending' ? 'Start' : 'Pause'}
            >
              {todo.status === 'pending' ? <Play className="h-3.5 w-3.5" /> : <Pause className="h-3.5 w-3.5" />}
            </Button>
          )}
          <Button
            variant="ghost"
            size="icon"
            className="h-7 w-7"
            onClick={() => onEdit(todo)}
            title="Edit"
          >
            <Pencil className="h-3.5 w-3.5" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="h-7 w-7 text-destructive hover:text-destructive"
            onClick={() => onDelete(todo.id)}
            title="Delete"
          >
            <Trash2 className="h-3.5 w-3.5" />
          </Button>
        </div>
      </div>
    </div>
  );
}
