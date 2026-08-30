import { ref, computed } from 'vue'

export type SupportedLocale = 'de' | 'en' | 'es' | 'fr'
export type SupportedCurrency = 'EUR' | 'USD' | 'GBP' | 'CHF' | 'JPY'

export interface LanguageOption {
  code: SupportedLocale
  name: string
  flag: string
}

export interface CurrencyOption {
  code: SupportedCurrency
  symbol: string
  name: string
}

export const SUPPORTED_LANGUAGES: LanguageOption[] = [
  { code: 'de', name: 'Deutsch', flag: '🇩🇪' },
  { code: 'en', name: 'English', flag: '🇬🇧' },
  { code: 'es', name: 'Español', flag: '🇪🇸' },
  { code: 'fr', name: 'Français', flag: '🇫🇷' },
]

export const SUPPORTED_CURRENCIES: CurrencyOption[] = [
  { code: 'EUR', symbol: '€', name: 'Euro (€)' },
  { code: 'USD', symbol: '$', name: 'US Dollar ($)' },
  { code: 'GBP', symbol: '£', name: 'Pound (£)' },
  { code: 'CHF', symbol: 'CHF', name: 'Swiss Franc' },
  { code: 'JPY', symbol: '¥', name: 'Yen (¥)' },
]

const STORAGE_LANG_KEY = 'restgeld_language'
const STORAGE_CURR_KEY = 'restgeld_currency'

export const translations: Record<SupportedLocale, Record<string, string>> = {
  de: {
    // Header & Navigation
    'header.online': 'Online',
    'header.offline': 'Offline',
    'header.sync_pending': '{count} ungesynct',
    'header.status_title': 'System- & Sync-Status',
    'header.network': 'Netzwerk-Verbindung',
    'header.connected': 'Verbunden (Online)',
    'header.api_health': 'API Service Health',
    'header.auth_health': 'Auth Service Health',
    'header.billing_health': 'Billing Service Health',
    'header.db_health': 'PostgreSQL Datenbank',
    'header.db_connected': 'Verbunden & bereit',
    'header.db_disconnected': 'Nicht erreichbar',
    'header.queue': 'Offline-Outbox Queue',
    'header.no_pending': 'Keine ausstehenden Syncs',
    'header.streak_tooltip': 'Streak & Spartage ansehen',
    'header.proj_tooltip': 'Monatsende-Prognose anzeigen',
    'header.settings_tooltip': 'Einstellungen',

    // Hero & Budget
    'budget.available_today': 'HEUTE VERFÜGBAR',
    'budget.savings': 'Ersparnis',
    'budget.deficit': 'Überzug',
    'budget.from_total': 'von {total} heute',
    'budget.base_label': 'Basis: {amount}',
    'budget.puffer_plus': '+{amount} Spar-Puffer',
    'budget.overdrawn': '-{amount} überzogen',
    'budget.empty_today': 'Kein Tagesbudget mehr übrig',
    'budget.on_track': 'Perfekt im Plan',

    // Numpad / Book Expense
    'numpad.title': 'Ausgabe buchen',
    'numpad.amount_label': 'Betrag',
    'numpad.note_label': 'Notiz (optional)',
    'numpad.note_placeholder': 'Notiz (z. B. Kaffee, Mittagessen)',
    'numpad.btn_cancel': 'Abbrechen',
    'numpad.btn_save': 'Speichern',
    'numpad.btn_saving': 'Wird gebucht...',
    'numpad.impact_available': 'Heute verfügbar: {amount}',
    'numpad.impact_remaining': 'Heute verfügbar: {current} ➔ Verbleibt danach: {diff}',
    'numpad.impact_exceeds': 'Heute verfügbar: {current} ➔ Überzieht Tagesbudget um {diff}',

    // Recent Expenses
    'recent.title': 'Letzte Ausgaben',
    'recent.show_all': 'Alle anzeigen',
    'recent.empty': 'Noch keine Ausgaben heute',
    'recent.default_note': 'Ausgabe',
    'recent.delete_title': 'Löschen',

    // Settings
    'settings.title': 'Einstellungen',
    'settings.monthly_budget': 'Monatsbudget',
    'settings.period_days': 'Periodendauer (Tage)',
    'settings.desired_daily': 'Wunsch Tages-Budget',
    'settings.calculated_monthly': 'Berechnetes Monatsbudget:',
    'settings.save_btn': 'Speichern',
    'settings.saved_msg': '✓ Einstellungen gespeichert',
    'settings.reset_heading': 'Neue Periode ab heute starten',
    'settings.reset_desc': 'Startet deinen Abrechnungszyklus ab heute bei Tag 1 mit dem konfigurierten Budget.',
    'settings.reset_btn': 'Neue Periode ab heute starten',
    'settings.reset_confirm_title': 'Sicher?',
    'settings.reset_confirm_body': 'Wirklich ab heute neu starten?',
    'settings.reset_confirm_btn': 'Ja, ab heute starten',
    'settings.reset_cancel_btn': 'Abbrechen',
    'settings.theme_heading': 'Design & Akzentfarbe',
    'settings.sound_heading': 'Audio & Haptik-Feedback',
    'settings.sound_desc': 'Subtile Klicks & Sound-Effekte bei Tastenanschlägen und Buchungen.',
    'settings.sound_toggle': 'Sound-Effekte aktivieren',
    'settings.language_heading': 'Sprache / Language',
    'settings.currency_heading': 'Währung / Currency',
    'settings.backup_heading': 'Daten & Archiv',
    'settings.backup_desc': 'Sichere deine Ausgaben oder wirf einen Blick in frühere Perioden.',
    'settings.export_json': 'JSON Backup',
    'settings.export_csv': 'CSV (Excel)',
    'settings.import_backup': 'Backup importieren (JSON/CSV)',
    'settings.importing': 'Importiere...',
    'settings.archive_trigger': '📜 Frühere Monate / Archiv ansehen',
    'settings.about_heading': 'App-Info & Philosophie',
    'settings.about_desc': 'Erfahre mehr über die Prinzipien von Restgeld und Shortcuts.',
    'settings.about_trigger': 'ℹ️ Über Restgeld öffnen',

    // Expenses Modal
    'expenses.title': 'Alle Ausgaben',
    'expenses.empty': 'Keine Ausgaben in dieser Periode vorhanden.',
    'expenses.page_info': 'Seite {page} von {total}',
    'expenses.delete': 'Löschen',
    'expenses.all_expenses': 'Alle anzeigen',
    'expenses.recent_title': 'Letzte Ausgaben',

    // Archive Modal
    'archive.title': 'Perioden-Archiv',
    'archive.subtitle': 'Monatsberichte & historische Ausgaben',
    'archive.empty': 'Noch keine vergangenen Perioden archiviert.',
    'archive.kpi_budget': 'Budget',
    'archive.kpi_spent': 'Ausgaben',
    'archive.kpi_avg': 'Ø / Tag',
    'archive.kpi_count': 'Buchungen',
    'archive.savings': 'Ersparnis:',
    'archive.total_spent': 'Gesamtausgaben:',
    'archive.view_report': 'Abschlussbericht ansehen',
    'archive.back_to_list': '← Zurück zum Archiv',
    'archive.report_details': 'Abschlussbericht & Buchungen',
    'archive.retry_btn': 'Erneut versuchen',
    'archive.loading': 'Lade Archiv...',

    // Streaks & Projection
    'streak.title': 'Spar-Streak',
    'streak.current': 'Aktuelle Spar-Streak',
    'streak.longest': 'Rekord',
    'streak.zero_days': 'Null-Euro',
    'streak.days_unit': 'Tage',
    'streak.subtitle': '🎯 {count} Spartage diese Woche',
    'streak.in_budget': '{count} Tage im Budget!',
    'projection.title': 'Monatsende-Prognose',
    'projection.savings': 'Voraussichtliche Ersparnis:',
    'projection.deficit': 'Voraussichtlicher Überzug:',
    'projection.total': 'Voraussichtliche Gesamtausgaben:',
    'projection.daily_avg': 'Ø Tagesausgabe bisher:',

    // Spending Trend
    'trend.title': 'Tages-Verlauf',
    'trend.legend_ok': 'im Budget',
    'trend.legend_savings': 'Null-Ausgaben-Tag',
    'trend.legend_over': 'über Budget',
    'trend.avg': 'Ø {amount} / Tag',

    // Auth & SaaS
    'auth.title': 'Anmelden / Registrieren',
    'auth.subtitle': 'Verbinde dein Konto für Cloud-Sync auf all deinen Geräten.',
    'auth.email_label': 'E-Mail-Adresse',
    'auth.email_placeholder': 'deine@email.de',
    'auth.send_link': 'Magic Link senden',
    'auth.sending': 'Sende Link...',
    'auth.passkey_btn': '🔐 Mit Passkey / Biometrie anmelden',
    'auth.register_passkey_btn': '✨ Dieses Gerät als Passkey registrieren',
    'auth.logged_in_as': 'Angemeldet als {email}',
    'auth.logout': 'Abmelden',
    'auth.delete_account': 'Konto löschen',
    'auth.pro_badge': 'PRO TIER',
    'auth.free_badge': 'FREE TIER',
    'auth.upgrade_pro': 'Auf Pro upgraden (Stripe)',
    // Categories (Smart Zero-Bloat)
    'category.coffee': 'Kaffee & Bäckerei',
    'category.food': 'Essen & Gastro',
    'category.groceries': 'Supermarkt & Einkauf',
    'category.transport': 'Mobilität & Transport',
    'category.leisure': 'Freizeit & Ausgehen',
    'category.shopping': 'Shopping & Technik',
    'category.health': 'Gesundheit & Sport',
    'category.other': 'Sonstiges',

    // Monitoring & Observability
    'monitoring.title': 'System Observability & Monitoring',
    'monitoring.cluster_status': 'Cluster Status',
    'monitoring.healthy': 'Alle Systeme operational',
    'monitoring.degraded': 'Teilweise beeinträchtigt',
    'monitoring.critical': 'Kritischer Systemausfall',
    'monitoring.live_telemetry': 'Live Telemetrie',
    'monitoring.goroutines': 'Goroutines',
    'monitoring.memory': 'RAM Allokation',
    'monitoring.uptime': 'Betriebszeit',
    'monitoring.services_heading': 'Microservices & Endpunkte',
    'monitoring.refresh': 'Neu prüfen',
    'monitoring.auto_refresh': 'Auto-Refresh (5s)',
    'monitoring.metrics_link': 'Prometheus Metriken',

    // About
    'about.title': 'Über Restgeld',
    'about.tagline': 'Minimalistischer Daily Allowance Tracker',
    'about.philosophy': 'Achtsames Sparen ohne Schnickschnack. Fokus auf das wesentliche Tagesbudget.',
    'about.open_source': 'Open Source Software (MIT Lizenz)',
    'about.close': 'Schließen',

    // Footer & Misc
    'footer.tagline': 'Track daily. Stay in budget.',
    'footer.commit': 'Commit',
    'common.close': 'Schließen',
    'common.back': 'Zurück',
    'toast.login_success': '✓ Erfolgreich eingeloggt',
    'toast.passkey_success': '✓ Passkey erfolgreich registriert',
    'toast.expense_deleted': '✓ Ausgabe gelöscht',
    'toast.offline_queued': '✓ Offline gespeichert (wird synchronisiert)',
  },
  en: {
    // Header & Navigation
    'header.online': 'Online',
    'header.offline': 'Offline',
    'header.sync_pending': '{count} pending',
    'header.status_title': 'System & Sync Status',
    'header.network': 'Network Connection',
    'header.connected': 'Connected (Online)',
    'header.api_health': 'API Service Health',
    'header.auth_health': 'Auth Service Health',
    'header.billing_health': 'Billing Service Health',
    'header.db_health': 'PostgreSQL Database',
    'header.db_connected': 'Connected & Ready',
    'header.db_disconnected': 'Unreachable',
    'header.queue': 'Offline Outbox Queue',
    'header.no_pending': 'No pending syncs',
    'header.streak_tooltip': 'View Streak & Savings',
    'header.proj_tooltip': 'View Month-End Projection',
    'header.settings_tooltip': 'Settings',

    // Hero & Budget
    'budget.available_today': 'AVAILABLE TODAY',
    'budget.savings': 'Savings',
    'budget.deficit': 'Deficit',
    'budget.from_total': 'from {total} today',
    'budget.base_label': 'Base: {amount}',
    'budget.puffer_plus': '+{amount} Savings Buffer',
    'budget.overdrawn': '{amount} overdrawn',
    'budget.empty_today': 'No daily budget remaining',
    'budget.on_track': 'Perfect on Track',

    // Numpad / Book Expense
    'numpad.title': 'Add Expense',
    'numpad.amount_label': 'Amount',
    'numpad.note_label': 'Note (optional)',
    'numpad.note_placeholder': 'Note (e.g. Coffee, Lunch)',
    'numpad.btn_cancel': 'Cancel',
    'numpad.btn_save': 'Save',
    'numpad.btn_saving': 'Saving...',
    'numpad.impact_available': 'Available today: {amount}',
    'numpad.impact_remaining': 'Available today: {current} ➔ Remaining after: {diff}',
    'numpad.impact_exceeds': 'Available today: {current} ➔ Exceeds daily budget by {diff}',

    // Recent Expenses
    'recent.title': 'Recent Expenses',
    'recent.show_all': 'View All',
    'recent.empty': 'No expenses recorded today',
    'recent.default_note': 'Expense',
    'recent.delete_title': 'Delete',

    // Settings
    'settings.title': 'Settings',
    'settings.monthly_budget': 'Monthly Budget',
    'settings.period_days': 'Period Duration (Days)',
    'settings.desired_daily': 'Desired Daily Budget',
    'settings.calculated_monthly': 'Calculated Monthly Total:',
    'settings.save_btn': 'Save Settings',
    'settings.saved_msg': '✓ Settings saved',
    'settings.reset_heading': 'Start New Period',
    'settings.reset_desc': 'Resets day to 1 and archives all current expenses.',
    'settings.reset_btn': 'Start New Period',
    'settings.reset_confirm_title': 'Are you sure?',
    'settings.reset_confirm_body': 'All expenses of the current period will be closed!',
    'settings.reset_confirm_btn': 'Yes, Start New Period',
    'settings.reset_cancel_btn': 'Cancel',
    'settings.theme_heading': 'Design & Color Theme',
    'settings.sound_heading': 'Audio & Haptic Feedback',
    'settings.sound_desc': 'Subtle clicks & sound effects when typing and saving expenses.',
    'settings.sound_toggle': 'Enable sound effects',
    'settings.language_heading': 'Language / Sprache',
    'settings.currency_heading': 'Currency / Währung',
    'settings.backup_heading': 'Data & Archive',
    'settings.backup_desc': 'Backup your expenses or explore past monthly reports.',
    'settings.export_json': 'JSON Backup',
    'settings.export_csv': 'CSV (Excel)',
    'settings.import_backup': 'Import Backup (JSON/CSV)',
    'settings.importing': 'Importing...',
    'settings.archive_trigger': '📜 View Previous Months / Archive',
    'settings.about_heading': 'App Info & Philosophy',
    'settings.about_desc': 'Learn more about Restgeld principles and shortcuts.',
    'settings.about_trigger': 'ℹ️ Open About Restgeld',

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
    'archive.kpi_budget': 'Budget',
    'archive.kpi_spent': 'Spent',
    'archive.kpi_avg': 'Ø / Day',
    'archive.kpi_count': 'Items',
    'archive.savings': 'Savings:',
    'archive.total_spent': 'Total Spent:',
    'archive.view_report': 'View Summary Report',
    'archive.back_to_list': '← Back to Archive',
    'archive.report_details': 'Detailed Report',
    'archive.retry_btn': 'Try Again',
    'archive.loading': 'Loading archive...',

    // Streaks & Projection
    'streak.title': 'Saving Streak',
    'streak.current': 'Current Streak',
    'streak.longest': 'Best Streak',
    'streak.zero_days': 'Zero-Spend Days',
    'streak.days_unit': 'Days',
    'streak.subtitle': '🎯 {count} savings days this week',
    'streak.in_budget': '{count} days within budget!',
    'projection.title': 'Month-End Projection',
    'projection.savings': 'Expected Savings:',
    'projection.deficit': 'Expected Deficit:',
    'projection.total': 'Expected Total Spent:',
    'projection.daily_avg': 'Ø Daily Spend so far:',

    // Spending Trend
    'trend.title': 'Spending Trend',
    'trend.legend_ok': 'Within Budget',
    'trend.legend_savings': 'Zero Spend',
    'trend.legend_over': 'Over Budget',
    'trend.avg': 'Ø {amount} / day',

    // Auth & SaaS
    'auth.title': 'Sign In / Register',
    'auth.subtitle': 'Connect your account for cloud sync across all devices.',
    'auth.email_label': 'Email Address',
    'auth.email_placeholder': 'your@email.com',
    'auth.send_link': 'Send Magic Link',
    'auth.sending': 'Sending link...',
    'auth.passkey_btn': '🔐 Sign In with Passkey / Biometrics',
    'auth.register_passkey_btn': '✨ Register This Device as Passkey',
    'auth.logged_in_as': 'Signed in as {email}',
    'auth.logout': 'Sign Out',
    'auth.delete_account': 'Delete Account',
    'auth.pro_badge': 'PRO TIER',
    'auth.free_badge': 'FREE TIER',
    'auth.upgrade_pro': 'Upgrade to Pro (Stripe)',
    'auth.manage_sub': 'Manage Subscription',
    // Categories (Smart Zero-Bloat)
    'category.coffee': 'Coffee & Bakery',
    'category.food': 'Food & Dining',
    'category.groceries': 'Groceries & Markets',
    'category.transport': 'Transport & Travel',
    'category.leisure': 'Leisure & Nightlife',
    'category.shopping': 'Shopping & Tech',
    'category.health': 'Health & Fitness',
    'category.other': 'Other',

    // Monitoring & Observability
    'monitoring.title': 'System Observability & Monitoring',
    'monitoring.cluster_status': 'Cluster Status',
    'monitoring.healthy': 'All systems operational',
    'monitoring.degraded': 'Partially degraded',
    'monitoring.critical': 'Critical service outage',
    'monitoring.live_telemetry': 'Live Telemetry',
    'monitoring.goroutines': 'Goroutines',
    'monitoring.memory': 'RAM Allocated',
    'monitoring.uptime': 'Uptime',
    'monitoring.services_heading': 'Microservices & Endpoints',
    'monitoring.refresh': 'Refresh',
    'monitoring.auto_refresh': 'Auto-Refresh (5s)',
    'monitoring.metrics_link': 'Prometheus Metrics',

    // About
    'about.title': 'About Restgeld',
    'about.tagline': 'Minimalist Daily Allowance Tracker',
    'about.philosophy': 'Mindful saving without unnecessary bloat. Pure focus on your daily allowance.',
    'about.open_source': 'Open Source Software (MIT License)',
    'about.close': 'Close',

    // Footer & Misc
    'footer.tagline': 'Save mindfully every day.',
    'footer.commit': 'Commit',
    'common.close': 'Close',
    'common.back': 'Back',
    'toast.login_success': '✓ Successfully logged in',
    'toast.passkey_success': '✓ Passkey successfully registered',
    'toast.expense_deleted': '✓ Expense deleted',
    'toast.offline_queued': '✓ Saved offline (will sync automatically)',
  },
  es: {
    // Header & Navigation
    'header.online': 'En línea',
    'header.offline': 'Sin conexión',
    'header.sync_pending': '{count} pendiente',
    'header.status_title': 'Estado del Sistema y Sync',
    'header.network': 'Conexión de Red',
    'header.connected': 'Conectado (En línea)',
    'header.api_health': 'Salud del Servicio API',
    'header.auth_health': 'Salud del Servicio Auth',
    'header.billing_health': 'Salud del Servicio Billing',
    'header.db_health': 'Base de Datos PostgreSQL',
    'header.db_connected': 'Conectado y Listo',
    'header.db_disconnected': 'Inaccesible',
    'header.queue': 'Cola de Salida Offline',
    'header.no_pending': 'Sin sincronizaciones pendientes',
    'header.streak_tooltip': 'Ver Racha y Ahorro',
    'header.proj_tooltip': 'Ver Proyección de Fin de Mes',
    'header.settings_tooltip': 'Ajustes',

    // Hero & Budget
    'budget.available_today': 'DISPONIBLE HOY',
    'budget.savings': 'Ahorro',
    'budget.deficit': 'Déficit',
    'budget.from_total': 'de {total} hoy',
    'budget.base_label': 'Base: {amount}',
    'budget.puffer_plus': '+{amount} Margen de Ahorro',
    'budget.overdrawn': '{amount} excedido',
    'budget.empty_today': 'Sin presupuesto restante para hoy',
    'budget.on_track': 'Perfecto en Plan',

    // Numpad / Book Expense
    'numpad.title': 'Añadir Gasto',
    'numpad.amount_label': 'Monto',
    'numpad.note_label': 'Nota (opcional)',
    'numpad.note_placeholder': 'Nota (ej. Café, Almuerzo)',
    'numpad.btn_cancel': 'Cancelar',
    'numpad.btn_save': 'Guardar',
    'numpad.btn_saving': 'Guardando...',
    'numpad.impact_available': 'Disponible hoy: {amount}',
    'numpad.impact_remaining': 'Disponible hoy: {current} ➔ Quedará después: {diff}',
    'numpad.impact_exceeds': 'Disponible hoy: {current} ➔ Excede el presupuesto diario por {diff}',

    // Recent Expenses
    'recent.title': 'Gastos Recientes',
    'recent.show_all': 'Ver todos',
    'recent.empty': 'No hay gastos registrados hoy',
    'recent.default_note': 'Gasto',
    'recent.delete_title': 'Eliminar',

    // Settings
    'settings.title': 'Ajustes',
    'settings.monthly_budget': 'Presupuesto Mensual',
    'settings.period_days': 'Duración del Período (Días)',
    'settings.desired_daily': 'Presupuesto Diario Deseado',
    'settings.calculated_monthly': 'Total Mensual Calculado:',
    'settings.save_btn': 'Guardar Ajustes',
    'settings.saved_msg': '✓ Ajustes guardados',
    'settings.reset_heading': 'Iniciar Nuevo Período',
    'settings.reset_desc': 'Reinicia el día a 1 y archiva todos los gastos actuales.',
    'settings.reset_btn': 'Iniciar Nuevo Período',
    'settings.reset_confirm_title': '¿Estás seguro?',
    'settings.reset_confirm_body': '¡Todos los gastos del período actual se cerrarán!',
    'settings.reset_confirm_btn': 'Sí, Iniciar Período',
    'settings.reset_cancel_btn': 'Cancelar',
    'settings.theme_heading': 'Diseño y Color',
    'settings.sound_heading': 'Audio y retroalimentación háptica',
    'settings.sound_desc': 'Clics y efectos de sonido sutiles al pulsar teclas y guardar.',
    'settings.sound_toggle': 'Activar efectos de sonido',
    'settings.language_heading': 'Idioma / Language',
    'settings.currency_heading': 'Moneda / Currency',
    'settings.backup_heading': 'Datos y Archivo',
    'settings.backup_desc': 'Haz una copia de seguridad de tus gastos o explora informes anteriores.',
    'settings.export_json': 'Copia JSON',
    'settings.export_csv': 'CSV (Excel)',
    'settings.import_backup': 'Importar Copia (JSON/CSV)',
    'settings.importing': 'Importando...',
    'settings.archive_trigger': '📜 Ver Meses Anteriores / Archivo',
    'settings.about_heading': 'Información y Filosofía',
    'settings.about_desc': 'Aprende más sobre los principios y atajos de Restgeld.',
    'settings.about_trigger': 'ℹ️ Abrir Acerca de Restgeld',

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
    'archive.kpi_budget': 'Presupuesto',
    'archive.kpi_spent': 'Gastado',
    'archive.kpi_avg': 'Ø / Día',
    'archive.kpi_count': 'Elementos',
    'archive.savings': 'Ahorro:',
    'archive.total_spent': 'Gasto Total:',
    'archive.view_report': 'Ver Informe Resumido',
    'archive.back_to_list': '← Volver al Archivo',
    'archive.report_details': 'Informe Detallado',
    'archive.retry_btn': 'Reintentar',
    'archive.loading': 'Cargando archivo...',

    // Streaks & Projection
    'streak.title': 'Racha de Ahorro',
    'streak.current': 'Racha Actual',
    'streak.longest': 'Mejor Racha',
    'streak.zero_days': 'Días Sin Gasto',
    'streak.days_unit': 'Días',
    'streak.subtitle': '🎯 {count} días de ahorro esta semana',
    'streak.in_budget': '¡{count} días en presupuesto!',
    'projection.title': 'Proyección de Fin de Mes',
    'projection.savings': 'Ahorro Previsto:',
    'projection.deficit': 'Déficit Previsto:',
    'projection.total': 'Gasto Total Previsto:',
    'projection.daily_avg': 'Ø Gasto Diario hasta ahora:',

    // Spending Trend
    'trend.title': 'Tendencia de Gastos',
    'trend.legend_ok': 'En Presupuesto',
    'trend.legend_savings': 'Cero Gastos',
    'trend.legend_over': 'Sobre Presupuesto',
    'trend.avg': 'Ø {amount} / día',

    // Auth & SaaS
    'auth.title': 'Iniciar Sesión / Registro',
    'auth.subtitle': 'Conecta tu cuenta para sincronizar en la nube en todos tus dispositivos.',
    'auth.email_label': 'Correo Electrónico',
    'auth.email_placeholder': 'tu@email.com',
    'auth.send_link': 'Enviar Magic Link',
    'auth.sending': 'Enviando enlace...',
    'auth.passkey_btn': '🔐 Iniciar con Passkey / Biometría',
    'auth.register_passkey_btn': '✨ Registrar Este Dispositivo como Passkey',
    'auth.logged_in_as': 'Sesión iniciada como {email}',
    'auth.logout': 'Cerrar Sesión',
    'auth.delete_account': 'Eliminar Cuenta',
    'auth.pro_badge': 'NIVEL PRO',
    'auth.free_badge': 'NIVEL GRATUITO',
    'auth.upgrade_pro': 'Mejorar a Pro (Stripe)',
    'auth.manage_sub': 'Gestionar Suscripción',
    // Categories (Smart Zero-Bloat)
    'category.coffee': 'Café y Panadería',
    'category.food': 'Comida y Restaurantes',
    'category.groceries': 'Supermercado y Compras',
    'category.transport': 'Transporte y Movilidad',
    'category.leisure': 'Ocio y Salidas',
    'category.shopping': 'Compras y Tecnología',
    'category.health': 'Salud y Deporte',
    'category.other': 'Otros',

    // Monitoring & Observability
    'monitoring.title': 'Observabilidad y Monitoreo del Sistema',
    'monitoring.cluster_status': 'Estado del Clúster',
    'monitoring.healthy': 'Todos los sistemas operativos',
    'monitoring.degraded': 'Parcialmente degradado',
    'monitoring.critical': 'Corte crítico de servicio',
    'monitoring.live_telemetry': 'Telemetría en Vivo',
    'monitoring.goroutines': 'Goroutines',
    'monitoring.memory': 'RAM Asignada',
    'monitoring.uptime': 'Tiempo de Actividad',
    'monitoring.services_heading': 'Microservicios y Endpoints',
    'monitoring.refresh': 'Actualizar',
    'monitoring.auto_refresh': 'Auto-Actualizar (5s)',
    'monitoring.metrics_link': 'Métricas Prometheus',

    // About
    'about.title': 'Acerca de Restgeld',
    'about.tagline': 'Rastreador Minimalista de Dinero Diario',
    'about.philosophy': 'Ahorro consciente sin complicaciones. Enfocado puramente en tu asignación diaria.',
    'about.open_source': 'Software de Código Abierto (Licencia MIT)',
    'about.close': 'Cerrar',

    // Footer & Misc
    'footer.tagline': 'Ahorra de forma consciente cada día.',
    'footer.commit': 'Commit',
    'common.close': 'Cerrar',
    'common.back': 'Volver',
    'toast.login_success': '✓ Sesión iniciada con éxito',
    'toast.passkey_success': '✓ Passkey registrado con éxito',
    'toast.expense_deleted': '✓ Gasto eliminado',
    'toast.offline_queued': '✓ Guardado offline (se sincronizará)',
  },
  fr: {
    // Header & Navigation
    'header.online': 'En ligne',
    'header.offline': 'Hors ligne',
    'header.sync_pending': '{count} en attente',
    'header.status_title': 'État du Système et Synchro',
    'header.network': 'Connexion Réseau',
    'header.connected': 'Connecté (En ligne)',
    'header.api_health': 'Santé du Service API',
    'header.auth_health': 'Santé du Service Auth',
    'header.billing_health': 'Santé du Service Billing',
    'header.db_health': 'Base de Données PostgreSQL',
    'header.db_connected': 'Connecté et Prêt',
    'header.db_disconnected': 'Inaccessible',
    'header.queue': 'File d\'attente Hors-ligne',
    'header.no_pending': 'Aucune synchronisation en attente',
    'header.streak_tooltip': 'Voir la Série & Économies',
    'header.proj_tooltip': 'Voir la Projection Fin de Mois',
    'header.settings_tooltip': 'Paramètres',

    // Hero & Budget
    'budget.available_today': 'DISPONIBLE AUJOURD\'HUI',
    'budget.savings': 'Économies',
    'budget.deficit': 'Déficit',
    'budget.from_total': 'sur {total} aujourd\'hui',
    'budget.base_label': 'Base: {amount}',
    'budget.puffer_plus': '+{amount} Marge d\'Épargne',
    'budget.overdrawn': '{amount} dépassé',
    'budget.empty_today': 'Aucun budget quotidien restant',
    'budget.on_track': 'Parfaitement dans les Clous',

    // Numpad / Book Expense
    'numpad.title': 'Ajouter une Dépense',
    'numpad.amount_label': 'Montant',
    'numpad.note_label': 'Note (optionnel)',
    'numpad.note_placeholder': 'Note (ex. Café, Déjeuner)',
    'numpad.btn_cancel': 'Annuler',
    'numpad.btn_save': 'Enregistrer',
    'numpad.btn_saving': 'Enregistrement...',
    'numpad.impact_available': 'Disponible aujourd\'hui: {amount}',
    'numpad.impact_remaining': 'Disponible aujourd\'hui: {current} ➔ Reste ensuite: {diff}',
    'numpad.impact_exceeds': 'Disponible aujourd\'hui: {current} ➔ Dépasse le budget quotidien de {diff}',

    // Recent Expenses
    'recent.title': 'Dépenses Récentes',
    'recent.show_all': 'Tout afficher',
    'recent.empty': 'Aucune dépense enregistrée aujourd\'hui',
    'recent.default_note': 'Dépense',
    'recent.delete_title': 'Supprimer',

    // Settings
    'settings.title': 'Paramètres',
    'settings.monthly_budget': 'Budget Mensuel',
    'settings.period_days': 'Durée de la Période (Jours)',
    'settings.desired_daily': 'Budget Quotidien Souhaité',
    'settings.calculated_monthly': 'Total Mensuel Calculé:',
    'settings.save_btn': 'Enregistrer les Paramètres',
    'settings.saved_msg': '✓ Paramètres enregistrés',
    'settings.reset_heading': 'Démarrer une Nouvelle Période',
    'settings.reset_desc': 'Réinitialise le jour à 1 et archive toutes les dépenses actuelles.',
    'settings.reset_btn': 'Démarrer une Nouvelle Période',
    'settings.reset_confirm_title': 'Êtes-vous sûr ?',
    'settings.reset_confirm_body': 'Toutes les dépenses de la période actuelle seront clôturées !',
    'settings.reset_confirm_btn': 'Oui, Nouvelle Période',
    'settings.reset_cancel_btn': 'Annuler',
    'settings.theme_heading': 'Thème & Couleurs',
    'settings.sound_heading': 'Audio et retour haptique',
    'settings.sound_desc': 'Clics et effets sonores subtils lors de la saisie et enregistrement.',
    'settings.sound_toggle': 'Activer les effets sonores',
    'settings.language_heading': 'Langue / Language',
    'settings.currency_heading': 'Devise / Currency',
    'settings.backup_heading': 'Données & Archives',
    'settings.backup_desc': 'Sauvegardez vos dépenses ou consultez les rapports mensuels passés.',
    'settings.export_json': 'Sauvegarde JSON',
    'settings.export_csv': 'CSV (Excel)',
    'settings.import_backup': 'Importer une Sauvegarde (JSON/CSV)',
    'settings.importing': 'Importation...',
    'settings.archive_trigger': '📜 Voir les Mois Précédents / Archives',
    'settings.about_heading': 'Info App & Philosophie',
    'settings.about_desc': 'En savoir plus sur les principes et raccourcis de Restgeld.',
    'settings.about_trigger': 'ℹ️ Ouvrir À Propos',

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
    'archive.kpi_budget': 'Budget',
    'archive.kpi_spent': 'Dépensé',
    'archive.kpi_avg': 'Ø / Jour',
    'archive.kpi_count': 'Éléments',
    'archive.savings': 'Économies:',
    'archive.total_spent': 'Total Dépensé:',
    'archive.view_report': 'Voir le Rapport Récapitulatif',
    'archive.back_to_list': '← Retour aux Archives',
    'archive.report_details': 'Rapport Détaillé',
    'archive.retry_btn': 'Réessayer',
    'archive.loading': 'Chargement des archives...',

    // Streaks & Projection
    'streak.title': 'Série d\'Économies',
    'streak.current': 'Série Actuelle',
    'streak.longest': 'Meilleure Série',
    'streak.zero_days': 'Jours Sans Dépense',
    'streak.days_unit': 'Jours',
    'streak.subtitle': '🎯 {count} jours d\'économie cette semaine',
    'streak.in_budget': '{count} jours dans le budget !',
    'projection.title': 'Projection de Fin de Mois',
    'projection.savings': 'Économies Prévisibles:',
    'projection.deficit': 'Déficit Prévisible:',
    'projection.total': 'Total Dépensé Prévisible:',
    'projection.daily_avg': 'Ø Dépense Quotidienne jusqu\'ici:',

    // Spending Trend
    'trend.title': 'Tendance des Dépenses',
    'trend.legend_ok': 'Dans le Budget',
    'trend.legend_savings': 'Zéro Dépense',
    'trend.legend_over': 'Hors Budget',
    'trend.avg': 'Ø {amount} / jour',

    // Auth & SaaS
    'auth.title': 'Connexion / Inscription',
    'auth.subtitle': 'Connectez votre compte pour la synchronisation cloud sur tous vos appareils.',
    'auth.email_label': 'Adresse E-mail',
    'auth.email_placeholder': 'votre@email.fr',
    'auth.send_link': 'Envoyer un Magic Link',
    'auth.sending': 'Envoi du lien...',
    'auth.passkey_btn': '🔐 Connexion par Passkey / Biométrie',
    'auth.register_passkey_btn': '✨ Enregistrer cet Appareil comme Passkey',
    'auth.logged_in_as': 'Connecté en tant que {email}',
    'auth.logout': 'Déconnexion',
    'auth.delete_account': 'Supprimer le Compte',
    'auth.pro_badge': 'NIVEAU PRO',
    'auth.free_badge': 'NIVEAU GRATUIT',
    'auth.upgrade_pro': 'Passer à Pro (Stripe)',
    'auth.manage_sub': 'Gérer l\'Abonnement',
    // Categories (Smart Zero-Bloat)
    'category.coffee': 'Café & Boulangerie',
    'category.food': 'Nourriture & Restaurants',
    'category.groceries': 'Courses & Supermarché',
    'category.transport': 'Transports & Mobilité',
    'category.leisure': 'Loisirs & Sorties',
    'category.shopping': 'Shopping & High-Tech',
    'category.health': 'Santé & Bien-être',
    'category.other': 'Divers',

    // Monitoring & Observability
    'monitoring.title': 'Observabilité & Surveillance du Système',
    'monitoring.cluster_status': 'État du Cluster',
    'monitoring.healthy': 'Tous les systèmes opérationnels',
    'monitoring.degraded': 'Partiellement dégradé',
    'monitoring.critical': 'Panne critique du service',
    'monitoring.live_telemetry': 'Télémétrie en Direct',
    'monitoring.goroutines': 'Goroutines',
    'monitoring.memory': 'RAM Allouée',
    'monitoring.uptime': 'Temps de Fonctionnement',
    'monitoring.services_heading': 'Microservices & Points de Terminaison',
    'monitoring.refresh': 'Actualiser',
    'monitoring.auto_refresh': 'Actualisation Auto (5s)',
    'monitoring.metrics_link': 'Métriques Prometheus',

    // About
    'about.title': 'À propos de Restgeld',
    'about.tagline': 'Suivi Minimaliste de l\'Argent Quotidien',
    'about.philosophy': 'Économies conscientes sans superflu. Focus pur sur votre allocation quotidienne.',
    'about.open_source': 'Logiciel Open Source (Licence MIT)',
    'about.close': 'Fermer',

    // Footer & Misc
    'footer.tagline': 'Économisez consciemment chaque jour.',
    'footer.commit': 'Commit',
    'common.close': 'Fermer',
    'common.back': 'Retour',
    'toast.login_success': '✓ Connexion réussie',
    'toast.passkey_success': '✓ Passkey enregistré avec succès',
    'toast.expense_deleted': '✓ Dépense supprimée',
    'toast.offline_queued': '✓ Sauvegardé hors ligne (sera synchronisé)',
  },
}

const currentLocale = ref<SupportedLocale>('de')
const currentCurrency = ref<SupportedCurrency>('EUR')

function detectBrowserLocale(): SupportedLocale {
  if (typeof navigator === 'undefined') return 'de'
  const lang = (navigator.language || 'de').toLowerCase()
  if (lang.startsWith('en')) return 'en'
  if (lang.startsWith('es')) return 'es'
  if (lang.startsWith('fr')) return 'fr'
  return 'de'
}

export interface CategoryStat {
  key: string
  icon: string
  name: string
  total: number
  percentage: number
  count: number
}

export function detectCategoryKey(note?: string): { key: string; icon: string } {
  if (!note || note.trim().length === 0) return { key: 'other', icon: '💶' }
  const clean = note.trim().toLowerCase()

  if (/kino|cinema|movie|party|club|konzert|concert|festival|event|feiern|bier|beer|bar|drinks|cocktail/.test(clean)) return { key: 'leisure', icon: '🎉' }
  if (/kaffee|coffee|espresso|cappuccino|latte|bäcker|bakery|croissant|donut|tee|tea/.test(clean)) return { key: 'coffee', icon: '☕' }
  if (/essen|food|lunch|dinner|mittag|abendessen|pizza|döner|burger|sushi|pasta|restaurant|imbiss|snack|kebab/.test(clean)) return { key: 'food', icon: '🍔' }
  if (/supermarkt|einkauf|groceries|rewe|aldi|lidl|edeka|dm|rossmann|kaufland|penny|market/.test(clean)) return { key: 'groceries', icon: '🛒' }
  if (/tanken|benzin|fuel|diesel|bahn|zug|train|bus|ticket|uber|taxi|bolt|fahrt|mvv|bvg|flight|flug/.test(clean)) return { key: 'transport', icon: '🚗' }
  if (/amazon|zalando|kleidung|clothes|tech|apple|electronics|gadget|paket|shopping/.test(clean)) return { key: 'shopping', icon: '📦' }
  if (/apotheke|pharma|arzt|doctor|gym|fitness|sport|climbing|training|medikament/.test(clean)) return { key: 'health', icon: '💊' }

  return { key: 'other', icon: '💶' }
}

/**
 * Automatically detect an emoji category tag from free-text notes
 * without any user configuration, dropdowns, or friction!
 */
export function detectCategoryIcon(note?: string): string {
  return detectCategoryKey(note).icon
}

export function calculateCategoryBreakdown(
  expenses: { amount: number; note?: string }[],
  t: (key: string) => string
): CategoryStat[] {
  if (!expenses || expenses.length === 0) return []

  const totalSpent = expenses.reduce((acc, curr) => acc + (curr.amount || 0), 0)
  const map: Record<string, { icon: string; total: number; count: number }> = {}

  for (const exp of expenses) {
    const { key, icon } = detectCategoryKey(exp.note)
    if (!map[key]) {
      map[key] = { icon, total: 0, count: 0 }
    }
    map[key].total += exp.amount || 0
    map[key].count += 1
  }

  const result: CategoryStat[] = Object.entries(map).map(([key, data]) => {
    const percentage = totalSpent > 0 ? Math.round((data.total / totalSpent) * 100) : 0
    return {
      key,
      icon: data.icon,
      name: t(`category.${key}`),
      total: data.total,
      percentage,
      count: data.count,
    }
  })

  return result.sort((a, b) => b.total - a.total)
}

export function useI18n() {
  function setLocale(locale: SupportedLocale) {
    currentLocale.value = locale
    try {
      localStorage.setItem(STORAGE_LANG_KEY, locale)
    } catch {
      // Ignore storage errors
    }
  }

  function setCurrency(currency: SupportedCurrency) {
    currentCurrency.value = currency
    try {
      localStorage.setItem(STORAGE_CURR_KEY, currency)
    } catch {
      // Ignore storage errors
    }
  }

  function initI18n() {
    try {
      const savedLang = localStorage.getItem(STORAGE_LANG_KEY) as SupportedLocale | null
      if (savedLang && ['de', 'en', 'es', 'fr'].includes(savedLang)) {
        currentLocale.value = savedLang
      } else {
        currentLocale.value = detectBrowserLocale()
      }

      const savedCurr = localStorage.getItem(STORAGE_CURR_KEY) as SupportedCurrency | null
      if (savedCurr && ['EUR', 'USD', 'GBP', 'CHF', 'JPY'].includes(savedCurr)) {
        currentCurrency.value = savedCurr
      }
    } catch {
      // Ignore storage errors
    }
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

  const currencySymbol = computed(() => {
    const opt = SUPPORTED_CURRENCIES.find((c) => c.code === currentCurrency.value)
    return opt?.symbol || '€'
  })

  function formatCurrency(val: number): string {
    const digits = currentCurrency.value === 'JPY' ? 0 : 2
    return val.toLocaleString(localeCode.value, {
      minimumFractionDigits: digits,
      maximumFractionDigits: digits,
    })
  }

  function formatMoney(val: number): string {
    const formatted = formatCurrency(val)
    if (currentCurrency.value === 'USD') return `$${formatted}`
    if (currentCurrency.value === 'GBP') return `£${formatted}`
    if (currentCurrency.value === 'JPY') return `¥${formatted}`
    if (currentCurrency.value === 'CHF') return `${formatted} CHF`
    return `${formatted} €`
  }

  return {
    currentLocale,
    currentCurrency,
    currencySymbol,
    languages: SUPPORTED_LANGUAGES,
    currencies: SUPPORTED_CURRENCIES,
    setLocale,
    setCurrency,
    initI18n,
    t,
    formatCurrency,
    formatMoney,
    detectCategoryIcon,
    detectCategoryKey,
    calculateCategoryBreakdown: (exps: { amount: number; note?: string }[]) => calculateCategoryBreakdown(exps, t),
  }
}
