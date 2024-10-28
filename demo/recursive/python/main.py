import os
import sys
import sketch
import random
from io import BytesIO


def read_directory(path):
    # os.path.
    for dirent in os.scandir(path):
        if dirent.is_dir():
            d = sketch.Directory(
                name=dirent.name, entries=list(read_directory(dirent.path))
            )
            yield sketch.DirectoryEntry.Directory(d)
        elif dirent.is_file():
            f = sketch.File(name=dirent.name)
            yield sketch.DirectoryEntry.File(f)
        else:
            print("Ignoring special file ", dirent.name, file=sys.stderr)


def writeDirectoryEntry(dirent, path):
    if isinstance(dirent, sketch.DirectoryEntry.File):
        path = os.path.join(path, dirent.value.name)
        print(path)
    elif isinstance(dirent, sketch.DirectoryEntry.Directory):
        path = os.path.join(path, dirent.value.name)
        print(path)
        for d in dirent.value.entries:
            writeDirectoryEntry(d, path)
    else:
        print("Ignoring special file ", dirent.name, file=sys.stderr)


def insertTree(root, value):
    if root is None:
        return sketch.BinaryTree(value=value)

    if value < root.value:
        root.left = insertTree(root.left, value)
    else:
        root.right = insertTree(root.right, value)

    return root


def walkTree(root, fn):
    if root is None:
        return
    walkTree(root.left, fn)
    fn(root.value)
    walkTree(root.right, fn)


def main():
    root = sketch.BinaryTree()
    for i in range(32):
        root = insertTree(root, random.randint(-50, 50))
    walkTree(root, lambda x: print(x, end=" "))
    print()

    llist = sketch.LinkedList(value="Hello")
    llist.next = sketch.LinkedList(value="World")
    llist.next.next = sketch.LinkedList(value="!!!")

    stream = BytesIO()
    with sketch.BinaryMyProtocolWriter(stream) as writer:
        writer.write_tree(root)
        writer.write_ptree(root)

        writer.write_list(llist)

        writer.write_cwd(read_directory("."))

        writer.close()

    # print(len(stream.getvalue()))
    # print(str(stream.getvalue()))

    stream = BytesIO(stream.getvalue())
    with sketch.BinaryMyProtocolReader(stream) as reader:
        walkTree(reader.read_tree(), lambda x: print(x, end=" "))
        print()
        walkTree(reader.read_ptree(), lambda x: print(x, end=" "))
        print()

        head = reader.read_list()
        while head is not None:
            print(head.value, end=" ")
            head = head.next
        print()

        for dirent in reader.read_cwd():
            writeDirectoryEntry(dirent, ".")

        reader.close()

    print("Finished!")


if __name__ == "__main__":
    main()
