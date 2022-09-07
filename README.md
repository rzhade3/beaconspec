# beaconspec

beaconspec is a Golang implementation of a [beacon specification](https://gbv.github.io/beaconspec/beacon.html) parser. The Beacon spec is a common format for URL dumps.

It is designed to be used as a library to parse Beacon spec files.

### Usage

```golang
import "github.com/rzhade3/beaconspec"

func main() {
    metadata, err := beaconspec.ParseMetadata("/path/to/beacon_file.txt")
    if err != nil {
        panic(err)
    }

    // Since the file might be massive, parse the file yourself and pass
    // the contents of a single line to ParseLine
    line := "foo|bar|baz"
    record, err := beaconspec.ParseLine(line, metadata)
    if err != nil {
        panic(err)
    }
    fmt.Println(record.Source)
    fmt.Println(record.Annotation)
    fmt.Println(record.Target)
}
```
