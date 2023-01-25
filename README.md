# Test Task solution for Tages.

Test task from Tages company.

## Installation

Use the git cli to clone repository.

```bash
git clone https://github.com/zagart47/imageserver.git
```

## Usage
### Start server
```bash
go run api/server.go
```

## You need use the client to work with this server!!!
```html
https://github.com/zagart47/imageclient
```

### Upload Image RPC
```
rpc Upload(stream UploadRequest) returns (UploadResponse) {}

message UploadRequest {
  string filename = 1;
  bytes fragment = 2;
}
message UploadResponse {
  string name = 1;
}

```

### Get Image list RPC
```
rpc GetFiles(GetFilesRequest) returns (GetFilesResponse) {}

message GetFilesRequest {}
message GetFilesResponse{
  repeated File info = 1;
}
message File {
  string file_name = 1;
  string created = 2;
  string updated = 3;
}
```


### Download Image RPC:
```
rpc Download(DownloadRequest) returns (stream DownloadResponse) {}

message DownloadRequest {
  string filename = 1;
}

message DownloadResponse {
  string filename = 1;
  bytes fragment = 2;
}
```