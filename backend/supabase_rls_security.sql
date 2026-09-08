-- 004: RLS absichern (Supabase Security Advisory - Priority 1 critical)
--
-- Problem: 7 Tabellen in public haben RLS deaktiviert -> ueber PostgREST mit dem
--          oeffentlichen anon/publishable-Key kann JEDER alle Daten lesen/aendern
--          (users, auth_sessions, magic_links, webauthn_credentials, expenses, periods).
--
-- Loesung (Backend-heavy App, keine PostgREST-Nutzung):
--   * RLS auf allen Tabellen aktivieren.
--   * KEINE Policies fuer die Rollen `anon` / `authenticated` anlegen.
--     -> Ohne Policy verweigert RLS diesen Rollen JEDEN Zugriff (SELECT/INSERT/UPDATE/DELETE).
--   * Die privilegierte Backend-Rolle (Tabellen-Owner bzw. service_role) umgeht RLS
--     und funktioniert unveraendert weiter.
--
-- HINWEIS: Vor Ausfuehrung sicherstellen, dass KEIN App-Service ueber die
--          anon/publishable Keys auf die Datenbank zugreift (nur service_role /
--          Owner / eigene DB-Rolle). Sonst blockt dieser Fix die App.

-- 1) RLS aktivieren
ALTER TABLE public.users                ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.periods              ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.expenses             ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.magic_links          ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.auth_sessions        ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.webauthn_credentials ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.schema_migrations    ENABLE ROW LEVEL SECURITY;

-- 2) Keine Policies fuer anon/authenticated -> PostgREST-Zugriff vollstaendig blocken.
--    (Bewusst leer. Ohne Policies = ALLER Zugriff fuer anon/authenticated verweigert.)

-- 3) Optional: schema_migrations aus PostgREST-Exposition entfernen
--    (nur anon/authenticated betreffen - sie haben eh keinen Zugriff mehr)
REVOKE ALL ON public.schema_migrations FROM anon, authenticated;
