package indices

// Index Type
type Index interface{}

type Indices []Index

func (idx1 Indices) Equal(idx2 Indices) bool {
	if len(idx1) != len(idx2) {
		return false
	}
	for i := range idx1 {
		if idx1[i] != idx2[i] {
			return false
		}
	}
	return true
}
