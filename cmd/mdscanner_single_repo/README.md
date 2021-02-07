TODO: Write about what you're doing here

- To build and run:
```shell
docker build . -t mdscanner_single_repo && docker run -e repo=https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/alb-ingress.md mdscanner_single_repo
```

- Module setup: `go mod init DanielCalvo/mdscanner/cmd/mdscanner_cmdline` <- I am not sure if this is correct