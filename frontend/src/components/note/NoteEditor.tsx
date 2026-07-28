import { useState } from 'react';
import { X, Save } from 'lucide-react';
import type { Note } from '@/types/models';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Button } from '@/components/ui/button';
import { TagChip } from '@/components/common/TagChip';

/**
 * Data submitted from the note editor form.
 */
export interface NoteEditorInput {
  /** Title of the note. */
  title: string;
  /** Body content of the note. */
  content: string;
  /** List of tag names attached to the note. */
  tagNames: string[];
}

/** Props for the NoteEditor component. */
interface NoteEditorProps {
  /** Callback invoked on form submit with the editor's data. */
  onSubmit: (input: NoteEditorInput) => void;
  /** Callback to cancel editing and close the editor. */
  onCancel: () => void;
  /** Existing note used to pre-fill fields when editing. */
  initialNote?: Note | null;
}

/**
 * Form for creating or editing a note with title, content area and tag management.
 * Pre-fills fields from `initialNote` when provided.
 */
export default function NoteEditor({ onSubmit, onCancel, initialNote }: NoteEditorProps) {
  const [title, setTitle] = useState(initialNote?.title || '');
  const [content, setContent] = useState(initialNote?.content || '');
  const [tagInput, setTagInput] = useState('');
  const [tags, setTags] = useState<string[]>(initialNote?.tags?.map(t => t.name) || []);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({ title: title.trim(), content: content.trim(), tagNames: tags });
  };

  const addTag = () => {
    const tag = tagInput.trim();
    if (tag && !tags.includes(tag)) {
      setTags(prev => [...prev, tag]);
      setTagInput('');
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4 rounded-lg border bg-card p-4">
      <Input
        type="text"
        value={title}
        onChange={e => setTitle(e.target.value)}
        placeholder="Note title"
        className="font-medium"
      />

      <Textarea
        value={content}
        onChange={e => setContent(e.target.value)}
        placeholder="Write your note..."
        rows={8}
        className="resize-y"
      />

      <div className="flex gap-2">
        <Input
          type="text"
          value={tagInput}
          onChange={e => setTagInput(e.target.value)}
          onKeyDown={e => { if (e.key === 'Enter') { e.preventDefault(); addTag(); } }}
          placeholder="Add tag + Enter"
          className="flex-1"
        />
      </div>

      {tags.length > 0 && (
        <div className="flex flex-wrap gap-1">
          {tags.map(tag => (
            <TagChip key={tag} name={tag} onRemove={() => setTags(prev => prev.filter(t => t !== tag))} />
          ))}
        </div>
      )}

      <div className="flex justify-end gap-2">
        <Button type="button" variant="outline" onClick={onCancel}>
          <X className="h-4 w-4 mr-1" />
          Cancel
        </Button>
        <Button type="submit">
          <Save className="h-4 w-4 mr-1" />
          Save Note
        </Button>
      </div>
    </form>
  );
}
