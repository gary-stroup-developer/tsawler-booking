package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gary-stroup-developer/tsawler-booking/pkg/config"
	"github.com/gary-stroup-developer/tsawler-booking/pkg/models"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	// no need to create a template cache every single time i run this func
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		// get the template cache from the app config
		tc, _ = CreateTemplateCache()
	}
	// get requested template from cache
	t, ok := tc[tmpl]

	if !ok {
		log.Fatal("error getting template from cache")
	}
	//using buffer will allow you to know with mroe certainty where error is coming from
	buf := new(bytes.Buffer)

	err := t.Execute(buf, nil)
	if err != nil {
		log.Println(err)
	}
	//render the template
	_, err = buf.WriteTo(w)

	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files named *.page.tmpl from ../../templates
	pages, err := filepath.Glob("../templates/*.gohtml")

	if err != nil {
		return myCache, err
	}

	// range through all file ending with *.pages.tmpl
	for _, page := range pages {
		// name will be base of file passed in i.e. home.page.tmpl
		name := filepath.Base(page)
		// ts = template set. the current page *template.Template is binded to ts
		ts := template.Must(template.New(name).ParseFiles(page))

		// // checks to see if there are layout templates
		// matches, err := filepath.Glob("../templates/*.layout.gohtml")
		// if err != nil {
		// 	return myCache, err
		// }

		// // the page may need the layout template so ParseGlob binds the files to the ts variable
		// if len(matches) > 0 {
		// 	ts, err = ts.ParseGlob("../templates/*.layout.gohtml")
		// 	if err != nil {
		// 		return myCache, err
		// 	}
		// }

		myCache[name] = ts

	}
	return myCache, nil
}
