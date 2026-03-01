import { useEffect, useRef, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import MDEditor from '@uiw/react-md-editor';
import { createPost, deletePost, getPost, listPosts, updatePost, type Post } from '../api/client';
import CoverUpload from '../components/CoverUpload';
import StatusBar from '../components/StatusBar';
import './Write.css';

export default function Write() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const editId = searchParams.get('id');

  const [title, setTitle] = useState('');
  const [tagsInput, setTagsInput] = useState('');
  const [content, setContent] = useState('');
  const [coverImageUrl, setCoverImageUrl] = useState('');
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const [previewing, setPreviewing] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [sidebarPosts, setSidebarPosts] = useState<Post[]>([]);
  const importRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    document.title = 'Posting....';
  }, []);

  // Guard: redirect if no token
  useEffect(() => {
    if (!localStorage.getItem('token')) {
      navigate('/admin');
    }
  }, [navigate]);

  // Load sidebar posts
  useEffect(() => {
    listPosts()
      .then(setSidebarPosts)
      .catch(() => {});
  }, []);

  // Load existing post when editing
  useEffect(() => {
    if (!editId) return;
    getPost(Number(editId)).then((post) => {
      if (!post) return;
      setTitle(post.title);
      setTagsInput(post.tags?.join(', ') ?? '');
      setContent(post.content);
      setCoverImageUrl(post.cover_image_url ?? '');
    });
  }, [editId]);

  function handleNewPost() {
    setTitle('');
    setTagsInput('');
    setContent('');
    setCoverImageUrl('');
    navigate('/write');
  }

  function parseTags(raw: string): string[] {
    return raw
      .split(',')
      .map((t) => t.trim())
      .filter(Boolean);
  }

  async function handlePublish() {
    setError('');
    setSaving(true);
    try {
      const data = {
        title,
        content,
        tags: parseTags(tagsInput),
        cover_image_url: coverImageUrl || undefined,
      };

      if (editId) {
        await updatePost(Number(editId), data);
      } else {
        await createPost(data);
      }
      navigate('/');
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete() {
    if (!editId) return;
    try {
      await deletePost(Number(editId));
      setSidebarPosts((sp) => sp.filter((p) => p.id !== Number(editId)));
      setTitle('');
      setTagsInput('');
      setContent('');
      setCoverImageUrl('');
      navigate('/write');
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setShowDeleteModal(false);
    }
  }

  function handleImport(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    const reader = new FileReader();
    reader.onload = (ev) => {
      const text = ev.target?.result;
      if (typeof text === 'string') setContent(text);
    };
    reader.readAsText(file);
    e.target.value = '';
  }

  const filename = title
    ? title.toLowerCase().replace(/\s+/g, '-') + '.md'
    : 'untitled.md';

  const tags = parseTags(tagsInput);

  return (
    <div className="write-layout" data-color-mode="dark">
      <div className="write-body">
        <div className="write-container">
          <CoverUpload value={coverImageUrl} onChange={setCoverImageUrl} />

          <div className="write-meta-row">
            <input
              className="write-title"
              placeholder="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
            />
            <input
              className="write-tags"
              placeholder="tags, comma, separated"
              value={tagsInput}
              onChange={(e) => setTagsInput(e.target.value)}
            />
          </div>

          {tags.length > 0 && (
            <div className="write-tag-pills">
              {tags.map((t) => (
                <span key={t} className="write-tag-pill">
                  {t}
                </span>
              ))}
            </div>
          )}

          <div className="write-editor">
            <MDEditor
              value={content}
              onChange={(v) => setContent(v ?? '')}
              preview={previewing ? 'preview' : 'edit'}
              height="100%"
              visibleDragbar={false}
            />
          </div>

          <div className="write-toolbar">
            {editId && (
              <button
                className="write-btn write-btn-delete"
                onClick={() => setShowDeleteModal(true)}
              >
                [ delete ]
              </button>
            )}
            <input
              ref={importRef}
              type="file"
              accept=".md"
              style={{ display: 'none' }}
              onChange={handleImport}
            />
            <button
              className="write-btn"
              onClick={() => importRef.current?.click()}
            >
              [ import .md ]
            </button>

            {error && <span className="write-error">{error}</span>}

            <div className="write-toolbar-right">
              <button
                className="write-btn write-btn-preview"
                onClick={() => setPreviewing((p) => !p)}
              >
                {previewing ? '[ edit ]' : '[ preview ]'}
              </button>
              <button
                className="write-btn write-btn-publish"
                onClick={handlePublish}
                disabled={saving || !title || !content}
              >
                {saving ? 'saving...' : editId ? '[ update ]' : '[ publish ]'}
              </button>
            </div>
          </div>
        </div>

        <div className="write-sidebar">
          <div className="write-sidebar-header" onClick={handleNewPost}>posts/</div>
          {sidebarPosts.map((post) => (
            <div
              key={post.id}
              className={`write-sidebar-row ${editId === String(post.id) ? 'active' : ''}`}
              onClick={() => navigate(`/write?id=${post.id}`)}
            >
              <span className="write-sidebar-title">{post.title}</span>
              <span className="write-sidebar-date">{post.created_at.slice(0, 10)}</span>
            </div>
          ))}
        </div>
      </div>

      {showDeleteModal && (
        <div className="delete-modal-overlay">
          <div className="delete-modal">
            <p className="delete-modal-question">Delete this post?</p>
            <p className="delete-modal-title">"{title}"</p>
            <div className="delete-modal-actions">
              <button
                className="delete-modal-cancel"
                onClick={() => setShowDeleteModal(false)}
              >
                [ cancel ]
              </button>
              <button className="delete-modal-confirm" onClick={handleDelete}>
                [ delete ]
              </button>
            </div>
          </div>
        </div>
      )}

      <StatusBar left={filename} center="Markdown" right="WRITE" />
    </div>
  );
}
