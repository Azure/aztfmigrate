package tf

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
)

const terraformBinary = "terraform"

// FindTerraform finds the path to the terraform executable whose version meets the min/max version constraint.
// It first tries to find from the local OS PATH. If there is no match, it will then download the release of the minVersion from hashicorp to the tfDir.
func FindTerraform(ctx context.Context, minVersion, maxVersion *version.Version) (string, error) {

	// Initialize the workspace
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("error finding the user cache directory: %w", err)
	}
	rootDir := filepath.Join(cacheDir, "azapi2azurerm")
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return "", fmt.Errorf("creating workspace root %q: %w", rootDir, err)
	}
	tfDir := filepath.Join(rootDir, "terraform")
	if err := os.MkdirAll(tfDir, 0755); err != nil {
		return "", fmt.Errorf("creating terraform cache dir %q: %w", tfDir, err)
	}

	var terraformPath string
	opts := []tfinstall.ExecPathFinder{
		tfinstall.ExactPath(filepath.Join(tfDir, terraformBinary)),
		tfinstall.LookPath(),
		tfinstall.ExactVersion(maxVersion.String(), tfDir),
	}

	// go through the options in order
	// until a valid terraform executable is found
	for _, opt := range opts {
		p, err := opt.ExecPath(ctx)
		if err != nil {
			return "", fmt.Errorf("unexpected error: %w", err)
		}

		if p == "" {
			// strategy did not locate an executable - fall through to next
			continue
		}

		v, err := getTerraformVersion(ctx, p)
		if err != nil {
			return "", fmt.Errorf("error getting terraform version for executable found at path %s: %w", p, err)
		}

		if v.LessThan(minVersion) || v.GreaterThan(maxVersion) {
			continue
		}

		terraformPath = p
		break
	}

	if terraformPath == "" {
		return "", fmt.Errorf("could not find terraform executable")
	}

	return terraformPath, nil
}

func getTerraformVersion(ctx context.Context, execPath string) (*version.Version, error) {
	tf, err := tfexec.NewTerraform(os.TempDir(), execPath)
	if err != nil {
		return nil, fmt.Errorf("error running NewTerraform with path %q: %w", execPath, err)
	}
	ver, _, err := tf.Version(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("error running terraform version: %w", err)
	}
	return ver, nil
}
