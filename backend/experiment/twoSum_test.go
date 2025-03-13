package experiment

import "testing"

func TestTwoSum(t *testing.T) {
    testCases := []struct {
        nums   []int
        target int
        want  []int
    }{
        // Test case 1: Simple case where two numbers add up to the target.
        {[]int{2, 7, 11, 15}, 9, []int{0, 1}},
        // Test case 2: Multiple possible pairs but returns the first occurrence.
        {[]int{3, 2, 4}, 6, []int{1, 2}},
        // Test case 3: Sum of same numbers where multiple occurrences exist.
        {[]int{3, 3}, 6, []int{0, 1}},
        // Test case 4: Negative numbers sum to target.
        {[]int{-1, -2}, -3, []int{0, 1}},
        // Test case 5: No pair found, should return empty slice.
        {[]int{1}, 3, []int{}},
    }

    for _, tc := range testCases {
        result := twoSum(tc.nums, tc.target)
        if !equalSlices(result, tc.want) {
            t.Errorf("TwoSum(%v, %d) = %v, want %v", tc.nums, tc.target, result, tc.want)
        }
    }
}

// Helper function to check if two slices are equal.
func equalSlices(a, b []int) bool {
    if len(a) != len(b) {
        return false
    }
    for i := 0; i < len(a); i++ {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}