import { useState } from 'react';

export default function NewPhoto() {
  const [file, setFile] = useState<File | null>(null);
  const [tags, setTags] = useState('');


  return (
    <form className="upload-form">
      <h2>Upload a Photo</h2>

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

      <button type="submit" className="upload-btn" disabled={!file}>
        Upload
      </button>
    </form>
  );
}
