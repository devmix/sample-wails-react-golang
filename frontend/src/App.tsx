import { useState } from 'react';
import Sidebar from '@/components/layout/Sidebar';
import TodosPage from '@/pages/TodosPage';
import NotesPage from '@/pages/NotesPage';
import SettingsPage from '@/pages/SettingsPage';

type Page = 'todos' | 'notes' | 'settings';

export default function App() {
  const [activePage, setActivePage] = useState<Page>('todos');

  const renderPage = () => {
    switch (activePage) {
      case 'todos': return <TodosPage />;
      case 'notes': return <NotesPage />;
      case 'settings': return <SettingsPage />;
    }
  };

  return (
    <div className="flex h-screen bg-background">
      <Sidebar activePage={activePage} onNavigate={setActivePage} />
      <main className="flex-1 overflow-y-auto p-6">
        {renderPage()}
      </main>
    </div>
  );
}
