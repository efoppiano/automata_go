package generator

type Generator interface {
	Random() float64
	RandInt(a int, b int) int
	Poi(l float64) int
}
