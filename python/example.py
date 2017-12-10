from metaparticle import containerize


@containerize('brendanburns', options={'ports': ['80', '8080']})
def container_with_port():
    print('hello container_with_port')


@containerize(
    'brendanburns', options={'name': 'testcontainer', 'publish': True})
def hihi():
    print('hello worldd')


if __name__ == '__main__':
    hihi()
    container_with_port()
