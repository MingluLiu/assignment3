package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
	"encoding/json"
     "io/ioutil"
     "gopkg.in/mgo.v2/bson"
     "io"
     "net/url"
    "math/rand"
     "bytes"
    "github.com/julienschmidt/httprouter"
	// "github.com/MingluLiu/assignment3"
	// "github.com/anweiss/uber-api-golang/uber"
)
var RandomID int

type Arguments struct{
    Name string `json:"name"    bson:"name"`
    Address string `json:"address" bson:"address"`
    City string `json:"city"    bson:"city"`
    State string `json:"state"   bson:"state"`
    Zip string `json:"zip"     bson:"zip"`
}

type Response struct{
    Id bson.ObjectId `json:"id"      bson:"_id,omitempty"`
    Name string `json:"name"    bson:"name"`
    Address string `json:"address" bson:"address"`
    City string `json:"city"    bson:"city"`
    State string `json:"state"   bson:"state"`
    Zip string `json:"zip"     bson:"zip"`
    Coordinate struct {
        Lat float64 `json:"lat" bson:"lat"`
        Lng float64 `json:"lng" bson:"lng"`
    } `json:"coordinate"        bson:"coordinate"`
}

type TripArguments struct{
	Starting_from_location_id string `json:"starting_from_location_id"`
	Location_ids []string
}

type TripResponse struct{
	Id     bson.ObjectId 				  `json:"_id" bson:"_id,omitempty"`
	Status string  						  `json:"status"`
	Starting_from_location_id   string    `json:"starting_from_location_id"`
	Best_route_location_ids []string
	Total_uber_costs int			  `json:"total_uber_costs"`
	Total_uber_duration int			  `json:"total_uber_duration"`
	Total_distance float64				  `json:"total_distance"`
}

type TripPutOutput struct{
	Id     bson.ObjectId 				  `json:"_id" bson:"_id,omitempty"`
	Status string  						  `json:"status"`
	Starting_from_location_id   string    `json:"starting_from_location_id"`
	Next_destination_location_id   string    `json:"next_destination_location_id"`
	Best_route_location_ids []string
	Total_uber_costs int			  `json:"total_uber_costs"`
	Total_uber_duration int			  `json:"total_uber_duration"`
	Total_distance float64			  `json:"total_distance"`
	Uber_wait_time_eta int 			  `json:"uber_wait_time_eta"`

}

type Struct_for_put struct{
	trip_route []string
	trip_visits map[string]int
}

type Final_struct struct{
	theMap map[string]Struct_for_put
}

type PriceEstimates struct {
	Prices         []PriceEstimate `json:"prices"`
}
// type PriceEstimates struct {
// 	StartLatitude  float64
// 	StartLongitude float64
// 	EndLatitude    float64
// 	EndLongitude   float64
// 	Prices         []PriceEstimate `json:"prices"`
// }

type PriceEstimate struct {
	ProductId       string  `json:"product_id"`
	CurrencyCode    string  `json:"currency_code"`
	DisplayName     string  `json:"display_name"`
	Estimate        string  `json:"estimate"`
	LowEstimate     int     `json:"low_estimate"`
	HighEstimate    int     `json:"high_estimate"`
	SurgeMultiplier float64 `json:"surge_multiplier"`
	Duration        int     `json:"duration"`
	Distance        float64 `json:"distance"`
}

type UberOutput struct{
	Cost int
	Duration int
	Distance float64
}

type UberETA struct{
  Request_id string `json:"request_id"`
  Status string	    `json:"status"`
  Vehicle string	`json:"vehicle"`
  Driver string		`json:"driver"`
  Location string	`json:"location"`
  ETA int			`json:"eta"`
  SurgeMultiplier float64 `json:"surge_multiplier"`
}

type Internal_data struct{
	Id string               `json:"_id" bson:"_id,omitempty"`
	Trip_visited []string  `json:"trip_visited"`
	Trip_not_visited []string  `json:"trip_not_visited"`
	Trip_completed int        `json:"trip_completed"`
}

func GenerateID() int{
    if RandomID == 0{
        for RandomID == 0{
            RandomID = rand.Intn(99999)
        }
    }else{
        RandomID = RandomID + 1
    }
    return RandomID
}


func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Create",
		"POST",
		"/locations",
		Create,
	},
	Route{
		"Query",
		"GET",
		"/locations/{location_id}",
		Query,
	},
	Route{
		"Update",
		"PUT",
		"/locations/{location_id}",
		Update,
	},
	Route{
		"Delete",
		"DELETE",
		"/locations/{location_id}",
		Delete,
	},
	// Route{
	// 	"TripCreate",
	// 	"POST",
	// 	"/trips",
	// 	TripCreate,
	// },
	// Route{
	// 	"TripQuery",
	// 	"GET",
	// 	"/trips/{trip_id}",
	// 	TripQuery,
	// },
	// Route{
	// 	"TripUpdate",
	// 	"PUT",
	// 	"/trips/{trip_id}",
	// 	TripUpdate,
	// },
}

type LocationController struct {
		session *mgo.Session
	}


type LocationAllInfo struct {
	Results []struct {
    AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry struct {
			Location struct {
			       Lat float64 `json:"lat"`
			       Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport  struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

func GoogleAPI(Address string) (Response, error) {

	url1 :=  "http://maps.google.com/maps/api/geocode/json?address="
	url2 := url.QueryEscape(Address)
	url3 := "&sensor=false"
	fullUrl := url1 + url2 +url3

	fmt.Println(fullUrl)
	var locationAllInfo Response
	var l LocationAllInfo
	
	res, err := http.Get(fullUrl)
	if err!=nil {
		fmt.Println("GoogleAPI: http.Get",err)
		return locationAllInfo,err
	}
	defer res.Body.Close()

	body,err := ioutil.ReadAll(res.Body)
	if err!=nil {
		fmt.Println("GoogleAPI: ioutil.ReadAll",err)
		return locationAllInfo,err
	}

	err = json.Unmarshal(body, &l)

	if err!=nil {
		fmt.Println("GoogleAPI: json.Unmarshal",err)
		return locationAllInfo,err
	}

	locationAllInfo.Coordinate.Lat = l.Results[0].Geometry.Location.Lat;
	locationAllInfo.Coordinate.Lng = l.Results[0].Geometry.Location.Lng;

	return locationAllInfo,nil

}

func Create(w http.ResponseWriter, r *http.Request) {
    var req Arguments
    var res Response

    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }

    if err := json.Unmarshal(body, &req); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
        fmt.Println("Unmarshal Json Error.", body)
        return
    }
    fullAddress := req.Address+","+req.City+","+req.State+","+req.Zip

    Information,err := GoogleAPI(fullAddress);

    res.Address = req.Address;
    res.City = req.City;
    res.State = req.State;
    res.Zip = req.Zip;
    res.Name = req.Name;
    res.Coordinate.Lat = Information.Coordinate.Lat
    res.Coordinate.Lng = Information.Coordinate.Lng
    MongoCreate(&res)

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)

    if err := json.NewEncoder(w).Encode(res); err != nil {
        panic(err)
    }
    return
}

func Delete(w http.ResponseWriter, r *http.Request) {

    var res Response

    vars := mux.Vars(r)
    res.Id = bson.ObjectIdHex(vars["location_id"])

    err := MongoDelete(res.Id)
    if err != nil {

        fmt.Printf(err.Error())
        return
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)

    return
}

func Update(w http.ResponseWriter, r *http.Request) {
    var req Arguments
    var res Response
    vars := mux.Vars(r)
    res.Id = bson.ObjectIdHex(vars["location_id"])

    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }

    if err := json.Unmarshal(body, &req); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) 
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
        fmt.Println("Unmarshal Json Error.", body)
        return
    }

    fullAddress := req.Address+","+req.City+","+req.State+","+req.Zip

    Information,err := GoogleAPI(fullAddress);

    res.Address = req.Address;
    res.City = req.City;
    res.State = req.State;
    res.Zip = req.Zip;
    res.Coordinate.Lat = Information.Coordinate.Lat
    res.Coordinate.Lng = Information.Coordinate.Lng

    Information, err = MongoUpdate(res)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    res.Name = Information.Name

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)

    if err := json.NewEncoder(w).Encode(res); err != nil {
        panic(err)
    }

}

func Query(w http.ResponseWriter, r *http.Request) {

    var res Response
    var err error
    vars := mux.Vars(r)
    res.Id = bson.ObjectIdHex(vars["location_id"])


    res, err = MongoQuery(res.Id )
    if err != nil {

        fmt.Printf(err.Error())
        return
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)

    if err := json.NewEncoder(w).Encode(res); err != nil {
        panic(err)
    }
}

func MongoCreate(Information *Response) {
//mongodb://<dbuser>:<dbpassword>@ds039484.mongolab.com:39484/<db_name>
	sess, err := mgo.Dial("mongodb://minglu:liu273@ds057862.mongolab.com:57862/mingluliumongodb")
	if err != nil {
		fmt.Printf("MongoDB connection error %v\n", err)
		panic(err)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})
	db := sess.DB("mingluliumongodb").C("Location")

	Information.Id = bson.NewObjectId()

	err = db.Insert(&Information)
	if err != nil {
		fmt.Printf("MongoDB create error %v\n", err)
		os.Exit(1)
	}

	var results []Response
	err = db.Find(bson.M{"_id": Information.Id}).Sort("-timestamp").All(&results)

	if err != nil {
		panic(err)
	}

	err = db.Find(bson.M{}).Sort("-timestamp").All(&results)

	if err != nil {
		panic(err)
	}
}

func MongoDelete(Id bson.ObjectId) (error) {

	sess, err := mgo.Dial("mongodb://minglu:liu273@ds057862.mongolab.com:57862/mingluliumongodb")
	if err != nil {
		fmt.Printf("MongoDB connection error  %v\n", err)
		panic(err)
		return err;
	}
	defer sess.Close()
	
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("mingluliumongodb").C("Location")

	err = collection.Remove(bson.M{"_id": Id})
	if err != nil {
		panic(err)
		return err
	}

	return nil
}

func MongoQuery(Id bson.ObjectId) (Response, error) {

    var result []Response
	sess, err := mgo.Dial("mongodb://minglu:liu273@ds057862.mongolab.com:57862/mingluliumongodb")
	if err != nil {
		fmt.Printf("MongoDB connection error %v\n", err)
		panic(err)
		return result[0], err
	}
	defer sess.Close()
	
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("mingluliumongodb").C("Location")

	err = collection.Find(bson.M{"_id": Id}).All(&result)

	if err != nil {
		panic(err)
		return result[0], err
	}

	return result[0],nil
}

func MongoUpdate(Information Response) (Response, error) {

	var LocationInfo Response

	sess, err := mgo.Dial("mongodb://minglu:liu273@ds057862.mongolab.com:57862/mingluliumongodb")
	if err != nil {
		fmt.Printf("MongoDB connection error %v\n", err)
		panic(err)
		return LocationInfo,err;
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("mingluliumongodb").C("Location")

	colQuerier := bson.M{"_id": Information.Id}

	change := bson.M{"$set": bson.M{"address": Information.Address,
									"city": Information.City,
									"state": Information.State,
									"zip": Information.Zip,
									"coordinate": bson.M{"lat":Information.Coordinate.Lat,
											"lng":Information.Coordinate.Lng}}}
	err = collection.Update(colQuerier, change)
	if err != nil {
		panic(err)
		return LocationInfo,err
	}
	Response,error := MongoQuery(Information.Id)
	return Response,error
}
type Service struct{}


func Get_uber_price(startLat, startLng, endLat, endLng string) UberOutput{
	client := &http.Client{}
	reqURL := fmt.Sprintf("https://sandbox-api.uber.com/v1/estimates/price?start_latitude=%s&start_longitude=%s&end_latitude=%s&end_longitude=%s&server_token=btQIERoY3u8zIP_dmYceMF1jTbfg2U7vTi9N-0aA", startLat, startLng, endLat, endLng)
	fmt.Println("URL formed: "+ reqURL)
	req, err := http.NewRequest("GET", reqURL , nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error in sending req to Uber: ", err);	
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error in reading response: ", err);	
	}

	var res PriceEstimates
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("error in unmashalling response: ", err);	
	}

	var uberOutput UberOutput
	uberOutput.Cost = res.Prices[0].LowEstimate
	uberOutput.Duration = res.Prices[0].Duration
	uberOutput.Distance = res.Prices[0].Distance

	return uberOutput

}


func Get_uber_eta(startLat, startLng, endLat, endLng string) int{

	var jsonStr = []byte(`{"start_latitude":"` + startLat + `","start_longitude":"` + startLng + `","end_latitude":"` + endLat + `","end_longitude":"` + endLng + `","product_id":"04a497f5-380d-47f2-bf1b-ad4cfdcb51f2"}`)
	reqURL := "https://sandbox-api.uber.com/v1/requests"
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicmVxdWVzdCJdLCJzdWIiOiI3NzQ3OTY5ZS01YWZkLTQ3NjUtODI5Ni04NGUwZmNiZTZiMjMiLCJpc3MiOiJ1YmVyLXVzMSIsImp0aSI6IjEyMzdmZjZhLTUxMGEtNGEzMS05MjgzLWMxODhlZmM2MDk3ZiIsImV4cCI6MTQ1MDg0NDM2MCwiaWF0IjoxNDQ4MjUyMzU5LCJ1YWN0IjoiTW12aWpNRFFNMDN0dUhTblpvTHIyUkZ4bFB5aEZaIiwibmJmIjoxNDQ4MjUyMjY5LCJhdWQiOiJsNUxsVW5XX3lSbGFwVVpMdTd6ZDlIbjA0dVZHOUxMdyJ9.fubHpBsKrbCIYfeK_kVLP1TOOH4NLqwVDGgLVNfEe0Bj7tHOCd3dfOEWKxdVJdFF7rsm7vcvZOpXDdvCsV_NJiiSVLKoNRt0hDMM5LxBI_WhISfY_SLqdD3M8reRwCJPMhg8vwmdRidSlor6NAfsupBt_b84LqK-SeHcOqpFhR8qgtkkmZTug93tpvQ6P0AZ7qGPLTpxqhEdW-l553XNpp5QdWTI3YzR223PeFfpAVsxZweCFlQtim5TytFNblS5SLdKdjmslIjFgKHzpyn4xBUvMw0gQ3ZElAE9tQIAXyGloKg3c1fTiv2AP6OdE1lkKqdcQObBClPjN1k5LsODAQ")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error in sending req to Uber: ", err);	
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error in reading response: ", err);	
	}

	var res UberETA
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("error in unmashalling response: ", err);	
	}
	eta:= res.ETA
	return eta
	
}

func (uc LocationController) TripCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var ta TripArguments
	var tr TripResponse
	var cost_array []int
	var duration_array []int
	var distance_array []float64
	cost_total := 0
	duration_total := 0
	distance_total := 0.0

	json.NewDecoder(r.Body).Decode(&ta)	

	starting_id:= bson.ObjectIdHex(ta.Starting_from_location_id)
	var start Response
	if err := uc.session.DB("mingluliumongodb").C("Location").FindId(starting_id).One(&start); err != nil {
       	w.WriteHeader(404)
        return
    }
    start_Lat := start.Coordinate.Lat
    start_Lng := start.Coordinate.Lng

    for len(ta.Location_ids)>0{
	
			for _, loc := range ta.Location_ids{
				id := bson.ObjectIdHex(loc)
				var res Response
				if err := uc.session.DB("mingluliumongodb").C("Location").FindId(id).One(&res); err != nil {
		       		w.WriteHeader(404)
		        	return
		    	}
		    	loc_Lat := res.Coordinate.Lat
		    	loc_Lng := res.Coordinate.Lng
		    	
		    	getUberResponse := uber.Get_uber_price(start_Lat, start_Lng, loc_Lat, loc_Lng)
		    	fmt.Println("Uber Response is: ", getUberResponse.Cost, getUberResponse.Duration, getUberResponse.Distance );
		    	cost_array = append(cost_array, getUberResponse.Cost)
		    	duration_array = append(duration_array, getUberResponse.Duration)
		    	distance_array = append(distance_array, getUberResponse.Distance)
		    	
			}
			fmt.Println("Cost Array", cost_array)

			min_cost:= cost_array[0]
			var indexNeeded int
			for index, value := range cost_array {
		        if value < min_cost {
		            min_cost = value 
		            indexNeeded = index
		        }
		    }

			cost_total += min_cost
			duration_total += duration_array[indexNeeded]
			distance_total += distance_array[indexNeeded]

			tr.Best_route_location_ids = append(tr.Best_route_location_ids, ta.Location_ids[indexNeeded])

			starting_id = bson.ObjectIdHex(ta.Location_ids[indexNeeded])
			if err := uc.session.DB("mingluliumongodb").C("Location").FindId(starting_id).One(&start); err != nil {
       			w.WriteHeader(404)
        		return
    		}
    		ta.Location_ids = append(ta.Location_ids[:indexNeeded], ta.Location_ids[indexNeeded+1:]...)

    		start_Lat = start.Coordinate.Lat
    		start_Lng = start.Coordinate.Lng

    		cost_array = cost_array[:0]
    		duration_array = duration_array[:0]
    		distance_array = distance_array[:0]

	}


	Last_loc_id := bson.ObjectIdHex(tr.Best_route_location_ids[len(tr.Best_route_location_ids)-1])
	var resp Response
	if err := uc.session.DB("mingluliumongodb").C("Location").FindId(Last_loc_id).One(&resp); err != nil {
		w.WriteHeader(404)
		return
	}
	last_loc_Lat := resp.Coordinate.Lat
	last_loc_Lng := resp.Coordinate.Lng

	ending_id:= bson.ObjectIdHex(ta.Starting_from_location_id)
	var end Response
	if err := uc.session.DB("mingluliumongodb").C("Location").FindId(ending_id).One(&end); err != nil {
       	w.WriteHeader(404)
        return
    }
    end_Lat := end.Coordinate.Lat
    end_Lng := end.Coordinate.Lng
		    	
	getUberResponse_last := Get_uber_price(last_loc_Lat, last_loc_Lng, end_Lat, end_Lng)


	tr.Id = bson.NewObjectId()
	tr.Status = "planning"
	tr.Starting_from_location_id = ta.Starting_from_location_id
	tr.Total_uber_costs = cost_total + getUberResponse_last.Cost
	tr.Total_distance = distance_total + getUberResponse_last.Distance
	tr.Total_uber_duration = duration_total + getUberResponse_last.Duration
	

	uc.session.DB("mingluliumongodb").C("Trips").Insert(tr)

	uj, _ := json.Marshal(tr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}



func (uc LocationController) TripUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	var theStruct Struct_for_put
	var final Final_struct
	var Uber uber
	final.theMap = make(map[string]Struct_for_put)

	var tPO TripPutOutput
	var internal Internal_data

	id := p[0].Value
	if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        return
    }
    oid := bson.ObjectIdHex(id)
	if err := uc.session.DB("mingluliumongodb").C("Trips").FindId(oid).One(&tPO); err != nil {
        w.WriteHeader(404)
        return
    }


	theStruct.trip_route = tPO.Best_route_location_ids
    theStruct.trip_route = append([]string{tPO.Starting_from_location_id}, theStruct.trip_route...)
    fmt.Println("The route array is: ", theStruct.trip_route)
    theStruct.trip_visits = make(map[string]int)

    // theStruct.trip_route = list_location_ids
    var trip_visited []string 
    var trip_not_visited []string

  	if err := uc.session.DB("mingluliumongodb").C("Trip_internal_data").FindId(id).One(&internal); err != nil {
    	for index, loc := range theStruct.trip_route{
    		if index == 0{
    		// fmt.Println("Coming here.....................")
    			theStruct.trip_visits[loc] = 1
    			trip_visited = append(trip_visited, loc)
    		}else{
    			theStruct.trip_visits[loc] = 0
    			trip_not_visited = append(trip_not_visited, loc)
    		}
    	}
    	internal.Id = id
    	internal.Trip_visited = trip_visited
    	internal.Trip_not_visited = trip_not_visited
    	internal.Trip_completed = 0
    	uc.session.DB("mingluliumongodb").C("Trip_internal_data").Insert(internal)

    }else {
    	for _, loc_id := range internal.Trip_visited {
    		theStruct.trip_visits[loc_id] = 1
    	}
    	for _, loc_id := range internal.Trip_not_visited {
    		theStruct.trip_visits[loc_id] = 0
    	}
    }


  	fmt.Println("Trip visit map ", theStruct.trip_visits)
  	final.theMap[id] = theStruct


  	last_index := len(theStruct.trip_route) - 1
  	trip_completed := internal.Trip_completed
  	// last_elem = theStruct.trip_route[last_index]
  		// fmt.Println("Trip completed ==", trip_completed)
  	if trip_completed == 1 {
  		fmt.Println("Entering the trip completed if statement")
  		// tpost.Status = "completed"
  		tPO.Status = "completed"

		uj, _ := json.Marshal(tPO)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		fmt.Fprintf(w, "%s", uj)
		return
	}

	for i, location := range theStruct.trip_route{
	  	if  (theStruct.trip_visits[location] == 0){
	  		tPO.Next_destination_location_id = location
	  		nextoid := bson.ObjectIdHex(location)
			var res Response
			if err := uc.session.DB("mingluliumongodb").C("Location").FindId(nextoid).One(&res); err != nil {
        		w.WriteHeader(404)
        		return
    		}
    		nlat := res.Coordinate.Lat
    		nlng:= res.Coordinate.Lng

	  		if i == 0 {
	  			starting_point := theStruct.trip_route[last_index]
	  			startingoid := bson.ObjectIdHex(starting_point)
				var res Response
				if err := uc.session.DB("mingluliumongodb").C("Location").FindId(startingoid).One(&res); err != nil {
        			w.WriteHeader(404)
        			return
    			}
    			slat := res.Coordinate.Lat
    			slng:= res.Coordinate.Lng


	  			eta := UberETA.Get_uber_eta(slat, slng, nlat, nlng)
	  			tPO.Uber_wait_time_eta = eta
	  			trip_completed = 1
	  		}else {
	  			starting_point2 := theStruct.trip_route[i-1]
	  			startingoid2 := bson.ObjectIdHex(starting_point2)
				var res Response
				if err := uc.session.DB("mingluliumongodb").C("Location").FindId(startingoid2).One(&res); err != nil {
        			w.WriteHeader(404)
        			return
    			}
    			slat := res.Coordinate.Lat
    			slng:= res.Coordinate.Lng
	  			eta := UberETA.Get_uber_eta(slat, slng, nlat, nlng)
	  			tPO.Uber_wait_time_eta = eta
	  		}	

	  		fmt.Println("Starting Location: ", tPO.Starting_from_location_id)
	  		fmt.Println("Next destination: ", tPO.Next_destination_location_id)
	  		theStruct.trip_visits[location] = 1
	  		if i == last_index {
	  			theStruct.trip_visits[theStruct.trip_route[0]] = 0
	  		}
	  		break
	  	}
	}

	trip_visited  = trip_visited[:0]
	trip_not_visited  = trip_not_visited[:0]
	for location, visit := range theStruct.trip_visits{
		if visit == 1 {
			trip_visited = append(trip_visited, location)
		}else {
			trip_not_visited = append(trip_not_visited, location)
		} 
	}

	internal.Id = id
	internal.Trip_visited = trip_visited
	internal.Trip_not_visited = trip_not_visited
	fmt.Println("Trip Visisted", internal.Trip_visited)
	fmt.Println("Trip Not Visisted", internal.Trip_not_visited)
	internal.Trip_completed = trip_completed

	c := uc.session.DB("mingluliumongodb").C("Trip_internal_data")
	id2 := bson.M{"_id": id}
	err := c.Update(id2, internal)
	if err != nil {
		panic(err)
	}

    uj, _ := json.Marshal(tPO)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)

}

func (uc LocationController) TripQuery(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("trip_id")
	if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        return
    }
    oid := bson.ObjectIdHex(id)
	var res TripResponse
	if err := uc.session.DB("mingluliumongodb").C("Trips").FindId(oid).One(&res); err != nil {
        w.WriteHeader(404)
        return
    }
	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

func getSession() *mgo.Session {
	s, err := mgo.Dial("mongodb://minglu:liu273@ds057862.mongolab.com:57862/mingluliumongodb")
	if err != nil {
		panic(err)
	}
	s.SetMode(mgo.Monotonic, true)
	return s
}

func main() {
	r := httprouter.New()
	// router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", r))
}
