package hello

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"google.golang.org/appengine/aetest"

	"github.com/jarcoal/httpmock"
)

var adress = "https://api.github.com/repos/git/git/contributors"

func TestGetContributors(t *testing.T) {

	repo1 := Repo{}

	repo1.Project = "git/git"
	repo1.Name = "git"
	repo1.Owner = Owner{Login: "git"}
	repo1.Contributors = adress
	repo1.Languages = "https://api.github.com/repos/git/git/languages"

	/*
		Project      string `json:"full_name"`
		Name         string `json:"name"`
		Owner        Owner  `json:"owner"`
		Contributors string `json:"contributors_url"`
		Languages    string `json:"languages_url"`
	*/
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	raw, err := ioutil.ReadFile("./contributors.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	httpmock.RegisterResponder("GET", adress,
		httpmock.NewStringResponder(200, string(raw)))

	req, e := aetest.NewInstance(nil)
	if e != nil {
		t.Fatal("Can't create instance")
	}
	defer req.Close()
	r, e := req.NewRequest("GET", adress, nil)
	if e != nil {
		t.Fatal("Can't create request")
	}

	recorder := httptest.NewRecorder()

	cons := getContributors(recorder, r, repo1)

	mockCons := Contributors{}
	mockCons.Users = append(mockCons.Users, User{"gitster", 18497})
	mockCons.Users = append(mockCons.Users, User{"peff", 2702})
	mockCons.Users = append(mockCons.Users, User{"spearce", 1401})
	mockCons.Users = append(mockCons.Users, User{"torvalds", 1116})
	mockCons.Users = append(mockCons.Users, User{"dscho", 1104})
	mockCons.Users = append(mockCons.Users, User{"pclouds", 1100})
	mockCons.Users = append(mockCons.Users, User{"mhagger", 913})
	mockCons.Users = append(mockCons.Users, User{"jrn", 775})
	mockCons.Users = append(mockCons.Users, User{"rscharfe", 748})
	mockCons.Users = append(mockCons.Users, User{"jnareb", 512})
	mockCons.Users = append(mockCons.Users, User{"chriscool", 498})
	mockCons.Users = append(mockCons.Users, User{"j6t", 448})
	mockCons.Users = append(mockCons.Users, User{"felipec", 415})
	mockCons.Users = append(mockCons.Users, User{"avar", 400})
	mockCons.Users = append(mockCons.Users, User{"stefanbeller", 365})
	mockCons.Users = append(mockCons.Users, User{"npitre", 346})
	mockCons.Users = append(mockCons.Users, User{"paulusmack", 342})
	mockCons.Users = append(mockCons.Users, User{"trast", 336})
	mockCons.Users = append(mockCons.Users, User{"jiangxin", 302})
	mockCons.Users = append(mockCons.Users, User{"drafnel", 293})
	mockCons.Users = append(mockCons.Users, User{"mjg", 289})
	mockCons.Users = append(mockCons.Users, User{"moy", 282})
	mockCons.Users = append(mockCons.Users, User{"tronical", 252})
	mockCons.Users = append(mockCons.Users, User{"bk2204", 232})
	mockCons.Users = append(mockCons.Users, User{"szeder", 225})
	mockCons.Users = append(mockCons.Users, User{"pasky", 222})
	mockCons.Users = append(mockCons.Users, User{"bmwill", 217})
	mockCons.Users = append(mockCons.Users, User{"raalkml", 213})
	mockCons.Users = append(mockCons.Users, User{"artagnon", 196})
	mockCons.Users = append(mockCons.Users, User{"devzero2000", 193})

	if !reflect.DeepEqual(mockCons, cons) {
		t.Fail()
	}

}
