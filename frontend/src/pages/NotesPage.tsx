import { useState } from 'react';
import { Plus } from 'lucide-react';
import { useNotes } from '@/hooks/useNotes';
import NoteList from '@/components/note/NoteList';
import NoteEditor, { type NoteEditorInput } from '@/components/note/NoteEditor';
import SearchBar from '@/components/common/SearchBar';
import { Button } from '@/components/ui/button';
import type { Note } from '@/types/models';

/** Main notes page component. Displays a searchable note list with create/edit/delete actions. */
export default function NotesPage() {
  const { notes, loading, createNote, updateNote, deleteNote, search } = useNotes();
  const [editingNote, setEditingNote] = useState<Note | null>(null);
  const [showEditor, setShowEditor] = useState(false);

  /** Open the editor in create mode with no pre-filled note. */
  const handleNew = () => {
    setEditingNote(null);
    setShowEditor(true);
  };

  /** Open the editor in edit mode for an existing note. */
  const handleEdit = (note: Note) => {
    setEditingNote(note);
    setShowEditor(true);
  };

  /** Submit handler that either creates a new note or updates an existing one. */
  const handleSubmit = async (input: NoteEditorInput) => {
    if (editingNote) {
      await updateNote(editingNote.id, input as any);
    } else {
      await createNote(input);
    }
    setShowEditor(false);
    setEditingNote(null);
  };

  /** Handle note search; skips empty queries. */
  const handleSearch = async (query: string) => {
    if (!query.trim()) return;
    await search(query);
  };

  return (
    <div className="space-y-6">
      {/* Header with title and "New Note" button */}
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-bold">Notes</h2>
        {!showEditor && (
          <Button onClick={handleNew} size="sm">
            <Plus className="h-4 w-4 mr-1" /> New Note
          </Button>
        )}
      </div>

      {!showEditor && <SearchBar onSearch={handleSearch} placeholder="Search notes..." />}

      {showEditor ? (
        <NoteEditor onSubmit={handleSubmit} onCancel={() => { setShowEditor(false); setEditingNote(null); }} initialNote={editingNote} />
      ) : loading ? (
        <div className="py-8 text-center text-sm text-muted-foreground">Loading...</div>
      ) : (
        <NoteList notes={notes} onEdit={handleEdit} onDelete={deleteNote} />
      )}
    </div>
  );
}
