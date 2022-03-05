package lib

import application "github.com/bitwormhole/starter/application"

// ExportConfigForObjectBucketLib  ...
func ExportConfigForObjectBucketLib(cb application.ConfigBuilder) error {
	return autoGenConfig(cb)
}
