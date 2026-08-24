/* ==========================================================================
   Restgeld Landing Page - Interactivity & Live Calculations
   ========================================================================== */

document.addEventListener('DOMContentLoaded', () => {
  initCalculator();
  initMockupDemo();
});

/* --------------------------------------------------------------------------
   1. Interactive Budget Calculator
   -------------------------------------------------------------------------- */
function initCalculator() {
  const inputMonthly = document.getElementById('input-monthly-budget');
  const inputDays = document.getElementById('input-period-days');

  const valMonthly = document.getElementById('val-monthly-budget');
  const valDays = document.getElementById('val-period-days');

  const resDaily = document.getElementById('res-daily-budget');
  const resWeekly = document.getElementById('res-weekly-budget');
  const resSmartDaily = document.getElementById('res-smart-daily');

  if (!inputMonthly || !inputDays) return;

  function update() {
    const monthly = parseFloat(inputMonthly.value) || 0;
    const days = parseInt(inputDays.value, 10) || 30;

    // Update Slider Labels
    valMonthly.textContent = `${monthly.toLocaleString('de-DE')} €`;
    valDays.textContent = `${days} Tage`;

    // Calculate Metrics
    const daily = monthly / days;
    const weekly = daily * 7;

    // Smart Daily (assuming ~2 zero-spend days per 7 days)
    const zeroDaysFraction = (2 / 7) * days;
    const activeDays = Math.max(1, days - zeroDaysFraction);
    const smartDaily = monthly / activeDays;

    // Format Currencies
    resDaily.textContent = formatCurrency(daily);
    resWeekly.textContent = formatCurrency(weekly);
    resSmartDaily.textContent = formatCurrency(smartDaily);
  }

  inputMonthly.addEventListener('input', update);
  inputDays.addEventListener('input', update);

  // Initial Calculation
  update();
}

function formatCurrency(num) {
  return new Intl.NumberFormat('de-DE', {
    style: 'currency',
    currency: 'EUR',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(num);
}

/* --------------------------------------------------------------------------
   2. Interactive Phone Mockup Demo
   -------------------------------------------------------------------------- */
function initMockupDemo() {
  const amountEl = document.getElementById('mockup-amount');
  const badgeEl = document.getElementById('mockup-badge');
  const toastEl = document.getElementById('mockup-toast');
  const resetBtn = document.getElementById('mockup-reset');
  const chips = document.querySelectorAll('.mockup-chip');

  if (!amountEl || !badgeEl || !toastEl) return;

  const INITIAL_AMOUNT = 14.20;
  const INITIAL_SAVINGS = 42.50;
  let currentAmount = INITIAL_AMOUNT;
  let currentSavings = INITIAL_SAVINGS;
  let toastTimer = null;

  function updateMockupUI() {
    amountEl.textContent = formatCurrency(currentAmount);
    
    // Animate Number Pulse
    amountEl.classList.remove('pulsing');
    void amountEl.offsetWidth; // Force reflow
    amountEl.classList.add('pulsing');

    if (currentAmount < 0) {
      amountEl.classList.add('negative');
      badgeEl.className = 'hero-badge-pill badge-negative';
      badgeEl.innerHTML = `<span>${formatCurrency(currentSavings)} Puffer</span>`;
    } else {
      amountEl.classList.remove('negative');
      badgeEl.className = 'hero-badge-pill badge-positive';
      badgeEl.innerHTML = `<span>+${formatCurrency(currentSavings)} Spar-Puffer</span>`;
    }
  }

  function showToast(msg) {
    if (toastTimer) clearTimeout(toastTimer);
    toastEl.textContent = msg;
    toastEl.classList.add('active');

    toastTimer = setTimeout(() => {
      toastEl.classList.remove('active');
    }, 2000);
  }

  chips.forEach(chip => {
    chip.addEventListener('click', () => {
      const expense = parseFloat(chip.dataset.amount) || 0;
      const label = chip.dataset.label || 'Ausgabe';

      currentAmount -= expense;
      currentSavings -= expense;

      updateMockupUI();
      showToast(`${label} gebucht (-${formatCurrency(expense)})`);
    });
  });

  if (resetBtn) {
    resetBtn.addEventListener('click', () => {
      currentAmount = INITIAL_AMOUNT;
      currentSavings = INITIAL_SAVINGS;
      updateMockupUI();
      showToast('Demo zurückgesetzt ↺');
    });
  }
}
