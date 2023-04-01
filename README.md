# minirh

minirh is a tiny windows executable resource rewriter.

## USAGE

```txt
minirh.exe - windows executable resource rewriter

Usage of minirh.exe:
  -bin string
        resource binary filepath
  -in string
        input executable filepath
  -lang value
        language id (default 0x0409 (1033))
  -out string
        output executable filepath
  -res value
        resource id
  -type value
        resource type id
```

Use it like this:

```
minirh.exe -in foo.exe -out foo_modified.exe -type RT_RCDATA -res YOURDATA -lang 0x0411 -bin someresource.bin
```
