package handlers

import (
	"fmt"
	"strings"
	"io/ioutil"
	"html/template"
)

/*
var PAGES []{string}{}

func (pghldr *Handler) LoadPages(){
	// read the whole Directory
	pages,err := ioutil.ReadDir(pghldr.TemplatesDir + "pages/")
	if err != nil {
		return fmt.Errorf("Error reading body Directory: %q",err)
	}
	// For each, read the pages and index it with the name.
	for _, page := range pages {
		name,data
	}
}
*/
// Creates a whole new base template for serving
func (pghldr *PagesHolder) LoadBase() error {
	baseHTML, err := ioutil.ReadFile(pghldr.TemplatesDir +"base.tmpl")
	if err != nil {
		return fmt.Errorf("Error loading base template: %q",err)
	}
	sidebar, err := ioutil.ReadFile(pghldr.TemplatesDir + "sidebar.tmpl")
	if err != nil {
		return fmt.Errorf("Error loading sidebar template: %q",err)
	}
	notification, err := ioutil.ReadFile(pghldr.TemplatesDir + "notifications.tmpl")
	if err != nil {
		return  fmt.Errorf("Error loading notifications template: %q",err)
	}
	shortcuts, err := ioutil.ReadFile(pghldr.TemplatesDir + "shortcuts.tmpl")
	if err != nil {
		return  fmt.Errorf("Error loading shortcuts template: %q",err)
	}

	// Replace placeholders in the base HTML with actual content
	combinedHTML := strings.ReplaceAll(string(baseHTML), "{{.SIDEBAR}}", string(sidebar))
	combinedHTML = strings.ReplaceAll(combinedHTML, "{{.NOTIFICATIONS}}", string(notification))
	combinedHTML = strings.ReplaceAll(combinedHTML, "{{.SHORTCUTS}}", string(shortcuts))
	pghldr.Base = combinedHTML
	//var tpl = new(template.Template)
	var tpl = template.New("base")
	tpl,err = tpl.Parse(string(combinedHTML))
	if err != nil {
		return fmt.Errorf("Error parsing combined html to a template: %q",err)
	}
	pghldr.Tpl = tpl
	return  nil
}

func (pghldr *PagesHolder) GetATemplate(name,templFile string) (*template.Template,error) {
	body,err := ioutil.ReadFile(pghldr.TemplatesDir + "pages/" + templFile)
	if err != nil {
		return nil,fmt.Errorf("Error getting template %s: %q",templFile,err)
	}
	sTpl := strings.ReplaceAll(pghldr.Base,"{{.BODY}}",string(body))
	var tpl = template.New(name)
	tpl,err = tpl.Parse(string(sTpl))
	if err != nil{
		return nil,fmt.Errorf("Error parsing body to template: %q",err)
	}
	return tpl,nil
}

// change this to be loaded on start up and stored to avoid loading everytime
func (pghldr *PagesHolder) GetAStaticTemplate(name,templFile string) (*template.Template,error) {
	body,err := ioutil.ReadFile(pghldr.TemplatesDir + templFile)
	if err != nil {
		return nil,fmt.Errorf("Error getting static template %s: %q",templFile,err)
	}
	var tpl = template.New(name)
	tpl,err = tpl.Parse(string(body))
	if err != nil{
		return nil,fmt.Errorf("Error parsing static body to template: %q",err)
	}
	return tpl,nil
}
