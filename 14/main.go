// https://adventofcode.com/2024/day/14
// go run 14/main.go 14/14.txt

package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/NimbleMarkets/ollamatea"
	"github.com/ollama/ollama/api"
	ollama "github.com/ollama/ollama/api"
	ansitoimage "github.com/pavelpatrin/go-ansi-to-image"
)

type Point struct{ X, Y int }

type PointF64 struct{ X, Y float64 }

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

///////////////////////////////////////////////////////////////////////////////

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

type RobotHeatMap struct {
	Map      [][]int
	RoomSize Point
}

func MakeRobotHeatMap(robots []Robot, roomSize Point) RobotHeatMap {
	hm := RobotHeatMap{
		Map:      make([][]int, roomSize.Y),
		RoomSize: roomSize,
	}
	for y := 0; y < roomSize.Y; y++ {
		hm.Map[y] = make([]int, roomSize.X)
		for x := 0; x < roomSize.X; x++ {
			hm.Map[y][x] = 0
		}
	}
	for _, robot := range robots {
		hm.Map[robot.Pos.Y][robot.Pos.X] += 1
	}
	return hm
}

func (hm RobotHeatMap) View() string {
	var sb strings.Builder
	for y := 0; y < hm.RoomSize.Y; y++ {
		for x := 0; x < hm.RoomSize.X; x++ {
			val := hm.Map[y][x] % 10
			if val == 0 {
				sb.WriteByte(' ')
			} else {
				sb.WriteString(strconv.Itoa(val))
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

// returns meanVal, centroid, variance
func (hm RobotHeatMap) GetMetrics() (float64, PointF64, PointF64) {
	// calculate centroid
	var centroid PointF64
	var meanVal, count float64
	for y := 0; y < hm.RoomSize.Y; y++ {
		for x := 0; x < hm.RoomSize.X; x++ {
			val := float64(hm.Map[y][x])
			if val == 0 {
				continue
			}
			centroid.X += float64(x)
			centroid.Y += float64(y)
			meanVal += val
			count++
		}
	}
	centroid.X /= count
	centroid.Y /= count
	meanVal /= count

	// calculate center of mass
	var variance PointF64
	for y := 0; y < hm.RoomSize.Y; y++ {
		for x := 0; x < hm.RoomSize.X; x++ {
			val := float64(hm.Map[y][x])
			if val == 0 {
				continue
			}
			dx := val * (float64(x) - centroid.X)
			dy := val * (float64(y) - centroid.Y)
			variance.X += dx * dx
			variance.Y += dy * dy
		}
	}
	variance.X = math.Sqrt(variance.X)
	variance.Y = math.Sqrt(variance.Y)

	return meanVal, centroid, variance
}

///////////////////////////////////////////////////////////////////////////////

func DoesOllamaThinkTheresAChristmasTrees(hm RobotHeatMap, stepNum int) (string, bool) {
	// Use OllamaTeas's machinery to convert to image
	hmView := hm.View()
	convertConfig := ansitoimage.DefaultConfig
	convertConfig.PageRows = hm.RoomSize.Y
	convertConfig.PageCols = hm.RoomSize.X
	pngBytes, err := ollamatea.ConvertTerminalTextToImage(string(hmView), &convertConfig)

	ollamaURL, err := url.Parse("http://localhost:11434")
	if err != nil {
		return "", false
	}

	systemPrompt := `"Does the supplied image contain framed image of the outline of an evergreen tree?" 
	Include YES or NO and a brief reason why.  You must be absolutely sure, you cannot be wrong about this. 
	Christmas trees are evergreen trees and may have lights and ornaments but that is not necessary.`

	ollamaClient := ollama.NewClient(ollamaURL, http.DefaultClient)
	req := &ollama.GenerateRequest{
		Model:  "llama3.2-vision:11b-instruct-q8_0", //"llama3.2-vision",
		Prompt: "does the image contain the shape of an evergreen christmas tree?",
		System: systemPrompt,
		Images: []api.ImageData{pngBytes},
	}
	os.WriteFile(fmt.Sprintf("out.%d.png", stepNum), pngBytes, 0644)

	var sb strings.Builder
	respFunc := func(resp ollama.GenerateResponse) error {
		sb.WriteString(resp.Response)
		return nil
	}

	ctx := context.Background()
	err = ollamaClient.Generate(ctx, req, respFunc)
	if err != nil {
		return "", false
	}

	response := sb.String()
	if strings.Contains(strings.ToUpper(response), "YES") {
		return response, true
	}
	return response, false
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
	// fmt.Println(MakeRobotHeatMap(robots, roomSize).View())
	ul, ur, ll, lr := QuadrantScores(robots, roomSize)
	fmt.Println("ul:", ul, "ur:", ur, "ll:", ll, "lr:", lr)

	safetyFactor := ul * ur * ll * lr
	fmt.Println("14.1:", safetyFactor)

	// part 2
	robots = NewRobots(string(robotsData))
	stepsToTree := 0
	for {
		hm := MakeRobotHeatMap(robots, roomSize)
		_, _, stddev := hm.GetMetrics()
		if stddev.X < 450 && stddev.Y < 450 { // emperically determined
			break
		}
		Operate(robots, roomSize, 1)
		stepsToTree++
	}
	fmt.Print(MakeRobotHeatMap(robots, roomSize).View(), "\n", stepsToTree, "\n")

	// part 2 Ollama-version
	fmt.Print("\nOllama version\n")
	robots = NewRobots(string(robotsData))
	// CHEATING!  =) Stepping forward a bunch since we knew answer above
	llmStepsToTree := stepsToTree - 3
	Operate(robots, roomSize, llmStepsToTree)
	for {
		hm := MakeRobotHeatMap(robots, roomSize)

		start := time.Now()
		response, isTree := DoesOllamaThinkTheresAChristmasTrees(hm, llmStepsToTree)
		duration := time.Since(start)

		fmt.Printf("%s\n\nStep %d Ollama took %0.2fs\n\n",
			response, llmStepsToTree, duration.Seconds())
		if isTree {
			break
		}
		Operate(robots, roomSize, 1)
		llmStepsToTree++
	}
	fmt.Println("14.2llm:", llmStepsToTree)
}
