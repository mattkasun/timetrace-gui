# :alarm_clock: TimeTrace GUI 
## Motivation
Timetrace (https://github.com/dominikbraun/timetrace) is an excellent cli tool for keeping track of time spent on various projects. However, if you are working on a restricted corporate network that does not permit the installation of software, it is difficult to use, especially if said network also restricts outgoing connections so that you cannot ssh to a machine which has timetrace installed.

One option would be to run timetrace in termux (https://github.com/termux/termux-app) on your phone but I decided to build a web front end to timetrace that can be accessed from a web browser on your phone or desktop.

You can use timetrace-gui and timetrace at the same time if you use the native binary. Commands entered in one tool will show up in the other. 

- Language 
  - go with html/templates
- CSS 
  - w3.schools.com
- Material Icons
  - fonts.google.com/icons
- Stop Watch Icon
  - freepngimg.com/

## Installation

### Docker
- GUI only:  
  - `docker run --name timetrace-gui -p 8090:8090 -v $HOME/timetrace/data:/root/.timetrace nusak/timetrace-gui:latest`
  - GUI available at localhost:8090 

- GUI and timetrace:
  - `cd compose; docker-compose up -d`
  - GUI available at localhost:8090 
  - timetrace commands `docker exec timetrace timetrace <timetrace command and flags>` 
    - `docker exec timetrace timetrace version` or `docker exec timetrace timetrace status`


### Native binary
Download the appropriate binary from releases
- x86 
- i386 
- windows 
- mac 
- arm64 (pi4) 

Extract it to a directory in your path.  Run timetrace-gui and point your browser at localhost:8090.



## Usage
Run the timetrace binary (or docker container) and point your browser at localhost:8090.

On first run a user creation dialog will be presented.  Use the user/password you enter for future logins.

Port forwarding or setting up a reverse proxy to enable access from the intenet is left as an exercise for the reader.

## Screenshots
### Web
![browser](https://github.com/mattkasun/timetrace-gui/raw/master/screenshots/web.png "TimeTrace-GUI with Browser")

### Mobile
![phone](https://github.com/mattkasun/timetrace-gui/raw/master/screenshots/mobile.png "TimeTrace-GUI with Phone")

## RoadMap

- [x] Build binaries for all architectues
- [ ] Restore deleted projects
- [ ] Edit project
- [ ] Edit record
- [x] Reports
- [ ] Edit Configuration
- [ ] Users

