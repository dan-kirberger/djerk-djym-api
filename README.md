# Djerk Djym API app

Go training wheels app

###TODO:
 - [ ] Separate files for routes and stuff
 - [ ] Consolidate response handling junk
 - [ ] Error handling stuff (return 500, not panic)
 - [ ] Log request info on error
 - [ ] Get rid of globalContext thingy
 - [ ] Talk to Mongo cloud thingy from app engine or wherever
 - [ ] Deploy to app engine from drone
 - [ ] Separate driver app - static web app?
 - [ ] Request logging (timing, url??)
 - [ ] Parameterize mongo info (done for tests now, not for app)
 - [ ] Vendoring/module option? Don't really want to vendor the 18MB mongo repo
 
###Stuff:
* Code: here
* CI:  [Drone](https://cloud.drone.io/dan-kirberger/djerk-djym-api)
* Live App: https://djerk-djym-api.appspot.com/ (Google App Engine)