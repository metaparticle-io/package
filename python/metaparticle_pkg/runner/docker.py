import os
import subprocess

class DockerRunner:
    def run(self, img, name, options):
        # Launch docker container
        command = ['docker', 'run', '-d', '--name', name]
        if options.get('ports', None) is not None and len(options['ports']) > 0:
            command.append('-p')
            for i in options['ports']:
                command.append("{}:{}".format(i, i))
        command.append(img)
        subprocess.run(command, check=True)

    def logs(self, name):
        # Attach to logs
        # TODO: make this streaming...
        subprocess.run(['docker', 'logs', '-f', name], check=True)

    def cancel(self, name):
        subprocess.run(['docker', 'kill', name], check=True)
        subprocess.run(['docker', 'rm', name], check=True)
