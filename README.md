# krita notes helper

`knot`, standing for **k**rita **not**es, is a simple command line tool for organizing handwritten digital notes using krita. It views projects as containing "batches" of pages:

* project
    * batch-0
        - page-0.kra
        - page-1.kra
        - ...
    * batch-1
        - page-0.kra
        - ...
    * ...

knot allows you to create and quickly open these projetcs, add batches, add pages to batches and export batches to pdf with simple commands. It has been compiled and tested for linux, and requires the following dependencies:
```
krita
imagemagick
ghostscript
```

## Installation
This repository includes a compiled binary which may work for you. You can clone it wherever you like:
```sh
$ git clone https://github.com/RayOfSunDull/knot
```
You may use the install script to install it:
```
$ cd path_to_repo
$ chmod +x ./install.sh
$ ./install.sh
```
If you don't want to run the script, what it does is basically this: 
* move the executable from `path_to_repo/bin` to `~/bin`
* move the template directory from `path_to_repo/templates` to `~/.config/knot`
* move `path_to_repo/projects.json` to `~/.config/knot` if it exists, otherwise create it

You may also compile it using the `go` compiler:
```
$ cd path_to_repo
$ ./build.sh
```
There are no build dependencies other than the go standard library.

## Basic Usage
### silent mode

By default, many knot commands will open created files. If you don't want this, add this flag:
```sh
$ knot -s [commands]
```
### Initialising and accessing a project
```sh
$ knot -i project_name
```
will create a project called `project_name` in `$PWD/project_name`. It will get registered to knot:
```sh
$ knot -l
registered projects:
        project <project_name> in </path/project_name>
```
And you can open it from any shell using:
```sh
$ knot -o project_name
```
And also deregister it:
```sh
$ knot -d project_name
```
By default, it will open `nautilus` on the project directory and all the `.kra` files in the last batch. There is currently no way to change the file explorer short of recompiling the tool.

### Managing batches
Once your working directory is **inside** the project directory somewhere, you may run:
```sh
$ knot -b
```
to add a new batch of notes. It will create one with the next appropriate number. If you want to specify a batch number, use:
```sh
$ knot -sb batch_number
```
You may open all the `.kra` pages in a batch using:
```sh
$ knot -ob batch_number
```
You may add a page to the `latest` batch using:
```sh
$ knot -p
```
You may add a page to a specified batch using:
```sh
$ knot -sp batch_number
```
You may export the pages in the `latest` batch to pdf using:
```sh
$ knot -e
```
It will be created at the top level of the batch and named `batch_name.pdf`. You may export the pages in a specified batch to pdf using:
```sh
$ knot -se batch_number
```
There are a few more configurable options, such as batch names and the ability to generate batches in a subdirectory instead of the top level of the project. Please refer to `knot -h` for info on all commands.
