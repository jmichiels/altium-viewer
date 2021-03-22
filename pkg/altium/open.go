//go:generate rice embed-go

package altium

import (
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/zserge/lorca"
	"html/template"
	"log"
	"net"
	"net/http"
)

type templateData struct {
	Id   string
	Name string
}

func OpenProject(id string, name string) error {
	ui, err := lorca.New("", "", 800, 600)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	lis, err := net.Listen("tcp", "")
	if err != nil {
		return err
	}
	defer lis.Close()

	// Serve static assets.
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(rice.MustFindBox("assets").HTTPBox())))

	indexTemplate, err := template.New("index.html").Parse(rice.MustFindBox("templates").MustString("index.html.tmpl"))
	if err != nil {
		return err
	}

	// Serve index.html
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		indexTemplate.Execute(res, &templateData{
			Id:   id,
			Name: name,
		})
	})

	errc := make(chan error, 1)
	go func() {
		errc <- http.Serve(lis, nil)
	}()

	if err := ui.Load(fmt.Sprintf("http://%s", lis.Addr())); err != nil {
		return err
	}
	select {
	case err := <-errc:
		return err
	case <-ui.Done():
	}
	return nil
}
