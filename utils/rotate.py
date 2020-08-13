from PIL import Image, ImageOps

img = Image.open('../data/textures/bunny_texture2.jpg')
img = ImageOps.flip(img)
img.save('../data/textures/bunny_texture2.jpg')