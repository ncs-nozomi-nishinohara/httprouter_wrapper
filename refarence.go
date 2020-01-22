package httprouter_wrapper

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/julienschmidt/httprouter"
)

// githubタイプのreadmeTohtml
func renderWithGitHub(md []byte) ([]byte, error) {
	client := github.NewClient(nil)
	opt := &github.MarkdownOptions{Mode: "gfm", Context: "google/go-github"}
	body, _, err := client.Markdown(context.Background(), string(md), opt)
	return []byte(body), err
}

func refarence(w http.ResponseWriter, r *http.Request, pm httprouter.Params) {
	buf, err := ioutil.ReadFile(HandlerSetting.Readme.Filename)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var html5 = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Golang Template Test</title>
</head>
<body>
%s
</body>
</html>
`
	html, err := renderWithGitHub(buf)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Contet-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	strhtml := fmt.Sprintf(html5, html)
	w.Write([]byte(strhtml))
}
