
var amqp = require('amqp');

var connection = amqp.createConnection({ host: 'amqp://guest:guest@'+process.argv[2]+':5672' });

var publish = function(exchange, callback) {
  console.log("Publishing message")
  exchange.publish("#", {ts: Date.now()}, {mandatory: true}, function(err) {
    console.log('done', err)
  });
  setTimeout(publish.bind(publish, exchange, callback), 100);
};

// Wait for connection to become established.
connection.on('ready', function () {
  console.log("Ready")
  // Use the default 'amq.topic' exchange
  var exc = connection.exchange('amq.topic', {type: 'fanout', durable: true, autoDelete: false}, function (exchange) {
    console.log('Exchange ' + exchange.name + ' is open');
    publish(exchange);
  });
});

connection.on('error', function (err) {
  console.error("Error: " + err.message);
});
