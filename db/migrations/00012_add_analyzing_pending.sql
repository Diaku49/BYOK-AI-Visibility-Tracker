-- +goose Up
ALTER TABLE scans
ADD COLUMN analysis_queued_at TIMESTAMPTZ;

ALTER TABLE scans
DROP CONSTRAINT IF EXISTS scans_status_check;

ALTER TABLE scans
ADD CONSTRAINT scans_status_check CHECK (
    status IN ('pending', 'running', 'analyzing_pending', 'analyzing', 'completed', 'failed', 'canceled')
);

-- +goose Down
UPDATE scans
SET status = 'running',
    analysis_queued_at = NULL
WHERE status = 'analyzing_pending';

ALTER TABLE scans
DROP CONSTRAINT IF EXISTS scans_status_check;

ALTER TABLE scans
ADD CONSTRAINT scans_status_check CHECK (
    status IN ('pending', 'running', 'analyzing', 'completed', 'failed', 'canceled')
);

ALTER TABLE scans
DROP COLUMN IF EXISTS analysis_queued_at;
