# Event log producer

카프카 스트림 개발 및 테스트 목적의 Fake event log 를 생성하여 카프카 브로커로 전송하는 프로그램이다.

사용 언어는 GO 1.21로 개발되었다.

## Build
다음 명령으로 소스 코드를 빌드하여 log-producer 실행 파일을 생성한다.

```shell
$ go build -o log-producer cmd/log-producer/main.go
```

## Run
log-producer 실행 파일을 -c 옵선과 함께 환경 설정 파일을 지정하여 실행한다.
프로그램 종료는 프로그램을 실행한 터미널에서 CTRL+C를 눌러 프로그램을 종료한다. 

```shell
./log-producer -c config/config.yml
```

## Configuration

config/config.yml 파일에 프로그램 실행에 필요한 옵션들을 기술한다.

```yaml
number-of-worker: 3
kafka:
  bootstrap-servers: localhost:9092
  topic-name: EVENT_LOG
  partitions: 3
  properties:
    linger.ms: 9
    batch.size: 1000000
```

| 설정                      | 설명                                                                                                |
|-------------------------|---------------------------------------------------------------------------------------------------|
| number-of-worker        | 프로듀서의 갯수                                                                                          |
| kafak.bootstrap-servers | 카프카 브로커의 주소                                                                                       |
| kafka.topic-name        | 생성된 이벤트 로그를 전송할 토픽 명                                                                              |
| kafka.partitions        | kafka.topic-name에 설정한 토픽의 파티션 갯수                                                                  |
| kafka.propertis         | 프로듀서에 관한 설정으로 https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md 의 내용을 참고한다. |


## 이벤트 로그
생성되는 이벤트 로그 아래와 같은 형식의 JSON 문자열을 생성하여 카프카로 전송된다.
이벤트 로그는 log_type으로 REQ,RES 두 가지가 있으며 차이점은 RES에 response 항목이 추가된 형태이다.

```json
{
  "log_type": "REQ",
  "trace_id": "40df93c5f8efbef5191070c7ebb38053",
  "span_id": "0",
  "service": "SP016739",
  "operation": "/api/iqokLnL/Lizeth",
  "caller": {
    "channel": "https://joxrotd.edu/XFoWLdG.js",
    "channelIp": "50.89.16.1"
  },
  "host": {
    "name": "https://tekdlij.biz/wdNXOVM.jpg",
    "ip": "70.24.229.49"
  },
  "destination": {
    "name": "http://adantkl.info/ZSshvBx.js",
    "ip": "75.250.187.116"
  },
  "user": {
    "id": "Everette",
    "ip": "165.193.221.2",
    "agent": "ios"
  },
  "event_dt": 1706859987489
}
```

```json
{
  "log_type": "RES",
  "trace_id": "40df93c5f8efbef5191070c7ebb38053",
  "span_id": "0",
  "service": "SP016739",
  "operation": "/api/iqokLnL/Lizeth",
  "caller": {
    "channel": "https://joxrotd.edu/XFoWLdG.js",
    "channelIp": "50.89.16.1"
  },
  "host": {
    "name": "https://tekdlij.biz/wdNXOVM.jpg",
    "ip": "70.24.229.49"
  },
  "destination": {
    "name": "http://adantkl.info/ZSshvBx.js",
    "ip": "75.250.187.116"
  },
  "user": {
    "id": "Everette",
    "ip": "165.193.221.2",
    "agent": "ios"
  },
  "event_dt": 1706859988722,
  "response": {
    "type": "S",
    "status": 200,
    "duration": 2683
  }
}

```