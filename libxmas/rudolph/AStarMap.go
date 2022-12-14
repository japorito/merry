package rudolph

import (
	"math"

	"github.com/japorito/merry/libxmas/stockings"
)

type MapNode[T any] struct {
	Value          T
	TravelScore    int
	HeuristicScore int
	Self           Coordinate
	Predecessor    Coordinate
}

type IAStarMap[T any] interface {
	TravelCost(src, dst Coordinate) int
	Heuristic(src, dst Coordinate) int
	GetConnectedNodes(src Coordinate) []Coordinate
	Get(coords Coordinate) *MapNode[T]
	SetMap(m [][]*MapNode[T])
	GetMap() [][]*MapNode[T]
}

type AStarMap[T any] struct {
	IAStarMap[T]
	MapNodes [][]*MapNode[T]
}

func (m *AStarMap[T]) Get(coords Coordinate) *MapNode[T] {
	return m.MapNodes[coords.Row][coords.Col]
}

func abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func ManhattanDistance(a, b Coordinate) int {
	return abs(a.Col-b.Col) + abs(a.Col-b.Col)
}

func DifferenceDistance[T stockings.Integral](a T, b T) int {
	return int(b - a)
}

func AStarInit[T any](end Coordinate, initialData [][]T, emptyMap IAStarMap[T]) {
	mapNodes := make([][]*MapNode[T], len(initialData))
	for row := range initialData {
		mapNodes[row] = make([]*MapNode[T], len(initialData[row]))
		for col := range initialData[row] {
			selfCoords := Coordinate{Col: col, Row: row}
			mapNodes[row][col] = &MapNode[T]{
				Value:       initialData[row][col],
				TravelScore: math.MaxInt,
				Self:        selfCoords,
			}
		}
	}

	emptyMap.SetMap(mapNodes)

	// may not be able to calculate heuristic until the map is fully filled out
	for row := range initialData {
		for col := range initialData[row] {
			node := emptyMap.Get(Coordinate{Col: col, Row: row})
			node.HeuristicScore = emptyMap.Heuristic(node.Self, end)
		}
	}
}

func AStarReset[T any](aStarMap IAStarMap[T]) {
	for _, row := range aStarMap.GetMap() {
		for _, node := range row {
			node.TravelScore = math.MaxInt
			node.Predecessor = Coordinate{-1, -1}
		}
	}
}

func AStarRun[T any](start, end Coordinate, aStarMap IAStarMap[T]) {
	startNode := aStarMap.Get(start)
	startNode.TravelScore = 0

	endNode := aStarMap.Get(end)

	processQueue := stockings.NewMinPriorityQueue(32, func(item *MapNode[T]) int {
		return item.TravelScore + item.HeuristicScore
	})
	processQueue.Add(startNode)

	var current *MapNode[T]
	for current = processQueue.GetNext(); current != endNode; current = processQueue.GetNext() {
		nextScore := current.TravelScore + 1

		for _, neighborCoords := range aStarMap.GetConnectedNodes(current.Self) {
			neighbor := aStarMap.Get(neighborCoords)
			if neighbor.TravelScore > nextScore {
				neighbor.TravelScore = nextScore
				neighbor.Predecessor = current.Self

				if processQueue.Has(neighbor) {
					processQueue.TryIncreasePriority(neighbor)
				} else {
					processQueue.Add(neighbor)
				}
			}
		}

		if processQueue.Size() == 0 {
			// dead end! No routes from starting point
			break
		}
	}
}
