import { ERROR_GENERIC, ERROR_INVALID_URL } from '../utils/errors';

const API_URL = `${import.meta.env.VITE_API_URL}/api/urls`;

interface ShortenResponse {
    original_url: string;
    short_code: string;
    short_url: string;
    slug_type: string;
    error?: string;
}

export async function shortenUrl(url: string): Promise<string> {
    try {
        const res = await fetch(API_URL, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ url }),
        });
        const data: ShortenResponse = await res.json();
        if (!res.ok) {
            throw new Error(data.error === 'Invalid URL' ? ERROR_INVALID_URL : ERROR_GENERIC);
        }
        return data.short_url || '';
    } catch (error: any) {
        throw new Error(error.message ?? ERROR_GENERIC);
    }
}
