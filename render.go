package main 

import (
	// Basic
	"fmt"
	"bufio"
	"os"
	"path/filepath"

	// Strings 
	"strings"
	"strconv"
	
	// Image manipulation
	"image"
	"image/color"
	"image/png"
	"github.com/disintegration/imaging"

	// Math
	"math"
	"math/rand"
)

// 
func randColor() color.RGBA {
	r := uint8(255 * rand.Float64())
	g := uint8(255 * rand.Float64())
	b := uint8(255 * rand.Float64())
	return color.RGBA{r, g, b, 255}
}

// Vec2i integer
type Vec2i struct {
	x int
	y int
}

func newVec2i(x, y int) Vec2i {
	v := Vec2i{x: x, y: y}
	return v 
}

func barycentric(pts *[]*Vec2i, P *Vec2i) Vec3f{
	v0 := newVec3f(float64((*pts)[1].x - (*pts)[0].x), float64((*pts)[2].x - (*pts)[0].x), float64((*pts)[0].x - P.x))
	v1 := newVec3f(float64((*pts)[1].y - (*pts)[0].y), float64((*pts)[2].y - (*pts)[0].y), float64((*pts)[0].y - P.y))
	u := cross(&v0, &v1)
	if math.Abs(u.z) < 1 {
		return newVec3f(-1, 1, 1)
	}
	return newVec3f(1.0 - (u.x + u.y) / u.z, u.x / u.z, u.y / u.z)
}

// Vec3f float
type Vec3f struct {
	x float64
	y float64 
	z float64
}

func cross(v0, v1 *Vec3f) Vec3f {
	// Compute v0 X v1
	x := v0.y * v1.z - v0.z * v1.y 
	y := v0.z * v1.x - v0.x * v1.z 
	z := v0.x * v1.y - v0.y * v1.x 
	return newVec3f(x, y, z)
}

func dot(v0, v1 *Vec3f) float64 {
	return v0.x * v1.x + v0.y * v1.y + v0.z * v1.z
}

func (u *Vec3f) subtract(v *Vec3f) Vec3f{
	return newVec3f(u.x - v.x, u.y - v.y, u.z - v.z)
}

func newVec3f(x, y, z float64) Vec3f {
	v := Vec3f{x: x, y: y, z: z}
	return v
}


func (v *Vec3f) normalizeMinMax2D(m *Model) (float64, float64){
	x := (v.x - m.min_x) / (m.max_x - m.min_x) - 0.5
	y := (v.y - m.min_y) / (m.max_y - m.min_y) - 0.5
	return x, y
}

func (v *Vec3f) normalizeL2() {
	norm := math.Sqrt(v.x * v.x + v.y * v.y + v.z * v.z)
	v.x /= norm
	v.y /= norm 
	v.z /= norm
}

// Model
type Model struct {
	vertices []Vec3f
	faces [][]int

	min_x float64
	min_y float64
	max_x float64
	max_y float64
}

func newModel() Model{
	m := Model{min_x: 1e10, min_y: 1e10, max_x: -1e10, max_y: -1e10}
	return m
}

func (m *Model) aspectRatio() float64 {
	dx := m.max_x - m.min_x
	dy := m.max_y - m.min_y
	return dx / dy
}

func (m *Model) addVertex(v *Vec3f){
	m.vertices = append(m.vertices, *v)
}

func (m *Model) addFace(f *[]int){
	m.faces = append(m.faces, *f)
}

func (m *Model) nFaces() int {
	return len(m.faces)
}

func (m *Model) nVertices() int {
	return len(m.vertices)
}

func (m *Model) setMinMax(x, y float64) {
	m.min_x = math.Min(m.min_x, x)
	m.min_y = math.Min(m.min_y, y)
	m.max_x = math.Max(m.max_x, x)
	m.max_y = math.Max(m.max_y, y)
}

func parseObj(filePath string) Model {
	file, _ := os.Open(filePath)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	model := newModel()
	for scanner.Scan() {
		tok := strings.Split(scanner.Text(), " ")
		if len(tok) > 0 {
			if tok[0] == "v" {
				x, _ := strconv.ParseFloat(tok[1], 64)
				y, _ := strconv.ParseFloat(tok[2], 64)
				z, _ := strconv.ParseFloat(tok[3], 64)
				model.setMinMax(x, y)
				v := newVec3f(x, y, z)
				model.addVertex(&v)
			} else if tok[0] == "f" {
				var vs []int
				for i:=1; i<len(tok); i++ {
					v, _ := strconv.Atoi(strings.Split(tok[i], "/")[0])
					vs = append(vs, v - 1)
				}
				model.addFace(&vs)
			}
		}
	}
	return model
}

// Rendering 

func line(v0 *Vec2i, v1 *Vec2i, img *image.RGBA, color *color.RGBA) {
	x0, y0 := v0.x, v0.y
	x1, y1 := v1.x, v1.y

	var steep bool = false
	if math.Abs(float64(x0 - x1)) < math.Abs(float64(y0 - y1)) {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
		steep = true
	}
	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}

	var dx float64 = float64(x1 - x0)
	var dy float64 = float64(y1 - y0)
	var derr float64 = math.Abs(dy / dx)
	var err float64 = 0.0
	var y int = y0

	for x:=x0; x <= x1; x++ {
		if steep {
			img.Set(y, x, *color)
		} else {
			img.Set(x, y, *color)
		}
		err += derr 
		if err > 0.5 {
			if y1 > y0 {
				y += 1
			} else {
				y -= 1
			}
			err -= 1
		}
	}
}

func triangle(v0 *Vec2i, v1 *Vec2i, v2 *Vec2i, img *image.RGBA, color *color.RGBA, width int, height int) {
	if v0.y > v1.y {
		v0, v1 = v1, v0
	}
	if v0.y > v2.y {
		v0, v2 = v2, v0
	}
	if v1.y > v2.y {
		v1, v2 = v2, v1
	}
	
	pts := []*Vec2i{v0, v1, v2}

	bboxmin := newVec2i(width - 1, height - 1)
	bboxmax := newVec2i(0, 0)

	for i:=0; i<len(pts); i++ {
		bboxmin.x = int(math.Min(float64(bboxmin.x), float64(pts[i].x)))
		bboxmin.y = int(math.Min(float64(bboxmin.y), float64(pts[i].y)))
		bboxmax.x = int(math.Max(float64(bboxmax.x), float64(pts[i].x)))
		bboxmax.y = int(math.Max(float64(bboxmax.y), float64(pts[i].y)))
	}

	P := Vec2i{}
	for P.x=bboxmin.x; P.x<bboxmax.x; P.x++ {
		for P.y=bboxmin.y; P.y<bboxmax.y; P.y++ {
			v := barycentric(&pts, &P)
			if v.x < 0 || v.y < 0 || v.z < 0 {
				continue
			}
			img.Set(P.x, P.y, *color)
		}
	}
}

func renderWireframe(model *Model, img *image.RGBA, color *color.RGBA, width int, height int) {
	// fill
	for i:=0; i<model.nFaces(); i++ {
		face := model.faces[i]
		for j:=0; j<len(face); j++ {
			v0 := model.vertices[face[j]]
			v1 := model.vertices[face[(j+1)%len(face)]]

			x0, y0 := v0.normalizeMinMax2D(model)  // normalize w.r.t min, max boundaries
			x1, y1 := v1.normalizeMinMax2D(model)

			scale := 1.5
			x0 = (x0 + 0.5 * scale) * float64(width)  / scale
			y0 = (y0 + 0.5 * scale) * float64(height) / scale
			x1 = (y1 + 0.5 * scale) * float64(width)  / scale
			y1 = (y0 + 0.5 * scale) * float64(height) / scale

			line(&Vec2i{x:int(x0), y:int(y0)}, &Vec2i{x:int(x1), y:int(y1)}, img, color)
		}
	}
}

func renderTriangleMesh(model *Model, img *image.RGBA, fillColor *color.RGBA, width int, height int) {
	// fill

	lightDir := newVec3f(0, 0, -1)
	lightDir.normalizeL2()

	for i:=0; i<model.nFaces(); i++ {
		face := model.faces[i]
		
		var screenCoords [3]Vec2i
		var worldCoords [3]Vec3f
		for j:=0; j<3; j++ {
			v := model.vertices[face[j]]
			x, y := v.normalizeMinMax2D(model)  // normalize w.r.t min, max boundaries

			scale := 1.5
			x = (x + 0.5 * scale) * float64(width)  / scale
			y = (y + 0.5 * scale) * float64(height) / scale
			
			screenCoords[j] = newVec2i(int(x), int(y))
			worldCoords[j] = v
		}
		v0 := worldCoords[2].subtract(&worldCoords[0])
		v1 := worldCoords[1].subtract(&worldCoords[0])
		n := cross(&v0, &v1)
		n.normalizeL2()
		I := dot(&n, &lightDir)
		if I > 0 {
			r, g, b := I * float64(fillColor.R), I * float64(fillColor.G), I * float64(fillColor.B)
			fill := color.RGBA{uint8(r), uint8(g), uint8(b), fillColor.A}
			triangle(&screenCoords[0], &screenCoords[1], &screenCoords[2], img, &fill, width, height)
		}
	}
}

func main() {
	// Parse .obj file
	relPath := "./obj/bunny.obj"
	absPath, _ := filepath.Abs(relPath)
	model := parseObj(absPath)

	fmt.Println("Number of faces: ", model.nFaces())
	fmt.Println("Number of vertices: ", model.nVertices())
	ratio := model.aspectRatio()

	// Create canvas
	height := 1000
	width := int(ratio * float64(height))

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Render
	//renderWireframe(&model, img, &color.RGBA{0, 0, 0, 255}, width, height)
	renderTriangleMesh(&model, img, &color.RGBA{255, 255, 255, 255}, width, height)

	// Save
	f, _ := os.Create("./results/test.png")
	png.Encode(f, imaging.FlipV(img))
}