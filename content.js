$(document).ready(function () {
	jQuery.ajax("http://pleskac.org:1337/5439").done(
                function(data){
			image_urls = JSON.parse(data);
			jQuery.each(image_urls, function(){
				var imageURL = resizeToLarge(this.guid, this.meta_value);
				
				var imageString = '<div class="caption">';
				imageString += '<img class="dynamic_image" src="' + imageURL + '" />';
				if(this.post_excerpt){
					imageString += '<span>' + this.post_excerpt + '</span>';
				}
				imageString += '</div>';
				$('#post_content').append(imageString);
			});
	});
});

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
