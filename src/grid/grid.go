package grid

import (
	"errors"
	"fmt"
	"go_automata/src/generator"
	"go_automata/src/utils"
)

type Grid struct {
	rows      int
	cols      int
	grid      [][]interface{}
	generator generator.Generator
}

func NewGrid(rows, cols int, generator generator.Generator) *Grid {
	grid := make([][]interface{}, rows)
	for i := range grid {
		grid[i] = make([]interface{}, cols)
	}
	return &Grid{rows, cols, grid, generator}
}

func (g *Grid) IsFill(row, col int) bool {
	if row < 0 || row >= g.rows {
		return false
	}
	if col < 0 || col >= g.cols {
		return false
	}
	return g.grid[row][col] != nil
}

func (g *Grid) Fill(row, col int, v interface{}) error {
	if g.IsFill(row, col) {
		return errors.New("cell already filled")
	}
	g.grid[row][col] = v
	return nil
}

func (g *Grid) Clear(row, col int) {
	if !g.IsFill(row, col) {
		panic("Attempted to clear an empty cell")
	}
	g.grid[row][col] = nil
}

func (g *Grid) GetValue(row, col int) interface{} {
	if !g.IsFill(row, col) {
		panic("Element not found")
	}
	return g.grid[row][col]
}

func (g *Grid) CalcDistToNext(row, col int, f func(interface{}) bool, maxChecks int) int {
	if f == nil {
		f = func(interface{}) bool { return true }
	}
	if maxChecks == 0 {
		maxChecks = g.cols - col
	}
	for i := col + 1; i < min(g.cols, col+maxChecks+1); i++ {
		if g.IsFill(row, i) && f(g.GetValue(row, i)) {
			return i - col - 1
		}
	}
	return -1
}

func (g *Grid) CalcDistToPrev(row, col int, f func(interface{}) bool, maxChecks int) int {
	if f == nil {
		f = func(interface{}) bool { return true }
	}
	if maxChecks == 0 {
		maxChecks = col
	}
	for i := col - 1; i > max(-1, col-maxChecks-1); i-- {
		if g.IsFill(row, i) && f(g.GetValue(row, i)) {
			return col - i - 1
		}
	}
	return -1
}

func (g *Grid) CalcDistToVerticallyNext(row, col int, f func(interface{}) bool, maxChecks int) int {
	if f == nil {
		f = func(interface{}) bool { return true }
	}
	if maxChecks == 0 {
		maxChecks = g.rows - row
	}
	for i := row + 1; i < min(g.rows, row+maxChecks+1); i++ {
		if g.IsFill(i, col) && f(g.GetValue(i, col)) {
			return i - row - 1
		}
	}
	return -1
}

func (g *Grid) CalcDistToVerticallyPrev(row, col int, f func(interface{}) bool, maxChecks int) int {
	if f == nil {
		f = func(interface{}) bool { return true }
	}
	if maxChecks == 0 {
		maxChecks = row
	}
	for i := row - 1; i > max(-1, row-maxChecks-1); i-- {
		if g.IsFill(i, col) && f(g.GetValue(i, col)) {
			return row - i - 1
		}
	}
	return -1
}

func (g *Grid) GetPrev(row, col int, f func(interface{}) bool, maxChecks int) interface{} {
	if maxChecks == 0 {
		maxChecks = col
	}
	dist := g.CalcDistToPrev(row, col, f, maxChecks)
	if dist == -1 {
		return nil
	}
	return g.GetValue(row, col-dist-1)
}

func (g *Grid) GetVerticallyPrev(row, col int, f func(interface{}) bool, maxChecks int) interface{} {
	if maxChecks == 0 {
		maxChecks = row
	}
	dist := g.CalcDistToVerticallyPrev(row, col, f, maxChecks)
	if dist == -1 {
		return nil
	}
	return g.GetValue(row-dist-1, col)
}

func (g *Grid) GetNext(row, col int, f func(interface{}) bool, maxChecks int) interface{} {
	if maxChecks == 0 {
		maxChecks = g.cols - col
	}
	dist := g.CalcDistToNext(row, col, f, maxChecks)
	if dist == -1 {
		return nil
	}
	return g.GetValue(row, col+dist+1)
}

func (g *Grid) GetVerticallyNext(row, col int, f func(interface{}) bool, maxChecks int) interface{} {
	if maxChecks == 0 {
		maxChecks = g.rows
	}
	dist := g.CalcDistToVerticallyNext(row, col, f, maxChecks)
	if dist == -1 {
		return nil
	}
	return g.GetValue(row+dist+1, col)
}

func (g *Grid) getCellsWithValue() []interface{} {
	values := []interface{}{}
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			if g.IsFill(i, j) {
				values = append(values, g.GetValue(i, j))
			}
		}
	}
	return values
}

func (g *Grid) Apply(f func(interface{})) {
	values := g.getCellsWithValue()
	for len(values) > 0 {
		pos := g.generator.RandInt(0, len(values))
		value := values[pos]
		values = append(values[:pos], values[pos+1:]...)
		f(value)
	}
}

func (g *Grid) Rows() int {
	return g.rows
}

func (g *Grid) Cols() int {
	return g.cols
}

func (g *Grid) Plot(f func(*utils.Point, interface{}) string, bounds *utils.Rectangle) {
	if bounds == nil {
		bounds = utils.NewRectangle(g.rows, g.cols)
	}

	plotColNumbers(bounds)

	for row := bounds.StartRow(); row <= bounds.EndRow(); row++ {
		plotRowNumber(row)
		point := &utils.Point{X: row, Y: 0}
		for col := bounds.StartCol(); col <= bounds.EndCol(); col++ {
			if !bounds.IsInside(point) {
				continue
			}
			var s string
			if g.IsFill(row, col) {
				s = f(point, g.GetValue(row, col))
			} else {
				s = f(point, nil)
			}
			print(s)
		}
		println()
	}
}

func plotColNumbers(bounds *utils.Rectangle) {
	print("  ")
	for col := bounds.StartCol(); col <= bounds.EndCol(); col++ {
		if col%2 == 0 {
			print("\033[91m")
		} else {
			print("\033[94m")
		}
		fmt.Printf("%2d", col)
		print("\033[0m")
	}
	println()
}

func plotRowNumber(row int) {
	if row%2 == 0 {
		print("\033[91m")
	} else {
		print("\033[94m")
	}
	fmt.Printf("%2d", row)
	print("\033[0m")
}
