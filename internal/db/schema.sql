CREATE TYPE "job_status" AS ENUM ('queued', 'in_progress', 'completed', 'failed');

CREATE TABLE "jobs" (
    id SERIAL PRIMARY KEY,
    type TEXT NOT NULL,
    payload JSONB,
    status "job_status" NOT NULL DEFAULT 'queued',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
