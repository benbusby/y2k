module State
    @enum Mode none print def_fn
    const mode::Base.RefValue{Mode} = Ref(none)
    const print_string::Base.RefValue{String} = Ref("")
    const pause::Base.RefValue{Bool} = Ref(false)

    function reset_state()
        if mode[] == print
            println(print_string[])
        end
        mode[] = none
    end

    function start_print()
        mode[] = print
    end

    function define_function()
        mode[] = def_fn
    end

    function proceed()
        # Release pause before proceeding to next file
        pause[] = false
    end
end
