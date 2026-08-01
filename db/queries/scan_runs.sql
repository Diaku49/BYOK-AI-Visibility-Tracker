-- name: GetScanRunByID :one
SELECT *
FROM scan_runs
WHERE id = $1;

-- name: UpdateScanRunStateByID :one
UPDATE scan_runs
SET
    status = sqlc.arg(status),
    error = sqlc.narg(error),
    started_at = CASE
        WHEN sqlc.arg(status) = 'running' THEN now()
        ELSE started_at
    END,
    finished_at = CASE
        WHEN sqlc.arg(status) IN ('completed', 'failed', 'canceled') THEN now()
        ELSE finished_at
    END
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateScanRun :one
UPDATE scan_runs
SET
    scan_id = sqlc.arg(scan_id),
    engine_id = sqlc.arg(engine_id),
    prompt_id = sqlc.arg(prompt_id),
    provider_key_id = sqlc.arg(provider_key_id),
    try_number = sqlc.arg(try_number),
    status = sqlc.arg(status),
    answer_text = sqlc.narg(answer_text),
    raw_response = sqlc.narg(raw_response),
    brand_mentioned = sqlc.narg(brand_mentioned),
    brand_domain_cited = sqlc.narg(brand_domain_cited),
    competitors_mentioned = sqlc.narg(competitors_mentioned),
    cited_domains = sqlc.narg(cited_domains),
    error = sqlc.narg(error),
    started_at = sqlc.narg(started_at),
    finished_at = sqlc.narg(finished_at)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateScanRunAnalysis :exec
UPDATE scan_runs
SET
    brand_mentioned = sqlc.arg(brand_mentioned),
    brand_domain_cited = sqlc.arg(brand_domain_cited),
    competitors_mentioned = sqlc.arg(competitors_mentioned),
    cited_domains = sqlc.arg(cited_domains)
WHERE id = sqlc.arg(id);

-- name: GetScansForAnalysis :many
WITH ready_scans AS (
    SELECT s.id
    FROM scans s
    WHERE s.status = 'running'
      AND NOT EXISTS (
          SELECT 1 FROM scan_runs sr
          WHERE sr.scan_id = s.id
            AND sr.status NOT IN ('completed', 'failed')
      )
    FOR UPDATE OF s SKIP LOCKED
), claimed AS (
    UPDATE scans s
    SET status = 'analyzing'
    FROM ready_scans rs
    WHERE s.id = rs.id
    RETURNING s.id, s.project_id
)
SELECT
    c.id AS scan_id,
    c.project_id,
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
FROM claimed c
JOIN projects p ON p.id = c.project_id
JOIN scan_runs sr ON sr.scan_id = c.id
JOIN prompts pr ON pr.id = sr.prompt_id
JOIN provider_keys pk ON pk.id = sr.provider_key_id
ORDER BY c.id, sr.created_at ASC;

-- name: GetScansForWorkers :many
WITH claimed AS (
    SELECT s.id
    FROM scans s
    WHERE s.status = 'pending'
       OR (s.status = 'running' AND s.started_at <= now() - interval '15 minutes')
    ORDER BY s.created_at ASC
    FOR UPDATE OF s SKIP LOCKED
), updated AS (
    UPDATE scans s
    SET status = 'running',
        started_at = now()
    FROM claimed c
    WHERE s.id = c.id
    RETURNING s.id, s.project_id, s.status, s.tries_per_prompt
)
SELECT
    p.brand_name AS brand_name,
    p.domain AS brand_domain,
    s.id AS scan_id,
    sr.id AS scan_run_id,
    sr.engine_id,
    sr.prompt_id,
    sr.provider_key_id,
    sr.try_number,
    sr.status AS scan_run_status,
    s.status AS scan_status,
    s.tries_per_prompt,
    p.language,
    p.region,
    pr.text AS prompt_text,
    pr.active AS prompt_active,
    pk.encrypted_key,
    pk.key_nonce,
    pk.active AS provider_key_active,
    pk.monthly_run_limit,
    pk.monthly_runs_used
FROM updated s
JOIN projects p
    ON p.id = s.project_id
JOIN scan_runs sr
    ON sr.scan_id = s.id
    AND sr.status IN ('pending', 'running')
JOIN prompts pr
    ON pr.id = sr.prompt_id
JOIN provider_keys pk
    ON pk.id = sr.provider_key_id
ORDER BY sr.created_at ASC;
