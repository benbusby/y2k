(define variable-types (make-hash-table))
(hash-table-set! variable-types 00 'none)
(hash-table-set! variable-types 01 'string)
(hash-table-set! variable-types 02 'numeric)

(define-record-type :variable
  (make-variable size str val)
  variable?
  (size variable-size set-variable-size!)
  (str variable-str set-variable-str!)
  (val variable-val set-variable-val!))
