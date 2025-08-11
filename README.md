# sandbox

## How to start

Copy config

```shell
cp ./configs/config.yaml.example ./configs/config.yaml 
```

Run application

```shell
make run
```

Send request

```shell
curl -X GET http://127.0.0.1:8080/api/v1/health
```

```shell
make stress-test
```