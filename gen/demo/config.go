package demo

import "github.com/bitwormhole/starter/application"

// ExportConfigForObjectBucketDemo ...
func ExportConfigForObjectBucketDemo(cb application.ConfigBuilder) error {
	return autoGenConfig(cb)
}
