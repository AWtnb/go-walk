package dir

import (
	"strings"
)

func getDepth(path string) int {
	return strings.Count(strings.TrimSuffix(path, SEP), SEP)
}

type DirMember struct {
	rootDepth int
	maxDepth  int
}

func (dm *DirMember) Init(path string, depth int) {
	dm.rootDepth = getDepth(path)
	dm.maxDepth = depth
}

func (dm DirMember) MaxDepth() int {
	return dm.maxDepth
}

func (dm DirMember) IsSkippableDepth(path string) bool {
	return 0 < dm.maxDepth && dm.maxDepth < getDepth(path)-dm.rootDepth
}

func (dm DirMember) FilterByDepth(paths []string) (filteredPaths []string) {
	if dm.maxDepth < 0 {
		filteredPaths = paths
		return
	}
	for i := 0; i < len(paths); i++ {
		p := paths[i]
		if dm.IsSkippableDepth(p) {
			continue
		}
		filteredPaths = append(filteredPaths, p)
	}
	return
}
