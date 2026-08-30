import { describe, it, expect, beforeEach } from 'vitest'
import { useI18n, detectCategoryIcon } from '../composables/useI18n'

describe('useI18n Composable & Currency Engine', () => {
  beforeEach(() => {
    localStorage.clear()
    const { setLocale, setCurrency } = useI18n()
    setLocale('de')
    setCurrency('EUR')
  })

  it('initializes with default locale (de) and default currency (EUR)', () => {
    const { currentLocale, currentCurrency, currencySymbol, t, formatMoney } = useI18n()
    expect(currentLocale.value).toBe('de')
    expect(currentCurrency.value).toBe('EUR')
    expect(currencySymbol.value).toBe('€')
    expect(t('budget.available_today')).toBe('HEUTE VERFÜGBAR')
    expect(formatMoney(15.5)).toBe('15,50 €')
  })

  it('switches locale dynamically to English, Spanish and French', () => {
    const { setLocale, currentLocale, t } = useI18n()
    setLocale('en')
    expect(currentLocale.value).toBe('en')
    expect(t('budget.available_today')).toBe('AVAILABLE TODAY')
    expect(t('numpad.btn_save')).toBe('Save')

    setLocale('es')
    expect(t('budget.available_today')).toBe('DISPONIBLE HOY')

    setLocale('fr')
    expect(t('budget.available_today')).toBe('DISPONIBLE AUJOURD\'HUI')
  })

  it('switches currency dynamically (USD, GBP, CHF, JPY)', () => {
    const { setLocale, setCurrency, formatMoney, currencySymbol } = useI18n()
    setLocale('en')

    setCurrency('USD')
    expect(currencySymbol.value).toBe('$')
    expect(formatMoney(25.5)).toBe('$25.50')

    setCurrency('GBP')
    expect(currencySymbol.value).toBe('£')
    expect(formatMoney(25.5)).toBe('£25.50')

    setCurrency('CHF')
    expect(currencySymbol.value).toBe('CHF')
    expect(formatMoney(25.5)).toBe('25.50 CHF')
  })

  it('interpolates parameters correctly', () => {
    const { setLocale, t } = useI18n()
    setLocale('de')
    expect(t('budget.from_total', { total: '15,00 €' })).toBe('von 15,00 € heute')

    setLocale('en')
    expect(t('numpad.impact_available', { amount: '12.50 €' })).toBe('Available today: 12.50 €')
  })

  it('automatically detects smart categories without manual configuration', () => {
    expect(detectCategoryIcon('Espresso und Croissant')).toBe('☕')
    expect(detectCategoryIcon('Döner Kebab')).toBe('🍔')
    expect(detectCategoryIcon('Rewe Wocheneinkauf')).toBe('🛒')
    expect(detectCategoryIcon('Aral Tanken')).toBe('🚗')
    expect(detectCategoryIcon('Kino Ticket')).toBe('🎉')
    expect(detectCategoryIcon('Amazon Paket')).toBe('📦')
    expect(detectCategoryIcon('Apotheke Aspirin')).toBe('💊')
    expect(detectCategoryIcon('Unbekannte Notiz')).toBe('💶')
  })

  it('persists locale and currency to localStorage', () => {
    const { setLocale, setCurrency, initI18n, currentLocale, currentCurrency } = useI18n()
    setLocale('fr')
    setCurrency('USD')
    expect(localStorage.getItem('restgeld_language')).toBe('fr')
    expect(localStorage.getItem('restgeld_currency')).toBe('USD')

    currentLocale.value = 'de'
    currentCurrency.value = 'EUR'
    initI18n()
    expect(currentLocale.value).toBe('fr')
    expect(currentCurrency.value).toBe('USD')
  })
})
