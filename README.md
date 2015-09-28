# PISC
Position Independent Source Code. A small, stack-based, concatenative language.

## About 

This is currently a small side project to see how far a small stack-based scripting language can go in spare time. Inspired by Factor (80%), Forth (5%) and Python (5%) and Tcl (5%) (not 100% inspired yet.)
Plans are currently for an interperated language that uses a "arrays and hashtables" approach to data stucturing.

## TODO:

If you can understand what is going on, please submit a pull request to add something that I'm missing. I'm not trying to compete with [factor](http://factorcode.org) or [forth](http://www.forth.com/forth/) for performance/features, but I trying out their style of programming to see how it goes. Ports of this interpretor to the language of your choice are welcome as well. 

With that in mind, things that I will be adding (or accepting PRs for) as [time](http://www.catb.org/jargon/html/C/copious-free-time.html) allows:

  - More tests for different combinators (if, loop, while)
  - A standard library build from a minimal core (this is a lot of the things below)
  - Stack shuffling combinators (see the ones in factor):
  -- drop ( x -- )
  -- 2drop ( x y -- )
  -- 3drop ( x y z -- )
  -- nip ( x y -- y )
  -- 2nip ( x y z -- z )
  -- dup  ( x -- x x )
  -- 2dup ( x y -- x y x y )
  -- 3dup ( x y z -- x y z )
  -- 2over ( x y z -- x y z x y )
  -- pick ( x y z -- x y z x )
  -- swap ( x y -- y x )
  - Add a way for modules to be added without a lot of major modifications in the core loop of the interp.
  - Math words. A lot is needed here, and in double/int versions ( >, <, =, -, +, div, mod ) 
  - String manipulation words (concat, >upper, >lower, int>string, double>string, etc)
  - Array and map manipulating words (ways to define them, literals, member access, so on.)
  - A basic compliation step for reducing cache misses in the core loop (transforming words into constant lookup codes, so that word lookup isn't proportional to the number of words in the core loop)
  - STDIN/STDOUT words.
  - Regex facilties. 
  - File i/o words
  - A plan for multiprocessing. (I want to pull from TCL on this one if I can)
  - Combinators for quotations, like bi and tri. 
  - A plan for a module system.
  - Syscalls
  - shellout words. 
  - struct words (when this thing allows for partial compilation or some such thing.)
  - Bindings to awesome libraries (SDL, Tk, ImageMagick)


.pisc is the file extension for these files. 
