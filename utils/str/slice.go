package str

func SliceDiff(s1, s2 []string) []string {
	if len(s1) == 0 || len(s2) == 0 {
		return s1
	}
	m := make(map[string]bool)
	for _, v := range s2 {
		m[v] = true
	}
	var res []string
	for _, v := range s1 {
		if _, found := m[v]; !found {
			res = append(res, v)
		}
	}
	return res
}

func SliceUniq(s []string) []string {
	if len(s) == 0 {
		return s
	}
	m := make(map[string]bool, len(s))
	res := make([]string, 0, len(m))
	for _, v := range s {
		if _, found := m[v]; found {
			continue
		} else {
			res = append(res, v)
			m[v] = true
		}
	}
	return res
}

func SliceFilter(s []string, predictFunc func(string) bool) []string {
	if len(s) == 0 {
		return s
	}
	res := make([]string, 0, len(s))
	for _, v := range s {
		if predictFunc(v) {
			res = append(res, v)
		}
	}
	return res
}
