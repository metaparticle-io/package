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
        subprocess.call(command)

    def logs(self, name):
        # Attach to logs
        # TODO: make this streaming...
        subprocess.call(['docker', 'logs', '-f', name])

    def cancel(self, name):
        subprocess.call(['docker', 'kill', name])
        subprocess.call(['docker', 'rm', name])
