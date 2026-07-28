import { useState } from 'react';
import { Plus } from 'lucide-react';
import { useTodos } from '@/hooks/useTodos';
import { useTags } from '@/hooks/useTags';
import TodoForm from '@/components/todo/TodoForm';
import TodoList from '@/components/todo/TodoList';
import SearchBar from '@/components/common/SearchBar';
import { Button } from '@/components/ui/button';
import type { Todo, TodoFilter } from '@/types/models';

/** Main todos page component. Displays a searchable todo list with create/edit/delete actions. */
export default function TodosPage() {
  const { todos, loading, createTodo, updateTodo, deleteTodo, search } = useTodos();
  const { tags: availableTags } = useTags();
  const [filter, setFilter] = useState<TodoFilter>({});
  const [formOpen, setFormOpen] = useState(false);
  const [editingTodo, setEditingTodo] = useState<Todo | null>(null);

  /** Handle status update for a todo item. */
  const handleUpdate = async (id: number, input: Partial<Todo>) => {
    const statusInput = input.status ? { status: input.status } : {};
    await updateTodo(id, statusInput as any);
  };

  /** Handle todo search; skips empty queries. */
  const handleSearch = async (query: string) => {
    if (!query.trim()) return;
    await search(query);
  };

  return (
    <div className="space-y-6">
      {/* Header with title and "New Todo" button */}
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-bold">Todos</h2>
        <Button onClick={() => { setEditingTodo(null); setFormOpen(true); }} size="sm">
          <Plus className="h-4 w-4 mr-1" /> New Todo
        </Button>
      </div>

      <SearchBar onSearch={handleSearch} placeholder="Search todos..." />

      {loading ? (
        <div className="py-8 text-center text-sm text-muted-foreground">Loading...</div>
      ) : (
        <TodoList
          todos={todos}
          filter={filter}
          onFilterChange={setFilter}
          onUpdate={handleUpdate}
          onDelete={deleteTodo}
          onEdit={(todo) => { setEditingTodo(todo); setFormOpen(true); }}
          availableTags={availableTags}
        />
      )}

      <TodoForm
        open={formOpen}
        onClose={() => { setFormOpen(false); setEditingTodo(null); }}
        onCreate={async (input) => { await createTodo(input); }}
        onUpdate={async (id, input) => { await updateTodo(id, input as any); }}
        todo={editingTodo}
      />
    </div>
  );
}
