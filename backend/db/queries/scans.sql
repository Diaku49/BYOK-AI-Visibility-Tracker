-- name: CreateScan :one
INSERT INTO scans (
    id,
    project_id,
    status,
    total_runs,
    completed_runs,
    failed_runs
)
SELECT
    sqlc.arg(id),
    p.id,
    sqlc.arg(status),
    sqlc.arg(total_runs),
    sqlc.arg(completed_runs),
    sqlc.arg(failed_runs)
FROM projects p
WHERE p.id = sqlc.arg(project_id)
  AND p.user_id = sqlc.arg(user_id)
RETURNING id;

-- name: GetScanByID :one
SELECT *
FROM scans
WHERE id = $1;

-- name: UpdateScanStateByID :one
UPDATE scans
SET status = sqlc.arg(status),
    error = sqlc.narg(error),
    finished_at = COALESCE(now(), finished_at)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateScan :one
UPDATE scans
SET
    project_id = sqlc.arg(project_id),
    status = sqlc.arg(status),
    total_runs = sqlc.arg(total_runs),
    completed_runs = sqlc.arg(completed_runs),
    failed_runs = sqlc.arg(failed_runs),
    summary = sqlc.narg(summary),
    error = sqlc.narg(error),
    started_at = sqlc.narg(started_at),
    finished_at = sqlc.narg(finished_at)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: GetScansForAnalysis :many
WITH ready_scans AS (
    SELECT s.id, s.project_id
    FROM scans s
    WHERE (
        s.status = 'running'
        OR (s.status = 'analyzing' AND s.started_at <= now() - interval '15 minutes')
        OR (
            s.status = 'analyzing_pending'
            AND s.analysis_queued_at <= now() - interval '15 minutes'
        )
    )
      AND NOT EXISTS (
          SELECT 1
          FROM scan_runs sr
          WHERE sr.scan_id = s.id
            AND sr.status IN ('pending', 'running')
    )
    ORDER BY s.created_at ASC
    LIMIT 15
    FOR UPDATE OF s SKIP LOCKED
), queued_scans AS (
    UPDATE scans s
    SET status = 'analyzing_pending',
        analysis_queued_at = now()
    FROM ready_scans rs
    WHERE s.id = rs.id
    RETURNING s.id, s.project_id
)
SELECT
    s.id AS scan_id,
    s.project_id,
    p.brand_name,
    p.domain AS brand_domain,
    sr.id AS scan_run_id,
    sr.engine_id,
    sr.prompt_id,
    sr.provider_key_id,
    sr.status AS scan_run_status,
    sr.answer_text,
    sr.raw_response,
    pr.text AS prompt_text,
    pk.encrypted_key,
    pk.key_nonce
FROM queued_scans s
JOIN projects p ON p.id = s.project_id
JOIN scan_runs sr ON sr.scan_id = s.id
JOIN prompts pr ON pr.id = sr.prompt_id
JOIN provider_keys pk ON pk.id = sr.provider_key_id
ORDER BY s.id, sr.created_at ASC;

-- name: ClaimScanForAnalysis :one
UPDATE scans
SET started_at = now(),
    status = 'analyzing',
    analysis_queued_at = NULL
WHERE scans.id = sqlc.arg(id)
  AND status = 'analyzing_pending'
  AND NOT EXISTS (
      SELECT 1
      FROM scan_runs sr
      WHERE sr.scan_id = scans.id
        AND sr.status IN ('pending', 'running')
  )
RETURNING id;

-- name: IncrementScanCompletedRuns :exec
UPDATE scans
SET completed_runs = completed_runs + 1
WHERE id = sqlc.arg(id);

-- name: IncrementScanFailedRuns :exec
UPDATE scans
SET failed_runs = failed_runs + 1
WHERE id = sqlc.arg(id);
