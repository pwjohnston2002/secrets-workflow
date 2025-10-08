# Runbook: Ephemeral State Management

This runbook explains the concept and implementation of ephemeral state management for Terraform within this project's CI/CD workflows.

## 1. The Principle of Ephemeral State

Just as we use ephemeral, short-lived credentials via OIDC, we should also treat our Terraform state as ephemeral for temporary infrastructure. The state file is critical for the lifecycle of a single workflow run (`plan` -> `apply` -> `destroy`), but it should not persist after the run is complete.

This aligns with our core philosophy of **"no long-lived infrastructure,"** which extends to the metadata associated with that infrastructure.

## 2. State Management Patterns

There are two primary patterns for ephemeral state demonstrated in this project.

### Pattern 1: The Default Local Backend (Simple but Fragile)

By default, if no backend is configured, Terraform creates a `terraform.tfstate` file in the local working directory.

*   **How it works:** `terraform init` creates the file, and subsequent commands (`plan`, `apply`, `destroy`) read from and write to it.
*   **Pros:** Zero configuration required.
*   **Cons:** In a CI environment, this is fragile. If the runner experiences an issue or if you need to inspect the state for debugging, accessing this file is difficult. It's an implicit dependency on the runner's transient filesystem.

### Pattern 2: The MinIO-in-CI Backend (Robust and Isolated)

This is the recommended and more robust pattern for ephemeral state. We run an S3-compatible object storage service (MinIO) as a container directly within the CI job.

*   **How it works:** The workflow starts a MinIO container as a `service`. Terraform is then configured to use this service as a standard S3 remote backend.
*   **Pros:**
    *   **Isolation:** Each workflow run gets its own private, clean state backend. There is no risk of state conflicts between runs.
    *   **Robustness:** It behaves like a true remote backend, providing better locking and consistency than a local file.
    *   **Zero Cost:** The MinIO container is created and destroyed with the CI job, incurring no persistent storage costs.
    *   **Realistic Simulation:** It perfectly mimics the pattern of using a real S3 bucket for state, making the CI environment a high-fidelity representation of a production setup.

## 3. Implementation in the `terraform-ephemeral.yml` Workflow

Our reusable workflow implements the MinIO-in-CI pattern automatically. Here’s how the pieces fit together.

### Step 1: Start the MinIO Service Container

A `services` block in the workflow file starts the MinIO container. GitHub Actions handles networking, so it's available at the hostname `minio`.

```yaml
services:
  minio:
    image: minio/minio:latest
    ports:
      - 9000:9000
    command: server /data
    env:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
```

### Step 2: Generate the Backend Configuration

A step dynamically creates a `backend.tf` file, configuring Terraform to use the `minio` service. A unique bucket name is generated for each run to guarantee isolation.

```bash
# From the "Generate Terraform Backend Config" step
BUCKET_NAME="tf-state-${{ github.run_id }}"

cat <<EOF > backend.tf
terraform {
  backend "s3" {
    bucket                      = "$BUCKET_NAME"
    key                         = "terraform.tfstate"
    region                      = "us-east-1" # Ignored by MinIO
    endpoint                    = "http://minio:9000"
    skip_credentials_validation = true
    skip_region_validation      = true
    force_path_style            = true
  }
}
EOF
```

### Step 3: Authenticate Terraform with MinIO

The `terraform init`, `apply`, and `destroy` steps are given environment variables that the S3 backend provider uses for authentication. These credentials match the ones defined in the service container.

```yaml
- name: Terraform Init, Plan, and Apply
  env:
    AWS_ACCESS_KEY_ID: minioadmin
    AWS_SECRET_ACCESS_KEY: minioadmin
  run: |
    terraform init -reconfigure
    # ...
```
The MinIO container automatically shuts down with the job, ensuring no residual data or credentials remain between runs.

This completes the pattern, providing a secure, isolated, and truly ephemeral state backend for every workflow run.

See also:
 - [Runbook: OIDC Federation](./oidc_federation.md)
 - [Reusable Workflow: terraform-ephemeral.yml](../../.github/workflows/terraform-ephemeral.yml)