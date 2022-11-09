module Interpreter

    @enum State none print def_fn
    const time_ext = ".time"
    const space_ext = ".space"
    const current_state::Base.RefValue{State} = Ref(none)
    const print_string::Base.RefValue{String} = Ref("")
    const pause::Base.RefValue{Bool} = Ref(false)

    function reset_state()
        if current_state[] == print
            println(print_string[])
        end
        current_state[] = none
    end

    function start_print()
        current_state[] = print
    end

    function define_function()
        current_state[] = def_fn
    end

    function proceed()
        # Release pause before proceeding to next file
        pause[] = false
    end

    p_chars = " abcdefghijklmnopqrstuvwxyz!@#\$%^&*()"

    reserved = Dict(
        0 => reset_state,
        1 => start_print,
        2 => proceed,
        3 => define_function,
    )

    function run_space_interpreter(file::String)
        lines = readlines(file)
        for line in lines
            if pause[] || current_state[] == none || length(line) == 0
                if !(length(line) in keys(reserved))
                    continue
                end
                reserved[length(line)]()
            elseif current_state[] == print
                print_string[] = print_string[] * p_chars[length(line)]
            end
        end
    end

    function run_time_interpreter(file::String)
        unix_time::Float64 = Base.stat(file).mtime
        time_str::String = "0" * string(trunc(BigInt, unix_time))
        idx::Int32 = 1

        while idx + 1 < length(time_str) + 1
            window::Int32 = parse(Int32, time_str[idx:idx+1])
            if pause[] || current_state[] == none || window == 0
                if !(window in keys(reserved))
                    continue
                end

                reserved[window]()
            elseif current_state[] == print
                print_string[] = print_string[] * p_chars[window]
            end

            idx += 2
        end
    end

    function process_file(file::String)
        file_ext::String = last(Base.Filesystem.splitext(file))

        if file_ext == time_ext
            run_time_interpreter(file)
        elseif file_ext == space_ext
            run_space_interpreter(file)
        end
    end

    function main()
        path = ARGS[1]
        if isdir(path)
            files::Vector{String} = sort(readdir(path))
            for file in files
                process_file(path * file)
                pause[] = true
            end
        elseif isfile(path)
            process_file(file)
        else
            error("Invalid file or dir $path")
        end

    end

    if abspath(PROGRAM_FILE) == @__FILE__
        main()
    end

end #module
