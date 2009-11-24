
  **  Overview  **

GoSpec is a BDD framework for the Go programming language.

Project web page and source code is at http://github.com/orfjackal/gospec

For discussion, use the Go Nuts mailing list, at least until GoSpec has so many
users that it requires its own mailing list:
http://groups.google.com/group/golang-nuts


  **  Project Goals  **

With this framework you should be able to:

  - Organize the tests/specs in BDD style. Although usually two levels of
    nesting is enough for that, GoSpec supports unlimitedly nested specs.

  - Write the names of the specs as strings. You are not restricted to using
    only those characters that are allowed in method names.

  - Write the specs using a fluent API, so that the code is be easily readable.

  - Execute the specs concurrently on multiple CPU cores. Running the specs
    quickly (i.e. less than 10-20 seconds) is a must for using TDD, so being
    able to take advantage of all processing power is important.

  - Specs are isolated from the side-effects of their sibling specs. Each spec
    will see only the side-effects of its parent specs. In effect, the parent
    specs work similar to the "before" test code in many test frameworks. This
    will make it easier to write tests which are executed safely in isolation.


  **  License  **

Copyright (c) 2009 Esko Luontola <www.orfjackal.net>
This software is released under the Apache License 2.0.
The license text is at http://www.apache.org/licenses/LICENSE-2.0

