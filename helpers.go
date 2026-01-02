package ELH

import (
	"os"
	"fmt"
	"time"
	"bytes"
	"errors"
	"strconv"
	"strings"
	"reflect"
	"net/http"
	"path/filepath"
)

// render src with specific registery
func RenderWithRegistry(src string, registry map[string]Runner, r *http.Request, wr http.ResponseWriter) (string, error) {
	return parseAndRun(src, registry, r, wr)
}

// wrapper that uses the DefaultRegistry.
func Render(src string, r *http.Request, wr http.ResponseWriter) (string, error) {
	return RenderWithRegistry(src, DefaultRegistry(), r, wr)
}

func BldApiStruct(req *http.Request, wr http.ResponseWriter) API {
	params := map[string]string{}
	for k, _ := range req.URL.Query() {
		params[k] = req.URL.Query().Get(k)
	}
	api := API{
		Req: reqApi{
			Method: req.Method,
			IsTLS: false,
			URL: urlApi{
				Host: req.Host,
				URI: req.RequestURI,
				Path: req.URL.Path,
				Params: params,
			},
		},
	}
	
	if req.TLS != nil { api.Req.IsTLS = true }

	return api
}

func readStruct(v interface{}, path string) (reflect.Value, error) {
	if path == "" {
		return reflect.Value{}, fmt.Errorf("%w: empty path", ErrInvalidPath)
	}
	rv := reflect.ValueOf(v)
	if !rv.IsValid() || rv.Kind() == reflect.Ptr && rv.IsNil() {
		return reflect.Value{}, fmt.Errorf("%w: nil or invalid input", ErrInvalidPath)
	}

	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	parts := strings.Split(path, ".")
	for _, p := range parts {
		if p == "" {
			return reflect.Value{}, fmt.Errorf("%w: empty path segment", ErrInvalidPath)
		}
		
		switch rv.Kind() {
		 case reflect.String:
			rv = reflect.Value(rv).Elem().FieldByName(p)
		 case reflect.Struct:
			rv = rv.FieldByName(p)
			if !rv.IsValid() {
				err := fmt.Errorf("%w: field %q not found", ErrInvalidPath, p)
				return reflect.Value{}, err
			}
		 case reflect.Map:
			k := reflect.ValueOf(p)
			rv = rv.MapIndex(k)
			if !rv.IsValid() {
				err := fmt.Errorf("%w: key %q isn't in map", ErrInvalidPath, p)
				return reflect.Value{}, err
			}
		 case reflect.Slice, reflect.Array:
			i, err := strconv.Atoi(p)
			if err != nil {
				err = fmt.Errorf("%w: %q is NaN", ErrInvalidPath, p)
				return reflect.Value{}, err
			}
			if i < 0 || i >= rv.Len() {
				err = fmt.Errorf("%w: index %q out of range", ErrInvalidPath, i)
				return reflect.Value{}, err
			}
			rv = rv.Index(i)
		 case reflect.Ptr:
			for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
				if rv.IsNil() { break }
				rv = rv.Elem()
			}
		 default:
			err := fmt.Errorf("%w: can't traverse into %s"+
						"(not yet implemented)", ErrInvalidPath, rv.Kind())
			return reflect.Value{}, err
		}
	}
	return rv, nil
}

func MkReg(caller string, cmd string, args []string, timeout int, env []string) map[string]Runner {
	reg := map[string]Runner {
		caller: &ExternalRunner {
			CmdName: cmd,
			Args:    args,
			Timeout: time.Duration(timeout) * time.Second,
			Env:     env,
		},
	}
	return reg
}

func MkRegDefaults(cmd string, args []string) map[string]Runner {
	reg := MkReg(cmd, cmd, args, 5, os.Environ())
	return reg
}

func Serve(w http.ResponseWriter, r *http.Request) (string, error) {
	//get the requested file
	file := r.URL.Path
	if file == "/" {
		file = "index"
	} else if file[len(file)-1:] ==  "/" {
		file = fmt.Sprintf("%sindex", string(file[1:]))
		file, _ = checkIsDir(file)
	} else {
		file = file[1:]
		file, _ = checkIsDir(file)
	}

	//get the extension of the requested file
	ext := filepath.Ext(file)	
	if ext == "" { //if there is no ext
		//check against list of ext which can
		//  have no ext in url
		for i := 0; i < len(suppNoExt); i++ {
			checkFile := fmt.Sprintf("%s%s", file, suppNoExt[i])
			_, err := os.Stat(checkFile)
			if err == nil { //if the file exists
				file = checkFile //assume it's the correct one
				ext = suppNoExt[i]
				break
			} else if !errors.Is(err, os.ErrNotExist) {
				return file, errors.New("cannot check if file exists! Schrodinger's file:  "+err.Error())
			}
		}
	}
	if fileExists(file) {
		fileByte, err := os.ReadFile(file)
		if err != nil {
			return file, errors.New("read file:  "+err.Error())
		}
		fileStr := string(fileByte)
		var result string
		//if the file is elh, parse it
		if ext == ".elh" {
			result, err = Render(fileStr, r, w)
			if err != nil {
				var fileStr string
				if Server.ErrPage != nil {
					fileStr, err = Render(string(Server.ErrPage), r, w) 
					if err != nil {
						http.Error(w, "500 server err", 500)
						return "500", err
					}
				} else { fileStr = "500 server err" }
				http.Error(w, fileStr, 500)
				file = "500 server err" 
				return file, errors.New("elh failed; "+err.Error())
			}
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintln(w, result)
		} else {
			fileReader := bytes.NewReader(fileByte)
			http.ServeContent(w, r, file, time.Now(), fileReader)
		}
	} else {
		var err error
		var fileStr string
		if Server.ErrPage != nil {
			fileStr, err = Render(string(Server.ErrPage), r, w) 
			if err != nil {
				http.Error(w, "500 server err", 500)
				return "500", err
			}
		} else { fileStr = "404 forbidden" }
		http.Error(w, fileStr, 404)
		file = "404 forbidden" 
	}
	return file, nil
}


//returns error so Go doesn't panic, but error
//  is ignored when call fn
//   (it's handled later, after fn call) 
func checkIsDir(file string) (string, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return file, err
	}
	if fileInfo.IsDir() {
		file = fmt.Sprintf("%s/index", file)
	}
	return file, nil
}

func ServeWithRegistry(w http.ResponseWriter, r *http.Request, registry map[string]Runner) (string, error) {
	//get the requested file
	file := r.URL.Path
	//set the system path of file
	if file == "/" {
		file = "index"
	} else if file[len(file)-1:] ==  "/" {
		file = fmt.Sprintf("%sindex", string(file[1:]))
		file = filepath.Join(Server.WebDir, file)
		file, _ = checkIsDir(file)
	} else {
		file = file[1:]
		file = filepath.Join(Server.WebDir, file)
		file, _ = checkIsDir(file)
	}

	//get the extension of the requested file
	ext := filepath.Ext(file)	
	if ext == "" { //if there is no ext
		//check against list of ext which can
		//  have no ext in url
		for i := 0; i < len(suppNoExt); i++ {
			checkFile := fmt.Sprintf("%s%s", file, suppNoExt[i])
			_, err := os.Stat(checkFile)
			if err == nil { //if the file exists
				file = checkFile //assume it's the correct one
				ext = suppNoExt[i]
				break
			} else if !errors.Is(err, os.ErrNotExist) {
				return file, errors.New("cannot check if file exists! Schrodinger's file:  "+err.Error())
			}
		}
	}
	if fileExists(file) {
		fileByte, err := os.ReadFile(file)
		if err != nil {
			return file, errors.New("read file:  "+err.Error())
		}
		fileStr := string(fileByte)
		var result string
		//if the file is elh, parse it
		if ext == ".elh" {
			result, err = RenderWithRegistry(fileStr, registry, r, w)
			if err != nil {
				var fileStr string
				if Server.ErrPage != nil {
					fileStr, err = Render(string(Server.ErrPage), r, w) 	
					if err != nil {
						http.Error(w, "500 server err", 500)
						return "500", err
					}
				} else { fileStr = "500 server err" }
				http.Error(w, fileStr, 500)
				file = "500 server err" 
				return file, errors.New("elh failed; "+err.Error())
			}
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintln(w, result)
		} else {
			fileReader := bytes.NewReader(fileByte)
			http.ServeContent(w, r, file, time.Now(), fileReader)
		}
	} else {
		var err error
		var fileStr string
		if Server.ErrPage != nil {
			fileStr, err = Render(string(Server.ErrPage), r, w) 
			if err != nil {
				http.Error(w, "500 server err", 500)
				return "500", err
			}
		} else { fileStr = "404 forbidden" }
		http.Error(w, fileStr, 404)
		file = "404 forbidden" 
	}
	return file, nil
}

func HttpServer(w http.ResponseWriter, r *http.Request) {
	var resp string
	//get the requested file
	file := filepath.Join(Server.WebDir, r.URL.Path)
	if file == Server.WebDir {
		file = "index"
	} else if file[len(file)-1:] ==  "/" {
		file = fmt.Sprintf("%sindex", string(file[1:]))
		file, _ = checkIsDir(file)
	} else {
		file = file[1:]
		file, _ = checkIsDir(file)
	}

	//get the extension of the requested file
	ext := filepath.Ext(file)	
	if ext == "" { //if there is no ext
		//check against list of ext which can
		//  have no ext in url
		for i := 0; i < len(suppNoExt); i++ {
			checkFile := fmt.Sprintf("%s%s", file, suppNoExt[i])
			_, err := os.Stat(checkFile)
			if err == nil { //if the file exists
				file = checkFile //assume it's the correct one
				ext = suppNoExt[i]
				break
			} else if !errors.Is(err, os.ErrNotExist) {
				http.Error(w, "cannot check if file exists! Schrodinger's file:  "+err.Error(), http.StatusInternalServerError)
				resp = "500; Schrodinger's file"
			}
		}
	}
	if fileExists(file) {
		fileByte, err := os.ReadFile(file)
		if err != nil {
			http.Error(w, "read file:  "+err.Error(), http.StatusInternalServerError)
			resp = "500; err reading file"
		}
		fileStr := string(fileByte)
		var result string
		//if the file is elh, parse it
		if ext == ".elh" {
			result, err = Render(fileStr, r, w)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				fmt.Fprintf(w, "There appears to be an error in the `.elh` file %s", file)
				resp = "500; problem with elh file"
			}
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintln(w, result)
		} else {
			fileReader := bytes.NewReader(fileByte)
			http.ServeContent(w, r, file, time.Now(), fileReader)
		}
	} else {
		var err error
		var fileStr string
		if Server.ErrPage != nil {
			fileStr, err = Render(string(Server.ErrPage), r, w)
			if err != nil {
				http.Error(w, "500 server err", 500)
				return
			}
		} else { fileStr = "404 forbidden" }
		http.Error(w, fileStr, 404)
		file = "404 forbidden" 
	}

	//colorize response log string
	if resp == "" { resp = "\033[32m"+file+"\033[0m"
	} else { resp = "\033[31m"+resp+"\033[0m" }

	//colorize file log string
	file = "\033[35m"+file+"\033[0m"

	//check if logger is set
	if Server.Log.Func != nil {
		//build string
		logStr := "\033[1m[req]:\033[0m "
		logStr += file+" | "
		logStr += "\033[1m[resp]:\033[0m "+resp
		//log it
		Server.Log.Func(logStr)
	}
}

func UseDefault(name string) *ExternalRunner {
	reg := DefaultRegistry()[name]
	if reg == nil {
		if Server.Log.Func != nil {
			Server.Log.Func("runner not in registry:  \n"+name)
			os.Exit(1)
		}
		fmt.Fprintf(syserr, "runner not in registry:  %s\n", name)
		os.Exit(1)
	}
	return reg.GetRunner()
}

func ToArgs(src string, by string) []string {
	if len(src) <= 1 { return []string{} }
	if src[len(src)-1] == ' ' { src = src[:len(src)] }
	if len(src) <= 1 { return []string{} }
	type argParser struct{
		quot bool
		esc bool
		pos int
		mem string
		out []string
		in []string
	};p := argParser{
		quot: false,
		esc: false,
		pos: 0,
		mem: "",
		out: []string{},
		in: strings.Split(src, ""),
	}

	var foo func(p argParser) []string
	foo = func(p argParser) []string {
		if p.pos >= len(p.in) {
			p.out = append(p.out, p.mem)
			return p.out
		}
		cur := p.in[p.pos]
		if p.esc || p.quot {
			if p.quot && (cur == `"` || cur == `'`) {
				p.pos++ ; p.quot = false ; return foo(p)
			}
			p.mem += cur
			p.pos++
			if p.esc { p.esc = false }
			return foo(p)
		}
		switch cur {
     case `\`: p.esc = true
	   case `"`, `'`: p.quot = true
		 case by:
			p.out = append(p.out, p.mem)
			p.mem = ""
		 default: p.mem += cur
		}
		p.pos++
		return foo(p)
	}
	return foo(p)
}

func (opts *RunOpts) Unwrap() (string, string, *os.File, *http.Request, http.ResponseWriter) {
	code := opts.Code
	tmpDir := opts.TmpDir
	tmpFile := opts.TmpFile
	req := opts.Req
	wr := opts.Wr
	return code, tmpDir, tmpFile, req, wr
}

func splitLines(src string, split_on []string) ([]string) {
	if split_on == nil { split_on = []string{} }
	split_on = append(split_on, "\n")
	type pa struct {
		mem string
		out []string
		pos int
		in []string
		esc bool
		split []string
	};p := pa{
		in: strings.Split(src, ""),
		split: split_on,
	}
	var foo func(p pa) []string
	foo = func(p pa) []string {
		if p.pos >= len(p.in) {
			p.mem = strings.TrimSpace(p.mem)
			if len(p.mem) != 0 {
				p.out = append(p.out, p.mem)
			}
			return p.out
		}
		switch p.in[p.pos] {
		 case "'", `"`:
			p.esc = !p.esc ; p.mem += p.in[p.pos]
		 default:
			var s bool
			for _, c := range p.split {
				if p.in[p.pos] == c {
					p.mem = strings.TrimSpace(p.mem)
					if len(p.mem) != 0 {
						p.out = append(p.out, p.mem)
					}
					p.mem = "" ; s = true ; break
				}
			};if !s { p.mem += p.in[p.pos] }
		}
		p.pos++
		return foo(p)
	}
	return foo(p)
}
