(define variables '())

(define variable-types
  '((00 'none) (01 'string) (02 'numeric)))

(define (divide num1 num2)
  (exact->inexact (/ num1 num2)))

(define variable-ops-numeric
  '((01 ,+) (02 ,-) (03 ,*) (04 ,divide)))

; TODO - Need to decide if string ops are worth it
;(define variable-ops-string
  ;'((00 ,string-append)))

(define variable-ops
  '((numeric ,variable-ops-numeric)))

(define (get-var-ops var)
    (eval (cadadr (assoc (variable-type var) variable-ops))))

(define-record-type :variable
  (make-variable size str val type)
  variable?
  (size variable-size set-variable-size!)
  (type variable-type set-variable-type!)
  (str variable-str set-variable-str!)
  (val variable-val set-variable-val!))
