import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import StreakCard from '../components/StreakCard.vue'

describe('StreakCard', () => {
  it('rendert nichts wenn streak prop nicht vorhanden ist', () => {
    const wrapper = mount(StreakCard, {
      props: { streak: undefined },
    })
    expect(wrapper.find('.streak-card').exists()).toBe(false)
  })

  it('rendert Streak-Daten und Flammen-Status korrekt', () => {
    const wrapper = mount(StreakCard, {
      props: {
        streak: {
          currentStreak: 4,
          longestStreak: 7,
          noSpendDays: 2,
          underBudgetDays: 5,
        },
      },
    })

    expect(wrapper.find('.streak-card').exists()).toBe(true)
    expect(wrapper.find('.streak-count').text()).toBe('4 Tage')
    expect(wrapper.find('.streak-flame-wrap').classes()).toContain('flame-active')

    const miniBadges = wrapper.findAll('.mini-badge')
    expect(miniBadges.length).toBe(2)
    // Null-Euro Tage: 2
    expect(miniBadges[0].find('.badge-val').text()).toBe('2')
    // Rekord: 7
    expect(miniBadges[1].find('.badge-val').text()).toBe('7')
  })

  it('zeigt Tag im Singular bei 1 Tag Streak', () => {
    const wrapper = mount(StreakCard, {
      props: {
        streak: {
          currentStreak: 1,
          longestStreak: 1,
          noSpendDays: 0,
          underBudgetDays: 1,
        },
      },
    })

    expect(wrapper.find('.streak-count').text()).toBe('1 Tag')
  })
})
