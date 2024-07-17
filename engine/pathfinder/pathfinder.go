package pathfinder

import (
	"fmt"
	"strings"

	fibHeap "github.com/starwander/GoFibonacciHeap"
	"github.com/victorbetoni/go-streams/streams"
)

type DistancesMap map[string]float64
type PreviousMap map[string]string
type Path []string

type Node struct {
	Territory    string
	Holding      string
	Parent       *Node
	Distance     float64
	AllyTax      float64
	Tax          float64
	Cheapest     bool
	BorderClosed bool
	Connections  []string
	Allies       []string
}

func (n *Node) Tag() string {
	return n.Territory
}

func (n *Node) Key() float64 {
	return n.Distance
}

type Pathfinder struct {
	Root      Node
	Target    Node
	Nodes     []*Node
	nodes     map[string]*Node
	visited   map[string]bool
	previous  PreviousMap
	distances DistancesMap
	path      []string
}

func (p *Pathfinder) djikstra() {

	p.distances = make(map[string]float64)
	p.visited = make(map[string]bool)
	p.previous = make(map[string]string)
	p.nodes = make(map[string]*Node)

	for _, node := range p.Nodes {
		p.distances[node.Tag()] = node.Distance
	}

	p.distances[p.Root.Territory] = 0

	heap := fibHeap.NewFibHeap()
	heap.Insert(p.Root, 0)

	for heap.Num() != 0 {

		minTag, _ := heap.ExtractMin()
		tag := fmt.Sprint(minTag)
		if _, ok := p.visited[tag]; !ok {
			continue
		}

		from := p.nodes[tag]
		for _, conn := range from.Connections {
			if _, ok := p.visited[conn]; !ok {
				weight := 1.0
				territory := p.nodes[conn]
				ally := streams.StreamOf[string](territory.Allies...).AnyMatch(func(val string) bool {
					return strings.EqualFold(val, p.Root.Holding)
				})

				if p.Root.Holding != territory.Holding {
					if territory.BorderClosed {
						weight = 99999
					} else if territory.Cheapest {
						if ally {
							weight += territory.AllyTax
						} else {
							weight += territory.Tax
						}
					}
				}

				to := p.nodes[conn]
				newDist := p.distances[from.Territory]
				if newDist < p.distances[conn] {
					to.Distance = from.Distance + weight
					to.Parent = from
					p.distances[conn] = newDist
					p.previous[conn] = from.Territory
					heap.Insert(to, to.Distance)
				}
			}
		}

	}
}

func (p *Pathfinder) Route() ([]string, float64, float64, bool) {

	p.djikstra()

	path := make([]string, 0)
	possible := true
	composedTax, tax := 0.0, 0.0
	current := p.previous[p.Target.Territory]

	for current != "" && !strings.EqualFold(current, p.Root.Territory) {
		currentTerr := p.nodes[current]
		if !strings.EqualFold(currentTerr.Holding, p.Root.Holding) {

			ally := streams.StreamOf[string](currentTerr.Allies...).AnyMatch(func(val string) bool {
				return strings.EqualFold(val, p.Root.Holding)
			})

			if currentTerr.BorderClosed {
				possible = false
			} else {
				if ally {
					tax += currentTerr.AllyTax
					composedTax *= (1 - (currentTerr.AllyTax / 100))
				} else {
					tax += currentTerr.Tax
					composedTax *= (1 - (currentTerr.Tax / 100))
				}
			}
		}

		path = append(path, current)
		current = p.previous[current]
	}

	reversed := streams.StreamOf[string](path...).Reversed()

	return reversed, tax, composedTax, possible
}
