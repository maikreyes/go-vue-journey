import { ApiURL, TopUrl } from '../config/config';
import type { Stock } from '../models/stock.model';
import type { StocksResponse } from '../models/stocksResponse.model';

export type StocksFilter = 'all' | 'up' | 'down' | 'equal'

export async function fetchStocks(options?: {
    nextPage?: string | null
    filter?: StocksFilter
    ticker?: string | null
}): Promise<StocksResponse> {
    const url = new URL(ApiURL)
    const nextPage = options?.nextPage ?? null
    const filter = options?.filter ?? 'all'
    const ticker = options?.ticker ?? null

    if (nextPage) url.searchParams.set('next_page', nextPage)
    if (filter !== 'all') url.searchParams.set('filter', filter)
    if (ticker) url.searchParams.set('ticker', ticker)

    const response = await fetch(url);
    if (!response.ok) {
        throw new Error(`Failed to fetch stocks: ${response.status}`)
    }
    const data: StocksResponse = await response.json();
    return data;
}

export async function fetchTopStocks(): Promise<Stock[]> {
    const response = await fetch(TopUrl);
    if (!response.ok) {
        throw new Error(`Failed to fetch top stocks: ${response.status}`)
    }
    const data: Stock[] = await response.json();
    return data;
}