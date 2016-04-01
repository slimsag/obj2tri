package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("usage: obj2tri <input.obj>")
	}

	objs, err := Read(args[0])
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(strings.TrimSuffix(args[0], filepath.Ext(args[0])) + ".tri")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "1\n"); err != nil {
		log.Fatal(err)
	}
	for _, obj := range objs {
		for _, group := range obj.Groups {
			for i := 0; i < len(group.Vertexes)/3; i++ {
				i3 := i * 3
				v := group.Vertexes[i3 : i3+3]
				n := group.Normals[i3 : i3+3]

				var t []float32
				var tf = "ground.tga"
				if len(group.Textcoords) > 0 {
					i2 := i * 2
					t = group.Textcoords[i2 : i2+2]
					tf = filepathBase(group.Material.Texturefile)
				} else {
					t = []float32{0, 0}
				}
				// vx, vy, vz, nx, ny, nz, uvx, uvy, name
				if _, err := fmt.Fprintf(f, "%v %v %v %v %v %v %v %v %s\n", v[0], v[1], v[2], n[0], n[1], n[2], t[0], t[1], tf); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func filepathBase(f string) string {
	fv := strings.Split(f, "/")
	if len(fv) > 1 {
		return fv[len(fv)-1]
	}
	fv = strings.Split(f, `\`)
	return fv[len(fv)-1]
}
