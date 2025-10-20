# Integration & Validation Tests

This directory contains the integration and end-to-end tests for the `secrets-workflow` repository, primarily using Go and the Terratest framework.

## Purpose

The tests in this suite validate the core principles of the project by:
-   Provisioning real infrastructure using the repository's Terraform modules and examples.
-   Asserting that the provisioned resources are configured correctly.
-   Executing a full `destroy` lifecycle.
-   Verifying that no resources are left behind (the "teardown guarantee").

## Go Version

This test suite requires **Go ≥ 1.24** to support the latest testing modules and language features.