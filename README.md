# Test pulsar using k6

## Build k6 with extension 

```bash

# install xk6 using go or download(https://github.com/grafana/xk6/releases)
❯ go install go.k6.io/xk6/cmd/xk6@latest

# build k6
❯ xk6 build --with gitlab.mypaas.com.cn/fast/k6x/pls=.

```

## Run tests

```bash
# --vus virtual users 
# --duration  how lang continue testing

❯ PULSAR_TOPIC=localtest PULSAR_ADDR=localhost:6650 ./k6 run test_producer.js --duration 10s --vus 2

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: test_producer.js
     output: -

  scenarios: (100.00%) 1 scenario, 2 max VUs, 40s max duration (incl. graceful stop):
           * default: 2 looping VUs for 10s (gracefulStop: 30s)

INFO[0010] teardown!!                                    source=console

running (10.1s), 0/2 VUs, 2142 complete and 0 interrupted iterations
default ✓ [======================================] 2 VUs  10s

     ✓ is send

     █ teardown

     checks.........................: 100.00% ✓ 2142       ✗ 0
     data_received..................: 0 B     0 B/s
     data_sent......................: 0 B     0 B/s
     iteration_duration.............: avg=9.31ms min=102.58µs med=9.19ms max=18.47ms p(90)=10.67ms p(95)=11.69ms
     iterations.....................: 2142    212.791617/s
     pulsar.publish.error.count.....: 0       0/s
     pulsar.publish.message.bytes...: 70 MB   7.0 MB/s
     pulsar.publish.message.count...: 2142    212.791617/s
     vus............................: 2       min=2        max=2
     vus_max........................: 2       min=2        max=2
```

