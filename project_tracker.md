# project_tracker.md — Secrets Workflow

## Project Complexity Level

- [x] **Standard**: Multi-phase project with moderate complexity

---

## 1) Project Overview

- **Project Name:** Secrets Workflow
- **Repository:** `secrets-workflow`
- **Primary Objective:** Establish a secure, reusable, ephemeral-first methodology and reference implementation set for secrets handling, Terraform state handling, federated authentication, and safe cloud-access workflows.
- **Current Role:** Supporting security/infrastructure capability repo for downstream portfolio projects and workflow consumers.
- **Current Version:** v0.1-draft

### Success Criteria
- Reusable patterns exist for ephemeral auth, ephemeral state, and encrypted or runtime-generated secret handling.
- Reference modules/examples/tests demonstrate the patterns credibly.
- Repo guidance is honest, operational, and aligned with implementation.
- Downstream consumers can adopt workflows/modules through documented, versioned contracts.

### Hard Constraints
- No long-lived infrastructure as the default operating model.
- No committed plaintext secrets.
- Cloud-facing applies require explicit IAM preflight discipline.
- Security posture must favor least privilege, ephemeral credentials, and teardown guarantees.

### Out of Scope
- Operating a long-lived centralized secrets platform as the repo’s default model.
- General application delivery unrelated to secrets/state workflow patterns.
- Remediation of legacy hardcoded secrets outside this repo.

---

## 2) Context Engineering Foundation

### Context Engineering Completeness Check

- [x] **User Persona Definition:** SRE / DevSecOps transition, portfolio-oriented, cost-conscious
- [x] **Domain Context Loaded:** Ephemeral infra, CI/CD security, Terraform, OIDC, `sops/age`, managed access
- [x] **Tool/Resource Context:** GitHub Actions, AWS-first patterns, Terratest, Terraform, docs/runbooks/preflight
- [x] **Historical Context:** Repo emerged from a need to learn and demonstrate safer secrets/API key/state handling during AWS-oriented portfolio work
- [x] **Edge Case Documentation:** Local ephemeral credentials, `age` key sharing, teardown failures, cloud federation drift
- [x] **Success Pattern Recognition:** OIDC, encrypted examples, managed instance access, ephemeral state backends, reusable workflows

### User / Project Context

- **Role / Perspective:** Mid-career SRE moving toward security-focused infrastructure and cloud engineering, using the repo both as learning vehicle and portfolio artifact.
- **Primary Goal:** Create something reusable, safe, and demonstrably aligned with ephemeral-first security practice.
- **Constraint Profile:** Cost-conscious, security-sensitive, cautious about secret handling and API key management.
- **Human-AI Working Style:** Architecture and risk discussions can be collaborative; destructive or cloud-applying actions remain high-scrutiny decisions.

### Knowledge Framework

- **Known Knowns:** Terraform, CI/CD, AWS-oriented workflows, SRE/infra operational concerns
- **Known Unknowns:** Advanced `age` key management, cross-cloud OIDC portability, least-privilege hardening details, deeper observability patterns
- **Unknown Unknown Triggers:** Workflow/doc drift, hidden persistent state, incomplete teardown guarantees, consumer contract ambiguity
- **Project History:** This repo is a focused extraction of a capability needed to support broader portfolio and workflow work, not a standalone app product.

---

## 3) Execution Rules

If work is performed from this tracker:
- Tasks should be handled in listed order unless explicitly reprioritized.
- New implementation work should be added to the appropriate phase before it is treated as in scope.
- AGENTS.md governs boundaries and operational rules; this tracker governs sequence, progress, and context evolution.
- If repo reality conflicts with older planning text, note the drift explicitly rather than silently inheriting stale assumptions.

---

## 4) Repository Scope and Integration Context

### In Scope
- Ephemeral secrets and state management patterns
- OIDC-based cloud authentication examples and modules
- Managed/keyless compute access patterns such as SSM-style access
- Reusable workflows and examples for downstream adoption
- Tests and docs that prove create → validate → destroy behavior

### Out of Scope
- Persistent secrets-management platforms as core runtime dependency
- Long-lived infra normalized for convenience
- General-purpose app delivery unrelated to the capability boundary

### Relationship to Other Repos
- This repo is a foundational supporting capability.
- It is intended to publish reusable workflows, examples, and modules that other repos can consume.
- It should remain semantically stable enough for downstream reuse, which makes contract clarity more important than feature sprawl.

---

## 5) Phase Plan

### Phase 1 — Foundational Patterns
**Goal:** Establish the core capability boundary and reference implementations.

Status summary:
- [x] OIDC federation pattern work exists
- [x] `sops + age` examples/patterns exist
- [x] Local dev and secure instance-access runbooks exist
- [x] Initial principles/invariants are documented
- [x] Minimal Terraform modules exist for OIDC and SSM-related access patterns

Remaining focus:
- [ ] Ensure docs accurately match what is implemented today
- [ ] Remove or clearly mark aspirational references that are not yet present in repo reality

### Phase 1.5 — Validation & Test Credibility
**Goal:** Demonstrate that patterns are not just documented, but testable and teardown-oriented.

Status summary:
- [x] Terratest-based validation exists
- [x] Smoke/integration tests exist for core patterns
- [x] Secure instance-access validation is represented

Remaining focus:
- [ ] Tighten zero-residual-resource validation where currently implicit
- [ ] Make expected failure/canary testing more explicit where useful

### Phase 2 — Workflow Integration & Hardening
**Goal:** Make the repo safer, more reusable, and more honest as a downstream dependency.

Open tasks:
- [ ] Finalize reusable workflow contracts for downstream consumption
- [ ] Integrate or verify security scans and secret scanning in CI reality
- [ ] Improve observability/audit evidence for OIDC sessions and ephemeral state access
- [ ] Document downstream consumption patterns cleanly
- [ ] Harden IAM/test permissions toward least privilege
- [ ] Reconcile AWS-first implementation with any cloud-agnostic claims

### Phase 3 — Advanced Evolution
**Goal:** Improve maturity without violating the repo’s capability boundary.

Open tasks:
- [ ] Advance `age` key management strategy documentation
- [ ] Strengthen teardown guarantees and failure-path validation
- [ ] Explore advanced but still ephemeral access patterns where appropriate
- [ ] Evaluate cross-cloud examples only when they are concrete enough to document honestly

---

## 6) Adaptive Task Framework

### Task: Keep Operational Docs Honest
- **State:** `🔄` In Progress
- **Context Dependencies:** README, instruction file, runbooks, workflows, module/test reality
- **Goal:** Ensure repo-level docs describe implemented behavior accurately and clearly separate present-state from future-state plans
- **Success Criteria:**
  - [ ] Broken or mismatched file references corrected
  - [ ] Missing workflow references either implemented or relabeled as planned
  - [ ] Cloud-agnostic claims scoped accurately
- **Failure Modes:** Trust erosion, confusing downstream consumers, portfolio overstatement
- **AI Collaboration Notes:** AI should flag drift plainly and propose minimal corrective text

### Task: Stabilize Reusable Contract Surface
- **State:** `[ ]` To Do
- **Context Dependencies:** workflows, modules, examples, README, release/versioning expectations
- **Goal:** Define what downstream repos can rely on without ambiguity
- **Success Criteria:**
  - [ ] Workflow inputs/outputs documented accurately
  - [ ] Versioning expectations documented
  - [ ] Examples reflect real supported paths
- **Failure Modes:** Consumer breakage, undocumented assumptions, accidental interface churn
- **AI Collaboration Notes:** Prefer contract clarity over feature expansion

### Task: Harden Validation Story
- **State:** `[ ]` To Do
- **Context Dependencies:** tests, preflight docs, teardown expectations, CI design
- **Goal:** Improve confidence that ephemeral-first claims are validated, not just asserted
- **Success Criteria:**
  - [ ] Teardown guarantees are evidenced clearly
  - [ ] Scan/test expectations are documented and executable
  - [ ] Failure paths are not silently ignored
- **Failure Modes:** Security theater, stranded resources, misleading portfolio claims
- **AI Collaboration Notes:** AI should prefer evidence, not optimism

### Task: Improve Key Management Guidance
- **State:** `[ ]` To Do
- **Context Dependencies:** `sops/age` usage model, local dev runbook, team-sharing scenarios
- **Goal:** Move from basic workable guidance to a more mature and explicit key-management strategy
- **Success Criteria:**
  - [ ] Solo vs team guidance separated clearly
  - [ ] Advanced options documented without pretending they are implemented
  - [ ] Residual risks stated plainly
- **Failure Modes:** Unsafe key sharing, overcomplicated bootstrap process, false sense of maturity
- **AI Collaboration Notes:** Distinguish recommended-now from future-options

---

## 7) Risks and Mitigation

| Risk | Impact | Mitigation |
|---|---|---|
| Workflow/doc drift | High | Treat mismatches as defects; correct or relabel aspirational content |
| Incomplete teardown on failure paths | High | Strengthen destroy validation and explicit cleanup expectations |
| Long-lived bootstrap secrets creeping back in | High | Prefer OIDC and ephemeral credentials; document exceptions loudly |
| Overstated cloud-agnostic positioning | Medium | Describe current repo as AWS-reference-first until other clouds are implemented credibly |
| Consumer contract ambiguity | Medium | Version interfaces, document inputs/outputs, provide examples |
| Key management confusion for `age` | High | Separate current workable guidance from future hardening strategies |

---

## 8) Dependencies

- Terraform modules and examples in this repo
- Test harness in `tests/`
- Runbooks and preflight guides in `docs/`
- GitHub Actions / CI workflow definitions
- Cloud IAM/OIDC configuration for real-world validation
- Downstream repos that may consume versioned contracts

---

## 9) AI Collaboration Framework

### Preferred Interaction Style
- Research assistant for discovery and comparison
- Pair programmer for doc and workflow refinement
- Devil’s advocate for security and operational risk review

### Context Refresh Triggers
- At phase transitions
- When repo structure, workflow contracts, or safety invariants change
- When docs and implementation appear to diverge

### Decision Boundaries
- AI may analyze architecture, risks, documentation, contracts, tests, and migration strategy
- Human approval is required for destructive cloud actions, persistent-boundary changes, or repo-structure changes with broad downstream impact

### Assistance Preferences
- Prefer direct analysis over flattery
- Prefer explicit uncertainty over bluffing
- Prefer identifying drift, risk, and hidden assumptions early

---

## 10) Context Evolution Tracking

- **Context Version:** v1.1
- **Last Context Audit:** 2026-03-07
- **Context Debt:**
  - Some docs describe workflows or security posture more completely than the current repo implements
  - Cloud-agnostic framing is ahead of implementation
  - Key-management maturity is still partial
- **Proactive Context Gaps:**
  - Exact downstream workflow contracts need tighter definition
  - Explicit residual-resource validation can be strengthened
  - Observability patterns need more concrete implementation evidence
- **Context Compression Events:**
  - Older planning text should be summarized when superseded by repo-verified reality
- **Evolution Triggers:**
  - New release tags
  - Workflow contract changes
  - Security model changes
  - Migration of legacy instruction content into stable repo-standard docs

---

## 11) Cross-Domain Learning Integration

- **Transferable Insights:** Ephemeral-first, explicit boundaries, least-privilege defaults, and teardown evidence generalize across portfolio repos.
- **Domain-Specific Adaptations:** This repo needs stronger operational and security precision than a generic project tracker because it functions as a reusable infra/security substrate.
- **Pattern Recognition:** Context contracts tend to bloat when they mix mission, operations, planning, and history into one file.
- **Template Evolution Notes:** This repo is the migration case proving why AGENTS.md and `project_tracker.md` should have distinct jobs.
- **Methodology Transfers:** The modular context framework remains useful, but should be split into thinner operational front doors plus linked depth docs.

---

## 12) Versioning & Evolution Log

- **Template Version:** repo-specific adaptation based on the user’s generalized tracker framework
- **Project Version:** v0.1-draft
- **Major Context Updates:**
  - 2026-03-07: Migrated legacy instruction/tracker responsibilities into split `AGENTS.md` and `project_tracker.md`
- **Template Improvements Applied:**
  - Repo operations moved to AGENTS
  - Context evolution/history retained in tracker
  - Legacy context-contract overload reduced
- **Cross-Project Learnings Applied:**
  - Use AGENTS as operational front door
  - Keep sequencing and context debt in tracker
  - Link outward instead of embedding everything everywhere

---

## 13) Immediate Next Actions

- [x] Review `AGENTS.md` against actual repo paths and commands one final time
- [x] Decide whether to rename legacy `instructions-secrets_workflow.md` to archival/supporting-doc status
- [ ] Reconcile README and repo docs with current implementation reality
- [ ] Define the minimum stable downstream contract for first tagged release
- [ ] Preserve this repo as a supporting capability boundary until cross-repo comparison is complete
