package template

const defaultTemplate = `+++
date = '{{ .createAt }}'
title = '{{ .title }}'
categories = {{ .categories }}
tags = {{ .tags }}
+++`
