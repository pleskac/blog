$(window).load(function() {

});

$(document).ready(function () {
	var pathname = window.location.pathname;
	var url = getURL(pathname)

	jQuery.ajax(url).done(
	        function(data){
		post = JSON.parse(data);
		// First append the text
		$('#post_content').append(post.Content);

		// Next, append the photos
		jQuery.each(post.Pictures, function(){
			var imageURL = resizeToLarge(this.Guid, this.Meta_value);

			var imageString = '<img class="caption" src="' + imageURL + '" data-caption="' + this.Post_excerpt + '" />';

			$('#post_content').append(imageString);
		});

		// we need to wait until the <img> are inserted
		$('img.caption').captionjs({
        	'class_name' : 'captionjs',  // Class name assigned to each <figure>
        	'schema'     : true,         // Use schema.org markup (i.e., itemtype, itemprop)
        	'mode'       : 'animated',    // default | static | animated | hide
        	'debug_mode' : true         // Output debug info to the JS console
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
		var imageName = (jibberish.toString().substring(start,end - 1));
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
