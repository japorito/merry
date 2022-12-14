/*
Copyright © 2022 Jacob Saporito
*/
package cmd

import (
	"fmt"
	"math"
	"time"

	"github.com/japorito/merry/libxmas/rudolph"
	"github.com/japorito/merry/libxmas/toybag"
	"github.com/japorito/merry/libxmas/xmas"
	"github.com/spf13/cobra"
)

const (
	sand    = '▒'
	rock    = '█'
	air     = ' '
	origin  = 'O'
	errchar = '⚠'
)

type caveMap struct {
	cave           [][]rune
	minCol, minRow int
}

func (cave *caveMap) checkCoords(coord rudolph.Coordinate) bool {
	row := coord.Row - cave.minRow
	col := coord.Col - cave.minCol

	return row >= 0 && col >= 0 && row < len(cave.cave) && col < len(cave.cave[row])
}

func (cave *caveMap) Free(coord rudolph.Coordinate) bool {
	return cave.checkCoords(coord) && cave.Get(coord) == air
}

func (cave *caveMap) Init(min, max rudolph.Coordinate) {
	cave.minCol, cave.minRow = min.Col, min.Row
	cave.cave = make([][]rune, (max.Row-min.Row)+1)
	for i := range cave.cave {
		cave.cave[i] = make([]rune, (max.Col-min.Col)+1)
		for j := range cave.cave[i] {
			cave.cave[i][j] = air
		}
	}
}

func (cave *caveMap) Print() {
	for _, row := range cave.cave {
		fmt.Println(string(row))
	}
}

func (cave *caveMap) Draw(coords [2]rudolph.Coordinate, char rune) {
	if coords[1].Row < coords[0].Row || coords[1].Col < coords[0].Col {
		coords[0], coords[1] = coords[1], coords[0]
	}
	start, end := coords[0], coords[1]

	for i := start.Row; i <= end.Row; i++ {
		for j := start.Col; j <= end.Col; j++ {
			cave.Put(rudolph.Coordinate{Row: i, Col: j}, char)
		}
	}
}

func (cave *caveMap) TryGet(coord rudolph.Coordinate) (rune, bool) {
	if cave.checkCoords(coord) {
		return cave.cave[coord.Row-cave.minRow][coord.Col-cave.minCol], true
	}

	return errchar, false
}

func (cave *caveMap) Get(coord rudolph.Coordinate) rune {
	return cave.cave[coord.Row-cave.minRow][coord.Col-cave.minCol]
}

func (cave *caveMap) TryPut(coord rudolph.Coordinate, putting rune) bool {
	if cave.checkCoords(coord) {
		if occupant := cave.Get(coord); occupant == air {
			cave.Put(coord, putting)
			return true
		}
	}

	return false
}

func (cave *caveMap) Put(coord rudolph.Coordinate, putting rune) {
	row := coord.Row - cave.minRow
	col := coord.Col - cave.minCol

	cave.cave[row][col] = putting
}

func (cave *caveMap) Drop(dropC rudolph.Coordinate) bool {
	row, col := dropC.Row, dropC.Col

	for i := row; ; {
		nextrow := i + 1
		below, ok := cave.TryGet(rudolph.Coordinate{Row: nextrow, Col: col})
		if !ok {
			return false
		}

		switch below {
		case air:
			//continue dropping loop
			i++
		case sand, rock:
			left, right := rudolph.Coordinate{Row: nextrow, Col: col - 1}, rudolph.Coordinate{Row: nextrow, Col: col + 1}
			if !cave.checkCoords(left) {
				// left out of bounds!
				return false
			} else if cave.Free(rudolph.Coordinate{Row: nextrow, Col: col - 1}) {
				//try left
				col--
			} else if !cave.checkCoords(right) {
				//right out of bounds!
				return false
			} else if cave.Free(rudolph.Coordinate{Row: nextrow, Col: col + 1}) {
				//try right
				col++
			} else {
				return cave.TryPut(rudolph.Coordinate{Row: i, Col: col}, sand)
			}
		}
	}
}

func findCaveDimensions(lines [][2]rudolph.Coordinate, points ...rudolph.Coordinate) (bottomleft, topright rudolph.Coordinate) {
	minCol, minRow := math.MaxInt, math.MaxInt
	maxCol, maxRow := math.MinInt, math.MinInt

	updateMinMax := func(rock rudolph.Coordinate) {
		if rock.Col < minCol {
			minCol = rock.Col
		}
		if rock.Col > maxCol {
			maxCol = rock.Col
		}
		if rock.Row < minRow {
			minRow = rock.Row
		}
		if rock.Row > maxRow {
			maxRow = rock.Row
		}
	}

	for _, pair := range lines {
		updateMinMax(pair[0])
		updateMinMax(pair[1])
	}

	for _, point := range points {
		updateMinMax(point)
	}

	return rudolph.Coordinate{Row: minRow, Col: minCol},
		rudolph.Coordinate{Row: maxRow, Col: maxCol}
}

func parseRockLines(input [][]string) [][2]rudolph.Coordinate {
	rocklines := make([][2]rudolph.Coordinate, 0, len(input)*2)

	for _, rockline := range input {
		for i := 0; i+2 < len(rockline); i += 2 {
			start, _ := rudolph.ParseCoordinate(rockline[i])
			end, _ := rudolph.ParseCoordinate(rockline[i+2])

			rocklines = append(rocklines, [2]rudolph.Coordinate{start, end})
		}
	}

	return rocklines
}

func setupCaveMap(input [][]string, sandOrigin rudolph.Coordinate) *caveMap {
	rocklines := parseRockLines(input)
	min, max := findCaveDimensions(rocklines, sandOrigin)
	cave := caveMap{}
	cave.Init(min, max)

	for _, rockline := range rocklines {
		cave.Draw(rockline, rock)
	}
	cave.TryPut(sandOrigin, origin)

	return &cave
}

func setupCaveMap2(input [][]string, sandOrigin rudolph.Coordinate) *caveMap {
	rocklines := parseRockLines(input)
	min, max := findCaveDimensions(rocklines, sandOrigin)

	height := (max.Row - min.Row) + 2
	floorStart := rudolph.Coordinate{Row: height, Col: min.Col - height}
	floorEnd := rudolph.Coordinate{Row: height, Col: max.Col + height}
	rocklines = append(rocklines, [2]rudolph.Coordinate{floorStart, floorEnd})
	min, max = findCaveDimensions(rocklines, sandOrigin)

	cave := caveMap{}
	cave.Init(min, max)

	for _, rockline := range rocklines {
		cave.Draw(rockline, rock)
	}
	cave.TryPut(sandOrigin, origin)

	return &cave
}

// day14Cmd represents the day14 command
var day14Cmd = &cobra.Command{
	Use:   "day14 path/to/input/file",
	Short: "AoC Day 14",
	Long:  `Advent of Code Day 14: Regolith Reservoir`,
	Run: func(cmd *cobra.Command, args []string) {
		if input := toybag.ReadAsTokenizedStringSlice(args...); input != nil {
			fmt.Printf("%d rock formations read.\n", len(input))

			defer xmas.PrintHolidayMessage(time.Now())

			sandOrigin := rudolph.Coordinate{Row: 0, Col: 500}

			if Parts.Has(1) {
				fmt.Println("Part 1 running...")
				cave := setupCaveMap(input, sandOrigin)

				i := 0
				for cave.Drop(sandOrigin) {
					i++
					if i > 1000 {
						break
					}
				}

				fmt.Printf("**%d** piles of sand came to rest.\n", i)
			}

			if Parts.Has(2) {
				fmt.Println("Part 2 running...")
				cave := setupCaveMap2(input, sandOrigin)

				i := 0
				for cave.Drop(sandOrigin) {
					i++
				}

				fmt.Printf("**%d** total piles of sand came to rest from floor.\n", i+1) //code stops right _before_ origins
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(day14Cmd)
}
