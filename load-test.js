import http from 'k6/http';
import { check } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export let options = {
    vus: 10,           // Virtual Users
    duration: '30m',   // Test duration
    rps: 30            // Requests per second
};


function getRandomIntInRange(min, max) {
    return Math.floor(Math.random() * (max - min)) + min;
}

export default function () {

    console.log(getRandomIntInRange(5, 15));
    let randomName = randomString(getRandomIntInRange(5, 200));

    let res = http.get(`http://envoy:8081/app?name=${randomName}`);

    check(res, {
        'status is 200': (r) => r.status === 200,
    });
}