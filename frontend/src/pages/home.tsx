export default function Home() {
  // Placeholder data — replace with API call later
  const photos = [
    { path: 'https://picsum.photos/600/400?random=1', tags: ['nature', 'sunset'], createdAt: '2026-05-13T14:30:00Z' },
    { path: 'https://picsum.photos/600/400?random=2', tags: ['city', 'night'], createdAt: '2026-05-12T10:00:00Z' },
    { path: 'https://picsum.photos/600/400?random=3', tags: ['food'], createdAt: '2026-05-11T08:15:00Z' },
  ];

  if (photos.length === 0) {
    return <div className="empty-state">No photos yet. Upload your first one!</div>;
  }

  return (
    <div className="feed">
      {photos.map((photo, i) => (
        <div className="feed-card" key={i}>
          <img src={photo.path} alt="photo" />
          <div className="feed-card-info">
            <div className="feed-card-tags">
              {photo.tags.map((tag) => (
                <span className="tag" key={tag}>#{tag}</span>
              ))}
            </div>
            <div className="feed-card-date">
              {new Date(photo.createdAt).toLocaleDateString()}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
