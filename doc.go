/*
Package contexttest is a collection of tests for context implementations.

The standard library, alongside a context interface, provides context implementation
which gives the context implementation a behavior, such as a perceived immutability
of the context object, safe concurrent access and cancel propagation.

This package gathers testing functions for context implementations
in order to make sure that they follow the standard library behavior.

Usage:

See example in ./standardlib_test.go.
*/
package contexttest

