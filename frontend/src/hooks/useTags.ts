import { useState, useEffect } from 'react';
import * as TagService from '@bindings/sample-wails-react-golang/internal/app/service/tagservice';
import type { Tag } from '../types/models';

/**
 * React hook that fetches all available tags on mount.
 *
 * @returns An object containing:
 *   - `tags`    – current list of tags
 *   - `loading` – true while the initial fetch is in progress
 */
export function useTags() {
  const [tags, setTags] = useState<Tag[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    TagService.ListAll().then(result => {
      setTags(result ? Array.from(result) : []);
      setLoading(false);
    }).catch(() => setLoading(false));
  }, []);

  return { tags, loading };
}
