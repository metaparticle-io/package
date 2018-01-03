from metaparticle_pkg import Containerize
import time

package_repo = 'repo'
package_name = 'something'
sleep_time = 10


@Containerize(
    package={'name': package_name, 'repository': package_repo},
    runtime={'ports': [80, 8080]}
)
def container_with_port():
    print('hello container_with_port')

    for i in range(sleep_time):
        print('Sleeping ... {} sec'.format(i))
        time.sleep(1)


@Containerize(
    package={'name': package_name, 'repository': package_repo, 'publish': True}
)
def hihi():
    print('hello world!')

    for i in range(sleep_time):
        print('Sleeping ... {} sec'.format(i))
        time.sleep(1)


if __name__ == '__main__':
    hihi()
    container_with_port()
