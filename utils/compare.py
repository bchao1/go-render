import numpy as np
from PIL import Image
from matplotlib import pyplot as plt

for shading in ['flat', 'gouraud', 'phong']:
    img = np.asarray(Image.open('../results/shading/{}.png'.format(shading)))
    img = img[210:510, 150:450, :]
    Image.fromarray(img).save("../results/shading/{}_detail.png".format(shading))