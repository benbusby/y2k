(import srfi-69)
(import
  srfi-69
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
     (displayln (variable-val (hash-table-ref variables printvar)))
     (set! printvar -1)])

  (set! mode 'none))

(define (print-str) (update-mode 'print-str))
(define (print-var) (update-mode 'print-var))
(define (set-var)   (update-mode 'set-var))
(define (continue)  (set! mode prev-mode))

; Define table for mapping program input->commands
; Keys intentionally include a leading 0 to visually
; match how the interpreter will see it when parsing.
(define commands (make-hash-table))
(hash-table-set! commands 00 reset)
(hash-table-set! commands 01 continue)
(hash-table-set! commands 02 print-str)
(hash-table-set! commands 03 print-var)
(hash-table-set! commands 04 set-var)

; Define a table for storing variables
(define variables (make-hash-table))

; Establish variables for tracking progress while
; parsing timestamp values involved with creating a
; new variable.
(define new-var-step   0)
(define new-var-id    -1)
(define new-var-type '())
(define new-var-steps (make-hash-table))
(hash-table-set! new-var-steps 00 (lambda (val)
  ; Set new variable "name" (a 2-digit value to use
  ; when referencing the variable in a program)
  (set! new-var-id val)
  (set! new-var-step (+ new-var-step 1))))
(hash-table-set! new-var-steps 01 (lambda (val)
  ; Set the new variable data type
  ; TODO -- document types
  (set! new-var-type (hash-table-ref variable-types val))
  (set! new-var-step (+ new-var-step 1))))
(hash-table-set! new-var-steps 02 (lambda (val)
  ; Use the two digit size provided here to create and
  ; insert the new variable.
  (let ([new-var (make-variable val "" 0)])
    (hash-table-set! variables new-var-id new-var)
    (set! next-w-size 1))

  ; Now that we have the size, we can switch to single digit timestamp parsing
  ; until we've parsed the full number of digits for the variable.
  (set! new-var-step (+ new-var-step 1))))
(hash-table-set! new-var-steps 03 (lambda (val)
  (let ([new-var (hash-table-ref variables new-var-id)])
    (cond
      [(equal? new-var-type 'numeric)
       (set-variable-str! new-var (string-append (variable-str new-var) (->string val)))
       (cond
         [(= (string-length (variable-str new-var)) (variable-size new-var))
          ; Reset window size back to 2 digits
          (set! next-w-size 2)

          ; Convert the string value to a numeric value
          (set-variable-val! new-var (string->number (variable-str new-var)))

          ; Reset all values pertaining to the creation of new variables
          (set!-values
            (new-var-step new-var-id new-var-type)
            (values 0 -1 '()))
          (reset)])]))))

(define (update-printstr index)
  (let ([char (string-ref printable index)])
    (set! printstr (string-append printstr (string char)))))

(define (end-of-exp cmd)
  (and pause (= cmd 0)))

(define (y2k-dirlist input)
  (sort! (find-files input test: ext-exp limit: 0) string<?))

(define (parse-ts time-str index)
  (unless (= index (string-length time-str))
    (let ()
      (set! w-size (if (not (= w-size next-w-size)) next-w-size w-size))
      (define cmd (string->number (substring time-str index (+ index w-size))))
      (cond
        ; Refer to the main commands table if a mode has not yet been set,
        ; an escape sequence was entered, or we're beginning to parse a new file.
        [(or (equal? mode 'none) (end-of-exp cmd) new-file)
         ((hash-table-ref commands cmd))
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
             ((hash-table-ref new-var-steps new-var-step) cmd)]
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
