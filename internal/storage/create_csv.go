package storage

import (
	"encoding/csv"
	"log"
	"os"
	"strings"

	"github.com/mystpen/parser-test/internal/model"
)

func CreateCSV(InfluencersInfo *[]model.InfluencerInfo) error {
	file, err := os.Create("influencers.csv")
	if err != nil {
		log.Println("Error creating CSV file:", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Rank", "Account", "Name", "Avatar Image", "Category", "Subscribers", "Country", "Eng. (Auth.)", "Eng. (Avg.)"}
	if err := writer.Write(header); err != nil {
		log.Println("Error writing header:", err)
		return err
	}

	// Write data rows
	for _, info := range *InfluencersInfo {
		record := []string{
			info.Rank,
			info.Account,
			info.Name,
			info.AvatarImage,
			strings.Join(info.Category, ", "),
			info.Subscribers,
			info.Country,
			info.EngAuth,
			info.EngAvg,
		}
		if err := writer.Write(record); err != nil {
			log.Println("Error writing record:", err)
			return err
		}
	}
	return nil
}
