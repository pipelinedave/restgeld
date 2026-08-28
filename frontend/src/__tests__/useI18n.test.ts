import { describe, it, expect, beforeEach } from 'vitest'
import { useI18n } from '../composables/useI18n'

describe('useI18n Composable', () => {
  beforeEach(() => {
    localStorage.clear()
    const { setLocale } = useI18n()
    setLocale('de')
  })

  it('initializes with default locale (de)', () => {
    const { currentLocale, t } = useI18n()
    expect(currentLocale.value).toBe('de')
    expect(t('budget.available_today')).toBe('HEUTE VERFÜGBAR')
  })

  it('switches locale dynamically to English', () => {
    const { setLocale, currentLocale, t } = useI18n()
    setLocale('en')
    expect(currentLocale.value).toBe('en')
    expect(t('budget.available_today')).toBe('AVAILABLE TODAY')
    expect(t('numpad.btn_save')).toBe('Save')
  })

  it('switches locale dynamically to Spanish and French', () => {
    const { setLocale, t } = useI18n()
    setLocale('es')
    expect(t('budget.available_today')).toBe('DISPONIBLE HOY')

    setLocale('fr')
    expect(t('budget.available_today')).toBe('DISPONIBLE AUJOURD\'HUI')
  })

  it('interpolates parameters correctly', () => {
    const { setLocale, t } = useI18n()
    setLocale('de')
    expect(t('budget.from_total', { total: '15,00' })).toBe('von 15,00 € heute')

    setLocale('en')
    expect(t('numpad.impact_available', { amount: '12.50' })).toBe('Available today: 12.50 €')
  })

  it('persists locale selection to localStorage', () => {
    const { setLocale, initI18n, currentLocale } = useI18n()
    setLocale('fr')
    expect(localStorage.getItem('restgeld_language')).toBe('fr')

    currentLocale.value = 'de'
    initI18n()
    expect(currentLocale.value).toBe('fr')
  })
})
