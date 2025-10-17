package homebrew

import (
	"context"
	"testing"
)

// BenchmarkSearch benchmarks the parallel search implementation
func BenchmarkSearch(b *testing.B) {
	client, err := NewClient()
	if err != nil {
		b.Skip("Skipping: brew not found")
	}

	// Pre-populate cache
	client.formulaeCache = make([]FormulaListItem, 10000)
	for i := 0; i < 10000; i++ {
		client.formulaeCache[i] = FormulaListItem{
			Name: "package-" + string(rune(i)),
			Desc: "Description for package " + string(rune(i)),
		}
	}

	client.casksCache = make([]CaskListItem, 5000)
	for i := 0; i < 5000; i++ {
		client.casksCache[i] = CaskListItem{
			Token: "cask-" + string(rune(i)),
			Name:  []string{"Cask " + string(rune(i))},
			Desc:  "Description for cask " + string(rune(i)),
		}
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = client.Search(ctx, "package")
	}
}

// BenchmarkLoadFormulaeAndCasks benchmarks the parallel loading implementation
func BenchmarkLoadFormulaeAndCasks(b *testing.B) {
	client, err := NewClient()
	if err != nil {
		b.Skip("Skipping: brew not found")
	}

	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.loadFormulaeAndCasks(ctx)
	}
}

// BenchmarkSearchLargeDataset benchmarks search with a realistic dataset
func BenchmarkSearchLargeDataset(b *testing.B) {
	client, err := NewClient()
	if err != nil {
		b.Skip("Skipping: brew not found")
	}

	// Simulate realistic Homebrew data (~7000 formulae, ~6000 casks)
	client.formulaeCache = make([]FormulaListItem, 7000)
	for i := 0; i < 7000; i++ {
		client.formulaeCache[i] = FormulaListItem{
			Name: "formula-" + string(rune(i)),
			Desc: "A formula description with some keywords like git, node, python",
		}
	}

	client.casksCache = make([]CaskListItem, 6000)
	for i := 0; i < 6000; i++ {
		client.casksCache[i] = CaskListItem{
			Token: "cask-" + string(rune(i)),
			Name:  []string{"Cask Application " + string(rune(i))},
			Desc:  "A cask description with some keywords",
		}
	}

	ctx := context.Background()
	searches := []string{"git", "node", "python", "docker", "visual"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		term := searches[i%len(searches)]
		_, _, _ = client.Search(ctx, term)
	}
}
