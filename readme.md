# go-render

This is a toy project to learn Go and rendering at the same time.

<p align="center">
    <img width="400" src="./results/light3/light3.gif">
    <img width="400" src="./results/camera2/camera.gif">
</p>

The OG [tinyrenderer](https://github.com/ssloy/tinyrenderer) project helped me alot. It's some great stuff, and I highly recommend everyone check it out.

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

## `.obj` sources
- https://www.prinmath.com/csci5229/OBJ/index.html
- https://people.sc.fsu.edu/~jburkardt/data/obj/obj.html
- https://groups.csail.mit.edu/graphics/classes/6.837/F03/models/