# Runbook: Ephemeral State Management

This runbook explains **ephemeral Terraform state management** in this repository and separates:

- what is implemented today in `.github/workflows/terraform-ephemeral.yml`, and
- stronger hardening guidance that is valuable but not yet fully implemented.

## 1. The Principle of Ephemeral State

Just as we use ephemeral, short-lived credentials via OIDC, we should also treat our Terraform state as ephemeral for temporary infrastructure. The state file is critical for the lifecycle of a single workflow run (`plan` → `apply` → `destroy`), but it should not persist after the run is complete.

This aligns with our core philosophy of **"no long-lived infrastructure,"** which extends to the metadata associated with that infrastructure.

## 2. State Management Patterns

There are two primary patterns for ephemeral state demonstrated in this project.

### Pattern 1: The Default Local Backend (Simple but Fragile)

By default, if no backend is configured, Terraform creates a `terraform.tfstate` file in the local working directory.

* **How it works:** `terraform init` creates the file, and subsequent commands (`plan`, `apply`, `destroy`) read from and write to it.
* **Pros:** Zero configuration required.
* **Cons:** In a CI environment, this is fragile. If the runner experiences an issue or if you need to inspect the state for debugging, accessing this file is difficult. It's an implicit dependency on the runner's transient filesystem.

### Pattern 2: The MinIO-in-CI Backend (Current Project Pattern)

This is the current recommended pattern for ephemeral state in CI. We run an S3-compatible object storage service (MinIO) as a container directly within the workflow job.

* **How it works:** The workflow starts a MinIO container as a `service`. Terraform is then configured to use this service as a standard S3 remote backend.
* **Pros:**

  * **Isolation:** Each workflow run gets its own private, clean state backend. There is no risk of state conflicts between runs.
  * **Robustness:** It behaves like a true remote backend, providing better consistency than a local file.
  * **Zero Cost:** The MinIO container is created and destroyed with the CI job, incurring no persistent storage costs.
  * **Realistic Simulation:** It mimics the pattern of using a real S3 bucket for state, making CI a high-fidelity representation of production.
  * **Closer to remote-backend behavior than local state:** while still ephemeral to the workflow run.

### Current Implementation Status (Reality Check)

The reusable workflow at `.github/workflows/terraform-ephemeral.yml` currently implements:

- MinIO service-backed state (`endpoint = "http://minio:9000"`).
- A per-run bucket name (`tf-state-${{ github.run_id }}`).
- Workflow-level `concurrency` to reduce parallel collision risk.
- Conditional teardown (`terraform destroy`) when `apply=true` and `destroy_on_complete=true`.

The same workflow currently **does not yet implement all hardening controls** described later in this runbook. Today it still uses:

- static MinIO credentials (`minioadmin`),
- `minio/minio:latest` rather than a pinned image digest,
- host port publishing (`9000:9000`),
- no explicit post-run local backend/state file scrubbing step.

> **Note:** The S3 backend provides state locking only when used with DynamoDB. Because we run in a single job and do not provision DynamoDB, we serialize runs via a workflow `concurrency` group to avoid parallel state access.

## 3. Target-State Hardening Guidance

The following pattern is **target-state guidance** for strengthening the current workflow. Treat it as a hardened reference pattern to adopt incrementally, not as a statement of current implementation parity.

### 1) Minimal Global Permissions + Concurrency Guard

```yaml
permissions:
  contents: read

concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: false
```

### 2) Generate Per-run MinIO Credentials

```yaml
    steps:
      - id: secrets
        name: Generate per-run MinIO credentials
        run: |
          USER=minio-${RANDOM}${RANDOM}
          PASS=$(openssl rand -hex 24)
          echo "MINIO_USER=$USER" >> $GITHUB_OUTPUT
          echo "MINIO_PASS=$PASS" >> $GITHUB_OUTPUT
```

### 3) Start MinIO Service (Pinned Digest, No Host Ports)

> **Note:** MinIO service containers are supported only on GitHub-hosted Linux runners.

```yaml
services:
  minio:
    # Pinned digest for supply-chain security. See Docker Hub for latest.
    image: minio/minio@sha256:2061b442f48a280793020291c1273ed5b364d0458032b339953a18a57b224445
    env:
      MINIO_ROOT_USER: ${{ steps.secrets.outputs.MINIO_USER }}
      MINIO_ROOT_PASSWORD: ${{ steps.secrets.outputs.MINIO_PASS }}
    # Health checks ensure the service is ready before steps run.
    options: >-
      --health-cmd="curl -f http://localhost:9000/minio/health/ready || exit 1"
      --health-interval=2s --health-retries=30 --health-timeout=2s
```

> For internal networking we do **not** publish host ports. The container is reachable as `http://minio:9000` from workflow steps only.

### 4) Backend Configuration Generation

```bash
BUCKET_NAME="tf-state-${{ github.run_id }}"

cat <<EOF > backend.tf
terraform {
  backend "s3" {
    bucket                      = "$BUCKET_NAME"
    key                         = "terraform.tfstate"
    region                      = "us-east-1"
    endpoint                    = "http://minio:9000"
    skip_credentials_validation = true
    skip_region_validation      = true
    force_path_style            = true
  }
}
EOF
```

### 5) Authenticate and Run Terraform

```yaml
- name: Terraform Init, Plan, and Apply
  env:
    AWS_ACCESS_KEY_ID: ${{ steps.secrets.outputs.MINIO_USER }}
    AWS_SECRET_ACCESS_KEY: ${{ steps.secrets.outputs.MINIO_PASS }}
  run: |
    terraform init -reconfigure
    terraform plan -out=tfplan
    terraform apply -auto-approve tfplan
```

### 6) Teardown and Cleanup (Defense in Depth)

```yaml
- name: Always destroy (teardown guarantee)
  if: always()
  env:
    AWS_ACCESS_KEY_ID: ${{ steps.secrets.outputs.MINIO_USER }}
    AWS_SECRET_ACCESS_KEY: ${{ steps.secrets.outputs.MINIO_PASS }}
  run: terraform destroy -auto-approve || true

- name: Scrub backend and state
  if: always()
  run: |
    rm -f backend.tf || true
    rm -rf .terraform || true
    rm -f .terraform.lock.hcl || true
    rm -f terraform.tfstate terraform.tfstate.backup || true
```

### 7) Optional User Rotation for Defense in Depth

If the `services.env` block cannot consume generated secrets directly, initialize MinIO with static creds and then rotate them immediately after the container is ready.

```yaml
- name: Wait for MinIO readiness
  run: |
    for i in $(seq 1 60); do
      curl -sf http://minio:9000/minio/health/ready && exit 0
      sleep 1
    done
    echo "MinIO not ready in time" >&2
    exit 1

- name: Rotate to random user for Terraform
  run: |
    curl -LO https://dl.min.io/client/mc/release/linux-amd64/mc && chmod +x mc && sudo mv mc /usr/local/bin/
    mc alias set ci http://minio:9000 minioadmin minioadmin
    USER=${{ steps.secrets.outputs.MINIO_USER }}
    PASS=${{ steps.secrets.outputs.MINIO_PASS }}
    mc admin user add ci "$USER" "$PASS"
```

## 4. Target-State Security and Shift-Left Principles

* **Pin all GitHub Actions by commit SHA** (use Dependabot to update SHAs).
* **Set minimal permissions:** `contents: read` globally, elevate per-job only.
* **No host port publishing**—networking is internal to the job.
* **Generate random credentials per run**, scoped to steps.
* **Avoid uploading any Terraform state** or backend files as artifacts.
* **Use default Terraform log level**, avoid `TF_LOG=TRACE` leaks.
* **Pin container digests** instead of `:latest`.
* **Planned workflow:** enable `.github/workflows/security-scans.yml` (gitleaks, tfsec, checkov) for this branch **after it is implemented**.
* **Use pre-commit hooks** for `terraform validate` and `gitleaks`.

## 5. Summary: Current State vs Target State

Current workflow reality in this repo provides ephemeral-in-job MinIO-backed state and conditional teardown, but does not yet implement every control in the hardened reference.

Target-state hardening aims to ensure:

* Supply-chain integrity (pinned versions and images)
* Scoped, per-run credentials
* No leaked state or logs
* Serialized, isolated runs
* Defense-in-depth cleanup

As hardening work is implemented, update this runbook and `.github/workflows/terraform-ephemeral.yml` together so documented guarantees and actual behavior stay aligned.

See also:

* [Runbook: OIDC Federation](./oidc_federation.md)
* [Reusable Workflow: terraform-ephemeral.yml](../../.github/workflows/terraform-ephemeral.yml)
* [Runbook: Local Development Setup](./local_dev_setup.md) 
