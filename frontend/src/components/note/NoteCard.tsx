import { Pencil, Trash2 } from 'lucide-react';
import type { Note } from '@/types/models';
import { Card, CardHeader, CardContent, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { TagList } from '@/components/common/TagChip';
import { formatDate } from '@/lib/utils';

/** Props for the NoteCard component. */
interface NoteCardProps {
  /** The note object to display. */
  note: Note;
  /** Callback to open the editor for this note. */
  onEdit: (note: Note) => void;
  /** Callback to delete this note by its id. */
  onDelete: (id: number) => void;
}

/**
 * Card component that displays a note's title, content preview, tags and last-updated date.
 * Edit and Delete buttons appear on hover.
 */
export default function NoteCard({ note, onEdit, onDelete }: NoteCardProps) {
  return (
    <Card className="group transition-all hover:shadow-md">
      <CardHeader className="p-4 pb-2">
        <h3 className="text-sm font-semibold">{note.title || 'Untitled'}</h3>
      </CardHeader>
      <CardContent className="p-4 pt-0">
        <p className="whitespace-pre-wrap text-xs text-muted-foreground line-clamp-4">{note.content}</p>
        {note.tags && note.tags.length > 0 && (
          <div className="mt-2">
            <TagList tags={note.tags} />
          </div>
        )}
      </CardContent>
      <CardFooter className="flex items-center justify-between p-4 pt-0">
        <span className="text-xs text-muted-foreground">
          {formatDate(note.updatedAt)}
        </span>
        <div className="flex gap-1 opacity-0 transition-opacity group-hover:opacity-100">
          <Button variant="ghost" size="icon" className="h-7 w-7" onClick={() => onEdit(note)} title="Edit">
            <Pencil className="h-3.5 w-3.5" />
          </Button>
          <Button variant="ghost" size="icon" className="h-7 w-7 text-destructive hover:text-destructive" onClick={() => onDelete(note.id)} title="Delete">
            <Trash2 className="h-3.5 w-3.5" />
          </Button>
        </div>
      </CardFooter>
    </Card>
  );
}
