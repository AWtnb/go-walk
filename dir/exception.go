package dir

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const SEP = string(os.PathSeparator)

func getCommonRoot(paths []string) string {
	if len(paths) < 2 {
		return paths[0]
	}
	slices.SortFunc(paths, func(a, b string) int {
		la := len(strings.Split(a, SEP))
		lb := len(strings.Split(b, SEP))
		if la < lb {
			return -1
		}
		if lb < la {
			return 1
		}
		return 0
	})
	shortest := paths[0]
	second := paths[1]
	if strings.HasPrefix(second, shortest+SEP) {
		return shortest
	}
	return filepath.Dir(shortest)
}

type WalkException struct {
	names []string
}

func (wex *WalkException) SetName(s string) {
	wex.names = append(wex.names, strings.TrimSpace(s))
}

func (wex *WalkException) SetNames(s string, sep string) {
	if len(s) < 1 {
		return
	}
	for _, elem := range strings.Split(s, sep) {
		wex.SetName(elem)
	}
}

func (wex WalkException) Contains(name string) bool {
	for _, n := range wex.names {
		if n == name {
			return true
		}
	}
	return false
}

func (wex WalkException) isSkippable(path string, root string) bool {
	rel := strings.TrimPrefix(strings.TrimPrefix(path, root), SEP)
	if strings.HasPrefix(rel, ".") || strings.Contains(rel, SEP+".") {
		return true
	}
	for _, n := range wex.names {
		if strings.HasPrefix(rel, n) || strings.Contains(rel, SEP+n+SEP) || strings.HasSuffix(rel, n) {
			return true
		}
	}
	return false
}

func (wex WalkException) Filter(paths []string) (filtered []string) {
	if len(wex.names) < 1 {
		filtered = paths
		return
	}
	root := getCommonRoot(paths)
	for i := 0; i < len(paths); i++ {
		p := paths[i]
		if wex.isSkippable(p, root) {
			continue
		}
		filtered = append(filtered, p)
	}
	return
}
