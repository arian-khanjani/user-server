
# user-server
A simple gRPC server practice using MongoDB with repository design pattern

## Layers
- Controller
- Repository Layer

## Features
- CRUD operations
- Indexing
- gRPC registry support
- Graceful shutdown

### Database
To spin up the database, use the provided `docker-compose.yml` file:
```bash
docker-compose up -d
```

### Generate .proto files 
To generate the protocol buffer files using `buf` (should already be installed):
```bash
buf generate && buf generate --template buf.gen.tag.yaml
```

