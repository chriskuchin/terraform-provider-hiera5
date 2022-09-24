// Package catch provides the controlled recovery of panics using a special error implementation that has a cause
// (another error). The package is meant to be used in situation where it is tedious or impossible to propagate errors
// as return values.
package catch
