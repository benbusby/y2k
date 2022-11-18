(import
  (chicken file)
  (chicken file posix)
  (chicken irregex)
  (chicken process-context)
  (chicken sort)
  (chicken string))

(load "variables.scm")
(load "utils.scm")

(define get-ts file-modification-time)
(define input (list-ref (command-line-arguments) 0))
(define ext-exp (irregex "^.*.y2k$"))
(define printable " abcdefghijklmnopqrstuvwxyz!@#$%^&*()")

;; Variables for handling different interpreter modes
(define mode 'none)
(define prev-mode 'none)
(define pause #f)
(define new-file #t)
(define w-size 2)
(define next-w-size 2)

(define (update-mode new-mode)
  (set! mode new-mode)
  (set! prev-mode new-mode))

(define (reset #!optional (hard-reset #f))
  ; Reset back to base mode
  (set! mode 'none)

  ; "Hard resets" should also ensure that the previous mode should be
  ; erased and the window size should be returned to its original size.
  (cond [hard-reset (set! prev-mode 'none) (set! next-w-size 2)]))

;; Define mapping of program input->commands.
;; Keys intentionally include a leading 0 to visually
;; match how the interpreter will see it when parsing.
(define commands
  '(,reset ; --------------------------------- 00 : Reset
    ,(lambda () (set! mode prev-mode)) ; ----- 01 : Continue from previous file
    ,(lambda () (update-mode 'print-str)) ; -- 02 : Print string to console
    ,(lambda () (update-mode 'print-var)) ; -- 03 : Print variable to console
    ,(lambda () (update-mode 'set-var)) ; ---- 04 : Create a new variable
    ,(lambda () (update-mode 'mod-var)))) ; -- 05 : Modify an existing variable

;; Establish variables for tracking progress while
;; parsing timestamp values involved with creating a
;; new variable.
(define new-var-step   0)
(define new-var-id    -1)
(define new-var-type '())
(define new-var-steps
  '(,(lambda (val)
      ; Set new variable "name" (a 2-digit value to use
      ; when referencing the variable in a program)
      (set! new-var-id val)
      (inc new-var-step))
    ,(lambda (val)
      ; Set the new variable data type
      ; TODO -- document types
      (set! new-var-type (cadadr (assoc val variable-types)))
      (inc new-var-step))
    ,(lambda (val)
      ; Use the two digit size provided here to create and
      ; insert the new variable.
      (let ([new-var (make-variable val "" 0 new-var-type)])
        (set! variables (append variables `((,new-var-id ,new-var))))
        (set! next-w-size 1)

      ; Now that we have the size, we can switch to single digit timestamp
      ; parsing until we've parsed the full number of digits for the variable.
      (inc new-var-step)))
    ,(lambda (val)
      (let ([new-var (cadr (assoc new-var-id variables))])
        (cond
          [(equal? (variable-type new-var) 'numeric)
           (set-variable-str!
             new-var
             (string-append (variable-str new-var) (number->string val)))
           (cond
             [(= (string-length (variable-str new-var)) (variable-size new-var))
              ; Convert the string value to a numeric value
              (set-variable-val! new-var (string->number (variable-str new-var)))

              ; Reset all values pertaining to the creation of new variables
              (set!-values (new-var-step new-var-id new-var-type) (values 0 -1 '()))
              (reset #t)])]
          [else (error "Unsupported variable type")])))))

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
          (set! next-w-size 1)
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
  (set! w-size next-w-size)
  (unless (> (+ index w-size) (string-length time-str))
    (let ([cmd (string->number (substring time-str index (+ index w-size)))])
      (cond
        ; Refer to the main commands table if a mode has not yet been set,
        ; an escape sequence was entered, or we're beginning to parse a new file.
        [(or (equal? mode 'none) (end-of-exp cmd) new-file)
         ((eval (cadr (list-ref commands cmd))))
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
             (displayln (variable-val (cadr (assoc cmd variables))))
             (reset)]
            ; "set-var" mode initiates a chain of processing steps that allow sequential
            ; 2-digit codes to set up and store a variable value. Refer to the documentation
            ; for more detail, but the tldr sequence is:
            ;
            ; set-var -> var name -> var type -> var size -> var value
            ;
            ; For example, to set variable "01" to 100, you would use the following timestamp
            ; 0401020310 -> 010...
            [(equal? mode 'set-var)
             ((eval (cadr (list-ref new-var-steps new-var-step))) cmd)]
            [(equal? mode 'mod-var)
             ((eval (cadr (list-ref mod-var-steps mod-var-step))) cmd)]
            [else (error "Unknown mode")])])
      (cond
        ; If we're using the default window size, we can assume that a "00"
        ; command should be interpreted as a pause in execution.
        [(and (> w-size 1) (= cmd 0))
         (set! pause #t)])
      (parse-ts time-str (+ index w-size)))))

; Run the interpreter using the specified file or directory
(cond
  [(directory-exists? input)
   (for-each
     (lambda (f)
       ; If the current window size is the default (2), prepend a 0 to make
       ; sure the initial command is always interpreted without throwing off
       ; the rest of the commands in the time string. Each file always starts
       ; with a 0X command.
       (define time-str (if (= w-size 2)
                          (string-append "0" (->string (get-ts f)))
                          (->string (get-ts f))))

       ; Parse the file and enable the new-file flag to indicate that
       ; any in-progress parsing (commands split between multiple files)
       ; should be paused until the next file is read.
       (parse-ts time-str 0)
       (set! new-file #t))
     (y2k-dirlist input))
   (reset)]
  [else (error "Dir not found")])
