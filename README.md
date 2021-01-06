# ReToots!

A lightweight service for fetching comments to a particular Mastdodon status
and showing them as comments on, for instance, a static website.

While one could directly fetch the mentions via the Mastodon API within the
user's browser, I wanted to have some more control over what comments are shown
on my website and therefore have a layer between the blog and the upstream
Mastodon server.

Right now, retoots is just a proxy with some small normalizations for easier
consumption but I plan to add a couple features eventually:

- Caching
- Blocking of comments from specific accounts

This repository also doesn't include any client implementation. If you want to
see an example, you can find it in the source code of my blog
[here](https://github.com/zerok/zerokspot.com/blob/a1ed110b87c750b7f9934039653c12c529887a1c/assets/js/main.js#L84).
