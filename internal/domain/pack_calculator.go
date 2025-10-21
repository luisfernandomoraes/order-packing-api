package domain

import (
	"sort"
	"sync"
)

// PackResult represents the calculation result containing the order details,
// total items to be shipped, and the distribution of packs.
type PackResult struct {
	Order      int         `json:"order"`
	TotalItems int         `json:"total_items"`
	Packs      map[int]int `json:"packs"`
	PackSizes  []int       `json:"pack_sizes_used"`
}

// PackCalculator is responsible for calculating the optimal pack combination
// to fulfill customer orders while minimizing items and packs sent.
type PackCalculator struct {
	mu        sync.RWMutex
	packSizes []int
}

// NewPackCalculator creates a new calculator instance with the given pack sizes.
// The pack sizes are automatically sorted in ascending order for optimization.
func NewPackCalculator(sizes []int) *PackCalculator {
	sortedSizes := make([]int, len(sizes))
	copy(sortedSizes, sizes)
	sort.Ints(sortedSizes)

	return &PackCalculator{
		packSizes: sortedSizes,
	}
}

// solution represents a possible pack combination during the calculation process.
type solution struct {
	totalItems     int
	packsBySize    map[int]int
	totalPackCount int
}

// Calculate computes the optimal pack combination for the given order quantity
// using dynamic programming to ensure the best solution according to business rules:
//
// Rules (in priority order):
//  1. Only whole packs can be sent (packs cannot be broken)
//  2. Send the minimum number of items to fulfill the order
//  3. Send the minimum number of packs (within rule 2 constraints)
//
// Examples:
//
//	order = 1    -> TotalItems: 250,   Packs: {250: 1}
//	order = 251  -> TotalItems: 500,   Packs: {500: 1}
//	order = 501  -> TotalItems: 750,   Packs: {500: 1, 250: 1}
//	order = 12001-> TotalItems: 12250, Packs: {5000: 2, 2000: 1, 250: 1}
func (pc *PackCalculator) Calculate(order int) PackResult {
	packSizes := pc.GetPackSizes()

	if order <= 0 {
		return PackResult{
			Order:      order,
			TotalItems: 0,
			Packs:      make(map[int]int),
			PackSizes:  packSizes,
		}
	}

	if len(packSizes) == 0 {
		return PackResult{
			Order:      order,
			TotalItems: 0,
			Packs:      make(map[int]int),
			PackSizes:  packSizes,
		}
	}

	largestPack := packSizes[len(packSizes)-1]
	searchLimit := order + largestPack

	optimalSolutions := make(map[int]*solution)
	optimalSolutions[0] = &solution{
		totalItems:     0,
		packsBySize:    make(map[int]int),
		totalPackCount: 0,
	}

	pc.buildOptimalSolutions(optimalSolutions, searchLimit, packSizes)

	return pc.findBestSolutionForOrder(optimalSolutions, order, searchLimit, packSizes)
}

// buildOptimalSolutions fills the dynamic programming table with optimal solutions.
func (pc *PackCalculator) buildOptimalSolutions(optimalSolutions map[int]*solution, limit int, packSizes []int) {
	for currentQuantity := 1; currentQuantity <= limit; currentQuantity++ {
		for _, packSize := range packSizes {
			if currentQuantity >= packSize {
				previousSolution := optimalSolutions[currentQuantity-packSize]
				if previousSolution == nil {
					continue
				}

				newSolution := pc.createSolutionWithPack(previousSolution, packSize)
				currentBestSolution := optimalSolutions[currentQuantity]

				if pc.isBetterSolution(newSolution, currentBestSolution) {
					optimalSolutions[currentQuantity] = newSolution
				}
			}
		}
	}
}

// createSolutionWithPack creates a new solution by adding one pack to an existing solution.
func (pc *PackCalculator) createSolutionWithPack(baseSolution *solution, packSize int) *solution {
	newPacks := make(map[int]int)
	for size, quantity := range baseSolution.packsBySize {
		newPacks[size] = quantity
	}
	newPacks[packSize]++

	return &solution{
		totalItems:     baseSolution.totalItems + packSize,
		packsBySize:    newPacks,
		totalPackCount: baseSolution.totalPackCount + 1,
	}
}

// isBetterSolution determines if the new solution is better than the current one.
// Priority: fewer items first, then fewer packs.
func (pc *PackCalculator) isBetterSolution(newSolution, currentSolution *solution) bool {
	if currentSolution == nil {
		return true
	}

	if newSolution.totalItems < currentSolution.totalItems {
		return true
	}

	if newSolution.totalItems == currentSolution.totalItems &&
		newSolution.totalPackCount < currentSolution.totalPackCount {
		return true
	}

	return false
}

// findBestSolutionForOrder searches for the first valid solution that meets or exceeds the order.
func (pc *PackCalculator) findBestSolutionForOrder(
	optimalSolutions map[int]*solution,
	order int,
	searchLimit int,
	packSizes []int,
) PackResult {
	for quantity := order; quantity <= searchLimit; quantity++ {
		if solution := optimalSolutions[quantity]; solution != nil {
			return PackResult{
				Order:      order,
				TotalItems: solution.totalItems,
				Packs:      solution.packsBySize,
				PackSizes:  packSizes,
			}
		}
	}

	return PackResult{
		Order:      order,
		TotalItems: 0,
		Packs:      make(map[int]int),
		PackSizes:  packSizes,
	}
}

// UpdatePackSizes updates the available pack sizes and re-sorts them.
func (pc *PackCalculator) UpdatePackSizes(sizes []int) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	sortedSizes := make([]int, len(sizes))
	copy(sortedSizes, sizes)
	sort.Ints(sortedSizes)
	pc.packSizes = sortedSizes
}

// GetPackSizes returns the currently configured pack sizes.
func (pc *PackCalculator) GetPackSizes() []int {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	result := make([]int, len(pc.packSizes))
	copy(result, pc.packSizes)
	return result
}

// GetTotalPackCount returns the total number of packs in this result.
func (pr *PackResult) GetTotalPackCount() int {
	totalPacks := 0
	for _, quantity := range pr.Packs {
		totalPacks += quantity
	}
	return totalPacks
}

// GetSurplus returns the number of extra items being sent beyond the order.
func (pr *PackResult) GetSurplus() int {
	return pr.TotalItems - pr.Order
}
