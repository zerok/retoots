# Getting started with retoots

retoots for now consists of a server program that acts as a go-between your website (or your users) and your Mastodon account. In order to get started with it, you should get a retoots package from our [releases][r] page for your platform.

[r]: https://github.com/zerok/retoots/releases

Once you have downloaded a matching archive and extracted it, run the `retoots` binary with the following arguments to specify on which port you want the server to listen on and what Mastodon accounts should be proxied by it:

```
$ ./retoots --addr localhost:8000 \
    --allowed-root-accounts "zerok@chaos.social"
```

retoots will now start on port 8000 and you can, for instance, request all the
interactions with <https://chaos.social/@zerok/105600404041298688> by doing a
GET request to the `/api/v1/interactions` endpoint:


```
$ curl "http://localhost:8000/api/v1/interactions?status=https://chaos.social/@zerok/105600404041298688" | jq .
{
  "descendants": [
    {
      "id": "105604551282351868",
      "content": "<p><span class=\"h-card\"><a href=\"https://chaos.social/@zerok\" class=\"u-url mention\" rel=\"nofollow noopener noreferrer\" target=\"_blank\">@<span>zerok</span></a></span> As I haven't been able to attend a Golang meetup for quite some time, I'm looking forward to attending remotely next time!</p>",
      "created_at": "2021-01-23T10:26:11Z",
      "url": "https://fosstodon.org/@totoroot/105604551260393403",
      "account": {
        "id": "221659",
        "username": "totoroot",
        "acct": "totoroot@fosstodon.org",
        "avatar": "https://chaos.social/system/accounts/avatars/000/221/659/original/30d2525840b524ab.png",
        "url": "https://fosstodon.org/@totoroot"
      }
    }
  ],
  "favorites_by": [
    {
      "id": "221659",
      "username": "totoroot",
      "acct": "totoroot@fosstodon.org",
      "avatar": "https://chaos.social/system/accounts/avatars/000/221/659/original/30d2525840b524ab.png",
      "url": "https://fosstodon.org/@totoroot"
    },
    {
      "id": "173779",
      "username": "mwfc",
      "acct": "mwfc@chaos.social",
      "avatar": "https://chaos.social/system/accounts/avatars/000/173/779/original/10787cdaf4242dcb.png",
      "url": "https://chaos.social/@mwfc"
    }
  ]
}
```

If you want to use this server via XHR/Fetch-API from within a web browser, you
can also specify allowed origin hosts (CORS) using the `--allowed-origins`
flag.

