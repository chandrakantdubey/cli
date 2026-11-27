## Try it in the terminal

### Turn a captured package answer into evidence

Inspect the header and first two rows of each capture. In the Debian file, compare the name and version for similar packages. A version includes a distribution revision, so keep the entire field. In Fedora's file, compare the two paths beginning with `/usr/bin/runner`. The executable's owner comes from the exact path.

Try a deliberately broad search in each capture, then replace it with an equality test on field 1. Only one row should remain in each case. To compare with your own system, `command -v rpm` and `command -v dpkg-query` tell you which query tools are available. The scored lab uses the captures and does not need either tool or an installation.
