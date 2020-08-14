from PIL import Image

img_list = []

for i in range(-10, 11):
    path = "../results/light/{}.png".format(i)
    img_list.append(Image.open(path))
img_list += img_list[::-1]
img_list[0].save('light.gif', save_all=True, append_images=img_list[1:], duration=100, loop=0)