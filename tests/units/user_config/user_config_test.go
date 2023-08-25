package user_config_test

import (
	"kermoo/modules/logger"
	"kermoo/modules/user_config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadUserConfig(t *testing.T) {
	logger.MustInitLogger("fatal")

	root := "../../.."

	tt := []struct {
		name     string
		filename string
		isError  bool
		errMsg   string
		stdin    string
	}{
		// {
		// 	name:     "filename is empty",
		// 	filename: "",
		// 	isError:  true,
		// 	errMsg:   "provided filename is empty",
		// },
		{
			name:     "filename is stdin and valid json",
			filename: "-",
			isError:  false,
			stdin:    "{\"schemaVersion\":\"1\",\"process\":{\"exit\":{\"after\":\"10ms to 1s100ms\",\"code\":2}}}",
		},
		{
			name:     "filename is stdin and valid yaml",
			filename: "-",
			isError:  false,
			stdin:    "schemaVersion: \"1\"\nprocess:\n  exit:\n    after:\n      10ms to 1s100ms\n    code: 2",
		},
		{
			name:     "valid json file",
			filename: root + "/tests/units/stubs/valid.json",
			isError:  false,
		},
		{
			name:     "valid yaml file",
			filename: root + "/tests/units/stubs/valid.yaml",
			isError:  false,
		},
		{
			name:     "invalid json file",
			filename: root + "/tests/units/stubs/invalid.json",
			isError:  true,
			errMsg:   "unable to unmarshal json content",
		},
		{
			name:     "invalid yaml file",
			filename: root + "/tests/units/stubs/invalid.yaml",
			isError:  true,
			errMsg:   "invalid yaml configuration",
		},
		{
			name:     "non-existent file",
			filename: root + "/tests/units/stubs/non_existent.json",
			isError:  true,
			errMsg:   "unable to read file",
		},
		{
			name:     "valid json but does not match user_config_type",
			filename: root + "/tests/units/stubs/invalid_structure.json",
			isError:  true,
			errMsg:   "schema version is not supported",
		},
		{
			name:     "valid yaml but does not match user_config_type",
			filename: root + "/tests/units/stubs/invalid_structure.yaml",
			isError:  true,
			errMsg:   "schema version is not supported",
		},
		{
			name:     "stdin is not available",
			filename: "-",
			isError:  true,
			errMsg:   "stdin is not available to read from",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Temporarily replace os.Stdin if we're testing with filename == "-"
			if tc.filename == "-" && tc.stdin != "" {
				tmpfile, err := os.CreateTemp("", "stdin")
				if err != nil {
					t.Fatalf("Failed to create temporary file: %v", err)
				}
				defer os.Remove(tmpfile.Name())

				tmpfile.WriteString(tc.stdin)
				tmpfile.Seek(0, 0) // rewind

				oldStdin := os.Stdin
				defer func() { os.Stdin = oldStdin }() // Restore original Stdin
				os.Stdin = tmpfile
			}

			uc, err := user_config.LoadUserConfig(tc.filename)

			if tc.isError {
				require.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), tc.errMsg)
				}
			} else {
				require.NoError(t, err)

				pc, err := uc.GetPreparedConfig()

				require.NoError(t, err, "prepared config is problematic")

				err = pc.Validate()

				require.NoError(t, err, "prepared config is invalid")
			}

		})
	}
}
