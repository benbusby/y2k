(define variable-types
  '((00 'none) (01 'string) (02 'numeric)))

(define-record-type :variable
  (make-variable size str val)
  variable?
  (size variable-size set-variable-size!)
  (str variable-str set-variable-str!)
  (val variable-val set-variable-val!))
