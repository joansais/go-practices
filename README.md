# Go Practices

My programming exercices for learning Go and exploring some of its underlying design principles.

## Package 'wiki'
This package contains an implementation of the [Writing Web Applications](https://golang.org/doc/articles/wiki/) tutorial, a simple wiki. This implementation supports markdown syntax, page linking, page renaming, and extensibility of storage strategies and markdown syntaxes (via the `PageStore` and `SyntaxHandler` interfaces).

## Package 'wikix'
This is an alternative implementation of the `wiki` package using exceptions instead of error codes. Although using panic-recover in this way goes against Go's error handling conventions and philosophy (\*), I think this is a useful exercise for comparing both error handling strategies (see the _controversy_ section below).

(\*) To be exact, the use of panic-recover for error handling *internally within* a package is an [accepted practice](https://github.com/golang/go/wiki/PanicAndRecover#Usage_in_a_Package), actually used in some places of the standard library's implementation. This `wikix` package, however, exposes interfaces that throw exceptions instead of returning errors. Potential implementations of these interfaces in other packages (let's say, a `wikisql.DbPageStore`) would violate the convention that explicit panics must not cross package boundaries. Therefore the interface itself (being an exposed interface) is violating the convention.

Another difference with the `wiki` package is the use of the [go-check](http://labix.org/gocheck) testing framework instead of the standard _testing_ library. In this case, the benefits are pretty obvious IMO.

## Package 'exception'
An exception-like alternative idiom for panic-recover, used by the `wikix` package.

# Exceptions vs. error values controversy

Here's my summary of the exceptions vs. error values (or error codes) controversy, and how it applies to Go:

* The main criticism of exceptions is that they complicate reasoning about the code behavior, because of the hidden control-flow paths they introduce (the exceptional paths). Reasoning about the exceptional paths is especially important when "undo" actions are needed (e.g. releasing acquired resources, or leaving the object in a consistent state). While this reasoning is necessary in both error handling strategies, some argue that error values naturally "force" the programmer to do such reasoning (by making the exceptional paths explicit), while exceptions require more discipline to take into account these hidden paths. In other words, detractors of exceptions claim that they provide a false sense of correctness and thus are more prone to result in improper error handling than error values.

* Correctly handling the "undo" actions can be hard and error-prone (in both strategies) without the help of an appropriate technique or syntax, like [RAII](https://en.wikipedia.org/wiki/Resource_Acquisition_Is_Initialization) in C++. Note that Java's *try-finally* clause is considered [insufficient](http://www.cs.virginia.edu/~weimer/p/weimer-toplas2008.pdf) for this purpose. In Go, such mechanism is provided by [deferred functions](http://blog.golang.org/defer-panic-and-recover), which also work well with panics. In this sense, exception handling in Go would probably have been less problematic than in Java with *try-finally*.

* The main advantage of exceptions is that they produce code that is more expressive and understandable for the normal (non-exceptional) control flow. In many cases, errors simply need to be propagated to the caller, with no "undo" actions or error conversions required. Exception propagation does not add extra code to the normal control flow, while error value propagation requires boilerplate code, impacting thus readability. (However, proponents of error values would claim that this is not boilerplate code, but code that explicitly expresses that the programmer's intention was to propagate the error as is.)

* Generally speaking, handling errors properly is hard, no matter the strategy used. Beyond that, there are good arguments in favor of exceptions, and good arguments in favor of error values. Also, the specific features and design of each programming language must be taken into account when comparing both strategies. In the end, it is almost a matter of taste of the language designers. Go [advocates](https://golang.org/doc/faq#exceptions) the use of error values. Still, exceptions are supported (with panic-recover), but by convention they should be used in a very limited way. Strong supporters of exceptions could of course ignore this convention in application-level code (like in the `wikix` example above), even in large code bases and across packages, but this doesn't seem a good idea, since 1) the language design does not facilitate the exception idiom, and 2) the rest of the community will be coding all frameworks and libraries you will need to reuse using error values.

Some references:
* [Exception Handling: A False Sense of Security](http://ptgmedia.pearsoncmg.com/images/020163371x/supplements/Exception_Handling_Article.html)
* [Lessons Learned from Specifying Exception-Safety for the C++ Standard Library](http://www.boost.org/community/exception_safety.html)
* [Exceptional Situations and Program Reliability (PDF)](http://www.cs.virginia.edu/~weimer/p/weimer-toplas2008.pdf)
* [Cleaner, more elegant, and wrong](http://blogs.msdn.com/b/oldnewthing/archive/2004/04/22/118161.aspx)

