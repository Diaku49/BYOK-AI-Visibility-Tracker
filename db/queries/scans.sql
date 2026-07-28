-- name: CreateScan :one
INSERT INTO scans (
    id,
    project_id,
    status,
    tries_per_prompt,
    total_runs,
    completed_runs,
    failed_runs,
    summary,
    error,
    started_at,
    finished_at
)
VALUES (
    sqlc.arg(id),
    sqlc.arg(project_id),
    sqlc.arg(status),
    sqlc.arg(tries_per_prompt),
    sqlc.arg(total_runs),
    sqlc.arg(completed_runs),
    sqlc.arg(failed_runs),
    sqlc.narg(summary),
    sqlc.narg(error),
    sqlc.narg(started_at),
    sqlc.narg(finished_at)
)
RETURNING *;

-- name: GetScanByID :one
SELECT *
FROM scans
WHERE id = $1;

-- name: UpdateScan :one
UPDATE scans
SET
    project_id = sqlc.arg(project_id),
    status = sqlc.arg(status),
    tries_per_prompt = sqlc.arg(tries_per_prompt),
    total_runs = sqlc.arg(total_runs),
    completed_runs = sqlc.arg(completed_runs),
    failed_runs = sqlc.arg(failed_runs),
    summary = sqlc.narg(summary),
    error = sqlc.narg(error),
    started_at = sqlc.narg(started_at),
    finished_at = sqlc.narg(finished_at)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: IncrementScanCompletedRuns :exec
UPDATE scans
SET completed_runs = completed_runs + 1
WHERE id = sqlc.arg(id);

-- name: IncrementScanFailedRuns :exec
UPDATE scans
SET failed_runs = failed_runs + 1
WHERE id = sqlc.arg(id);
