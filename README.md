# redmine-org-feed-proxy

A proxy for https://redmine.org/ to subscribe feeds from Slack.

## Usage

```
$ PORT=4567 go run . &
$ curl http://localhost:4567/projects/redmine/repository/revisions.atom
```

## Why proxy?

`/feed subscribe http://redmine.org/projects/redmine/repository/revisions.atom` fails with an error of `That feed (http://redmine.org/projects/redmine/repository/revisions.atom) could not be saved. Try again?`.

(I haven't check the cause.)

## License

CC0
