package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Aksh-Bansal-dev/bingelist/pkg/db"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "will-be-random-later"
)

func setup() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/login")
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func googleCallbackHandler(w http.ResponseWriter, r *http.Request, database *gorm.DB) {
	log.Println("/login/redirect")
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	var data map[string]string
	json.Unmarshal(content, &data)
	log.Println("NewUser:", data)

	token, err := db.AddUser(database, data["email"])
	if err != nil {
		http.Error(w, "Something went wrong in DB", 400)
		log.Println(err)
		return
	}

	fmt.Fprintf(w, `
		<h1>Loading...</h1>	
		<script>
			localStorage.setItem("token", "%s");
			window.location.href = "/";
		</script>
	`, token)
}

func getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}
