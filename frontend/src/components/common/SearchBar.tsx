import { useState, useCallback } from 'react';
import { Search } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { debounce } from '@/utils/debounce';

/** Props for the SearchBar component. */
interface SearchBarProps {
  /** Callback that receives the search query string after debounce. */
  onSearch: (q: string) => void;
  /** Placeholder text shown in the input field. Defaults to "Search...". */
  placeholder?: string;
}

/**
 * Text input with a magnifying-glass icon and debounced search callback.
 * Fires `onSearch` 300 ms after the user stops typing.
 */
export default function SearchBar({ onSearch, placeholder = 'Search...' }: SearchBarProps) {
  const [value, setValue] = useState('');

  const handleSearch = useCallback(
    debounce((q: string) => onSearch(q), 300),
    [onSearch]
  );

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    setValue(val);
    handleSearch(val);
  };

  return (
    <div className="relative">
      <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
      <Input
        type="text"
        value={value}
        onChange={handleChange}
        placeholder={placeholder}
        className="pl-10"
      />
    </div>
  );
}
