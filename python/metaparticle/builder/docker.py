import os


def build(img):
    os.system('docker build -t {} .'.format(img))


def publish(img):
    os.system('docker push {}'.format(img))
