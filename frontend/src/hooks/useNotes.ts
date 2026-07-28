import { useState, useEffect, useCallback } from 'react';
import * as NoteService from '@bindings/sample-wails-react-golang/internal/app/service/noteservice';
import type { Note, NoteCreateInput, NoteUpdateInput } from '../types/models';

/** Ensures every note has a non-null tags array by defaulting to an empty array. */
function normalizeNote(n: Note): Note {
  return { ...n, tags: n.tags || [] };
}

/**
 * React hook that manages the full lifecycle of notes (list, create, update, delete, search).
 *
 * Fetches all notes on mount and keeps local state in sync with mutations.
 *
 * @returns An object containing:
 *   - `notes`       – current list of notes
 *   - `loading`     – true while the initial fetch is in progress
 *   - `error`       – last error message, or null if no error occurred
 *   - `fetchNotes`  – re-fetch all notes from the backend
 *   - `createNote`  – create a new note and prepend it to local state
 *   - `updateNote`  – update an existing note by id
 *   - `deleteNote`  – remove a note by id
 *   - `search`      – search notes by query string
 */
export function useNotes() {
  const [notes, setNotes] = useState<Note[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  /** Fetch all notes from the backend. */
  const fetchNotes = useCallback(async () => {
    try {
      setLoading(true);
      const result = await NoteService.List();
      setNotes(result ? result.map(normalizeNote) : []);
      setError(null);
    } catch (e: unknown) {
      setError((e as Error)?.message || 'Failed to load notes');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchNotes();
  }, [fetchNotes]);

  /** Create a new note. Returns the created note or null on failure. */
  const createNote = async (input: NoteCreateInput): Promise<Note | null> => {
    try {
      const result = await NoteService.Create(input);
      if (result) {
        const normalized = normalizeNote(result);
        setNotes(prev => [normalized, ...prev]);
        return normalized;
      }
      return null;
    } catch (e: unknown) {
      setError((e as Error)?.message || 'Failed to create note');
      return null;
    }
  };

  /** Update an existing note by id. Returns the updated note or null on failure. */
  const updateNote = async (id: number, input: NoteUpdateInput): Promise<Note | null> => {
    try {
      const result = await NoteService.Update(id, input);
      if (result) {
        const normalized = normalizeNote(result);
        setNotes(prev => prev.map(n => n.id === id ? normalized : n));
        return normalized;
      }
      return null;
    } catch (e: unknown) {
      setError((e as Error)?.message || 'Failed to update note');
      return null;
    }
  };

  /** Delete a note by id. Returns true on success, false on failure. */
  const deleteNote = async (id: number): Promise<boolean> => {
    try {
      await NoteService.Delete(id);
      setNotes(prev => prev.filter(n => n.id !== id));
      return true;
    } catch (e: unknown) {
      setError((e as Error)?.message || 'Failed to delete note');
      return false;
    }
  };

  /** Search notes by a free-text query. Returns matching notes or an empty array on failure. */
  const search = async (query: string): Promise<Note[]> => {
    try {
      const result = await NoteService.Search(query);
      return result ? result.map(normalizeNote) : [];
    } catch (e: unknown) {
      setError((e as Error)?.message || 'Search failed');
      return [];
    }
  };

  return { notes, loading, error, fetchNotes, createNote, updateNote, deleteNote, search };
}
