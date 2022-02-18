package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:  "convert",
	Args: cobra.ExactArgs(2),
}

func init() {
	rootCmd.AddCommand(convertCmd)
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
