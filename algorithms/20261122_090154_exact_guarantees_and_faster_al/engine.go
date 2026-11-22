package main

// Longest returns the length of the longest subsequence of s.seq that
// contains at most s.k direction changes (i.e., at most k+1 monotone runs).
// The algorithm runs in O(n²·k) time and O(n·k) space, where n = len(seq).
func (s *RollerCoasterSolver) Longest() int {
	n := len(s.seq)
	if n == 0 {
		return 0
	}
	if s.k >= n-1 {
		// Trivially we can take the whole sequence.
		return n
	}

	// dp[i][c][d] = longest length of a subsequence ending at index i,
	// using exactly c direction changes, and whose last step has direction d.
	// c ranges 0..k, d is 0 (up) or 1 (down).
	dp := make([][][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([][]int, s.k+1)
		for c := 0; c <= s.k; c++ {
			dp[i][c] = []int{0, 0} // [up, down]
		}
	}

	maxLen := 1 // any single element forms a valid subsequence

	// Initialise each position as a subsequence of length 1 with 0 changes.
	for i := 0; i < n; i++ {
		dp[i][0][DirUp] = 1
		dp[i][0][DirDown] = 1
	}

	// Main DP transition.
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			if s.seq[j] == s.seq[i] {
				// Equal values do not create a monotone step; skip.
				continue
			}
			var dir Direction
			if s.seq[j] < s.seq[i] {
				dir = DirUp
			} else {
				dir = DirDown
			}
			opp := DirDown
			if dir == DirDown {
				opp = DirUp
			}
			for cPrev := 0; cPrev <= s.k; cPrev++ {
				// Extend a subsequence that previously ended with the same direction.
				if dp[j][cPrev][dir] > 0 {
					if dp[i][cPrev][dir] < dp[j][cPrev][dir]+1 {
						dp[i][cPrev][dir] = dp[j][cPrev][dir] + 1
						if dp[i][cPrev][dir] > maxLen {
							maxLen = dp[i][cPrev][dir]
						}
					}
				}
				// Extend a subsequence that ended with the opposite direction,
				// incurring a direction change.
				if dp[j][cPrev][opp] > 0 {
					cNew := cPrev + 1
					if cNew <= s.k {
						if dp[i][cNew][dir] < dp[j][cPrev][opp]+1 {
							dp[i][cNew][dir] = dp[j][cPrev][opp] + 1
							if dp[i][cNew][dir] > maxLen {
								maxLen = dp[i][cNew][dir]
							}
						}
					}
				}
			}
		}
	}
	return maxLen
}
