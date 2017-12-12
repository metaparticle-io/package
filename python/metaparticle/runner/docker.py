import os


def ports(ports):
    try:
        port_string = ""
        for port in ports:
            port_string += " -p {port}:{port}".format(port=port)
        return port_string
    except KeyError:
        return ' '


def run(img, name, options):
    # Launch docker container
    os.system('docker run --rm --name {} {} -d {}'.format(
        name,
        ports(options.ports),
        img
    ))
    # Attach to logs
    os.system('docker logs -f {}'.format(name))


def cancel(name):
    os.system('docker kill {}'.format(name))
    os.system('docker rm {}'.format(name))
