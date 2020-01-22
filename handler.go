package httprouter_wrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/yaml.v2"
)

var (
	HandlerSetting *RouterWrapperHandler
)

// wrapperハンドラーエラーメソッド
func (w *RouterWrapperHandler) Error() string {
	if w.error_ != nil {
		return w.error_.Error()
	}
	return fmt.Sprintf("%s error key %s", w.Filename, w.key)
}

// wrapperhandleコンストラクタ
func NewRouterWrapperHandler(filename string, readme ReadMe) *RouterWrapperHandler {
	if readme.Filename == "" {
		readme.Filename = "README.md"
	}
	return &RouterWrapperHandler{
		Readme:   readme,
		Filename: filename,
	}
}

// レシーバーメソッドコンストラクタ
func construct(r interface{}, methodname string) func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	return func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		var param0 = reflect.ValueOf(rw)
		var param1 = reflect.ValueOf(req)
		var param2 = reflect.ValueOf(ps)
		var value = reflect.ValueOf(r).MethodByName(methodname)
		value.Call([]reflect.Value{param0, param1, param2})
	}
}

// コンストラクタ
func New(w *RouterWrapperHandler) {
	HandlerSetting = w
	var servicename string
	var router = httprouter.New()
	if w.Readme.Refarence {
		router.GET("/refarence", refarence)
	}
	var routerf = reflect.ValueOf(router)
	var buf, _ = ioutil.ReadFile(w.Filename)
	var cfg = make(map[string]interface{})
	var err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		w.error_ = err
		return
	}

	var service map[interface{}]interface{}
	for k, v := range cfg {
		servicename = k
		service = v.(map[interface{}]interface{})
		break
	}
	port_, ok := service["port"]
	if !ok {
		w.error_ = nil
		w.key = "port"
		return
	}
	var port string
	switch reflect.ValueOf(port_).Kind() {
	case reflect.String:
		port = port_.(string)
	case reflect.Int:
		port = strconv.Itoa(port_.(int))
	case reflect.Float64:
		port = strconv.Itoa(int(port_.(float64)))
	}
	paths, ok := service["paths"]
	if !ok {
		w.error_ = nil
		w.key = "paths"
		return
	}
	for k, v := range paths.(map[interface{}]interface{}) {
		var url = k
		var maps = v.(map[interface{}]interface{})["methods"].(map[interface{}]interface{})
		for httpmethod_, methodvalue := range maps {
			httpmethod := strings.ToUpper(httpmethod_.(string))
			log.Printf("%s %s", httpmethod, url)
			var param0 = reflect.ValueOf(url)
			method := methodvalue.(map[interface{}]interface{})
			var funcname = method["func"].(string)
			var param1 = reflect.ValueOf(construct(w.Router, funcname))
			routerf.MethodByName(httpmethod).Call([]reflect.Value{param0, param1})
		}
	}
	w.port = ":" + port
	w.Handler = Log(router)
	if w.Readme.Write {
		CreateReadme(w.Readme.Filename, servicename, service)
	}
	log.Printf("%s Service Start", servicename)

	w.ListenServe = func() error {
		return http.ListenAndServe(w.port, w.Handler)
	}
}

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
			method := attribute["attribute"].(map[interface{}]interface{})
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
	fmt.Fprint(file, base)
}
