package ghdc

import (
	"go/doc"
	"go/token"
	"go/printer"
	
	"bytes"
	"fmt"
)

type Builder struct {
	Target bytes.Buffer
	
	temp  bytes.Buffer
	fset *token.FileSet
	words map[string]string
}
var mprinter = &printer.Config{Tabwidth: 8, Indent: 1}
func (bdr *Builder) toString(aste interface{}) (string) {
	bdr.temp.Reset()
	mprinter.Fprint(&bdr.temp,bdr.fset,aste) // XXX: we ignore any error here.
	return bdr.temp.String()
}

func ne(s string) string {
	if s=="" { return "-" }
	return s
}

func (bdr *Builder) renderIndexValues(list []*doc.Value, title string) {
	if len(list)==0 { return }
	fmt.Fprintf(&bdr.Target,"<li><a href=\"#pkg-%s\">%s</a></li>",title,title)
}
func (bdr *Builder) renderIndexFuncs(list []*doc.Func, prefix string, ism bool) {
	for _,fn := range list {
		
		//receiver := ""
		//if ism { receiver = fmt.Sprintf("(%s) ",fn.Recv) }
		/*
		fmt.Fprintf(&bdr.Target,
			"<li><a href=\"#%s%s\">func %s%s</a></li>\n",
			prefix,
			fn.Name,
			receiver,
			fn.Name)
		*/
		fmt.Fprintf(&bdr.Target,
			"<li><a href=\"#%s%s\">%s</a></li>\n",
			prefix,
			fn.Name,
			bdr.toString(fn.Decl) )
	}
}
func (bdr *Builder) renderIndexTypes(list []*doc.Type) {
	for _,tp := range list {
		fmt.Fprintf(&bdr.Target,
			"<li><a href=\"#%s\">type %s</a><br><ul>\n",tp.Name,tp.Name)
		bdr.renderIndexFuncs (tp.Funcs  , ""         , false )
		bdr.renderIndexFuncs (tp.Methods, tp.Name+".", true  )
		fmt.Fprintf(&bdr.Target,"</ul></li>\n")
	}
}

func (bdr *Builder) renderValues(list []*doc.Value, title string, l3 bool) {
	if len(list)==0 { return }
	if l3 {
		fmt.Fprintf(&bdr.Target,"<h3>%s</h3>\n",title)
	} else {
		fmt.Fprintf(&bdr.Target,"<h2 id=\"pkg-%s\">%s</h2>\n",title,title)
	}
	for _,val := range list {
		fmt.Fprintf(&bdr.Target,"<div>\n")
		doc.ToHTML(&bdr.Target,bdr.toString(val.Decl)+"\n"+ne(val.Doc),bdr.words)
		fmt.Fprintf(&bdr.Target,"</div>\n")
	}
}

func (bdr *Builder) renderFuncs(list []*doc.Func, prefix string, ism bool) {
	for _,fn := range list {
		receiver := ""
		if ism { receiver = fmt.Sprintf("(%s) ",fn.Recv) }
		fmt.Fprintf(&bdr.Target,
			"<h2 id=\"%s%s\">func %s%s</h2>\n",
			prefix,
			fn.Name,
			receiver,
			fn.Name)
		fmt.Fprintf(&bdr.Target,"<div>\n")
		doc.ToHTML(&bdr.Target,bdr.toString(fn.Decl)+"\n"+ne(fn.Doc),bdr.words)
		fmt.Fprintf(&bdr.Target,"</div>\n")
	}
}

func (bdr *Builder) renderTypes(list []*doc.Type) {
	for _,tp := range list {
		fmt.Fprintf(&bdr.Target,
			"<h2 id=\"%s\">type %s</h2>\n",tp.Name,tp.Name)
		fmt.Fprintf(&bdr.Target,"<div>\n")
		doc.ToHTML(&bdr.Target,bdr.toString(tp.Decl)+"\n"+ne(tp.Doc),bdr.words)
		fmt.Fprintf(&bdr.Target,"</div>\n")
		bdr.renderValues(tp.Consts , "Constants", true  )
		bdr.renderValues(tp.Vars   , "Variables", true  )
		bdr.renderFuncs (tp.Funcs  , ""         , false )
		bdr.renderFuncs (tp.Methods, tp.Name+".", true  )
	}
}

func (bdr *Builder) fillTypeMap(pd *doc.Package) {
	bdr.words = make(map[string]string)
	for _,tp := range pd.Types {
		bdr.words[tp.Name] = "#"+tp.Name
	}
}

func (bdr *Builder) renderHeader(pd *doc.Package) {
	fmt.Fprintf(&bdr.Target,"<h1>Package %s</h1>\n",pd.Name)
	fmt.Fprintf(&bdr.Target,"<code>import %q</code>",pd.ImportPath)
	fmt.Fprintf(&bdr.Target,"<h2>Overview</h2>\n")
	fmt.Fprintf(&bdr.Target,"<div>\n")
	doc.ToHTML(&bdr.Target,pd.Doc,nil)
	fmt.Fprintf(&bdr.Target,"</div>\n")
}
func (bdr *Builder) renderIndex(pd *doc.Package) {
	fmt.Fprintf(&bdr.Target,"<h2>Index</h2>\n<ul>")
	bdr.renderIndexValues(pd.Consts, "Constants")
	bdr.renderIndexValues(pd.Vars,   "Variables")
	bdr.renderIndexFuncs (pd.Funcs,  "",false)
	bdr.renderIndexTypes (pd.Types)
	fmt.Fprintf(&bdr.Target,"</ul>\n")
}
func (bdr *Builder) renderBodies(pd *doc.Package) {
	
	bdr.renderValues(pd.Consts, "Constants", false)
	bdr.renderValues(pd.Vars,   "Variables", false)
	bdr.renderFuncs (pd.Funcs,  "",false)
	bdr.renderTypes (pd.Types)
}
func (bdr *Builder) renderFooter(pd *doc.Package) {
	if len(pd.Imports)==0 { return }
	fmt.Fprintf(&bdr.Target,"<h2>Dependencies</h2><ul>\n")
	for _,imp := range pd.Imports {
		fmt.Fprintf(&bdr.Target,"<li><code>import %q</code></li>",imp)
	}
	fmt.Fprintf(&bdr.Target,"</ul>\n")
}
func (bdr *Builder) Generate(fset *token.FileSet, pd *doc.Package) {
	bdr.fset = fset
	//bdr.fillTypeMap(pd)
	bdr.renderHeader(pd)
	bdr.renderIndex (pd)
	bdr.renderBodies(pd)
	bdr.renderFooter(pd)
}



// ----
