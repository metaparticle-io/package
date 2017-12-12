import subprocess

class DockerBuilder:
    def build(self, img):
        subprocess.call(['docker', 'build', '-t', img, '.'])

    def publish(self, img):
        subprocess.call(['docker', 'push', img])
