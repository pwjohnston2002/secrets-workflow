# Secrets Workflow Project Tracker

## Project Complexity Level

- [X] **Standard**: Multi-phase project with moderate complexity

---

## 1. Project Overview
**Project Name:** Secrets Workflow  
**Repository:** `secrets_workflow`  
**Objective:** Establish a robust, secure, and automated **ephemeral secrets and state management methodology and reference implementations** to support the Workflow Toolbelt project and all dependent repositories. This project will deliver foundational patterns for just-in-time credential generation, temporary state handling, and secure configuration, ensuring no long-lived cloud infrastructure for secrets or state.

### Goals
- Build a secure, **ephemeral-first** secrets and state management methodology.
- Provide **reference implementations** for OIDC-based authentication, MinIO-in-CI state, and sops+age secrets.
- Integrate seamlessly with the Workflow Toolbelt and other repos via **reusable workflows and templates**.
- Enable environment-agnostic workflows with consistent, **just-in-time** secret and state handling.
- Enforce compliance and least-privilege security practices through **policy-as-code and automated scanning**.

---

## Context Engineering Foundation

### Context Engineering Completeness Check

- [ ] **User Persona Definition**: Complete profile with expertise level, goals, constraints (Self: SRE/DevSecOps, Portfolio Builder, Cost-Conscious)
- [ ] **Domain Context Loaded**: Ephemeral infrastructure, CI/CD security, Terraform, OIDC, sops/age
- [ ] **Tool/Resource Context**: GitHub Actions, AWS/GCP/Azure OIDC, MinIO, sops, age, Terratest, Checkov, tfsec, Gitleaks
- [ ] **Historical Context**: Previous discussions on "no long-lived infra," "ephemeral creds," and the need for a foundational repo (`secrets.md`)
- [ ] **Edge Case Documentation**: Handling local dev with ephemeral patterns, sops key management, cross-cloud OIDC setup
- [ ] **Success Pattern Recognition**: OIDC for auth, MinIO for ephemeral state, sops for encrypted config, reusable workflows for consumption

---

Current and planned projects require secure handling of API keys, database credentials, tokens, and sensitive configuration data. At present, secret storage is fragmented and manual, creating:
- Risk of credential leaks.
- Lack of version control separation between code and secrets.
- No automated rotation or auditing for ephemeral credentials.

The `secrets_workflow` repository addresses these gaps by creating an extensible, automated **methodology and reference implementations for ephemeral secrets and state management**.

---

## 3. Scope
### In-Scope
- Implementation of **ephemeral secrets and state management patterns**.
- Support for local development and CI/CD environments (no long-lived cloud resources).
- CI/CD integration with Workflow Toolbelt via **reusable workflows**.
- Patterns for **just-in-time credential generation and destruction**.
- Audit trails for ephemeral access and usage (via cloud provider logs).

### Out-of-Scope
- Hardcoded secrets in legacy repos (separate remediation plan).
- **Deployment or maintenance of a long-lived, centralized secrets manager (e.g., HashiCorp Vault, AWS Secrets Manager).**
- Full PKI certificate management (focus on ephemeral TLS for testing).

---

## 3.5 Integration with Workflow Toolbelt

This repository is a foundational dependency for the `workflow-toolbelt`. The integration points are as follows:

- **Consumption Model:** The `workflow-toolbelt` will consume the reusable workflows (e.g., `terraform-ephemeral.yml`) from this repository via `workflow_call`.
- **Dependency Sequence:**
    1.  **`secrets-workflow` (This Repo):** Must first implement and version the core reusable workflows for OIDC auth, ephemeral state, and security scanning (Phase 1 & 2).
    2.  **`workflow-toolbelt`:** Can then integrate these versioned, reusable workflows to provide standardized CI/CD pipelines for projects that use it.
- **Cross-Repo Milestone:** Phase 1 of this project is a direct dependency for starting Phase 2 of the `workflow-toolbelt` project.
- **Shared Principles:** The `pre-commit` hooks and security scanning tools (`gitleaks`, `tfsec`, `checkov`) defined here will become the standard enforced by the `workflow-toolbelt`'s CI templates.

---

## 4. Proposed Architecture
### 4.1 Core Components
- **Ephemeral Authentication:** OIDC Federation (GitHub Actions → AWS STS AssumeRole, Azure AD Workload Identity, GCP Workload Identity Federation) for just-in-time, short-lived cloud credentials.
- **Ephemeral State Management:**
    - Terraform `local` backend with encrypted artifact handoff in CI.
    - MinIO-in-CI (S3-compatible container spun up as a sidecar, then destroyed) for temporary remote state.
- **Ephemeral Secrets Storage:**
    - `sops+age` for encrypting static configuration examples and sensitive files in-repo (decrypted at runtime).
    - Runtime generation of secrets (e.g., Terraform `random_password`, `tls_private_key`) for throwaway environments.
    - Environment injection via CI (job-scoped env vars, temporary files).
- **Reusable CI Components:**
    - GitHub Actions `workflow_call` for init/plan/apply/test/destroy pipelines.
    - Composite Actions for common tasks (e.g., sops decryption, tool setup).
- **Policy-as-Code & Scanning:** `checkov`, `tfsec`, `gitleaks` integrated into CI/CD and `pre-commit` hooks.

### 4.1.1 Local Development Experience
- **Local Ephemeral Credentials:** Developers will use tools like `aws-vault` or `Leapp` to generate short-lived credentials for local Terraform runs, mimicking the OIDC flow in CI.
- **Local State:** For local runs, the default will be the standard Terraform `local` backend (`terraform.tfstate` file), which should be included in `.gitignore`.
- **Local Secrets (`sops`):** Developers will need access to a shared `age` private key to decrypt configuration files. The initial recommendation is to manage this key via a secure team password manager (e.g., 1Password), with a runbook detailing the setup.

### 4.1.2 Observability & Auditing
- **Ephemeral Logging:** CI/CD pipelines will be configured to output key audit events.
    - **OIDC Session Logging:** Log the unique session name/ID from the cloud provider's STS service to trace actions back to a specific workflow run.
    - **MinIO Access Logs:** When using MinIO-in-CI, capture and stream its access logs as a job artifact for debugging state-related issues.
    - **Secret Usage:** Log when `sops` decryption occurs.
- **Log Aggregation:** For demonstration purposes, logs will be output to the CI console and stored as job artifacts. In a real-world scenario, these would be shipped to a central logging platform (e.g., Splunk, Datadog).

### 4.2 Data Flow
1.  **Trigger Pipeline:** A CI/CD pipeline (e.g., GitHub Actions) is triggered.
2.  **Ephemeral Authentication:** The CI runner uses OIDC federation to obtain short-lived cloud credentials directly from the cloud provider (e.g., AWS STS). No static credentials are stored.
3.  **Ephemeral State Setup:**
    -   For Terraform, either a `local` backend is used, or a MinIO container is spun up as a service within the CI job to act as a temporary S3-compatible backend.
4.  **Ephemeral Secret Injection:**
    -   Static configuration (e.g., example API keys) encrypted with `sops+age` is decrypted at runtime within the CI job.
    -   Dynamic secrets (e.g., database passwords) are generated by Terraform during the `apply` phase.
    -   Secrets are injected as environment variables or temporary files, never persisted.
5.  **Infrastructure Provisioning & Testing:**
    -   IaC (e.g., Terraform) uses the ephemeral credentials and secrets to provision cloud resources.
    -   Automated tests (e.g., Terratest) and security scans (`checkov`, `tfsec`, `gitleaks`) run against the provisioned infrastructure.
6.  **Teardown:**
    -   `terraform destroy` is executed to remove all provisioned cloud resources.
    -   Ephemeral state files (local or MinIO) are deleted.
    -   Short-lived cloud credentials expire automatically.
7.  **Safety Nets:** `pre-commit` hooks and CI scans prevent accidental secret leaks or insecure configurations from being committed.

---

### Phase 1 — Foundational Patterns (MVP)
- [ ] **Implement OIDC federation for CI/CD authentication** (e.g., GitHub Actions to AWS STS).
- [ ] **Establish `sops+age` for encrypting static configuration examples** and demonstrating in-repo secret management.
- [ ] **Define patterns for Terraform `local` backend and MinIO-in-CI** for ephemeral state management.
- [ ] **Create runbook for local development setup** (`docs/runbooks/local_dev_setup.md`) covering `aws-vault` and `sops+age` key access.
- [ ] Create initial `PRINCIPLES.md` outlining the "no long-lived infra" philosophy.

### Phase 1.5 — Validation & Testing
- [ ] Implement **Terratest examples** demonstrating create→assert→destroy patterns for ephemeral infrastructure.
- [ ] Build a **test harness** for validating that the `destroy` step leaves zero cloud resources.

### Phase 2 — Workflow Integration & Hardening
- [ ] **Integrate reusable workflows** (`workflow_call`) for ephemeral Terraform init/plan/apply/destroy into the Workflow Toolbelt.
- [ ] Add **dynamic secret injection via runtime generation and sops decryption** into CI/CD pipelines.
- [ ] Set up **least-privilege IAM roles** for OIDC federation, scoped to ephemeral operations.
- [ ] Document usage patterns for developers for consuming ephemeral patterns.
- [ ] **Implement Observability:** Add logging for OIDC session IDs and MinIO access to CI workflows.
- [ ] Integrate **policy-as-code tools** (`checkov`, `tfsec`) into CI workflows.
- [ ] Integrate **`gitleaks`** into CI and `pre-commit` hooks.
- [ ] Develop **basic cross-cloud OIDC→STS examples** (GCP and Azure) to prove pattern portability.

### Phase 3 — Advanced Features & Key Management
- [ ] Create **Terraform modules** for common ephemeral patterns (e.g., `oidc-role-aws`, `ephemeral-minio`).
- [ ] Refine **teardown guarantees** with advanced Terratest validation.
- [ ] **Evolve Key Management:**
    - [ ] Propose and document advanced `age` key management strategies (e.g., bootstrapping with OIDC-signed secrets).
    - [ ] Explore using a cloud KMS to encrypt/decrypt the `age` private key, removing the need to share it directly.

---

## 6. Deliverables
- `secrets_workflow` repository with **reference implementations for ephemeral secrets and state management**.
- **Reusable CI workflows** (`workflow_call`) for ephemeral builds (init/plan/apply/destroy).
- **Terraform modules** demonstrating ephemeral patterns (e.g., OIDC roles, MinIO-in-CI backend).
- **Documentation:**
    - **Developer Onboarding Runbooks** (e.g., `docs/runbooks/local_dev_setup.md`).
    - **Principles & Guides** (e.g., `PRINCIPLES.md`, `ephemeral_state.md`, `secrets_with_sops.md`).
    - **Architectural Decision Records (ADRs)** for key design choices.

---

## 7. Risks & Mitigation
| **Risk** | **Impact** | **Mitigation** |
|-----------|------------|----------------|
| **Misconfiguration of OIDC trust policies** | High | Strict IAM policies, automated scanning (`checkov`), Terratest validation. |
| **Accidental `sops` key exposure** | High | Strong local key management practices, `gitleaks` in pre-commit/CI, clear documentation. |
| **MinIO-in-CI stability/performance** | Medium | Use official Docker images, resource limits, robust error handling in CI. |
| **Incomplete resource teardown** | High | Comprehensive `terraform destroy` in CI, Terratest assertions for zero remaining resources. |
| **Developer onboarding friction** | Medium | Provide clear documentation, examples, and pre-configured templates. |
| **Complexity of ephemeral patterns** | Medium | Start with minimal patterns, iterative expansion, focus on clear examples. |

---

## 8. Dependencies
- Workflow Toolbelt repo (`workflow_toolbelt_project_tracker`)
- Cloud provider IAM configurations
- CI/CD pipelines (GitHub Actions)
- `sops`, `age`, `MinIO`

---

## 9. Success Metrics
- **Security:** 0 leaked credentials.
- **Adoption:** 100% new repos consume ephemeral patterns via reusable workflows.
- **Automation:** 100% ephemeral credentials generated and destroyed automatically per run.
- **Observability:** Full audit trails for all ephemeral sessions and secret access events.
- **Cost-Efficiency:** Zero cloud costs when projects are not actively running.

---

## 10. Open Questions
- **Initial `age` key management:** For a solo developer, the key can be stored locally and backed up in a password manager. For teams, the recommended starting point is sharing the private key via a secure, audited system like a team password manager (e.g., 1Password).
- **Advanced `age` key management:** How do we evolve beyond a shared private key for larger teams? Future phases should investigate using cloud KMS to encrypt/decrypt the `age` key or using a tool that can bootstrap encryption with OIDC identities.
- What are the best practices for handling sensitive data that *must* persist between ephemeral runs (e.g., database backups)? (Likely out of scope for this repo, but a common question).
- How to ensure cross-cloud OIDC configurations are consistent and secure?
- What level of detail is needed for `sops` examples in a public repo without exposing sensitive patterns?

---

## 11. Next Steps
- Finalize the core ephemeral patterns (OIDC, MinIO-in-CI, sops+age).
- Create initial repository scaffold for `secrets_workflow`.
- Implement initial reusable CI workflows for ephemeral Terraform apply/destroy.
- Integrate with Workflow Toolbelt by consuming the new reusable workflows.
- Review security and compliance requirements with stakeholders.

---

## Context-Aware Collaboration Framework

### Context Interaction Optimization

- **Preferred Interaction Style:** [Research assistant | Pair programmer | Devil's advocate | Creative partner]
- **Context Refresh Triggers:** [When to reload full context vs. incremental updates]
- **Suggestion Integration Process:** [How to evaluate and incorporate external/tooling-driven recommendations]
- **Decision Boundaries:** [What can be automated vs. what requires human judgment]
- **Assistance Preferences:** [Proactive suggestions, reactive responses, or mixed]

---

## Context Evolution Tracking

### Dynamic Context Management

- **Context Version:** v1.0
- **Last Context Audit:** 2025-08-21
- **Context Debt:** [Areas where context is incomplete or outdated]
- **Proactive Context Gaps:** [AI-identified areas needing attention]
- **Context Compression Events:** [When/how to summarize older context]
- **Evolution Triggers:** [What events cause this document to change?]

---

## Cross-Domain Learning Integration

### Cross-Domain Learning Patterns

- **Transferable Insights:** [Lessons that apply across domains]
- **Domain-Specific Adaptations:** [How this project differs from your template]
- **Pattern Recognition:** [Recurring themes across your projects]
- **Template Evolution Notes:** [How this project is improving your framework]
- **Methodology Transfers:** [Techniques from other domains that apply here]

---

## Versioning & Evolution Log

### Template & Project Evolution

- **Template Version:** v1.0 (for this specific tracker)
- **Project Version:** v0.1 (initial draft)
- **Major Context Updates:**
  - 2025-08-21: Initial creation based on `secrets.md` and `v3.0` template.
- **Template Improvements Applied:**
  - 2025-08-21: Incorporated `Project Complexity Level`, `Context Engineering Foundation`, `AI Collaboration Framework`, `Context Evolution Tracking`, `Cross-Domain Learning Integration`.
- **Cross-Project Learnings Applied:**
  - 2025-08-21: Applied lessons from `workflow_toolbelt_project_tracker.md` regarding phased implementation and clear task definitions.

---

## Appendix A: Testing Strategy Matrix

This matrix outlines the testing strategy for the core components of the ephemeral workflow.

| Component / Feature | Phase | Test Type | Tool | Goal |
|---|---|---|---|---|
| **OIDC Federation** | 1.5 | Integration | Terratest | Verify that a Terraform `apply` using OIDC credentials can successfully create and destroy a simple cloud resource (e.g., S3 bucket). |
| **OIDC Federation (Mocked)** | 2 | Unit | Go Mocks | (Optional) Mock the cloud provider's STS endpoint to speed up CI runs and test failure conditions without incurring API call costs. |
| **MinIO-in-CI** | 1.5 | Integration | Terratest | Confirm that Terraform can initialize against the MinIO service container, write state, and that the state is accessible during the run. |
| **MinIO-in-CI Stability** | 2 | Stress Test | CI Workflow | Run multiple parallel jobs that use the MinIO service to check for race conditions or performance bottlenecks. |
| **`sops+age` Decryption** | 1.5 | Unit | Shell Script / Go Test | Verify that an encrypted file can be successfully decrypted in the CI environment using a test key. |
| **Full Teardown Guarantee** | 2 | End-to-End | Terratest | After `terraform destroy`, use the cloud provider's SDK within the test to list resources and assert that no resources created by the test remain. |
| **Policy-as-Code Scanners** | 2 | Integration | CI Workflow | Run scanners (`checkov`, `tfsec`) against intentionally insecure Terraform code to ensure the pipeline fails as expected. |
