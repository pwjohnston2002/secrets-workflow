# AGENTS.md — Operational Contract for secrets-workflow

This file defines how AI coding agents must operate in this repository.
It is a governance document, not a brainstorming canvas.

If instructions here conflict with direct user instructions in-session, follow the user.
If instructions here conflict with `project_tracker.md`, follow the tracker order but respect all boundaries in this file.

---

## 0) Mission Definition

**Primary Objective**  
`secrets-workflow` exists to provide reusable, ephemeral-first security and infrastructure patterns for secrets handling, Terraform state handling, federated authentication, and safe cloud-facing workflows. It is a supporting capability repo: a reference implementation and operational substrate for portfolio projects and downstream workflow consumers, not a product application.

**Definition of Done**
A task is complete only when:
- The intended behavior or documentation outcome is implemented and matches repo principles.
- Relevant validation passes (`terraform validate`, tests, teardown-oriented checks, and applicable scanners).
- Documentation is updated if behavior, workflow contracts, or operational expectations changed.
- `project_tracker.md` is updated when task status, sequence, or scope changes.

**Non-Goals**
- Building or normalizing long-lived infrastructure for secrets or state.
- Introducing a centralized persistent secrets platform as part of this repo’s default model.
- Adding speculative features that are not tied to ephemeral-first workflow patterns.

---

## 1) Companion Planning File (Optional but Enforced if Present)

If `project_tracker.md` exists:

- Work tasks strictly in listed order.
- If a change is required but not listed, add or update a tracker entry before implementing.
- Do not skip or reorder tracker phases without explicit approval.
- AGENTS.md defines rules; `project_tracker.md` defines sequence.

---

## 2) Operational Constraints

### Default Behavior
- Make the smallest viable change.
- Preserve the repo’s role as a reusable security/infra pattern library unless explicitly instructed otherwise.
- Prefer modification over replacement.
- Avoid speculative refactors.
- Keep examples, modules, workflows, tests, and docs aligned.
- Link to deeper docs instead of duplicating large explanations here.

### Assumptions Policy
- Do not invent missing requirements.
- State assumptions explicitly when repo reality is incomplete or drifted.
- Mark inferred decisions clearly.
- Treat README, runbooks, tracker, and workflow files as separate sources of truth that may drift.
- Validate workflow names, commands, and file paths against the repo before asserting them.

---

## 3) Repository Structure

```text
.
├── docs/
│   ├── policies/        # Formal contribution and branching rules
│   ├── preflight/       # Mandatory cloud deployment readiness gates
│   └── runbooks/        # Operational procedures and local-dev guidance
├── examples/            # Example Terraform configurations and usage patterns
├── modules/             # Reusable Terraform modules (OIDC, SSM access, etc.)
├── tests/               # Terratest and smoke validation
├── README.md            # Public-facing repo overview
├── PRINCIPLES.md        # Core philosophy and non-negotiable rules
└── project_tracker.md   # Ordered project execution and context evolution
```

If uncertain where something belongs, mirror existing patterns and keep infra logic, examples, tests, and docs separated.

---

## 4) Critical Commands (Must Be Used Exactly or Verified Before Substitution)

### Setup
- Terraform formatting/validation should use the repo’s existing Terraform roots and module layout.
- Terratest execution should run from `tests/` using Go tooling.
- Secrets examples must remain encrypted or illustrative only; never commit real credentials.

### Development / Validation
Use the smallest applicable command set for the files changed. Verify exact paths before claiming success.

Common commands used in this repo include:
- Terraform format: `terraform fmt -recursive`
- Terraform validation: run `terraform validate` within the relevant example/module root after init
- Go tests / Terratest smoke tests: `cd tests && go test ./...`

### Cloud Apply Safety Gate
Before recommending or performing any real cloud deployment workflow (`terraform apply`, CloudFormation deploy, CDK deploy, or equivalent), treat the IAM preflight gates in `docs/preflight/aws_iam_deployment_readiness.md` as mandatory.

Record or require results using:
- `docs/preflight/aws_iam_deployment_readiness.md`
- `docs/preflight/iam_preflight_run_template.md`

Do not recommend applying changes to AWS until that checklist is explicitly walked and marked pass/fail.

### Before Final Output
Run or account for the relevant subset of:
- `terraform fmt -recursive`
- `terraform validate` in affected Terraform roots
- `cd tests && go test ./...`
- applicable secret/IaC scanning expectations when workflow or security behavior changes

Do not substitute “similar” commands unless the repo actually uses them.

---

## 5) Boundaries

### ✅ Always Allowed
- Edit repo documentation in `README.md`, `PRINCIPLES.md`, `docs/`, and repo-local guidance files.
- Edit Terraform examples and modules when the change stays within the repo’s ephemeral-first mission.
- Add or update tests in `tests/`.
- Improve clarity, validation, or safety language in workflow-supporting docs.

### ⚠️ Ask First
- Public workflow contract changes for downstream consumers.
- Dependency additions.
- CI/CD behavior changes that alter release or security guarantees.
- Cross-module refactors that materially change repo structure.
- Changes that broaden this repo from reference capability into general application/infrastructure management.

### 🛑 Dual Control Protocol
This protocol applies to high-blast-radius actions, including:
- Destructive or irreversible Terraform operations (`apply`, `destroy`) against non-demo, non-ephemeral, or real AWS resources.
- Deleting or mutating Terraform state, encrypted secret material, or access configuration in ways that could cause loss, lockout, or security regression.
- Changing IAM, OIDC, SSM, or other access-control behavior in a way that could break access or reduce security.
- Disabling or weakening preflight gates, teardown guarantees, auditability, or other safety controls.

If such an action is required:
1. **STOP.**
2. Explain the blast radius in plain language.
3. Request explicit human approval using this exact string: `APPROVE DESTRUCTIVE ACTION`.
4. Wait for that exact string before proceeding.

### 🚫 Never Do
- Commit directly to the `main` branch.
- Commit secrets, credentials, decrypted secret material, or long-lived access keys.
- Normalize long-lived infrastructure for secrets or state.
- Weaken teardown guarantees.
- Bypass required validation or preflight checks.
- Disable tests or scanners just to force a pass.
- Present aspirational workflow files as implemented if they are not present in the repo.

## Branch Hygiene and Task Isolation

Treat each independent repo change as a separate task branch and pull request.

### Required workflow

- Start every new task from the latest `origin/main` (not just a possibly stale local `main`).
- Create a fresh branch from `origin/main` for each independent change.
- Keep each branch scoped to one recommendation, fix, or tightly related change set.
- Merge the PR to `main` before starting the next independent task.
- After merge, delete the completed task branch before beginning the next task.

### Never do

- Never continue a new task on top of a previous task branch unless explicitly instructed.
- Never mix multiple unrelated recommendations into one branch.
- Never proceed if the current branch already contains unrelated changes.
- Never open a PR whose diff includes work from earlier completed tasks.

### Agent verification steps

Before making changes for a new task, run and verify the equivalent of:

- `git fetch origin --prune`
- confirm `origin/main` exists
- confirm the task branch is based on current `origin/main`
- confirm `git diff --stat origin/main...HEAD` contains only the intended task scope

Do not rely only on local `main` if remote refs are missing or stale.
Do not create a new task branch from current `HEAD` unless the user explicitly instructs a stacked-branch workflow.

If any of these checks fail, STOP and tell the user what is wrong before proceeding.

### Pushback rule

If the user asks for a new change but the current branch/worktree appears to contain unrelated work, the agent must pause and say so explicitly.

The agent should instruct the user to:
- merge or discard the existing work
- return to current `main`
- create a fresh branch for the next task

Do not silently continue on a contaminated branch.
If `origin/main` is unavailable, stale, or not the actual base of the current task branch, STOP and tell the user.
Do not silently continue by branching from the current `HEAD`.

### Contaminated Branch Recovery

If the current task branch contains unrelated files or prior-task changes:

- Stop and do not continue editing on that branch.
- Do not present the branch as ready for merge.
- Create a fresh branch from current `origin/main`.
- Recover only the intended change by:
  - cherry-picking the intended commit, if isolated cleanly, or
  - copying only the intended file(s) into the fresh branch.
- Re-verify the diff against `origin/main` before opening or updating a PR.

### Review discipline

Prefer one PR per independent change.
Prefer small, reviewable diffs over bundled cleanup.

Before asking for PR review, verify that `git diff --stat origin/main...HEAD` includes only the files required for the current task.
If unrelated files appear, STOP and treat the branch as contaminated.

If a task expands beyond its original scope, stop and ask whether to:
- narrow the task, or
- finish the current branch and create a follow-up branch

---

## 6) Coding & Repo Standards (High Signal Only)

### Naming / Structure
- Prefer small, explicit, purpose-driven modules and examples.
- Keep reusable workflow and module contracts stable and clearly documented.
- Keep docs honest about what exists today versus future-state plans.

### Security / Error Handling
- Prefer ephemeral credentials, OIDC federation, encrypted-at-rest examples, and auditable managed access patterns.
- Treat drift between docs and implementation as a defect to call out, not to gloss over.

### Logging / Evidence
- When changing workflow or infra behavior, preserve or improve traceability, teardown evidence, and validation clarity.

For deeper standards, see:
- `PRINCIPLES.md`
- `docs/policies/branching_policy.md`
- `docs/runbooks/`
- `docs/preflight/`

AGENTS.md must not duplicate those documents.

---

## 7) Failure & Recovery Protocol

If blocked:
1. Attempt a bounded fix twice.
2. If still failing, stop.
3. Summarize:
   - What was attempted
   - What failed
   - Likely causes
   - Safe next paths forward
4. Preserve repo safety invariants while awaiting direction.

Do not enter infinite correction loops.
Do not bypass failing tests, teardown checks, or cloud safety gates.

---

## 8) Change Discipline

When making changes:
- Keep diffs minimal.
- Explain why the change is required.
- Update tests when behavior changes.
- Update documentation if interfaces, workflows, safety gates, or repo structure change.
- Update `project_tracker.md` when task ordering, scope, or completion state changes. Do not bundle tracker updates into unrelated implementation branches unless the tracker itself is part of the intended scope.

If commands, structure, or workflow expectations change:
- Update `AGENTS.md` in the same change set.

---

## 9) Instruction Budget Discipline

This file must remain concise.

- Do not embed large documentation.
- Link to deeper docs instead.
- Remove obsolete instructions promptly.
- Avoid redundant constraints.
- Keep this file operational; keep richer context, history, and evolution notes in `project_tracker.md` or supporting docs.

---

## 10) Non-Negotiable Invariants

- No long-lived infrastructure.
- All secrets must be ephemeral or encrypted with `sops + age`.
- All changes must be versioned, reviewable, and compatible with teardown guarantees.
- Managed, auditable access patterns are preferred over static credential or exposed-key workflows.

---

Last Updated: 2026-03-08 
Owner: Repository maintainer
