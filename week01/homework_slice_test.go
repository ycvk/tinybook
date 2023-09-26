package week01

import "testing"

func TestIntDel(t *testing.T) {
	ints := []int{1, 2, 3, 4, 5}
	t.Logf("ints:%v, len:%d, cap:%d", ints, len(ints), cap(ints))
	SliceDelIdx(&ints, 2)
	t.Logf("ints:%v, len:%d, cap:%d", ints, len(ints), cap(ints))
}

func TestStringDel(t *testing.T) {
	strings := []string{"a", "b", "c", "d", "e"}
	t.Logf("strings:%v, len:%d, cap:%d", strings, len(strings), cap(strings))
	SliceDelIdx(&strings, 2)
	t.Logf("strings:%v, len:%d, cap:%d", strings, len(strings), cap(strings))
}

func TestSliceReduce(t *testing.T) {
	i := make([]int, 5, 257)
	i[0] = 1
	i[1] = 2
	i[2] = 3
	t.Logf("i:%v, len:%d, cap:%d", i, len(i), cap(i))
	SliceDelIdx(&i, 1)
	t.Logf("i:%v, len:%d, cap:%d", i, len(i), cap(i))

	i = append(i, 4)
	SliceDelIdx(&i, 1)
	t.Logf("i:%v, len:%d, cap:%d", i, len(i), cap(i))

	SliceDelIdx(&i, 0)
	t.Logf("i:%v, len:%d, cap:%d", i, len(i), cap(i))
}

func TestSliceReduce2(t *testing.T) {
	ints := make([]int, 1, 1)
	ints[0] = 1
	t.Logf("ints:%v, len:%d, cap:%d", ints, len(ints), cap(ints))
	SliceDelIdx(&ints, 0)
	t.Logf("ints:%v, len:%d, cap:%d", ints, len(ints), cap(ints))
	ints = append(ints, 2)
	t.Logf("ints:%v, len:%d, cap:%d", ints, len(ints), cap(ints))
}
