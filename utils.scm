(define (displayln val) (display val) (display "\n"))
(define-syntax inc
  (syntax-rules ()
    ((_ x) (begin (set! x (+ x 1)) x))))
