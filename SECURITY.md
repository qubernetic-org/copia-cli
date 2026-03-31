# Security Policy

## Scope

Copia CLI handles authentication tokens and communicates with Copia/Gitea API instances. Security issues can arise if the tool:

- Exposes or leaks authentication tokens (config file permissions, logs, error messages)
- Sends credentials to unintended hosts
- Allows command injection via user input or API responses
- Stores secrets in plaintext without appropriate file permissions
- Connects to API endpoints without TLS verification

If you discover any of the above, or any other security-relevant behavior, please report it.

## Reporting a Vulnerability

**Do not open a public GitHub issue for security vulnerabilities.**

Instead, use one of the following:

- **GitHub Security Advisories**: Report privately via [GitHub's Security Advisory feature](https://github.com/qubernetic-org/copia-cli/security/advisories/new) on this repository.
- **Email**: Send a detailed report to **[INSERT SECURITY EMAIL]**.

Please include:

- A description of the issue and its potential impact
- Steps to reproduce or a scenario demonstrating the problem
- Your OS and `copia --version` output

## Response Timeline

- **Acknowledgment**: Within 3 business days of receiving your report
- **Initial assessment**: Within 7 business days
- **Resolution or mitigation**: Targeted within 30 days, depending on severity

We will keep you informed of progress and credit reporters in the fix unless anonymity is requested.

## Supported Versions

Only the latest released version receives security updates.
