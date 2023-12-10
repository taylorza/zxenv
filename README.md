# zxenv [![zxenv](https://github.com/taylorza/zxenv/actions/workflows/go.yml/badge.svg)](https://github.com/taylorza/zxenv/actions/workflows/go.yml)
Simple command line utility to install a development environment and create starter projects for the ZX Spectrum Next

**NOTE** You do need to manually install the [OpenAL library](https://www.openal.org/downloads/)

## Create a new development environment
The following command will create a development environment in the C:\NextDev folder.
Your development environment can be setup in any folder you like, just run the `zxenv init` command in the target folder.

```
C:\NextDev>zxenv init
```

Without any other arguments a CSpect development environment with a Core 3 compatible 128mb SD Card will be created. To see other options run `zxenv` without any arguments.

## Create a new project
In your newly minted development environment you can run the `zxenv init <projectname>` to create a new project.
You can name the project anything you want, but do not put spaces in the name. For example the following will create a project called `mygame`

```
C:\NextDev>zxenv new mygame
```

The default project is a NEX project, you can also create DOT commands and Drivers using the `-t` or `--type` flag to specify the types of project to create.

### Create a DOT Command project
The command below will create a new project using the template for a DOT command
```
C:\NextDev>zxenv new -tDOT mydot
```

### Create a Driver project
The command below will create a new driver project, the template project is a printer driver that flashes the border color for each character printed.
```
C:\NextDev>zxenv new -tDRV mydriver
```

## Serial Debugging
If you are planing to do serial debugging, the VSCode launch template includes everything you need to debug directly on the Next/N-GO via a serial cable. 
Of course, you will require a custom serial cable that can be plugged into one of the joystick ports. You can find instructions on building a serial cable [HERE](https://amaninhistechnoshed.com/a-man-in-his-technoshed/coding) 

You will also need to update the launch script to let the debugger know what port to use. Locate the following section in the `.vscode/launch.json` file and replace `<SERIAL PORT>` with the port assigned to the serial interface on your machine.

```
"zxnext": 
{
    "serial": "<SERIAL PORT>"
},
```

For Windows that might look something like this (remember your port will probably differ)

```
"zxnext": 
{
    "serial": "COM2"
},
```

On a Mac it would look something like this

```
"zxnext": 
{
    "serial": "/dev/tty.usbserial-AQ02PCVO"
},
```

## What does it do
When you run the `init` command the tool will create the folder structure described below, download the tools and move them to the appropriate location in the folder structure.

|Folder      |Description                                                                  |
|------------|-----------------------------------------------------------------------------|
|.tmp        |Temporary folder used to store the downloaded items and extract the zip files|
|emulator    |The emulator and all the related files like roms etc. can be found in this folder.|
|sdcard      |The SD card image used by the emulator|

Running the `new` command creates a folder with the name of your project that includes a simple sample program as well as the VS Code templates required to launch the emulator, build and debug the application. The default launch action will run the build task before launching the application to ensure that you are running the latest code. 

## What is installed
See [sources.json](https://github.com/taylorza/zxenv/blob/main/sources.json)

sjasmplus builds for Linux and Mac are hosted in this repository and built by me from the original source code without any changes

## What plans do I have for this
This does what I needed for now, but I would like to do add the following support, in no particular order

- [x] Improve the handling of command line arguments so that more options can be speficied
- [x] Support command line options selecting the SD Card image to install
- [x] Move projects into a project folder under the development environment
- [x] Add support for ZEsarUX emulator
- [x] Add support for an external config file that can override the links used to get the resources
- [x] Add support for setting up the environment on a Mac 
- [x] Add support for setting up the environment on Linux
- [x] Cleanup the code so that it makes more intelligent decisions about paths
- [ ] What ever else comes to mind... any ideas?
