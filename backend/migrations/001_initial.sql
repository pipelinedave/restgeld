CREATE TABLE IF NOT EXISTS periods (
    id TEXT PRIMARY KEY,
    start_date TIMESTAMPTZ NOT NULL,
    month_days INTEGER NOT NULL,
    base_budget NUMERIC(10,2) NOT NULL,
    monthly_total NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    period_id TEXT NOT NULL REFERENCES periods(id) ON DELETE CASCADE,
    amount NUMERIC(10,2) NOT NULL,
    note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_expenses_period ON expenses(period_id);
CREATE INDEX IF NOT EXISTS idx_expenses_created ON expenses(created_at DESC);
