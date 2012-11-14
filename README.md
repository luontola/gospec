
GoSpec
======

GoSpec is a [BDD](http://dannorth.net/introducing-bdd)-style testing framework for the [Go programming language](http://golang.org/). It allows writing self-documenting tests/specs, and executes them in parallel and safely isolated.

Source code is available at <http://github.com/orfjackal/gospec>

For discussion, use the [golang-nuts mailing list](http://groups.google.com/group/golang-nuts), at least until GoSpec has so many users that it requires its own mailing list. You may also contact GoSpec's developer, [Esko Luontola](http://github.com/orfjackal), by email.


Quick Start
-----------

First you must have Go installed on your machine, as instructed in [Installing Go](http://golang.org/doc/install.html).


### Install and Update

Download GoSpec using the `go get` tool:

    go get "github.com/orfjackal/gospec/src/gospec"

See `go help get` for more instructions on using the tool.

See "Version History" for any additional upgrade notes.


### Sample Project

Make a copy of the [hello-world-template] directory to get started. You can run its tests with the `go test` command. All test files must end with `_test.go` and all specs must be listed in `all_specs_test.go`.


### Running Specs

You can use the `go test` command to run GoSpec's specs. The integration with gotest requires a couple of lines of boilerplate: you'll need to write a gotest test method, where you list all your specs and call GoSpec. See [all_specs_test.go] in the [examples] directory for an example. Also all your specs must be in files whose names end with `_test.go`.

See [gotest's documentation](http://golang.org/doc/code.html#Testing) for instructions on how to use gotest.

GoSpec adds one additional parameter to gotest. Use the `-print-all` parameter to print a list of all specs: `go test -print-all` Otherwise only the failing specs are printed. The list of all specs can be useful as documentation.


### Writing Specs

The following imports are needed. The first imports the `gospec.Context` interface and the second is needed for using GoSpec's expectation matchers (`Equals`, `IsTrue`, `IsNil`, `Not()`, `Contains` etc.) without having to prefix them with the package name. (In a future GoSpec version the matchers will be moved to their own package.)

    import "github.com/orfjackal/gospec/src/gospec"
    import . "github.com/orfjackal/gospec/src/gospec"

The specs are written as functions which take `gospec.Context` as a parameter. You can call the methods of `Context` to declare expectations and nested specs.

For examples on how to write specs, see the files in the [examples] directory.

- [fib_test.go] is an example of the recommended style for writing specs for the class in [fib.go].

- [stack_test.go] is an example of the recommended style for writing specs for the class in [stack.go].

- [expectation_syntax_test.go] explains the syntax for writing expectations.

- [execution_model_test.go] explains GoSpec's runtime model, for example how the specs are isolated from each other and executed in parallel.

<!---
Replace /master/ with /gospec-1.x.x/ before making a release.
-->
[hello-world-template]:       http://github.com/orfjackal/gospec/tree/master/src/hello-world-template/
[examples]:                   http://github.com/orfjackal/gospec/tree/master/src/examples/
[all_specs_test.go]:          http://github.com/orfjackal/gospec/blob/master/src/examples/all_specs_test.go
[fib.go]:                     http://github.com/orfjackal/gospec/blob/master/src/examples/fib.go
[fib_test.go]:                http://github.com/orfjackal/gospec/blob/master/src/examples/fib_test.go
[stack.go]:                   http://github.com/orfjackal/gospec/blob/master/src/examples/stack.go
[stack_test.go]:              http://github.com/orfjackal/gospec/blob/master/src/examples/stack_test.go
[expectation_syntax_test.go]: http://github.com/orfjackal/gospec/blob/master/src/examples/expectation_syntax_test.go
[execution_model_test.go]:    http://github.com/orfjackal/gospec/blob/master/src/examples/execution_model_test.go


Version History
---------------

**1.x.x (2012-xx-xx)**

- ...

**1.3.9 (2012-03-28)**

*UPGRADE NOTES:* Check your imports - when using the `go` tool they are different than when using the old hand-written Makefiles.

- Build using the `go` tool instead of Makefiles
- Upgraded to Go 1 (weekly.2012-02-07)

**1.3.8 (2011-08-04)**

- Upgraded to Go release.r59 (weekly.2011-07-07)

**1.3.7 (2011-07-02)**

- Upgraded to Go release.r58 (weekly.2011-06-23)

**1.3.6 (2011-05-04)**

- Upgraded to Go release.r57.1 (weekly.2011-04-27)

**1.3.5 (2011-01-21)**

- Upgraded to Go release.2011-01-20

**1.3.4 (2010-10-15)**

- Upgraded to Go release.2010-10-13

**1.3.3 (2010-10-11)**

- Fixed an occasional off-by-one in exception stack trace line numbers

**1.3.2 (2010-10-01)**

- Upgraded to Go release.2010-09-29

**1.3.1 (2010-09-11)**

- [Issue 754](http://code.google.com/p/go/issues/detail?id=754) was fixed in Go release.2010-09-06, so line numbers in GoSpec's stack traces are now correct
- Fixed an occasional off-by-one in exception stack trace line numbers
- Upgraded to Go release.2010-09-06

**1.3.0 (2010-09-06)**

*UPGRADE NOTES:* If you have written custom matchers, their result parameters' types have changed. Also the error messages are expected to be in a slightly different format. See [expectation_syntax_test.go] or GoSpec's built-in matchers for examples.

- New error message format
- Workaround for a bug in gedit 2.28.0 which caused stack traces to be sometimes non-clickable
- Improved the stack traces to hide GoSpec internals also for root specs
- Upgraded to Go release.2010-08-25

**1.2.0 (2010-04-29)**

*UPGRADE NOTES:* In your spec suite, replace `r.AddSpec("SomeSpec", SomeSpec)` with `r.AddSpec(SomeSpec)`.

- Recover from panics in specs and report their stack traces
- Retrieve the spec function names using reflection, to avoid some boilerplate in the spec suite
- Changes to error messages: function names and full file paths are now shown in the stack traces
- Improved documentation and provided a hello world project template
- Removed the deprecated `c.Then()` syntax
- Upgraded to Go release.2010-04-13

**1.1.0 (2010-03-08)**

*UPGRADE NOTES:* In all your specs, replace `*gospec.Context` with `gospec.Context` in the spec's parameters. Add `import . "gospec"` to the imports and change every assertion of the old `c.Then(x).Should.Equal(y)` syntax to use the new `c.Expect(x, Equals, y)` syntax.

- New expectation syntax. The old `c.Then()` syntax is deprecated and will be removed later.
- New matchers: IsSame, IsNil, IsTrue, IsFalse, ContainsAll, ContainsAny, ContainsExactly, ContainsInOrder, ContainsInPartialOrder
- Added Fibonacci numbers example
- Added instructions about the style of naming and organizing specs
- Minor changes to the print format of error messages
- Upgraded to Go release.2010-02-04

**1.0.0 (2009-12-30)**

- Initial release


Project Goals
-------------

The following are *a must*, because they enable using [specification-style](http://blog.orfjackal.net/2010/02/three-styles-of-naming-tests.html) the way I prefer:

- **Unlimited Nesting** - The specs can be organized into a nested hierarchy. This makes it possible to apply [One Assertion Per Test](http://www.artima.com/weblogs/viewpost.jsp?thread=35578) which [isolates the reason for a failure](http://agileinaflash.blogspot.com/2009/02/first.html), because the specs are very fine-grained. Many unit testing tools allow only 2 levels of nesting (e.g. JUnit) and a few only 1 level (e.g. gotest), but for specification-style at least 3 levels are needed (e.g. JDave), and once you have 3 levels you might as well implement unlimited levels with the same abstraction.

- **Isolated Execution** - The specs must be [isolated from the side-effects](http://agileinaflash.blogspot.com/2009/02/first.html) of their sibling specs. Each spec will see only the side-effects of its parent specs. In effect, the parent specs work similar to the "before" (and "after") test code in many test frameworks, and *by default* none of the specs can see its siblings (there will be a way to override the default). Without this isolation, it would be harder to write reliable side-effect free specs, which in turn would force the specs to be organized differently than what was desired.

- **No Forced Words** - [Getting the words right](http://behaviour-driven.org/GettingTheWordsRight) was the starting point for BDD, so it is absurd that almost all of the BDD frameworks force the programmer to use fixed words (describe, it, should, given, when, then etc.) which incline the programmer to write spec names as sentences which begin or end with those words. You should be able to choose yourself the best possible words that fit a situation. GoSpec uses the syntax `c.Specify("name", ...)` for all levels in the specs, which leads to the word `Specify` becoming *background noise*, so that you ignore it and it does not force you to start your sentences with any particular word (using a meaningless word such as "Spec" would also be a good choice, as long as it is easy to pronounce when communicating with others).

The following are *nice-to-haves*, which make it more pleasant to use the framework:

- **Plain Text Names** - You can use any Unicode characters in the spec names, because they are declared as strings. Using only those characters that are allowed in method names would be too limiting and hard to read.

- **Fluent API** - The syntax for writing specs should be easily readable. It should be obvious that what an assert does, and which is the *expected* and which the *actual* value. Also writing the specs should be easy, requiring as little syntax as possible, but readability has always higher priority than writability.

- **Parallel Execution** - Running the specs [quickly](http://agileinaflash.blogspot.com/2009/02/first.html) (i.e. less than 10-20 seconds) is a must for using TDD, so being able to take advantage of all processing power is important, and multiple CPU cores is the only way to go fast in the foreseen future. GoSpec executes the specs using as much parallelism as possible (one goroutine for each leaf spec), so that it would be possible to utilize all available CPU cores (just remember to set GOMAXPROCS).


License
-------

Copyright © 2009-2012 Esko Luontola <<http://www.orfjackal.net>>  
This software is released under the Apache License 2.0.  
The license text is at <http://www.apache.org/licenses/LICENSE-2.0>
