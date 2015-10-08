
var amqp = require('amqp');

var connection = amqp.createConnection({ host: 'amqp://guest:guest@localhost:5672' });

// Wait for connection to become established.
connection.on('ready', function () {
  // Use the default 'amq.topic' exchange
  console.log("Ready")
  connection.queue('test123', {durable: true, autoDelete: false}, function (q) {
      // Catch all messages
      q.bind('#');

      // Receive messages
      q.subscribe(function (message, headers, deliveryInfo, messageObject) {
        // Print messages to stdout
        console.log(message);
        //messageObject.acknowledge();
      });
  });
});

connection.on('error', function (err) {
  console.error("Error: " + err.message)
});
