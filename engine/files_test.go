package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/goccy/go-yaml"
)

type TestPoint struct {
	Amount float64 `yaml:"amount"`
	End    string  `yaml:"end"`
}

type TestCase struct {
	Name           string      `yaml:"name"`
	Now            string      `yaml:"now"`
	History        string      `yaml:"history"`
	TestPoints     []TestPoint `yaml:"tests"`
	ExpectedExpiry string      `yaml:"expiry"`
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

					testSingleTariff(t, path, tariffDescr, testCases)
				})
			}
		})
	}
}

func testSingleTariff(t *testing.T, path string, tariffDescr []byte, testCases []TestCase) {
	// Iterate over each test case

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {

			tariff, err := ParseTariffDefinition(tariffDescr)
			if err != nil {
				t.Fatalf("failed to parse tariff definition: %v", err)
			}
			now, err := time.ParseInLocation("2006-01-02T15:04:05", testCase.Now, time.Local)
			if err != nil {
				t.Fatalf("failed to parse now time: %v", err)
			}

			// Load parking right history if specified
			var history AssignedRights
			if testCase.History != "" {
				filename := filepath.Join(path, testCase.History)
				history, err = LoadHistoryFromFile(filename)
				if err != nil {
					t.Fatalf("failed to load history from file: %v", err)
				}
			}

			// Compute the tariff table
			table := tariff.Compute(now, history)

			//Display JSON output
			json, err := table.ToJson()
			if err != nil {
				t.Fatalf("failed to convert table to JSON: %v", err)
			}
			fmt.Println(string(json))

			// Check the expiry date
			expectedExpiry, err := time.ParseInLocation("2006-01-02T15:04:05", testCase.ExpectedExpiry, time.Local)
			if err == nil && !expectedExpiry.IsZero() {
				if table.ExpiryDate.Before(now) {
					t.Errorf("Invalid expiry date: %v is before now time: %v", table.ExpiryDate, now)
				}
				if table.ExpiryDate != expectedExpiry {
					t.Errorf("Expiry date mismatch: got %v, expected %v", table.ExpiryDate, expectedExpiry)
				}
			}

			// Iterate over each testpoints in the current test case
			for _, test := range testCase.TestPoints {
				end, err := time.ParseInLocation("2006-01-02T15:04:05", test.End, time.Local)
				if err != nil {
					t.Fatalf("failed to parse end time: %v", err)
				}
				if end.Before(now) {
					t.Fatalf("invalid test case, end time is before now time: %v < %v", end, now)
				}
				amount := table.AmountForDuration(end.Sub(now))
				if amount.Simplify() != Amount(test.Amount) {
					t.Errorf("Amount mismatch: got %f, expected %f", amount.Simplify(), test.Amount)
				}
			}
		})
	}
}

func LoadHistoryFromFile(filename string) (AssignedRights, error) {
	historyData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return LoadAssignedRightHistoryFromJSON(historyData)
}
