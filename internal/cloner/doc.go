// Package cloner provides functionality to clone entries from one env file
// into another, with optional key renaming via prefix stripping and addition,
// key filtering, and overwrite control.
//
// Typical usage:
//
//	out, err := cloner.Apply(src, dst, cloner.Options{
//		StripPrefix: "PROD_",
//		AddPrefix:   "STAGING_",
//		Overwrite:   true,
//	})
package cloner
