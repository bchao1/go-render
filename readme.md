# go-render

A simple renderer written in Go.

<p align="center">
    <img width="400" src="./results/light3/light3.gif">
    <img width="400" src="./results/camera2/camera.gif">
</p>

<p align="center">
    <img width="800" src="./results/dragon.png">
</p>

The OG [tinyrenderer](https://github.com/ssloy/tinyrenderer) project helped me alot. It's amazing stuff, and I highly recommend everyone check it out.

## Usage
The code is pretty self-contained. I only used a 3rd-party library `imaging` to flip images vertically. 
<br>
<br>
To do your custom render:
```
go run render.go <path to .obj file> <path to texture file>
```
For example, see `run.sh`:
```
go run data/obj/bunny_2.obj data/textures/bunny_texture.jpg
```

Of course, you can first build `render.go`, and then run the executable:
```
go build render.go
./render <path to .obj file> <path to texture file>
```

### Customization
You could also play with some other parameters (light direction, camera position, spectral lighting, and etc) in `render.go`.
- `eye`: The camera position. Default is (0, 0, 1).
- `center`: Position the camera is looking at Default is (0, 0, 0).
- `up`: The vertical axis of the camera. Default is (0, 1, 0), which is the y-axis.
- `lightDir`: Light ray direction. Default is (0, 0, -1), which means the light is parellel to user's eye.
- `specCoeff`: Weight of spectral lighting. (For the shiny dragon above, I used 20.0)
- `imageHeight`: Pixel height of output image.
- `background`: Whether to color output image background black. If not, then output image has transparent background.
- `outFile`: Render output file path.
- `defaultFill`: If no texture is specified, use this color.

## Demo 

> Basics

|||
|--|--|
|Wireframe|Triangle rasterization|
|![img](./results/basic/wireframe.png)|![img](./results/basic/triangle_color.png)|

> Shading

|Flat shading|Gouraud shading|Phong shading|
|--|--|--|
|![img](./results/shading/flat.png)|![img](./results/shading/gouraud.png)|![img](./results/shading/phong.png)|
|![img](./results/shading/flat_detail.png)|![img](./results/shading/gouraud_detail.png)|![img](./results/shading/phong_detail.png)|

> Perspective

|||||
|--|--|--|--|
|![img](results/project/project_5.0.png)|![img](results/project/project_2.0.png)|![img](results//project/project_1.5.png)|![img](results/project/project_1.0.png)|

> Textures

Kudos to the author of [this article](https://blenderartists.org/t/uv-unwrapped-stanford-bunny-happy-spring-equinox/1101297) for providing custom Stanford bunny texture files.

|Colored|Terracotta|
|--|--|
|![img](./results/textures/bunny_color.png)|![img](./results/textures/bunny_terracotta.png)|

> Camera

|||||
|--|--|--|--|
|![img](./results/camera/1.png)|![img](./results/camera/2.png)|![img](./results/camera/3.png)|![img](./results/camera/4.png)|
|![img](./results/camera/8.png)|![img](./results/camera/7.png)|![img](./results/camera/6.png)|![img](./results/camera/5.png)|

> Light

|||||
|--|--|--|--|
|![img](./results/light/-10.png)|![img](./results/light/-5.png)|![img](./results/light/-2.png)|![img](./results/light/-1.png)|
![img](./results/light/10.png)|![img](./results/light/5.png)|![img](./results/light/2.png)|![img](./results/light/1.png)|

> Specular lighting

The stronger specular lighting is, the more "glossy" the object surface becomes. I simply used uniform power for each pixel since specular intensity is not specified in my texture files.
|No specular|Some specular (used)|Intense specular|
|--|--|--|
|![img](./results/specular/none.png)|![img](./results/specular/moderate.png)|![img](./results/specular/strong.png)|

## `.obj` sources
- https://www.prinmath.com/csci5229/OBJ/index.html
- https://people.sc.fsu.edu/~jburkardt/data/obj/obj.html
- https://groups.csail.mit.edu/graphics/classes/6.837/F03/models/
- https://casual-effects.com/data/
- https://github.com/alecjacobson/common-3d-test-models