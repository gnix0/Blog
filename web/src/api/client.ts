const BASE = '/v1';

function getToken(): string | null {
  return localStorage.getItem('token');
}

async function request<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    ...(options.headers as Record<string, string>),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  if (!(options.body instanceof FormData)) {
    headers['Content-Type'] = 'application/json';
  }

  const res = await fetch(`${BASE}${path}`, { ...options, headers });

  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(body.error ?? res.statusText);
  }

  if (res.status === 204) return undefined as T;
  return res.json();
}

// Auth
export function login(username: string, password: string) {
  return request<{ token: string }>('/auth/token', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  });
}

// Posts
export interface Post {
  id: number;
  title: string;
  content: string;
  user_id: number;
  tags: string[];
  cover_image_url?: string;
  created_at: string;
  updated_at: string;
}

export function listPosts() {
  return request<Post[]>('/posts');
}

export function getPost(id: number) {
  return request<Post>(`/posts/${id}`);
}

export function createPost(data: {
  title: string;
  content: string;
  tags: string[];
  cover_image_url?: string;
}) {
  return request<Post>('/posts', { method: 'POST', body: JSON.stringify(data) });
}

export function updatePost(
  id: number,
  data: Partial<{ title: string; content: string; tags: string[]; cover_image_url: string }>
) {
  return request<Post>(`/posts/${id}`, { method: 'PATCH', body: JSON.stringify(data) });
}

export function deletePost(id: number) {
  return request<void>(`/posts/${id}`, { method: 'DELETE' });
}

// Uploads
export async function uploadImage(file: File): Promise<{ url: string }> {
  const form = new FormData();
  form.append('image', file);
  return request<{ url: string }>('/uploads', { method: 'POST', body: form });
}
