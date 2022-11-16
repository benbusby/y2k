(import
  (chicken file)
  (chicken file posix)
  (chicken irregex)
  (chicken process-context)
  (chicken sort)
  (chicken string))

(load "variables.scm")

(define get-ts file-modification-time)
(define (displayln val) (display val) (display "\n"))

(define input (list-ref (command-line-arguments) 0))
(define ext-exp (irregex "^.*.y2k$"))
(define printable " abcdefghijklmnopqrstuvwxyz!@#$%^&*()")
(define printstr "")
(define printvar -1)

; Mutable mode vars for switching between interpreter states
(define mode 'none)
(define prev-mode 'none)
(define pause #f)
(define new-file #t)
(define w-size 2)
(define next-w-size 2)

(define (update-mode new-mode)
  (set! mode new-mode)
  (set! prev-mode new-mode))

(define (reset)
  (cond
    [(> (string-length printstr) 0)
     (displayln printstr)]
    [(> printvar 0)
     (displayln (variable-val (cadr (assoc printvar variables))))
     (set! printvar -1)])

  (set! mode 'none))

(define (print-str) (update-mode 'print-str))
(define (print-var) (update-mode 'print-var))
(define (set-var)   (update-mode 'set-var))
(define (continue)  (set! mode prev-mode))

; Define mapping of program input->commands.
; Keys intentionally include a leading 0 to visually
; match how the interpreter will see it when parsing.
(define commands
  '((00 ,reset)
    (01 ,continue)
    (02 ,print-str)
    (03 ,print-var)
    (04 ,set-var)))

; Define a table for storing variables
(define variables '())

; Establish variables for tracking progress while
; parsing timestamp values involved with creating a
; new variable.
(define new-var-step   0)
(define new-var-id    -1)
(define new-var-type '())
(define new-var-steps
  '((00 ,(lambda (val)
          ; Set new variable "name" (a 2-digit value to use
          ; when referencing the variable in a program)
          (set! new-var-id val)
          (set! new-var-step (+ new-var-step 1))))
    (01 ,(lambda (val)
          ; Set the new variable data type
          ; TODO -- document types
          (set! new-var-type (cadadr (assoc val variable-types)))
          (set! new-var-step (+ new-var-step 1))))
    (02 ,(lambda (val)
          ; Use the two digit size provided here to create and
          ; insert the new variable.
          (let ([new-var (make-variable val "" 0)])
            (set! variables (append variables `((,new-var-id ,new-var))))
            (set! next-w-size 1))

          ; Now that we have the size, we can switch to single digit timestamp
          ; parsing until we've parsed the full number of digits for the variable.
          (set! new-var-step (+ new-var-step 1))))
    (03 ,(lambda (val)
          (let ([new-var (cadr (assoc new-var-id variables))])
            (cond
              [(equal? new-var-type 'numeric)
               (set-variable-str! new-var (string-append (variable-str new-var) (number->string val)))
               (cond
                 [(= (string-length (variable-str new-var)) (variable-size new-var))
                  ; Reset window size back to 2 digits
                  (set! next-w-size 2)

                  ; Convert the string value to a numeric value
                  (set-variable-val! new-var (string->number (variable-str new-var)))

                  ; Reset all values pertaining to the creation of new variables
                  (set!-values (new-var-step new-var-id new-var-type) (values 0 -1 '()))
                  (set! prev-mode 'none)
                  (reset)])]
              [else (error "Unsupported variable type")]))))))

(define (update-printstr index)
  (let ([char (string-ref printable index)])
    (set! printstr (string-append printstr (string char)))))

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
         ((eval (cadadr (assoc cmd commands))))
         (set! new-file #f)]
        [else
          (set! pause #f)
          (cond
            ; "print-str" mode prints a sequence of characters, accessed using
            ; 2-digit inputs, until commanded to stop ("0000") or the last file
            ; is reached.
            [(equal? mode 'print-str)
             (update-printstr cmd)]
            ; "print-var" mode accepts a single 2-digit input that will be used
            ; to print the value of the variable that matches the requested 2-digit
            ; value.
            [(equal? mode 'print-var)
             (set! printvar cmd)
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
             ((eval (cadadr (assoc new-var-step new-var-steps))) cmd)]
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
