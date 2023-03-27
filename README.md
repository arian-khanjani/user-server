
# user-server
A simple gRPC server practice using MongoDB with repository design pattern


## Layers
- Controller
- Repository Layer


## Generate .proto files 
To generate the protocol buffer files:
```bash
buf generate && buf generate --template buf.gen.tag.yaml
```


## Features
- CRUD Operations
- Indexing
- gRPC registry support

