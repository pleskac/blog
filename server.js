var mysql      = require('mysql');
var connection = mysql.createConnection({
  host     : 'localhost',
  user     : 'root',
  password : 'rootroot', //not a secret, lol
  database : 'wordpress',
});

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
    	//Valid post. Return it.
		var sql = 'SELECT wp_posts.guid,wp_posts.post_excerpt,wp_posts.id,wp_postmeta.meta_value ';
		sql +=	  'FROM wp_posts ';
		sql +=    'LEFT JOIN wp_postmeta ';
		sql +=    'ON wp_posts.ID = wp_postmeta.post_id ';
		sql +=    'WHERE wp_posts.post_parent = ' + postInt + ' AND wp_posts.post_status = "publish" ';
		sql +=    'AND wp_posts.post_type = "attachment" '
		sql +=    'AND wp_postmeta.meta_key = "_wp_attachment_metadata"';

		var query = connection.query(sql.toString(), function(err, results){
			if(!err){
				response.setHeader("Access-Control-Allow-Origin", "http://pleskac.org");
				response.writeHeader(200, {"Content-Type": "text/plain"});
				response.end(JSON.stringify(results));
			}
			else{
				response.status(404).send('Not found');
			}
		});
    } else if(my_path = "/blog"){
    	//Title page
    	var sql = 'SELECT id,post_title FROM wp_posts WHERE post_type = "post"';
    	var query = connection.query(sql.toString(), function(err, results){
			if(!err){
				response.setHeader("Access-Control-Allow-Origin", "http://pleskac.org");
				response.writeHeader(200, {"Content-Type": "text/plain"});
				response.end(JSON.stringify(results));
			}
			else{
				response.status(500).send('Error getting list of posts');
			}
		});

	} else{
    	//Return error 
		response.status(500).send('Failed to parse post id.');
		response.end();
    }  
}).listen(1337);  
sys.puts("Server Running on 1337");   

