package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", index)
	e.GET("/pages/:page", index)
	e.Start(":8000")
}

func index(c echo.Context) error {
	outData := struct {
		Menu  []string `json:"menu"`
		Pages []string `json:"pages"`
	}{}
	outData.Menu = make([]string, 0)
	outData.Pages = make([]string, 0)

	page := c.Param("page")

	files, _ := ioutil.ReadDir("pages")
	for _, file := range files {
		if !file.IsDir() {
			outData.Menu = append(outData.Menu, file.Name())
		}
		if page != "" && page == file.Name() {
			fbody, _ := ioutil.ReadFile("pages" + "/" + file.Name())
			outData.Pages = append(outData.Pages, string(fbody))
		} else if page == "" {
			fbody, _ := ioutil.ReadFile("pages" + "/" + file.Name())
			outData.Pages = append(outData.Pages, string(fbody))
		}

	}

	rendered := renderContent(outData)

	return c.HTMLBlob(http.StatusOK, rendered)
}

func renderContent(input interface{}) []byte {
	jsonInput, _ := json.Marshal(input)

	tmplEngine := exec.Command("php", "-f", "index.php")
	tmplEngine.Stdin = strings.NewReader(string(jsonInput))
	rendered, _ := tmplEngine.Output()

	return rendered
}
