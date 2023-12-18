# scale-codec

Golang implementation of parity-scale-codec without reflect

#### Generating Enums

Create the enum in a `.scale` file

```
// simple_enum.scale

enum MyScaleEncodedEnum {
    Int(uint64)
    Bool(bool)
}
```

Download the `enum_script` CLI tool, and include the following script to generate the enums following the enums grammar

```
//go:generate enum_script simple_enum.scale main
```

The tool will generate a `.go` file with the same name, the file contains the enum definitions and method to scale encode/decode the enum

For more info check the following directory `tests/enums`
