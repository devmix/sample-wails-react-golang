import { useState, useEffect, useCallback } from 'react';
import * as TodoService from '@bindings/sample-wails-react-golang/internal/app/service/todoservice';
import type { Todo, TodoFilter, TodoCreateInput, TodoUpdateInput } from '../types/models';

/** Ensures every todo has a non-null tags array by defaulting to an empty array. */
function normalizeTodo(t: Todo): Todo {
  return { ...t, tags: t.tags || [] };
}

/**
 * React hook that manages the full lifecycle of todos (list, create, update, delete, search).
 *
 * Fetches all todos on mount and keeps local state in sync with mutations.
 *
 * @returns An object containing:
 *   - `todos`       – current list of todos
 *   - `loading`     – true while the initial fetch is in progress
 *   - `error`       – last error message, or null if no error occurred
 *   - `fetchTodos`  – re-fetch todos with an optional filter
 *   - `createTodo`  – create a new todo and append it to local state
 *   - `updateTodo`  – update an existing todo by id
 *   - `deleteTodo`  – remove a todo by id
 *   - `search`      – search todos by query string
 */
export function useTodos() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  /** Fetch todos from the backend, optionally filtered. */
  const fetchTodos = useCallback(async (filter?: TodoFilter) => {
    try {
      setLoading(true);
      const result = await TodoService.List(filter ?? null);
      setTodos(result ? result.map(normalizeTodo) : []);
      setError(null);
    } catch (e: unknown) {
      setError((e as Error)?.message || 'Failed to load todos');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchTodos();
  }, [fetchTodos]);

  /** Create a new todo. Returns the created todo or null on failure. */
  const createTodo = async (input: TodoCreateInput): Promise<Todo | null> => {
    try {
      const result = await TodoService.Create(input);
      if (result) {
        const normalized = normalizeTodo(result);
        setTodos(prev => [...prev, normalized]);
        return normalized;
      }
      return null;
    } catch (e: unknown) {
      setError((e as Error)?.message || 'Failed to create todo');
      return null;
    }
  };

  /** Update an existing todo by id. Returns the updated todo or null on failure. */
  const updateTodo = async (id: number, input: TodoUpdateInput): Promise<Todo | null> => {
    try {
      const result = await TodoService.Update(id, input);
      if (result) {
        const normalized = normalizeTodo(result);
        setTodos(prev => prev.map(t => t.id === id ? normalized : t));
        return normalized;
      }
      return null;
    } catch (e: unknown) {
      setError((e as Error)?.message || 'Failed to update todo');
      return null;
    }
  };

  /** Delete a todo by id. Returns true on success, false on failure. */
  const deleteTodo = async (id: number): Promise<boolean> => {
    try {
      await TodoService.Delete(id);
      setTodos(prev => prev.filter(t => t.id !== id));
      return true;
    } catch (e: unknown) {
      setError((e as Error)?.message || 'Failed to delete todo');
      return false;
    }
  };

  /** Search todos by a free-text query. Returns matching todos or an empty array on failure. */
  const search = async (query: string): Promise<Todo[]> => {
    try {
      const result = await TodoService.Search(query);
      return result ? result.map(normalizeTodo) : [];
    } catch (e: unknown) {
      setError((e as Error)?.message || 'Search failed');
      return [];
    }
  };

  return { todos, loading, error, fetchTodos, createTodo, updateTodo, deleteTodo, search };
}
