package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Test runner for Aura application
// This script runs all tests with proper categorization and reporting

func main() {
	fmt.Println("🧪 Aura Test Suite Runner")
	fmt.Println("==========================")

	// Track overall results
	var results []TestResult
	startTime := time.Now()

	// 1. Unit Tests
	fmt.Println("\n📦 Running Unit Tests...")
	results = append(results, runTestSuite("Unit Tests", []string{
		"./config/...",
		"./llm/...",
		"./models/...",
	}, false))

	// 2. Integration Tests
	fmt.Println("\n🔗 Running Integration Tests...")
	results = append(results, runTestSuite("Integration Tests", []string{
		"./app/...",
	}, false))

	// 3. End-to-End Tests (with timeout)
	fmt.Println("\n🎯 Running End-to-End Tests...")
	results = append(results, runTestSuite("E2E Tests", []string{
		"./e2e/...",
	}, false))

	// 4. Benchmark Tests
	fmt.Println("\n⚡ Running Benchmark Tests...")
	results = append(results, runBenchmarks())

	// Summary
	duration := time.Since(startTime)
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("📊 Test Summary (Total time: %v)\n", duration)
	fmt.Println(strings.Repeat("=", 50))

	totalTests := 0
	totalPassed := 0
	totalFailed := 0

	for _, result := range results {
		fmt.Printf("%-20s: ", result.Name)
		if result.Success {
			fmt.Printf("✅ PASS (%d tests)\n", result.TestCount)
			totalPassed += result.TestCount
		} else {
			fmt.Printf("❌ FAIL (%d tests)\n", result.TestCount)
			totalFailed += result.TestCount
		}
		totalTests += result.TestCount
	}

	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Total: %d tests, %d passed, %d failed\n", totalTests, totalPassed, totalFailed)

	if totalFailed > 0 {
		fmt.Printf("\n❌ Some tests failed. See output above for details.\n")
		os.Exit(1)
	}
	fmt.Printf("\n✅ All tests passed! 🎉\n")
}

type TestResult struct {
	Name      string
	Success   bool
	TestCount int
	Duration  time.Duration
	Output    string
}

func runTestSuite(name string, packages []string, isShort bool) TestResult {
	start := time.Now()

	args := []string{"test", "-v"}
	if isShort {
		args = append(args, "-short")
	}
	args = append(args, packages...)

	cmd := exec.Command("go", args...)
	output, err := cmd.CombinedOutput()

	duration := time.Since(start)

	// Count tests from output
	testCount := countTests(string(output))

	success := err == nil

	if !success {
		fmt.Printf("❌ %s failed:\n%s\n", name, string(output))
	} else {
		fmt.Printf("✅ %s passed (%d tests in %v)\n", name, testCount, duration)
	}

	return TestResult{
		Name:      name,
		Success:   success,
		TestCount: testCount,
		Duration:  duration,
		Output:    string(output),
	}
}

func runBenchmarks() TestResult {
	start := time.Now()

	cmd := exec.Command("go", "test", "-bench=.", "-benchmem", "./config", "./llm", "./models", "./app", "./e2e")
	output, err := cmd.CombinedOutput()

	duration := time.Since(start)

	benchCount := countBenchmarks(string(output))
	success := err == nil

	if !success {
		fmt.Printf("❌ Benchmarks failed:\n%s\n", string(output))
	} else {
		fmt.Printf("✅ Benchmarks completed (%d benchmarks in %v)\n", benchCount, duration)
	}

	return TestResult{
		Name:      "Benchmarks",
		Success:   success,
		TestCount: benchCount,
		Duration:  duration,
		Output:    string(output),
	}
}

func countTests(output string) int {
	count := 0
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "=== RUN") {
			count++
		}
	}
	return count
}

func countBenchmarks(output string) int {
	count := 0
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Benchmark") && strings.Contains(line, "ns/op") {
			count++
		}
	}
	return count
}
