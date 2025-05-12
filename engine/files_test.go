package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/goccy/go-yaml"
)

type Tests struct {
	Amount float64 `yaml:"amount"`
	End    string  `yaml:"end"`
}

type TestCase struct {
	Name  string  `yaml:"name"`
	Now   string  `yaml:"now"`
	Tests []Tests `yaml:"tests"`
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func TestTariffs(t *testing.T) {
	testDir := "./testdata"

	// Iterate over all subdirectories in the test directory
	dirs, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatalf("failed to read test directory: %v", err)
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		// Iterate over all files in the subdirectory
		path := filepath.Join(testDir, dir.Name())
		files, err := os.ReadDir(path)
		if err != nil {
			t.Fatalf("failed to read test directory: %v", err)
		}
		t.Run(dir.Name(), func(t *testing.T) {
			for _, file := range files {
				if filepath.Ext(file.Name()) != ".yaml" {
					continue
				}
				testName := fileNameWithoutExtension(file.Name())
				t.Run(testName, func(t *testing.T) {
					tariffFile := filepath.Join(path, testName) + ".yaml"
					testFile := filepath.Join(path, testName) + ".tests"

					tariffDescr, err := os.ReadFile(tariffFile)
					if err != nil {
						t.Fatalf("failed to read yaml file: %v", err)
					}

					var testCases []TestCase
					testData, err := os.ReadFile(testFile)
					if err != nil {
						t.Fatalf("failed to read test file: %v", err)
					}
					err = yaml.Unmarshal(testData, &testCases)
					if err != nil {
						t.Fatalf("failed to unmarshal yaml data: %v", err)
					}

					testSingleTariff(t, tariffDescr, testCases)
				})
			}
		})
	}
}

func testSingleTariff(t *testing.T, tariffDescr []byte, testCases []TestCase) {
	// Iterate over each test case
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tariff, err := ParseTariffDefinition(tariffDescr)
			if err != nil {
				t.Errorf("failed to parse tariff definition: %v", err)
			}
			now, err := time.ParseInLocation("2006-01-02T15:04:05", testCase.Now, time.Local)
			if err != nil {
				t.Errorf("failed to parse now time: %v", err)
			}

			// Iterate over each test in the test case
			for _, test := range testCase.Tests {
				end, err := time.ParseInLocation("2006-01-02T15:04:05", test.End, time.Local)
				if err != nil {
					t.Errorf("failed to parse end time: %v", err)
				}
				if end.Before(now) {
					t.Errorf("invalid test case, end time is before now time: %v < %v", end, now)
				}
				out := tariff.Compute(now, []AssignedRight{})
				amount := out.AmountForDuration(end.Sub(now))
				if amount.Simplify() != Amount(test.Amount) {
					t.Errorf("Amount mismatch: got %f, expected %f", amount.Simplify(), test.Amount)
				}
			}
		})
	}
}
