# secrets-workflow

This is a foundational repository that establishes the principles, patterns, and reusable components for managing secrets and state in an **ephemeral, cloud-agnostic, and secure** manner across all portfolio projects.

## 📜 Core Philosophy

The core principle of this repository is **"no long-lived infrastructure."** We avoid persistent cloud services for secrets or state management (like AWS Secrets Manager or dedicated S3 buckets) in favor of just-in-time, ephemeral solutions.

This approach ensures:
- **Zero Cost When Idle:** No cloud resources are left running after a CI/CD pipeline or local test completes.
- **Enhanced Security:** We rely on short-lived credentials (via OIDC) and in-repo encrypted files (`sops+age`), minimizing the attack surface.
- **Cloud-Agnosticism:** The patterns are designed to be portable across AWS, GCP, and Azure.
- **Reproducibility:** Every environment is built from code and torn down completely, guaranteeing consistency.

## ✨ What This Repo Provides

- **Reusable GitHub Actions Workflows:** For ephemeral Terraform execution (`init/plan/apply/destroy`), security scanning, and more.
- **Reference Implementations:** Code examples for OIDC federation, MinIO-in-CI for state, and `sops+age` for secrets.
- **Documentation & Runbooks:** Clear guides on principles, local development setup, and architectural decisions.
- **Terraform Modules:** For common ephemeral patterns.

## ⚠️ Disclaimer

The principles, patterns, and reference implementations in this repository are for educational and portfolio purposes. They are derived from publicly available documentation, open-source tools (like `sops`, `age`, `MinIO`, and `Terraform`), and represent standard, modern industry practices for cloud security and ephemeral infrastructure.

This project is a personal initiative developed on personal time using personal resources. The content is **non-proprietary** and does not reflect or include any confidential information or intellectual property from any employer, past or present. The project is governed by the MIT License, which permits broad use but provides the software "AS IS" without warranty.
