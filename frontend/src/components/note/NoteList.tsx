import type { Note } from '../../types/models';
import NoteCard from './NoteCard';

/** Props for the NoteList component. */
interface NoteListProps {
  /** Array of note objects to display in a responsive grid. */
  notes: Note[];
  /** Callback to open the editor for a given note. */
  onEdit: (note: Note) => void;
  /** Callback to delete a note by its id. */
  onDelete: (id: number) => void;
}

/**
 * Responsive grid of `NoteCard` components.
 * Shows an empty-state message when no notes are present.
 */
export default function NoteList({ notes, onEdit, onDelete }: NoteListProps) {
  return (
    <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
      {notes.length === 0 ? (
        <p className="col-span-full py-8 text-center text-sm text-muted-foreground">No notes yet</p>
      ) : (
        notes.map(note => (
          <NoteCard key={note.id} note={note} onEdit={onEdit} onDelete={onDelete} />
        ))
      )}
    </div>
  );
}
