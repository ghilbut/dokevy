# Local development infrastructures

## A. Elastic APM environment

Version: 8.6.2

### Run elasticsearch and kibana

```shell
cd es
docker-compose up -d
```

### Run apm-server

```shell
cd es
docker-compose -f docker-compose-apm.yml up -d
```

### Setup apm-server on kibana UI

http://localhost:5601

1. Select `[Observability] APM` menu →
2. Click `Add the APM integration` button
3. Click `Add Elastic APM` button
4. Fill inputs 
   * **Integration name:** APM
   * **Host:** apm-server:8200
   * **URL:** http\://apm-server:8200
   * **New agent policy name:** APM Agent policy
6. Click `Save and continue` button
7. Click `Add Elastic Agent later` button

### Run fleet server

## B.

---

나중에 참고해 볼 것
* https://github.com/elastic/apm-server/blob/v8.6.2/docker-compose.yml