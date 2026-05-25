# sp — social poster

Posts newly-published blog posts to social media. A small CLI that, for each
new markdown post, parses the frontmatter, skips drafts and posts without a
description, formats a per-platform message (description + hashtags + URL,
truncated to each platform's limit), and posts to **Bluesky**, **Mastodon**,
and **X (Twitter)**.

Each network is attempted only when its credentials are present in the
environment; a failure on one network never aborts the others.

## Install

```sh
go install github.com/jbrodriguez/sp@latest
```

## Usage

```sh
sp "data/posts/202545/index.md"
# or several files (newline- or space-separated, as one or many args):
sp "data/posts/202545/index.md
data/posts/202546/index.md"
```

Post paths are resolved relative to the current working directory, so run `sp`
from your site repo root (cover images referenced in frontmatter are resolved
relative to the post file).

## Configuration

| Variable | Default | Purpose |
| --- | --- | --- |
| `SITE_URL` | `https://jbrio.net` | Base URL used to build the post link (`<SITE_URL>/posts/<slug>/`). |
| `DEPLOYMENT_DELAY_MS` | `60000` | Delay before posting, to let a deploy finish propagating. |

Credentials (each network is skipped if its vars are unset):

- **Bluesky** — `BLUESKY_IDENTIFIER`, `BLUESKY_PASSWORD`
- **Mastodon** — `MASTODON_INSTANCE_URL`, `MASTODON_ACCESS_TOKEN`
- **X / Twitter** — `TWITTER_API_KEY`, `TWITTER_API_SECRET`, `TWITTER_ACCESS_TOKEN`, `TWITTER_ACCESS_SECRET`

## A post is announced only when its frontmatter has

- `status: published`
- a non-empty `description`
