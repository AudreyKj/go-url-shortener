import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { ERROR_INVALID_URL } from '../utils/errors';
import { shortenUrl } from '../api/shorten';
import { vi } from 'vitest';
import type { MockInstance } from 'vitest';
import App from '../App';


const MOCK_URL = 'https://example.com/long-url/path/abc?query=123';
const MOCK_SHORT_URL = 'http://short.url/abc123';
const MOCK_ERROR = ERROR_INVALID_URL;

vi.mock('../api/shorten', () => ({
  shortenUrl: vi.fn(),
}));

describe('App', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders the URL shortener form', () => {
    render(<App />);
    expect(screen.getByPlaceholderText(/paste your link/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /shorten/i })).toBeInTheDocument();
  });

  it('shortens a URL and displays the result', async () => {
  (shortenUrl as unknown as MockInstance).mockResolvedValueOnce(MOCK_SHORT_URL);

    render(<App />);
    const input = screen.getByPlaceholderText(/paste your link/i);
    const button = screen.getByRole('button', { name: /shorten/i });

    fireEvent.change(input, { target: { value: MOCK_URL } });
    fireEvent.click(button);

    await waitFor(() => {
      expect(screen.getByTestId('shortened-url-label')).toBeInTheDocument();
      expect(screen.getByText(MOCK_SHORT_URL)).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /copy/i })).toBeInTheDocument();
    });
  });

  it('displays an error message when API fails', async () => {
  (shortenUrl as unknown as MockInstance).mockRejectedValueOnce(new Error(MOCK_ERROR));

    render(<App />);
    const input = screen.getByPlaceholderText(/paste your link/i);
    const button = screen.getByRole('button', { name: /shorten/i });

    fireEvent.change(input, { target: { value: MOCK_URL } });
    fireEvent.click(button);

    expect(await screen.findByText(/invalid url/i)).toBeInTheDocument();
  });
});

