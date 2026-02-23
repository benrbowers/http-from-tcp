# HTTP From TCP boot.dev course

In this repo, I've completed [this boot.dev course](https://www.boot.dev/courses/learn-http-protocol-golang).

The purpose of this course was to use the Go programming language build a simple HTTP/1.1 parser and server: **starting from TCP.**

I had an absolute blast completing this course, and I would recommend it to anyone wanting to diver deeper into Go, the HTTP protocol, interfaces, and/or testing.

## What I've Learned
- General shape an HTTP/1.1 message
  - [request-line (if request)/status-line (if response)] CRLF
  - [field-line] CRLF x [number of headers]
    - field-line = field-value: field-name
  - CRLF
  - [message-body]
- Finally learned what CRLF stands for: "Carriage Return, Line Feed"
  - Good way to remember "\r\n" in your code is "Registered Nurse"
- Chunked encoding
  - By default, an HTTP server expects an HTTP message to include the exact length **in bytes** of the body up front in the "Content-Length" header.
    - If the body is off by even a single byte, an error is thrown.
  - With _chunked encoding,_ the size of the body isn't known up front, and so it gets sent down in **chunks.**
  - Each "chunk-data" block is preceded by "chunk-size" line, and much like Content-Length, they must match exactly.
  - Generally, after the body has been completely streamed down and the exact content length is known, it will be appended as a **trailer.** (Usually as "X-Content-Length")
    - Trailers are only present in chunked responses, and serve the _exact same purpose as headers,_ except that they come at the very end.
    - Trailers must be announced up front in a comma separated list using the "Trailer" header.
  - Users of HTTP can enable chunked encoding simply by omitting the "Content-Length" header and then including "Transfer-Encoding: chunked".
    - It is compatible with most if not all of the different "Content-Type" options. You can stream pretty much anything you can turn into bytes.
- [io](https://pkg.go.dev/io) package interfaces
  - Simply by writing functions that make use of the relatively simple io.Reader and io.Writer interfaces, your application will easily be able to "interface" with the many different types throughout both the Go standard library and the wider community.
  - For example, both os.File and net.Conn implement io.ReadWriter
- Testing with the [testify](https://github.com/stretchr/testify) package
  - This is the procedural alternative to the "table-driven" tests that are generally used in Go.
  - While I somewhat miss how quickly you can throw together new cases with table-driven tests, I must say that the procedural approach is much more **flexible.**
    - With testify, tests don't have to fit into the predefined mold of your table shape, and you never even have to worry about _altering_ the "mold," because there is no table parsing function to fiddle with.

## Important RFC’s
This is copied straight from the course, but these were cool to read through:
- [RFC 7231](https://datatracker.ietf.org/doc/html/rfc7231) – An active and widely referenced RFC.
- [RFC 9112](https://datatracker.ietf.org/doc/html/rfc9112) – Easier to read than RFC 7231, relies on understanding from RFC 9110.
- [RFC 9110](https://datatracker.ietf.org/doc/html/rfc9110) – Covers HTTP "semantics."
- [RFC 2616](https://datatracker.ietf.org/doc/html/rfc2616) – Deprecated by RFC 7231.
