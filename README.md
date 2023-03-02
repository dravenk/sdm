## SDM

Software Deployment Manager

#### Install

Install the latest version
```shell
go install github.com/dravenk/sdm@latest
```

#### Usage

Initialize the configuration file with custom parameters.
```shell
sdm init --minport 8100 --maxport 8200 --image dravenk/dp:10-fpm --dir ../apps --app=d1,d2,d3
```

Create and start containers.
```shell
sdm up
```

Stop and remove containers, networks.
```shell
sdm down
```

Delete configuration files.
```shell
sdm clean
```
