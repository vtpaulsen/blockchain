package setupbooter

import "strconv"

// Source:
// https://github.com/philopon/go-toposort

type Graph struct {
	nodes   []string
	outputs map[string]map[string]int
	inputs  map[string]int
}

func NewGraph(cap int) *Graph {
	return &Graph{
		nodes:   make([]string, 0, cap),
		inputs:  make(map[string]int),
		outputs: make(map[string]map[string]int),
	}
}

func (g *Graph) addNode(name string) bool {
	g.nodes = append(g.nodes, name)

	if _, ok := g.outputs[name]; ok {
		return false
	}
	g.outputs[name] = make(map[string]int)
	g.inputs[name] = 0
	return true
}

func (g *Graph) addEdge(from, to string) bool {
	m, ok := g.outputs[from]
	if !ok {
		return false
	}

	m[to] = len(m) + 1
	g.inputs[to]++

	return true
}

func (g *Graph) unsafeRemoveEdge(from, to string) {
	delete(g.outputs[from], to)
	g.inputs[to]--
}

func (g *Graph) toposort() ([]string, bool) {
	L := make([]string, 0, len(g.nodes))
	S := make([]string, 0, len(g.nodes))

	for _, n := range g.nodes {
		if g.inputs[n] == 0 {
			S = append(S, n)
		}
	}

	for len(S) > 0 {
		var n string
		n, S = S[0], S[1:]
		L = append(L, n)

		ms := make([]string, len(g.outputs[n]))
		for m, i := range g.outputs[n] {
			ms[i-1] = m
		}

		for _, m := range ms {
			g.unsafeRemoveEdge(n, m)

			if g.inputs[m] == 0 {
				S = append(S, m)
			}
		}
	}

	N := 0
	for _, v := range g.inputs {
		N += v
	}

	if N > 0 {
		return L, false
	}

	return L, true
}

// Extension of https://github.com/philopon/go-toposort
// with our own topological sorting function that sorts PeerInstances
func topologicalSort(peers []PeerInstance) []PeerInstance {
	graph := NewGraph(len(peers))

	// Add all nodes to the graph
	for i := range peers {
		node := i + 1
		graph.addNode(strconv.Itoa(node))
	}
	// Add all the edges
	for i, v := range peers {
		a := i + 1
		b := v.ShouldConnectTo
		if b != 0 {
			graph.addEdge(strconv.Itoa(a), strconv.Itoa(b))
		}
	}
	sortedStr, _ := graph.toposort()

	// Reverse the returned list
	for i, j := 0, len(sortedStr)-1; i < j; i, j = i+1, j-1 {
		sortedStr[i], sortedStr[j] = sortedStr[j], sortedStr[i]
	}

	// Construct the resulting []PeerInstance by looking at the ordering of the topologically-sorted sortedStr
	var sortedPeers []PeerInstance
	for _, v := range sortedStr {
		i, err := strconv.Atoi(v)
		logOnError(err, "Error while parsing int from string '"+v+"'")

		sortedPeers = append(sortedPeers, peers[i-1])
	}
	return sortedPeers
}
