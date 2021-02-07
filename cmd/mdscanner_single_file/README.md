TODO: Write about what you're doing here

- To build and run:
```shell
docker build . -t mdscanner_single_file && docker run -e file-url=https://github.com/kubernetes/community/blob/master/README.md mdscanner_single_file
```

- Module setup, I am not sure if this is correct:
```shell
go mod init DanielCalvo/mdscanner/cmd/mdscanner_single_file
go mod init
```
