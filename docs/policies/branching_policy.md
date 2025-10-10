# Branching & Chunking Policy

Scope: Applies to the Secrets Workflow repository. Keeps changes small, secure, and ephemeral-first.

## Purpose
- Make the branching model explicit and enforceable.
- Ensure security and teardown guarantees before merge.
- Keep history clean and releases predictable.

## Branch Model
- Trunk-based on `main`; short-lived branches only.
- No long-lived environment branches.
- `main` is protected: requires passing checks and 1+ review.

## Branch Naming
- Pattern: `feat|fix|docs|chore|security/<scope>/<tracker-id>-<slug>`
  - Examples:
    - `feat/minio/SW-142-add-ttl-bucket`
    - `security/ci/SW-201-tighten-gitleaks-baseline`
    - `fix/aws/SW-155-adjust-oidc-audience`
- Scope tags (examples): `aws`, `gcp`, `azure`, `minio`, `ci`, `docs`, `workflow`, `module`.
- Every branch references a tracker item from `secrets_workflow_project_tracker.md` (e.g., `SW-142`).

## Commits
- Conventional Commits: `type(scope): message`
  - Example: `feat(minio): add TTL for ephemeral buckets`
- Keep commits focused; no mixing unrelated changes.

## Pull Requests
- One logical change per PR.
- Target size: ≤ 400 LOC delta and ≤ 8 files; split otherwise.
- Must reference tracker ID (e.g., `SW-142`) in title or description.
- Update docs when behavior/flows change:
  - `secrets_workflow-README.md` (public overview)
  - `docs/runbooks/` (operator steps)
- Prefer squash-merge with a clean, conventional title.

## Required CI Gates (must pass)
- Pre-commit hooks (formatting, linting where configured).
- Terraform validation: `terraform init -backend=false` and `terraform validate`.
- Security/IaC scanners: `gitleaks`, `tfsec`, `checkov`.
- Unit/Terratest (when present).
- Teardown verification for ephemeral infra (see Teardown Guarantee).

Relevant reusable workflows in this repo:
- `.github/workflows/terraform-ephemeral.yml` (ephemeral Terraform)
- `.github/workflows/security-scans.yml` (gitleaks/tfsec/checkov)

## Ephemeral Environments
- Per-branch, isolated infra for validation.
- Naming: `<repo>-<branch-slug>-<shortsha>`.
- TTL enforced by workflow; must auto-destroy on PR close/merge.

## Secrets Hygiene
- No plaintext secrets in repo.
- All sensitive values must use `sops+age`.
- `gitleaks` failure blocks merge.

## Teardown Guarantee
- PR must include a link to CI teardown logs confirming successful destroy.
- Any teardown failure blocks merge until resolved.

## Releases
- Merge to `main` → tag SemVer release.
- Update published version for reusable workflows as applicable
  (e.g., bump `.github/workflows/terraform-ephemeral.yml` if changed).
- Update `secrets_workflow_project_tracker.md` and changelog.

## Hotfixes
- Branch from the release tag: `hotfix/<tracker-id>-<slug>`.
- PR into `main`; tag a patch release.
- Backport only if explicitly required.

## Merge & Cleanup
- Prefer squash-merge to keep history linear.
- Delete branches on merge.
- Ensure ephemeral resources are destroyed (CI verified) before deletion.

## Local Workflow (suggested)
1. Create branch per naming convention.
2. Run pre-commit and local checks: `gitleaks`, `terraform validate`.
3. Keep PR small; update docs alongside code.
4. Confirm ephemeral env stands up and tears down in CI.

---

Enforcement notes: Use GitHub Protected Branches to require status checks and reviews; configure required checks to include the security scans, Terraform validation, and ephemeral teardown jobs.

