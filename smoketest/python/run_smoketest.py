#! /usr/bin/env python3

import sys

from smoketest import NDJsonMyProtocolWriter, Header, Sample, DateTime

with NDJsonMyProtocolWriter(sys.stdout) as w:
    w.write_header(Header(subject="me"))
    w.write_samples(
        [
            Sample(timestamp=DateTime.now(), data=[1, 2, 3]),
            Sample(timestamp=DateTime.now(), data=[4, 5, 6]),
        ]
    )
