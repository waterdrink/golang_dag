package golang_dag

import (
	"container/list"
	"errors"
)

var (
	ErrCycle           = errors.New("dag: cycle between edges")
	ErrEdgeExists      = errors.New("dag: edge already exists")
	ErrVertexExists    = errors.New("dag: vertex already exists")
	ErrVertexNotExists = errors.New("dag: vertex does not exist")
)

type vertex struct {
	Id       string
	Value    interface{}
	Parents  []*vertex
	Children []*vertex
}

func newVertex(id string, value interface{}) *vertex {
	return &vertex{
		Id:       id,
		Value:    value,
		Parents:  make([]*vertex, 0),
		Children: make([]*vertex, 0),
	}
}

func (v *vertex) removeChild(ChildId string) {
	for i := len(v.Children) - 1; i >= 0; i-- {
		if ChildId == v.Children[i].Id {
			copy(v.Children[i:], v.Children[i+1:])
			v.Children[len(v.Children)-1] = nil
			v.Children = v.Children[:len(v.Children)-1]
		}
	}
}

func (v *vertex) removeParent(parentId string) {
	for i := len(v.Parents) - 1; i >= 0; i-- {
		if parentId == v.Parents[i].Id {
			copy(v.Parents[i:], v.Parents[i+1:])
			v.Parents[len(v.Parents)-1] = nil
			v.Parents = v.Parents[:len(v.Parents)-1]
		}
	}
}

func (v *vertex) isEqual(v2 *vertex) bool {
	if v.Id != v2.Id {
		return false
	}
	if len(v.Parents) != len(v2.Parents) {
		return false
	}
	for i := range v.Parents {
		if v.Parents[i].Id != v2.Parents[i].Id {
			return false
		}
	}
	if len(v.Children) != len(v2.Children) {
		return false
	}
	for i := range v.Children {
		if v.Children[i].Id != v2.Children[i].Id {
			return false
		}
	}
	return true
}

// A directed acyclic graph implementation
// Not thread safe, caller is responsible to ensure mutex
type DAG struct {
	Vertexes map[string]*vertex
}

func NewDAG() *DAG {
	return &DAG{
		Vertexes: make(map[string]*vertex),
	}
}

func (dag *DAG) AddVertex(vertexId string, value interface{}) error {
	if _, ok := dag.Vertexes[vertexId]; ok {
		return ErrVertexExists
	}

	dag.Vertexes[vertexId] = newVertex(vertexId, value)
	return nil
}

func (dag *DAG) RemoveVertex(vertexId string) {
	vertex, ok := dag.Vertexes[vertexId]
	if !ok {
		return
	}

	for _, parent := range vertex.Parents {
		parent.removeChild(vertexId)
	}
	for _, child := range vertex.Children {
		child.removeParent(vertexId)
	}
	delete(dag.Vertexes, vertexId)
}

func (dag *DAG) AddEdge(fromVertexId, toVertexId string) error {
	if fromVertexId == toVertexId {
		return ErrCycle
	}
	var from, to *vertex
	var ok bool

	if from, ok = dag.Vertexes[fromVertexId]; !ok {
		return ErrVertexNotExists
	}

	if to, ok = dag.Vertexes[toVertexId]; !ok {
		return ErrVertexNotExists
	}

	for _, childVertex := range from.Children {
		if childVertex == to {
			return ErrEdgeExists
		}
	}

	if dag.DepthFirstSearch(toVertexId, fromVertexId) {
		return ErrCycle
	}

	from.Children = append(from.Children, to)
	to.Parents = append(to.Parents, from)
	return nil
}

func (dag *DAG) RemoveEdge(fromVertexId, toVertexId string) error {
	var from, to *vertex
	var ok bool

	if from, ok = dag.Vertexes[fromVertexId]; !ok {
		return ErrVertexNotExists
	}

	if to, ok = dag.Vertexes[toVertexId]; !ok {
		return ErrVertexNotExists
	}

	to.removeParent(fromVertexId)
	from.removeChild(toVertexId)
	return nil
}

func (dag *DAG) EdgeExists(fromVertexId, toVertexId string) (bool, error) {
	var from, to *vertex
	var ok bool

	if from, ok = dag.Vertexes[fromVertexId]; !ok {
		return false, ErrVertexNotExists
	}

	if to, ok = dag.Vertexes[toVertexId]; !ok {
		return false, ErrVertexNotExists
	}

	// quick return
	if len(to.Parents) == 0 {
		return false, nil
	}

	for _, childVertex := range from.Children {
		if childVertex == to {
			return true, nil
		}
	}

	return false, nil
}

func (dag *DAG) GetVertex(id string) *vertex {
	if v, ok := dag.Vertexes[id]; ok {
		return v
	}

	return nil
}

func (dag *DAG) DepthFirstSearch(fromVertexId, toVertexId string) bool {
	found := map[string]bool{}
	dag.dfs(found, fromVertexId)
	return found[toVertexId]
}

func (dag *DAG) dfs(found map[string]bool, vertexId string) {
	vertex, ok := dag.Vertexes[vertexId]
	if !ok {
		return
	}
	for _, child := range vertex.Children {
		if !found[child.Id] {
			found[child.Id] = true
			dag.dfs(found, child.Id)
		}
	}
}

func (dag *DAG) IsEqual(dag2 *DAG) bool {
	if len(dag.Vertexes) != len(dag2.Vertexes) {
		return false
	}
	for vId, v := range dag.Vertexes {
		v2, ok := dag2.Vertexes[vId]
		if !ok {
			return false
		}
		if !v.isEqual(v2) {
			return false
		}
	}
	return true
}

// shallow Copy
func (dag *DAG) Copy() *DAG {
	new := NewDAG()
	for _, v := range dag.Vertexes {
		new.Vertexes[v.Id] = &vertex{
			Id:    v.Id,
			Value: v.Value,
		}
	}
	for _, v := range dag.Vertexes {
		for _, child := range v.Children {
			new.AddEdge(v.Id, child.Id)
		}
	}
	return new
}

func (dag *DAG) TopologicalSort() []*vertex {
	copy := dag.Copy()

	sort := []*vertex{}
	for {
		for _, v := range copy.Vertexes {
			if len(v.Parents) != 0 {
				continue
			}
			for _, child := range v.Children {
				child.removeParent(v.Id)
			}
			delete(copy.Vertexes, v.Id)
			sort = append(sort, v)
		}
		if len(copy.Vertexes) == 0 {
			break
		}
	}

	return sort
}

func (dag *DAG) TopologicalSortStable() []*vertex {
	copy := dag.Copy()
	noParentsVertexes := newSortedVertexes()
	length := len(copy.Vertexes)
	sort := make([]*vertex, 0, length)
	if length == 0 {
		return sort
	}

	for {
		for _, v := range copy.Vertexes {
			if len(v.Parents) != 0 {
				continue
			}
			noParentsVertexes.add(v)
			delete(copy.Vertexes, v.Id)
		}
		firstNoParentsVertex := noParentsVertexes.popFront()
		sort = append(sort, firstNoParentsVertex)
		if len(sort) == length {
			break
		}
		for _, child := range firstNoParentsVertex.Children {
			child.removeParent(firstNoParentsVertex.Id)
		}
	}

	return sort
}

type sortedVertexes struct {
	*list.List
}

func newSortedVertexes() *sortedVertexes {
	l := list.New()
	return &sortedVertexes{l}
}

func (s *sortedVertexes) add(v *vertex) {
	for e := s.Front(); e != nil; e = e.Next() {
		if v.Id < e.Value.(*vertex).Id {
			s.InsertBefore(v, e)
			return
		}
	}
	s.PushBack(v)
}

func (s *sortedVertexes) popFront() *vertex {
	e := s.Front()
	if nil == e {
		return nil
	}
	s.Remove(e)
	return e.Value.(*vertex)
}

func (dag *DAG) Print() (str string) {
	for _, v := range dag.Vertexes {
		if len(v.Parents) == 0 {
			str = str + dag.print(v, "") + "\n"
		}
	}
	return str
}

func (dag *DAG) print(root *vertex, prefix string) string {
	str := prefix + root.Id + "\n"
	for i, child := range root.Children {
		// If last iteration, don't add a pipe character
		if i == len(root.Children)-1 {
			str = str + dag.print(child, prefix+"    ")
		} else {
			str = str + dag.print(child, prefix+"    |")
		}
	}
	return str
}
