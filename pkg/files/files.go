package files

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type UntypedJson map[string]interface{}

var (
	packageJson UntypedJson
	nxJson      UntypedJson
)

func loadJsonFile(name string) UntypedJson {
	packageJson, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
	}
	defer packageJson.Close()

	byteValue, _ := io.ReadAll(packageJson)

	var result map[string]interface{}

	json.Unmarshal([]byte(byteValue), &result)

	return result
}

func GetPackageJson() map[string]interface{} {
	if packageJson == nil {
		packageJson = loadJsonFile("package.json")
	}

	return packageJson
}

func GetNxJson() UntypedJson {
	if nxJson == nil {
		nxJson = loadJsonFile("nx.json")
	}

	return nxJson
}
