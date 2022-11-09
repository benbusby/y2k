module Interpreter
    include("state.jl")
    using .State

    const time_ext = ".time"
    const space_ext = ".space"
    const p_chars = " abcdefghijklmnopqrstuvwxyz!@#\$%^&*()"

    const reserved = Dict(
        0 => State.reset_state,
        1 => State.start_print,
        2 => State.proceed,
        3 => State.define_function,
    )

    function run_space_interpreter(file::String)
        lines = readlines(file)
        for line in lines
            if State.pause[] || State.mode[] == State.none || length(line) == 0
                if !(length(line) in keys(reserved))
                    continue
                end
                reserved[length(line)]()
            elseif State.mode[] == State.print
                State.print_string[] = State.print_string[] * p_chars[length(line)]
            end
        end
    end

    function run_time_interpreter(file::String)
        unix_time::Float64 = Base.stat(file).mtime
        time_str::String = "0" * string(trunc(BigInt, unix_time))
        idx::Int32 = 1

        while idx + 1 < length(time_str) + 1
            window::Int32 = parse(Int32, time_str[idx:idx+1])
            if State.pause[] || State.mode[] == State.none || window == 0
                if !(window in keys(reserved))
                    continue
                end

                reserved[window]()
            elseif State.mode[] == State.print
                State.print_string[] = State.print_string[] * p_chars[window]
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
                State.pause[] = true
            end
        elseif isfile(path)
            process_file(path)
        else
            error("Invalid file or dir $path")
        end

    end

    if abspath(PROGRAM_FILE) == @__FILE__
        main()
    end

end
