package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"strconv"
	"time"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	colorful "github.com/lucasb-eyer/go-colorful"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type point struct {
	x float64
	y float64
}

type triangle [3]*point

func newPoint(x, y float64) *point {
	return &point{x, y}
}

func (p *point) String() string {
	return fmt.Sprintf("(%v, %v)", p.x, p.y)
}

// x and y are the boundaries of the shape
func generatePoints(numSegments int, x, y float64) []*point {
	points := make([]*point, 2)

	bucketSize := x / float64(numSegments)

	rndY := func() float64 {
		return (.4)*y*rand.Float64() + y*.3
	}

	points[0] = newPoint(0, y)
	points[1] = newPoint(0, rndY())
	for i := 0; i < numSegments; i++ {
		var x float64
		for x < (.2)*bucketSize+points[i+1].x {
			x = float64(i)*bucketSize + bucketSize*rand.Float64()
		}

		points = append(points, newPoint(x, rndY()))
	}

	penult := newPoint(x, rndY())
	end := newPoint(x, y)

	points = append(points, penult, end)
	return points
}

func convertToTriangles(points []*point, x, y float64) []*triangle {
	// ignore the beginning and end points (they close the shape)

	triangles := make([]*triangle, 0)

	for i := 1; i < len(points)-2; i++ {
		topA := points[i]
		topB := points[i+1]

		// left triangle
		left := &triangle{
			newPoint(topA.x, y),
			topA,
			topB,
		}

		// right triangle
		right := &triangle{
			topB,
			newPoint(topB.x, y),
			newPoint(topA.x, y),
		}

		triangles = append(triangles, left, right)
	}

	return triangles
}

var (
	x = 2436
	y = 1125
	n = 15
)

func drawShape(canvas *image.RGBA, points []*point) {
	gc := draw2dimg.NewGraphicContext(canvas)
	gc.SetStrokeColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.SetLineWidth(5)

	gc.MoveTo(points[0].x, points[0].y)
	for _, p := range points[1:] {

		gc.LineTo(p.x, p.y)
	}

	gc.Close()

	// cornflower blue best color
	gc.SetFillColor(color.RGBA{0x54, 0x95, 0xed, 0xff})
	gc.FillStroke()

	triangles := convertToTriangles(points, float64(x), float64(y))

	for _, t := range triangles {
		gc.BeginPath()

		gc.MoveTo(t[0].x, t[0].y)
		gc.LineTo(t[1].x, t[1].y)
		gc.LineTo(t[2].x, t[2].y)
		gc.Close()

		gc.SetStrokeColor(colorful.FastHappyColor())
		gc.Stroke()
	}

}

func labelPoints(canvas *image.RGBA, points []*point) {
	gc := draw2dimg.NewGraphicContext(canvas)
	gc.SetFontData(draw2d.FontData{Name: "helvetica"})
	gc.SetFontSize(32)
	gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})

	for i, p := range points {
		x := 0.0
		if i > len(points)-3 {
			x = p.x - 50
		} else {
			x = p.x + 5
		}

		gc.FillStringAt(strconv.Itoa(i), x, p.y-5)
	}
}

func main() {
	canvas := image.NewRGBA(image.Rect(0, 0, x, y))
	gc := draw2dimg.NewGraphicContext(canvas)

	gc.SetFillColor(color.RGBA{0x22, 0x22, 0x22, 0xff})
	gc.ClearRect(0, 0, x, y)

	points := generatePoints(n, float64(x), float64(y))

	drawShape(canvas, points)

	labelPoints(canvas, points)

	draw2dimg.SaveToPngFile("test.png", canvas)
}
