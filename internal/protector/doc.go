// Package protector enforces write-protection rules on .env file keys.
//
// Keys can be protected by explicit name or by key prefix. When Apply is
// called with a base and a next EnvFile, it reports any violations where a
// protected key has been removed or its value changed.
//
// By default a violation causes an error to be returned so that callers
// (e.g. reconcile or patch pipelines) can abort early. Setting
// Options.AllowOverride to true downgrades violations to warnings while
// still returning the full list for audit purposes.
package protector
