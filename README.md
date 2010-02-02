
GoSpec
======

GoSpec is a [BDD](http://dannorth.net/introducing-bdd)-style testing framework for the [Go programming language](http://golang.org/). It allows writing self-documenting tests/specs, and executes them concurrently and safely isolated.

Source code is available at GoSpec's [project page](http://github.com/orfjackal/gospec) in GitHub.

For discussion, use the [golang-nuts mailing list](http://groups.google.com/group/golang-nuts), at least until GoSpec has so many users that it requires its own mailing list. You may also contact GoSpec's developer, [Esko Luontola](http://github.com/orfjackal), by email.


Quick Start
-----------

First you must have Go installed on your machine, as instructed in [Installing Go](http://golang.org/doc/install.html).


### Install and Update

Checkout GoSpec's source code from its Git repository and install it.

    git clone git://github.com/orfjackal/gospec.git
    cd gospec/src
    make install

When new versions of GoSpec are released, you can update it by pulling the latest version (from the `release` branch) and then installing it as above.

    git pull
    make install


### Running Specs

You can use the [gotest command](http://golang.org/cmd/gotest/) to run GoSpec's specs. The integration with gotest requires a couple of lines of boilerplate. Basically you will need to create a gotest test method, where you list all your specs and call GoSpec. See [all_specs_test.go] in the [examples] directory for an example.

    func TestAllSpecs(t *testing.T) {
        r := gospec.NewRunner()
        r.AddSpec("SomeSpec", SomeSpec)
        r.AddSpec("AnotherSpec", AnotherSpec)
        // ...
        gospec.MainGoTest(r, t)
    }

See [gotest's documentation](http://golang.org/doc/code.html#Testing) for instructions on how to use gotest.

GoSpec adds one additional parameter to gotest. Use the `-print-all` parameter to print a list of all specs. By default only the failing specs are printed. The list of all specs can be useful as documentation.


### Writing Specs

Because of using gotest, all your specs must be in files whose names end with `_test.go` and they must be added to your test suite as explained above.

Each group of specs is a method which takes `gospec.Context` as a parameter. You can call the methods on `Context` to make assertions and to declare nested specs.

For examples on how to write specs, see the files in the [examples] directory.

- [stack_test.go] is an example of how you might write a spec for the class in [stack.go].

- [fib_test.go] is an example of how you might write a spec for the class in [fib.go].

- [assert_examples_test.go] explains all the available assertions.

- [execution_model_test.go] explains GoSpec's runtime model, for example how the specs are isolated from each other and executed concurrently.


[examples]:                http://github.com/orfjackal/gospec/tree/gospec-1.x.x/examples/
[all_specs_test.go]:       http://github.com/orfjackal/gospec/blob/gospec-1.x.x/examples/all_specs_test.go
[stack.go]:                http://github.com/orfjackal/gospec/blob/gospec-1.x.x/examples/stack.go
[stack_test.go]:           http://github.com/orfjackal/gospec/blob/gospec-1.x.x/examples/stack_test.go
[fib.go]:                  http://github.com/orfjackal/gospec/blob/gospec-1.x.x/examples/fib.go
[fib_test.go]:             http://github.com/orfjackal/gospec/blob/gospec-1.x.x/examples/fib_test.go
[assert_examples_test.go]: http://github.com/orfjackal/gospec/blob/gospec-1.x.x/examples/assert_examples_test.go
[execution_model_test.go]: http://github.com/orfjackal/gospec/blob/gospec-1.x.x/examples/execution_model_test.go


Version History
---------------

**1.x.x (2010-xx-xx)**

- NOTE: In all your specs, replace `*gospec.Context` with `gospec.Context` in the spec's parameters.

- Added Fibonacci numbers example
- Added instructions about the style of naming and organizing specs
- Small changes to the print format of error messages

**1.0.0 (2009-12-30)**

- Initial release


Project Goals
-------------

With this framework you should be able to:

- Write tests/specs in BDD style. This means that you can use descriptive names and organize the specs into a nested hierarcy. Although usually two levels of nesting is enough for most situations, GoSpec supports unlimitedly nested specs, which all are based on the same abstraction (even the top-level method).

- Write specs using a fluent API, so that the code is easily readable. It should be obvious that what an assert does, and which is the *expected* and which the *actual* value.

- Write the names of the specs as strings. You are not restricted to using only those characters that are allowed in method names. All Unicode characters are allowed.

- Write the names of the specs in any style that fits the situation. The framework should not restrict you to using specific words as a prefix or suffix to the sentences (such as *describe*, *it*, *should*). You should be able to choose the best possible words yourself. GoSpec uses the syntax `c.Specify("name", ...)` for all levels in the specs, which leads the word `Specify` into becoming noise, so that it does not force you to start your sentences with any particular word.

- Write specs which are [isolated](http://agileinaflash.blogspot.com/2009/02/first.html) from the side-effects of their sibling specs. Each spec will see only the side-effects of its parent specs. In effect, the parent specs work similar to the "before" (and "after") test code in many test frameworks, and *by default* none of the specs can see its siblings (there will be a way to override the default). This will make it easier to write tests which are executed safely in isolation.

- Execute the specs concurrently on multiple CPU cores. Running the specs [quickly](http://agileinaflash.blogspot.com/2009/02/first.html) (i.e. less than 10-20 seconds) is a must for using TDD, so being able to take advantage of all processing power is important, and multiple CPU cores is the only way to go fast in the foreseen future.


License
-------

Copyright © 2009-2010 Esko Luontola <<http://www.orfjackal.net>>  
This software is released under the Apache License 2.0.  
The license text is at <http://www.apache.org/licenses/LICENSE-2.0>
