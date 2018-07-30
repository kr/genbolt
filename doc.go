/*

Command genbolt generates code for conveniently
reading and writing objects in a bolt database.
It reads a set of Go type definitions
describing the layout of data in a bolt database,
and generates code for reading and writing that data.

Each struct is a bucket. Maps and slices are
also buckets. Fields with numeric types or
string or []byte are values stored in the bucket.

For example, consider this code.

	package db

	type Root struct {
		Users  []*User
		Config *Config
	}

	type User struct {
		Name string
	}

	type Config struct {
		RateLimit int64
	}

Here, Root is the root bucket.
Field Users leads to a bucket indexed by
an automatically incrementing uint64,
holding all user records.
Type User is a bucket representing a single user.
Field Config leads to the single Config bucket,
holding a single number.

It is conventional to put a +build ignore directive
in this file, so it can live in the same directory
as the generated code without its symbols conflicting.

*/
package main
