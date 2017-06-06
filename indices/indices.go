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

// AND function looks for common indexes between the two arguments of indexes
func (idx1 Indices) And(idx2 Indices) Indices {
	d1 := map[Index]bool{}
	d2 := map[Index]bool{}
	for _, v := range idx1 {
		d1[v] = true
	}
	for _, v := range idx2 {
		d2[v] = true
	}
	ret := Indices{}
	for k := range d1 {
		_, ok := d2[k]
		if ok {
			ret = append(ret, k)
			delete(d1, k)
			delete(d2, k)
		}
	}
	for k := range d2 {
		_, ok := d1[k]
		if ok {
			ret = append(ret, k)
			delete(d1, k)
			delete(d2, k)
		}
	}
	return ret
}

// OR function looks for indexes that are in first argument or second aregument of indexes
func (idx1 Indices) Or(idx2 Indices) Indices {
	d := map[Index]bool{}
	for _, v := range idx1 {
		d[v] = true
	}
	for _, v := range idx2 {
		d[v] = true
	}
	ret := Indices{}
	for k := range d {
		ret = append(ret, k)
	}
	return ret
}
