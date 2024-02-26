package main

import (
	"fmt"

	"github.com/JonasBernard/min-cost-max-flow/graph"
	"github.com/JonasBernard/min-cost-max-flow/network"
	"github.com/JonasBernard/min-cost-max-flow/util"
)

func __main() {
	children := []Child{
		{"Mia", "Tanz", "Tuch", "Akro"},
		{"Noah", "Tanz", "Akro", "Tuch"},
		{"Jonas", "Tanz", "Akro", "Tuch"},
		{"Max", "Tuch", "Tanz", "Akro"},
		{"Johanna", "Akro", "Tanz", "Tuch"},
		{"Sarah", "Tanz", "Jonglage", "Akro"},
	}

	tanz := Workshop{"Tanz", 1}
	tuch := Workshop{"Tuch", 1}
	akro := Workshop{"Akro", 2}
	jonglage := Workshop{"Jonglage", 2}

	getWorkshop := func(Name string) Workshop {
		switch s := Name; s {
		case "Tanz":
			return tanz
		case "Tuch":
			return tuch
		case "Akro":
			return akro
		case "Jonglage":
			return jonglage
		default:
			return jonglage
		}
	}

	workshops := []Workshop{tanz, tuch, akro, jonglage}

	// modelling as a network problem

	childNodes := util.MapSlice(children, func(c *Child) graph.Vertex[MatchNode[Child, WorkshopSlot]] {
		return graph.V(&MatchNode[Child, WorkshopSlot]{
			Name:      c.Name,
			IsRight:   false,
			LeftValue: *c,
		})
	})

	workshopNodes := util.FlatMapSlice(workshops, func(w *Workshop) []graph.Vertex[MatchNode[Child, WorkshopSlot]] {
		workshopSlots := make([]graph.Vertex[MatchNode[Child, WorkshopSlot]], 0, w.Capacity)
		for i := 0; i < w.Capacity; i++ {
			theSlot := WorkshopSlot{
				Workshop: *w,
				Nr:       i + 1,
			}
			workshopSlots = append(workshopSlots, graph.V(&MatchNode[Child, WorkshopSlot]{
				Name:       fmt.Sprintf("%v (%v)", w.Name, i+1),
				IsRight:    true,
				RightValue: theSlot,
			}))
		}
		return workshopSlots
	})

	// the capacity cannot be known beforehand because it depends on what the choices are
	allEdges := make([]*graph.WeightedDirectedEdge[MatchNode[Child, WorkshopSlot]], 0, 4*len(childNodes)+len(workshopNodes))

	source := graph.V(&MatchNode[Child, WorkshopSlot]{
		Name: "S", IsSource: true,
	})

	sink := graph.V(&MatchNode[Child, WorkshopSlot]{
		Name: "T", IsSink: true,
	})

	for _, childNode := range childNodes {
		allEdges = append(allEdges, &graph.WeightedDirectedEdge[MatchNode[Child, WorkshopSlot]]{
			VertexFrom: source,
			VertexTo:   childNode,
			Weight:     1,
			Capacity:   1,
		})

		w1 := getWorkshop(childNode.Node.LeftValue.W1)
		w2 := getWorkshop(childNode.Node.LeftValue.W2)
		w3 := getWorkshop(childNode.Node.LeftValue.W3)

		correspondingNodesW1 := util.FilterSlice(workshopNodes, func(node graph.Vertex[MatchNode[Child, WorkshopSlot]]) bool {
			return node.Node.RightValue.Workshop == w1
		})

		correspondingNodesW2 := util.FilterSlice(workshopNodes, func(node graph.Vertex[MatchNode[Child, WorkshopSlot]]) bool {
			return node.Node.RightValue.Workshop == w2
		})

		correspondingNodesW3 := util.FilterSlice(workshopNodes, func(node graph.Vertex[MatchNode[Child, WorkshopSlot]]) bool {
			return node.Node.RightValue.Workshop == w3
		})

		allOtherSlots := util.FilterSlice(workshopNodes, func(node graph.Vertex[MatchNode[Child, WorkshopSlot]]) bool {
			w := node.Node.RightValue.Workshop
			return w != w1 && w != w2 && w != w3
		})

		for _, slot := range correspondingNodesW1 {
			newEdge := graph.WeightedDirectedEdge[MatchNode[Child, WorkshopSlot]]{
				VertexFrom: childNode,
				VertexTo:   slot,
				Weight:     1,
				Capacity:   1,
			}
			allEdges = append(allEdges, &newEdge)
		}

		for _, slot := range correspondingNodesW2 {
			newEdge := graph.WeightedDirectedEdge[MatchNode[Child, WorkshopSlot]]{
				VertexFrom: childNode,
				VertexTo:   slot,
				Weight:     2,
				Capacity:   1,
			}
			allEdges = append(allEdges, &newEdge)
		}

		for _, slot := range correspondingNodesW3 {
			newEdge := graph.WeightedDirectedEdge[MatchNode[Child, WorkshopSlot]]{
				VertexFrom: childNode,
				VertexTo:   slot,
				Weight:     4,
				Capacity:   1,
			}
			allEdges = append(allEdges, &newEdge)
		}

		for _, slot := range allOtherSlots {
			newEdge := graph.WeightedDirectedEdge[MatchNode[Child, WorkshopSlot]]{
				VertexFrom: childNode,
				VertexTo:   slot,
				Weight:     10,
				Capacity:   1,
			}
			allEdges = append(allEdges, &newEdge)
		}
	}

	for _, slot := range workshopNodes {
		allEdges = append(allEdges, &graph.WeightedDirectedEdge[MatchNode[Child, WorkshopSlot]]{
			VertexFrom: slot,
			VertexTo:   sink,
			Weight:     0,
			Capacity:   1,
		})
	}

	allVertices := append(childNodes, workshopNodes...)
	allVertices = append(allVertices, source)
	allVertices = append(allVertices, sink)

	network := network.WeigthedNetwork[MatchNode[Child, WorkshopSlot]]{
		WeigthedDirectedGraph: graph.WeigthedDirectedGraph[MatchNode[Child, WorkshopSlot]]{
			Vertices: allVertices,
			Edges:    allEdges,
		},
		Source: source,
		Sink:   sink,
	}

	// fmt.Print(network)

	flow := network.MinCostMaxFlow()

	// util.PrintMap(flow)

	// interpret flow again

	matchingEdges := util.FilterMapBoth(flow, func(wde *graph.WeightedDirectedEdge[MatchNode[Child, WorkshopSlot]], f float64) bool {
		if f == 0 {
			return false
		}
		if wde.VertexFrom.Node == source.Node || wde.VertexTo.Node == sink.Node {
			return false
		}
		return true
	})

	if len(matchingEdges) != len(children) {
		fmt.Printf("\n\nThere is no perfect solution (beacuse the capacities of the workshop do not suffice).\n")
		fmt.Printf("Here is the best it can get (missing %v child(ren)):\n", len(children)-len(matchingEdges))
	} else {
		print("\n\nFound solution:\n")
	}

	for m := range matchingEdges {
		fmt.Printf("Assing %v to slot %v of workshop %v\n",
			m.VertexFrom.Node.LeftValue.Name,
			m.VertexTo.Node.RightValue.Nr,
			m.VertexTo.Node.RightValue.Workshop.Name)
	}

	for _, w := range workshops {
		fmt.Println()
		fmt.Printf("Kids of workshop %v (max %v):\n", w.Name, w.Capacity)
		for e := range util.FilterMapBoth(matchingEdges, func(wde *graph.WeightedDirectedEdge[MatchNode[Child, WorkshopSlot]], f float64) bool {
			return wde.VertexTo.Node.RightValue.Workshop == w
		}) {
			fmt.Printf("%v\n", e.VertexFrom.Node.LeftValue.Name)
		}
	}
}
