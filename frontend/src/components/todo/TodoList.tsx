import type { Todo, TodoFilter } from '@/types/models';
import TodoItem from './TodoItem';
import TodoFilters from './TodoFilters';

/** Props for the TodoList component. */
interface TodoListProps {
  /** Full list of todos to render and filter. */
  todos: Todo[];
  /** Current active filters applied to the todo list. */
  filter: TodoFilter;
  /** Callback to update the active filter state. */
  onFilterChange: (f: TodoFilter) => void;
  /** Callback to update a specific todo's fields. */
  onUpdate: (id: number, input: Partial<Todo>) => void;
  /** Callback to delete a todo by its id. */
  onDelete: (id: number) => void;
  /** Callback to open the edit form for a given todo. */
  onEdit: (todo: Todo) => void;
  /** Optional list of tags available for filtering. Defaults to an empty array. */
  availableTags?: { id: number; name: string }[];
}

/**
 * Container that renders filter controls and a filtered list of todo items.
 * Filtering is done client-side based on status, priority and tag name.
 */
export default function TodoList({ todos, filter, onFilterChange, onUpdate, onDelete, onEdit, availableTags = [] }: TodoListProps) {
  const filtered = todos.filter(todo => {
    if (filter.status && todo.status !== filter.status) return false;
    if (filter.priority && todo.priority !== filter.priority) return false;
    if (filter.tagName) {
      const hasTag = todo.tags?.some(t => t.name === filter.tagName);
      if (!hasTag) return false;
    }
    return true;
  });

  return (
    <div>
      <TodoFilters filter={filter} onFilterChange={onFilterChange} availableTags={availableTags} />

      <div className="mt-4 space-y-2">
        {filtered.length === 0 ? (
          <p className="py-8 text-center text-sm text-muted-foreground">No todos found</p>
        ) : (
          filtered.map(todo => (
            <TodoItem key={todo.id} todo={todo} onUpdate={onUpdate} onDelete={onDelete} onEdit={onEdit} />
          ))
        )}
      </div>
    </div>
  );
}
