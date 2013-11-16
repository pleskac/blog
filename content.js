$(document).ready(function () {
	jQuery.ajax("http://pleskac.org:1337/5439").done(
                function(data){
			if (window.console) console.log(data);
			image_urls = JSON.parse(data);
			jQuery.each(image_urls, function(){
				var imageString = '<div class="caption">';
				imageString += '<img class="dynamic_image" src="' + this.guid + '" />';
				imageString += '<span>' + this.post_excerpt + '</span>';
				imageString += '</div>';
				//$('#post_content').append('<div class="caption">'); 
				//$('#post_content').append('<img class="dynamic_image" src="' + this.guid + '" />');
				//$('#post_content').append('<span>' + this.post_excerpt + '</span>');
                        	if (window.console) console.log(imageString);
				$('#post_content').append(imageString);
			});
	});
});

function resizeToLarge(uri){
	uri = uri.replace('.jpg','');
	uri.append('-1024x682.jpg');
	return uri;
}
