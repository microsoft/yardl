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

    # A Header value
    header: Header

    # A stream of Samples
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

    # The time the sample was taken
    timestamp: datetime

    # A vector of integers
    data: int*
```

After running `yardl generate`, you can write code like the following to write
data to standard out in a compact binary format:

```cpp
#include <iostream>

#include "generated/binary/protocols.h"

int main() {
  playground::binary::MyProtocolWriter writer(std::cout);

  writer.WriteHeader({"123"});

  writer.WriteSamples({std::chrono::system_clock::now(), {1, 2, 3}});
  writer.WriteSamples({std::chrono::system_clock::now(), {4, 5, 6, 7}});

  // signal the end of the samples stream
  writer.EndSamples();
}
```

And then another executable can read it in from standard in:

```cpp
#include <iostream>

#include "generated/binary/protocols.h"

int main() {
  playground::binary::MyProtocolReader reader(filename);

  playground::Header header;
  reader.ReadHeader(header);

  std::cout << "Read Header.subject: " << header.subject << std::endl;

  playground::Sample sample;
  while (reader.ReadSamples(sample)) {
    std::cout << "Read Sample.data.size(): " << sample.data.size() << std::endl;
  }

  return 0;
}
```

:::

<!--@include: ../parts/intro-core.md-->
