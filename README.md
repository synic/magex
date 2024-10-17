# MAGEX

Extension functions and utilities to augment the Go mage build system:
https://magefile.org/

## Functions

#### `ModuleVersion`

Pulls the version of a given module from the `go.mod` file. This is intended to
be used in pair with `MaybeInstallTool` so that the `go.mod` file can be the
source of truth for a command line utility. For example, templ
(https://templ.guide) is a module that you import in your code, AND a cli
codegen utility, and then need to be the same version. Example:

```go
func Codegen() error {
  version, err := magex.ModuleVersion("github.com/a-h/templ")

  if err != nil {
    return err
  }

  path, err := magex.MaybeInstallTool("templ", "github.com/a-h/templ/cmd/templ", version)

  if err != nil {
    return err
  }

  return sh.Run(path, "generate")
}
```

#### `MaybeInstallTool`

Checks to see if a go command line tool is installed, and if it's not, it
installs it.

```go
func Dev() error {
  path, err := magex.MaybeInstallTool("air", "github.com/cosmtrek/air", "v1.49.0")

  if err != nil {
    return err
  }

  return sh.RunV(path)
}
```

#### `MaybeInstallToolToDestination`

Checks to see if a go command line tool exists at the given location. If not,
it installs it. Useful if you want to install to `./bin` and avoid changing the
host `GOPATH`.

```go
func Dev() error {
  path, err := magex.MaybeInstallToolToDestination(
    "air", "github.com/cosmtrek/air", "v1.49.0", "bin",
  )

  if err != nil {
    return err
  }

  return sh.RunV(path)
}
```
