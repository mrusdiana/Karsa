const dateFormat = new Intl.DateTimeFormat('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
const dateTimeFormat = new Intl.DateTimeFormat('id-ID', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
const numberFormat = new Intl.NumberFormat('id-ID')
const dayLength = 24 * 60 * 60 * 1000

const startOfDay = (value) => {
  const date = new Date(value)
  return new Date(date.getFullYear(), date.getMonth(), date.getDate())
}

const parseDay = (value) => (/^\d{4}-\d{2}-\d{2}$/.test(value) ? new Date(`${value}T00:00:00`) : new Date(value))

const millions = (value) => numberFormat.format(Math.round((value / 1_000_000) * 10) / 10)

export const formatNumber = (value) => numberFormat.format(Number(value) || 0)
export const formatDate = (value) => (value ? dateFormat.format(parseDay(value)) : '')
export const formatDateTime = (value) => (value ? dateTimeFormat.format(new Date(value)) : '')

export function formatSalary(job) {
  const values = [job.salary_min, job.salary_max].filter((value) => value != null)
  if (!job.show_salary || values.length === 0) return null
  const [low, high] = [Math.min(...values), Math.max(...values)]
  return low === high ? `Rp ${millions(low)} jt/bulan` : `Rp ${millions(low)}–${millions(high)} jt/bulan`
}

export function formatRelative(value, now = new Date()) {
  const days = Math.round((startOfDay(now) - startOfDay(value)) / dayLength)
  if (days <= 0) return 'Hari ini'
  if (days < 7) return `${days} hari lalu`
  if (days < 30) return `${Math.floor(days / 7)} minggu lalu`
  return formatDate(value)
}

export function deadlineInfo(deadline, now = new Date()) {
  const days = Math.round((parseDay(deadline) - startOfDay(now)) / dayLength)
  if (days < 0) return { text: 'Sudah ditutup', urgent: false, days }
  if (days === 0) return { text: 'Ditutup hari ini', urgent: true, days }
  return { text: `Ditutup dalam ${days} hari`, urgent: days <= 3, days }
}

export const initials = (name = '') =>
  name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((word) => word[0].toUpperCase())
    .join('')

export const scoreLevel = (score) => (score >= 75 ? 'high' : score >= 50 ? 'medium' : 'low')

export const SCORE_LABEL = { high: 'Cocok tinggi', medium: 'Cocok sedang', low: 'Cocok rendah' }

export const JOB_TYPE = { full_time: 'Penuh waktu', part_time: 'Paruh waktu', internship: 'Magang', contract: 'Kontrak' }
export const JOB_LEVEL = { internship: 'Magang', entry: 'Entry level', junior: 'Junior', mid: 'Mid level' }
export const WORK_MODE = { onsite: 'Di kantor', hybrid: 'Hybrid', remote: 'Remote' }
export const EDUCATION = ['SMA/SMK', 'D3', 'S1', 'S2']
export const COMPANY_SIZE = ['1-10', '11-50', '51-200', '201-500', '>500']

export const APPLICATION_STATUS = {
  applied: { label: 'Dilamar', icon: 'Send', tone: 'primary' },
  reviewed: { label: 'Ditinjau', icon: 'Eye', tone: 'info' },
  interview: { label: 'Interview', icon: 'CalendarClock', tone: 'warning' },
  accepted: { label: 'Diterima', icon: 'CheckCircle2', tone: 'success' },
  rejected: { label: 'Ditolak', icon: 'XCircle', tone: 'danger' },
  withdrawn: { label: 'Ditarik', icon: 'Undo2', tone: 'neutral' },
}

export const COMPANY_STATUS = {
  pending: { label: 'Menunggu verifikasi', icon: 'Clock', tone: 'warning' },
  verified: { label: 'Terverifikasi', icon: 'BadgeCheck', tone: 'success' },
  rejected: { label: 'Ditolak', icon: 'XCircle', tone: 'danger' },
}

export const INVITATION_STATUS = {
  pending: { label: 'Menunggu', icon: 'Clock', tone: 'warning' },
  accepted: { label: 'Sudah melamar', icon: 'Send', tone: 'info' },
  declined: { label: 'Ditolak', icon: 'XCircle', tone: 'neutral' },
  expired: { label: 'Kedaluwarsa', icon: 'Hourglass', tone: 'neutral' },
}

export const JOB_STATUS = {
  draft: { label: 'Draft', icon: 'FilePen', tone: 'neutral' },
  open: { label: 'Dibuka', icon: 'CheckCircle2', tone: 'success' },
  closed: { label: 'Ditutup', icon: 'Lock', tone: 'neutral' },
}

export const COMPANY_MOVES = {
  applied: ['reviewed', 'rejected'],
  reviewed: ['interview', 'rejected'],
  interview: ['accepted', 'rejected'],
}

export const MISSING_LABEL = {
  headline: 'Headline',
  location: 'Lokasi',
  education: 'Pendidikan',
  graduation_year: 'Tahun lulus',
  skills: 'Minimal 1 skill',
  cv: 'CV (PDF)',
}
