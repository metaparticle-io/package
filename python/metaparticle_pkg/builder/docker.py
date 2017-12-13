import subprocess

class DockerBuilder:
    def build(self, img):
        subprocess.check_call(['docker', 'build', '-t', img, '.'])

    def publish(self, img):
        subprocess.check_call(['docker', 'push', img])
