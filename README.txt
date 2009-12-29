
  **  Overview  **

GoSpec is a BDD framework for the Go programming language.

Project web page and source code is at http://github.com/orfjackal/gospec

For discussion, use the Go Nuts mailing list, at least until GoSpec has so many
users that it requires its own mailing list:
http://groups.google.com/group/golang-nuts


  **  Project Goals  **

With this framework you should be able to:

  - Write tests/specs in BDD style.  Although usually two levels of nesting is
    enough for that, GoSpec supports unlimitedly nested specs.

  - Write the names of the specs as strings.  You are not restricted to using
    only those characters that are allowed in method names.

  - Write the names of the specs in any style that fits the situation.  The
    framework should not restrict you to using specific words as a prefix or
    suffix to the sentences (such as "describe", "it", "should").  You should be
    able to choose the best possible words yourself.

  - Write specs using a fluent API, so that the code is easily readable.

  - Write specs which are isolated [1] from the side-effects of their sibling
    specs.  Each spec will see only the side-effects of its parent specs.  In
    effect, the parent specs work similar to the "before" test code in many test
    frameworks, and by default none of the specs can see its siblings.  This
    will make it easier to write tests which are executed safely in isolation.

  - Execute the specs concurrently on multiple CPU cores.  Running the specs
    quickly [1] (i.e. less than 10-20 seconds) is a must for using TDD, so being
    able to take advantage of all processing power is important, and multiple
    CPU cores is the only way to go fast in the foreseen future.

[1] http://agileinaflash.blogspot.com/2009/02/first.html


  **  License  **

Copyright © 2009 Esko Luontola <www.orfjackal.net>
This software is released under the Apache License 2.0.
The license text is at http://www.apache.org/licenses/LICENSE-2.0
