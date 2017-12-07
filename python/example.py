from main import containerize


@containerize('willnewby', options={'name': 'testcontainer', 'publish': True})
def hihi():
    print('hello worldd')


if __name__ == '__main__':
    hihi()
