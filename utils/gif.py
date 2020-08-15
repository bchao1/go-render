import os
from PIL import Image

img_list = []
folder = '../results/light3/'

for i in range(len(os.listdir(folder))):
    path = folder + "{}.png".format(i)
    img_list.append(Image.open(path))
img_list += img_list[::-1]
img_list[0].save('light3.gif', save_all=True, append_images=img_list[1:], duration=100, loop=0)