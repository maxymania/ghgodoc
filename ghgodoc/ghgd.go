package main

import (
	"go/build"
	//"go/doc"
	
	"github.com/maxymania/ghgodoc/ghim"
	"github.com/maxymania/ghgodoc/ghdc"
	"flag"
	"log"
	
	"os"
	"path/filepath"
	"strings"
)

var pspkg = flag.String("pkg","","golang package")
var ghpdir = ""

func p2d(pn string) string{
	if filepath.Separator=='/' { return pn }
	pna := []byte(pn)
	for i,b := range pna {
		if b!='/' { continue }
		pna[i] = filepath.Separator
	}
	return string(pna)
}

func document() {
	pkg,err := build.Default.Import(*pspkg,"",build.FindOnly)
	if err!=nil { log.Fatal(err) }
	fset,apkg,err := ghim.Parse(pkg,*pspkg)
	if err!=nil { log.Fatal(err) }
	
	targp := filepath.Join(ghpdir,p2d(*pspkg))
	err = os.MkdirAll(targp,0755)
	if err!=nil { log.Fatal(err) }
	
	if txt,err := os.Create(filepath.Join(targp,"pkg.txt")); err==nil {
		pkgtxt := apkg.Doc
		i := strings.Index(pkgtxt,".")
		if i>0 { pkgtxt = pkgtxt[:i+1] }
		txt.Write([]byte(pkgtxt))
		txt.Close()
	}
	if _,err := os.Stat(filepath.Join(targp,"list.html")); err!=nil && os.IsNotExist(err) {
		if lst,err := os.Create(filepath.Join(targp,"list.html")); err==nil {
			lst.Write([]byte(`<!-- No list yet -->`))
			lst.Close()
		}
	}
	
	dcb := new(ghdc.Builder)
	dcb.Target.WriteString("---\nlayout: godoc\n")
	dcb.Target.WriteString("title: "+apkg.Name+"\n")
	dcb.Target.WriteString("gopkg: "+apkg.ImportPath+"\n")
	dcb.Target.WriteString("---\n")
	dcb.Generate(fset,apkg)
	dcb.Target.WriteString("\n{% include_relative list.html %}\n")
	fobj,err := os.Create(filepath.Join(targp,"index.html"))
	if err!=nil { log.Fatal(err) }
	defer fobj.Close()
	dcb.Target.WriteTo(fobj)
}

func main() {
	flag.Parse()
	ghpdir = os.Getenv("GHPKG")
	fi,err := os.Stat(ghpdir)
	if err!=nil { log.Fatal(err) }
	if !fi.IsDir() { log.Fatal(fi,"is not a directory") }
	
	if *pspkg!="" {
		document()
	}
}

