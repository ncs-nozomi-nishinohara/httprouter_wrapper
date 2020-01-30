package httprouter_wrapper

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/ncs-nozomi-nishinohara/httprouter_wrapper/v2/wrapper_utils"
	"gopkg.in/yaml.v2"
)

var (
	HandlerSetting     *wrapper_utils.RouterWrapperHandler
	methodName_to_func map[string]reflect.Value
)

func init() {
	methodName_to_func = make(map[string]reflect.Value)
}

// wrapperhandleコンストラクタ
func NewRouterWrapperHandler(filename string, readme wrapper_utils.ReadMe) *wrapper_utils.RouterWrapperHandler {
	if readme.Filename == "" {
		readme.Filename = "README.md"
	}
	return &wrapper_utils.RouterWrapperHandler{
		Readme:   readme,
		Filename: filename,
	}
}

// レシーバーメソッドコンストラクタ
func construct(r interface{}, methodname string) func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	return func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		var value reflect.Value
		var param0 = reflect.ValueOf(rw)
		var param1 = reflect.ValueOf(req)
		var param2 = reflect.ValueOf(ps)
		value, ok := methodName_to_func[methodname]
		if !ok {
			value = reflect.ValueOf(r).MethodByName(methodname)
			methodName_to_func[methodname] = value
		}
		value.Call([]reflect.Value{param0, param1, param2})
	}
}

// コンストラクタ
func New(w *wrapper_utils.RouterWrapperHandler) {
	HandlerSetting = w
	var servicename string
	var router = httprouter.New()
	var routerf = reflect.ValueOf(router)
	var buf, _ = ioutil.ReadFile(w.Filename)
	var cfg = make(map[string]interface{})
	var err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		w.SetError(err)
		return
	}

	var service map[interface{}]interface{}
	for k, v := range cfg {
		servicename = k
		service = v.(map[interface{}]interface{})
		break
	}
	environment, ok := service["environment"]
	if ok {
		for key, value := range environment.(map[interface{}]interface{}) {
			if os.Getenv(key.(string)) == "" {
				os.Setenv(key.(string), value.(string))
			}
		}
	}
	migrationflg := false
	migrationenv := os.Getenv("MIGRATION")
	if migrationenv != "" && strings.ToUpper(migrationenv) == "Y" {
		migrationflg = true
	}
	w.Migration = migrationflg
	migration, ok := service["migration"]
	if ok {
		var driver, dirname string
		driver = migration.(map[interface{}]interface{})["driver"].(string)
		dirname = migration.(map[interface{}]interface{})["dirname"].(string)
		wrapper_utils.Migration(driver, dirname)
	}
	if w.Readme.Refarence {
		router.GET("/refarence", refarence)
	}
	if w.Readme.Write {
		wrapper_utils.CreateReadme(w.Readme.Filename, servicename, service)
	}

	port_, ok := service["port"]
	if !ok {
		w.SetError(nil)
		w.SetKey("port")
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
	w.SetPort(port)

	paths, ok := service["paths"]
	if !ok {
		w.SetError(nil)
		w.SetKey("paths")
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

	w.Handler = wrapper_utils.Log(router)
	log.Printf("%s Service Start", servicename)

	w.ListenServe = func() error {
		return http.ListenAndServe(w.GetPort(), w.Handler)
	}
}
