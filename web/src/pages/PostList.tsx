import { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { listPosts, type Post } from '../api/client';
import Navbar from '../components/Navbar';
import StatusBar from '../components/StatusBar';
import './PostList.css';

function formatDate(iso: string): string {
  return iso.slice(0, 10);
}

function toSlug(title: string): string {
  return title.toLowerCase().replace(/\s+/g, '_').replace(/[^a-z0-9_]/g, '');
}

export default function PostList() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [cursor, setCursor] = useState(0);
  const [error, setError] = useState('');
  const navigate = useNavigate();
  const listRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    document.title = 'Home';
  }, []);

  useEffect(() => {
    listPosts()
      .then(setPosts)
      .catch((e) => setError((e as Error).message));
  }, []);

  useEffect(() => {
    function handleKey(e: KeyboardEvent) {
      if ((e.target as HTMLElement).tagName === 'INPUT') return;

      if (e.key === 'j') {
        setCursor((c) => Math.min(c + 1, posts.length - 1));
      } else if (e.key === 'k') {
        setCursor((c) => Math.max(c - 1, 0));
      } else if (e.key === 'l' || e.key === 'Enter') {
        if (posts[cursor]) navigate(`/posts/${posts[cursor].id}/${toSlug(posts[cursor].title)}`);
      }
    }

    window.addEventListener('keydown', handleKey);
    return () => window.removeEventListener('keydown', handleKey);
  }, [posts, cursor, navigate]);

  useEffect(() => {
    const row = listRef.current?.querySelector(`[data-index="${cursor}"]`);
    row?.scrollIntoView({ block: 'nearest' });
  }, [cursor]);

  return (
    <div className="postlist-layout">
      <Navbar />
      <div className="postlist-container" ref={listRef}>
        <div className="postlist-parent">..</div>

        {error && <div className="postlist-error">{error}</div>}

        {posts.map((post, i) => (
          <div
            key={post.id}
            data-index={i}
            className={`postlist-row ${i === cursor ? 'active' : ''}`}
            onClick={() => navigate(`/posts/${post.id}/${toSlug(post.title)}`)}
            onMouseEnter={() => setCursor(i)}
          >
            <span className="postlist-cursor">{i === cursor ? '▶' : ' '}</span>
            <span className="postlist-index">{i + 1}</span>
            <span className="postlist-title">{post.title}</span>
            <span className="postlist-date">{formatDate(post.created_at)}</span>
            <span className="postlist-tags">
              {post.tags?.length > 0 && `[${post.tags.join(', ')}]`}
            </span>
          </div>
        ))}
      </div>

      <StatusBar
        left="gustavo's blog/"
        center={`${posts.length} posts`}
        right="NORMAL"
      />
    </div>
  );
}
