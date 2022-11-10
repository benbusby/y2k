#lang racket/base
(require racket/string)

(define input (vector-ref (current-command-line-arguments) 0))
(define printable " abcdefghijklmnopqrstuvwxyz!@#$%^&*()")
(define printstr "")

; Mutable "mode" for switching between interpreter states
(define mode (box 'none))
(define pause #f)

; Define commands for modifying interpreter mode
(define (reset)
  (cond
    [(> (string-length printstr) 0)
     (displayln (string-trim printstr))
     (set! printstr "")])
  (set-box! mode 'none))
(define (begin-print)
  (set-box! mode 'begin-print))

; Define table for mapping program input->commands
(define commands (make-hash))
(hash-set! commands 0 reset)
(hash-set! commands 1 begin-print)

(define count-substring (compose length regexp-match*))

(define (update-printstr line)
  (let ([char (string-ref printable (count-substring " " line))])
    (set! printstr (string-append printstr (string char)))))

(define (end-of-exp line)
  (and pause (= (string-length line) 0)))

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
             (update-printstr line)]
            [else (error "Unknown mode")])])
      (cond
        [(= (string-length line) 0)
         (set! pause #t)])
      (line-iter file))
    (reset)))

; Run the interpreter against the specified file or directory
(cond
    [(directory-exists? input) (displayln "dir")]
    [(file-exists? input)
      (line-iter (open-input-file input))]
    [else (error "File not found")])
