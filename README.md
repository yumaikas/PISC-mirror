# PISC
Position Independent Source Code. A small, stack-based, concatenative language.

## About 

This is currently a small side project to see how far a small stack-based scripting language can go in spare time. Inspired by Factor (80%), Forth (5%) and Python (5%) and Tcl (5%) (not 100% inspired yet.)
Plans are currently for an interperated language that uses a "arrays and hashtables" approach to data stucturing.

## TODO:

If you can understand what is going on, please submit a pull request to add something that I'm missing. I'm not trying to compete with [factor](http://factorcode.org) or [forth](http://www.forth.com/forth/) for performance/features, but I trying out their style of programming to see how it goes. Ports of this interpretor to the language of your choice are welcome as well. 

With that in mind, things that I will be adding (or accepting PRs for) as [time](http://www.catb.org/jargon/html/C/copious-free-time.html) allows:

  - More tests for different combinators (if, loop, while)
  - Stack shuffling combinators (see the ones in factor):
  -- drop
  -- 2drop
  -- 3drop
  -- nip
  -- 2nip
  -- dup
  -- 2dup
  -- 3dup
  -- 2over
  -- pick
  -- swap
  - Array and map manipulating words
  - A standard library build from a minimal core
  - Math words. A lot is needed here, and in double/int versions. 
  - 

.pisc is the file extension for these files. I don't have a standard library in place yet.