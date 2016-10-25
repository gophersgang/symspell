package symspell

import "unicode/utf8"

// EditMatch ...
func EditMatch(a, b string, ed int) bool {
	f := make([]int, utf8.RuneCountInString(b)+1)

	for j := range f {
		f[j] = j
	}

	for _, ca := range a {
		j := 1
		fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
		f[0]++
		x := ed + 1
		for _, cb := range b {
			mn := min3(sel(ca != cb, fj1+1, fj1), f[j]+1, f[j-1]+1)
			fj1, f[j] = f[j], mn // save f[j] to fj1(j is about to increase), update f[j] to mn
			j++
			x = min(x, mn)
		}
		if x > ed {
			return false
		}
	}
	return f[len(f)-1] <= ed
}

// Levenshtein ...
func Levenshtein(a, b string) int {
	f := make([]int, utf8.RuneCountInString(b)+1)

	for j := range f {
		f[j] = j
	}

	for _, ca := range a {
		j := 1
		fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
		f[0]++
		for _, cb := range b {
			mn := min(f[j]+1, f[j-1]+1) // delete & insert
			if cb != ca {
				mn = min(mn, fj1+1) // change
			} else {
				mn = min(mn, fj1) // matched
			}

			fj1, f[j] = f[j], mn // save f[j] to fj1(j is about to increase), update f[j] to mn
			j++
		}
	}

	return f[len(f)-1]
}

func sel(cond bool, a, b int) int {
	if cond {
		return a
	}
	return b
}
func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func min3(a, b, c int) int {
	if a > b {
		return min(b, c)
	}
	return min(a, c)
}
