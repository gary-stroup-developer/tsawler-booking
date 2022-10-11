package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

//var app *config.AppConfig

func RenderTemplate(w http.ResponseWriter, tmpl string) {
	// get the template cache from the app config
	// no need to create a template cache every single time i run this func

	tc, err := CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	// get requested template from cache
	t, ok := tc[tmpl]

	if !ok {
		log.Fatal(err)
	}
	//using buffer will allow you to know with mroe certainty where error is coming from
	buf := new(bytes.Buffer)

	err = t.Execute(buf, nil)
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
	pages, err := filepath.Glob("../../templates/*.page.tmpl")

	if err != nil {
		return myCache, err
	}

	// range through all file ending with *.pages.tmpl
	for _, page := range pages {
		// name will be base of file passed in i.e. home.page.tmpl
		name := filepath.Base(page)
		// ts = template set. the current page *template.Template is binded to ts
		ts := template.Must(template.New(name).ParseFiles(page))

		// checks to see if there are layout templates
		matches, err := filepath.Glob("../../templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		// the page may need the layout template so ParseGlob binds the files to the ts variable
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("../../templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts

	}
	return myCache, nil
}
