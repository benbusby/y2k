#lang racket/base
(require racket/format)
(require racket/string)
(require racket/include)

(require "variables.rkt")

(define get-ts file-or-directory-modify-seconds)

(define input (vector-ref (current-command-line-arguments) 0))
(define ext ".y2k")
(define printable " abcdefghijklmnopqrstuvwxyz!@#$%^&*()")
(define printstr "")
(define printvar -1)

; Mutable "mode" for switching between interpreter states
(define mode (box 'none))
(define prev-mode (box 'none))
(define pause #f)
(define new-file #t)
(define w-size 2)
(define next-w-size 2)

; Define commands for modifying interpreter mode
(define (update-mode new-mode)
  (set-box! mode new-mode)
  (set-box! prev-mode new-mode))

(define (reset)
  (cond
    [(> (string-length printstr) 0)
     (displayln (string-trim printstr))
     (set! printstr "")]
    [(> printvar 0)
     (displayln (numeric-val (hash-ref variables printvar)))
     (set! printvar -1)])

  (set-box! mode 'none))
(define (print-str) (update-mode 'print-str))
(define (print-var) (update-mode 'print-var))
(define (set-var)   (update-mode 'set-var))
(define (continue)  (set-box! mode (unbox prev-mode)))

; Define table for mapping program input->commands
; Documented with leading 0 to match how the interpreter
; will see it when parsing.
(define commands (make-hash))
(hash-set! commands 00 reset)
(hash-set! commands 01 continue)
(hash-set! commands 02 print-str)
(hash-set! commands 03 print-var)
(hash-set! commands 04 set-var)

; Define table for storing variables
(define variables (make-hash))

; Define table for variable types
(define variable-types (make-hash))
(hash-set! variable-types 00 'none)
(hash-set! variable-types 01 'string)
(hash-set! variable-types 02 numeric)

; Set up values for tracking progress while parsing
; the timestamp to declare a new variable.
(define new-var-step 0)
(define new-var-id   -1)
(define new-var-type '())
(define new-var-steps (make-hash))
(hash-set! new-var-steps 00 (lambda (val)
  ; Set new variable "name" (a 2-digit code to use when
  ; referring to the variable)
  (set! new-var-id val)
  (set! new-var-step (+ new-var-step 1))))
(hash-set! new-var-steps 01 (lambda (val)
  ; Set the new variable data type --
  (set! new-var-type (hash-ref variable-types val))
  (set! new-var-step (+ new-var-step 1))))
(hash-set! new-var-steps 02 (lambda (val)
  (cond
    [(equal? (object-name new-var-type) 'numeric)
     (let ([new-var (new-var-type val "" 0)])
       (hash-set! variables new-var-id new-var)
       (set! next-w-size 1))])
  (set! new-var-step (+ new-var-step 1))))
(hash-set! new-var-steps 03 (lambda (val)
  (let ([new-var (hash-ref variables new-var-id)])
    (cond
      [(equal? (object-name new-var-type) 'numeric)
       (cond
         [(= (string-length (numeric-str new-var)) (numeric-digits new-var))
          ; Reset window size back to 2 digits
          (set! next-w-size 2)

          ; Convert the string value to a numeric value and store within the
          ; variable struct
          (set-numeric-val! new-var (string->number (numeric-str new-var)))

          ; Reset all values related creating new variables
          (set!-values
            (new-var-step new-var-id new-var-type)
            (values 0 -1 '()))
          (reset)]
         [else
           (set-numeric-str! new-var (string-append (numeric-str new-var) (~v val)))])
       ]))))

(define (update-printstr idx)
  (let ([char (string-ref printable idx)])
    (set! printstr (string-append printstr (string char)))))

(define (end-of-exp cmd)
  (and pause (= cmd 0)))

(define (y2k-dirlist path)
  (filter
    (lambda (f)
      (string-suffix? (path->string f) ext))
    (directory-list path)))

(define (parse-ts time-str idx)
  (unless (= idx (string-length time-str))
    (set! w-size (if (not (= w-size next-w-size)) next-w-size w-size))
    (define cmd (string->number (substring time-str idx (+ idx w-size))))
    (cond
      ; Refer back to the main commands table if a mode has not been
      ; set, an escape sequence was entered, or we're starting a new
      ; file.
      [(or (equal? (unbox mode) 'none) (end-of-exp cmd) new-file)
       ((hash-ref commands cmd))
       (set! new-file #f)]
      [else
        (set! pause #f)
        (cond
          ; "print-str" mode prints a sequence of characters, accessed using
          ; 2-digit commands, until commanded to stop ("0000") or the last file
          ; is reached.
          [(equal? (unbox mode) 'print-str)
           (update-printstr cmd)]
          ; "print-var" mode accepts a single 2-digit command that will be used
          ; to print the value of the variable that matches the requested 2-digit
          ; value.
          [(equal? (unbox mode) 'print-var)
           (set! printvar cmd)
           (reset)]
          ; "set-var" mode initiates a chain of processing steps that allow sequential
          ; 2-digit codes to set up and store a numeric value. Refer to the documentation
          ; for more detail, but the tldr sequence is:
          ;
          ; set-var -> var name -> var type -> var size -> var value
          ;
          ; For example, to set variable "01" to 100, you would use the following timestamp
          ; 0401020310 -> 010...
          [(equal? (unbox mode) 'set-var)
           ((hash-ref new-var-steps new-var-step) cmd)]
          [else (error "Unknown mode")])])
    (cond
      ; If we're using the default window size, we can assume
      ; that a 0 or "00" command should be interpreted as a
      ; pause in execution.
      [(and (> w-size 1) (= cmd 0))
       (set! pause #t)])
    (parse-ts time-str (+ idx w-size))))

; Run the interpreter against the specified file or directory
(cond
    [(directory-exists? input)
     (for-each
       (lambda (f)
         ; Build path to be absolute ("path/to/file.y2k")
         (define path (build-path input f))

         ; Prepend a 0 to ensure initial command is interpreted
         ; without disturbing the rest of the time string.
         ; Each file always starts with a single digit command.
         (define time-str (string-append "0" (~v (get-ts path))))

         ; Parse the file and enable the new-file flag to indicate that
         ; any in-progress parsing (commands split between multiple files)
         ; should be paused until the next file is read.
         (parse-ts time-str 0)
         (set! new-file #t))
       (y2k-dirlist input))
     (reset)]
    ; TODO: Should single files be supported?
    ;[(file-exists? input)
     ;(line-iter (open-input-file input))]
    [else (error "File not found")])
