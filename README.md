# fungjai

## Local Deploy Backend

### ขั้นเตรียม

1. สร้าง python environment โดยรันคำสั่ง

``` sh
    python3 -m venv ENV
```

2. activate python environment

``` sh
    source ENV/bin/activate
```

3. ทำการติดตั้ง packages โดยรันคำสั่ง

``` sh
    pip install -r requirements.txt
```

4. ติดตั้ง docker compose สามารถศึกษาขั้นตอนได้จาก ที่นี่ >> https://docs.docker.com/get-docker/

### Deploy

1. run command

``` sh
    docker-compose up -d
```

**หลังรันคำสั่ง จะได้ 3 containers** 

- clickhouse (sql database)

- minio (bucket storage)

- metabase (visualize tool)

### วิธีเช็คว่าทุกอย่างทำงานได้ถูกต้อง

- สามารถเปิด minio ui ได้ที่

```
http://localhost:9001/
```

- สามารถเปิด clickhouse ui ได้ที่

```
http://localhost:8123/play
```

- สามารถเปิด metabase ui ได้ที่

```
http://localhost:3000/setup
```
### Docker Data
เพื่อความสะดวกในการเข้าถึง table ใน clickhouse และการเข้าถึง file ใน minio bucket ได้ทำการ mount volume ของ docker container ไว้ที่ local ดังด้านล่าง

``` yml
      volumes:
        - ./docker_data/clickhouse/server:/var/lib/clickhouse
        
```

ซึ่งจะอยู่ภายใน folder docker_data/ ดังด้านล่าง

![docker_data](backend/imgs/docker_data_folder.png)

หากต้องการเริ่ม docker container โดยละทิ้งค่าที่เคยตั้งค่าไว้ทั้งหมด ให้ลบ
- docker_data/clickhouse
- docker_data/minio/.minio.sys

## DBT Deploy
1. เข้าไปใน folder clickhose_dbt/

1. run command เพื่อเช็คว่าสามารถ connect ไปยัง clickhouse ได้

``` sh
   dbt debug
```

3. run command เพื่อทำงานตามการแก้ไขใหม่

``` sh
   dbt run
```