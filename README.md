# go-replace-tag
Replace tags in golang structs

## Install

```
$ go get github.com/adityadani/go-replace-tag
```

## Usage

```
$ go-replace-tag -input <path-to-pb.go>
```

Add the following comments to the struct fields which need a tag to be replaced
`// @replace-tag <tagkey>:"<new tagvalue>"`

Example:

```
api.proto
syntax = "proto3";

// Example for replace tags
message Example {
  // Field1
  string Field1 = 0;
  // Field2
  // @replace-tag json:"-"
  string Field2 = 1;
}
```

Generate the api.pb.go file
```
$ protoc api.proto
```

Replace the tags
```
$ go-replace-tag -input api.pb.go
```

Output file: 

```
api.pb.go

type Example struct {
     // Field1
     Field1 string `protobuf:"bytes,15,opt,name=field1" json:"field1,omitempty"`
     // Field2
     // @replace-tag json:"-"
     Field2 string `protobuf:"bytes,15,opt,name=field2" json:"-"`
}
```
