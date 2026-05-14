import { useState, useEffect, useCallback } from 'react';

interface Photo {
  Path: string;
  Tags: string[];
  CreatedAt: string;
  Data?: string;
}

interface PhotosResponse {
  photos: Photo[];
  next_cursor?: string;
}

interface HomeProps {
  newPhotos: Photo[];
}

export default function Home({ newPhotos}: HomeProps) {
  const [photos, setPhotos] = useState<Photo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filterTags, setFilterTags] = useState('');

  const fetchPhotos = useCallback(async (tags: string) => {
    setLoading(true);
    setError(null);
    try {
      const params = new URLSearchParams();
      if (tags.trim()) {
        params.set('tags', tags.trim());
      }
      const url = `http://127.0.0.1:8000/photos?${params.toString()}`;
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error('Network response was not okay');
      }
      const data: PhotosResponse = await response.json();
      setPhotos(data.photos || []);
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('An unexpected error occurred');
      }
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchPhotos('');
  }, [fetchPhotos]);

  useEffect(() => {
    if (newPhotos.length > 0) {
      const activeTags = filterTags
        .split(',')
        .map((t) => t.trim().toLowerCase())
        .filter(Boolean);

      const matching = activeTags.length > 0
        ? newPhotos.filter((photo) =>
            photo.Tags?.some((tag) => activeTags.includes(tag.toLowerCase()))
          )
        : newPhotos;

      if (matching.length > 0) {
        setPhotos((prev) => [...matching, ...prev]);
      }
    }
  }, [newPhotos, filterTags]);

  const handleFilter = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    fetchPhotos(filterTags);
  };

  return (
    <>
      <form className="filter-bar" onSubmit={handleFilter}>
        <input
          type="text"
          placeholder="Filter by tags (e.g. sunset, beach)"
          value={filterTags}
          onChange={(e) => setFilterTags(e.target.value)}
        />
        <button type="submit">Filter</button>
        {filterTags && (
          <button
            type="button"
            onClick={() => { setFilterTags(''); fetchPhotos(''); }}
          >
            Clear
          </button>
        )}
      </form>

      {loading && <div className="empty-state">Loading...</div>}
      {error && <div className="empty-state">Error: {error}</div>}
      {!loading && !error && photos.length === 0 && (
        <div className="empty-state">No photos found.</div>
      )}

      {!loading && !error && photos.length > 0 && (
        <div className="feed">
          {photos.map((photo, i) => (
            <div className="feed-card" key={i}>
              <img
                src={photo.Data ? `data:image/jpeg;base64,${photo.Data}` : photo.Path}
                alt={photo.Tags?.join(', ') || 'uploaded content'}
              />
              <div className="feed-card-info">
                <div className="feed-card-tags">
                  {photo.Tags?.map((tag) => (
                    <span className="tag" key={tag}>#{tag}</span>
                  ))}
                </div>
                <div className="feed-card-date">
                  {new Date(photo.CreatedAt).toLocaleDateString()}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </>
  );
}
