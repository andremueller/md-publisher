# md-publisher

The `md-publisher` golang tool shall simplify the process of publishing local HTML (and
using pandoc Markdown) files to https://medium.com.

Local images are first uploaded to medium and then the article is published as
draft.

Install the `md-publisher` tool with
```bash
go get github.com/andremueller/md-publisher
```

For using the `md-publisher` tool a integration token is required. Hereto you need an medium.com account and must have published an article. Then you can simply create a token on your settings page https://medium.com/me/settings.
Create a TOML configuration file in `$HOME/.config/md-publisher/md-publisher.conf` with the following content:

```TOML
# md-publisher.conf configuration file shall be found in
# $HOME/.config/md-publisher/md-publisher.conf

[medium]
# Settings for medium.com

# Create an integration token in your medium.com account on your settings
# page https://medium.com/me/settings and enter it here
MediumAccessToken="YOUR_ACCESS_TOKEN"
```

After that you should be able to upload a local HTML file with

```bash
md-publisher publish my_file.html
```

Then you should find the newly article within your drafts including local images.

# Dependencies

| Dependency                          | License                    |
| ----------------------------------- | -------------------------- |
| github.com/urfave/cli/v2            | MIT License                |
| github.com/sirupsen/logrus          | MIT License                |
| github.com/Medium/medium-sdk-go     | Apache License Version 2.0 |
| github.com/PuerkitoBio/goquery      | BSD 3 Clause License       |
| github.com/yuin/goldmark            | MIT License                |
| github.com/litao91/goldmark-mathjax | MIT License                |
|                                     |                            |
|                                     |                            |
|                                     |                            |

# License

MIT License

**The project is currently in an experimental state.**
So please don't blame me if something is not working. However, you are welcome to contribute to this project.