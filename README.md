# Code generation

Application generates structures and simple functions by the json data.

## Incoming data

See the `./data.json` file.
Incoming data is a set of different entity types of the sport-related entities.

## Running
Execute `go generate ./...` in the root of the project. 

## Outcome

See the `./entity` folder.
The generated result of the application running is:
- structures with all possible fields. 
If some entities come with a few fields, and other entities has more, 
then the resulted structure will collect all of them.
- simple functions based on fields type. 
As example, it generated entity function IsActive for a bool field `isActive.
