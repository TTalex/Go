##Introduction
The [films](https://github.com/TTalex/Go/tree/master/src/films) project is a tool that helps ranking movies available on the replay platform of [Canal+](http://replay.mycanal.fr/cplus/selection) by [Imdb](http://imdb.com) scores.

##Used API
* [OMDb API](http://www.omdbapi.com/) to retrieve imdb scores per movie.

##Installation
1. Make sur Golang is installed
2. Clone the repo

  ```
  git clone https://github.com/TTalex/Go.git
  ```
  
3. Set up the $GOPATH environement variable

  ```
  export GOPATH=`pwd`/Go
  ```

##Running it
1. Run the webserver

  ```
  cd Go
  go run films
  ```
  
2. Access [localhost:8080](http://localhost:8080)

##Known bugs
Because of the mismatching between the French movie titles on the Canal+ platform and the English OMDb, movies data might not be found.
