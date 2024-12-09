// https://adventofcode.com/2024/day/9
// go run 9/main.go 9/9.txt

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func fillSlice[T any](t T, size int) []T {
	s := make([]T, size)
	for i := range s {
		s[i] = t
	}
	return s
}

///////////////////////////////////////////////////////////////////////////////

type Filesystem struct {
	puzzle  string
	diskMap []byte
	fileMap []int // stores id-1, zero is freespace
}

func NewFilesystem(puzzle string) *Filesystem {
	filesystem := Filesystem{puzzle: puzzle}
	// extract all the lines
	for _, c := range puzzle {
		filesystem.diskMap = append(filesystem.diskMap, byte(c)-byte('0'))
	}
	return &filesystem
}

func (fs *Filesystem) MakeFileMap() {
	var fileMap []int
	curID := 0
	for i, b := range fs.diskMap {
		fillVal := 0
		if i%2 == 0 {
			fillVal = curID + 1 // we store id+1
			curID += 1
		}
		fileMap = append(fileMap, fillSlice(fillVal, int(b))...)
	}
	fs.fileMap = fileMap
}

func (fs *Filesystem) Defrag() {
	// forward scan filemap
	lastJ := len(fs.fileMap) - 1 // last searched end
	for i := 0; i < lastJ; i++ {
		// not free space? continue
		if fs.fileMap[i] != 0 {
			continue
		}
		// at free space, search from back for filled block
		for j := lastJ; j > i; j-- {
			valJ := fs.fileMap[j]
			if valJ != 0 {
				fs.fileMap[i] = valJ
				fs.fileMap[j] = 0
				lastJ = j - 1 // next time start before us
				break
			}
		}
	}
}

func (fs *Filesystem) CalcChecksum() int {
	sum := 0
	for i, val := range fs.fileMap {
		if val != 0 {
			sum += i * (val - 1) // ID is val-1
		}
	}
	return sum
}

///////////////////////////////////////////////////////////////////////////////

func (fs *Filesystem) DiskMapView() string {
	var sb strings.Builder
	for _, b := range fs.diskMap {
		sb.WriteByte(b + '0')
	}
	return sb.String()
}

func (fs *Filesystem) FileMapView() string {
	var sb strings.Builder
	for _, val := range fs.fileMap {
		if val == 0 {
			sb.WriteByte('.')
		} else {
			r := byte(val-1) + '0'
			if r < 32 || r > 126 {
				r = '#'
			}
			sb.WriteByte(r)
		}
	}

	return sb.String()
}

func (fs *Filesystem) View() string {
	var style = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder())
	return lipgloss.JoinVertical(lipgloss.Left,
		style.Render(fs.DiskMapView()),
		style.Render(fs.FileMapView()))
}

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Open and read data file
	diskMapData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err.Error())
		os.Exit(1)
	}
	fs := NewFilesystem(string(diskMapData))
	if fs == nil {
		fmt.Fprintf(os.Stderr, "Bad disk map data\n")
		os.Exit(1)
	}

	// part 1
	fs.MakeFileMap()
	//fmt.Println(fs.View())
	fs.Defrag()
	//fmt.Println(fs.View())
	fmt.Println("9.1:", fs.CalcChecksum())
}
