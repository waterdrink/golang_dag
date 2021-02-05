package golang_dag

import (
	"fmt"
	"testing"
)

func TestDAG_TopologicalSort(t *testing.T) {
	dag := NewDAG()
	dag.AddVertex("v-1", 1)
	dag.AddVertex("v-2", 2)
	dag.AddVertex("v-3", 3)
	dag.AddVertex("v-4", 4)
	dag.AddEdge("v-1", "v-2")
	dag.AddEdge("v-1", "v-3")
	dag.AddEdge("v-1", "v-4")
	dag.AddEdge("v-2", "v-3")
	dag.AddEdge("v-3", "v-4")
	vv := dag.TopologicalSort()
	if len(vv) != 4 {
		t.Fatalf("wrong dag vertex length %v", len(vv))
	}
	actual := ""
	for _, v := range vv {
		actual = actual + " " + v.Id
	}
	expect := " v-1 v-2 v-3 v-4"
	if expect != actual {
		t.Fatalf("wrong dag sort %v", actual)
	}
	fmt.Println(dag.Print())
}

func TestDAG_TopologicalSort_Stable(t *testing.T) {
	dag := NewDAG()
	dag.AddVertex("v-1", 1)
	dag.AddVertex("v-2", 2)
	dag.AddVertex("v-3", 3)
	dag.AddVertex("v-4", 4)
	dag.AddVertex("v-5", 5)
	dag.AddVertex("v-6", 6)
	dag.AddEdge("v-4", "v-2")
	dag.AddEdge("v-2", "v-1")
	vv := dag.TopologicalSortStable()
	if len(vv) != 6 {
		t.Fatalf("wrong dag vertex length %v", len(vv))
	}
	actual := ""
	for _, v := range vv {
		actual = actual + " " + v.Id
	}
	expect := " v-3 v-4 v-2 v-1 v-5 v-6"
	if expect != actual {
		t.Fatalf("wrong dag sort %v", actual)
	}
	fmt.Println(dag.Print())
}

func BenchmarkDAG_TopologicalSort(b *testing.B) {
	for n := 0; n < b.N; n++ {
		dag := NewDAG()
		dag.AddVertex("v-1", 1)
		dag.AddVertex("v-2", 2)
		dag.AddVertex("v-3", 3)
		dag.AddVertex("v-4", 4)
		dag.AddVertex("v-5", 5)
		dag.AddVertex("v-6", 6)
		dag.AddVertex("v-7", 7)
		dag.AddVertex("v-8", 8)
		dag.AddVertex("v-9", 9)
		dag.AddVertex("v-10", 10)
		dag.TopologicalSort()
	}
}

func BenchmarkDAG_TopologicalSort_Stable(b *testing.B) {
	for n := 0; n < b.N; n++ {
		dag := NewDAG()
		dag.AddVertex("v-1", 1)
		dag.AddVertex("v-2", 2)
		dag.AddVertex("v-3", 3)
		dag.AddVertex("v-4", 4)
		dag.AddVertex("v-5", 5)
		dag.AddVertex("v-6", 6)
		dag.AddVertex("v-7", 7)
		dag.AddVertex("v-8", 8)
		dag.AddVertex("v-9", 9)
		dag.AddVertex("v-10", 10)
		dag.TopologicalSortStable()
	}
}
