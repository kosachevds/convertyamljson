package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var convertCmd = &cobra.Command{
	Use:  "convert",
	Args: cobra.ExactArgs(2),
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
