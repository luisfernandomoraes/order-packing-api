package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPackCalculator(t *testing.T) {
	tests := []struct {
		name           string
		inputSizes     []int
		expectedSorted []int
	}{
		{
			name:           "should sort sizes in ascending order",
			inputSizes:     []int{5000, 250, 1000, 500, 2000},
			expectedSorted: []int{250, 500, 1000, 2000, 5000},
		},
		{
			name:           "should handle already sorted sizes",
			inputSizes:     []int{1, 5, 10, 25},
			expectedSorted: []int{1, 5, 10, 25},
		},
		{
			name:           "should handle single size",
			inputSizes:     []int{100},
			expectedSorted: []int{100},
		},
		{
			name:           "should handle duplicate sizes",
			inputSizes:     []int{250, 500, 250, 1000},
			expectedSorted: []int{250, 250, 500, 1000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator := NewPackCalculator(tt.inputSizes)

			assert.NotNil(t, calculator)
			assert.Equal(t, tt.expectedSorted, calculator.packSizes)
		})
	}
}

func TestPackCalculator_Calculate(t *testing.T) {
	defaultSizes := []int{250, 500, 1000, 2000, 5000}

	tests := []struct {
		name               string
		packSizes          []int
		order              int
		expectedTotalItems int
		expectedPacks      map[int]int
	}{
		{
			name:               "order 1 item - should use smallest pack",
			packSizes:          defaultSizes,
			order:              1,
			expectedTotalItems: 250,
			expectedPacks:      map[int]int{250: 1},
		},
		{
			name:               "order exactly 250 - perfect match",
			packSizes:          defaultSizes,
			order:              250,
			expectedTotalItems: 250,
			expectedPacks:      map[int]int{250: 1},
		},
		{
			name:               "order 251 - should use 500 (fewer items than 2x250)",
			packSizes:          defaultSizes,
			order:              251,
			expectedTotalItems: 500,
			expectedPacks:      map[int]int{500: 1},
		},
		{
			name:               "order 501 - should use 500+250",
			packSizes:          defaultSizes,
			order:              501,
			expectedTotalItems: 750,
			expectedPacks:      map[int]int{500: 1, 250: 1},
		},
		{
			name:               "order 12001 - complex combination",
			packSizes:          defaultSizes,
			order:              12001,
			expectedTotalItems: 12250,
			expectedPacks:      map[int]int{5000: 2, 2000: 1, 250: 1},
		},
		{
			name:               "order 1000 - exact match",
			packSizes:          defaultSizes,
			order:              1000,
			expectedTotalItems: 1000,
			expectedPacks:      map[int]int{1000: 1},
		},
		{
			name:               "order 5000 - exact match largest pack",
			packSizes:          defaultSizes,
			order:              5000,
			expectedTotalItems: 5000,
			expectedPacks:      map[int]int{5000: 1},
		},
		{
			name:               "order 0 - should return empty",
			packSizes:          defaultSizes,
			order:              0,
			expectedTotalItems: 0,
			expectedPacks:      map[int]int{},
		},
		{
			name:               "order negative - should return empty",
			packSizes:          defaultSizes,
			order:              -100,
			expectedTotalItems: 0,
			expectedPacks:      map[int]int{},
		},
		{
			name:               "money example - 100 cents (4 quarters)",
			packSizes:          []int{1, 5, 10, 25},
			order:              100,
			expectedTotalItems: 100,
			expectedPacks:      map[int]int{25: 4},
		},
		{
			name:               "money example - 28 cents",
			packSizes:          []int{1, 5, 10, 25},
			order:              28,
			expectedTotalItems: 28,
			expectedPacks:      map[int]int{25: 1, 1: 3},
		},
		{
			name:               "order 2500 - multiple of smallest and largest",
			packSizes:          defaultSizes,
			order:              2500,
			expectedTotalItems: 2500,
			expectedPacks:      map[int]int{2000: 1, 500: 1},
		},
		{
			name:               "order 7501 - requires multiple large packs",
			packSizes:          defaultSizes,
			order:              7501,
			expectedTotalItems: 7750,
			expectedPacks:      map[int]int{5000: 1, 2000: 1, 500: 1, 250: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator := NewPackCalculator(tt.packSizes)
			result := calculator.Calculate(tt.order)

			assert.Equal(t, tt.order, result.Order)
			assert.Equal(t, tt.expectedTotalItems, result.TotalItems)
			assert.Equal(t, tt.expectedPacks, result.Packs)
			assert.Equal(t, calculator.packSizes, result.PackSizes)
		})
	}
}

func TestPackCalculator_UpdatePackSizes(t *testing.T) {
	tests := []struct {
		name           string
		initialSizes   []int
		newSizes       []int
		expectedSorted []int
	}{
		{
			name:           "should update and sort pack sizes",
			initialSizes:   []int{250, 500},
			newSizes:       []int{100, 200, 300},
			expectedSorted: []int{100, 200, 300},
		},
		{
			name:           "should sort new unsorted sizes",
			initialSizes:   []int{250, 500},
			newSizes:       []int{500, 100, 300},
			expectedSorted: []int{100, 300, 500},
		},
		{
			name:           "should handle empty update gracefully",
			initialSizes:   []int{250, 500},
			newSizes:       []int{},
			expectedSorted: []int{},
		},
		{
			name:           "should handle single size update",
			initialSizes:   []int{250, 500, 1000},
			newSizes:       []int{750},
			expectedSorted: []int{750},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator := NewPackCalculator(tt.initialSizes)
			calculator.UpdatePackSizes(tt.newSizes)

			assert.Equal(t, tt.expectedSorted, calculator.packSizes)
		})
	}
}

func TestPackCalculator_GetPackSizes(t *testing.T) {
	tests := []struct {
		name     string
		sizes    []int
		expected []int
	}{
		{
			name:     "should return pack sizes",
			sizes:    []int{250, 500, 1000},
			expected: []int{250, 500, 1000},
		},
		{
			name:     "should return sorted sizes",
			sizes:    []int{1000, 250, 500},
			expected: []int{250, 500, 1000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator := NewPackCalculator(tt.sizes)
			result := calculator.GetPackSizes()

			assert.Equal(t, tt.expected, result)
		})
	}

	t.Run("should return copy to prevent external modification", func(t *testing.T) {
		calculator := NewPackCalculator([]int{250, 500, 1000})
		sizes := calculator.GetPackSizes()

		// Modify returned slice
		sizes[0] = 999

		// Verify original is not modified
		assert.Equal(t, 250, calculator.packSizes[0])
	})
}

func TestPackResult_GetTotalPackCount(t *testing.T) {
	tests := []struct {
		name          string
		packs         map[int]int
		expectedTotal int
	}{
		{
			name:          "single pack type",
			packs:         map[int]int{250: 1},
			expectedTotal: 1,
		},
		{
			name:          "multiple pack types",
			packs:         map[int]int{500: 1, 250: 1},
			expectedTotal: 2,
		},
		{
			name:          "multiple quantities of same pack",
			packs:         map[int]int{5000: 2, 2000: 1, 250: 1},
			expectedTotal: 4,
		},
		{
			name:          "empty packs",
			packs:         map[int]int{},
			expectedTotal: 0,
		},
		{
			name:          "many packs",
			packs:         map[int]int{25: 4},
			expectedTotal: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PackResult{Packs: tt.packs}
			total := result.GetTotalPackCount()

			assert.Equal(t, tt.expectedTotal, total)
		})
	}
}

func TestPackResult_GetSurplus(t *testing.T) {
	tests := []struct {
		name            string
		order           int
		totalItems      int
		expectedSurplus int
	}{
		{
			name:            "surplus with 501 order",
			order:           501,
			totalItems:      750,
			expectedSurplus: 249,
		},
		{
			name:            "exact match - no surplus",
			order:           250,
			totalItems:      250,
			expectedSurplus: 0,
		},
		{
			name:            "large surplus",
			order:           1,
			totalItems:      250,
			expectedSurplus: 249,
		},
		{
			name:            "complex order surplus",
			order:           12001,
			totalItems:      12250,
			expectedSurplus: 249,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PackResult{
				Order:      tt.order,
				TotalItems: tt.totalItems,
			}
			surplus := result.GetSurplus()

			assert.Equal(t, tt.expectedSurplus, surplus)
		})
	}
}

func TestPackCalculator_EdgeCases(t *testing.T) {
	tests := []struct {
		name               string
		packSizes          []int
		order              int
		expectedTotalItems int
		minExpectedPacks   int
	}{
		{
			name:               "single pack size available",
			packSizes:          []int{100},
			order:              350,
			expectedTotalItems: 400,
			minExpectedPacks:   4,
		},
		{
			name:               "very large order",
			packSizes:          []int{250, 500, 1000, 2000, 5000},
			order:              100000,
			expectedTotalItems: 100000,
			minExpectedPacks:   20,
		},
		{
			name:               "prime number order",
			packSizes:          []int{250, 500, 1000, 2000, 5000},
			order:              1013,
			expectedTotalItems: 1250,
			minExpectedPacks:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator := NewPackCalculator(tt.packSizes)
			result := calculator.Calculate(tt.order)

			assert.GreaterOrEqual(t, result.TotalItems, tt.order, "should fulfill order")
			assert.Equal(t, tt.expectedTotalItems, result.TotalItems)
			assert.GreaterOrEqual(t, result.GetTotalPackCount(), tt.minExpectedPacks)
		})
	}
}

func TestPackCalculator_BusinessRules(t *testing.T) {
	t.Run("Rule 1: Only whole packs can be sent", func(t *testing.T) {
		calculator := NewPackCalculator([]int{250, 500, 1000})
		result := calculator.Calculate(501)

		// Verify all pack counts are whole numbers
		for _, count := range result.Packs {
			assert.Greater(t, count, 0)
		}
	})

	t.Run("Rule 2: Send minimum items to fulfill order", func(t *testing.T) {
		calculator := NewPackCalculator([]int{250, 500, 1000})
		result := calculator.Calculate(251)

		// 1x500 = 500 is better than 2x250 = 500 (same items, fewer packs)
		// But both are better than 1x1000 = 1000 (more items)
		assert.Equal(t, 500, result.TotalItems)
		assert.LessOrEqual(t, result.TotalItems-result.Order, 250)
	})

	t.Run("Rule 3: Send minimum packs (within rule 2)", func(t *testing.T) {
		calculator := NewPackCalculator([]int{250, 500, 1000})
		result := calculator.Calculate(251)

		// Should prefer 1x500 over 2x250 (both give 500 items)
		assert.Equal(t, 1, result.GetTotalPackCount())
		assert.Equal(t, map[int]int{500: 1}, result.Packs)
	})
}

func TestPackCalculator_Calculate_RequireValid(t *testing.T) {
	calculator := NewPackCalculator([]int{250, 500, 1000, 2000, 5000})

	t.Run("result should always fulfill order", func(t *testing.T) {
		orders := []int{1, 100, 251, 500, 501, 999, 1000, 5001, 12001}

		for _, order := range orders {
			result := calculator.Calculate(order)

			require.GreaterOrEqual(t, result.TotalItems, order,
				"Order %d: total items %d should be >= order", order, result.TotalItems)
		}
	})
}

// Benchmark tests
func BenchmarkCalculate_SmallOrder(b *testing.B) {
	calculator := NewPackCalculator([]int{250, 500, 1000, 2000, 5000})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		calculator.Calculate(501)
	}
}

func BenchmarkCalculate_MediumOrder(b *testing.B) {
	calculator := NewPackCalculator([]int{250, 500, 1000, 2000, 5000})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		calculator.Calculate(12001)
	}
}

func BenchmarkCalculate_LargeOrder(b *testing.B) {
	calculator := NewPackCalculator([]int{250, 500, 1000, 2000, 5000})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		calculator.Calculate(50001)
	}
}

func TestPackCalculator_Calculate_WithEmptyPackSizes(t *testing.T) {
	tests := []struct {
		name               string
		initialSizes       []int
		updatedSizes       []int
		order              int
		expectedTotalItems int
		expectedPacks      map[int]int
	}{
		{
			name:               "should handle empty pack sizes with positive order",
			initialSizes:       []int{250, 500},
			updatedSizes:       []int{},
			order:              100,
			expectedTotalItems: 0,
			expectedPacks:      map[int]int{},
		},
		{
			name:               "should handle empty pack sizes with zero order",
			initialSizes:       []int{},
			updatedSizes:       []int{},
			order:              0,
			expectedTotalItems: 0,
			expectedPacks:      map[int]int{},
		},
		{
			name:               "should handle empty pack sizes with negative order",
			initialSizes:       []int{},
			updatedSizes:       []int{},
			order:              -50,
			expectedTotalItems: 0,
			expectedPacks:      map[int]int{},
		},
		{
			name:               "should handle empty pack sizes after update",
			initialSizes:       []int{100, 200, 300},
			updatedSizes:       []int{},
			order:              500,
			expectedTotalItems: 0,
			expectedPacks:      map[int]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator := NewPackCalculator(tt.initialSizes)
			if tt.updatedSizes != nil {
				calculator.UpdatePackSizes(tt.updatedSizes)
			}

			result := calculator.Calculate(tt.order)

			assert.Equal(t, tt.order, result.Order)
			assert.Equal(t, tt.expectedTotalItems, result.TotalItems)
			assert.Equal(t, tt.expectedPacks, result.Packs)
		})
	}
}

func TestPackCalculator_Calculate_OptimalItemMinimization(t *testing.T) {
	tests := []struct {
		name               string
		packSizes          []int
		order              int
		expectedTotalItems int
		expectedPacks      map[int]int
		expectedSurplus    int
	}{
		{
			name:               "should prefer fewer items with same pack count",
			packSizes:          []int{3, 5},
			order:              6,
			expectedTotalItems: 6,
			expectedPacks:      map[int]int{3: 2},
			expectedSurplus:    0,
		},
		{
			name:               "should minimize items even with different configurations",
			packSizes:          []int{7, 11},
			order:              14,
			expectedTotalItems: 14,
			expectedPacks:      map[int]int{7: 2},
			expectedSurplus:    0,
		},
		{
			name:               "should prioritize fewer items over fewer packs",
			packSizes:          []int{3, 7, 13},
			order:              9,
			expectedTotalItems: 9,
			expectedPacks:      map[int]int{3: 3},
			expectedSurplus:    0,
		},
		{
			name:               "should find exact match with multiple small packs",
			packSizes:          []int{4, 7},
			order:              12,
			expectedTotalItems: 12,
			expectedPacks:      map[int]int{4: 3},
			expectedSurplus:    0,
		},
		{
			name:               "should minimize surplus with mixed pack sizes",
			packSizes:          []int{6, 9},
			order:              15,
			expectedTotalItems: 15,
			expectedPacks:      map[int]int{6: 1, 9: 1},
			expectedSurplus:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator := NewPackCalculator(tt.packSizes)
			result := calculator.Calculate(tt.order)

			assert.Equal(t, tt.expectedTotalItems, result.TotalItems)
			assert.Equal(t, tt.expectedPacks, result.Packs)
			assert.Equal(t, tt.expectedSurplus, result.GetSurplus())
		})
	}
}

func TestPackCalculator_Calculate_EdgeCasesWithLargePacks(t *testing.T) {
	tests := []struct {
		name               string
		packSizes          []int
		order              int
		expectedTotalItems int
		expectedPacks      map[int]int
		minSurplus         int
	}{
		{
			name:               "should handle very large pack with small order",
			packSizes:          []int{10000},
			order:              1,
			expectedTotalItems: 10000,
			expectedPacks:      map[int]int{10000: 1},
			minSurplus:         9999,
		},
		{
			name:               "should handle single large pack size",
			packSizes:          []int{1000},
			order:              500,
			expectedTotalItems: 1000,
			expectedPacks:      map[int]int{1000: 1},
			minSurplus:         500,
		},
		{
			name:               "should handle prime pack sizes",
			packSizes:          []int{17, 23, 29},
			order:              50,
			expectedTotalItems: 51,
			expectedPacks:      map[int]int{17: 3},
			minSurplus:         1,
		},
		{
			name:               "should handle coprime pack sizes",
			packSizes:          []int{13, 17},
			order:              30,
			expectedTotalItems: 30,
			expectedPacks:      map[int]int{13: 1, 17: 1},
			minSurplus:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator := NewPackCalculator(tt.packSizes)
			result := calculator.Calculate(tt.order)

			assert.Equal(t, tt.order, result.Order)
			assert.Equal(t, tt.expectedTotalItems, result.TotalItems)
			assert.Equal(t, tt.expectedPacks, result.Packs)
			assert.GreaterOrEqual(t, result.GetSurplus(), tt.minSurplus)
			assert.GreaterOrEqual(t, result.TotalItems, tt.order)
		})
	}
}
