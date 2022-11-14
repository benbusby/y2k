#lang racket/base
(require racket/format)
(require racket/string)

(define get-ts file-or-directory-modify-seconds)

(define input (vector-ref (current-command-line-arguments) 0))
(define ext ".y2k")
(define printable " abcdefghijklmnopqrstuvwxyz!@#$%^&*()")
(define printstr "")

; Mutable "mode" for switching between interpreter states
(define mode (box 'none))
(define prev-mode (box 'none))
(define pause #f)
(define new-file #t)

; Define commands for modifying interpreter mode
(define (update-mode new-mode)
  (set-box! mode new-mode)
  (set-box! prev-mode new-mode))

(define (reset)
  (cond
    [(> (string-length printstr) 0)
     (displayln (string-trim printstr))
     (set! printstr "")])
  (update-mode 'none))
(define (begin-print)
  (update-mode 'begin-print))
(define (continue)
  (set-box! mode (unbox prev-mode)))

; Define table for mapping program input->commands
; Documented with leading 0 to match how the interpreter
; will see it when parsing.
(define commands (make-hash))
(hash-set! commands 00 reset)
(hash-set! commands 01 continue)
(hash-set! commands 02 begin-print)

(define count-substring (compose length regexp-match*))

(define (update-printstr idx)
  (let ([char (string-ref printable idx)])
    (set! printstr (string-append printstr (string char)))))

(define (end-of-exp line)
  (and pause (= (string-length line) 0)))

(define (end-of-input cmd)
  (and pause (= cmd 0)))

(define (sb-directory-list path)
  (filter
    (lambda (f)
      (string-suffix? (path->string f) ext))
    (directory-list path)))

(define (line-iter file)
  (let ((line (read-line file 'any)))
    (unless (eof-object? line)
      (cond
        [(or (equal? (unbox mode) 'none) (end-of-exp line))
         ((hash-ref commands (count-substring " " line)))]
        [else
          (set! pause #f)
          (cond
            [(equal? (unbox mode) 'begin-print)
             (update-printstr (count-substring " " line))]
            [else (error "Unknown mode")])])
      (cond
        [(= (string-length line) 0)
         (set! pause #t)])
      (line-iter file))
    (reset)))

(define (parse-ts time-str idx)
  (unless (= idx (string-length time-str))
    (define cmd (string->number (substring time-str idx (+ idx 2))))
    (cond
      [(or (equal? (unbox mode) 'none) (end-of-input cmd) new-file)
       ((hash-ref commands cmd))
       (set! new-file #f)]
      [else
        (set! pause #f)
        (cond
          [(equal? (unbox mode) 'begin-print)
           (update-printstr cmd)]
          [else (error "Unknown mode")])])
    (cond
      [(= cmd 0)
       (set! pause #t)])
    (parse-ts time-str (+ idx 2))))

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
       (sb-directory-list input))
     (reset)]
    [(file-exists? input)
     (line-iter (open-input-file input))]
    [else (error "File not found")])
