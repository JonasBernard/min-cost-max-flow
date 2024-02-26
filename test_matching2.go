package main

import (
	"fmt"

	"github.com/JonasBernard/min-cost-max-flow/matching"
	"github.com/JonasBernard/min-cost-max-flow/util"
)

type MatchNode[L any, R any] struct {
	Name       string
	IsRight    bool
	IsSource   bool
	IsSink     bool
	LeftValue  L
	RightValue R
}

func (n MatchNode[L, R]) String() string {
	return n.Name
}

type Child struct {
	Name string
	W1   string
	W2   string
	W3   string
}

func (c Child) String() string {
	return c.Name
}

type WorkshopSlot struct {
	Workshop Workshop
	Nr       int
}

func (w WorkshopSlot) String() string {
	return fmt.Sprintf("%v (slot %v)", w.Workshop.Name, w.Nr)
}

type Workshop struct {
	Name     string
	Capacity int
}

func main() {
	children := []Child{
		{"Mia", "Jonglage", "Tuch", "Akro"},
		{"Noah", "Jonglage", "Tuch", "Akro"},
		{"Jonas", "Jonglage", "Tuch", "Akro"},
		{"Max", "Tuch", "Tanz", "Akro"},
		{"Johanna", "Akro", "Tanz", "Tuch"},
		{"Sarah", "Tanz", "Jonglage", "Akro"},
		{"Felix", "Jonglage", "Tanz", "Tuch"},
	}

	tanz := Workshop{"Tanz", 2}
	tuch := Workshop{"Tuch", 2}
	akro := Workshop{"Akro", 2}
	jonglage := Workshop{"Jonglage", 1}

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

	workshops := []Workshop{jonglage, akro, tuch, tanz}

	workshopSlots := util.FlatMapSlice(workshops, func(w *Workshop) []WorkshopSlot {
		slots := make([]WorkshopSlot, 0, w.Capacity)
		for i := 0; i < w.Capacity; i++ {
			slot := WorkshopSlot{
				Workshop: *w,
				Nr:       i + 1,
			}
			slots = append(slots, slot)
		}
		return slots
	})

	matchingProblem := matching.MatchingProblem[Child, WorkshopSlot]{
		Lefts:  children,
		Rights: workshopSlots,
	}

	matchingEdgesArray, err := matchingProblem.SolveMany(5, func(c Child, w WorkshopSlot) (connect bool, weight float64) {
		w1 := getWorkshop(c.W1)
		w2 := getWorkshop(c.W2)
		w3 := getWorkshop(c.W3)

		if w1 == w.Workshop {
			return true, 1
		}
		if w2 == w.Workshop {
			return true, 2
		}
		if w3 == w.Workshop {
			return true, 5
		}

		return true, 10
	})

	if err != nil {
		fmt.Printf("%v\n", err)
	}

	// matchingEdgesArray := [][]matching.MatchingEdge[Child, WorkshopSlot]{matchingEdges}

	for i, matchingEdges := range matchingEdgesArray {
		fmt.Printf("---\nSolution %v\n", i)
		for _, m := range matchingEdges {
			fmt.Printf("Assing %v to slot %v of workshop %v\n",
				m.Left.Name,
				m.Right.Nr,
				m.Right.Workshop.Name)
		}

		for _, w := range workshops {
			fmt.Println()
			fmt.Printf("Kids of workshop %v (max %v):\n", w.Name, w.Capacity)
			for _, e := range util.FilterSlice(matchingEdges, func(e matching.MatchingEdge[Child, WorkshopSlot]) bool {
				return e.Right.Workshop == w
			}) {
				fmt.Printf("%v\n", e.Left.Name)
			}
		}
	}
}
