import { useState, useEffect } from 'react';

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

export default function Home() {
  const [photos, setPhotos] = useState<Photo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPhotos = async () => {
      try {
        const response = await fetch('http://127.0.0.1:8000/photos');
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
    };

    fetchPhotos();
  }, []);

  if (loading) return  <div className="empty-state">Loading...</div>;
  if (error) return  <div className="empty-state">Error: {error}</div>;

  if (photos.length === 0) {
    return <div className="empty-state">No photos yet. Upload your first one!</div>;
  }

  return (
    <div className="feed">
      {photos.map((photo, i) => (
        <div className="feed-card" key={i}>
          <img
            src={photo.Data ? `data:image/jpeg;base64,${photo.Data}` : photo.Path}
            alt="photo"
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
  );
}
