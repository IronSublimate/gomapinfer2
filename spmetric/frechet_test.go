package spmetric

import (
	"github.com/mitroadmaps/gomapinfer/common"

	"math"
	"testing"
)

func TestComputeFrechetDistance(t *testing.T) {
	f := func(a []common.Point, b []common.Point, expected float64) {
		d := ComputeFrechetDistance(a, b)
		if math.Abs(d-expected) > 0.001 {
			t.Errorf("d(%v, %v) got %v expected %v", a, b, d, expected)
		}
	}

	// single point tests
	f([]common.Point{{1, 1}}, []common.Point{{1, 1}}, 0)
	f([]common.Point{{1, 1}}, []common.Point{{3, 1}}, 2)

	// loop path against single point
	path := []common.Point{
		{0, 0},
		{1, 0},
		{1, 1},
		{0, 1},
		{0, 0},
	}
	f(path, []common.Point{{0, 0}}, math.Sqrt(2))
	f(path, []common.Point{{0.5, 0.5}}, math.Sqrt(2)/2)

	// two paths, one with detour
	directPath := []common.Point{
		{0, 0},
		{1, 0},
		{1, 1},
		{2, 1},
	}
	detourPath := []common.Point{
		{0, 0},
		{1, 0},
		{1, 1},
		{1, 2},
		{2, 2},
		{2, 1},
	}
	f(directPath, directPath, 0)
	f(directPath, detourPath, 1)
	f(detourPath, directPath, 1)

	// two paths, one loops back
	straightPath := []common.Point{
		{0, 0},
		{2, 0},
		{4, 0},
	}
	loopPath := []common.Point{
		{0, 0},
		{4, 0},
		{0, 0},
		{4, 0},
	}
	f(loopPath, loopPath, 0)
	f(straightPath, loopPath, 2)
}

func TestGetClosestPath(t *testing.T) {
	graph := &common.Graph{}
	v11 := graph.AddNode(common.Point{1, 1}) //0
	v12 := graph.AddNode(common.Point{1, 2}) //1
	v31 := graph.AddNode(common.Point{3, 1}) //2
	v32 := graph.AddNode(common.Point{3, 2}) //3
	v51 := graph.AddNode(common.Point{5, 1}) //4
	v52 := graph.AddNode(common.Point{5, 2}) //5
	graph.AddBidirectionalEdge(v11, v12)     //0,1
	graph.AddBidirectionalEdge(v11, v31)     //0,3
	graph.AddBidirectionalEdge(v31, v32)     //2,3
	graph.AddBidirectionalEdge(v31, v51)     //2,4
	graph.AddBidirectionalEdge(v32, v52)     //3,5
	graph.AddBidirectionalEdge(v51, v52)     //4,5

	radius := 10.0

	nps := make(map[int]NodePaths)
	addDis := func(i int, j int) {
		p1 := graph.Nodes[i].Point
		p2 := graph.Nodes[j].Point
		d := p1.Distance(p2)
		nps[i].Distances[j] = d
		nps[j].Distances[i] = d
	}
	for i := 0; i < len(graph.Nodes); i += 1 {
		bp := make(map[int]int)
		dis := make(map[int]float64)
		for j := 0; j < len(graph.Nodes); j += 1 {
			if i == j {
				continue
			}
			bp[j] = i
			dis[j] = math.Inf(1)
		}
		nps[i] = NodePaths{bp, dis}
	}
	addDis(0, 1)
	addDis(0, 3)
	addDis(2, 3)
	addDis(2, 4)
	addDis(3, 5)
	addDis(4, 5)

	f := func(inPath []common.Point, expected []*common.Node, d float64) {
		g2 := NodePathsGraph{
			graph, nps,
		}
		outPath, gotD := GetClosestPath(g2, inPath, radius)
		var outNodes []*common.Node
		outNodes = append(outNodes, outPath.Start.Edge.Src)
		outNodes = append(outNodes, outPath.Path...)
		outNodes = append(outNodes, outPath.End.Edge.Dst)
		if len(outNodes) != len(expected) {
			t.Errorf("GetClosestPath(%v) got %v expected %v", inPath, outNodes, expected)
			return
		}
		var points []common.Point
		for i := range outNodes {
			if outNodes[i] != expected[i] {
				t.Errorf("GetClosestPath(%v) got %v expected %v", inPath, outNodes, expected)
				return
			}
			points = append(points, outNodes[i].Point)
		}
		if math.Abs(gotD-d) > 0.001 {
			t.Errorf("GetClosestPath(%v) got distance %v expected %v", inPath, gotD, d)
			return
		}
	}
	println(f)

	//f([]common.Point{
	//	{1, 2},
	//	{3, 2},
	//}, []*common.Node{
	//	v12,
	//	v11,
	//	v31,
	//	v32,
	//}, 1)
	//
	//f([]common.Point{
	//	{0.8, 2.2},
	//	{0.8, 0.8},
	//	{3.2, 0.8},
	//	{2.8, 2.2},
	//}, []*common.Node{
	//	v12,
	//	v11,
	//	v31,
	//	v32,
	//}, math.Sqrt(2)/5)
}
