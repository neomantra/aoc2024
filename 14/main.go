// https://adventofcode.com/2024/day/14
// go run 14/main.go 14/14.txt

package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Point struct{ X, Y int }

type Robot struct {
	Pos Point
	Vel Point
}

func NewRobots(puzzle string) []Robot {
	var numberRegexp = regexp.MustCompile(`p=([-\d]*),([-\d]*) v=([-\d]*),([-\d]*)`)
	robots := []Robot{}
	for _, line := range strings.Split(puzzle, "\n") {
		match := numberRegexp.FindStringSubmatch(line)
		x, _ := strconv.Atoi(match[1])
		y, _ := strconv.Atoi(match[2])
		vx, _ := strconv.Atoi(match[3])
		vy, _ := strconv.Atoi(match[4])
		robots = append(robots, Robot{Point{x, y}, Point{vx, vy}})
	}
	return robots
}

func RobotHeatMap(robots []Robot, roomSize Point) string {
	heatMap := make([][]int, roomSize.Y)
	for y := 0; y < roomSize.Y; y++ {
		heatMap[y] = make([]int, roomSize.X)
		for x := 0; x < roomSize.X; x++ {
			heatMap[y][x] = 0
		}
	}
	for _, robot := range robots {
		heatMap[robot.Pos.Y][robot.Pos.X] += 1
	}
	var sb strings.Builder
	for y := 0; y < roomSize.Y; y++ {
		for x := 0; x < roomSize.X; x++ {
			val := heatMap[y][x] % 10
			if val == 0 {
				sb.WriteByte('.')
			} else {
				sb.WriteString(strconv.Itoa(val))
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

///////////////////////////////////////////////////////////////////////////////

func Operate(robots []Robot, roomSize Point, steps int) {
	for step := 0; step < steps; step++ {
		for i := 0; i < len(robots); i++ {
			robot := &robots[i]
			robot.Pos.X = (robot.Pos.X + robot.Vel.X)
			robot.Pos.Y = (robot.Pos.Y + robot.Vel.Y)

			if robot.Pos.X < 0 {
				robot.Pos.X = roomSize.X - (-robot.Pos.X % roomSize.X)
			} else if robot.Pos.X >= roomSize.X {
				robot.Pos.X = robot.Pos.X % roomSize.X
			}
			if robot.Pos.Y < 0 {
				robot.Pos.Y = roomSize.Y - (-robot.Pos.Y % roomSize.Y)
			} else if robot.Pos.Y >= roomSize.Y {
				robot.Pos.Y = robot.Pos.Y % roomSize.Y
			}
		}
	}
}

const UL, UR, LL, LR, NOPE = 1, 2, 3, 4, 0

func QuandrantOf(p Point, roomSize Point) int {
	if p.X < roomSize.X/2 && p.Y < roomSize.Y/2 {
		return UL
	} else if p.X > roomSize.X/2 && p.Y < roomSize.Y/2 {
		return UR
	} else if p.X < roomSize.X/2 && p.Y > roomSize.Y/2 {
		return LL
	} else if p.X > roomSize.X/2 && p.Y > roomSize.Y/2 {
		return LR
	}
	return NOPE
}

// returns (UL, UR, LL, LR)
func QuadrantScores(robots []Robot, roomSize Point) (ul int, ur int, ll int, lr int) {
	scores := make([]int, 5)
	for i := 0; i < len(robots); i++ {
		robotPos := robots[i].Pos
		quadrant := QuandrantOf(robotPos, roomSize)
		scores[quadrant]++
	}
	return scores[UL], scores[UR], scores[LL], scores[LR]
}

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Open and read data file
	robotsData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}

	// part 1
	robots := NewRobots(string(robotsData))
	if robots == nil {
		fmt.Fprintf(os.Stderr, "Bad robot data\n")
		os.Exit(1)
	}
	roomSize := Point{101, 103}
	Operate(robots, roomSize, 100)

	fmt.Println(RobotHeatMap(robots, roomSize))
	ul, ur, ll, lr := QuadrantScores(robots, roomSize)
	fmt.Println("ul:", ul, "ur:", ur, "ll:", ll, "lr:", lr)

	safetyFactor := ul * ur * ll * lr
	fmt.Println("14.1:", safetyFactor)
}
