# LLRB [![GoDoc](https://godoc.org/github.com/mars9/llrb?status.svg)](https://godoc.org/github.com/mars9/llrb)

Package llrb implements an immutable Left-Leaning Red-Black tree as
described by Robert Sedgewick. More details relating to the
implementation are available at the following locations:

* [http://www.cs.princeton.edu/~rs/talks/LLRB/LLRB.pdf](http://www.cs.princeton.edu/~rs/talks/LLRB/LLRB.pdf)
* [http://www.cs.princeton.edu/~rs/talks/LLRB/Java/RedBlackBST.java](http://www.cs.princeton.edu/~rs/talks/LLRB/Java/RedBlackBST.java)
* [http://www.teachsolaisgames.com/articles/balanced_left_leaning.html](http://www.teachsolaisgames.com/articles/balanced_left_leaning.html)

The immutable version of the llrb tree is obviously going to be slower
than the mutable version but should offer higher read availability.
Immutability is achieved by branch copying.
