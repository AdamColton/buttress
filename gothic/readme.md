## Gothic
This tool will help run Gothic projects.

A Gothic project has three basic parts; an app, a blueprint and their shared
code.

The app is the thing being developed, it's just a Go application.

The blueprint is code generator. It is also a Go application, but it runs like a
script; it's run once and it populates the generated portions of the
application.

The shared code is anything shared between both app and blueprint. Generally
this is just config values.

### gothic.json
The root of a gothic project will have a file name gothic.json. This defines the
project configuration.

```json
{
	"Name": "Project Name",
	"App": "app directory [optional] defaults to 'app'",
	"Blueprint": "blueprint directory [optional] defaults to 'blueprint'",
	"Tests": [
		"path/to/test/dir/one",
		"path/to/test/dir/two"
	]
}
```

When the blueprint project is run a file will be created name
gothic.instance.json. This contains information about the local instance of the
project. Right now it just contains the last time the blueprint ran. If no files
in the blueprint directory have changed, it will be skipped.