package starterobjectbucket

import (
	"embed"

	"github.com/bitwormhole/starter"
	"github.com/bitwormhole/starter-object-bucket/gen/demo"
	"github.com/bitwormhole/starter-object-bucket/gen/lib"
	"github.com/bitwormhole/starter/application"
	"github.com/bitwormhole/starter/collection"
)

const (
	theModuleName     = "github.com/bitwormhole/starter-object-bucket"
	theModuleVersion  = "v0.1.1"
	theModuleRevision = 2
)

////////////////////////////////////////////////////////////////////////////////

//go:embed src/main/resources
var theModuleSrcMainRes embed.FS

// Module 导出模块【github.com/bitwormhole/starter-object-bucket】
func Module() application.Module {

	mb := application.ModuleBuilder{}
	mb.Name(theModuleName).Version(theModuleVersion).Revision(theModuleRevision)
	mb.Resources(collection.LoadEmbedResources(&theModuleSrcMainRes, "src/main/resources"))
	mb.OnMount(lib.ExportConfigForObjectBucketLib)

	mb.Dependency(starter.Module())

	return mb.Create()
}

////////////////////////////////////////////////////////////////////////////////

//go:embed src/demo/resources
var theModuleSrcDemoRes embed.FS

// ModuleForDemo 导出模块【github.com/bitwormhole/starter-object-bucket】
func ModuleForDemo() application.Module {

	mb := application.ModuleBuilder{}
	mb.Name(theModuleName + "#demo").Version(theModuleVersion).Revision(theModuleRevision)
	mb.Resources(collection.LoadEmbedResources(&theModuleSrcDemoRes, "src/demo/resources"))
	mb.OnMount(demo.ExportConfigForObjectBucketDemo)

	mb.Dependency(Module())

	return mb.Create()
}

////////////////////////////////////////////////////////////////////////////////
