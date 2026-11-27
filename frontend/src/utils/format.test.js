import { test } from 'node:test'
import assert from 'node:assert/strict'
import { deadlineInfo, formatRelative, formatSalary, initials, scoreLevel } from './format.js'

test('formatSalary', () => {
  assert.equal(formatSalary({ show_salary: true, salary_min: 5_000_000, salary_max: 8_000_000 }), 'Rp 5–8 jt/bulan')
  assert.equal(formatSalary({ show_salary: true, salary_min: 4_500_000, salary_max: null }), 'Rp 4,5 jt/bulan')
  assert.equal(formatSalary({ show_salary: true, salary_min: 6_000_000, salary_max: 6_000_000 }), 'Rp 6 jt/bulan')
  assert.equal(formatSalary({ show_salary: false, salary_min: 5_000_000, salary_max: 8_000_000 }), null)
  assert.equal(formatSalary({ show_salary: true, salary_min: null, salary_max: null }), null)
})

test('formatRelative', () => {
  const now = new Date(2026, 10, 24, 9)
  assert.equal(formatRelative(new Date(2026, 10, 24, 1), now), 'Hari ini')
  assert.equal(formatRelative(new Date(2026, 10, 21, 23), now), '3 hari lalu')
  assert.equal(formatRelative(new Date(2026, 10, 10), now), '2 minggu lalu')
})

test('deadlineInfo', () => {
  const now = new Date(2026, 10, 24, 15)
  assert.deepEqual(deadlineInfo('2026-11-23', now), { text: 'Sudah ditutup', urgent: false, days: -1 })
  assert.deepEqual(deadlineInfo('2026-11-24', now), { text: 'Ditutup hari ini', urgent: true, days: 0 })
  assert.equal(deadlineInfo('2026-11-27', now).urgent, true)
  assert.equal(deadlineInfo('2026-11-28', now).text, 'Ditutup dalam 4 hari')
})

test('initials dan scoreLevel', () => {
  assert.equal(initials('Nadia Putri Ayu'), 'NP')
  assert.equal(initials(''), '')
  assert.equal(scoreLevel(75), 'high')
  assert.equal(scoreLevel(50), 'medium')
  assert.equal(scoreLevel(49), 'low')
})
