# Contributing to ECRYPTO

Thank you for your interest! Here's how to contribute:

## 🐛 Reporting Bugs

1. Check [existing issues](https://github.com/pandarudra/ecrypto/issues)
2. Create a new issue with:
   - Steps to reproduce
   - Expected vs actual behavior
   - OS version, Go version
   - Command used + error output

## ✨ Suggesting Features

Open an issue with:

- Clear description of the feature
- Use case / why it's needed
- Proposed implementation (optional)

## 🔧 Pull Requests

1. Fork the repo
2. Create a branch: `git checkout -b feature/my-feature`
3. Make changes
4. Test thoroughly:
   ```powershell
   go test ./...
   go build -o ecrypto.exe
   .\ecrypto.exe  # Test interactive mode
   ```
5. Commit: `git commit -m "feat: add amazing feature"`
6. Push: `git push origin feature/my-feature`
7. Open PR with description

## 📝 Coding Standards

- Follow Go conventions (`gofmt`, `golint`)
- Add comments for exported functions
- Write tests for new features
- Keep functions small and focused

## 🔒 Security

Report security issues privately to: rudrapanda8206@gmail.com

---

Thank you! 🎉
