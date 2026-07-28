import { useEffect, useState } from 'react';
import * as AppService from '@bindings/sample-wails-react-golang/internal/app/service/appservice';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';

/** Settings page component. Displays application info and placeholder for keyboard shortcuts. */
export default function SettingsPage() {
  const [version, setVersion] = useState('dev');

  useEffect(() => {
    AppService.Version().then(setVersion);
  }, []);

  return (
    <div className="space-y-6">
      <h2 className="text-xl font-bold">Settings</h2>

      <Card>
        <CardHeader>
          <CardTitle className="text-sm font-semibold">Application Info</CardTitle>
        </CardHeader>
        <CardContent>
          <dl className="space-y-2 text-sm">
            <div className="flex justify-between">
              <dt className="text-muted-foreground">Version</dt>
              <dd className="font-medium">{version}</dd>
            </div>
            <div className="flex justify-between">
              <dt className="text-muted-foreground">Storage</dt>
              <dd className="font-medium">SQLite (Local)</dd>
            </div>
          </dl>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-sm font-semibold">Keyboard Shortcuts</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">Coming soon...</p>
        </CardContent>
      </Card>
    </div>
  );
}
