addpath("./generated");

filename = "sandbox_mat.bin";

data = [ ...
    [892.37889483 - 9932.485937837j; 73.383672763878 - 33.3394472537j], ...
    [3883.22890980 + 373.4933837j; 56985.39384393 - 33833.3330128474373j], ...
    [283.383672763878 - 33.3394472537j; 3883.22890980 + 373.4933837j], ...
];

w = sandbox.binary.HelloWorldWriter(filename);
w.write_data(data);
w.end_data();
w.close();

r = sandbox.binary.HelloWorldReader(filename);
while r.has_data()
    values = r.read_data();
    for n = 1:length(values)
        disp(values(n));
    end
end
r.close();

system(sprintf("hexdump -C %s", filename));
