import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import SpendingTrend from '../components/SpendingTrend.vue'

describe('SpendingTrend', () => {
  it('rendert nichts wenn keine stats vorhanden sind', () => {
    const wrapper = mount(SpendingTrend, {
      props: { stats: [], baseBudget: 15, currentDay: 1 },
    })
    expect(wrapper.find('.spending-trend').exists()).toBe(false)
  })

  it('rendert Balken fuer jeden Tag und berechnet Tagesschnitt', () => {
    const sampleStats = [
      { day: 1, date: '2026-08-01', spent: 10 },
      { day: 2, date: '2026-08-02', spent: 0 },
      { day: 3, date: '2026-08-03', spent: 20 },
    ]

    const wrapper = mount(SpendingTrend, {
      props: { stats: sampleStats, baseBudget: 15, currentDay: 3 },
    })

    expect(wrapper.find('.spending-trend').exists()).toBe(true)
    expect(wrapper.find('.trend-title').text()).toBe('Tages-Verlauf')

    const columns = wrapper.findAll('.bar-column')
    expect(columns.length).toBe(3)

    // Tagesschnitt: (10 + 0 + 20) / 3 = 10,00 €
    expect(wrapper.find('.average-badge').text()).toContain('10,00')

    // Tag 2 ist 0 € (bar-zero)
    expect(columns[1].find('.bar-fill').classes()).toContain('bar-zero')
    // Tag 1 <= 15 € (bar-good)
    expect(columns[0].find('.bar-fill').classes()).toContain('bar-good')
    // Tag 3 > 15 € (bar-over)
    expect(columns[2].find('.bar-fill').classes()).toContain('bar-over')
  })

  it('zeigt Detail-Informationen bei Klick auf einen Tag', async () => {
    const sampleStats = [
      { day: 1, date: '2026-08-01', spent: 8.5 },
      { day: 2, date: '2026-08-02', spent: 0 },
    ]

    const wrapper = mount(SpendingTrend, {
      props: { stats: sampleStats, baseBudget: 15, currentDay: 2 },
    })

    const columns = wrapper.findAll('.bar-column')
    await columns[0].trigger('click')

    expect(wrapper.find('.detail-day').text()).toContain('Tag 1 (01.08.):')
    expect(wrapper.find('.detail-spent').text()).toContain('8,50 €')

    // Klick auf Tag 2 (Spar-Tag)
    await columns[1].trigger('click')
    expect(wrapper.find('.detail-day').text()).toContain('Tag 2 (02.08.):')
    expect(wrapper.find('.detail-spent').text()).toContain('Spar-Tag!')
  })
})
