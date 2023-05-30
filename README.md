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

knot allows you to create and quickly open these projetcs, add batches, add pages to batches and export batches to pdf with simple commands. It has only been compiled and tested for linux, and only has ``krita`` as a runtime dependency.

## Installation (Linux)
This repository includes a compiled binary which may work for you. You can clone it wherever you like:
```sh
$ git clone https://github.com/RayOfSunDull/knot
```
You may use the build script to install it. You need to have ``go`` installed in your system to run it.
```sh
$ cd path_to_repo
$ go run build.go -no-build -install
```
This command will install the binaries and config files without building them from source. Run ``go run build.go -help`` for more options.

You may also build it from source, but since Python was added to the project this has become a bit more complicated. You need to have ``cython`` to convert the export script to a C file and ``gcc`` to compile it to a binary. First we need to set up the build area and necessary Python modules:
```sh
$ cd path_to_repo
$ mkdir aux # build directory for the C file
$ python -m venv export-venv # make venv for the modules
$ source export-venv/bin/activate # activate the venv
$ pip install -r requirements.txt # install modules
```
Depending on your distribution, you might need to use ``python3`` instead. Just make sure ``python --version`` gives you Python 3. Now everything is ready to run the build script:
```sh
$ go run build.go -install # this will build and install
```
On the go side, there are no packages required beyond the standard library.

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
    "ExportQuality": 100,
}
```
The string passed to each setting must be the name of the **command line utility** that opens the appropriate program. Currently it's not possible to configure your commands to accept extra options for the viewers. You may lower ``ExportQuality`` to save space, and this is recommended. Generaly the readability won't drop too much even if you set ``ExportQuality`` to 10 (implied 10%). Play around with it and find what best suits your needs.

### Roadmap (tentative)
* Port to Windows and maybe MacOS
* Switch to custom arg parser