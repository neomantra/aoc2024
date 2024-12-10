// https://adventofcode.com/2024/day/9
// go run 9/main.go 9/9.txt

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func minOf(x, y int) int {
	if x < y {
		return x
	}
	return y
}

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
	filesystem.makeFileMap()

	return &filesystem
}

func (fs *Filesystem) makeFileMap() {
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

func (fs *Filesystem) DefragBlock() {
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

///////////////////////////////////////////////////////////////////////////////

// returns -1 if not found
func (fs *Filesystem) searchForwardForFreespace(lastIndex int, spaceSize int) int {
	// search forward for free space
	for i := 0; i < lastIndex; i++ {
		if fs.fileMap[i] == 0 {
			// we found some free space, now search forward for the end of it
			count := 1
			freeIndex := i
			for i = i + 1; i < minOf(freeIndex+spaceSize, lastIndex); i++ {
				if fs.fileMap[i] != 0 {
					i--
					break
				}
				count++
			}
			if count >= spaceSize {
				return freeIndex
			}
		}
	}
	return -1
}

// returns true if we made a change
func (fs *Filesystem) DefragWholeFile() {
	// reverse scan for a file
	for i := len(fs.fileMap) - 1; i > 0; i-- {
		ival := fs.fileMap[i]
		if ival == 0 {
			continue
		}
		// we found a file... now search backward for its start
		fileSize, fileIndex := 1, i
		for i = i - 1; i > 0; i-- {
			if fs.fileMap[i] != ival {
				// we are at next file
				i++
				break
			}
			fileSize++
			fileIndex = i
		}
		// we are at a file, now search forward for free space for it
		freeIndex := fs.searchForwardForFreespace(fileIndex, fileSize)
		if freeIndex != -1 {
			// move file to free space
			for j := 0; j < fileSize; j++ {
				fs.fileMap[freeIndex+j] = fs.fileMap[fileIndex+j]
				fs.fileMap[fileIndex+j] = 0
			}
		}
	}
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
	fs.DefragBlock()
	fmt.Println("9.1:", fs.CalcChecksum())

	// part 1
	fs = NewFilesystem(string(diskMapData))
	//fmt.Println(fs.View())
	fs.DefragWholeFile()
	//fmt.Println(fs.View())
	fmt.Println("9.2:", fs.CalcChecksum())
}
