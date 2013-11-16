/*var http = require('http');
http.createServer(function (req, res) {
  res.writeHead(200, {'Content-Type': 'text/plain'});
  res.end('Hello World\n');
}).listen(1337, "127.0.0.1");
console.log('Server running at http://127.0.0.1:1337/');

*/


var sys = require("sys"),  
my_http = require("http");  
path = require("path"),  
url = require("url"),  
filesys = require("fs"); 
my_http.createServer(function(request,response){  
    var my_path = url.parse(request.url).pathname;  
    var full_path = path.join(process.cwd(),my_path);  
    var postInt = parseInt(my_path.substring(1));
    if(postInt > 0){
        sys.puts(postInt);
    }
    //sys.puts(full_path);  
    //select the files based on my_path
    response.writeHeader(200, {"Content-Type": "text/plain"});  
    response.write("Hello World");  
    response.end();  
}).listen(1337);  
sys.puts("Server Running on 1337");   
