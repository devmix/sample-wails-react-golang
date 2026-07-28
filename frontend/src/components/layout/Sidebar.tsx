import { useState, useEffect } from 'react';
import { CheckSquare, FileText, Settings, Moon, Sun } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

/** Allowed navigation page identifiers. */
type Page = 'todos' | 'notes' | 'settings';

/** Props for the Sidebar component. */
interface SidebarProps {
  /** Currently active page key. */
  activePage: Page;
  /** Callback invoked when a navigation item is selected. */
  onNavigate: (page: Page) => void;
}

/**
 * Vertical sidebar with navigation links and a dark/light mode toggle.
 * Persists the theme preference in localStorage.
 */
export default function Sidebar({ activePage, onNavigate }: SidebarProps) {
  const [darkMode, setDarkMode] = useState(() => {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('theme') === 'dark' ||
        (!localStorage.getItem('theme') && window.matchMedia('(prefers-color-scheme: dark)').matches);
    }
    return false;
  });

  useEffect(() => {
    const root = document.documentElement;
    if (darkMode) {
      root.classList.add('dark');
      localStorage.setItem('theme', 'dark');
    } else {
      root.classList.remove('dark');
      localStorage.setItem('theme', 'light');
    }
  }, [darkMode]);

  const items: { key: Page; label: string; icon: any }[] = [
    { key: 'todos', label: 'Todos', icon: CheckSquare },
    { key: 'notes', label: 'Notes', icon: FileText },
    { key: 'settings', label: 'Settings', icon: Settings },
  ];

  return (
    <aside className="flex h-full w-56 flex-col border-border bg-muted/30">
      <div className="flex items-center gap-2 border-b border-border px-4 py-4">
        <CheckSquare className="h-7 w-7 text-primary" />
        <h1 className="text-lg font-bold">ToDo Notes</h1>
      </div>
      <nav className="flex-1 px-2 py-4">
        <ul className="space-y-1">
          {items.map(item => {
            const Icon = item.icon;
            return (
              <li key={item.key}>
                <Button
                  variant={activePage === item.key ? 'secondary' : 'ghost'}
                  className={cn(
                    'w-full justify-start gap-3 font-medium',
                    activePage === item.key && 'bg-secondary text-secondary-foreground'
                  )}
                  onClick={() => onNavigate(item.key)}
                >
                  <Icon className="h-5 w-5" />
                  {item.label}
                </Button>
              </li>
            );
          })}
        </ul>
      </nav>
      <div className="border-t border-border p-2">
        <Button
          variant="ghost"
          className="w-full justify-start gap-3 font-medium"
          onClick={() => setDarkMode(!darkMode)}
        >
          {darkMode ? (
            <>
              <Sun className="h-5 w-5" />
              Light Mode
            </>
          ) : (
            <>
              <Moon className="h-5 w-5" />
              Dark Mode
            </>
          )}
        </Button>
      </div>
    </aside>
  );
}
