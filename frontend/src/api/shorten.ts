export const API_URL = 'http://localhost:8080/api/urls';

export async function shortenUrl(url: string): Promise<string> {
    try {

        const res = await fetch(API_URL, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ url }),
        });
        if (!res.ok) throw new Error('Failed to shorten URL');
        const data = await res.json();
        return data.short_url || data.shortUrl || '';

    } catch (error) {
        throw new Error("Failed to shorten URL");
    }
}
