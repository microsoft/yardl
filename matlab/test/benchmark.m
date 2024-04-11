% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

function benchmark(varargin)
    if nargin ~= 2
        error("Usage: benchmark <scenario> <format>");
    end

    scenario = parse_scenario(varargin{1});
    format = parse_format(varargin{2});

    % Matlab support currently implemented for binary format only
    if format ~= "binary"
        return
    end

    addpath("../generated/");
    fprintf("%s\n", scenario(format));
end

function output_filename = get_output_filename()
    output_filename = "/tmp/benchmark_data.dat";
end

function format = parse_format(format)
    switch format
        case "binary"
            format = "binary";
        case "ndjson"
            format = "ndjson";
        case "hdf5"
            format = "hdf5";
        otherwise
            error("Invalid format");
    end
end

function scenario = parse_scenario(scenario)
    switch scenario
        case "float256x256"
            scenario = @benchmark_float256x256;
        case "floatvlen"
            scenario = @benchmark_float_vlen;
        case "smallint256x256"
            scenario = @benchmark_small_int_256x256;
        case "smallrecord"
            scenario = @benchmark_small_record;
        case "smallrecordbatched"
            scenario = @benchmark_small_record_batched;
        case "smalloptionalsbatched"
            scenario = @benchmark_small_optionals_batched;
        case "simplemrd"
            scenario = @benchmark_simple_mrd;
        otherwise
            error("Invalid scenario");
    end
end

function [writer, reader] = create_writer_reader(format, protocol)
    format_name = "test_model." + format + "." + protocol;
    writer_name = format_name + "Writer";
    reader_name = format_name + "Reader";
    writer = str2func(writer_name);
    reader = str2func(reader_name);
end

function w = create_validating_writer(testCase, format, protocol)
    writer_name = "test_model." + format + "." + protocol + "Writer";
    reader_name = "test_model." + format + "." + protocol + "Reader";
    test_writer_name = "test_model.testing.Test" + protocol + "Writer";

    create_writer = str2func(writer_name);
    create_reader = str2func(reader_name);
    create_test_writer = str2func(test_writer_name);

    w = create_test_writer(testCase, format, @(f) create_writer(f), @(f) create_reader(f));
end

function reps = scale_repetitions(repetitions, scale)
    reps = repetitions * scale;
end

function res = time_scenario(total_bytes_size, write_func, read_func)
    delete(get_output_filename());

    total_size_mi_byte = total_bytes_size / 1024 / 1024;

    write_start = tic;
    write_func();
    write_elapsed = toc(write_start);
    write_mibps = total_size_mi_byte / write_elapsed;

    read_start = tic;
    read_func();
    read_elapsed = toc(read_start);
    read_mibps = total_size_mi_byte / read_elapsed;

    s.write_mi_bytes_per_second = write_mibps;
    s.read_mi_bytes_per_second = read_mibps;
    s.roundtrip_duration_seconds = write_elapsed + read_elapsed;
    res = jsonencode(s);
end

function res = benchmark_float256x256(format)
    scale = 1;

    a = zeros(256, 256, 'single');
    for i = 1:numel(a)
        a(i) = i - eps;
    end

    repetitions = scale_repetitions(10000, scale);
    total_bytes = 4 * numel(a) * repetitions;

    [create_writer, create_reader] = create_writer_reader(format, 'BenchmarkFloat256x256');

    function write()
        w = create_writer(get_output_filename());
        for r = 1:repetitions
            w.write_float256x256({a});
        end
        w.close();
    end

    function read()
        r = create_reader(get_output_filename());
        count = 0;
        while r.has_float256x256()
            b = r.read_float256x256();
            count = count + 1;
        end
        assert(count == repetitions);
        r.close();
    end

    res = time_scenario(total_bytes, @write, @read);
end

function res = benchmark_float_vlen(format)
    scale = 1;

    a = zeros(256, 256, 'single');
    for i = 1:numel(a)
        a(i) = i - eps;
    end

    repetitions = scale_repetitions(10000, scale);
    total_bytes = 4 * numel(a) * repetitions;

    [create_writer, create_reader] = create_writer_reader(format, 'BenchmarkFloatVlen');

    function write()
        w = create_writer(get_output_filename());
        for r = 1:repetitions
            w.write_float_array({a});
        end
        w.close();
    end

    function read()
        r = create_reader(get_output_filename());
        count = 0;
        while r.has_float_array()
            b = r.read_float_array();
            count = count + 1;
        end
        assert(count == repetitions);
        r.close();
    end

    res = time_scenario(total_bytes, @write, @read);
end

function res = benchmark_small_int_256x256(format)
    scale = 1;

    a = zeros(256, 256, 'int32');
    for i = 1:numel(a)
        a(i) = int32(37);
    end

    repetitions = scale_repetitions(10, scale);
    total_bytes = 4 * numel(a) * repetitions;

    [create_writer, create_reader] = create_writer_reader(format, 'BenchmarkInt256x256');

    function write()
        w = create_writer(get_output_filename());
        for r = 1:repetitions
            w.write_int256x256({a});
        end
        w.close();
    end

    function read()
        r = create_reader(get_output_filename());
        count = 0;
        while r.has_int256x256()
            b = r.read_int256x256();
            count = count + 1;
        end
        assert(count == repetitions);
        r.close();
    end

    res = time_scenario(total_bytes, @write, @read);
end

function res = benchmark_small_record(format)
    scale = 1;

    record = test_model.SmallBenchmarkRecord(double(73278383.23123213), single(78323.2820379), single(-2938923.29882));

    repetitions = scale_repetitions(50000, scale);
    total_bytes = 16 * repetitions;

    [create_writer, create_reader] = create_writer_reader(format, 'BenchmarkSmallRecord');

    function write()
        w = create_writer(get_output_filename());
        for r = 1:repetitions
            w.write_small_record({record});
        end
        w.close();
    end

    function read()
        r = create_reader(get_output_filename());
        count = 0;
        while r.has_small_record()
            b = r.read_small_record();
            count = count + 1;
        end
        assert(count == repetitions);
        r.close();
    end

    res = time_scenario(total_bytes, @write, @read);
end

function res = benchmark_small_record_batched(format)
    % Batching has not yet been implemented for Matlab
    res = "";
end

function res = benchmark_small_optionals_batched(format)
    % Batching has not yet been implemented for Matlab
    res = "";
end

function res = benchmark_simple_mrd(format)
    scale = 1;

    acq = test_model.SimpleAcquisition();
    acq.data = complex(zeros(256, 32, 'single'));
    acq.trajectory = zeros(2, 32, 'single');

    value = test_model.AcquisitionOrImage.Acquisition(acq);

    repetitions = scale_repetitions(2500, scale);
    total_bytes = 66032 * repetitions;

    [create_writer, create_reader] = create_writer_reader(format, 'BenchmarkSimpleMrd');

    function write()
        w = create_writer(get_output_filename());
        for r = 1:repetitions
            w.write_data({value});
        end
        w.close();
    end

    function read()
        r = create_reader(get_output_filename());
        count = 0;
        while r.has_data()
            b = r.read_data();
            count = count + 1;
        end
        assert(count == repetitions);
        r.close();
    end

    res = time_scenario(total_bytes, @write, @read);
end
