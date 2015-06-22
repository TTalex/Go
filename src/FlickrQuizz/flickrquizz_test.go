package main

import "testing"

func TestFindcity(t *testing.T){
	Paris := &City{Name: "paris", Lat: 48.856614, Lng: 2.3522219}
	city := findcity("paris")
	if (city.Name == ""){
		t.Error("City not filled")
	}
	if (Paris.Name != city.Name){
		t.Error("Error on city Name")
	}
	if (Paris.Lat != city.Lat){
		t.Error("Error on city Lat")
	}
	if (Paris.Lng != city.Lng){
		t.Error("Error on city Lng")
	}
}

func TestFindcity_2(t *testing.T){
	city := findcity("")
	if (city.Name != ""){
		t.Error("Erro on city Name")
	}
	if (city.Lat != 0.0){
		t.Error("Error on city Lat")
	}
	if (city.Lng != 0.0){
		t.Error("Error on city Lng")
	}
}

func TestFindpicture(t *testing.T){
//	photo := &Photo{Url: "https://farm1.staticflickr.com/301/18875586420_59b4c8120e_n.jpg", Link: "https://www.flickr.com/photos/130654398@N05/18875586420", Lat: 48.861305, Lng: 2.288782, Maxpages: 362398}
	p := findpicture("Paris", 1)
	if (p.Url == ""){
		t.Error("Photo not filled")
	}
	/* Flickr search API is unstable, gives different results on same call
	if (photo.Url != p.Url){
		t.Error("Error on photo url")
	}
	if (photo.Link != p.Link){
		t.Error("Error on photo link")
	}
	if (photo.Lat != p.Lat){
		t.Error("Error on photo Lat")
	}
	if (photo.Lng != p.Lng){
		t.Error("Error on photo Lng")
	}
*/
}

func TestSerializehidden(t *testing.T){
	photo := Photo{Url: "https://farm1.staticflickr.com/301/18875586420_59b4c8120e_n.jpg", Link: "https://www.flickr.com/photos/130654398@N05/18875586420", Lat: 48.861305, Lng: 2.288782, Maxpages: 362398}
	var dscore float64
	var score float64
	var seed int64
	var counter int64
	dscore = 10.0
	score = 5.0
	seed = 1
	counter = 42
	str := serializehidden(photo, photo.Maxpages, dscore, score, seed, counter)
	plng, plat, mp, ds, s, se, c := unserializehidden(str)
	if (plng != photo.Lng){
		t.Error("Error on photo Lng")
	}
	if (plat != photo.Lat){
		t.Error("Error on photo Lat")
	}
	if (mp != photo.Maxpages){
		t.Error("Error on photo Maxpages")
	}
	if (ds != dscore){
		t.Error("Error on dscore")
	}
	if (s != score){
		t.Error("Error on score")
	}
	if (se != seed){
		t.Error("Error on seed")
	}
	if (c != counter){
		t.Error("Error on counter")
	}
}

func TestUnserializehidden(t *testing.T){
	str := "48.800000;2.280000;362498;10.000000;5.000000;1;42"
	plng, plat, mp, ds, s, se, c := unserializehidden(str)
	photo := Photo{Lat: plat, Lng: plng}
	st := serializehidden(photo, mp, ds, s, se, c)
	if (str != st){
		t.Error("Error in serialize")
	}
}
