# fungjai

Local Deploy Backend

1. แก้ folder vol plugin ใน docker-compose.yml เป็น path เครื่องตัวเอง

``` yml
    volumes:
      - /Users/users/fungjai/backend/plugins:/plugins
```

2. run command

``` sh
    docker-compose up -d
```

จะได้ 3 containers -> clickhouse (sql database), minio (bucket storage) และ metabase (visualize tool)

วิธีเช็คว่าทุกอย่าง up

- เข้า minio ui

```
http://localhost:9001/
```

- เข้า clickhouse ui

```
http://localhost:8123/play
```

- เข้า metabase ui

```
http://localhost:3000/setup
```
