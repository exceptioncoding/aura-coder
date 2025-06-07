# 🤝 Contributing to Aura

Thank you for your interest in contributing to Aura! This document provides guidelines and information for contributors.

## 🚀 Quick Start

1. **Fork the repository** on GitHub
2. **Clone your fork** locally
3. **Create a feature branch** from `develop`
4. **Make your changes** following our guidelines
5. **Run tests** to ensure everything works
6. **Submit a pull request** to the `develop` branch

## 📋 Development Setup

### Prerequisites
- Go 1.20 or later
- Git
- Make (optional, for using Makefile commands)

### Local Development
```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/aura.git
cd aura

# Install dependencies
go mod download

# Run tests
go run test_runner.go

# Run linting
golangci-lint run

# Build the application
go build -o aura .
```

## 🌳 Branching Strategy

We use **Git Flow** branching model:

- **`main`**: Production-ready code, protected branch
- **`develop`**: Integration branch for features, protected branch
- **`feature/*`**: New features (branch from `develop`)
- **`bugfix/*`**: Bug fixes (branch from `develop`)
- **`hotfix/*`**: Critical fixes (branch from `main`)
- **`release/*`**: Release preparation (branch from `develop`)

### Branch Naming Convention
```
feature/add-new-llm-provider
bugfix/fix-config-loading
hotfix/security-vulnerability
release/v1.2.0
```

## 🔄 Pull Request Process

### 1. Before Creating a PR
- [ ] Ensure your branch is up-to-date with `develop`
- [ ] Run all tests locally: `go run test_runner.go`
- [ ] Run linting: `golangci-lint run`
- [ ] Update documentation if needed
- [ ] Add/update tests for new functionality

### 2. Creating the PR
- [ ] Use the provided PR template
- [ ] Write a clear, descriptive title
- [ ] Provide detailed description of changes
- [ ] Link related issues
- [ ] Mark the PR as draft if work is in progress

### 3. PR Requirements
- [ ] **All CI checks must pass** ✅
- [ ] **At least one approving review** from a maintainer
- [ ] **No merge conflicts** with target branch
- [ ] **Test coverage maintained** or improved
- [ ] **Documentation updated** if applicable

### 4. After PR Approval
- [ ] Squash commits if requested
- [ ] Ensure CI is still passing
- [ ] Merge will be handled by maintainers

## 🧪 Testing Guidelines

### Test Categories
1. **Unit Tests**: Test individual functions/methods
2. **Integration Tests**: Test component interactions
3. **E2E Tests**: Test complete user workflows
4. **Benchmark Tests**: Performance testing

### Writing Tests
```go
func TestFunctionName(t *testing.T) {
    // Arrange
    input := setupTestData()
    
    // Act
    result, err := functionUnderTest(input)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Test Requirements
- [ ] All new code must have tests
- [ ] Test coverage should not decrease
- [ ] Tests must be deterministic and fast
- [ ] Use table-driven tests for multiple scenarios
- [ ] Mock external dependencies

## 📝 Code Style Guidelines

### Go Code Standards
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Use `goimports` for import organization
- Follow Go naming conventions
- Write clear, self-documenting code

### Code Organization
```
aura/
├── cmd/           # CLI commands
├── config/        # Configuration management
├── llm/           # LLM client implementations
├── models/        # Data models and types
├── ui/            # User interface components
├── app/           # Application logic
├── e2e/           # End-to-end tests
└── docs/          # Documentation
```

### Comments and Documentation
- Public functions/types must have comments
- Complex logic should be explained
- Use examples in documentation
- Keep comments up-to-date with code changes

## 🔒 Security Guidelines

### Security Best Practices
- [ ] Never commit API keys or secrets
- [ ] Validate all user inputs
- [ ] Use secure defaults
- [ ] Follow principle of least privilege
- [ ] Sanitize file paths and user data

### Reporting Security Issues
- **DO NOT** create public issues for security vulnerabilities
- Email security issues to: [security@aura-project.com]
- Include detailed reproduction steps
- Allow time for fix before public disclosure

## 🐛 Bug Reports

### Before Reporting
- [ ] Search existing issues
- [ ] Try latest version
- [ ] Check documentation
- [ ] Reproduce with minimal example

### Bug Report Template
```markdown
**Bug Description**
Clear description of the bug

**Steps to Reproduce**
1. Step one
2. Step two
3. Step three

**Expected Behavior**
What should happen

**Actual Behavior**
What actually happens

**Environment**
- OS: [e.g., Windows 11, macOS 13, Ubuntu 22.04]
- Go Version: [e.g., 1.21.0]
- Aura Version: [e.g., v1.2.0]

**Additional Context**
Any other relevant information
```

## ✨ Feature Requests

### Before Requesting
- [ ] Check if feature already exists
- [ ] Search existing feature requests
- [ ] Consider if it fits project scope
- [ ] Think about implementation complexity

### Feature Request Template
```markdown
**Feature Description**
Clear description of the proposed feature

**Use Case**
Why is this feature needed?

**Proposed Solution**
How should this feature work?

**Alternatives Considered**
Other approaches you've considered

**Additional Context**
Any other relevant information
```

## 📚 Documentation

### Documentation Types
- **README**: Project overview and quick start
- **API Documentation**: Code documentation
- **User Guide**: Detailed usage instructions
- **Developer Guide**: Contributing and development

### Documentation Standards
- Use clear, concise language
- Include code examples
- Keep documentation up-to-date
- Use proper markdown formatting
- Include screenshots for UI features

## 🏷️ Release Process

### Version Numbering
We use [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backwards compatible)
- **PATCH**: Bug fixes (backwards compatible)

### Release Workflow
1. Create release branch from `develop`
2. Update version numbers and changelog
3. Run full test suite
4. Create PR to `main`
5. After merge, tag release
6. Automated deployment via GitHub Actions

## 🎯 Code Review Guidelines

### For Authors
- [ ] Keep PRs focused and small
- [ ] Write clear commit messages
- [ ] Respond to feedback promptly
- [ ] Test your changes thoroughly
- [ ] Update documentation

### For Reviewers
- [ ] Review code logic and design
- [ ] Check for security issues
- [ ] Verify test coverage
- [ ] Ensure documentation is updated
- [ ] Be constructive and respectful

### Review Checklist
- [ ] Code follows project standards
- [ ] Tests are comprehensive
- [ ] Documentation is updated
- [ ] No security vulnerabilities
- [ ] Performance impact considered

## 🤝 Community Guidelines

### Code of Conduct
- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Assume good intentions
- Follow GitHub's community guidelines

### Communication Channels
- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and ideas
- **Pull Requests**: Code contributions
- **Email**: Security issues and private matters

## 🏆 Recognition

### Contributors
All contributors will be recognized in:
- README contributors section
- Release notes
- Project documentation

### Contribution Types
We value all types of contributions:
- Code contributions
- Bug reports
- Documentation improvements
- Feature suggestions
- Testing and feedback
- Community support

## 📞 Getting Help

### Resources
- **Documentation**: Check project docs first
- **GitHub Issues**: Search existing issues
- **GitHub Discussions**: Ask questions
- **Code Examples**: Look at existing code

### Contact
- **General Questions**: GitHub Discussions
- **Bug Reports**: GitHub Issues
- **Security Issues**: [security@aura-project.com]
- **Maintainers**: @maintainer1 @maintainer2

---

Thank you for contributing to Aura! Your contributions help make this project better for everyone. 🚀
