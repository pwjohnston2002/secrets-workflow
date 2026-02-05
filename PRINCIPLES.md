# Project Principles

This document outlines the core principles that guide the design, implementation, and usage of the `secrets-workflow` repository. These principles ensure that all patterns and modules adhere to a consistent, secure, and efficient philosophy.

---

## 1. No Long-Lived Infrastructure

**This is the foundational principle.** We avoid creating any cloud resource that persists indefinitely. This includes, but is not limited to:
-   Virtual machines for development or testing.
-   Persistent storage for state files (unless absolutely necessary and with a defined lifecycle).
-   Static secrets or credentials.
-   Open network ports (e.g., SSH, RDP).

Every resource must have a clear, automated path to destruction.

## 2. Ephemeral-First by Default

If a resource *can* be temporary, it *must* be temporary. All CI/CD pipelines, testing environments, and developer sandboxes are designed to be created on-demand and destroyed immediately after use. This minimizes attack surfaces, reduces cost, and prevents configuration drift.

## 3. Just-in-Time (JIT) Access

Access to resources and credentials is granted only when needed and for the shortest possible duration. We achieve this through patterns like OIDC federation, which exchanges a short-lived token for temporary cloud credentials scoped to a single job. This eliminates the need for static, shared secrets like `AWS_ACCESS_KEY_ID`.

## 4. Secure by Default

Security is not an add-on; it is the default state.
-   **Least Privilege:** IAM roles and policies are scoped as narrowly as possible.
-   **Zero Trust Networking:** No inbound ports are opened by default. Access is brokered through managed services like AWS SSM.
-   **In-Repo Encryption:** All sensitive static configuration committed to the repository must be encrypted using `sops+age`.

## 5. Zero Cost When Idle

A direct benefit of ephemeral infrastructure is cost efficiency. When no workflows or development tasks are active, the associated cloud infrastructure does not exist, and therefore, incurs zero cost. This is critical for portfolio projects and cost-conscious teams.

## 6. Automate the Entire Lifecycle

The entire lifecycle of a resource—provisioning, validation, and teardown—must be automated. This ensures that our principles are applied consistently and reliably, removing the potential for human error. All ephemeral environments must have a "teardown guarantee."

