package chem

func filterEmptyStrings(strings ...string) []string {
	out := make([]string, 0, len(strings))
	for _, str := range strings {
		if str != "" {
			out = append(out, str)
		}
	}
	return out
}
