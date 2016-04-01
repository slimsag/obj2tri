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

				_, tf := filepath.Split(group.Material.Texturefile)
				if tf == "" {
					tf = "ground.tga"
				}

				var t []float32
				if len(group.Textcoords) > 0 {
					i2 := i * 2
					t = group.Textcoords[i2 : i2+2]
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
