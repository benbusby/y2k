(module variables racket/base
  (provide (struct-out numeric))
  (struct numeric
          (digits
            [str #:mutable]
            [val #:mutable])
          #:transparent))
