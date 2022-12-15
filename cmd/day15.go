/*
Copyright © 2022 Jacob Saporito
*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/japorito/merry/libxmas/rudolph"
	"github.com/japorito/merry/libxmas/toybag"
	"github.com/japorito/merry/libxmas/xmas"
	"github.com/spf13/cobra"
)

type coordinatePair struct {
	sensor, beacon rudolph.Coordinate
	distance       int
}

func parseCoordinateNumber(entry string) int {
	num, _ := strconv.Atoi(strings.Trim(entry, "xy=,:"))
	return num
}

func parseSensorCoordinate(col, row string) rudolph.Coordinate {
	return rudolph.Coordinate{
		Row: parseCoordinateNumber(row),
		Col: parseCoordinateNumber(col),
	}
}

func parseSensorData(input [][]string) []coordinatePair {
	output := make([]coordinatePair, len(input))

	for i, sensor := range input {
		output[i] = coordinatePair{
			sensor: parseSensorCoordinate(sensor[2], sensor[3]),
			beacon: parseSensorCoordinate(sensor[8], sensor[9]),
		}

		output[i].distance = rudolph.ManhattanDistance(output[i].sensor, output[i].beacon)
	}

	return output
}

func findEmergencyBeacon(sensorData []coordinatePair, min, max rudolph.Coordinate) rudolph.Coordinate {
	for i := min.Row; i <= max.Row; i++ {
	ColLoop:
		for j := min.Col; j <= max.Col; j++ {
			for _, pair := range sensorData {
				pt := rudolph.Coordinate{Row: i, Col: j}
				dPt := rudolph.ManhattanDistance(pt, pair.sensor)

				if dPt <= pair.distance {
					// skip all columns eliminated by this sensor
					j = pair.sensor.Col + (pair.distance - abs(pair.sensor.Row-i))
					continue ColLoop
				}
			}
			return rudolph.Coordinate{Row: i, Col: j}
		}
	}

	return rudolph.ErrorCoordinate
}

func countExcludedZones(sensorData []coordinatePair, line int) int {
	exclusions := make(map[int]bool)
	for _, pair := range sensorData {
		dBeacon := rudolph.ManhattanDistance(pair.sensor, pair.beacon)

		nearestPtOnLine := rudolph.Coordinate{Row: line, Col: pair.sensor.Col}
		dLine := rudolph.ManhattanDistance(pair.sensor, nearestPtOnLine)

		if dLine <= dBeacon {
			availableD := dBeacon - dLine
			for i := nearestPtOnLine.Col - availableD; i <= nearestPtOnLine.Col+availableD; i++ {
				exclusions[i] = true
			}
		}
	}

	// remove existing beacons
	for _, pair := range sensorData {
		if pair.beacon.Row == line {
			delete(exclusions, pair.beacon.Col)
		}
	}

	return len(exclusions)
}

// day15Cmd represents the day15 command
var day15Cmd = &cobra.Command{
	Use:   "day15 path/to/input/file",
	Short: "AoC Day 15",
	Long:  `Advent of Code Day 15: Beacon Exclusion Zone`,
	Run: func(cmd *cobra.Command, args []string) {
		if input := toybag.ReadAsTokenizedStringSlice(args...); input != nil {
			fmt.Printf("%d sensor positions read.\n", len(input))

			defer xmas.PrintHolidayMessage(time.Now())

			sensorData := parseSensorData(input)

			if Parts.Has(1) {
				fmt.Println("Part 1 running...")

				fmt.Printf("The number of beacon-free squares on line 2,000,000 is **%d**.\n", countExcludedZones(sensorData, 2000000))
			}

			if Parts.Has(2) {
				fmt.Println("Part 2 running...")
				const maxrange int = 4000000
				min := rudolph.Coordinate{Row: 0, Col: 0}
				max := rudolph.Coordinate{Row: maxrange, Col: maxrange}

				beaconCoords := findEmergencyBeacon(sensorData, min, max)

				fmt.Printf("Emergency beacon tuning frequency: **%d**\n", beaconCoords.Col*maxrange+beaconCoords.Row)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(day15Cmd)
}
