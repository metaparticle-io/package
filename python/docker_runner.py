import os


def ports(options):
    try:
        port_string = ""
        print(options['ports'])
        for i in options['ports']:
            port_string += "-p {}:{}".format(i, i)
        return port_string
    except KeyError:
        return ' '


def run(img, name, options):
    # Launch docker container
    os.system('docker run --rm --name {} {} -d {}'.format(
        name,
        ports(options),
        img
    ))
    # Attach to logs
    os.system('docker logs -f {}'.format(name))


def cancel(name):
    os.system('docker kill {}'.format(name))
    os.system('docker rm {}'.format(name))
