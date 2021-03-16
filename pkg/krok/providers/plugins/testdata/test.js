var fs = require("fs");
var stdinBuffer = fs.readFileSync(0); // STDIN_FILENO = 0
console.log(stdinBuffer.toString());

console.log('Payload received. Showing arguments...');
var myArgs = process.argv.slice(2);
console.log('myArgs: ', myArgs);
