# moduledb_htwsaar_parser
moduledb_htwsaar_parser is a package that provides,
a html parser for the parsing of a specific module of the available 
htw saar courses.

## Usage
the packages Run function takes a http Client as an input, aswell 
as the url of the course that shall be parsed by this package. 

## Internal
- get available courses and their specific urls wich are stored in the html source code of the main page

- extract with the shortened version of the timetables id the specific link out of the created map

