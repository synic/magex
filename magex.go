package magex

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/sh"
	"golang.org/x/mod/modfile"
)

// ModuleVersionFromModFile pulls the version of modulePath from the go.mod file
// specified by modFilePath. This is useful in combination with MaybeInstallTool when
// the tool is paired with a module that is used in the code, for example,
// templ (https://templ.guide/) where you import the module, but it also has a cli
// codegen command and they need to be the same version. This allows your go.mod file
// to be the source of truth.
func ModuleVersionFromModFile(modFilePath, modulePath string) (string, error) {
	data, err := os.ReadFile(modFilePath)

	if err != nil {
		return "", nil
	}

	f, err := modfile.ParseLax(modFilePath, data, nil)

	if err != nil {
		return "", err
	}

	for _, info := range f.Require {
		if info.Mod.Path == modulePath {
			return info.Mod.Version, nil
		}
	}

	return "", fmt.Errorf("error parsing go.mod, module path \"%s\" not found", modulePath)
}

// ModuleVersion pulls the version of modulePath from "./go.mod" This is useful
// in combination with MaybeInstallTool when the tool is paired with a module
// that is used in the code, for example, templ (https://templ.guide/) where
// you import the module, but it also has a cli codegen command and they need
// to be the same version. This allows your go.mod file to be the source of
// truth.
func ModuleVersion(path string) (string, error) {
	return ModuleVersionFromModFile("go.mod", path)
}

// MaybeInstallTool checks to see if a go command line tool is installed, and if not,
// it installs it.
func MaybeInstallTool(name, path, version string) (string, error) {
	return maybeInstallTool(name, path, version, false)
}

// MaybeInstallToolV checks to see if a go command line tool is installed, and if not,
// it installs it. Any output is sent to stdout, similar to sh.RunV
func MaybeInstallToolV(name, path, version string) (string, error) {
	return maybeInstallTool(name, path, version, true)
}

func maybeInstallTool(name, path, version string, runV bool) (string, error) {
	var f func(cmd string, args ...string) error = sh.Run

	if runV {
		f = sh.RunV
	}

	binaryPath, err := exec.LookPath(name)

	if err == nil {
		return binaryPath, nil
	}

	err = f("go", "install", fmt.Sprintf("%s@%s", path, version))

	if err != nil {
		return "", err
	}

	binaryPath, err = exec.LookPath(name)

	if err != nil {
		return "", err
	}

	return binaryPath, nil
}
