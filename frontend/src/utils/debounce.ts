/**
 * Returns a debounced version of the provided function.
 * The wrapped function will only execute after `delay` milliseconds have elapsed
 * since the last invocation. Previous timers are cleared on each call.
 * @param fn - The function to debounce.
 * @param delay - Delay in milliseconds before executing the function.
 * @returns A debounced function with the same signature as `fn`.
 */
export function debounce<T extends (...args: any[]) => void>(fn: T, delay: number): T {
  let timer: ReturnType<typeof setTimeout>;
  return ((...args: Parameters<T>) => {
    clearTimeout(timer);
    timer = setTimeout(() => fn(...args), delay);
  }) as T;
}
