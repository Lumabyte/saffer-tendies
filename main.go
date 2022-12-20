package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TenderInfo struct {
	gorm.Model
	TenderID                  int                  `json:"id"`
	TenderNo                  string               `json:"tender_No"`
	Type                      string               `json:"type"`
	Delivery                  string               `json:"delivery"`
	Department                string               `json:"department"`
	CBrief                    bool                 `json:"cbrief"`
	CD                        string               `json:"cd"`                          // Should this be some sort of time format?
	DP                        string               `json:"dp"`                          // Should this be some sort of time format?
	DatePublished             string               `json:"date_Published"`              // Should this be some sort of time format?
	Brief                     string               `json:"brief"`                       // Should this be some sort of time format?
	ClosingDate               string               `json:"closing_Date"`                // Should this be some sort of time format?
	CompulsoryBriefingSession string               `json:"compulsory_briefing_session"` // Should this be some sort of time format?
	DepartmentID              string               `json:"departmentId"`
	ProvinceID                string               `json:"provinceId"`
	Status                    string               `json:"status"`
	Category                  string               `json:"category"`
	Description               string               `json:"description"`
	Province                  string               `json:"province"`
	ContactPerson             string               `json:"contactPerson"`
	Email                     string               `json:"email"`
	Telephone                 string               `json:"telephone"`
	Fax                       string               `json:"fax"`
	BriefingVenue             string               `json:"briefingVenue"`
	SupportingDocument        []SupportingDocument `json:"sd" gorm:"foreignKey:TendersID;references:TenderID"`
	Bidders                   string               `json:"bidders"`
	BiddersDoc                string               `json:"biddersdoc"`
	BiddersDocLink            string               `json:"biddersdoclink"`
	SuccessfulBidders         []SuccessfulBidders  `json:"successfulbidders" gorm:"foreignKey:TendersID;references:TenderID"`
	BF                        string               `json:"bf"`
	BC                        string               `json:"bc"`
	Reason                    string               `json:"reason"`
	ESubmissions              bool                 `json:"eSubmission"`
	Conditions                string               `json:"conditions"`
	Actions                   Actions              `json:"actions"  gorm:"foreignKey:TenderNo;references:TenderNo"`
}

func (TenderInfo) TableName() string {
	return "tenders"
}

type SupportingDocument struct {
	gorm.Model
	SupportDocumentID string `json:"supportDocumentID"`
	Filename          string `json:"fileName"`
	Extension         string `json:"extension"`
	TendersID         int    `json:"tendersID"`
	Active            bool   `json:"active"`
	UpdatedBy         string `json:"updatedBy"`
	DateModified      string `json:"dateModified"` // Should this be some sort of time format?
	Tenders           string `json:"tenders"`
}

func (SupportingDocument) TableName() string {
	return "supdocs"
}

type SuccessfulBidders struct {
	gorm.Model
	AwardID       int    `json:"awardID"`
	Company       string `json:"company"`
	ContactPerson string `json:"contactPerson"`
	ContactNumber string `json:"contactNumber"`
	TendersID     int    `json:"tendersID"`
	UpdatedBy     string `json:"updatedBy"`
	DateModified  string `json:"dateModified"`
	OCID          string `json:"ocid"`
	ReleaseID     int    `json:"releaseId"`
	SysStartTime  string `json:"sysStartTime"`
	SysEndTime    string `json:"sysEndTime"`
	Tenders       string `json:"tenders"`
}

func (SuccessfulBidders) TableName() string {
	return "successful_bidders"
}

type Actions struct {
	Authorized   bool   `json:"authorized"`
	TenderNo     string `json:"tender_No"`
	Notification bool   `json:"notification"`
	Bookmark     bool   `json:"bookmark"`
}

func (Actions) TableName() string {
	return "actions"
}

func LoadFromFile(filename string, data *[]TenderInfo) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln("Could not open JSON file ", filename)
	}
	defer f.Close()

	json_bytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln("Could not convert JSON into byte array")
	}

	json.Unmarshal(json_bytes, &data)
}

func LoadFromURL(json_url string, data *[]TenderInfo) {
	fmt.Println("Downloading from URL: ", json_url)
	resp, err := http.Get(json_url)
	if err != nil {
		log.Fatalln("Cannot retrieve JSON from URL", json_url)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatalln(err)
	}
}

func LoadTenderToDB(path string, db *gorm.DB) {
	var data []TenderInfo
	LoadFromURL(path, &data)
	result := db.CreateInBatches(&data, 100)
	if result.Error != nil {
		log.Fatalln("An error occurred during write to the DB")
	}
}

func main() {
	tender_urls := [4]string{
		"https://www.etenders.gov.za/Home/TenderOpportunities/?status=1",
		"https://www.etenders.gov.za/Home/TenderOpportunities/?status=2",
		"https://www.etenders.gov.za/Home/TenderOpportunities/?status=3",
		"https://www.etenders.gov.za/Home/TenderOpportunities/?status=4",
	}

	db, err := gorm.Open(sqlite.Open("new.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalln(err)
	}

	var (
		supdocs            SupportingDocument
		tenders            TenderInfo
		successful_bidders SuccessfulBidders
		actions            Actions
	)
	db.AutoMigrate(&supdocs)
	db.AutoMigrate(&tenders)
	db.AutoMigrate(&successful_bidders)
	db.AutoMigrate(&actions)

	fmt.Println("Loading tender data into the DB")

	// Load all tenders into the DB
	for index := range tender_urls {
		LoadTenderToDB(tender_urls[index], db)
		fmt.Printf("Loaded %s into the DB\n", tender_urls[index])
	}
}
