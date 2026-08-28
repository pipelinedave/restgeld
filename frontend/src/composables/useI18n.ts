import { ref, computed } from 'vue'

export type SupportedLocale = 'de' | 'en' | 'es' | 'fr'

export interface LanguageOption {
  code: SupportedLocale
  name: string
  flag: string
}

export const SUPPORTED_LANGUAGES: LanguageOption[] = [
  { code: 'de', name: 'Deutsch', flag: '🇩🇪' },
  { code: 'en', name: 'English', flag: '🇬🇧' },
  { code: 'es', name: 'Español', flag: '🇪🇸' },
  { code: 'fr', name: 'Français', flag: '🇫🇷' },
]

const STORAGE_KEY = 'restgeld_language'

export const translations: Record<SupportedLocale, Record<string, string>> = {
  de: {
    // Header
    'header.online': 'Online',
    'header.offline': 'Offline',
    'header.sync_pending': '{count} ausstehend',

    // Hero & Budget
    'budget.available_today': 'HEUTE VERFÜGBAR',
    'budget.savings': 'Ersparnis',
    'budget.deficit': 'Überzug',
    'budget.from_total': 'von {total} € heute',

    // Numpad / Book Expense
    'numpad.title': 'Ausgabe buchen',
    'numpad.amount_label': 'Betrag (€)',
    'numpad.note_label': 'Notiz (optional)',
    'numpad.note_placeholder': 'Notiz (z. B. Kaffee, Mittagessen)',
    'numpad.btn_cancel': 'Abbrechen',
    'numpad.btn_save': 'Speichern',
    'numpad.btn_saving': 'Wird gebucht...',
    'numpad.impact_available': 'Heute verfügbar: {amount} €',
    'numpad.impact_remaining': 'Heute verfügbar: {current} € ➔ Verbleibt danach: {diff} €',
    'numpad.impact_exceeds': 'Heute verfügbar: {current} € ➔ Überzieht Tagesbudget um {diff} €',

    // Settings
    'settings.title': 'Einstellungen',
    'settings.monthly_budget': 'Monatsbudget',
    'settings.period_days': 'Periodendauer (Tage)',
    'settings.desired_daily': 'Wunsch Tages-Budget',
    'settings.calculated_monthly': 'Berechnetes Monatsbudget:',
    'settings.save_btn': 'Einstellungen speichern',
    'settings.saved_msg': '✓ Einstellungen gespeichert',
    'settings.reset_heading': 'Neue Periode starten',
    'settings.reset_desc': 'Setzt den Tag auf 1 zurück und löscht alle Ausgaben der aktuellen Periode.',
    'settings.reset_btn': 'Neue Periode starten',
    'settings.reset_confirm_title': 'Sicher?',
    'settings.reset_confirm_body': 'Alle Ausgaben der aktuellen Periode werden gelöscht!',
    'settings.reset_confirm_btn': 'Ja, zurücksetzen',
    'settings.reset_cancel_btn': 'Abbrechen',
    'settings.theme_heading': 'Design & Farbwelt',
    'settings.language_heading': 'Sprache / Language',
    'settings.backup_heading': 'Daten & Backup',
    'settings.export_json': 'JSON-Backup exportieren',
    'settings.export_csv': 'CSV exportieren',
    'settings.import_backup': 'Backup importieren',

    // Expenses Modal
    'expenses.title': 'Ausgaben-Historie',
    'expenses.empty': 'Keine Ausgaben in dieser Periode vorhanden.',
    'expenses.page_info': 'Seite {page} von {total}',
    'expenses.delete': 'Löschen',
    'expenses.all_expenses': 'Alle anzeigen',
    'expenses.recent_title': 'Letzte Ausgaben',

    // Archive Modal
    'archive.title': 'Perioden-Archiv',
    'archive.subtitle': 'Monatsberichte & historische Ausgaben',
    'archive.empty': 'Noch keine vergangenen Perioden archiviert.',
    'archive.total_spent': 'Gesamtausgaben:',
    'archive.savings': 'Ersparnis:',
    'archive.count': 'Ausgaben:',
    'archive.avg_daily': 'Ø pro Tag:',
    'archive.view_report': 'Abschlussbericht ansehen',
    'archive.back_to_list': '← Zurück zum Archiv',

    // Streaks & Projection
    'streak.title': 'Spar-Streak',
    'streak.current': 'Aktuelle Streak',
    'streak.longest': 'Rekord-Streak',
    'streak.zero_days': 'Null-Euro-Tage',
    'streak.days_unit': 'Tage',
    'projection.title': 'Monatsende-Prognose',
    'projection.savings': 'Voraussichtliche Ersparnis:',
    'projection.deficit': 'Voraussichtlicher Überzug:',
    'projection.total': 'Voraussichtliche Gesamtausgaben:',
    'projection.daily_avg': 'Ø Tagesausgabe bisher:',

    // Spending Trend
    'trend.title': 'Ausgabenverlauf',
    'trend.legend_ok': 'Im Budget',
    'trend.legend_savings': 'Spar-Tag (€ 0)',
    'trend.legend_over': 'Über Budget',
    'trend.avg': 'Ø {amount} € / Tag',

    // Auth & SaaS
    'auth.title': 'Anmelden / Registrieren',
    'auth.subtitle': 'Verbinde dein Konto für Cloud-Sync auf allen Geräten.',
    'auth.email_label': 'E-Mail-Adresse',
    'auth.email_placeholder': 'deine@email.de',
    'auth.send_link': 'Magic Link senden',
    'auth.sending': 'Sende Link...',
    'auth.logged_in_as': 'Angemeldet als {email}',
    'auth.logout': 'Abmelden',
    'auth.delete_account': 'Konto löschen',
    'auth.pro_badge': 'PRO TIER',
    'auth.upgrade_pro': 'Auf Pro upgraden (Stripe)',
    'auth.manage_sub': 'Abonnement verwalten',

    // About
    'about.title': 'Über Restgeld',
    'about.tagline': 'Minimalistischer Daily Allowance Tracker',
    'about.philosophy': 'Achtsames Sparen ohne Schnickschnack. Fokus auf das wesentliche Tagesbudget.',
    'about.open_source': 'Open Source Software (MIT Lizenz)',
    'about.close': 'Schließen',

    // Footer & Misc
    'footer.tagline': 'Jeden Tag achtsam sparen.',
    'common.close': 'Schließen',
    'common.back': 'Zurück',
  },
  en: {
    // Header
    'header.online': 'Online',
    'header.offline': 'Offline',
    'header.sync_pending': '{count} pending',

    // Hero & Budget
    'budget.available_today': 'AVAILABLE TODAY',
    'budget.savings': 'Savings',
    'budget.deficit': 'Deficit',
    'budget.from_total': 'from {total} € today',

    // Numpad / Book Expense
    'numpad.title': 'Add Expense',
    'numpad.amount_label': 'Amount (€)',
    'numpad.note_label': 'Note (optional)',
    'numpad.note_placeholder': 'Note (e.g. Coffee, Lunch)',
    'numpad.btn_cancel': 'Cancel',
    'numpad.btn_save': 'Save',
    'numpad.btn_saving': 'Saving...',
    'numpad.impact_available': 'Available today: {amount} €',
    'numpad.impact_remaining': 'Available today: {current} € ➔ Remaining after: {diff} €',
    'numpad.impact_exceeds': 'Available today: {current} € ➔ Exceeds daily budget by {diff} €',

    // Settings
    'settings.title': 'Settings',
    'settings.monthly_budget': 'Monthly Budget',
    'settings.period_days': 'Period Duration (Days)',
    'settings.desired_daily': 'Desired Daily Budget',
    'settings.calculated_monthly': 'Calculated Monthly Total:',
    'settings.save_btn': 'Save Settings',
    'settings.saved_msg': '✓ Settings saved',
    'settings.reset_heading': 'Start New Period',
    'settings.reset_desc': 'Resets day to 1 and clears all expenses of the current period.',
    'settings.reset_btn': 'Start New Period',
    'settings.reset_confirm_title': 'Are you sure?',
    'settings.reset_confirm_body': 'All expenses of the current period will be deleted!',
    'settings.reset_confirm_btn': 'Yes, Reset',
    'settings.reset_cancel_btn': 'Cancel',
    'settings.theme_heading': 'Design & Color Theme',
    'settings.language_heading': 'Language / Sprache',
    'settings.backup_heading': 'Data & Backup',
    'settings.export_json': 'Export JSON Backup',
    'settings.export_csv': 'Export CSV',
    'settings.import_backup': 'Import Backup',

    // Expenses Modal
    'expenses.title': 'Expense History',
    'expenses.empty': 'No expenses recorded in this period.',
    'expenses.page_info': 'Page {page} of {total}',
    'expenses.delete': 'Delete',
    'expenses.all_expenses': 'View All',
    'expenses.recent_title': 'Recent Expenses',

    // Archive Modal
    'archive.title': 'Period Archive',
    'archive.subtitle': 'Monthly reports & historical expenses',
    'archive.empty': 'No archived periods yet.',
    'archive.total_spent': 'Total Spent:',
    'archive.savings': 'Savings:',
    'archive.count': 'Expenses:',
    'archive.avg_daily': 'Ø per day:',
    'archive.view_report': 'View Summary Report',
    'archive.back_to_list': '← Back to Archive',

    // Streaks & Projection
    'streak.title': 'Saving Streak',
    'streak.current': 'Current Streak',
    'streak.longest': 'Best Streak',
    'streak.zero_days': 'Zero-Spend Days',
    'streak.days_unit': 'Days',
    'projection.title': 'Month-End Projection',
    'projection.savings': 'Expected Savings:',
    'projection.deficit': 'Expected Deficit:',
    'projection.total': 'Expected Total Spent:',
    'projection.daily_avg': 'Ø Daily Spend so far:',

    // Spending Trend
    'trend.title': 'Spending Trend',
    'trend.legend_ok': 'Within Budget',
    'trend.legend_savings': 'Zero Spend (€ 0)',
    'trend.legend_over': 'Over Budget',
    'trend.avg': 'Ø {amount} € / day',

    // Auth & SaaS
    'auth.title': 'Sign In / Register',
    'auth.subtitle': 'Connect your account for cloud sync across all devices.',
    'auth.email_label': 'Email Address',
    'auth.email_placeholder': 'your@email.com',
    'auth.send_link': 'Send Magic Link',
    'auth.sending': 'Sending link...',
    'auth.logged_in_as': 'Signed in as {email}',
    'auth.logout': 'Sign Out',
    'auth.delete_account': 'Delete Account',
    'auth.pro_badge': 'PRO TIER',
    'auth.upgrade_pro': 'Upgrade to Pro (Stripe)',
    'auth.manage_sub': 'Manage Subscription',

    // About
    'about.title': 'About Restgeld',
    'about.tagline': 'Minimalist Daily Allowance Tracker',
    'about.philosophy': 'Mindful saving without unnecessary bloat. Pure focus on your daily allowance.',
    'about.open_source': 'Open Source Software (MIT License)',
    'about.close': 'Close',

    // Footer & Misc
    'footer.tagline': 'Save mindfully every day.',
    'common.close': 'Close',
    'common.back': 'Back',
  },
  es: {
    // Header
    'header.online': 'En línea',
    'header.offline': 'Sin conexión',
    'header.sync_pending': '{count} pendiente',

    // Hero & Budget
    'budget.available_today': 'DISPONIBLE HOY',
    'budget.savings': 'Ahorro',
    'budget.deficit': 'Déficit',
    'budget.from_total': 'de {total} € hoy',

    // Numpad / Book Expense
    'numpad.title': 'Añadir Gasto',
    'numpad.amount_label': 'Monto (€)',
    'numpad.note_label': 'Nota (opcional)',
    'numpad.note_placeholder': 'Nota (ej. Café, Almuerzo)',
    'numpad.btn_cancel': 'Cancelar',
    'numpad.btn_save': 'Guardar',
    'numpad.btn_saving': 'Guardando...',
    'numpad.impact_available': 'Disponible hoy: {amount} €',
    'numpad.impact_remaining': 'Disponible hoy: {current} € ➔ Quedará después: {diff} €',
    'numpad.impact_exceeds': 'Disponible hoy: {current} € ➔ Excede el presupuesto diario por {diff} €',

    // Settings
    'settings.title': 'Ajustes',
    'settings.monthly_budget': 'Presupuesto Mensual',
    'settings.period_days': 'Duración del Período (Días)',
    'settings.desired_daily': 'Presupuesto Diario Deseado',
    'settings.calculated_monthly': 'Total Mensual Calculado:',
    'settings.save_btn': 'Guardar Ajustes',
    'settings.saved_msg': '✓ Ajustes guardados',
    'settings.reset_heading': 'Iniciar Nuevo Período',
    'settings.reset_desc': 'Reinicia el día a 1 y borra todos los gastos del período actual.',
    'settings.reset_btn': 'Iniciar Nuevo Período',
    'settings.reset_confirm_title': '¿Estás seguro?',
    'settings.reset_confirm_body': '¡Se borrarán todos los gastos del período actual!',
    'settings.reset_confirm_btn': 'Sí, Reiniciar',
    'settings.reset_cancel_btn': 'Cancelar',
    'settings.theme_heading': 'Diseño y Color',
    'settings.language_heading': 'Idioma / Language',
    'settings.backup_heading': 'Datos y Copia de Seguridad',
    'settings.export_json': 'Exportar Copia JSON',
    'settings.export_csv': 'Exportar CSV',
    'settings.import_backup': 'Importar Copia',

    // Expenses Modal
    'expenses.title': 'Historial de Gastos',
    'expenses.empty': 'No hay gastos registrados en este período.',
    'expenses.page_info': 'Página {page} de {total}',
    'expenses.delete': 'Eliminar',
    'expenses.all_expenses': 'Ver todos',
    'expenses.recent_title': 'Gastos Recientes',

    // Archive Modal
    'archive.title': 'Archivo de Períodos',
    'archive.subtitle': 'Informes mensuales y gastos históricos',
    'archive.empty': 'Aún no hay períodos archivados.',
    'archive.total_spent': 'Gasto Total:',
    'archive.savings': 'Ahorro:',
    'archive.count': 'Gastos:',
    'archive.avg_daily': 'Ø por día:',
    'archive.view_report': 'Ver Informe Resumido',
    'archive.back_to_list': '← Volver al Archivo',

    // Streaks & Projection
    'streak.title': 'Racha de Ahorro',
    'streak.current': 'Racha Actual',
    'streak.longest': 'Mejor Racha',
    'streak.zero_days': 'Días de Cero Gastos',
    'streak.days_unit': 'Días',
    'projection.title': 'Proyección de Fin de Mes',
    'projection.savings': 'Ahorro Previsto:',
    'projection.deficit': 'Déficit Previsto:',
    'projection.total': 'Gasto Total Previsto:',
    'projection.daily_avg': 'Ø Gasto Diario hasta ahora:',

    // Spending Trend
    'trend.title': 'Tendencia de Gastos',
    'trend.legend_ok': 'En Presupuesto',
    'trend.legend_savings': 'Cero Gastos (€ 0)',
    'trend.legend_over': 'Sobre Presupuesto',
    'trend.avg': 'Ø {amount} € / día',

    // Auth & SaaS
    'auth.title': 'Iniciar Sesión / Registro',
    'auth.subtitle': 'Conecta tu cuenta para sincronizar en la nube en todos tus dispositivos.',
    'auth.email_label': 'Correo Electrónico',
    'auth.email_placeholder': 'tu@email.com',
    'auth.send_link': 'Enviar Magic Link',
    'auth.sending': 'Enviando enlace...',
    'auth.logged_in_as': 'Sesión iniciada como {email}',
    'auth.logout': 'Cerrar Sesión',
    'auth.delete_account': 'Eliminar Cuenta',
    'auth.pro_badge': 'NIVEL PRO',
    'auth.upgrade_pro': 'Mejorar a Pro (Stripe)',
    'auth.manage_sub': 'Gestionar Suscripción',

    // About
    'about.title': 'Acerca de Restgeld',
    'about.tagline': 'Rastreador Minimalista de Dinero Diario',
    'about.philosophy': 'Ahorro consciente sin complicaciones. Enfocado puramente en tu asignación diaria.',
    'about.open_source': 'Software de Código Abierto (Licencia MIT)',
    'about.close': 'Cerrar',

    // Footer & Misc
    'footer.tagline': 'Ahorra de forma consciente cada día.',
    'common.close': 'Cerrar',
    'common.back': 'Volver',
  },
  fr: {
    // Header
    'header.online': 'En ligne',
    'header.offline': 'Hors ligne',
    'header.sync_pending': '{count} en attente',

    // Hero & Budget
    'budget.available_today': 'DISPONIBLE AUJOURD\'HUI',
    'budget.savings': 'Économies',
    'budget.deficit': 'Déficit',
    'budget.from_total': 'sur {total} € aujourd\'hui',

    // Numpad / Book Expense
    'numpad.title': 'Ajouter une Dépense',
    'numpad.amount_label': 'Montant (€)',
    'numpad.note_label': 'Note (optionnel)',
    'numpad.note_placeholder': 'Note (ex. Café, Déjeuner)',
    'numpad.btn_cancel': 'Annuler',
    'numpad.btn_save': 'Enregistrer',
    'numpad.btn_saving': 'Enregistrement...',
    'numpad.impact_available': 'Disponible aujourd\'hui: {amount} €',
    'numpad.impact_remaining': 'Disponible aujourd\'hui: {current} € ➔ Reste ensuite: {diff} €',
    'numpad.impact_exceeds': 'Disponible aujourd\'hui: {current} € ➔ Dépasse le budget quotidien de {diff} €',

    // Settings
    'settings.title': 'Paramètres',
    'settings.monthly_budget': 'Budget Mensuel',
    'settings.period_days': 'Durée de la Période (Jours)',
    'settings.desired_daily': 'Budget Quotidien Souhaité',
    'settings.calculated_monthly': 'Total Mensuel Calculé:',
    'settings.save_btn': 'Enregistrer les Paramètres',
    'settings.saved_msg': '✓ Paramètres enregistrés',
    'settings.reset_heading': 'Démarrer une Nouvelle Période',
    'settings.reset_desc': 'Réinitialise le jour à 1 et efface toutes les dépenses de la période actuelle.',
    'settings.reset_btn': 'Démarrer une Nouvelle Période',
    'settings.reset_confirm_title': 'Êtes-vous sûr ?',
    'settings.reset_confirm_body': 'Toutes les dépenses de la période actuelle seront supprimées !',
    'settings.reset_confirm_btn': 'Oui, Réinitialiser',
    'settings.reset_cancel_btn': 'Annuler',
    'settings.theme_heading': 'Thème & Couleurs',
    'settings.language_heading': 'Langue / Language',
    'settings.backup_heading': 'Données & Sauvegarde',
    'settings.export_json': 'Exporter la Sauvegarde JSON',
    'settings.export_csv': 'Exporter CSV',
    'settings.import_backup': 'Importer une Sauvegarde',

    // Expenses Modal
    'expenses.title': 'Historique des Dépenses',
    'expenses.empty': 'Aucune dépense enregistrée dans cette période.',
    'expenses.page_info': 'Page {page} sur {total}',
    'expenses.delete': 'Supprimer',
    'expenses.all_expenses': 'Tout afficher',
    'expenses.recent_title': 'Dépenses Récentes',

    // Archive Modal
    'archive.title': 'Archives des Périodes',
    'archive.subtitle': 'Rapports mensuels & dépenses historiques',
    'archive.empty': 'Aucune période archivée pour le moment.',
    'archive.total_spent': 'Total Dépensé:',
    'archive.savings': 'Économies:',
    'archive.count': 'Dépenses:',
    'archive.avg_daily': 'Ø par jour:',
    'archive.view_report': 'Voir le Rapport Récapitulatif',
    'archive.back_to_list': '← Retour aux Archives',

    // Streaks & Projection
    'streak.title': 'Série d\'Économies',
    'streak.current': 'Série Actuelle',
    'streak.longest': 'Meilleure Série',
    'streak.zero_days': 'Jours Sans Dépense',
    'streak.days_unit': 'Jours',
    'projection.title': 'Projection de Fin de Mois',
    'projection.savings': 'Économies Prévisibles:',
    'projection.deficit': 'Déficit Prévisible:',
    'projection.total': 'Total Dépensé Prévisible:',
    'projection.daily_avg': 'Ø Dépense Quotidienne jusqu\'ici:',

    // Spending Trend
    'trend.title': 'Tendance des Dépenses',
    'trend.legend_ok': 'Dans le Budget',
    'trend.legend_savings': 'Zéro Dépense (€ 0)',
    'trend.legend_over': 'Hors Budget',
    'trend.avg': 'Ø {amount} € / jour',

    // Auth & SaaS
    'auth.title': 'Connexion / Inscription',
    'auth.subtitle': 'Connectez votre compte pour la synchronisation cloud sur tous vos appareils.',
    'auth.email_label': 'Adresse E-mail',
    'auth.email_placeholder': 'votre@email.fr',
    'auth.send_link': 'Envoyer un Magic Link',
    'auth.sending': 'Envoi du lien...',
    'auth.logged_in_as': 'Connecté en tant que {email}',
    'auth.logout': 'Déconnexion',
    'auth.delete_account': 'Supprimer le Compte',
    'auth.pro_badge': 'NIVEAU PRO',
    'auth.upgrade_pro': 'Passer à Pro (Stripe)',
    'auth.manage_sub': 'Gérer l\'Abonnement',

    // About
    'about.title': 'À propos de Restgeld',
    'about.tagline': 'Suivi Minimaliste de l\'Argent Quotidien',
    'about.philosophy': 'Économies conscientes sans superflu. Focus pur sur votre allocation quotidienne.',
    'about.open_source': 'Logiciel Open Source (Licence MIT)',
    'about.close': 'Fermer',

    // Footer & Misc
    'footer.tagline': 'Économisez consciemment chaque jour.',
    'common.close': 'Fermer',
    'common.back': 'Retour',
  },
}

const currentLocale = ref<SupportedLocale>('de')

function detectBrowserLocale(): SupportedLocale {
  if (typeof navigator === 'undefined') return 'de'
  const lang = (navigator.language || 'de').toLowerCase()
  if (lang.startsWith('en')) return 'en'
  if (lang.startsWith('es')) return 'es'
  if (lang.startsWith('fr')) return 'fr'
  return 'de'
}

export function useI18n() {
  function setLocale(locale: SupportedLocale) {
    currentLocale.value = locale
    try {
      localStorage.setItem(STORAGE_KEY, locale)
    } catch {
      // Ignore storage errors
    }
  }

  function initI18n() {
    try {
      const saved = localStorage.getItem(STORAGE_KEY) as SupportedLocale | null
      if (saved && ['de', 'en', 'es', 'fr'].includes(saved)) {
        currentLocale.value = saved
        return
      }
    } catch {
      // Ignore storage errors
    }
    currentLocale.value = detectBrowserLocale()
  }

  function t(key: string, params?: Record<string, string | number>): string {
    const dict = translations[currentLocale.value] || translations['de']
    let text = dict[key] || translations['de'][key] || key

    if (params) {
      Object.entries(params).forEach(([pKey, pValue]) => {
        text = text.replace(new RegExp(`\\{${pKey}\\}`, 'g'), String(pValue))
      })
    }
    return text
  }

  const localeCode = computed(() => {
    switch (currentLocale.value) {
      case 'en':
        return 'en-US'
      case 'es':
        return 'es-ES'
      case 'fr':
        return 'fr-FR'
      default:
        return 'de-DE'
    }
  })

  function formatCurrency(val: number): string {
    return val.toLocaleString(localeCode.value, {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    })
  }

  return {
    currentLocale,
    languages: SUPPORTED_LANGUAGES,
    setLocale,
    initI18n,
    t,
    formatCurrency,
  }
}
