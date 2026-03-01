import { BrowserRouter, Route, Routes } from 'react-router-dom';
import PostList from './pages/PostList';
import PostReader from './pages/PostReader';
import Admin from './pages/Admin';
import Write from './pages/Write';

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<PostList />} />
        <Route path="/posts/:id/:slug" element={<PostReader />} />
        <Route path="/admin" element={<Admin />} />
        <Route path="/write" element={<Write />} />
      </Routes>
    </BrowserRouter>
  );
}
