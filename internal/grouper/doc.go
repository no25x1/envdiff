// Package grouper splits an EnvFile's entries into named groups based on a
// key prefix. By default the underscore character ("_") is used as the
// delimiter so that keys such as DB_HOST and DB_PORT are placed together in
// a group named "DB".
//
// Usage:
//
//	groups := grouper.Apply(envFile, grouper.Options{
//		Delimiter:        "_",
//		IncludeUngrouped: true,
//	})
//	for _, g := range groups {
//		fmt.Println(g.Name, len(g.Entries))
//	}
package grouper
