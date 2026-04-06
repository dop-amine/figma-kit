#!/bin/zsh
PS1='%F{blue}❯%f '
h() { print -P "%F{cyan}▸ $1%f  $2" }
ok() { print -P "%F{green}✓%f $1" }
dim() { print -P "%F{8}$1%f" }
title() { print -P "%F{blue}◆ figma-kit%f  AI-powered Figma design" }
hr() { dim '────────────────────────────────────────────────────' }
