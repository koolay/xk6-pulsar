
// 1. init code
import {check} from 'k6';
import pulsar from 'k6/x/pulsar';

const data = `{"c":"zzzz","d":221,"er":false}`
const client = pulsar.createClient({url: `pulsar://${__ENV.PULSAR_ADDR}`})
const producer = pulsar.createProducer(client, {topic: __ENV.PULSAR_TOPIC})

export default function() {
  // 3. VU code
  let err = pulsar.publish(producer, data, {}, false);
  check(err, {
    "is send": err => err == null
  })
}

export function setup() {
  // 2. setup code

}

export function teardown(data) {
  // 4. teardown code
  pulsar.closeClient(client)
  pulsar.closeProducer(producer)
  console.log("teardown!!")
}
