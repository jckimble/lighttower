---
---
<html>
  <head>
    <meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<title>{{ site.github.project_title }} :: {{ site.github.project_tagline }}</title>
	<meta name="description" content="{{ site.github.project_title }}: {{ site.github.project_tagline }}">
	<meta name="author" content="James Kimble">
	<meta name="viewport" content="initial-scale=1, maximum-scale=1, user-scalable=no, minimal-ui">

	<link rel="stylesheet" href="http://fonts.googleapis.com/css?family=Roboto:400,500,700%7COpen+Sans:400,300">
	<link rel="stylesheet" href="style.css">
  </head>
  <body>
    <div class="docs-header">
  <header class="docs-masthead clearfix">
  <div class="container">
    <div class="column">
      <h1 class="docs-title">
		  <a href="/" data-ignore="push">{{ site.github.project_title }}</a>
      </h1>
      <nav class="docs-nav clearfix">
        <a class="docs-nav-trigger icon icon-bars pull-right js-docs-nav-trigger" href="#"></a>
        <div class="docs-nav-group">
          <a class="docs-nav-item" href="/">Home</a>
		  <a class="docs-nav-item" href="{{ site.github.issues_url }}">Issues</a>
        </div>
      </nav>
    </div>
  </div>
</header>

  <div class="docs-header-content">
	  <p class="docs-subtitle">{{ site.github.project_tagline }}</p>
	{% assign release = site.github.releases | first %}
	{% for asset in release.assets %}
	{% if asset.name == "lighttower-linux-amd64" %}
	  <a href="{{ asset.browser_download_url }}" class="btn btn-primary">Download {{ site.github.project_title }}</a>
	  <p class="version">Version: {{ release.tag_name }} - Published: {{ release.published_at | date: "%-m-%d-%Y" }}</p>
	  <p class="version size">Size: {{ asset.size | divided_by:1024.0 | divided_by:1024.0 | round:1 }} MB - Downloads: {{ asset.download_count }}</p>
	{% endif %}
	{% endfor %}
  </div>

  <div class="docs-header-bottom">
    <div class="docs-footer">
  <p class="docs-footer-text">Code licensed under the <a href="https://www.gnu.org/licenses/gpl-3.0.txt">GPLv3 License</a>. {{ site.github.project_title }} was created by <a href="https://github.com/jckimble">James Kimble</a>.</p>
</div>
  </div>
</div>
  </body>
</html>
