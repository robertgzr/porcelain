porcelain
============

Parses `git status --porcelain=v2 --branch` and outputs nicely formatted strings for your shell.

<img width="646" alt="screen_shot" src="https://user-images.githubusercontent.com/3930615/27802035-9c1b92d2-6021-11e7-9289-7b8a17164bf4.png">

The minimum git version for porcelain v2 with `--branch` is `v2.13.2`.
Otherwise you can use the old porcelain v1 based parser on the [`legacy` branch](https://github.com/robertgzr/porcelain/tree/legacy).

With a working Go environment do: `go get -u github.com/robertgzr/porcelain`

Binaries can be found [here](https://github.com/robertgzr/porcelain/releases).

## Output explained:

 ` <branch>@<commit> [↑/↓ <ahead/behind count>][untracked][unmerged][modified][dirty/clean]`

- `?`  : untracked files
- `‼`  : unmerged : merge in process
- `Δ`  : modified : unstaged changes

Definitions taken from: https://www.kernel.org/pub/software/scm/git/docs/gitglossary.html#def_dirty
- `✘`  : dirty : working tree contains uncommited but staged changes
- `✔`  : clean : working tree corresponds to the revision referenced by HEAD

### Notice
In some fonts individual characters may look different, so it is important to use a `powerline` font. They are avalable:

- [From Source](https://github.com/powerline/fonts)
- Via apt package: `apt install fonts-powerline`
- Via dnf package: `dnf install powerline-fonts`

## Usage

Run `porcelain` without any options to get the colorful output :)
For all supported options see `porcelain -h`.

I run this in ZSH to fill my `RPROMPT`, for this the terminal color codes need to be escaped.
Use the `-bash` and `-zsh` flags to do that.

To use it in your tmux statusline you can turn off colors with `no-colors` or switch to tmux color formatting `-tmux`.

If you're using [`tpm`](https://github.com/tmux-plugins/tpm) you can install it as a plugin:
```tmux
set -g @plugin 'robertgzr/porcelain'
```
And then add `#{porcelain}` to your statusline configuration.

This installs the latest version into the tpm plugin directory.

---

The screenshots use:
* [Solarized Dark](http://ethanschoonover.com/solarized) colorscheme
* [Iosevka](https://github.com/be5invis/Iosevka) font
