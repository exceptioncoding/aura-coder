# 🧪 Aura Testing Documentation

This document provides comprehensive information about the testing strategy, implementation, and execution for the Aura AI coding assistant.

## 📊 Testing Strategy Overview

Our testing strategy follows the **Testing Pyramid** approach, ensuring comprehensive coverage at all levels:

```
           /\
          /  \
         / E2E \    ← 10% End-to-End Tests
        /______\
       /        \
      /Integration\ ← 20% Integration Tests  
     /__________\
    /            \
   /  Unit Tests  \ ← 70% Unit Tests
  /________________\
```

## 🏗️ Test Architecture

### 1. **Unit Tests** (Foundation Layer)
**Location**: `*_test.go` files alongside source code  
**Purpose**: Test individual components in isolation  
**Coverage**: 70% of total tests

#### Config Package Tests (`config/config_test.go`)
- ✅ Configuration loading and saving
- ✅ YAML marshaling/unmarshaling
- ✅ Default configuration creation
- ✅ Path resolution and home directory handling
- ✅ Validation logic for providers and API keys
- ✅ Trusted directory management
- ✅ Error handling for file operations

#### LLM Package Tests (`llm/client_test.go`, `llm/client_api_test.go`)
- ✅ HTTP client initialization and configuration
- ✅ Message construction and system prompts
- ✅ Provider-specific API request/response handling
- ✅ **Real API testing with HTTP mocks** (OpenAI, Anthropic)
- ✅ **Authentication header validation** (Bearer tokens, API keys)
- ✅ **Request/response format verification** (JSON schemas)
- ✅ **Error handling for network issues and API errors**
- ✅ **HTTP status code handling** (401, 500, etc.)
- ✅ **Empty response and malformed data testing**
- ✅ JSON serialization/deserialization performance
- ✅ **Mock server testing for API interactions without real API keys**

#### Models Package Tests (`models/models_test.go`)
- ✅ Message type definitions and constants
- ✅ State management structures
- ✅ Bubble Tea message implementations
- ✅ Command type validations
- ✅ Data structure integrity

### 2. **Integration Tests** (Middle Layer)
**Location**: `app/app_test.go`  
**Purpose**: Test component interactions  
**Coverage**: 20% of total tests

#### App Integration Tests
- ✅ Application initialization with different configurations
- ✅ Screen navigation and state transitions
- ✅ Configuration detection and routing logic
- ✅ File system integration with temporary directories
- ✅ Error propagation across application layers
- ✅ Dependency injection and mocking strategies

### 3. **End-to-End Tests** (Top Layer)
**Location**: `e2e/aura_test.go`  
**Purpose**: Test complete user workflows  
**Coverage**: 10% of total tests

#### E2E Test Scenarios
- ✅ First-run experience (Configuration flow)
- ✅ Trusted directory workflows
- ✅ Security prompts for untrusted directories
- ✅ Configuration persistence across restarts
- ✅ Application lifecycle management
- ✅ Concurrent usage scenarios
- ✅ Error handling and recovery
- ✅ Performance characteristics

## 🛠️ Test Implementation Details

### Testing Dependencies
```go
require (
    github.com/stretchr/testify v1.8.4    // Assertions and test utilities
    github.com/spf13/afero v1.10.0        // Mock filesystem operations
)
```

### Mock Strategy
1. **HTTP Client Mocking**: Using `httptest` for LLM API calls
2. **File System Mocking**: Using `afero` for configuration tests
3. **Time Mocking**: For animation and timeout testing
4. **Dependency Injection**: Variables for external dependencies

### Test Data Management
- **testdata/**: Directory for test fixtures and sample configurations
- **Temporary directories**: Created per test for isolation
- **Environment variables**: Mocked for cross-platform testing
- **Configuration files**: Generated programmatically for each test

## 🚀 Running Tests

### Quick Test Execution
```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test ./config

# Run with coverage
go test -cover ./...

# Run only unit tests (fast)
go test -short ./...
```

### Using the Test Runner
```bash
# Run comprehensive test suite
go run test_runner.go

# This provides:
# - Categorized test execution
# - Performance timing
# - Detailed reporting
# - Benchmark results
```

### Test Categories

#### 1. Unit Tests Only
```bash
go test ./config ./llm ./models
```

#### 2. Integration Tests
```bash
go test ./app
```

#### 3. End-to-End Tests
```bash
go test ./e2e
```

#### 4. Benchmark Tests
```bash
go test -bench=. -benchmem ./...
```

### Continuous Integration
```bash
# CI/CD pipeline tests
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 📈 Test Coverage Targets

| Package | Target Coverage | Current Focus |
|---------|----------------|---------------|
| `config` | 95%+ | Core functionality |
| `llm` | 90%+ | API integrations |
| `models` | 85%+ | Data structures |
| `app` | 80%+ | Integration flows |
| `ui` | 70%+ | Component logic |
| **Overall** | **85%+** | **Comprehensive** |

## 🎯 Test Quality Metrics

### Performance Benchmarks
- **App Startup**: < 100ms
- **Config Loading**: < 10ms
- **API Mock Response**: < 1ms
- **Memory Usage**: Minimal allocations

### Reliability Targets
- **Flaky Test Rate**: < 1%
- **Test Execution Time**: < 30 seconds total
- **Coverage Regression**: 0% tolerance
- **Build Success Rate**: 99%+

## 🔧 Testing Best Practices

### Test Structure
```go
func TestFunctionName(t *testing.T) {
    // Arrange
    testData := setupTestData()
    
    // Act
    result, err := functionUnderTest(testData)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Table-Driven Tests
```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name     string
        input    Input
        expected Output
        wantErr  bool
    }{
        // Test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Test Helpers
```go
func setupTestEnvironment(t *testing.T) *TestEnv {
    t.Helper()
    // Setup code
    t.Cleanup(func() {
        // Cleanup code
    })
    return env
}
```

## 🐛 Debugging Tests

### Running Specific Tests
```bash
# Run specific test function
go test -run TestConfigLoad ./config

# Run tests matching pattern
go test -run "TestConfig.*" ./config

# Debug with verbose output
go test -v -run TestConfigLoad ./config
```

### Test Debugging Techniques
1. **Isolation**: Run single tests to identify issues
2. **Logging**: Add debug output in test functions
3. **Mocking**: Verify mock calls and interactions
4. **Environment**: Check test environment setup

## 📊 Test Reporting

### Coverage Reports
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Coverage by function
go tool cover -func=coverage.out
```

### Benchmark Reports
```bash
# Compare benchmarks
go test -bench=. -count=5 ./... > bench.txt
benchstat bench.txt
```

## 🔄 Test Maintenance

### Adding New Tests
1. **Identify test category** (Unit/Integration/E2E)
2. **Create test file** following naming conventions
3. **Implement test cases** with proper structure
4. **Add to CI pipeline** if needed
5. **Update documentation**

### Updating Existing Tests
1. **Maintain backward compatibility**
2. **Update test data** when APIs change
3. **Refactor common patterns** into helpers
4. **Keep tests focused** and isolated

### Test Cleanup
- **Remove obsolete tests** when features are removed
- **Consolidate duplicate tests** to reduce maintenance
- **Update test documentation** regularly
- **Monitor test performance** and optimize slow tests

## 🎨 Test Organization

```
aura/
├── config/
│   ├── config.go
│   └── config_test.go      # Unit tests
├── llm/
│   ├── client.go
│   └── client_test.go      # Unit tests with mocks
├── models/
│   ├── models.go
│   └── models_test.go      # Unit tests
├── app/
│   ├── app.go
│   └── app_test.go         # Integration tests
├── e2e/
│   └── aura_test.go        # End-to-end tests
├── testdata/               # Test fixtures
├── test_runner.go          # Test orchestration
└── TESTING.md             # This documentation
```

## 🚀 Future Testing Enhancements

### Planned Improvements
- [ ] **Visual Regression Tests**: UI component testing
- [ ] **API Contract Tests**: Provider API validation
- [ ] **Load Testing**: High-concurrency scenarios
- [ ] **Security Testing**: Input validation and sanitization
- [ ] **Cross-Platform Tests**: Windows/macOS/Linux validation

### Advanced Testing Tools
- [ ] **Fuzzing**: Automated input generation
- [ ] **Property-Based Testing**: Generative test cases
- [ ] **Mutation Testing**: Test quality validation
- [ ] **Performance Profiling**: Memory and CPU analysis

The comprehensive test suite ensures Aura maintains high quality, reliability, and performance while providing confidence for continuous development and deployment.
