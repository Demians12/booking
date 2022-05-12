package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/demians12/pkg/config"
	"github.com/demians12/pkg/models"
)

//funcmap is a map of function that we can use in our templates such as format a date, return the current year...
var functions = template.FuncMap{}

var app *config.AppConfig

//create sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

//rendertemplate renders template using html/template
func RenderTemplate(w http.ResponseWriter, html string, td *models.TemplateData) {
	//get the template cache from the app config

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache(w)
	}

	t, ok := tc[html]
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Fatal("Error writing template to browser", err)
	}
}

//createtemplatecache creates a template cache as a map
func CreateTemplateCache(w http.ResponseWriter) (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	//read all the files in templates and store in pages
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		//create a template set
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		//findind out if those templates match to layout template
		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts

	}

	return myCache, err
}
