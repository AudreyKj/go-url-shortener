

import { useState } from 'react';
import { shortenUrl } from './api/shorten';
import './App.css';

function App() {
  const [url, setUrl] = useState('');
  const [shortUrl, setShortUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [copied, setCopied] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setShortUrl('');
    setCopied(false);
    try {
      const result = await shortenUrl(url);
      setShortUrl(result);
    } catch (err: any) {
      console.log('err', err)
      setError(err.message || 'Error occurred');
    } finally {
      setLoading(false);
    }
  };

  const handleCopy = async () => {
    if (shortUrl) {
      try {
        await navigator.clipboard.writeText(shortUrl);
        setCopied(true);
        setTimeout(() => setCopied(false), 1500);
      } catch {
        setCopied(false);
      }
    }
  };

  return (
    <div className="container">
      <h2>Minimal AI URL Shortener</h2>
      <form onSubmit={handleSubmit} className="shorten-form">
        <input
          type="url"
          placeholder="Paste your link here..."
          value={url}
          onChange={e => setUrl(e.target.value)}
          required
        />
        <button type="submit" disabled={loading}>
          {loading ? 'Shortening...' : 'Shorten'}
        </button>
      </form>
      {shortUrl && (
        <div className="result">
            <span>Shortened URL:</span>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginTop: '0.5rem' }}>
              <a href={shortUrl} target="_blank" rel="noopener noreferrer">{shortUrl}</a>
              <button className="copy-btn" onClick={handleCopy} type="button">
                {copied ? 'Copied!' : 'Copy'}
              </button>
            </div>
        </div>
      )}
      {error && <div className="error">{error}</div>}
    </div>
  );
}

export default App;
