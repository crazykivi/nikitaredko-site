const numericDate = new Intl.DateTimeFormat('ru-RU', {
  day: '2-digit',
  month: '2-digit',
  year: 'numeric',
})

const ruDate = new Intl.DateTimeFormat('ru-RU', {
  day: 'numeric',
  month: 'long',
  year: 'numeric',
})

export function formatDateNumeric(iso: string): string {
  const date = new Date(iso)
  return Number.isNaN(date.getTime()) ? '' : numericDate.format(date)
}

export function formatDateRu(iso: string): string {
  const date = new Date(iso)
  return Number.isNaN(date.getTime()) ? '' : ruDate.format(date)
}