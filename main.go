package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// const advertised_url string = "https://www.etenders.gov.za/Home/TenderOpportunities/?status=1"
// const awarded_url string = "https://www.etenders.gov.za/Home/TenderOpportunities/?status=2"
const closed_url string = "https://www.etenders.gov.za/Home/TenderOpportunities/?status=3"

// const cancelled_url string = "https://www.etenders.gov.za/Home/TenderOpportunities/?status=4"

// Using this actually nested everything one object deeper, eg [{{}}]
// type ClosedTender struct {
// 	Tender
// 	SupportingDocument []SupportingDocument
// }

type Tender struct {
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
	BF                        string               `json:"bf"`
	BC                        string               `json:"bc"`
	Conditions                string               `json:"conditions"`
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

// var resp = []byte(`
// 	[{
// 		"id": 39043,
// 		"tender_No": "MN 112-2021",
// 		"type": "Request for Bid(Open-Tender)",
// 		"delivery": "NO 2 INDUSTRIA CRESCENT - KWADUKUZA - KWADUKUZA /STANGER - 4450",
// 		"department": "Kwadukuza Municipality",
// 		"cbrief": true,
// 		"cd": "Friday, 04 March 2022 - 12:00",
// 		"dp": "Friday, 04 November 2022",
// 		"date_Published": "2022-11-04T00:00:00",
// 		"brief": "Friday, 25 February 2022 - 13:00",
// 		"closing_Date": "2022-03-04T12:00:00",
// 		"compulsory_briefing_session": "2022-02-25T13:00:00",
// 		"status": "Closed",
// 		"category": "Services: Electrical",
// 		"description": "TENDER MN 112-2021 APPOINTMENT OF CONTRACTOR FOR MV SUBSTATIONS UPGRADE AND REFURBISHMENT A PERIOD OF THREE YRS",
// 		"province": "KwaZulu-Natal",
// 		"contactPerson": "DHANESH RAMPERSADH",
// 		"email": "dhaneshr@kwadukuza.gov.za",
// 		"telephone": "032-437-5115",
// 		"fax": "032-437-5087",
// 		"briefingVenue": "PMU BOARD ROOM NO 2 INDUSTRIA CRES LAVOIPIERRE BUILDING YARD COMPLEX {BACK ENTRY}IN/OUT CAR PARK",
// 		"sd": [
// 			{
// 				"supportDocumentID": "e0ac0207-1887-42bc-920b-0a7839d4e8ea",
// 				"fileName": "TENDER MN 112-2021 APPOINTMENT OF CONTRACTOR FOR MV SUBSTATION UPGRADE \u0026 REFURBISHMENT 3 YEARS.docx",
// 				"extension": ".docx",
// 				"tendersID": 39043,
// 				"active": true,
// 				"updatedBy": "vanessaps@kwadukuza.gov.za",
// 				"dateModified": "2022-11-04T08:24:46.0939341",
// 				"tenders": null
// 			}
// 		],
// 		"bf": " Yes",
// 		"bc": "Yes",
// 		"conditions": "SEE ATTACHED ALREADY ADVERTISED THIS TENDER RECREATED AS IT WAS CANCELLED IN ERROR"
// 	}]`,
// )

func (Tender) TableName() string {
	return "tenders"
}

func (SupportingDocument) TableName() string {
	return "supdocs"
}

func main() {

	// This code was used with resp above
	// var data []Tender
	// err := json.Unmarshal(resp, &data)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Printf("%#+v\n", data[0].SupportingDocument[0].SupportDocumentID)
	// post, _ := json.Marshal(data)
	// fmt.Println(string(post))

	resp, err := http.Get(closed_url)
	if err != nil {
		log.Fatalln("Got an error on that!")
	}
	defer resp.Body.Close()

	var data []Tender

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatalln(err)
	}

	// Commenting this out for when I wanna use the _FULL_ JSON data again
	// json_data, err := json.Marshal(data)
	// if err != nil {
	// 	log.Fatalln("Failed marshalling data to JSON")
	// }
	// // Just printing the value here so the compiler doesn't bitch
	// // To do: Write data to file
	// fmt.Printf("%T\n", string(json_data))

	db, err := gorm.Open(sqlite.Open("new.db"), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalln(err)
	}

	var supdocs SupportingDocument
	var tenders Tender
	db.AutoMigrate(&supdocs)
	db.AutoMigrate(&tenders)
	result := db.CreateInBatches(&data, 100)
	// db.First(&data)
	fmt.Println(result.Error)

	// json_data, err := json.Marshal(data[0])
	// if err != nil {
	// 	log.Fatalln("Failed marshalling data to JSON")
	// }
	// fmt.Printf("%+v\n", string(json_data))
}
