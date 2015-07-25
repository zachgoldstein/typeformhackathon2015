package main
import (
	"fmt"
	"net/http"
	"log"
	"html"
	"io/ioutil"
	"github.com/nlopes/slack"
	"time"
	"strings"
	"bytes"
	"strconv"
	"encoding/json"
	"github.com/melvinmt/firebase"
)


func main(){

	api := slack.New("xoxb-8194552449-GU4fHYqiPlcVEC9gLcNgwv2a")
	api.SetDebug(false)

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		log.Println("hit webhook")

		requestBody, err := ioutil.ReadAll(r.Body);

		if err != nil {
			log.Println("err ",err)
		}
		log.Print("request body ",string(requestBody))
		log.Print("request url ",r.URL)


		var results FormResults

		err = json.Unmarshal(requestBody, &results)
		if (err != nil) {
			log.Println("COULD NOT UNMARSHALL RESULTS ",err)
		}

		log.Println("UNMARSHALLED RESULSTs ", results)

		userId := r.URL.Query().Get("id")
		timestamp := r.URL.Query().Get("timestamp")
		chanId := r.URL.Query().Get("chanId")

		log.Println("timestamp ",timestamp, " chanId ",chanId)

		user, err := api.GetUserInfo(userId)
		if (err != nil) {
			fmt.Println("Could not find person with id ",userId )
		}

		storeVotedPersonInFirebase(user)

		storeResultInFirebase(results)

		dmPerson(user, "thank you for voting!! See the results at http://3d8d52d9.ngrok.io/", api) //ADD link to results

//		itemRef := slack.NewRefToMessage(chanId, timestamp)
//		err = api.AddReaction("+1", itemRef)
//		if (err != nil) {
//			fmt.Println("Could not post reaction, ", err)
//		}

		msg := fmt.Sprintf("%v just voted! See the results at http://3d8d52d9.ngrok.io/",user.Name)

		api.PostMessage("C085QFNM7", msg, slack.PostMessageParameters {
			AsUser: true,
		})


		//store results in db DONE
		//save user image to indicate they've voted DONE
		//send a dm to thank the person and show results page DONE
		//add reaction to message DONE

		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	go setupRTM(api)

	log.Println("Listening on 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupRTM (api *slack.Slack) {
	chSender := make(chan slack.OutgoingMessage)
	chReceiver := make(chan slack.SlackEvent)

	wsAPI, err := api.StartRTM("", "http://example.com")
	if err != nil {
		fmt.Errorf("%s\n", err)
	}

	go wsAPI.HandleIncomingEvents(chReceiver)
	go wsAPI.Keepalive(20 * time.Second)
	go func(wsAPI *slack.SlackWS, chSender chan slack.OutgoingMessage) {
		for {
			select {
			case msg := <-chSender:
				wsAPI.SendMessage(&msg)
			}
		}
	}(wsAPI, chSender)
	for {
		select {
		case msg := <-chReceiver:
			fmt.Print("Event Received: ")
			switch msg.Data.(type) {
				case slack.HelloEvent:
				// Ignore hello
				case *slack.MessageEvent:
				a := msg.Data.(*slack.MessageEvent)
				fmt.Printf("Message: %v\n", a)
				fmt.Printf("ChanID:",a.ChannelId)
				if (strings.Contains(a.Text, "U085QG8D7")) {
					createFormForEverybody( a, api )
				}
				case *slack.PresenceChangeEvent:
				a := msg.Data.(*slack.PresenceChangeEvent)
				fmt.Printf("Presence Change: %v\n", a)
				case slack.LatencyReport:
				a := msg.Data.(slack.LatencyReport)
				fmt.Printf("Current latency: %v\n", a.Value)
				case *slack.SlackWSError:
				error := msg.Data.(*slack.SlackWSError)
				fmt.Printf("Error: %d - %s\n", error.Code, error.Msg)
				default:
				fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}

func dmPerson (user *slack.User, message string,  api *slack.Slack) {
	if (user.Name == "typebot") {
		return
	}

	_,_,chanId, err := api.OpenIMChannel(user.Id)
	if (err != nil) {
		fmt.Println("IM err ",err)
	}
	api.PostMessage(chanId, message, slack.PostMessageParameters {
		AsUser: true,
	})
}

func dmAll(message string, api *slack.Slack) {

	users, err := api.GetUsers()
	if (err != nil) {
		fmt.Println("Channel err ",err)
	}

	for _, user := range users {
		if (user.Name == "typebot") {
			continue
		}

		_,_,chanId, err := api.OpenIMChannel(user.Id)
		if (err != nil) {
			fmt.Println("IM err ",err)
		}
		api.PostMessage(chanId, message, slack.PostMessageParameters {
			AsUser: true,
		})
	}
}

func createFormForEverybody(formMessage *slack.MessageEvent, api *slack.Slack) {
	users, err := api.GetUsers()
	if (err != nil) {
		fmt.Println("Channel err ",err)
	}

	for _, user := range users {
		//create user form DONE
		//send dm to fill out form DONE

		formResp := createFormForUser(formMessage, user.Id)
		if (user.Name == "typebot" || user.Id == "U085QG8D7") {
			continue
		}
		_, _, chanId, err := api.OpenIMChannel(user.Id)
		if (err != nil) {
			fmt.Println("IM err ",err)
		}

		fmt.Println("formResp ",formResp)
		formDM := fmt.Sprintf("A form has been created, go fill it out! %v", formResp.Links[1].Href)

		api.PostMessage(chanId, formDM, slack.PostMessageParameters {
			AsUser: true,
		})
	}
}

func createFormForUser(message *slack.MessageEvent, userId string) FormSubmissionResponse {
	msgComponents := strings.Split(message.Text, "|")
	log.Println("msgComponents ",msgComponents)

	if ( len(msgComponents) < 2 ) {
		log.Println("not enough data to create form")
		return FormSubmissionResponse{}
	}

	question := strings.Split(msgComponents[0], ":")[1]

	rawChoices := msgComponents[1:]

	choices := []Choice{}

	for _, choice := range rawChoices {
		choiceComponents := strings.Split(choice, ",")

		parseString := strings.Trim(choiceComponents[1], " ")
		imgId64, err := strconv.ParseInt(parseString, 10, 32)
		if (err != nil) {
			fmt.Println("err parsing ",err)
		}
		fmt.Println("imgId64",imgId64)
		imgId := int(imgId64)
		fmt.Println("imgId",imgId)
		choices = append(choices, *NewChoice(imgId, choiceComponents[0]))
	}

	url := "https://api.typeform.io/v0.3/forms"
	webhook := fmt.Sprintf("https://4c7e576f.ngrok.io/webhook?id=%v&timestamp=%v&chanId=%v",userId, message.Timestamp, message.ChannelId)


	formSubmission := NewFormSubmission(question, choices, webhook)

	// sample message => "Would you rather | Fight a horse sized duck?, 2542 | Fight a thousand duck sized horses?, 2543"

	rawSubmission, _ := json.Marshal(formSubmission)

	buf := bytes.NewBuffer(rawSubmission)

	req, err := http.NewRequest("POST", url, buf)
	if (err != nil) {
		fmt.Println("Could not build new request ",err)
	}

	req.Header.Add("x-api-token", "8effac447a26ff118ddcea4335db9a3e")

	res, err := http.DefaultClient.Do(req)
	if (err != nil) {
		fmt.Println("Could not execute request ",err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var resp FormSubmissionResponse

	err = json.Unmarshal(body, &resp)
	if (err != nil) {
		fmt.Println("Could not unmarshall form repsonse ",err)
	}

	return resp
}

func NewFormSubmission (question string, choices []Choice, webhook string) *FormSubmission {

	field := Field{
		Choices : choices,
		Question : question,
		Required : true,
		Type : "picture_choice",
	}

	return &FormSubmission{
		Fields : []Field{field},
		WebhookSubmitURL : webhook,
		Title : "Typeform hackathon question",
	}


}

type FormSubmission struct {
	Fields []Field `json:"fields"`
	Title            string `json:"title"`
	WebhookSubmitURL string `json:"webhook_submit_url"`
}

type Field struct {
	Choices []Choice `json:"choices"`
	Description string `json:"description"`
	Question    string `json:"question"`
	Required    bool   `json:"required"`
	Type        string `json:"type"`
}

type Choice struct {
	ImageID int    `json:"image_id"`
	Label   string `json:"label"`
}

func NewChoice(imageId int, label string) *Choice {
	return &Choice{
		ImageID: imageId,
		Label: label,
	}
}

type FormResults struct {
	Answers []struct {
		Data struct {
			Type  string `json:"type"`
			Value struct {
				Label string      `json:"label"`
				Other interface{} `json:"other"`
			} `json:"value"`
		} `json:"data"`
		FieldID float64 `json:"field_id"`
	} `json:"answers"`
	ID    string `json:"id"`
	Token string `json:"token"`
}

type FormSubmissionResponse struct {
	Fields []struct {
		AllowMultipleSelections bool `json:"allow_multiple_selections"`
		Choices                 []struct {
			Filename string `json:"filename"`
			Height   int    `json:"height"`
			ImageID  int    `json:"image_id"`
			Label    string `json:"label"`
			Width    int    `json:"width"`
		} `json:"choices"`
		Description string `json:"description"`
		ID          int    `json:"id"`
		Labels      bool   `json:"labels"`
		Question    string `json:"question"`
		Required    bool   `json:"required"`
		Type        string `json:"type"`
	} `json:"fields"`
	ID    string `json:"id"`
	Links []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links"`
	Title            string `json:"title"`
	WebhookSubmitURL string `json:"webhook_submit_url"`
}


func storeVotedPersonInFirebase(user *slack.User) {
	var err error

	url := "https://typeformhackathon.firebaseio.com/votedUsers"

	// Can also be your Firebase secret:
	authToken := "67LoKprWLZvah2E6HxscXsMDh10y9SbuAxvi1wCd"

	// Auth is optional:
	ref := firebase.NewReference(url).Auth(authToken)

	// Write the value to Firebase.
	if err = ref.Push(user); err != nil {
		panic(err)
	}

}

func storeResultInFirebase(results FormResults) {

	var err error

	url := "https://typeformhackathon.firebaseio.com/answers"

	// Can also be your Firebase secret:
	authToken := "67LoKprWLZvah2E6HxscXsMDh10y9SbuAxvi1wCd"

	// Auth is optional:
	ref := firebase.NewReference(url).Auth(authToken)

	// Write the value to Firebase.
	if err = ref.Push(results); err != nil {
		panic(err)
	}

	log.Println("WROTE TO FIREBASE")

//	// Now, we're going to retrieve the person.
//	personUrl := "https://SampleChat.firebaseIO-demo.com/users/fred"
//
//	personRef := firebase.NewReference(personUrl).Export(false)
//
//	fred := Person{}
//
//	if err = personRef.Value(fred); err != nil {
//		panic(err)
//	}
//
//	fmt.Println(fred.Name.First, fred.Name.Last) // prints: Fred Swanson

}


/*

{
   "title":"Typeform IO hackathon",
   "webhook_submit_url":"http://22c95451.ngrok.io/webhook",
   "fields":[
      {
         "type":"picture_choice",
         "question":"What would you rather fight?",
         "description":"Life or death decision",
         "required":true,
         "choices":[
            {
               "image_id":2542,
               "label":"A horse sized duck"
            },
            {FormSubmission
               "image_id":2543,
               "label":"A thousand duck sized horses"
            }
         ]
      }
   ]
}

*/


/*

{
   "answers":[
      {
         "data":{
            "type":"choice",
            "value":{
               "label":" Fight a thousand duck sized horses?",
               "other":null
            }
         },
         "field_id":1.272649e+06
      }
   ],
   "id":"S65RosIui4nvig",
   "token":"8ed184766f516a5885e95f526182415b"
}

*/