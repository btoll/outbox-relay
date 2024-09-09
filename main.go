package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	//go:embed tpl/*
	templateFiles embed.FS
	serviceName   *string
	image         *string
	dbName        *string
)

type Image struct {
	Name string
	Tag  string
}

type Service struct {
	Name  string
	Image *Image
}

func IfErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func UsageErr(s string) {
	fmt.Fprintln(os.Stderr, errors.New(s))
	flag.Usage()
	os.Exit(1)
}

func main() {
	serviceName = flag.String("name", "", "The name of the service")
	image = flag.String("image", "", "The name of the image and tag (i.e., `nginx:latest`) to use for the service")
	dbName = flag.String("dbName", "", "The name of the database")
	flag.Parse()

	if *serviceName == "" || *image == "" || *dbName == "" {
		UsageErr("[ERROR] Must define all parameters.")
	}

	imageInfo := strings.Split(*image, ":")
	if len(imageInfo) != 2 {
		UsageErr("[ERROR] Image must be in the format of `name:tag`.")
	}

	s := &Service{
		Name: *serviceName,
		Image: &Image{
			Name: imageInfo[0],
			Tag:  imageInfo[1],
		},
	}

	t := template.Must(template.ParseFS(templateFiles, "tpl/*"))
	for _, tpl := range t.Templates() {
		tplName := tpl.Name()
		tmpl, err := template.New(tplName).ParseFiles(fmt.Sprintf("tpl/%s", tplName))
		IfErr(err)
		dirPath := filepath.Join("build", s.Name)
		dirs := strings.Split(tplName, "_")
		filename := dirs[len(dirs)-1]
		// The `tplName` is made up of dirPath-filename.
		// For example:
		// 		base_deployment 		     ->  base/deployment
		//		overlays_beta_kustomization  ->  overlays/beta/kustomization
		// So, in the loop below, only construct the directory path from
		// 0 - N-1 (N-1 being the filename, of course).
		for _, d := range dirs[0 : len(dirs)-1] {
			dirPath = filepath.Join(dirPath, d)
		}
		err = os.MkdirAll(dirPath, os.ModePerm)
		IfErr(err)
		if filename != "env" {
			filename += ".yaml"
		}
		f, err := os.Create(filepath.Join(dirPath, filename))
		IfErr(err)
		IfErr(tmpl.Execute(f, *s))
	}
}
