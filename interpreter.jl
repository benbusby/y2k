module Interpreter

    @enum State none print def_fn
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

    string = " abcdefghijklmnopqrstuvwxyz!@#\$%^&*()"

    reserved = Dict(
        0 => reset_state,
        1 => start_print,
        2 => define_function
    )

    function main()
        file = ARGS[1]
        if !isfile(file)
            error("could not find file $file")
        end

        lines = readlines(file)
        for line in lines
            if current_state[] == none || length(line) == 0
                if !(length(line) in keys(reserved))
                    continue
                end
                reserved[length(line)]()
            elseif current_state[] == print
                print_string[] = print_string[] * string[length(line)]
            end
        end
    end

    if abspath(PROGRAM_FILE) == @__FILE__
        main()
    end

end #module
