import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { getPost, type Post } from '../api/client';
import Navbar from '../components/Navbar';
import StatusBar from '../components/StatusBar';
import './PostReader.css';

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

export default function PostReader() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [post, setPost] = useState<Post | null>(null);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!id) return;
    getPost(Number(id))
      .then((p) => {
        if (!p) setError('Post not found');
        else {
          setPost(p);
          document.title = p.title;
        }
      })
      .catch((e) => setError((e as Error).message));
  }, [id]);

  useEffect(() => {
    function handleKey(e: KeyboardEvent) {
      if ((e.target as HTMLElement).tagName === 'INPUT') return;
      if (e.key === 'h') navigate('/');
    }
    window.addEventListener('keydown', handleKey);
    return () => window.removeEventListener('keydown', handleKey);
  }, [navigate]);

  const slug = post?.title.toLowerCase().replace(/\s+/g, '-') ?? '';

  return (
    <div className="reader-layout">
      <Navbar />
      <div className="reader-container">
        <button className="reader-back" onClick={() => navigate('/')}>
          [..]
        </button>

        {error && <p className="reader-error">{error}</p>}

        {post && (
          <>
            {post.cover_image_url && (
              <img
                src={post.cover_image_url}
                alt="cover"
                className="reader-cover"
              />
            )}

            <h1 className="reader-title">{post.title}</h1>

            <div className="reader-meta">
              <span className="reader-date">{formatDate(post.created_at)}</span>
              {post.tags?.length > 0 && (
                <span className="reader-tags">
                  {post.tags.map((t) => (
                    <span key={t} className="reader-tag">
                      {t}
                    </span>
                  ))}
                </span>
              )}
            </div>

            <div className="reader-content">
              <ReactMarkdown
                remarkPlugins={[remarkGfm]}
                components={{
                  code({ className, children, ...props }) {
                    const match = /language-(\w+)/.exec(className || '');
                    const isInline = !match;
                    return isInline ? (
                      <code className="reader-inline-code" {...props}>
                        {children}
                      </code>
                    ) : (
                      <SyntaxHighlighter
                        style={vscDarkPlus}
                        language={match[1]}
                        PreTag="div"
                      >
                        {String(children).replace(/\n$/, '')}
                      </SyntaxHighlighter>
                    );
                  },
                }}
              >
                {post.content}
              </ReactMarkdown>
            </div>
          </>
        )}
      </div>

      <StatusBar
        left={slug ? `${slug}.md` : '...'}
        right="READ"
      />
    </div>
  );
}
