# Deprecated
# MOVED
This file has been migrated to AGENTS.md.

# instructions-secrets\_workflow\.md

📌 **NOTE:** This file is the **single point of reference** for project context.
Consult here first before interpreting other files.
When in doubt, AI should reread this file before interpreting repo contents.

📎 **Repo Root:** [Secrets Workflow on GitHub](https://github.com/pwjohnston2002/secrets-workflow)  
This link provides access to the full repository (README, workflows, modules, and docs).  
Use this as the entry point when scanning or reviewing the repo as a whole.

---

## 1. Project Overview


**Project Name:** Secrets Workflow
**Mission Statement:** Provide a secure, ephemeral-first secrets and state management library that integrates with the Workflow Toolbelt and all dependent portfolio projects.
**Scope & Deliverables:**

* [Link to `secrets_workflow-README.md`](./secrets_workflow-README.md) for the public-facing philosophy and principles.
* [Link to `secrets_workflow_project_tracker.md`](./secrets_workflow_project_tracker.md) for roadmap, milestones, and detailed execution plan.

---

## 2. File Map & Roles

| File                                        | Purpose                                                 | Notes                                   |
| ------------------------------------------- | ------------------------------------------------------- | --------------------------------------- |
| `secrets_workflow-README.md`                | Public overview & philosophy                            | For human readers first                 |
| `PRINCIPLES.md`                             | Core design philosophy & non-negotiable rules           | Guiding document for contributors & architects |
| `secrets_workflow_project_tracker.md`       | Roadmap, tasks, risks, testing                          | Execution guide for human + AI          |
| `docs/runbooks/`                            | Step-by-step guides (local dev, instance access)        | Onboarding and operational focus        |
| `docs/policies/`                            | Formal policies for branching, releases, etc.           | Authoritative rules for contribution    |
| `.github/workflows/terraform-ephemeral.yml` | Reusable CI workflow for ephemeral Terraform            | Published for consumption (`@v0.1.0`)   |
| `.github/workflows/security-scans.yml`      | Planned reusable CI workflow for IaC + secret scanning (not yet implemented) | Intended to include `gitleaks`, `tfsec`, `checkov` |
| `modules/`                                  | Minimal Terraform modules (OIDC roles, ephemeral MinIO) | Example-driven, portfolio-aligned       |

---

## 3. Context Engineering Snapshot

* **User Persona:** Mid-career SRE shifting to security-focused infra/DevSecOps; building portfolio; highly cost-conscious.
* **Domain Context:** Ephemeral infrastructure, CI/CD security, Terraform, OIDC federation, `sops/age`.
* **Tool Context:** GitHub Actions, AWS/GCP/Azure OIDC, MinIO, `sops`, `age`, Terratest, `checkov`, `tfsec`, `gitleaks`.
* **Historical Context:** Rooted in the guiding principle of *no long-lived infrastructure*. Extends from Workflow Toolbelt groundwork.
* **Edge Cases:** Local dev with ephemeral creds, sharing `age` keys securely, teardown failures, cross-cloud OIDC quirks.
* **Success Patterns:** OIDC for auth, MinIO for temporary state, `sops` for encrypted config, reusable workflows for adoption.

---

## 4. AI Collaboration Framework

* **Preferred AI Role:** collaborator, research assistant, and mentor. 
* **Decision Boundaries:**

  * Human and AI discuss architecture, risks, design patterns.
  * Human writes/commits code and repo structure.
* **Context Refresh Triggers:** At major milestones (Phases 1, 1.5, 2 in tracker).
* **Integration Notes:** All AI guidance should map back to tracker and README for alignment.

---

## 5. Operating Instructions

**Workflow for human+AI collaboration:**

1. Draft code/docs locally (Human + Code Assistant for pair-programming).
2. Run pre-commit hooks and local tests (`gitleaks`, `terraform validate`, `ruff`).
3. Push to branch; GitHub Actions CI runs ephemeral workflows.
4. Use AI (here) for design deep dives, debugging, and conceptual reviews.
5. Before handing to Codex for repo-wide review, prepare a **Handoff Note** (goal, scope, validation).

### IAM Deployment Preflight (Required for Cloud Apply)

Before any real cloud deployment (Terraform apply, CloudFormation deploy, CDK deploy), run the IAM preflight gates in:
`docs/preflight/aws_iam_deployment_readiness.md`.

Record the results using the template at `docs/preflight/iam_preflight_run_template.md`.

AI should treat the gates as a mandatory checklist and must not recommend applying changes to AWS until the checklist is explicitly walked and marked pass/fail.

**Definition of Done:**

* All tests/scans pass (`checkov`, `tfsec`, `gitleaks`, teardown).
* Documentation updated (`README`, runbooks, principles).
* Project tracker updated.
* Tagged release (semver) published.
* [`secrets_workflow_project_tracker.md`]( https://github.com/pwjohnston2002/secrets-workflow/blob/main/secrets_workflow_project_tracker.md ) updated to reflect progress (authoritative task-level detail). #will update link when project repo is public, not yet.

---

### Branching & Commit Scope

To ensure consistent development practices and clear communication, the branching and commit strategy adapts to the current project phase. This allows for rapid iteration during exploratory stages and precise control during production-level development.

For specific naming conventions, PR requirements, and CI gates, see the formal Branching Policy.

| Phase                       | Branch Style                                                                                         | Commit Granularity                                | AI Guidance                                           |
| --------------------------- | ---------------------------------------------------------------------------------------------------- | ------------------------------------------------- | ----------------------------------------------------- |
| **Exploratory / Pre-Prod**  | Feature-level (“chunk”) branches that deliver conceptual patterns (e.g. `implement-ephemeral-state`) | Multiple small commits capturing design evolution | Favor speed and learning. Allow broader scopes.       |
| **Stabilization / QA**      | Mid-sized branches focused on one cohesive subsystem                                                 | Commits should be self-testing and revertible     | Begin enforcing semantic commits and teardown checks. |
| **Production / Maintained** | Atomic branches tied to single stories or bugfixes                                                   | One atomic, deployable change per commit          | Require CI pass + teardown proof before merge.        |

The AI must infer the correct scope from the current **project phase** (tracked in `secrets_workflow_project_tracker.md`). If the phase is labeled *Pre-Prod* or *Pattern Validation*, use feature-level branches. If *Prod* or *Maintained*, recommend smaller, atomic branches.

---

## 6. Invariants (Non-Negotiable Rules)

* No long-lived infrastructure.
* All secrets must be ephemeral or encrypted with `sops+age`.
* All changes must be versioned, peer-reviewed, and pass teardown guarantees.

---

## 7. Context Evolution Notes

* **Last Updated:** 2025-09-18
* **Known Gaps:** Advanced `age` key management, cross-cloud OIDC best practices.
* **Recent Changes:** Linked Secrets Workflow into Workflow Toolbelt roadmap.
* **Cross-Project Insights:** Ephemeral-first, OIDC, teardown validation patterns will propagate to all portfolio repos.

---

## 8. Glossary

* **Ephemeral** = Short-lived, torn down after use.
* **OIDC (OpenID Connect)** = Standard for identity federation, used here for workload identities in CI/CD.
* **sops+age** = Toolchain for encrypting/decrypting secrets in-repo.
* **MinIO** = S3-compatible object storage, used as temporary backend for Terraform in CI.
* **Teardown Guarantee** = Assertion that no resources remain after `terraform destroy`.

---

### Appendix A: Handoff Note Template

* **Goal:** \<what’s being handed off>
* **Scope:** \<boundaries, assumptions>
* **Validation:** \<tests/scans that must pass>

---

This `instructions-secrets_workflow.md` is your **Context Contract**. It ties together the README, the project tracker, and your repo’s guiding rules so that AI (and you) can always navigate and collaborate without losing alignment.
