import { ApiURL, TopUrl } from '../config/config';
import type { Stock } from '../models/stock.model';
import type { StocksResponse } from '../models/stocksResponse.model';

export async function fetchStocks(limit: number, cursor?: string | null, filter: 'all' | 'up' | 'down' = 'all'): Promise<StocksResponse> {
    const url = new URL(ApiURL)
    url.searchParams.set('limit', String(limit))
    if (cursor) url.searchParams.set('cursor', cursor)
    url.searchParams.set('filter', filter)

    const response = await fetch(url);
    const data: StocksResponse = await response.json();
    return data;
}

export async function fetchTopStocks(): Promise<Stock[]> {
    const response = await fetch(TopUrl);
    const data: Stock[] = await response.json();
    return data;
}