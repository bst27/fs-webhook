# About
This app sends webhooks for filesystem changes. You can define an URL to receive webhooks and a filesystem path
to be monitored for changes.

# Usage
```
./fs-webhook -url "<URL>" -path "<PATH>"
```

# Example
Execute the application with the following config:
```
./fs-webhook -url "http://localhost/webhook-receiver" -path "/tmp/hello-world"
```

When you make a change to `hello-world` you will receive the following http post request at your defined url:
```
POST /webhook-receiver HTTP/1.1
Host: localhost
Content-Length: 56
Content-Type: application/json

{
   "action": "write",
   "path": "/tmp/hello-world"
}
```

# Build
To build executables for multiple platforms you can use the build script at `scripts/build.sh`.