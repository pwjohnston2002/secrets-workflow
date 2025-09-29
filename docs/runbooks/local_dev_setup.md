# Runbook: Local Development Setup

This runbook provides step-by-step instructions for setting up a local development environment to work with this repository's ephemeral and secure patterns.

## 1. Tool Installation

You will need the following tools installed on your local machine:

- **[sops](https://github.com/getsops/sops/releases):** For encrypting and decrypting secret files.
- **[age](https://github.com/FiloSottile/age/releases):** The backend encryption tool used by `sops` in this project.
- **Terraform:** For infrastructure as code.
- **aws-vault:** (Recommended for AWS) For managing temporary, short-lived AWS credentials locally, mimicking the OIDC flow in CI.

## 2. `sops+age` Key Setup

This project uses `sops` with `age` to encrypt secrets stored in the repository. You must have a local `age` keypair to decrypt them.

### Step 1: Generate an `age` Keypair

If you don't already have a key, generate one. This command creates a `keys.txt` file in the standard `sops` directory.

```bash
# For macOS / Linux
mkdir -p ~/.config/sops/age
age-keygen -o ~/.config/sops/age/keys.txt
```

Your public key will be printed to the console. The private key is stored securely in the `keys.txt` file.

**IMPORTANT:** Back up your private key (the `keys.txt` file) to a secure location, like a password manager. If you lose it, you will not be able to decrypt any files encrypted with its corresponding public key.

### Step 2: Add Your Public Key to the Project

To encrypt new secrets or re-encrypt existing ones, your `age` public key must be listed in the `.sops.yaml` file at the root of this repository. If you are the first person setting up the project, add your key. If you are joining, ensure your key is added by a maintainer.

### Step 3: Test Decryption

Verify that you can decrypt an existing secret file.

```bash
# This command decrypts the file and prints its contents to standard output
sops --decrypt examples/secrets/example.tfvars.enc
```

### Step 4: Encrypting a File

Thanks to the `.sops.yaml` configuration, encrypting a file is simple. `sops` will automatically use the public keys listed in the config.

```bash
# This command encrypts a new or modified tfvars file
sops --encrypt -i examples/secrets/example.tfvars
```
