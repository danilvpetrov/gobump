package pathx

// IsPathMajor reports if the given string complies with path major
// requirements.
//
// See this link for reference: https://go.dev/ref/mod#major-version-suffixes.
func IsPathMajor(in string) bool {
	switch len(in) {
	case 0, 1:
		return false
	case 2:
		return in[0] == 'v' && in[1] >= '2' && in[1] <= '9'
	default:
		if in[0] != 'v' || in[1] < '1' || in[1] > '9' {
			return false
		}

		for _, l := range in[2:] {
			if l < '0' || l > '9' {
				return false
			}
		}
	}

	return true
}
