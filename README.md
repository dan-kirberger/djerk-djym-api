# Djerk Djym API app

Go training wheels app

###TODO:
 - [ ] Separate files for routes and stuff
 - [ ] Consolidate response handling junk
 - [ ] Error handling stuff (return 500, not panic)
 - [ ] Log request info on error
 - [ ] Integration Test (stand up mongo in docker?) - unit test looks annoying with web/db stuff
 - [ ] Mongo connection pool/dont init the whole thing every request
 - [ ] Talk to Mongo cloud thingy
 - [ ] Talk to Mongo cloud thingy from app engine or wherever
 - [ ] Deploy to app engine from drone
 - [ ] Separate driver app - static web app?
 - [ ] Request logging (timing, url??)
 
###Stuff:
* Code: here
* CI:  [Drone](https://cloud.drone.io/dan-kirberger/djerk-djym-api)
* Live App: https://djerk-djym-api.appspot.com/ (Google App Engine)