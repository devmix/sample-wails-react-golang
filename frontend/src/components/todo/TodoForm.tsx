import { useState, useEffect } from 'react';
import * as Dialog from '@radix-ui/react-dialog';
import { X, Plus, Save } from 'lucide-react';
import type { Todo, TodoCreateInput } from '@/types/models';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Button } from '@/components/ui/button';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select';
import { TagChip } from '@/components/common/TagChip';

/**
 * Partial data used when updating an existing todo.
 * Only the fields that changed need to be provided.
 */
export interface TodoUpdateInput {
  /** New title for the todo. */
  title?: string;
  /** New description for the todo. */
  description?: string;
  /** New priority level. */
  priority?: 'low' | 'medium' | 'high' | 'urgent';
  /** ISO date string or null to clear the due date. */
  dueDate?: string | null;
  /** Replacement list of tag names. */
  tagNames?: string[];
}

/** Props for the TodoForm component. */
interface TodoFormProps {
  /** Whether the dialog is currently open. */
  open: boolean;
  /** Callback to close the dialog without submitting. */
  onClose: () => void;
  /** Callback invoked with new todo data when creating a todo. */
  onCreate: (input: TodoCreateInput) => void;
  /** Optional callback for updating an existing todo. */
  onUpdate?: (id: number, input: TodoUpdateInput) => void;
  /** Existing todo being edited; absent when creating a new todo. */
  todo?: Todo | null;
}

/**
 * Modal dialog form for creating or editing a todo item.
 * Supports title, description, priority selector, due date picker and tag management.
 */
export default function TodoForm({ open, onClose, onCreate, onUpdate, todo }: TodoFormProps) {
  const isEdit = !!todo && !!onUpdate;

  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [priority, setPriority] = useState<'low' | 'medium' | 'high' | 'urgent'>('medium');
  const [dueDate, setDueDate] = useState('');
  const [tagInput, setTagInput] = useState('');
  const [tags, setTags] = useState<string[]>([]);

  useEffect(() => {
    if (todo) {
      setTitle(todo.title || '');
      setDescription(todo.description || '');
      setPriority((todo.priority as any) || 'medium');
      setDueDate(todo.dueDate?.split('T')[0] || '');
      setTags(todo.tags?.map(t => t.name) || []);
    } else {
      setTitle('');
      setDescription('');
      setPriority('medium');
      setDueDate('');
      setTags([]);
    }
    setTagInput('');
  }, [todo, open]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    if (isEdit && todo) {
      onUpdate!(todo.id, {
        title: title.trim(),
        description: description.trim() || undefined,
        priority,
        dueDate: dueDate || null,
        tagNames: tags.length > 0 ? tags : undefined,
      });
    } else {
      onCreate({
        title: title.trim(),
        description: description.trim(),
        priority,
        dueDate: dueDate || undefined,
        tagNames: tags,
      });
    }

    onClose();
  };

  const addTag = () => {
    const tag = tagInput.trim();
    if (tag && !tags.includes(tag)) {
      setTags(prev => [...prev, tag]);
      setTagInput('');
    }
  };

  return (
    <Dialog.Root open={open} onOpenChange={onClose}>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 bg-black/50 z-40" />
        <Dialog.Content className="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-full max-w-lg rounded-lg border bg-card p-6 shadow-lg z-50">
          <div className="flex items-center justify-between mb-4">
            <Dialog.Title className="text-lg font-semibold">{isEdit ? 'Edit Todo' : 'New Todo'}</Dialog.Title>
            <button onClick={onClose} className="rounded-sm opacity-70 hover:opacity-100">
              <X className="h-5 w-5" />
            </button>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            <Input
              type="text"
              value={title}
              onChange={e => setTitle(e.target.value)}
              placeholder="Todo title *"
              required
            />

            <Textarea
              value={description}
              onChange={e => setDescription(e.target.value)}
              placeholder="Description (optional)"
              rows={2}
            />

            <div className="flex gap-4">
              <Select value={priority} onValueChange={v => setPriority(v as 'low' | 'medium' | 'high' | 'urgent')}>
                <SelectTrigger className="w-[160px]">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="low">Low</SelectItem>
                  <SelectItem value="medium">Medium</SelectItem>
                  <SelectItem value="high">High</SelectItem>
                  <SelectItem value="urgent">Urgent</SelectItem>
                </SelectContent>
              </Select>

              <Input
                type="date"
                value={dueDate}
                onChange={e => setDueDate(e.target.value)}
              />
            </div>

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

            <div className="flex justify-end gap-2 pt-2">
              <Button type="button" variant="outline" onClick={onClose}>
                Cancel
              </Button>
              <Button type="submit">
                {isEdit ? (
                  <>
                    <Save className="h-4 w-4 mr-1" /> Save Changes
                  </>
                ) : (
                  <>
                    <Plus className="h-4 w-4 mr-1" /> Add Todo
                  </>
                )}
              </Button>
            </div>
          </form>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
