package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/tc-hib/winres"
)

func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q: %w", path, err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read %q: %w", path, err)
	}
	return b, nil
}

type Number struct {
	Value *uint64
}

func (v Number) String() string {
	if v.Value == nil {
		return ""
	}
	return fmt.Sprintf("0x%04x (%d)", *v.Value, *v.Value)
}

func (v Number) Set(s string) error {
	u, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return err
	}
	*v.Value = u
	return nil
}

type Identifier struct {
	ID      *winres.Identifier
	ForType bool
}

func (v Identifier) String() string {
	if v.ID == nil || *v.ID == nil {
		return ""
	}
	return fmt.Sprint(*v.ID)
}

func (v Identifier) Set(s string) error {
	u, err := strconv.ParseUint(s, 0, 64)
	if err == nil && u <= 0xffff {
		*v.ID = winres.ID(u)
		return nil
	}
	if v.ForType {
		switch s {
		case "RT_CURSOR":
			*v.ID = winres.RT_CURSOR
		case "RT_BITMAP":
			*v.ID = winres.RT_BITMAP
			return nil
		case "RT_ICON":
			*v.ID = winres.RT_ICON
			return nil
		case "RT_MENU":
			*v.ID = winres.RT_MENU
			return nil
		case "RT_DIALOG":
			*v.ID = winres.RT_DIALOG
			return nil
		case "RT_STRING":
			*v.ID = winres.RT_STRING
			return nil
		case "RT_FONTDIR":
			*v.ID = winres.RT_FONTDIR
			return nil
		case "RT_FONT":
			*v.ID = winres.RT_FONT
			return nil
		case "RT_ACCELERATOR":
			*v.ID = winres.RT_ACCELERATOR
			return nil
		case "RT_RCDATA":
			*v.ID = winres.RT_RCDATA
			return nil
		case "RT_MESSAGETABLE":
			*v.ID = winres.RT_MESSAGETABLE
			return nil
		case "RT_GROUP_CURSOR":
			*v.ID = winres.RT_GROUP_CURSOR
			return nil
		case "RT_GROUP_ICON":
			*v.ID = winres.RT_GROUP_ICON
			return nil
		case "RT_VERSION":
			*v.ID = winres.RT_VERSION
			return nil
		case "RT_PLUGPLAY":
			*v.ID = winres.RT_PLUGPLAY
			return nil
		case "RT_VXD":
			*v.ID = winres.RT_VXD
			return nil
		case "RT_ANICURSOR":
			*v.ID = winres.RT_ANICURSOR
			return nil
		case "RT_ANIICON":
			*v.ID = winres.RT_ANIICON
			return nil
		case "RT_HTML":
			*v.ID = winres.RT_HTML
			return nil
		case "RT_MANIFEST":
			*v.ID = winres.RT_MANIFEST
			return nil
		}
	}
	if s != "" {
		*v.ID = winres.Name(s)
	}
	return nil
}

func process() error {
	var (
		binfile = flag.String("bin", "", "resource binary filepath")
		infile  = flag.String("in", "", "input executable filepath")
		outfile = flag.String("out", "", "output executable filepath")
		typeID  winres.Identifier
		resID   winres.Identifier
		langID  uint64 = 0x0409
	)
	flag.Var(Identifier{&typeID, true}, "type", "resource type id")
	flag.Var(Identifier{&resID, false}, "res", "resource id")
	flag.Var(Number{&langID}, "lang", "language id")
	flag.Parse()

	if typeID == nil {
		return fmt.Errorf("type is required")
	}
	if resID == nil {
		return fmt.Errorf("res is required")
	}
	if langID > 0xffff {
		return fmt.Errorf("%x is invalid language ID", langID)
	}

	bin, err := readFile(*binfile)
	if err != nil {
		return fmt.Errorf("failed to read resource binary file: %w", err)
	}
	r, err := readFile(*infile)
	if err != nil {
		return fmt.Errorf("failed to read executable file: %w", err)
	}
	rs, err := winres.LoadFromEXE(bytes.NewReader(r))
	if err != nil {
		return fmt.Errorf("failed to parse executable: %w", err)
	}
	if err = rs.Set(typeID, resID, uint16(langID), bin); err != nil {
		return fmt.Errorf("failed to set resource: %w", err)
	}
	w, err := os.Create(*outfile)
	if err != nil {
		return fmt.Errorf("failed to create executable file: %w", err)
	}
	defer w.Close()
	err = rs.WriteToEXE(w, bytes.NewReader(r))
	if err != nil {
		return fmt.Errorf("failed to write executable file: %w", err)
	}
	return nil
}

func main() {
	if err := process(); err != nil {
		fmt.Printf("%s - tiny windows executable resource rewriter\n\n", os.Args[0])
		flag.Usage()
		fmt.Println()
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
