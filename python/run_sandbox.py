#! /usr/bin/env python3

import os
import sandbox
import numpy as np


def main():
    with sandbox.BinaryP1Writer("test.bin") as w:
        # w.write_an_int(8)
        # w.write_a_stream([1, 2, 3, 4, 5])
        # w.write_optional(2)
        # w.write_union(None)
        # w.write_flag(sandbox.MyFlags.A | sandbox.MyFlags.B)
        # w.write_vec([1,2,3])
        # new_var = np.array([[1,2,3],[4,5,6]], dtype=np.uint32)
        # w.write_arr(new_var)
        # w.write_map({"a": 1, "b": 2})
        # w.write_point(sandbox.Point(x=1, y=2))
        # w.write_points([sandbox.Point(x=1, y=2), sandbox.Point(x=3, y=4)])
        # arr = np.array([(1, 2), (3, 4)], dtype=[('x', np.float32), ('y', np.float32)])
        # w.write_points(arr)
        # w.write_gen_rec(sandbox.MyRec(f1=2, f2=22.3, f3=sandbox.MyFlags.A | sandbox.MyFlags.B))
        # dt = np.dtype([('f1', np.int8), ('f2', np.float32), ('f3', np.uint8)])
        # arr = np.array([(2, 22.3, 3)], dtype=dt)
        # w.write_gen_rec(arr)
        # w.write_myint(np.int32(2))
        # w._write_image([1,2,3])
        # w.write_intimage([1,2,3])

        dtype = np.dtype([('points', [('x', '<i4'), ('y', '<i4')], (2,))])

        arr = np.array([([(1.1, 2.2),(3.3, 4.4)],)], dtype=dtype)
        w.write_complicated_arr(arr)

        pass


    os.system("hexdump -C test.bin")

    s = "asd"
    print(s[0])


if __name__ == "__main__":
    main()
