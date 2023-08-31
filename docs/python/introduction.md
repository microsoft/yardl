# Introduction

Yardl is a simple schema language and command-line tool that generates domain
types and serialization code.

:::details Simple Example
Given a Yardl definition like this:

```yaml
# This is an example protocol, which is defined as a Header value
# followed by a stream of zero or more Sample values
MyProtocol: !protocol
  sequence:
    header: Header
    samples: !stream
      items: Sample

# Header is a record with a single string field
Header: !record
  fields:
    subject: string

# Sample is a record made up of a datetime and
# a vector of integers
Sample: !record
  fields:
    timestamp: datetime
    data: int*
```

After running `yardl generate`, you can write code like the following to write
data to standard out in a compact binary format:

```python
import sys
from sandbox import BinaryMyProtocolWriter, Header, Sample, DateTime

def generate_samples():
    yield Sample(timestamp=DateTime.now(), data=[1, 2, 3])
    yield Sample(timestamp=DateTime.now(), data=[4, 5, 6])

with BinaryMyProtocolWriter(sys.stdout.buffer) as w:
    w.write_header(Header(subject="Me"))
    w.write_samples(generate_samples())
```

And then another script can read it in from standard in:

```python
import sys
from sandbox import BinaryMyProtocolReader

with BinaryMyProtocolReader(sys.stdin.buffer) as r:
    print(r.read_header())
    for sample in r.read_samples():
        print(sample)
```

:::

<!--@include: ../parts/intro-core.md-->
