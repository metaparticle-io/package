import subprocess

class DockerBuilder:
    def build(self, img):
        subprocess.run(['docker', 'build', '-t', img, '.'], check=True)

    def publish(self, img):
        subprocess.run(['docker', 'push', img], check=True)
