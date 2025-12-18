package parser

import (
	"fmt"
	"strings"
	"testing"
)

// Simple GRL rule for basic benchmarking
const simpleRule = `
rule SimpleRule "A simple rule" salience 10 {
    when
        Order.Status == "pending"
    then
        Order.Status = "processing";
        Retract("SimpleRule");
}
`

// Complex rule with multiple conditions
const complexRule = `
rule ComplexRule "A complex rule with many conditions" salience 100 {
    when
        Order.Status == "pending" &&
        Order.Total > 0 &&
        Order.ItemCount > 0 &&
        Customer.IsVerified == true &&
        Customer.Balance >= Order.Total
    then
        Order.Status = "approved";
        Order.ApprovedAt = Now();
        Customer.Balance = Customer.Balance - Order.Total;
        Log("Order approved: " + Order.ID);
        Retract("ComplexRule");
}
`

// generateRules creates n rules for benchmarking
func generateRules(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, `
rule Rule%d "Rule number %d" salience %d {
    when
        Data.Value%d == %d
    then
        Data.Processed%d = true;
        Retract("Rule%d");
}
`, i, i, 100-i, i, i*10, i, i)
	}
	return sb.String()
}

func BenchmarkParseString_SimpleRule(b *testing.B) {
	p := NewParser()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.ParseString("test.grl", simpleRule)
	}
}

func BenchmarkParseString_ComplexRule(b *testing.B) {
	p := NewParser()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.ParseString("test.grl", complexRule)
	}
}

func BenchmarkParseString_10Rules(b *testing.B) {
	p := NewParser()
	content := generateRules(10)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.ParseString("test.grl", content)
	}
}

func BenchmarkParseString_50Rules(b *testing.B) {
	p := NewParser()
	content := generateRules(50)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.ParseString("test.grl", content)
	}
}

func BenchmarkParseString_100Rules(b *testing.B) {
	p := NewParser()
	content := generateRules(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.ParseString("test.grl", content)
	}
}

func BenchmarkParseString_SyntaxError(b *testing.B) {
	p := NewParser()
	content := `
rule BadRule {
    when
    then
}
`
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.ParseString("test.grl", content)
	}
}

// BenchmarkParseWithANTLR benchmarks only the ANTLR parsing step
func BenchmarkParseWithANTLR_SimpleRule(b *testing.B) {
	p := NewParser()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.parseWithANTLR(simpleRule)
	}
}

func BenchmarkParseWithANTLR_100Rules(b *testing.B) {
	p := NewParser()
	content := generateRules(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.parseWithANTLR(content)
	}
}

// Benchmark memory allocations
func BenchmarkParseString_Allocs(b *testing.B) {
	p := NewParser()
	content := generateRules(10)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.ParseString("test.grl", content)
	}
}

// Benchmark parallel parsing
func BenchmarkParseString_Parallel(b *testing.B) {
	p := NewParser()
	content := generateRules(10)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.ParseString("test.grl", content)
		}
	})
}
