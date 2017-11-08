fileupload
========

[Go](http://golang.org) package for uploading one or more files.

The package can sort files into different directories by mimetype.

Special functions for uploading images, that save the images as jpegs and create thumbnail images.

External compiled dependency libvips

## Notes on installation

Tested on Arch linux

Install libvips/vips, a fast image manipulation program. Instructions on github page.

https://github.com/jcupitt/libvips

```bash
$ ./configure
$ make
$ sudo make install
$ sudo ldconfig
```

By default this will install files to "/usr/local"

---------------

You will need to edit an environmental variable (like $PATH, on linux usually stored in "~/.bash_profile")

1. Check if the variable exists
```bash
echo $PKG_CONFIG_PATH
```

2. If blank
```bash
export PKG_CONFIG_PATH="/usr/local/lib/pkgconfig"
```
Else (not blank)
```bash
export PKG_CONFIG_PATH="$PKG_CONFIG_PATH:/usr/local/lib/pkgconfig"
```

3. Update Bash
```bash
source ~/.bash_profile
```


The linker (GNU 'ld') might have trouble locating libvips.so....
As root, add a configuration file
```bash
# cd /etc/ld.so.conf.d
# touch libvips.conf (verify that it has the same permissions as the other files)
# nano libvips.conf
```
(add one line) "/usr/local/lib/"
```bash
# ldconfig (refreshes the linker cache)
```

-----------------------

Last, install the go package for interfacing with libvips:

go get -u gopkg.in/h2non/bimg.v1
