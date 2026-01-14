# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Which versions are eligible for
receiving such patches depends on the CVSS v3.0 Rating:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

The ecrypto team takes security bugs seriously. We appreciate your efforts to
responsibly disclose your findings, and will make every effort to acknowledge
your contributions.

### Where to Report

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via one of the following methods:

1. **GitHub Security Advisories** (Preferred)

   - Navigate to the Security tab in this repository
   - Click "Report a vulnerability"
   - Fill out the form with details

2. **Email**
   - Send an email to the repository maintainers
   - Use a descriptive subject line (e.g., "Security: [Brief Description]")

### What to Include

To help us triage and respond to your report quickly, please include:

- **Type of issue** (e.g., buffer overflow, SQL injection, cross-site scripting, etc.)
- **Full paths of source file(s)** related to the manifestation of the issue
- **Location of the affected source code** (tag/branch/commit or direct URL)
- **Step-by-step instructions** to reproduce the issue
- **Proof-of-concept or exploit code** (if possible)
- **Impact of the issue**, including how an attacker might exploit it
- **Your name/handle** (if you'd like to be credited)

### What to Expect

- **Acknowledgment**: We'll acknowledge receipt of your report within 48 hours
- **Initial Assessment**: We'll provide an initial assessment within 5 business days
- **Updates**: We'll keep you informed about our progress toward a fix
- **Disclosure**: We'll work with you to coordinate disclosure timing
- **Credit**: We'll credit you in the release notes (unless you prefer to remain anonymous)

## Security Best Practices

When using ecrypto, we recommend:

1. **Keep Updated**: Always use the latest version of ecrypto
2. **Secure Key Storage**: Never commit encryption keys to version control
3. **Strong Passwords**: Use strong passwords when generating keys
4. **Secure Transmission**: Use secure channels when sharing encrypted files
5. **Verify Integrity**: Verify checksums of downloaded binaries
6. **File Permissions**: Set appropriate permissions on encrypted files and keys
7. **Memory Management**: Be aware that keys may be temporarily stored in memory

## Known Security Considerations

### Encryption Algorithm

- ecrypto uses industry-standard encryption algorithms (AES-256-GCM)
- Keys are derived using Argon2id for password-based encryption

### Key Management

- Users are responsible for secure storage of encryption keys
- Lost keys cannot be recovered - encrypted data will be permanently inaccessible

### Side-Channel Attacks

- Like all software-based encryption, ecrypto may be vulnerable to side-channel attacks
- Consider the threat model for your specific use case

## Security Updates

Security updates will be released as soon as possible after a vulnerability is confirmed.
We recommend:

- Watch this repository for security advisories
- Subscribe to release notifications
- Check for updates regularly

## Scope

This security policy applies to:

- The ecrypto core application
- Official binaries and releases
- Documentation and examples

It does not apply to:

- Third-party forks or modifications
- Deprecated versions
- Issues in dependencies (please report to the respective projects)

## Bug Bounty Program

We currently do not offer a bug bounty program. However, we deeply appreciate security
researchers who help keep ecrypto and our users safe.

## Questions

If you have questions about this security policy, please open a GitHub discussion or
contact the maintainers.

---

Thank you for helping keep ecrypto and its users safe!
