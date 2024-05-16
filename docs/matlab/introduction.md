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
data to a file in a compact binary format:

```matlab
addpath("./matlab/");

Sample = @sandbox.Sample;
samples = [Sample(yardl.DateTime.now(), [1, 2, 3]), Sample(yardl.DateTime.now(), [4, 5, 6])];

outfile = "sandbox.bin";

w = sandbox.binary.MyProtocolWriter(outfile);
w.write_header(sandbox.Header("Me"));
w.write_samples(samples);
w.end_samples();
w.close();
```

And then another script can read it in from the file:

```matlab
addpath("./matlab/");

infile = "sandbox.bin";

r = sandbox.binary.MyProtocolReader(infile);
disp(r.read_header());
while r.has_samples()
    sample = r.read_samples();
    disp(sample.timestamp.to_datetime());
    disp(sample.data);
end
r.close();
```

:::

<!--@include: ../parts/intro-core.md-->
