$(document).ready(function () {
	var pathname = window.location.pathname;
	var url = getURL(pathname)

	jQuery.ajax(url).done(
	        function(data){
		post = JSON.parse(data);
		jQuery.each(post.Pictures, function(){
			var imageURL = resizeToLarge(this.Guid, this.Meta_value);
			var height = getImageHeight(imageURL);
			var imageString = '<div class="caption" style="height:' + height + 'px">';
			imageString += '<img class="dynamic_image" src="' + imageURL + '" />';
			if(this.Post_excerpt){
				imageString += '<span>' + this.Post_excerpt + '</span>';
			}
			imageString += '</div>';
			$('#post_content').append(imageString);
		});
	});
});


function getURL(url){
	var pattern = new RegExp('([0-9]+)');
	var match = url.toString().match(pattern);
	if(match != null){
		return 'http://pleskac.org:1337/' + match[0];
	}
	return '';
}

function resizeToLarge(uri, jibberish){
	var pattern = new RegExp('((s:21:").+jpg")');
	var match = jibberish.toString().match(pattern);
	if(match != null){
		var start = match.index + 6;
		var end = match.index + match[0].length;
		var imageName = (jibberish.toString().substring(start,end));
		uri = uri.replace(/\/([^\/]*)$/,'/');
		return uri + imageName;
	}
	return uri;
}

function getImageHeight(url){
	var pattern = new RegExp('x(\\d+)');
	var match = url.toString().match(pattern);
	if(match != null){
		return match[1];
	}
	return 23;

}
