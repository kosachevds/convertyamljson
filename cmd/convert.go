package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var convertCmd = &cobra.Command{
	Use:   "convert input_path output_path",
	Short: "Сonvert YAML to JSON and vice versa",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cmd.Flags().GetBool("dir")
		if err != nil {
			return err
		}
		if dir {
			return convertDir(args[0], args[1])
		}
		return convertFile(args[0], args[1])
	},
}

func init() {
	convertCmd.Flags().BoolP("dir", "d", false, "Recursively convert files in dir")
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

	var result []byte
	ext := strings.ToLower(filepath.Ext(input))
	switch ext {
	case ".yml":
		result, err = convertYAMLToJSON(inputBytes)
	case ".json":
		result, err = convertJSONToYAML(inputBytes)
	default:
		return fmt.Errorf("input file must have extension '.yml' or '.json'")
	}
	if err != nil {
		return err
	}

	file, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("creating file error: %v", err)
	}
	defer file.Close()
	_, err = file.Write(result)
	if err != nil {
		return fmt.Errorf("file writing error: %v", err)
	}
	return nil
}

func convertDir(input, output string) error {
	return filepath.WalkDir(input, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		resultPath := filepath.Join(output, path[len(input):])
		ext := strings.ToLower(filepath.Ext(path))
		lenWithoutExt := len(resultPath) - len(ext)
		switch ext {
		case ".yml":
			resultPath = resultPath[:lenWithoutExt] + ".json"
		case ".json":
			resultPath = resultPath[:lenWithoutExt] + ".yml"
		default:
			fmt.Printf("Cannot convert %s file\n", path)
			return nil
		}

		resultDir := filepath.Dir(resultPath)
		err = os.MkdirAll(resultDir, 0777)
		if err != nil {
			return fmt.Errorf("cannot create %s: %v", resultDir, err)
		}
		err = convertFile(path, resultPath)
		if err != nil {
			return fmt.Errorf("cannot convert %s: %v", path, err)
		}
		return nil
	})
}
