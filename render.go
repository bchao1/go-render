package main 

import (
	// Basic
	"fmt"
	"bufio"
	"os"
	// "path/filepath"

	// Strings 
	"strings"
	"strconv"
	
	// Image manipulation
	"image"
	"image/color"
	"image/png"
	_ "image/jpeg"
	"image/draw"
	"github.com/disintegration/imaging"

	// Math
	"math"
	"math/rand"
	
	//"sort"
)

func randColor() color.RGBA {
	r := uint8(255 * rand.Float64())
	g := uint8(255 * rand.Float64())
	b := uint8(255 * rand.Float64())
	return color.RGBA{r, g, b, 255}
}

func worldToScreen(v *Vec3f, model *Model, width int, height int, scale float64) Vec3f {
	x, y := v.x, v.y
	// Center align
	x = float64(width) * ((x - 0.5) / scale + 0.5)
	y = float64(height) * ((y - 0.5) / scale + 0.5)
	return newVec3f(float64(int(x)), float64(int(y)), v.z)
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

func barycentric(pts *[]*Vec3f, P *Vec3f) Vec3f{
	v0 := newVec3f(float64((*pts)[1].x - (*pts)[0].x), float64((*pts)[2].x - (*pts)[0].x), float64((*pts)[0].x - P.x))
	v1 := newVec3f(float64((*pts)[1].y - (*pts)[0].y), float64((*pts)[2].y - (*pts)[0].y), float64((*pts)[0].y - P.y))
	u := cross(&v0, &v1)
	if math.Abs(u.z) > 1e-2 {
		return newVec3f(1.0 - (u.x + u.y) / u.z, u.x / u.z, u.y / u.z)
	}
	return newVec3f(-1, 1, 1)
}

// Vec2f float

type Vec2f struct {
	x float64 
	y float64
}

func newVec2f(x, y float64) Vec2f {
	v := Vec2f{x: x, y: y}
	return v
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

func (u *Vec3f) subtract(v *Vec3f, inplace bool) Vec3f{
	if inplace {
		u.x -= v.x
		u.y -= v.y
		u.z -= v.z
		return *u 
	} else {
		return newVec3f(u.x - v.x, u.y - v.y, u.z - v.z)
	}
}

func (u *Vec3f) add(v *Vec3f, inplace bool) Vec3f{
	if inplace {
		u.x += v.x
		u.y += v.y
		u.z += v.z
		return *u 
	} else {
		return newVec3f(u.x + v.x, u.y + v.y, u.z + v.z)
	}
}

func (u *Vec3f) mul(m float64, inplace bool) Vec3f{
	if inplace {
		u.x *= m
		u.y *= m
		u.z *= m
		return *u
	} else {
		return newVec3f(u.x * m, u.y * m, u.z * m)
	}
}

func (u *Vec3f) div(m float64, inplace bool) Vec3f{
	if inplace {
		u.x /= m
		u.y /= m
		u.z /= m
		return *u
	} else {
		return newVec3f(u.x / m, u.y / m, u.z / m)
	}
}

func (v *Vec3f) project(c float64) Vec3f {
	return v.div(1 - v.z / c, false)
}

func newVec3f(x, y, z float64) Vec3f {
	v := Vec3f{x: x, y: y, z: z}
	return v
}

func (m *Model) centerAlignShift() {
	// Shift amount to center align
	dx, dy := -(m.max_x + m.min_x) / 2, -(m.max_y + m.min_y) / 2
	for i:=0; i<m.nVertices(); i++ {
		m.vertices[i].x += dx 
		m.vertices[i].y += dy
	}
	m.min_x += dx
	m.max_x += dx 
	m.min_y += dy 
	m.max_y += dy

	m.origVertices = append([]Vec3f{}, m.vertices...)

}

func (v *Vec3f) normalizeL2() {
	norm := math.Sqrt(v.x * v.x + v.y * v.y + v.z * v.z)
	v.x /= norm
	v.y /= norm 
	v.z /= norm
}

func (v *Vec3f) normalizeCenteredCube(m *Model) Vec3f{
	// Normalize to [0, 1]
	x := (v.x - m.min_x) / (m.max_x - m.min_x)
	y := (v.y - m.min_y) / (m.max_y - m.min_y)
	return newVec3f(x, y, v.z)
}

// Model
type Model struct {
	vertices []Vec3f
	origVertices []Vec3f

	faces [][]int

	vertexFaceNeighbors [][]int

	faceNormals []Vec3f
	vertexNormals []Vec3f

	textureCoordinates []Vec2f
	faceTextures [][]int

	min_x float64
	min_y float64
	max_x float64
	max_y float64
}

func newModel() Model{
	m := Model{min_x: 1e10, min_y: 1e10, max_x: -1e10, max_y: -1e10}
	return m
}

func (model *Model) computeFaceNormals() {
	model.faceNormals = make([]Vec3f, model.nFaces())
	for i:=0; i<model.nFaces(); i++ {
		face := model.faces[i]
		var worldCoords [3]Vec3f
		for j:=0; j<3; j++ {
			world_v := model.vertices[face[j]]
			worldCoords[j] = world_v
		}
		v0 := worldCoords[2].subtract(&worldCoords[0], false)
		v1 := worldCoords[1].subtract(&worldCoords[0], false)
		n := cross(&v0, &v1)
		//n.div(n.z, true)
		n.normalizeL2() // normalize!!
		model.faceNormals[i] = n
	}
}

func (model *Model) computeVertexNormals() {
	model.vertexNormals = make([]Vec3f, model.nVertices())
	for i:=0; i<model.nVertices(); i++ {
		n := newVec3f(0.0, 0.0, 0.0)
		nfaces := len(model.vertexFaceNeighbors[i])
		for j:=0; j<nfaces; j++ {
			f := model.vertexFaceNeighbors[i][j]
			n.add(&model.faceNormals[f], true)
		}
		n.div(float64(nfaces), true)
		n.normalizeL2() // normalize!!
		model.vertexNormals[i] = n
	}
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

func (m *Model) addTexture(t *Vec2f){
	m.textureCoordinates = append(m.textureCoordinates, *t)
}

func (m *Model) addFaceTexture(ft *[]int) {
	m.faceTextures = append(m.faceTextures, *ft)
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

func (m *Model) transformCoordinates(eye *Vec3f, center *Vec3f, up *Vec3f) {
	// transform world vertices (projection, rotation, etc)
	
	m.vertices = append([]Vec3f{}, m.origVertices...) // reset to world coordinates

	// Compute camera scene basis
	b3 := eye.subtract(center, false)
	b3.normalizeL2()
	b1 := cross(up, &b3)
	b1.normalizeL2()
	b2 := cross(&b3, &b1)
	b2.normalizeL2()

	//c := 0.5
	for i:=0; i<m.nVertices(); i++ {
		// z := m.vertices[i].z
		//m.vertices[i].div(1 - z / c, true) 
		v := &m.vertices[i]
		x := v.x * b1.x + v.y * b2.x + v.z * b3.x
		y := v.x * b1.y + v.y * b2.y + v.z * b3.y
		z := v.x * b1.z + v.y * b2.z + v.z * b3.z
		v.x, v.y, v.z = x, y, z
	}
	return
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
				model.vertexFaceNeighbors = append(model.vertexFaceNeighbors, make([]int, 0))
			} else if tok[0] == "f" {
				var vs []int
				var vtextures []int
				for i:=1; i<len(tok); i++ {
					indices := strings.Split(tok[i], "/") // coordinate indices
					v, _ := strconv.Atoi(indices[0])
					vs = append(vs, v - 1)
					if len(indices) > 1 { // also with textures
						vt, _ := strconv.Atoi(indices[1])
						vtextures = append(vtextures, vt - 1)
					}
				}
				model.addFace(&vs)
				model.addFaceTexture(&vtextures)
			} else if tok[0] == "vt" {
				// Texture coordinates
				x, _ := strconv.ParseFloat(tok[1], 64)
				y, _ := strconv.ParseFloat(tok[2], 64)
				tc := newVec2f(x, y) // texture coordinates
				model.addTexture(&tc)
			}
		}
	}
	// Neighbor faces of vertices
	for i:=0; i<model.nFaces(); i++ {
		//fmt.Println(model.faces[i])
		for j:=0; j<len(model.faces[i]); j++ {
			v := model.faces[i][j]
			model.vertexFaceNeighbors[v] = append(model.vertexFaceNeighbors[v], i)
		} 
	}

	model.centerAlignShift()

	return model
}

// Rendering 

func line(v0 *Vec3f, v1 *Vec3f, img *image.RGBA, color *color.RGBA) {
	x0, y0 := v0.x, v0.y
	x1, y1 := v1.x, v1.y

	var steep bool = false
	if math.Abs(x0 - x1) < math.Abs(y0 - y1) {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
		steep = true
	}
	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}

	var dx float64 = x1 - x0
	var dy float64 = y1 - y0
	var derr float64 = math.Abs(dy / dx)
	var err float64 = 0.0
	var y int = int(y0)

	for x:=int(x0); x <= int(x1); x++ {
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

func triangle(
	v0 *Vec3f, v1 *Vec3f, v2 *Vec3f,  // vertices
	vertexNormals *[]Vec3f, vertexTextures *[]Vec2f, faceNormal *Vec3f, lightDir *Vec3f,  // vectors
	img *image.RGBA, textureImage *image.Image, zbuffer *[]float64, 
	fillColor *color.RGBA, width int, height int, specCoeff float64) {
	
	pts := []*Vec3f{v0, v1, v2}

	bboxmin := newVec2f(math.Inf(1), math.Inf(1))
	bboxmax := newVec2f(math.Inf(-1), math.Inf(-1))
	clamp := newVec2f(float64(width - 1), float64(height - 1))

	for i:=0; i<len(pts); i++ {
		bboxmin.x = math.Max(0.0, math.Min(bboxmin.x, float64(pts[i].x)))
		bboxmin.y = math.Max(0.0, math.Min(bboxmin.y, float64(pts[i].y)))
		bboxmax.x = math.Min(clamp.x, math.Max(bboxmax.x, float64(pts[i].x)))
		bboxmax.y = math.Min(clamp.y, math.Max(bboxmax.y, float64(pts[i].y)))
	}

	P := Vec3f{}
	for P.x=bboxmin.x; P.x<bboxmax.x; P.x++ {
		for P.y=bboxmin.y; P.y<bboxmax.y; P.y++ {
			v := barycentric(&pts, &P)
			if v.x < 0 || v.y < 0 || v.z < 0 {
				continue
			}
			P.z = v.x * pts[0].z + v.y * pts[1].z + v.z * pts[2].z
			if (*zbuffer)[int(P.x + P.y * float64(width))] < P.z {
				(*zbuffer)[int(P.x + P.y * float64(width))] = P.z
				diff, spec := phongShading(vertexNormals, lightDir, &v)
				diff = math.Max(0.0, diff)
				spec = math.Max(0.0, spec)
				var fill color.RGBA
				if textureImage == nil {
					fill = getColor(*fillColor, diff, spec, specCoeff)
				} else {
					fill = getColorFromTexture(textureImage, vertexTextures, &v, diff, spec, specCoeff)
				}
				img.Set(int(P.x), int(P.y), fill)
			}
		}
	}
}

func renderWireframe(model *Model, img *image.RGBA, color *color.RGBA, width int, height int, scale float64) {
	// fill
	for i:=0; i<model.nFaces(); i++ {
		face := model.faces[i]
		for j:=0; j<len(face); j++ {
			world_v0 := model.vertices[face[j]]
			world_v1 := model.vertices[face[(j+1)%len(face)]]

			world_v0 = world_v0.normalizeCenteredCube(model)
			world_v1 = world_v1.normalizeCenteredCube(model)

			screen_v0 := worldToScreen(&world_v0, model, width, height, scale)
			screen_v1 := worldToScreen(&world_v1, model, width, height, scale)
			
			line(&screen_v0, &screen_v1, img, color)
		}
	}
}

func renderTriangleMesh(
	model *Model, img *image.RGBA, textureImage *image.Image, 
	fillColor *color.RGBA, lightDir *Vec3f, eye *Vec3f, center *Vec3f, up *Vec3f,
	width int, height int, scale float64, specCoeff float64) {
	// fill
	lightDir.normalizeL2()

	var zbuffer = make([]float64, width * height)
	for i:=0; i<len(zbuffer); i++ {
		zbuffer[i] = math.Inf(-1)
	}

	// Transform coordinates
	model.transformCoordinates(eye, center, up) // Do projections and other transformations, update coordinates

	// Compute normals
	model.computeFaceNormals()
	model.computeVertexNormals()

	for i:=0; i<model.nFaces(); i++ {
		face := model.faces[i]
		faceTexture := model.faceTextures[i]

		var screenCoords [3]Vec3f
		var vertexNormals = make([]Vec3f, 3)
		var vertexTextures = make([]Vec2f, 3)
		faceNormal := model.faceNormals[i]
		for j:=0; j<3; j++ {
			vs := face[j]

			world_v := model.vertices[vs]
			world_v = world_v.normalizeCenteredCube(model) // Normalize coordinates to [0, 1] cube by min/max value
			screenCoords[j] = worldToScreen(&world_v, model, width, height, scale) // project to screen
			vertexNormals[j] = model.vertexNormals[vs]

			if len(faceTexture) > 0 {
				vts := faceTexture[j]
				vertexTextures[j] = model.textureCoordinates[vts]
			}
		}
		// render triangle
		triangle(&screenCoords[0], &screenCoords[1], &screenCoords[2], &vertexNormals, &vertexTextures, &faceNormal, lightDir, img, textureImage, &zbuffer, fillColor, width, height, specCoeff)
	}
}

// Shading 
// 1. Gouraud Shading

func gouraudShading(vertexNormals *[]Vec3f, lightDir *Vec3f, barycentric *Vec3f) float64{
	I1 := dot(&(*vertexNormals)[0], lightDir)
	I2 := dot(&(*vertexNormals)[1], lightDir)
	I3 := dot(&(*vertexNormals)[2], lightDir)
	I := I1 * barycentric.x + I2 * barycentric.y + I3 * barycentric.z 
	return I
}

// 2. Phong Shading

func phongShading(vertexNormals *[]Vec3f, lightDir *Vec3f, barycentric *Vec3f) (float64, float64){
	// Compute diffuse and spectral lighting

	n := Vec3f{}
	n1 := (*vertexNormals)[0].mul(barycentric.x, false)
	n.add(&n1, true)
	n2 := (*vertexNormals)[1].mul(barycentric.y, false)
	n.add(&n2, true)
	n3 := (*vertexNormals)[2].mul(barycentric.z, false)
	n.add(&n3, true)
	n.normalizeL2()
	
	diff := math.Max(0.0, dot(&n, lightDir)) // diffuse intensity

	// Reflected light
	r := n.mul(2 * dot(&n, lightDir), false)
	r.subtract(lightDir, true)

	spec := math.Pow(math.Max(0.0, -r.z), 10)

	return diff, spec
}

// 3. Flat Shading
func flatShading(faceNormal *Vec3f, lightDir *Vec3f) float64 {
	n := dot(faceNormal, lightDir)
	return n
}

// Utils
func newImage(height int, aspectRatio float64, fill bool) (*image.RGBA, int, int){
	width := int(aspectRatio * float64(height))
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img :=image.NewRGBA(image.Rectangle{upLeft, lowRight})
	if fill {
		background := color.RGBA{0, 0, 0, 255}
		draw.Draw(img, img.Bounds(), &image.Uniform{background}, image.ZP, draw.Src)
	}
	return img, width, height
}

// image utils
func getPixelValue(img *image.Image, x int, y int) (uint8, uint8, uint8) {
	r, g, b, _ := (*img).At(x, y).RGBA()
	return uint8(r / 257), uint8(g / 257), uint8(b / 257)
}

func getColorFromTexture(img *image.Image, vertexTextures *[]Vec2f, barycentric *Vec3f, diff float64, spec float64, specCoeff float64) color.RGBA {
	width, height := getImageSize(img)
	x := (*vertexTextures)[0].x * barycentric.x + (*vertexTextures)[1].x * barycentric.y + (*vertexTextures)[2].x * barycentric.z 
	y := (*vertexTextures)[0].y * barycentric.x + (*vertexTextures)[1].y * barycentric.y + (*vertexTextures)[2].y * barycentric.z 
	r, g, b := getPixelValue(img, int(x * float64(width)), int(y * float64(height)))
	return getColor(color.RGBA{r, g, b, 255}, diff, spec, specCoeff)
} 

func getColor(fillColor color.RGBA, diff float64, spec float64, specCoeff float64) color.RGBA {
	r, g, b := fillColor.R, fillColor.G, fillColor.B
	
	coeff := diff + specCoeff * spec

	r = uint8(math.Min(5 + float64(r) * coeff, 255))
	g = uint8(math.Min(5 + float64(g) * coeff, 255))
	b = uint8(math.Min(5 + float64(b) * coeff, 255))

	return color.RGBA{r, g, b, 255}
}

func getImageSize(img *image.Image)(int, int) {
	// returns image width, height
	bounds := (*img).Bounds()
	return bounds.Max.X, bounds.Max.Y
}

func main() {
	// Parse .obj file
	if len(os.Args) == 1 {
		fmt.Println("Specify input file!")
		os.Exit(1)
	}

	objPath := os.Args[1]
	fmt.Println("Using .obj file: ", objPath)
	model := parseObj(objPath)

	var texturePath string
	if len(os.Args) > 2 {
		texturePath = os.Args[2]
		fmt.Println("Using texture file: ", texturePath)
	} else {
		fmt.Println("No texture file specified. Using default color.")
	}

	file, _ := os.Open(texturePath)
	if file == nil {
		fmt.Println("Texture file does not exist. Using default color.")
	}
	textureImage, _, _ := image.Decode(file)

	// Report
	fmt.Println("Number of faces: ", model.nFaces())
	fmt.Println("Number of vertices: ", model.nVertices())
	fmt.Println("Number of textures coordinates: ", len(model.textureCoordinates))
	fmt.Println("Number of face textures: ", len(model.faceTextures))

	// Settings
	// ===========================================
	eye := newVec3f(-1, 0, -1)
	center := newVec3f(0, 0, 0)
	up := newVec3f(0, 1, 0)
	lightDir := newVec3f(0, 0, -1)
	specCoeff := 20.0
	imageHeight := 1000
	outFile := "./results/dragon.png"
	defaultFill := color.RGBA{218, 165, 32, 255}
	// ===========================================

	// Rendering
	ratio := model.aspectRatio()
	img, width, height := newImage(imageHeight, ratio, false)
	if textureImage == nil {
		renderTriangleMesh(
			&model, img, nil, &defaultFill, 
			&lightDir, &eye, &center, &up, width, height, 1.5, specCoeff)
	} else {
		renderTriangleMesh(
			&model, img, &textureImage, nil, 
			&lightDir, &eye, &center, &up, width, height, 1.5, specCoeff)
	}
	// Save
	f, _ := os.Create(outFile)
	png.Encode(f, imaging.FlipV(img))
}