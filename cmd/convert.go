package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var convertCmd = &cobra.Command{
	Use:  "convert",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return convertFile(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
}

func convertYAMLToJSON(data []byte) ([]byte, error) {
	dataMap := make(map[string]interface{})
	err := yaml.Unmarshal(data, &dataMap)
	if err != nil {
		return nil, fmt.Errorf("YAML parsing error: %v", err)
	}
	dataMap = convertNestedMapKeysToString(dataMap)
	result, err := json.Marshal(dataMap)
	if err != nil {
		return nil, fmt.Errorf("JSON marshal error: %v", err)
	}
	return result, nil
}

func convertJSONToYAML(data []byte) ([]byte, error) {
	dataMap := make(map[string]interface{})
	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		return nil, fmt.Errorf("JSON parsing error: %v", err)
	}
	result, err := yaml.Marshal(dataMap)
	if err != nil {
		return nil, fmt.Errorf("YAML marshal error: %v", err)
	}
	return result, nil
}

func convertNestedMapKeysToString(data map[string]interface{}) map[string]interface{} {
	for key, val := range data {
		nestedMap, ok := val.(map[interface{}]interface{})
		if !ok {
			continue
		}
		stringKeyMap := make(map[string]interface{})
		for k, v := range nestedMap {
			stringKeyMap[fmt.Sprint(k)] = v
		}
		data[key] = convertNestedMapKeysToString(stringKeyMap)
	}
	return data
}

func convertFile(input, output string) error {
	inputBytes, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("file read error: %v", err)
	}

	lowerInput := strings.ToLower(input)
	var result []byte
	if strings.HasSuffix(lowerInput, ".yml") {
		result, err = convertYAMLToJSON(inputBytes)
	} else if strings.HasSuffix(lowerInput, ".json") {
		result, err = convertJSONToYAML(inputBytes)
	} else {
		return fmt.Errorf("input file must have extension '.yml' or '.json'")
	}
	if err != nil {
		return err
	}

	file, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("creating file error: %v", err)
	}
	_, err = file.Write(result)
	if err != nil {
		return fmt.Errorf("file writing error: %v", err)
	}
	return nil
}
