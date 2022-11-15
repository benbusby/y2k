(module variables racket/base
  (provide variable-types)
  (provide (struct-out variable))

  ; Define table for variable types
  (define variable-types (make-hash))
  (hash-set! variable-types 00 'none)
  (hash-set! variable-types 01 'string)
  (hash-set! variable-types 02 'numeric)

  (struct variable
          (size
            [str #:mutable]
            [val #:mutable])
          #:transparent))
