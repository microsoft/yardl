function output_filename = invoke_translator(input_filename, input_format, output_format)
    translator_path = "../../cpp/build/translator";
    % output_filename = tempname();
    output_filename = sprintf("%s.translated", input_filename);
    dump = "cat";
    if ispc
        dump = "type";
    end
    cmd = sprintf("%s %s | %s %s %s > %s", dump, input_filename, translator_path, input_format, output_format, output_filename);

    [status, cmdout] = system(cmd);
    if status ~= 0
        error("Failed to invoke translator: %s", cmdout);
    end
end
