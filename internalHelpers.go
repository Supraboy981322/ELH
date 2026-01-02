package ELH 

import (
	"os"
	"fmt"
	"reflect"
	"strings"
	"strconv"
)

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
