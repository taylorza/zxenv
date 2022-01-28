# zxenv [![zxenv](https://github.com/taylorza/zxenv/actions/workflows/go.yml/badge.svg)](https://github.com/taylorza/zxenv/actions/workflows/go.yml)
Simple command line utility to install a development environment and create starter projects for the ZX Spectrum Next

**NOTE** You do need to manually install the [OpenAL library](https://www.openal.org/downloads/)

## Create a new development environment
The following command will create a development environment in the C:\NextDev folder.
Your development environment can be setup in any folder you like, just run the `zxenv init` command in the target folder.

```
C:\NextDev>zxenv init
```

Without any other arguments a CSpect development environment with a 512MB SD Card will be created. To see other options run `zxenv` without any arguments.

## Create a new project
In your newly minted development environment you can run the `zxenv init <projectname>` to create a new project.
You can name the project anything you want, but do not put spaces in the name. For example the following will create a project called `mygame`

```
C:\NextDev>zxenv new mygame
```

## What does it do
When you run the `init` command the tool will create the folder structure described below, download the tools and move them to the appropriate location in the folder structure.

|Folder      |Description                                                                  |
|------------|-----------------------------------------------------------------------------|
|.tmp        |Temporary folder used to store the downloaded items and extract the zip files|
|emulator    |The emulator and all the related files like roms etc. can be found in this folder.|
|sdcard      |The only file in here is the sdcard image, currently this only installs the 512MB image|

Running the `new` command creates a folder with the name of your project that includes a simple sample program as well as the VS Code templates required to launch the emulator, build and debug the application. The default launch action will run the build task before launching the application to ensure that you are running the latest code. 

## What is installed
|Tool        |URL|
|------------|-----------------------------------------------------------------------------|
|CSpect      |http://www.javalemmings.com/public/zxnext/CSpect2_15_01.zip|
|DeZogPlugin |https://github.com/maziac/DeZogPlugin/releases/download/v2.1.0/DeZogPlugin.dll|
|dezog-conf  |https://raw.githubusercontent.com/maziac/DeZogPlugin/main/DeZogPlugin.dll.config|
|sjasmplus   |https://github.com/z00m128/sjasmplus/releases/download/v1.18.3/sjasmplus-1.18.3.win.zip|
|hdfmonkey   |http://uto.speccy.org/downloads/hdfmonkey_windows.zip|
|sdcard      |http://www.zxspectrumnext.online/cspect/tbbluemmc-512mb.zip|

## What plans do I have for this
This does what I needed for now, but I would like to do add the following support, in no particular order

- [x] Improve the handling of command line arguments so that more options can be speficied
- [x] Support command line options selecting the SD Card image to install
- [ ] Move projects into a project folder under the development environment
- [x] Add support for ZEsarUX emulator
- [ ] Add support for an external config file that can override the links used to get the resources
- [ ] Add support for setting up the environment on a Mac 
- [ ] Add support for setting up the environment on Linux
- [ ] Cleanup the code so that it makes more intelligent decisions about paths
- [ ] What ever else comes to mind... any ideas?
