porcelain
============

Parses `git status --porcelain` and outputs nicely formatted strings.

**Formatted output with `porcelain -fmt`**

![formatted output screenshot](http://i.imgur.com/d3Ckvbj.png)

![formatted output screenshot 2](http://i.imgur.com/xAnGH7C.png)

**Basic output with `porcelain -basic`**

![basic output screenshot](http://i.imgur.com/F1DnTOA.png)

```
   commit,branch,tracked_branch,ahead,behind,untracked,added,modified,deleted,renamed,copied
```

With a working Go environment do: `go get -u github.com/robertgzr/porcelain`

Binaries can be found [here](https://github.com/robertgzr/porcelain/releases).

---

The screenshots use:
* [Solarized Dark](http://ethanschoonover.com/solarized) colorscheme
* [Iosevka](https://github.com/be5invis/Iosevka) font
