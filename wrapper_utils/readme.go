package wrapper_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// readme生成関数
func CreateReadme(filename string, servicename string, m map[interface{}]interface{}) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// タイトル
	base := `# %s

`
	base = fmt.Sprintf(base, servicename)
	// 全体
	base += `## Describe

%s

`

	base = fmt.Sprintf(base, m["describe"].(string))
	paths := m["paths"]
	type Attribute struct {
		Method    string
		Describe  string
		Parametes string
		Json      bool
	}
	for k, v := range paths.(map[interface{}]interface{}) {
		var url = k
		var methods = v.(map[interface{}]interface{})["methods"].(map[interface{}]interface{})
		sets := []Attribute{}
		for httpmethod_, methodvalue := range methods {
			httpmethod := strings.ToUpper(httpmethod_.(string))
			httpmethod = strings.TrimSpace(httpmethod)
			attribute := methodvalue.(map[interface{}]interface{})
			method, ok := attribute["attribute"].(map[interface{}]interface{})
			if ok {
				set := Attribute{}
				set.Method = httpmethod
				describe, ok := method["describe"]
				if ok {
					set.Describe = describe.(string)
				}
				parameter, ok := method["parameter"]
				if ok {
					var buf bytes.Buffer
					err := json.Indent(&buf, []byte(parameter.(string)), "", "  ")
					if err == nil {
						set.Json = true
						set.Parametes = buf.String()
					} else {
						set.Parametes = parameter.(string)
					}
				}
				sets = append(sets, set)
			}
		}
		if len(sets) > 0 {
			base += fmt.Sprintf("<details><summary>%s</summary>\n\n", url)
			for _, set := range sets {
				base += fmt.Sprintf("<details><summary>%s</summary>\n\n", set.Method)
				base += "- descirbe\n\n"
				base += set.Describe + "\n\n"
				if set.Parametes != "" {
					base += "- Parameter\n\n"
					if set.Json {
						base += "```json:paramete.json\n%s\n```\n\n"
						base = fmt.Sprintf(base, set.Parametes)
					} else {
						base += set.Parametes + "\n\n"
					}
				}
				base += "</details>\n\n"
			}
			base += "</details>"
		}
	}
	if base != "" {
		fmt.Fprint(file, base)
	}
}
