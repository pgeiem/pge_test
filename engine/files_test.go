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
	files, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatalf("failed to read test directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" {
			testName := fileNameWithoutExtension(file.Name())
			tariffFile := filepath.Join(testDir, testName) + ".yaml"
			testFile := filepath.Join(testDir, testName) + ".tests"

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

			for _, testCase := range testCases {
				t.Run(testName+"-"+testCase.Name, func(t *testing.T) {
					tariff, err := ParseTariffDefinition(tariffDescr)
					if err != nil {
						t.Errorf("failed to parse tariff definition: %v", err)
					}
					now, err := time.ParseInLocation("2006-01-02T15:04:05", testCase.Now, time.Local)
					if err != nil {
						t.Errorf("failed to parse now time: %v", err)
					}

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
	}
}
