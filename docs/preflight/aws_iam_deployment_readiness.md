---

# AWS IAM Deployment Readiness

**Roles, Privileges, and Policy Preflight**

## Invocation Template (for AI + Human)

When preparing to deploy, provide:

1) Cloud/account context (account ID, region, environment)
2) Deployer principal ARN (CI role) + runtime principal ARNs (if any)
3) Deployment tool (Terraform/CDK/CloudFormation) + what will be created
4) Any KMS keys, buckets, repos, or secrets involved
5) Whether SCPs/boundaries/SSO are in play

Then run Gates 0–11 and record PASS/FAIL with notes.
Stop on first FAIL unless explicitly doing diagnostic work.

---

## 1. Meta — Why This Document Exists

Cloud deployments fail in two fundamentally different ways:

1. **The infrastructure is wrong**
2. **The permissions are wrong**

Historically, teams blur these failure modes. When a deployment breaks, engineers change Terraform, broaden IAM policies, or “just try admin,” without knowing which layer is actually broken.

This document exists to prevent that ambiguity.

It defines a **minimum, repeatable, pre-deployment contract** for AWS IAM roles, privileges, and policies that must be satisfied before any infrastructure change is executed.

It does *not* attempt to describe every possible IAM scenario.

Instead, it establishes:

* A shared mental model
* A fixed sequence of validation gates
* A common vocabulary for debugging
* A safety floor for least privilege

Think of this as **training wheels for authorization correctness**.

If these gates pass, failures are almost certainly **code or design problems**, not IAM problems.

If these gates fail, deployment must stop.

---

## 2. How To Use This Document

This document is organized into three conceptual layers:

1. **Operational Gates** (what must be checked)
2. **Pass Criteria** (what “ready” means)
3. **Foundational Context** (why these rules exist)

Only Layer 1 blocks deployments.
Layer 3 exists to educate and calibrate judgment.

---

# Layer 1 — AWS IAM Deployment Readiness Checklist

## Gate 0 — Identity & Execution Context

**Objective:** Confirm who is deploying and under what authority.

* [ ] `aws sts get-caller-identity` matches expected account, role, and ARN
* [ ] Deployment identity is a **role**, not a long-lived IAM user
* [ ] No root credentials in use
* [ ] MFA enforced for any human-assumed role
* [ ] Active permissions boundary or SCP (if used) is documented
* [ ] Identify principals: **deployer role** (CI/Terraform/CDK) vs **runtime role(s)** (EC2 instance profile / ECS task role / Lambda exec role)
* [ ] Confirm STS session duration is sufficient for worst-case deploy + destroy (and CI refresh/reauth strategy exists)
* [ ] If using IAM Identity Center (SSO): identify the permission set + account assignment and validate session duration and boundary behavior

**Why:** If you misidentify the principal, every later result is meaningless.

---

## Gate 1 — Role Trust Policy Integrity

**Objective:** Ensure roles can only be assumed by intended principals.

For each role involved:

* [ ] Trust policy specifies **exact principals** (service, account, or role ARN)
* [ ] No `"Principal": "*"`
* [ ] Cross-account trust explicitly documented
* [ ] ExternalId required for third-party access (if applicable)
* [ ] No unused principals remain

**Red Flags**

* Wildcard principals
* One role shared by unrelated services
* Trust policy and permission policy owned by different teams

**If using OIDC / federation (GitHub Actions, EKS IRSA, external IdPs):**

* [ ] Trust policy conditions include tight `sub` and `aud` constraints
* [ ] No broad subject patterns (org/repo wildcards) unless explicitly intended
* [ ] Token issuer URL matches the configured OIDC provider
* [ ] For EKS IRSA: service account namespace/name constraints are correct

---

## Gate 2 — Permission Boundary (Maximum Power)

**Objective:** Define the absolute ceiling of what the role can ever do.

* [ ] Permissions boundary exists **or** SCP provides equivalent constraint
* [ ] Boundary blocks:

  * `iam:*` except explicitly required
  * `sts:AssumeRole` except approved targets
* [ ] Boundary denies known privilege-escalation paths (see Gate 7)

**Verification (not just documentation):**

* [ ] Confirm whether AWS Organizations SCPs apply to this account/OU and link the applicable SCP(s)
* [ ] If using permissions boundaries, verify the boundary policy allows the required actions/resources for this deployment
* [ ] If assuming roles via STS, verify no restrictive **session policy** is attached (explicitly note whether you use one)


**Why**

Identity policy = what it *may* do
Boundary/SCP = what it *can never* do

---

## Gate 3 — Least-Privilege Policy Design

**Objective:** Ensure permissions are minimal and intentional.

For each policy statement:

* [ ] Actions limited to required verbs
* [ ] Resources not `"*"` unless AWS requires it
* [ ] Separate statements per resource type
* [ ] Conditions used where possible:

  * `aws:RequestedRegion`
  * `aws:RequestTag`
  * `aws:ResourceTag`
  * `aws:PrincipalArn`

**Design Rules**

* Many small statements > one giant statement
* Customer-managed > inline
* No AWS managed “FullAccess” policies

---

## Gate 3.5 — Explicit Deny Review (Hard Stops)

**Objective:** Identify explicit `Deny` statements that will override allows and cause “it’s allowed but still denied” failures.

* [ ] Scan identity policies, permissions boundaries, SCPs, and resource policies for explicit `Deny` affecting the deployment path
* [ ] For each deny: list the denied action(s), resource scope, and the intended rationale/owner


---

## Gate 4 — Authorization Preflight (Policy Simulation)

**Objective:** Prove permissions before running Terraform.

```bash
aws iam simulate-principal-policy \
  --policy-source-arn <ROLE_ARN> \
  --action-names <ACTIONS> \
  --resource-arns <RESOURCES>
```

* [ ] All evaluated actions return `allowed`
* [ ] ARNs quoted
* [ ] Resource type matches action
* [ ] Simulator results are **necessary but not sufficient** (resource policies/SCPs/boundaries can still deny)
* [ ] Validate against CloudTrail `AccessDenied` events during a dry-run deploy (treat CloudTrail as the ground-truth denial source)
* [ ] Confirm simulation inputs include correct resource ARNs and relevant condition context keys (tags, region, principal ARN, etc.)


**Minimum Critical Paths**

* IAM: create role, instance profile, attach, passrole
* Compute
* Networking
* Storage / logging / secrets

**Resource-based policy preflight (common “IAM says yes, AWS says no” cases):**

* [ ] KMS: confirm the principal is permitted by the **KMS key policy** (and grants if used)
* [ ] S3: confirm bucket policy + ACLs (if any) allow the principal
* [ ] ECR: confirm repository policy allows pull/push (cross-account especially)
* [ ] Lambda: confirm Lambda resource policy permits invoke (service/account principals)
* [ ] Secrets Manager / SSM Parameter Store (if cross-account): confirm resource policies allow access


---

## Gate 5 — Functional vs Authorization Separation

**Objective:** Never mix code correctness with permission correctness.

* [ ] Code verified once with admin-equivalent role
* [ ] Same code verified again with least-privilege role
* [ ] Differences captured as IAM deltas

**Rule**

> Never fix IAM by editing Terraform.
> Never fix Terraform by broadening IAM.

**Break-glass + rollback guardrails:**

* [ ] Break-glass role exists, is audited, and requires MFA + approvals
* [ ] Rollback permissions (destroy/undo) are verified for the least-privilege deployer role
* [ ] Any emergency broadening is time-boxed and tracked (ticket/PR), then removed

---

## Gate 6 — PassRole & Role Chaining Control

* [ ] `iam:PassRole` scoped to exact ARNs
* [ ] `iam:PassedToService` constrained
* [ ] No wildcard passrole
* [ ] Chained assume-role paths documented

---

## Gate 7 — Privilege Escalation Scan

Search for (and justify) any of the following:

* `iam:PassRole` (especially with broad resources)
* `iam:CreatePolicyVersion`
* `iam:SetDefaultPolicyVersion`
* `iam:AttachRolePolicy`
* `iam:PutRolePolicy`
* `iam:UpdateAssumeRolePolicy`
* `sts:AssumeRole` to high-privilege roles (or with wildcard resources)
* `iam:CreateAccessKey`
* `iam:UpdateLoginProfile`
* `iam:PutUserPolicy` / `iam:AttachUserPolicy` (if users exist)
* `iam:TagRole` / `iam:UntagRole` (tag-based controls can be bypassed if tagging is too open)

* [ ] Each occurrence is explicitly justified and tightly scoped
* [ ] None exist accidentally


---

## Gate 8 — Resource Lifecycle Coverage

For each resource type:

* [ ] Create
* [ ] Read/List
* [ ] Modify
* [ ] Delete

Destroy must succeed.

---

## Gate 9 — Tagging Enforcement

* [ ] Required tags defined
* [ ] RequestTag enforced
* [ ] ResourceTag enforced

---

## Gate 10 — Observability & Auditability

* [ ] CloudTrail enabled
* [ ] AccessDenied monitored
* [ ] IAM Access Analyzer enabled
* [ ] Policy changes logged

**KMS sanity (encryption is its own IAM universe):**

* [ ] Identify all KMS keys involved (S3, EBS, RDS, Secrets, SSM params, logs, etc.)
* [ ] Verify key policy grants the **deployer role** required actions (plan/apply/destroy)
* [ ] Verify key policy grants the **runtime role(s)** required actions (decrypt/data key usage)
* [ ] Confirm key is in the correct region and not pending deletion


---

## Gate 11 — Drift & Hygiene

* [ ] Unused roles removed
* [ ] Unused policies removed
* [ ] Last-used reviewed
* [ ] Inline policies eliminated

---

# Layer 2 — Deployment Pass Criteria

Deployment may proceed only if:

* All gates pass
* Authorization preflight succeeds
* Least-privilege role deploys **and destroys**
* No wildcard admin policies exist

---

# Layer 3 — Foundational Context (Conceptual Primer)

This section provides conceptual grounding for engineers who want to understand *why* the gates exist. It does not replace the checklist.

---

## IAM Roles vs Policies

* Role = identity you can become
* Policy = permissions attached to that identity
* Trust policy = who may assume
* Permission policy = what may be done

---

## Principle of Least Privilege

Grant only the permissions required to perform a task and nothing more.

In AWS this means:

* Precise actions
* Scoped resources
* Conditional access
* Continuous reduction

Least privilege is not a static state. It is an iterative process.

---

## Identity-Based vs Resource-Based Policies

* Identity-based → attached to roles
* Resource-based → attached to resources (S3, KMS, etc.)

Both participate in authorization evaluation.

---

## Common Anti-Patterns

* Wildcard actions
* Wildcard resources
* Hardcoded credentials
* Broad PassRole
* Shared mega-roles
* Root usage

---

## Troubleshooting Philosophy

* Logs over guesses
* Simulator over intuition
* First denial only
* One change per iteration

---

### Closing Thought

Authorization is code.

Treat IAM with the same discipline as application logic:

* Reviewed
* Tested
* Versioned
* Minimal

When IAM becomes boring, your cloud becomes safe.

---

#below is archived for context to inform the above.

<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

# you're an expert in Amazon web services particularly for cloud infrastructure operations and infrastructure and I have a question about roles and policies in AWS

AWS IAM roles and policies are key for secure cloud infrastructure operations, enabling granular access control without long-term credentials. Roles act as assumable identities, while policies define permissions attached to them or other entities.[^1_1]

## Core Concepts

IAM roles are AWS identities with specific permissions that anyone who needs them can assume, providing temporary security credentials instead of permanent keys like passwords or access keys. They support delegation for users, apps, services, or cross-account access, such as granting an EC2 instance access to S3 without embedding keys.[^1_2][^1_1]
Policies are JSON documents outlining allowed or denied actions on resources; they follow least privilege by specifying effects like "Allow" or "Deny," actions, resources, and optional conditions.[^1_3][^1_4]

## Types of Policies

Identity-based policies attach to users, groups, or roles and grant permissions to those identities; they include AWS managed, customer managed, or inline types.[^1_3]
Resource-based policies attach directly to resources like S3 buckets or IAM role trust policies, specifying principals (even cross-account) that can access them.[^1_2][^1_3]

## Roles in Practice

Every IAM role requires a trust policy (resource-based) defining who or what can assume it, such as an AWS service or another account's role, plus permissions policies (identity-based) for actual actions.[^1_5][^1_1]
In infrastructure ops, use service roles for EC2 or Lambda to access resources securely, avoiding user keys; service-linked roles are auto-created and managed by services.[^1_1]

## Best Practices

Grant least privilege: Start minimal and use IAM Access Analyzer to generate policies from access activity or validate existing ones. Prefer roles over users for workloads, rotate credentials via assumption, and review unused permissions regularly.[^1_6][^1_3]
<span style="display:none">[^1_10][^1_11][^1_12][^1_13][^1_14][^1_15][^1_7][^1_8][^1_9]</span>

<div align="center">⁂</div>

[^1_1]: https://www.apono.io/wiki/whats-the-difference-between-aws-iam-roles-and-policies/

[^1_2]: https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies.html

[^1_3]: https://www.strongdm.com/blog/aws-iam-roles-vs-policies

[^1_4]: https://delinea.com/blog/aws-iam-roles-and-policies-explained

[^1_5]: https://www.wiz.io/academy/cloud-security/aws-iam-roles

[^1_6]: https://docs.aws.amazon.com/IAM/latest/UserGuide/best-practices.html

[^1_7]: https://www.reddit.com/r/AWSCertifications/comments/15ohcfh/can_someone_make_clear_iam_policies_vs_iam_roles/

[^1_8]: https://www.pluralsight.com/resources/blog/cloud/comparing-aws-azure-and-google-cloud-iam-services

[^1_9]: https://www.youtube.com/watch?v=BSodkwWB-8s

[^1_10]: https://www.datadoghq.com/blog/iam-least-privilege/

[^1_11]: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles.html

[^1_12]: https://aws.amazon.com/blogs/security/optimize-aws-administration-with-iam-paths/

[^1_13]: https://stackgen.com/blog/what-you-need-to-know-about-aws-vs.-azure-iam-policies

[^1_14]: https://www.linkedin.com/pulse/key-differences-between-role-policy-iam-avinash-tietler-hrzdc

[^1_15]: https://spacelift.io/blog/aws-iam-best-practices


---

# explain the concept of least privilege as a general concept first and then as a concept applicable to deploying into AWS environments please

The principle of least privilege restricts access rights to the minimum levels necessary for users, processes, or systems to perform their tasks, minimizing potential damage from errors, breaches, or misuse.[^2_1]

## General Concept

This cybersecurity principle ensures entities only access specific data, resources, or functions required for their role, reducing the attack surface against threats like malware or unauthorized lateral movement.[^2_2][^2_4]
It applies universally—e.g., a backup process gets read-only file access but no software installation rights—and supports Zero Trust by assuming no inherent trust.[^2_3][^2_11]

## AWS Deployment Application

In AWS, implement least privilege via IAM by creating granular policies with precise actions, resources, and conditions, avoiding broad permissions like "AdministratorAccess."[^2_12]
For deployments, use IAM roles for EC2 instances or Lambda functions instead of embedding access keys; attach minimal policies, then refine with IAM Access Analyzer based on access logs.[^2_13][^2_14]
Regular audits, policy simulation, and service control policies (SCPs) in AWS Organizations enforce this at scale, preventing over-privileged accounts in infrastructure operations.[^2_15]
<span style="display:none">[^2_10][^2_5][^2_6][^2_7][^2_8][^2_9]</span>

<div align="center">⁂</div>

[^2_1]: https://csrc.nist.gov/glossary/term/least_privilege

[^2_2]: https://www.paloaltonetworks.com/cyberpedia/what-is-the-principle-of-least-privilege

[^2_3]: https://www.crowdstrike.com/en-us/cybersecurity-101/identity-protection/principle-of-least-privilege-polp/

[^2_4]: https://www.cyberark.com/what-is/least-privilege/

[^2_5]: https://www.cisecurity.org/insights/spotlight/ei-isac-cybersecurity-spotlight-principle-of-least-privilege

[^2_6]: https://identitymanagementinstitute.org/the-principle-of-least-privilege/

[^2_7]: https://www.ncsc.govt.nz/protect-your-organisation/principle-of-least-privilege/

[^2_8]: https://zeronetworks.com/blog/a-practical-guide-to-least-privilege-access-zero-trust-security-in-action

[^2_9]: https://en.wikipedia.org/wiki/Principle_of_least_privilege

[^2_10]: https://learn.microsoft.com/en-us/entra/identity-platform/secure-least-privileged-access

[^2_11]: https://www.strongdm.com/blog/aws-iam-roles-vs-policies

[^2_12]: https://docs.aws.amazon.com/IAM/latest/UserGuide/best-practices.html

[^2_13]: https://www.apono.io/wiki/whats-the-difference-between-aws-iam-roles-and-policies/

[^2_14]: https://www.datadoghq.com/blog/iam-least-privilege/

[^2_15]: https://spacelift.io/blog/aws-iam-best-practices


---

# using the concept of roles and privileges within AWS IAM can you please now enumerate what would be considered minimum best practices and then extend that to recommended best practices for things that should be checked every time someone deploys anything into AWS in regards to IAM roles and privileges please ?

Minimum best practices for AWS IAM roles and privileges focus on least privilege and temporary credentials to secure deployments. Extending to recommended checks ensures ongoing compliance during every deployment.[^3_1]

## Minimum Best Practices

- Use IAM roles over users for workloads and services, providing temporary credentials via AssumeRole to avoid long-term keys.[^3_2][^3_1]
- Apply least-privilege permissions: Grant only required actions, resources, and conditions in policies; start minimal and add as needed.[^3_3][^3_4]
- Define trust policies precisely on roles to limit who/what can assume them, such as specific services or accounts.[^3_5]


## Recommended Deployment Checks

- Verify no root user usage; enforce MFA on all IAM users and check for unused credentials to deactivate.[^3_1][^3_2]
- Run IAM policy simulator or Access Analyzer to test and generate least-privilege policies from access activity logs.[^3_6][^3_1]
- Scan for unused roles/policies/users, external/cross-account access, and apply permissions boundaries or SCPs for guardrails.[^3_7][^3_1]
- Confirm conditions in policies (e.g., IP restrictions, MFA) and avoid inline policies; prefer managed or customer-managed ones.[^3_2][^3_7]
<span style="display:none">[^3_10][^3_11][^3_12][^3_13][^3_14][^3_15][^3_16][^3_8][^3_9]</span>

<div align="center">⁂</div>

[^3_1]: https://docs.aws.amazon.com/IAM/latest/UserGuide/best-practices.html

[^3_2]: https://spacelift.io/blog/aws-iam-best-practices

[^3_3]: https://aws.amazon.com/iam/resources/best-practices/

[^3_4]: https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies.html

[^3_5]: https://www.apono.io/wiki/whats-the-difference-between-aws-iam-roles-and-policies/

[^3_6]: https://www.datadoghq.com/blog/iam-least-privilege/

[^3_7]: https://www.wiz.io/academy/cloud-security/aws-iam-best-practices

[^3_8]: https://www.qovery.com/blog/aws-production-deployment-checklist

[^3_9]: https://docs.aws.amazon.com/prescriptive-guidance/latest/terraform-aws-provider-best-practices/security.html

[^3_10]: https://documentation.commvault.com/v11/software/best_practices_for_iam_roles.html

[^3_11]: https://www.reddit.com/r/aws/comments/17z37nv/best_practices_checklist/

[^3_12]: https://www.sysdig.com/blog/26-aws-security-best-practices

[^3_13]: https://docs.aws.amazon.com/cdk/v2/guide/best-practices.html

[^3_14]: https://allcloud.io/blog/best-practice-for-ultra-secure-deployment-on-amazon-cloud/

[^3_15]: https://www.reddit.com/r/aws/comments/gm9cbf/best_practices_for_iam_policies/

[^3_16]: https://www.sentinelone.com/cybersecurity-101/cloud-security/aws-security-checklist/


---

# what are some known common anti patterns when deploying into AWS in regards to IAM roles and privileges ?

Common anti-patterns in AWS IAM roles and privileges during deployments often lead to over-privileged access, privilege escalation risks, and compliance issues.[^4_1][^4_2]

## Overly Broad Permissions

Granting wildcard actions like "*" or full service access (e.g., AmazonS3FullAccess) instead of specific resources and conditions violates least privilege.[^4_3][^4_1]
Using managed policies without customization expands blast radius unnecessarily, as they can change or grant unintended permissions.[^4_1]

## Misconfigured Trust Policies

Wildcard principals ("Principal": {"AWS": "*"}) allow any account to assume roles, enabling cross-account takeovers.[^4_4]
Listing multiple principals treated as OR logic instead of intended AND, letting unintended entities assume roles.[^4_4]

## Credential and User Mistakes

Hardcoding IAM access keys in code, configs, or apps exposes them to leaks; always use roles with AssumeRole for temporary creds.[^4_5]
Creating isolated users with direct broad permissions instead of groups/roles; forgetting to revoke post-task access.[^4_5]

## Other Pitfalls

Using root account for routine ops or skipping MFA on admins/users, risking full account compromise.[^4_5]
Overlooking iam:PassRole with "*" resources, allowing privilege escalation by passing roles to unauthorized services.[^4_6]
Neglecting regular audits for unused entities or privilege escalations via tools like IAM Access Analyzer.[^4_2][^4_1]
<span style="display:none">[^4_10][^4_11][^4_12][^4_13][^4_14][^4_15][^4_16][^4_7][^4_8][^4_9]</span>

<div align="center">⁂</div>

[^4_1]: https://summitroute.com/blog/2018/07/02/aws_managed_policies_are_an_anti_pattern/

[^4_2]: https://docs.aws.amazon.com/IAM/latest/UserGuide/best-practices.html

[^4_3]: https://dev.solita.fi/2020/12/22/aws-antipatterns.html

[^4_4]: https://www.token.security/blog/iam-role-trust-policies-misconfigurations-hiding-in-plain-sight

[^4_5]: https://www.growexx.com/blog/top-10-security-pitfalls-to-avoid-in-aws-deployment/

[^4_6]: https://aws.amazon.com/blogs/security/use-iam-access-analyzer-policy-generation-to-grant-fine-grained-permissions-for-your-aws-cloudformation-service-roles/

[^4_7]: https://aws.amazon.com/blogs/security/high-availability-iam-design-patterns/

[^4_8]: https://aws.amazon.com/blogs/compute/operating-lambda-anti-patterns-in-event-driven-architectures-part-3/

[^4_9]: https://aws.amazon.com/blogs/devops/best-practices-for-deploying-aws-devops-agent-in-production/

[^4_10]: https://speakerdeck.com/htan/anti-patterns-from-security-best-practices-in-aws-identity-and-access-management-aws-iam

[^4_11]: https://www.iansresearch.com/resources/all-blogs/post/security-blog/2021/12/17/centralized-iam-best-practices

[^4_12]: https://www.linkedin.com/posts/riyazsayyad_this-statement-shouldnt-exist-for-cloud-activity-7412332105562693633-qTjy

[^4_13]: https://www.concurrencylabs.com/blog/mistakes-application-growth-aws/

[^4_14]: https://docs.aws.amazon.com/cdk/v2/guide/best-practices.html

[^4_15]: https://sonraisecurity.com/access-2024/iam-in-the-cloud-successes-pitfalls-and-lessons-learned/

[^4_16]: https://allthingsincloud.com/aws-iam-permission-boundaries-explained-real-life-mistakes-and-practical-fixes-abfeb6a6d495


---

# speaking is an AWS expert in IAM roles and policies what are some recommendations you would give for deployment troubleshooting ?

Deployment troubleshooting for AWS IAM roles and policies requires systematic isolation of denies using logs, simulation, and minimal permissions iteration.[^6_1][^6_2]

## Enable Detailed Logging

Activate CloudTrail for IAM and relevant services (e.g., EC2, SSM); filter Event history by username or runid to pinpoint AccessDenied events with exact actions/resources.[^6_2][^6_1]
Use Terraform/CLI verbose flags (--debug) and AWS SDK logging to capture RequestID-tied errors for correlation.[^6_1]

## Simulate Permissions

Employ IAM Policy Simulator console/CLI for principal-policy tests on denied actions; test single then multi-actions, noting EvalDecision and MatchedStatements—add temporary iam:SimulatePrincipalPolicy if needed, then revoke.[^6_3][^6_1]
Leverage IAM Access Analyzer to generate/review policies from CloudTrail access activity, validating against logs.[^6_2]

## Iterate Least-Privilege Deltas

Start with "*" resources for Create actions (e.g., iam:CreateRole), refine post-success; apply policy changes, re-run tests, diff logs for progressed denies.[^6_1]
Tag resources with runid for console tracing (VPC/EC2/IAM); confirm via sts:GetCallerIdentity no boundaries/SCPs block.[^6_1]

## Common Fixes

Scope trust policies to exact principals/AURNs, avoiding "*"; handle Create quirks with broad allows initially.[^6_4][^6_5]
Audit unused entities quarterly; use permissions boundaries to cap role escalations during deploys.[^6_2]
<span style="display:none">[^6_6]</span>

<div align="center">⁂</div>

[^6_1]: transcript.txt

[^6_2]: https://docs.aws.amazon.com/IAM/latest/UserGuide/best-practices.html

[^6_3]: https://aws.amazon.com/blogs/security/use-iam-access-analyzer-policy-generation-to-grant-fine-grained-permissions-for-your-aws-cloudformation-service-roles/

[^6_4]: https://www.token.security/blog/iam-role-trust-policies-misconfigurations-hiding-in-plain-sight

[^6_5]: https://www.apono.io/wiki/whats-the-difference-between-aws-iam-roles-and-policies/

[^6_6]: iam_roles_and_priv_primer.md
