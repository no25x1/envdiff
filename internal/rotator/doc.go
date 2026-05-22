// Package rotator renames keys within a .env file according to a set of
// Mapping rules, preserving the original value and entry order.
//
// Typical use-case: migrating from a legacy naming convention to a new one
// without losing any configuration values.
//
//	result, err := rotator.Apply(file, rotator.Options{
//		Mappings: []rotator.Mapping{
//			{OldKey: "DB_PASS", NewKey: "DB_PASSWORD"},
//		},
//		FailOnMissing: true,
//	})
package rotator
