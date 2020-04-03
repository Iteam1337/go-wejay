package utils

type Pair struct {
	ID   string
	Size int
}
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Size < p[j].Size }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
