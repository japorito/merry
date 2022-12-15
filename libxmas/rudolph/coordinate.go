package rudolph

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Coordinate struct {
	Col, Row int
}

func (c *Coordinate) Equals(other Coordinate) bool {
	return c.Col == other.Col && c.Row == other.Row
}

func (c *Coordinate) String() string {
	return fmt.Sprintf("%d,%d", c.Col, c.Row)
}

var ErrorCoordinate = Coordinate{math.MinInt, math.MinInt}

func ParseCoordinate(coord string) (Coordinate, error) {
	parts := strings.Split(coord, ",")
	col, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return ErrorCoordinate, err
	}

	row, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return ErrorCoordinate, err
	}

	return Coordinate{Row: row, Col: col}, nil
}
