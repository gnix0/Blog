import { useRef, useState } from 'react';
import { uploadImage } from '../api/client';
import './CoverUpload.css';

interface CoverUploadProps {
  value: string;
  onChange: (url: string) => void;
}

export default function CoverUpload({ value, onChange }: CoverUploadProps) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState('');

  async function handleFile(file: File) {
    setError('');
    setUploading(true);
    try {
      const { url } = await uploadImage(file);
      onChange(url);
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setUploading(false);
    }
  }

  function handleDrop(e: React.DragEvent) {
    e.preventDefault();
    const file = e.dataTransfer.files[0];
    if (file) handleFile(file);
  }

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (file) handleFile(file);
  }

  return (
    <div
      className={`cover-upload ${value ? 'has-image' : ''}`}
      onDrop={handleDrop}
      onDragOver={(e) => e.preventDefault()}
      onClick={() => inputRef.current?.click()}
    >
      <input
        ref={inputRef}
        type="file"
        accept="image/jpeg,image/png,image/webp"
        style={{ display: 'none' }}
        onChange={handleChange}
      />
      {value ? (
        <img src={value} alt="cover" className="cover-preview" />
      ) : (
        <span className="cover-placeholder">
          {uploading ? 'uploading...' : '[ drop cover image or click ]'}
        </span>
      )}
      {error && <span className="cover-error">{error}</span>}
    </div>
  );
}
