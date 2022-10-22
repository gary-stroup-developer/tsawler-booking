package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gary-stroup-developer/tsawler-booking/internal/config"
	"github.com/gary-stroup-developer/tsawler-booking/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

var pathToTemplates = "../templates"

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)

	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
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
		log.Println("could not get template from cahce")
	}

	//t.ExecuteTemplate(w, tmpl, nil)
	//using buffer will allow you to know with mroe certainty where error is coming from
	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)

	if err != nil {
		log.Println(err)
	}
	//render the template
	_, err = buf.WriteTo(w)

	if err != nil {
		log.Println(err)
		fmt.Println("Error writing template to browser!", err)
		return err
	}
	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files named *.page.tmpl from ../../templates
	pages, _ := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pathToTemplates))

	//range through all file ending with *.pages.tmpl
	for _, page := range pages {
		// name will be base of file passed in i.e. home.page.tmpl
		name := filepath.Base(page)
		//ts = template set. the current page *template.Template is binded to ts
		ts := template.Must(template.New(name).ParseFiles(page))

		// checks to see if there are layout templates
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		// the page may need the layout template so ParseGlob binds the files to the ts variable
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	//log.Println(myCache)
	return myCache, nil
}
