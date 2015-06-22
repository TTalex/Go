/*
FlickrQuizz is a web-based game using Flickr and Google APIs where the player has to guess locations of random pictures within a city.
An initial configuration allows the player to:
* Select the city he wants to play within
* Use a specific seed if he wants to challenge his friends on the same set of pictures
* Select a preferred game mode, endless continues until the player makes too many errors in a row, the other modes limit the number of pictures

@Devs: Rember to specify your own Google and Flickr API keys !
*/
package main

import (
	"fmt"
	"net/http"
	"text/template"
	"os"
	"apicaller"
	"strconv"
	"time"
	"math"
	"net/url"
	"math/rand"
	"strings"
)

//The Contents struct is used to fill static html templates with dynamic content
type Contents struct{
	Image string
	Lat float64
	Lng float64
	Circle string
	Resultmsg string
	Animations string
	Score string
	Hidden string
}

//The Photo struct contains information about a picture retrieved from Flickr
type Photo struct{
	Url string
	Link string
	Lat float64
	Lng float64
	Maxpages int
}

//The City struct contains information about a city retrieved from Google Geo API
type City struct{
	Name string
	Lat float64
	Lng float64
}

//Uses the contents of a decoded json to fill a Photo struct
func decodepicture(json map[string]interface{}, p *Photo){
	photos := json["photos"].(map[string]interface{})
	maxpages := int(photos["pages"].(float64))
	if (photos["total"] == "0"){
		fmt.Println("No photos found")
		return
	}
	photo := photos["photo"].([]interface{})[0].(map[string]interface{})
	farmid := photo["farm"]
	serverid := photo["server"]
	photoid := photo["id"]
	photosecret := photo["secret"]
	userid := photo["owner"]
	lat, _ := strconv.ParseFloat(photo["latitude"].(string), 64)
	lng, _ := strconv.ParseFloat(photo["longitude"].(string), 64)

	photourl := fmt.Sprintf("https://farm%.0f.staticflickr.com/%s/%s_%s_n.jpg", farmid, serverid, photoid, photosecret)

	photolink := fmt.Sprintf("https://www.flickr.com/photos/%s/%s", userid, photoid)
	*p = Photo{Url: photourl, Link: photolink, Lat: lat, Lng: lng, Maxpages: maxpages}

}

//Calls the Flickr API to retrieve a picture information formatted in json
func findpicture(citytag string, pagenumber int) Photo{
	var p Photo
	api_key := "f02cd0b3f01902b08393489f9ad04eab"
	str := fmt.Sprintf("https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=%s&tags=%s&has_geo=true&extras=geo&per_page=1&format=json&nojsoncallback=1&page=%d",api_key, citytag, pagenumber)

	m, err := apicaller.Callapi(str)
	if (err != nil){
		fmt.Println(err)
		return p
	}
	decodepicture(m, &p)
	return p
}

// Uses the contents of a decoded json to fill a City struct
func decodecity(name string, json map[string]interface{}, city *City){
	if (json["status"] != "OK"){
		fmt.Println(json["status"])
		return
	}
	results := json["results"].([]interface{})
	firstresult := results[0].(map[string]interface{})
	geo := firstresult["geometry"].(map[string]interface{})
	location := geo["location"].(map[string]interface{})
	lat := location["lat"].(float64)
	lng := location["lng"].(float64)
	
	*city = City{Name: name, Lat: lat, Lng: lng}
}

//Calls the Google Geo API to retrieve a city information formatted in json from a city name
func findcity(name string) City{
	var city City
	api_key := "AIzaSyB3FzjaP3F-piuwV0mRdawQhX7R5ogJdx8"
	str := "https://maps.googleapis.com/maps/api/geocode/json?address="+name+"&key="+api_key
	m, err := apicaller.Callapi(str)
	if (err != nil){
		fmt.Println(err)
		return city
	}
	decodecity(name, m, &city)
	return city
}

//Outputs the string of an Html hidden input containing the argument as value
func hiddenstr(hidden string) string{
	return fmt.Sprintf("<input type=\"hidden\" name=\"hidden\" value=%s />", hidden)
}

//Generates the HTML page sent back to the user when the game is still in progress
func reply(w http.ResponseWriter, photo Photo, city City, circle string, resultmsg string, hidden string, astr string){
	t, _ := template.ParseFiles(os.Getenv("GOPATH")+"/lib/flickrquizz/result.html")
	var s string
	if (photo.Link == ""){
		s = "No photo found for this location."
	}else{
		s = "<a href=\""+photo.Link+"\"><img src=\""+photo.Url+"\"></img></a>"
	}
	hiddenstring := hiddenstr(hidden)

	contents := &Contents{Image: s, Lat: city.Lat, Lng: city.Lng, Circle: circle, Resultmsg: resultmsg, Hidden: hiddenstring, Animations: astr}
	t.ExecuteTemplate(w, "Contents", contents)
}

//Generates the HTML page sent back to the user when the game is over (win or lose)
func replygameover(w http.ResponseWriter, city City, circle string, score float64, astr string){
	t, _ := template.ParseFiles(os.Getenv("GOPATH")+"/lib/flickrquizz/resultgameover.html")
	scorestr := fmt.Sprintf("%.0f",score)
	contents := &Contents{Lat: city.Lat, Lng: city.Lng, Circle: circle, Score: scorestr, Animations: astr}
	t.ExecuteTemplate(w, "Contents", contents)
}

//Returns a string containing the javascript code to display a circle on a Google map element
func circlestr(color string, lat float64, lng float64) string{
	return fmt.Sprintf(`
			   var marker={
			       position:  new google.maps.LatLng(%f,%f),
			       map: map,
			       icon:{      
				   path: google.maps.SymbolPath.CIRCLE,
				   strokeColor: '%s',
				   strokeOpacity: 0.8,
				   strokeWeight: 2,
				   fillColor: '%s',
				   fillOpacity: 0.35,
				   scale: 10
			       }      
			   };
			   new google.maps.Marker(marker);
			   `,lat, lng, color, color)
}

//Returns a string containing the css code used to animate the score bar
func animationsstr(positive bool, value float64) string{
	var left string
	var right string
	var background string
	if (positive){
		left = "50%"
		right = "auto"
		background = "lightgreen, green"
	}else{
		left = "auto"
		right = "50%"
		background = "red, orange"
	}
	dostuff := fmt.Sprintf("0%%{width:0%%;} 100%%{width:%f%%;}",value)
	domorestuff := fmt.Sprintf("0%%{width:%f%%;} 50%%{width:%f%%;} 100%%{width:%f%%;}",value, math.Abs(value-1.0), value)
	return fmt.Sprintf(`
			   .wave{
			       left:%s;
			       right:%s;
			       background: -webkit-linear-gradient(left, %s);
			       background: -o-linear-gradient(right, %s);
			       background: -moz-linear-gradient(right %s);
			       background: linear-gradient(to right %s);
			   }
			   @keyframes dostuff{%s}
			   @-webkit-keyframes dostuff{%s}
			   @-moz-keyframes dostuff{%s}
			   @-o-keygrames dostuff{%s}

			   @keyframes domorestuff{%s}
			   @-webkit-keyframes domorestuff{%s}
			   @-moz-keyframes domorestuff{%s}
			   @-o-keygrames domorestuff{%s}
			   `,left, right, background, background, background, background, dostuff, dostuff, dostuff, dostuff, domorestuff, domorestuff, domorestuff, domorestuff)
}

//Returns true or false if the player guess was correct or not given the photo and city coordinates
func checkresult(action string, photolng float64, citylng float64) bool{
	return (photolng <= citylng && action == "West") || (photolng >= citylng && action == "East")
}

//Returns a string containing hidden fields serialized
func serializehidden(photo Photo, maxpages int, dscore float64, score float64, seed int64, counter int64) string{
	return fmt.Sprintf("%f;%f;%d;%f;%f;%d;%d",photo.Lng, photo.Lat, maxpages, dscore, score, seed, counter)
}

//Reverse function from serializehidden, returns each hidden field from the serialized string
func unserializehidden(s string) (float64, float64, int, float64, float64, int64, int64){
	split := strings.Split(s, ";")
	photolng, _ := strconv.ParseFloat(split[0], 64)
	photolat, _ := strconv.ParseFloat(split[1], 64)
	maxpages, _ := strconv.ParseInt(split[2], 10, 0)
	dscore, _ := strconv.ParseFloat(split[3], 64)
	score, _ := strconv.ParseFloat(split[4], 64)
	seed, _ := strconv.ParseInt(split[5], 10, 64)
	counter, _ := strconv.ParseInt(split[6], 10, 64)
	return photolng, photolat, int(maxpages), dscore, score, seed, counter
}

//Handles a request on the /play/ root. Retrieves information from both the initial config form (GET) and the dynamic form (POST)
//Also computes the progress of the player and keep track of scores and multipliers
func handler2(w http.ResponseWriter, r *http.Request) {
	var circle string
	var resultmsg string
	var dscore float64
	var astr string
	var score float64
	var seed int64
	var maxpages int
	var counter int64
	maxpages = 100
	counter = 0

	// GET form
	cityname := r.FormValue("city")
	cityname = url.QueryEscape(cityname)
	seed, _ = strconv.ParseInt(r.FormValue("seed"), 10, 0)
	city := findcity(cityname)
	gamemode, _ := strconv.ParseInt(r.FormValue("gamemode"), 10, 64)

	// POST form
	action := r.FormValue("action")
	hidden := r.FormValue("hidden")

	//The first time, the POST form should contain empty values, there is no need to compute any score
	//The next times the page is loaded, the POST form values help compute the current score
	if (action != "" && hidden != ""){
		var photolng float64
		var photolat float64
		photolng, photolat, maxpages, dscore, score, seed, counter = unserializehidden(hidden)
		counter += 1
		win := checkresult(action, photolng, city.Lng)
		if (win){
			//When the player guess the correct location of a picture, a green circle is drawn on the map at the picture location.
			//The score is increased with a potential multiplier
			//A result message informs the player of his success
			circle = circlestr("#00FF00", photolat, photolng)
			multiplier := dscore/10 + 1
			multiplier = math.Max(multiplier, 1)
			score += 1 * multiplier
			dscore += 10.0
			dscore = math.Min(dscore, 50.0)
			resultmsg = fmt.Sprintf("#%d: YAY :) - Score %.0f - Multiplier x%.0f",counter, score, multiplier)
		}else{
			//When the player guess the incorrect location of a picture, a red circle is drawn on the map at the picture location.
			//The dynamic score is decreased with a potential multiplier
			//A result message informs the player of his failure
			circle = circlestr("#FF0000", photolat, photolng)
			dscore = math.Min(dscore, 10)
			dscore -= 10.0
			if(score >= 50){
				dscore -= 10.0
			}
			if(score >= 100){
				dscore -= 10.0
			}
			dscore = math.Max(dscore, -50.0)
			resultmsg = fmt.Sprintf("#%d: NEY :( - Score %.0f - Multiplier x1",counter, score)
		}
		astr = animationsstr(dscore>0,math.Abs(dscore))
	}
	//Test if the gamemode is not endless
	if(gamemode>1){
		//The number of pictures have reach the game limit
		if(counter == gamemode){
			replygameover(w, city, circle, score, astr)
			return
		}
	}
	//Test if player can still play
	if(dscore>-50.0){
		rand.Seed(seed)
		//Flicker API limits the photo search results to 4000
		newseed := rand.Intn(maxpages%4000)

		photo := findpicture(cityname, newseed)
		hidden = serializehidden(photo, photo.Maxpages, dscore, score, int64(newseed), counter)
		reply(w, photo, city, circle, resultmsg, hidden, astr)
	}else{
		replygameover(w, city, circle, score, astr)
	}
}

//Handles a request on the / root.
func indexhandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(os.Getenv("GOPATH")+"/lib/flickrquizz/index.html")
	t.ExecuteTemplate(w, "Seed", rand.Int())
}

//The main function defines handlers for various paths and sets up the webserver
func main() {
	rand.Seed(int64((time.Now()).Nanosecond()))
	http.HandleFunc("/", indexhandler)
	http.HandleFunc("/play/", handler2)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(os.Getenv("GOPATH")+"/lib/flickrquizz/css"))))
	http.ListenAndServe(":8081", nil)
}
