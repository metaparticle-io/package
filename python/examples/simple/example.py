from metaparticle_pkg import Containerize


@Containerize(package={'name': 'something', 'repository': 'repo'}, runtime={'ports': [80, 8080]})
def container_with_port():
    print('hello container_with_port')


@Containerize(package={'name': 'something', 'repository': 'repo', 'publish': True})
def hihi():
    print('hello world!')


if __name__ == '__main__':
    hihi()
    container_with_port()
