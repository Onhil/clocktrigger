package main

import (
	"log"
	"time"

	"github.com/Onhil/paragliding/db"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// MongoDB TODO
type MongoDB struct {
	DatabaseURL           string
	DatabaseName          string
	TrackCollectionName   string
	WebhookCollectionName string
}

func main() {
	var webhooks []paragliding.Webhooks
	var tracks []paragliding.Track
	GlobalDB := &MongoDB{
		DatabaseURL:           "mongodb://admin:admin1@ds145562.mlab.com:45562/paragliding",
		DatabaseName:          "paragliding",
		TrackCollectionName:   "Tracks",
		WebhookCollectionName: "Webhooks",
	}

	for {
		session, err := mgo.Dial(GlobalDB.DatabaseURL)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()

		wc := session.DB(GlobalDB.DatabaseName).C(GlobalDB.WebhookCollectionName)
		tc := session.DB(GlobalDB.DatabaseName).C(GlobalDB.TrackCollectionName)
		tc.Find(nil).All(&tracks)
		tcount, err := tc.Count()
		err = wc.Find(nil).All(&webhooks)
		if err != nil {
			log.Fatal(err)
		}
		for i := range webhooks {
			startTime := time.Now()
			if webhooks[i].AddedSince >= webhooks[i].MinTriggerValue {
				count, err := paragliding.SendMessage(webhooks[i], tracks, tcount, startTime)
				if err != nil {
					log.Fatal(err)
				} else {
					// Updates PrevTracksCount to current track collection count
					err = wc.Update(bson.M{"_id": webhooks[i].ID}, bson.M{"$set": bson.M{"prevtrackscount": count, "addedsince": 0}})
					if err != nil {
						log.Fatal(err)
					}
				}

			}

		}
		time.Sleep(10 * time.Minute)
	}
}
