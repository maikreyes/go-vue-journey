export function parseMoney(value?: string | number | null): number {
  if (typeof value === 'number') return value
  if (!value) return 0

  const cleaned = value
    .replace(/\$/g, '')
    .replace(/,/g, '')
    .trim()

  const parsed = Number(cleaned)
  return isNaN(parsed) ? 0 : parsed
}
