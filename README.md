# TimeTrace GUI
## Motivation
Timetrace (https://github.com/dominikbraun/timetrace) is an excellent cli tool for keeping track of time spent on various projects. However, if you are working on a restricted corporate network that does not permit the installation of software, it is difficult to use, especially if said network also restricts outgoing connections so that you cannot ssh to a machine which has timetrace installed.

One option would be to run timetrace in termux (https://github.com/termux/termux-app) on your phone, but I decided to build a front end to timetrace that can be accessed from a web browser on your phone or desktop.

You can use timetrace-gui and timetrace at the same time if you use the native binary. Updates entered in one tool will show up in the other. 

- Language 
-- go with html/templates
- CSS 
-- w3.schools.com
- Material Icons
-- fonts.google.com/icons
- Stop Watch Icon
-- freepngimg.com/

## Installation
### Docker
~~ docker run --name netmaker-gui -p 8090:8090 -v $HOME/.timetrace:/root/.timetrace nusak/timetrace-gui:v0.0.1 ~~

Any new timetrace projects or records created will be owned by root.  Not an issue if you only use timetrace-gui but it will cause problems with timetrace itself.

Building and running a docker as an arbitrary user is tricky.  Not too difficult if you know the UID:GID ahead of time, but there is no guarentee that the user's UID:GID is 1000:1000.  Not sure how this works on windows or mac.

### Native binary
Download the appropriate binary 
- x86 (https://github.com/mattkasun/timetrace-gui/blob/master/timetrace-gui) 
- i386 - coming soon
- windows - comming soon
- mac - coming soon
- arm64 (pi4) - coming soon
- arm7 (pi3) - coming soon

copy it to a directory in your path.  Run timetrace-gui and point your browser at localhost:8090.

Setting up a reverse proxy to enable access from the intenet is left as an exercise for the reader.

On first run a user creation dialog will be presented.  Use the user/password you entered for future logins.

## Screenshots
### Web
![browser](https://github.com/mattkasun/timetrace-gui/raw/master/screenshots/web.png "TimeTrace-GUI with Browser")

### Mobile
![phone](https://github.com/mattkasun/timetrace-gui/raw/master/screenshots/mobile.png "TimeTrace-GUI with Phone")

## RoadMap
- [ ] Restore deleted projects
- [ ] Edit project
- [ ] Edit record
- [ ] Reports
- [ ] Edit Configuration
- [ ] Users
- [ ] Docker without permission issues
