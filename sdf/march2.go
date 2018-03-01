//-----------------------------------------------------------------------------
/*

Marching Squares

Convert an SDF2 boundary to a set of line segments.

*/
//-----------------------------------------------------------------------------

package sdf

//-----------------------------------------------------------------------------

type LineCache struct {
	base  V2        // base coordinate of line
	inc   V2        // dx, dy for each step
	steps V2i       // number of x,y steps
	val0  []float64 // SDF values for x line
	val1  []float64 // SDF values for x + dx line
}

func NewLineCache(base, inc V2, steps V2i) *LineCache {
	return &LineCache{base, inc, steps, nil, nil}
}

// Evaluate the SDF for a given Y line
func (l *LineCache) Evaluate(sdf SDF2, x int) {

	// Swap the layers
	l.val0, l.val1 = l.val1, l.val0

	ny := l.steps[1]
	dx, dy := l.inc.X, l.inc.Y

	// allocate storage
	if l.val1 == nil {
		l.val1 = make([]float64, ny+1)
	}

	// setup the loop variables
	idx := 0
	var p V2
	p.X = l.base.X + float64(x)*dx

	// evaluate the line
	p.Y = l.base.Y
	for y := 0; y < ny+1; y++ {
		l.val1[idx] = sdf.Evaluate(p)
		idx += 1
		p.Y += dy
	}
}

func (l *LineCache) Get(x, y int) float64 {
	if x == 0 {
		return l.val0[y]
	}
	return l.val1[y]
}

//-----------------------------------------------------------------------------

func MarchingSquares(sdf SDF2, box Box2, step float64) []Line2_PP {

	var lines []Line2_PP
	size := box.Size()
	base := box.Min
	steps := size.DivScalar(step).Ceil().ToV2i()
	inc := size.Div(steps.ToV2())

	// create the line cache
	l := NewLineCache(base, inc, steps)
	// evaluate the SDF for x = 0
	l.Evaluate(sdf, 0)

	nx, ny := steps[0], steps[1]
	dx, dy := inc.X, inc.Y

	var p V2
	p.X = base.X
	for x := 0; x < nx; x++ {
		// read the x + 1 layer
		l.Evaluate(sdf, x+1)
		// process all squares in the x and x + 1 layers
		p.Y = base.Y
		for y := 0; y < ny; y++ {
			x0, y0 := p.X, p.Y
			x1, y1 := x0+dx, y0+dy
			corners := [4]V2{
				V2{x0, y0},
				V2{x1, y0},
				V2{x1, y1},
				V2{x0, y1},
			}
			values := [4]float64{
				l.Get(0, y),
				l.Get(1, y),
				l.Get(1, y+1),
				l.Get(0, y+1),
			}
			lines = append(lines, ms_ToLines(corners, values, 0)...)
			p.Y += dy
		}
		p.X += dx
	}

	return lines
}

//-----------------------------------------------------------------------------

// generate the line segments for a square
func ms_ToLines(p [4]V2, v [4]float64, x float64) []Line2_PP {

	index := 0
	for i := 0; i < 4; i++ {
		if v[i] < x {
			index |= 1 << uint(i)
		}
	}

	if edgeTable[index] == 0 {
		return nil
	}

	var points [4]V2
	for i := 0; i < 4; i++ {
		bit := 1 << uint(i)
		if edgeTable[index]&bit != 0 {
			a := pairTable[i][0]
			b := pairTable[i][1]
			points[i] = ms_Interpolate(p[a], p[b], v[a], v[b], x)
		}
	}

	return nil
}

//-----------------------------------------------------------------------------

func ms_Interpolate(p1, p2 V2, v1, v2, x float64) V2 {
	if Abs(x-v1) < EPS {
		return p1
	}
	if Abs(x-v2) < EPS {
		return p2
	}
	if Abs(v1-v2) < EPS {
		return p1
	}
	t := (x - v1) / (v2 - v1)
	return V2{
		p1.X + t*(p2.X-p1.X),
		p1.Y + t*(p2.Y-p1.Y),
	}
}

//-----------------------------------------------------------------------------

// these are the vertex pairs for the edges
var ms_pairs = [][]int{
	{0, 1},
	{1, 2},
	{2, 3},
	{3, 0},
}

// 4 vertices -> 16 possible inside/outside combinations
// a 1 bit in the value indicates which edge has a line point
var ms_edges = [16]int{
	0, 0, 0, 0,
	0, 0, 0, 0,
	0, 0, 0, 0,
	0, 0, 0, 0,
}

var ms_lines = [][]int{
	{},
}

//-----------------------------------------------------------------------------
