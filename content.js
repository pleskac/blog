$(document).ready(function () {
	var pathname = window.location.pathname;

	if(pathname = "/blog/"){
		jQuery.ajax("http://pleskac.org:1337/blog").done(
	                function(data){
				posts = JSON.parse(data).reverse();
				jQuery.each(posts, function(){
					var link = '<a href="http://pleskac.org:1337/blog/' + this.Id + '" >' + this.Post_title + '</a><br />';
					$('#post_content').append(link);
				});
		});
	}
});