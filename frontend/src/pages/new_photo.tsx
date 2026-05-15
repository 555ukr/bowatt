import { useState } from 'react';

export default function NewPhoto() {
  const [file, setFile] = useState<File | null>(null);
  const [tags, setTags] = useState('');
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!file) return;

    setUploading(true);
    setError(null);
    setSuccess(false);

    const formData = new FormData();
    formData.append('photo', file);
    formData.append('tags', tags);

    try {
      const response = await fetch('http://127.0.0.1:8000/upload', {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        const text = await response.text();
        throw new Error(text || 'Upload failed');
      }

      setSuccess(true);
      setFile(null);
      setTags('');
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('An unexpected error occurred');
      }
    } finally {
      setUploading(false);
    }
  };

  return (
    <form className="upload-form" onSubmit={handleSubmit}>
      <h2>Upload a Photo</h2>

      {error && <div style={{ color: 'red', marginBottom: 12 }}>{error}</div>}
      {success && <div style={{ color: 'green', marginBottom: 12 }}>Photo uploaded successfully!</div>}

      <label htmlFor="photo">Choose file</label>
      <input
        id="photo"
        type="file"
        accept="image/*"
        onChange={(e) => setFile(e.target.files?.[0] || null)}
      />

      <label htmlFor="tags">Tags (comma-separated)</label>
      <input
        id="tags"
        type="text"
        placeholder="sunset, beach, vacation"
        value={tags}
        onChange={(e) => setTags(e.target.value)}
      />

      <button type="submit" className="upload-btn" disabled={!file || uploading}>
        {uploading ? 'Uploading...' : 'Upload'}
      </button>
    </form>
  );
}
