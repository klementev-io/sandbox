# sandbox

## How to start

Rename config.yaml.example -> config.yaml.

Envs values override config values.

```shell
make run
```

```shell
curl -X GET \
  http://127.0.0.1:8080/api/v1/health
```

```shell
make stress-test
```