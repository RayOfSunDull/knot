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

knot allows you to create and quickly open these projetcs, add batches, add pages to batches and export batches to pdf with simple commands. It has been compiled and tested for linux, and has the following runtime dependencies:
```
krita
nohup
```
In the future I intend to drop ``nohup`` as a dependency, but ``krita`` will be necessary, of course.

## Installation (Linux)
This repository includes a compiled binary which may work for you. You can clone it wherever you like:
```sh
$ git clone https://github.com/RayOfSunDull/knot
```
You may use the install script to install it:
```sh
$ cd path_to_repo
$ make install
```
If you don't want to run the script, what it does is basically this: 
* move the executable from `path_to_repo/bin` to `~/bin`
* move the template directory from `path_to_repo/templates` to `~/.config/knot`
* move `path_to_repo/projects.json` to `~/.config/knot` if it exists, otherwise create it

You may also compile it using the `go` compiler:
```sh
$ cd path_to_repo
$ go get github.com/signintech/gopdf
$ make full # this will build and install
```
The ``gopdf`` package is now a build dependency and thus must be installed. 

## Basic Usage
### Silent mode

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
By default, it will open `nautilus` on the project directory and all the `.kra` files in the last batch. The file viewer, as well as the pdf viewer can be changed. More info on that later.

### Managing batches
Once you've initialised or opened an already existing project, the temporary knot working directory is set to the directory of that project. You may check this using:
```sh
$ knot -pwd
```
You may also set the knot working directory manually
```sh
$ knot -wd working_dir
```
If the knot working directory is inside a project, you may run commands affecting it, for example:
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

### Configuring knot
You can configure the file explorer and pdf reader used by knot. The config directory can be accessed as such:
```sh
$ ls ~/.config/knot
config.json  projects.json  templates
```
If there is no ``config.json`` file, you should create it:
```sh
$ cd ~/.config/knot
$ touch config.json
```
Open config.json with your preferred editor and paste the settings:
```json
{
    "PDFReader": "your-preferred-pdf-reader",
    "FileExplorer": "your-preferred-file-explorer",
    "ExportCompression": 0,
    "LegacyExport": false
}
```
The string passed to each setting must be the name of the **command line utility** that opens the appropriate program. Currently it's not possible to configure your commands to accept extra options for the viewers. You may set ``ExportCompression`` to 1 or 2 for higher compression, but be weary as this requires ``ghostscript``. The ``LegacyExport`` option requires ``ghostscript`` and ``imagemagick`` and should generally be avoided unless the regular export doesn't work.

### Roadmap
* Allow for more configurable viewer commands
* Look into writing a compression script in Python for portability
* Make versions for Windows and MacOS