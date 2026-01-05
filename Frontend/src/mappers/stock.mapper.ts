import type { Stock } from '../models/stock.model';
import type { CardProps } from '../models/cardProps.model';

export function stockToCard(stock: Stock): CardProps {
  return {
    ticker: stock.ticker,
    targetFrom: stock.target_from,
    targetTo: stock.target_to,
    company: stock.company,
    action: stock.rating_to
  }
}
