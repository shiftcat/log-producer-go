
# Kafka 설치

- .bashrc 파일에 CONFLUENT_HOME 환경 변수로 Confluent 디렉토리 설정

```sql
export CONFLUENT_HOME=/app/kafka/confluent-7.1.10
export PATH=.:$PATH:$CONFLUENT_HOME/bin
# 수정된 bashrc 적용
. .bashrc
```


## Kafka 클러스터 ID 생성 및 포맷

```shell
./bin/kafka-storage random-uuid

rYyh1JowTzK61eOCZBp3qg

./bin/kafka-storage format -t rYyh1JowTzK61eOCZBp3qg -c ./etc/kafka/kraft/server.properties
```

## Kafka 시작/중지

```shell
./bin/kafka-server-start -daemon ./etc/kafka/kraft/server.properties

./bin/kafka-server-start ./etc/kafka/kraft/server.properties

./bin/kafka-server-stop -daemon ./etc/kafka/kraft/server.properties
```

## systemd 서비스 이용
systemd 서비스로 사용하기 위해서는 다음과 같이 service 파일 작성을 해야 합니다.
파일경로 및 파일명 : /etc/systemd/system/kafka.service

```shell
[Unit]
Description=Apache Kafka server (broker)

[Service]
Type=forking
User=root
Group=root
Environment='KAFKA_HEAP_OPTS=-Xms512m -Xmx512m'
ExecStart=/usr/local/kafka/bin/kafka-server-start.sh -daemon /usr/local/kafka/config/kraft/server.properties
ExecStop=/usr/local/kafka/bin/kafka-server-stop.sh
LimitNOFILE=16384:163840
Restart=on-abnormal

[Install]
WantedBy=multi-user.target
```

위와 같이 systemd 서비스 파일 작성이 완료 되었다면 처음에는 enable 설정을 진행한 다음 서비스를 시작 합니다.

```shell
## 서비스 enable 설정(1회 실행)
sudo systemctl daemon-reload
sudo systemctl enable kafka


## 시작
sudo systemctl start kafka

## 종료
sudo systemctl stop kafka

```

# Kafka console consumer

## 토픽 생성/삭제

```shell
# 토픽 목록 조회
./bin/kafka-topics \
--bootstrap-server localhost:9092 \
--list

# 토픽 상세보기
./bin/kafka-topics \
--bootstrap-server localhost:9092 \
--topic EVENT_LOG \
--describe 

# 토픽 생성
./bin/kafka-topics \
--bootstrap-server localhost:9092 \
--topic EVENT_LOG \
--partitions 3 \
--create 

./bin/kafka-topics \
--bootstrap-server localhost:9092 \
--topic TOPIC_LOGS \
--partitions 3 \
--create 

./bin/kafka-topics \
--bootstrap-server localhost:9092 \
--topic TOPIC_STAT \
--partitions 3 \
--create 

# 토픽 삭제
./bin/kafka-topics \
--bootstrap-server localhost:9092 \
--topic EVENT_LOG \
--delete 

```

## 콘솔 컨슈머

```shell
./bin/kafka-console-consumer \
--bootstrap-server localhost:9092 \
--topic EVENT_LOG \
--from-beginning \
--property print.key=true

./bin/kafka-console-consumer \
--bootstrap-server localhost:9092 \
--group myGroup \
--topic myTopic \
--from-beginning
```


### Schema Registry 기동

- schema registry 기동 스크립트 생성

```sql
vi registry_start.sh
$CONFLUENT_HOME/bin/schema-registry-start $CONFLUENT_HOME/etc/schema-registry/schema-registry.properties
chmod +x *.sh
```

### ksqlDB 기동

- [ksql-server.properties](http://ksql-server.properties) 환결 설정 파일 수정

```shell
./bin/ksql-server-start ./etc/ksqldb/ksql-server.properties
```

```sql
bootstrap.servers=localhost:9092
ksql.schema.registry.url=http://localhost:8081
```