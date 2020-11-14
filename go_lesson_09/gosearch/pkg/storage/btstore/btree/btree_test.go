package btree

import (
	"testing"
)

// Тип, реализующий интерфейс Evaluator
// (в дереве хранятся ссылки на его объекты)
type Elem struct {
	val int
}

// Реализация интерфейса Evaluator
func (e *Elem) Value() int {
	return e.val
}

func genTree() *Structure {
	var tree = Structure{}

	var elems = []Elem{
		{val: 8}, {val: 8},
		{val: 3}, {val: 10},
		{val: 1}, {val: 6}, {val: 14},
		{val: 4}, {val: 7}, {val: 13},
	}

	for _, elem := range elems {
		tree.Insert(&elem)
	}

	return &tree
}

func TestStructure_Insert(t *testing.T) {
	tree := genTree()

	b1 := tree.Root.L.L.Value
	b2 := tree.Root.L.R.L.Value
	b3 := tree.Root.L.R.R.Value
	b4 := tree.Root.R.R.L.Value

	var got = [4]int{b1, b2, b3, b4}
	var want = [4]int{1, 4, 7, 13}
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	gotLen := tree.Len
	wantLen := 9
	if gotLen != wantLen {
		t.Fatalf("got len=%d, want len=%d", gotLen, wantLen)
	}
}

func TestStructure_Find(t *testing.T) {
	tree := genTree()

	got := tree.Find(4).Value
	var want int = 4
	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	tests := []struct {
		name string
		val  int
		want *Vertex
	}{
		{
			name: "No such value (going left)",
			val:  0,
			want: nil,
		},
		{
			name: "No such value (going right)",
			val:  777,
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tree.Find(tt.val); got != nil {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
