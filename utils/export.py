import argparse
from PIL import Image


parser = argparse.ArgumentParser()

parser.add_argument('-o')
parser.add_argument('-q')
parser.add_argument('images', type=str, nargs='+')

args = parser.parse_args()


def convert_to_rgb(img_rgba):
    try:
        img_rgba.load()
        _, _, _, alpha = img_rgba.split()

        img_rgb = Image.new('RGB', img_rgba.size, (255, 255, 255))
        img_rgb.paste(img_rgba, mask = alpha)

        return img_rgb
    except IndexError:
        return img_rgba


images = [
    convert_to_rgb(Image.open(image))
    for image in args.images
]


images[0].save(
    args.o, 'PDF', optimize = True, quality = int(args.q),
    save_all = True, append_images = images[1:]
)