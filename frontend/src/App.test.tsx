import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { vi } from 'vitest';
import App from './App';

// Mock the shortenUrl API using Vitest
vi.mock('./api/shorten', () => ({
  shortenUrl: async () => 'http://short.url/abc123',
}));

describe('App', () => {
  it('renders the URL shortener form', () => {
    render(<App />);
    expect(screen.getByPlaceholderText(/paste your link/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /shorten/i })).toBeInTheDocument();
  });

  it('shortens a URL and displays the result', async () => {
    render(<App />);
    const input = screen.getByPlaceholderText(/paste your link/i);
    const button = screen.getByRole('button', { name: /shorten/i });
    fireEvent.change(input, { target: { value: 'https://example.com' } });
    fireEvent.click(button);
    await waitFor(() => {
      expect(screen.getByText(/shortened url:/i)).toBeInTheDocument();
      expect(screen.getByText('http://short.url/abc123')).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /copy/i })).toBeInTheDocument();
    });
  });
});
