package mbox

type patch struct {
	subject string
	content string
}

// ParseJSON parses the thread for generating patches.
func (t *Thread) ParsePatch() [][]string {
	res := [][]string{}
	patches := findPatch(t)
	for _, p := range patches {
		newPatch := []string{p.subject, p.content}
		res = append(res, newPatch)
	}
	return res
}

func findPatch(t *Thread) []patch {
	patches := []patch{}
	return patches
}
