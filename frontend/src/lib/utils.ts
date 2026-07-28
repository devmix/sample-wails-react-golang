import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

/**
 * Merges conditional Tailwind CSS classes using clsx and tailwind-merge.
 * Resolves conflicting utility classes so the last one wins.
 * @param inputs - Class values to merge (strings, arrays, objects, conditionals).
 * @returns A single merged class name string.
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * Formats an ISO date string into a locale-aware display date.
 * Returns '—' for null, undefined, or invalid dates.
 * @param dateStr - An ISO date string, null, or undefined.
 * @returns A formatted date string or '—' if the input is invalid.
 */
export function formatDate(dateStr: string | undefined | null) {
  if (!dateStr) return '—';
  const d = new Date(dateStr);
  return Number.isNaN(d.getTime()) ? '—' : d.toLocaleDateString();
}
