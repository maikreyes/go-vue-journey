import { ApiURL, TopUrl } from '../config/config';
import type { Stock } from '../models/stock.model';



export async function fetchStocks(): Promise<Stock[]> {
    const response = await fetch(ApiURL);
    const data: Stock[] = await response.json();
    return data;
}

export async function fetchTopStocks(): Promise<Stock[]> {
    const response = await fetch(TopUrl);
    const data: Stock[] = await response.json();
    return data;
}