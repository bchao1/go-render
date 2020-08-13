import numpy as np
from PIL import Image
from matplotlib import pyplot as plt


shading = 'phong'

img = np.asarray(Image.open('../results/shading/{}.png'.format(shading)))
img = img[150:450, 150:450, :]
Image.fromarray(img).save("{}_detail.png".format(shading))