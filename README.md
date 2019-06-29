# contexttest

[![Build Status](https://travis-ci.org/posener/contexttest.svg?branch=master)](https://travis-ci.org/posener/contexttest)
[![codecov](https://codecov.io/gh/posener/contexttest/branch/master/graph/badge.svg)](https://codecov.io/gh/posener/contexttest)
[![golangci](https://golangci.com/badges/github.com/posener/contexttest.svg)](https://golangci.com/r/github.com/posener/contexttest)
[![GoDoc](https://godoc.org/github.com/posener/contexttest?status.svg)](http://godoc.org/github.com/posener/contexttest)
[![goreadme](https://goreadme.herokuapp.com/badge/posener/contexttest.svg)](https://goreadme.herokuapp.com)

Package contexttest is a collection of tests for context implementations.

The standard library, alongside a context interface, provides context implementation
which gives the context implementation a behavior, such as a perceived immutability
of the context object, safe concurrent access and cancel propagation.

This package gathers testing functions for context implementations
in order to make sure that they follow the standard library behavior.

Usage:

See example in [./standardlib_test.go](./standardlib_test.go).

## Functions

### func [TestWithCancel](https://github.com/posener/contexttest/blob/master/cancel.go#L16)

`func TestWithCancel(wc WithCancel) func(t *testing.T)`

TestWithCancel tests the behavior of a context.WithCancel function.

### func [TestWithDeadline](https://github.com/posener/contexttest/blob/master/deadline.go#L23)

`func TestWithDeadline(withDeadline WithDeadline) func(t *testing.T)`

TestWithDeadline tests the behavior of a context.WithDeadline function.

### func [TestWithTimeout](https://github.com/posener/contexttest/blob/master/timeout.go#L13)

`func TestWithTimeout(withTimeout WithTimeout) func(t *testing.T)`

TestWithTimeout tests the behavior of a context.WithTimeout function.

### func [TestWithValue](https://github.com/posener/contexttest/blob/master/value.go#L16)

`func TestWithValue(wv WithValue) func(t *testing.T)`

TestWithValue tests the behavior of a context.WithValue function.


---

Created by [goreadme](https://github.com/apps/goreadme)
