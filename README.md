# gomon
Kind of Nodemon for Go TCP based applications

# Run
`go run main.go [...options]`

# Options
* `--help` to see the usage
* `--path` path fo directory to watch (Note: path must end with trailing `/`)
* `--main` main file to run using `go run <file>` inside the `--path` (Note: Default is `app.go`)
* `--port` port which will be used by the process

# Build
* Clone this repo
* cd into repo
* `go build .`

# Example
`./gomon --path ~/projects/api/ --main main.go --port 3000`

# This currently only supports UNIX. Windows support will be added sooner
