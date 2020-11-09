package btree

// Structure - это бинарное дерево
type Structure struct {
	Root *Vertex
	Len  int
}

// Vertex - это вершина графа, хранящая:
//  - "значение" хранимой в дереве сущности (Value)
//  - указатель на хранимую сущность (Item)
//  - указатели на вершины "левее" (L) и "правее" (R)
type Vertex struct {
	Value int
	Item  *Evaluator
	L     *Vertex
	R     *Vertex
}

// Evaluator абстрагирует сущность, которую необходимо поместить в дерево.
// Value() позволяет получить "значение" сущности (для возможности сравнения с другими вершинами).
type Evaluator interface {
	Value() int
}

// Insert реализует вставку сущности в дерево
func (t *Structure) Insert(e Evaluator) {
	// инициализируем новую (добавляемую) вершину
	var newValue = e.Value()
	var ins Vertex
	ins.Value = newValue
	ins.Item = &e

	// если в дереве нет вершин, добавляем новую вершину в корень
	if t.Root == nil {
		t.Len++
		t.Root = &ins
		return
	}

	// если в дереве есть вершины, обходим дерево вглубь (начиная с корня)
	var next = t.Root
	for {
		// если в дереве уже есть вершина с таким же значением - новую не добавляем
		if next.Value == newValue {
			break
		}

		var moveLeft = newValue < next.Value
		var moveRight = !moveLeft

		if moveLeft {
			if next.L == nil {
				next.L = &ins
				t.Len++
				break
			} else {
				next = next.L
				continue
			}
		}

		if moveRight {
			if next.R == nil {
				next.R = &ins
				t.Len++
				break
			} else {
				next = next.R
				continue
			}
		}
	}
}

// Find осуществляет поиск вершины по указанному значению
func (t *Structure) Find(value int) *Vertex {
	var next = t.Root

	for {
		if next.Value == value {
			return next
		}

		var moveLeft = value < next.Value
		var moveRight = !moveLeft

		if moveLeft {
			if next.L == nil {
				return nil
			}
			next = next.L
			continue
		}

		if moveRight {
			if next.R == nil {
				return nil
			}
			next = next.R
			continue
		}
	}
}
