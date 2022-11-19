(import
  (chicken file)
  (chicken file posix)
  (chicken irregex)
  (chicken process-context)
  (chicken sort)
  (chicken string))

(load "variables.scm")

(define get-ts file-modification-time)
(define input (list-ref (command-line-arguments) 0))
(define ext-exp (irregex "^.*.y2k$"))
(define printable " abcdefghijklmnopqrstuvwxyz!@#$%^&*()")
(define (displayln val) (display val) (display "\n"))

;; Variables for handling different interpreter modes
(define mode 'standby)
(define pause #f)
(define new-file #t)

;; Define mapping of program input->commands.
;; Keys intentionally include a leading 0 to visually
;; match how the interpreter will see it when parsing.
(define commands
  '(,(lambda () 'standby) ; ---- 00 : Reset
    ,(lambda () mode) ; -------- 01 : Continue from previous file
    ,(lambda () 'print-str) ; -- 02 : Print string to console
    ,(lambda () 'print-var) ; -- 03 : Print variable to console
    ,(lambda () 'set-var) ; ---- 04 : Create a new variable
    ,(lambda () 'mod-var))) ; -- 05 : Modify an existing variable

;; Establish a default variable that can be modified and stored as needed
;; throughout a program
(define new-var default-var)
(define (update-var val)
  (cond [(< (var-id new-var) 0)
         ; Set new variable "name"/id (just a 2-digit value to use
         ; when referencing the variable later in a program)
         (set-var-id! new-var val)]
        [(equal? (var-type new-var) 'none)
         ; Set the new variable data type
         ; TODO -- document types
         (set-var-type! new-var (cadadr (assoc val var-types)))]
        [(< (var-size new-var) 0)
         ; Use the two digit size provided here to determine how many digits
         ; or characters marks the variable as "complete".
         (set-var-size! new-var val)]
        [(not (= (string-length (var-str new-var)) (var-size new-var)))
         ; Convert values passed by the interpreter into appropriate values
         ; for the variable type.
         (cond [(equal? (var-type new-var) 'numeric)
                (for-each
                  (lambda (d)
                    (unless (var-str-comp new-var)
                      (set-var-str!
                        new-var
                        (string-append (var-str new-var) (string d)))))
                  (string->list (number->string val)))]
               [else (error "Unsupported data type")])

         ; Once the string length matches the specified variable size, we
         ; can finalize the variable value into a number (if numeric) or
         ; just copy the string over.
         (when (var-str-comp new-var)
           (cond [(equal? (var-type new-var) 'numeric)
                  (set-var-val! new-var (string->number (var-str new-var)))]
                 [(equal? (var-type new-var) 'string)
                  (set-var-val! new-var (var-str new-var))])
           (let ([new-var-id (var-id new-var)])
             (set! variables (append variables `((,new-var-id ,new-var))))
             (set! new-var default-var)
             (set! mode 'standby)))]))

;; Establish variables for tracking progress while
;; parsing timestamp values involved with modifying an
;; existing variable.
(define mod-target-var '())
(define mod-tmp-var '())
(define mod-var-step  0)
(define mod-with-other-var #f)
(define mod-operation '())
(define mod-var-steps
  '(,(lambda (val)
      ; Grab which variable the user is wanting to modify
      (set! mod-target-var (cadr (assoc val variables)))
      (inc mod-var-step))
    ,(lambda (val)
      ; Set the operator to use, depending on the type of the variable
      ; ("+" for numeric, "string-append" for strings)
      (set! mod-operation
        (eval (cadadr (assoc val (get-var-ops mod-target-var)))))
      (inc mod-var-step))
    ,(lambda (val)
      ; Check if another variable is being used in the calculation,
      ; or a raw number/string
      (set! mod-with-other-var (= val 01))
      (inc mod-var-step))
    ,(lambda (val)
      ; If another variable is being used, we can perform the
      ; calculation and return. Otherwise, we need to check the
      ; size of the item and continue parsing.
      (cond
        [mod-with-other-var
          ; Use the code from the timestamp to access the other
          ; variable
          (let [(other-var (cadr (assoc val variables)))]
            ; Update the original variable's value with the value
            ; of the second variable
            (set-variable-val!
              mod-target-var
              (mod-operation
                (variable-val mod-target-var)
                (variable-val other-var)))
            (reset #t))]
        [else
          ; Create a temporary variable to store the character codes/digits
          ; that will follow.
          (set! mod-tmp-var (make-variable val "" 0 (variable-type mod-target-var)))
          (inc mod-var-step)]))
    ,(lambda (val)
      ; Process single digit input
      (cond
        [(equal? (variable-type mod-target-var) 'numeric)
           (set-variable-str!
             mod-tmp-var
             (string-append (variable-str mod-tmp-var) (number->string val)))
           (cond
             [(= (string-length (variable-str mod-tmp-var)) (variable-size mod-tmp-var))
              (set-variable-val!
                mod-target-var
                (mod-operation
                  (variable-val mod-target-var)
                  (string->number (variable-str mod-tmp-var))))
              (reset #t)])]
        [else (error "Unsupported type")]))
    ))

(define (end-of-exp cmd)
  (and pause (= cmd 0)))

(define (y2k-dirlist input)
  (sort! (find-files input test: ext-exp limit: 0) string<?))

(define (parse-ts time-str index)
  (unless (> (+ index 2) (string-length time-str))
    (let ([cmd (string->number (substring time-str index (+ index 2)))])
      (cond
        ; Refer to the main commands table if a mode has not yet been set,
        ; an escape sequence was entered, or we're beginning to parse a new file.
        [(or (equal? mode 'standby) (end-of-exp cmd) new-file)
         (set! mode ((eval (cadr (list-ref commands cmd)))))
         (set! new-file #f)]
        [else
          (set! pause #f)
          (cond
            ; "print-str" mode prints a sequence of characters, accessed using
            ; 2-digit inputs, until commanded to stop ("0000") or the last file
            ; is reached.
            [(equal? mode 'print-str)
             (display (string-ref printable cmd))]
            ; "print-var" mode accepts a single 2-digit input that will be used
            ; to print the value of the variable that matches the requested 2-digit
            ; value.
            [(equal? mode 'print-var)
             (displayln (var-val (cadr (assoc cmd variables))))
             (set! mode 'standby)]
            ; "set-var" mode initiates a chain of processing steps that allow sequential
            ; 2-digit codes to set up and store a variable value. Refer to the documentation
            ; for more detail, but the tldr sequence is:
            ;
            ; set-var -> var name -> var type -> var size -> var value
            ;
            ; For example, to set variable "01" to 100, you would use the following timestamp
            ; 0401020310 -> 010...
            [(equal? mode 'set-var)
             (update-var cmd)]
            [(equal? mode 'mod-var)
             ((eval (cadr (list-ref mod-var-steps mod-var-step))) cmd)]
            [else (error "Unknown mode")])])
      (cond
        ; If we're using the default window size, we can assume that a "00"
        ; command should be interpreted as a pause in execution.
        [(= cmd 0)
         (set! pause #t)])
      (parse-ts time-str (+ index 2)))))

; Run the interpreter using the specified file or directory
(cond
  [(directory-exists? input)
   (for-each
     (lambda (f)
       ; Prepend a 0 to make an even number of digits that can easily be parsed.
       ; The first command of each file will always be <10, with 01 indicating the
       ; continuation of a previous file in order to not adversely affect
       ; multi-file operations.
       (define time-str (string-append "0" (number->string (get-ts f))))

       ; Parse the file and enable the new-file flag to indicate that
       ; any in-progress parsing (commands split between multiple files)
       ; should be paused until the next file is read.
       (parse-ts time-str 0)
       (set! new-file #t))
     (y2k-dirlist input))]
  [else (error "Dir not found")])
