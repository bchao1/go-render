package main 

import (
	"fmt"
	"bufio"
	"os"
	//"reflect"
	"strings"
	"strconv"
	
	"image"
	"image/color"
	"image/png"

	"math"
	"github.com/disintegration/imaging"
)
func line(x0 int, y0 int, x1 int, y1 int, img *image.RGBA) {
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
			img.Set(y, x, color.White)
		} else {
			img.Set(x, y, color.White)
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

type Vec3d struct {
	x float64
	y float64 
	z float64
}

func (v *Vec3d) normalize(m *Model) {
	(*v).x = ((*v).x - (*m).min_x) / ((*m).max_x - (*m).min_x) - 0.5
	(*v).y = ((*v).y - (*m).min_y) / ((*m).max_y - (*m).min_y) - 0.5
}

type Face struct {
	v []int
}

type Model struct {
	vertices []Vec3d
	faces [][]int

	min_x float64
	min_y float64
	max_x float64
	max_y float64
}

func (m *Model) addVertex(v *Vec3d){
	(*m).vertices = append((*m).vertices, *v)
}

func (m *Model) addFace(f *[]int){
	(*m).faces = append((*m).faces, *f)
}

func (m *Model) nFaces() int {
	return len((*m).faces)
}

func (m *Model) nVertices() int {
	return len((*m).vertices)
}

func (m *Model) setMinMax(x, y float64) {
	(*m).min_x = math.Min((*m).min_x, x)
	(*m).min_y = math.Min((*m).min_y, y)
	(*m).max_x = math.Max((*m).max_x, x)
	(*m).max_y = math.Max((*m).max_y, y)
}

func newVec3d(x, y, z float64) Vec3d {
	v := Vec3d{x: x, y: y, z: z}
	return v
}

func newModel() Model{
	m := Model{min_x: 1e10, min_y: 1e10, max_x: -1e10, max_y: -1e10}
	return m
}

func main() {
	file, _ := os.Open("./bunny.obj")
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
				v := newVec3d(x, y, z)
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
	fmt.Println("Number of faces: ", model.nFaces())
	fmt.Println("Number of vertices: ", model.nVertices())

	width := 1000
	height := 1000

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// fill
	for i:=0; i<model.nFaces(); i++ {
		face := model.faces[i]
		for j:=0; j<len(face); j++ {
			v0 := model.vertices[face[j]]
			v1 := model.vertices[face[(j+1)%len(face)]]

			v0.normalize(&model)  // normalize w.r.t min, max boundaries
			v1.normalize(&model)

			scale := 1.5
			x0 := int((v0.x + 0.5 * scale) * float64(width)  / scale)
			y0 := int((v0.y + 0.5 * scale) * float64(height) / scale)
			x1 := int((v1.x + 0.5 * scale) * float64(width)  / scale)
			y1 := int((v1.y + 0.5 * scale) * float64(height) / scale)

			line(x0, y0, x1, y1, img)
		}
	}
	f, _ := os.Create("image.png")
	png.Encode(f, imaging.FlipV(img))
}