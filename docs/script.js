/* ==========================================================================
   Restgeld Landing Page - Interactivity & Live Calculations
   ========================================================================== */

let currentSelectedCurrency = 'EUR'

document.addEventListener('DOMContentLoaded', () => {
  initCurrencySelector()
  initCalculator()
  initMockupDemo()
})

/* --------------------------------------------------------------------------
   0. Global Multi-Currency Support
   -------------------------------------------------------------------------- */
function initCurrencySelector() {
  const currencyButtons = document.querySelectorAll('.currency-chip')
  currencyButtons.forEach((btn) => {
    btn.addEventListener('click', () => {
      currencyButtons.forEach((b) => b.classList.remove('active'))
      btn.classList.add('active')
      currentSelectedCurrency = btn.dataset.currency || 'EUR'

      // Trigger recalculations
      if (typeof window.updateCalculator === 'function') {
        window.updateCalculator()
      }
      if (typeof window.updateMockupUI === 'function') {
        window.updateMockupUI()
      }
    })
  })
}

function formatMoneyLanding(num, currency = currentSelectedCurrency) {
  const locale = localStorage.getItem('restgeld_language') === 'en' ? 'en-US' : 'de-DE'
  const digits = currency === 'JPY' ? 0 : 2
  const formatted = num.toLocaleString(locale, {
    minimumFractionDigits: digits,
    maximumFractionDigits: digits,
  })

  if (currency === 'USD') return `$${formatted}`
  if (currency === 'GBP') return `£${formatted}`
  if (currency === 'JPY') return `¥${formatted}`
  if (currency === 'CHF') return `${formatted} CHF`
  return `${formatted} €`
}

/* --------------------------------------------------------------------------
   1. Interactive Budget Calculator
   -------------------------------------------------------------------------- */
function initCalculator() {
  const inputMonthly = document.getElementById('input-monthly-budget')
  const inputDays = document.getElementById('input-period-days')

  const valMonthly = document.getElementById('val-monthly-budget')
  const valDays = document.getElementById('val-period-days')

  const resDaily = document.getElementById('res-daily-budget')
  const resWeekly = document.getElementById('res-weekly-budget')
  const resSmartDaily = document.getElementById('res-smart-daily')

  if (!inputMonthly || !inputDays) return

  function update() {
    const monthly = parseFloat(inputMonthly.value) || 0
    const days = parseInt(inputDays.value, 10) || 30

    // Update Slider Labels
    valMonthly.textContent = formatMoneyLanding(monthly)
    valDays.textContent = `${days} Tage`

    // Calculate Metrics
    const daily = monthly / days
    const weekly = daily * 7

    // Smart Daily (assuming ~2 zero-spend days per 7 days)
    const zeroDaysFraction = (2 / 7) * days
    const activeDays = Math.max(1, days - zeroDaysFraction)
    const smartDaily = monthly / activeDays

    // Format Currencies
    resDaily.textContent = formatMoneyLanding(daily)
    resWeekly.textContent = formatMoneyLanding(weekly)
    resSmartDaily.textContent = formatMoneyLanding(smartDaily)
  }

  window.updateCalculator = update
  inputMonthly.addEventListener('input', update)
  inputDays.addEventListener('input', update)

  // Initial Calculation
  update()
}

/* --------------------------------------------------------------------------
   2. Interactive Phone Mockup Demo
   -------------------------------------------------------------------------- */
function initMockupDemo() {
  const amountEl = document.getElementById('mockup-amount')
  const badgeEl = document.getElementById('mockup-badge')
  const toastEl = document.getElementById('mockup-toast')
  const resetBtn = document.getElementById('mockup-reset')
  const chips = document.querySelectorAll('.mockup-chip')

  if (!amountEl || !badgeEl || !toastEl) return

  const INITIAL_AMOUNT = 14.2
  const INITIAL_SAVINGS = 42.5
  let currentAmount = INITIAL_AMOUNT
  let currentSavings = INITIAL_SAVINGS
  let toastTimer = null

  function updateMockupUI() {
    amountEl.textContent = formatMoneyLanding(currentAmount)

    // Animate Number Pulse
    amountEl.classList.remove('pulsing')
    void amountEl.offsetWidth // Force reflow
    amountEl.classList.add('pulsing')

    if (currentAmount < 0) {
      amountEl.classList.add('negative')
      badgeEl.className = 'hero-badge-pill badge-negative'
      badgeEl.innerHTML = `<span>${formatMoneyLanding(currentSavings)} Puffer</span>`
    } else {
      amountEl.classList.remove('negative')
      badgeEl.className = 'hero-badge-pill badge-positive'
      badgeEl.innerHTML = `<span>+${formatMoneyLanding(currentSavings)} Spar-Puffer</span>`
    }
  }

  window.updateMockupUI = updateMockupUI

  function showToast(msg) {
    if (toastTimer) clearTimeout(toastTimer)
    toastEl.textContent = msg
    toastEl.classList.add('active')

    toastTimer = setTimeout(() => {
      toastEl.classList.remove('active')
    }, 2000)
  }

  chips.forEach((chip) => {
    chip.addEventListener('click', () => {
      const expense = parseFloat(chip.dataset.amount) || 0
      const label = chip.dataset.label || 'Ausgabe'

      currentAmount -= expense
      currentSavings -= expense

      updateMockupUI()
      showToast(`${label} gebucht (-${formatMoneyLanding(expense)})`)
    })
  })

  if (resetBtn) {
    resetBtn.addEventListener('click', () => {
      currentAmount = INITIAL_AMOUNT
      currentSavings = INITIAL_SAVINGS
      updateMockupUI()
      showToast('Demo zurückgesetzt ↺')
    })
  }
}
