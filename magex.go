package magex

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/magefile/mage/sh"
	"golang.org/x/mod/modfile"
)

// ModuleVersion pulls the version of `name` from "./go.mod" This is useful
// in combination with MaybeInstallTool when the tool is paired with a module
// that is used in the code, for example, templ (https://templ.guide/) where
// you import the module, but it also has a cli codegen command and they need
// to be the same version. This allows your go.mod file to be the source of
// truth.
func ModuleVersion(name string) (string, error) {
	data, err := os.ReadFile("go.mod")

	if err != nil {
		return "", nil
	}

	f, err := modfile.ParseLax("go.mod", data, nil)

	if err != nil {
		return "", err
	}

	for _, info := range f.Require {
		if info.Mod.Path == name {
			return info.Mod.Version, nil
		}
	}

	return "", fmt.Errorf("error parsing go.mod, module path \"%s\" not found", name)
}

// MaybeInstallTool checks to see if a go command line tool is installed, and if not,
// it installs it.
func MaybeInstallTool(cmdName, modulePath, version string) (string, error) {
	binaryPath, err := exec.LookPath(cmdName)

	if err == nil {
		return binaryPath, nil
	}

	fmt.Printf("Installing %s@%s...", cmdName, version)
	err = sh.Run("go", "install", fmt.Sprintf("%s@%s", modulePath, version))
	fmt.Println("done")

	if err != nil {
		return "", err
	}

	binaryPath, err = exec.LookPath(cmdName)

	if err != nil {
		return "", err
	}

	return binaryPath, nil
}

// MaybeInstallToolToDest checks to see if the tool command exists at the destination,
// and if not, attempts to install it.
func MaybeInstallToolToDestination(
	cmdName, modulePath, version, destination string,
) (string, error) {
	destination, err := filepath.Abs(destination)
	cmd := path.Join(destination, cmdName)

	if err != nil {
		return "", fmt.Errorf("error converting %s to absolute path: %w", destination, err)
	}

	_, err = os.Stat(cmd)

	if err == nil {
		return cmd, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	fmt.Printf("Installing %s@%s to %s\n", cmdName, version, destination)
	err = sh.RunWith(
		map[string]string{"GOBIN": destination},
		"go",
		"install",
		fmt.Sprintf("%s@%s", modulePath, version),
	)

	if err != nil {
		return "", err
	}

	_, err = os.Stat(path.Join(destination, cmdName))

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	return cmd, nil
}
