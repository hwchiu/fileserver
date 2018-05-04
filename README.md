Fileserver [![Build Status](https://travis-ci.org/hwchiu/fileserver.svg?branch=master)](https://travis-ci.org/hwchiu/fileserver)
============
A simple www server which is written in golang supports the scan/read/write/delete function

Usage
=====
### Build the fileserver binary
You should install the golang environment into your system first.
Type the following command to build the binary
```sh
make build
```

### Run the fileserver
The fileserverlisten on `localhost:33333` by default and you should use the `-documentRoot` to change the root of all files path operations.
For example, if the `-documentRoot` is `/tmp`, the scan operation will scan all directories/files under /tmp.
Use the following command to run a fileserver
``` sh
./fileserver -documentRoot /tmp
```

After the fileserver is running, we can use the curl to send the HTTP request.
### Scan
the basic URL usage is `/scan/{path}` and the fileserver will scan all directories/files under `-documentRoot`/{path}.

For example 
`curl -X GET -i http://localhost:33333/scan` will scan the directory /tmp and `curl -X GET -i http://localhost:33333/scan/data` will scan /tmp/data/.
