# Runbook: Secure Instance Access

This runbook explains the project's default pattern for accessing compute instances (e.g., for debugging or interactive sessions). It adheres to the core principle of **"no long-lived infrastructure,"** which includes eliminating persistent network vulnerabilities like open SSH ports.

## 1. The Principle: No Open SSH Ports

Traditionally, developers access virtual machines using SSH, which requires:
- Opening port 22 on the instance's security group.
- Managing and distributing static SSH keys.

Both of these create a persistent security risk. An open port is a constant attack vector, and SSH keys can be lost, stolen, or mismanaged.

Our approach is to **keep all ports closed by default** and use the cloud provider's managed services for secure, on-demand access. This is a keyless, just-in-time pattern that is more secure and auditable.

## 2. Implementation: AWS Systems Manager (SSM) Session Manager

For AWS, the primary tool for this pattern is **AWS Systems Manager (SSM) Session Manager**. It provides a secure, browser-based or CLI-based shell and port-forwarding capabilities without needing to open any inbound ports on the instance.

### How It Works

1.  The SSM Agent, installed on the EC2 instance, makes an outbound connection to the AWS SSM service endpoint.
2.  An operator with the correct IAM permissions uses the AWS CLI or Console to request a session.
3.  The SSM service authenticates the operator and brokers a secure, encrypted tunnel to the agent on the instance.

This entire process happens without any inbound traffic from the internet to the instance.

### Prerequisites

For an instance to be accessible via Session Manager, it needs:

1.  **SSM Agent Installed:** Most modern Amazon Linux, Ubuntu, and Windows AMIs include this by default.
2.  **IAM Instance Profile:** The instance must have an IAM role attached that includes the `AmazonSSMManagedInstanceCore` managed policy. This grants the agent permission to communicate with the SSM service.
3.  **Network Access:** The instance must have outbound internet access to reach the SSM endpoints (e.g., via a NAT Gateway if in a private subnet).

> **Note:** A future Terraform module (`modules/instance-ssm-profile`) will provide a reusable component to configure this IAM role easily.

### Example 1: Starting a Shell Session

To get a standard command-line shell on an instance, use the `aws ssm start-session` command.

```bash
# Replace i-0123456789abcdef0 with your instance ID
aws ssm start-session --target i-0123456789abcdef0
```

This will drop you into a `sh` shell on the instance as the `ssm-user`.

### Example 2: Port Forwarding (e.g., for a database)

Session Manager can also forward a local port to a port on a remote host accessible from your EC2 instance. This is useful for securely connecting to a database or web server running in a private subnet.

```bash
# Forward your local port 8080 to port 3306 on a private RDS instance
aws ssm start-session \
    --target i-0123456789abcdef0 \
    --document-name AWS-StartPortForwardingSessionToRemoteHost \
    --parameters '{"host":["my-private-rds-instance.random-chars.us-east-1.rds.amazonaws.com"],"portNumber":["3306"],"localPortNumber":["8080"]}'
```

You can now connect to `localhost:8080` on your machine to access the database.

## 3. Auditing and Logging

All Session Manager actions are logged in **AWS CloudTrail**, providing a full audit trail of who accessed which instance, when, and what commands were run (if session logging to S3 or CloudWatch is enabled). This provides far superior observability compared to managing shared SSH keys.

## 4. Cloud-Agnostic Equivalents

This principle is not unique to AWS. Other cloud providers offer similar managed access services:

- **Google Cloud Platform (GCP):** IAP (Identity-Aware Proxy) TCP Forwarding
- **Microsoft Azure:** Azure Bastion

The goal is always the same: leverage the cloud provider's identity and access management to broker secure connections without exposing the underlying infrastructure directly.