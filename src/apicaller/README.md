##Introduction
The [apicaller](https://github.com/TTalex/Go/tree/master/src/apicaller) project is a Go library that helps when making requests to APIs using the Json format.

##Methods
###func Callapi(url string) (map[string]interface{}, error)
Single call to an API url specified via `url` parameter.
Returns a `map[string]interface{}` of the decoded Json response and an error if the API call failed.

###func Callapisem(url string, c chan bool) (map[string]interface{}, error)
Single call to an API url specified via `url` parameter. A channel acting as a semaphore to limit the number of concurrent calls is specified via `c`.
Blocks in case the semaphore is filled.
Returns a `map[string]interface{}` of the decoded Json response and an error if the API call failed.

##Installation
```
go install apicaller
```

##Usage
####Without API limits
```go
func main(){
	m, err := apicaller.Callapi("http://www.omdbapi.com/?t=Kill+Bill&y=&plot=short&r=json")
	if (err != nil){
		fmt.Println("Error")
		return
	}
	if m["Response"] == "True" {
		fmt.Println(m["imdbRating"])
	} else {
		fmt.Println("Error")
	}
}
```
