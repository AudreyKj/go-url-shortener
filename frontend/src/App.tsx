

import { useState } from 'react';
import { shortenUrl } from './api/shorten';
import './styles/App.css';

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
    } catch (err: unknown) {
      if(err instanceof Error) setError(err.message);
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
    <div className="url-shortener">
      <h2 className="url-shortener__title">Minimal AI URL Shortener</h2>
      <form onSubmit={handleSubmit} className="url-shortener__form">
        <input
          type="url"
          className="url-shortener__input"
          placeholder="Paste your link here..."
          value={url}
          onChange={e => setUrl(e.target.value)}
          required
        />
        <button type="submit" className="url-shortener__button" disabled={loading}>
          {loading ? 'Shortening...' : 'Shorten'}
        </button>
      </form>
      {shortUrl && (
        <div className="url-shortener__result">
            <span className="url-shortener__result-label" data-testid="shortened-url-label">Shortened URL:</span>
            <div className="url-shortener__result-link">
              <a className="url-shortener__result-anchor" href={shortUrl} target="_blank" rel="noopener noreferrer">{shortUrl}</a>
              <button className={`url-shortener__copy-btn${copied ? ' url-shortener__copy-btn--copied' : ''}`} onClick={handleCopy} type="button">
                {copied ? 'Copied' : 'Copy'}
              </button>
            </div>
        </div>
      )}
      {error && <div className="url-shortener__error">{error}</div>}
    </div>
  );
}

export default App;
