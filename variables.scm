(define variables '())

(define var-types
  '((00 'none) (01 'string) (02 'numeric)))

(define (divide num1 num2)
  (exact->inexact (/ num1 num2)))

(define var-ops-numeric
  '((01 ,+) (02 ,-) (03 ,*) (04 ,divide)))

; TODO - Need to decide if string ops are worth it
;(define variable-ops-string
  ;'((00 ,string-append)))

(define var-ops
  '((numeric ,var-ops-numeric)))

(define (get-var-ops var)
    (eval (cadadr (assoc (var-type var) var-ops))))

(define-record-type :var
  (make-var id size str val type)
  var?
  (id var-id set-var-id!)
  (size var-size set-var-size!)
  (type var-type set-var-type!)
  (str var-str set-var-str!)
  (val var-val set-var-val!))

(define (var-str-comp t-var)
  (= (string-length (var-str t-var)) (var-size t-var)))

(define default-var
  (make-var -1 -1 "" -1 'none))
