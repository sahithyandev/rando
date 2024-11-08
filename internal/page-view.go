package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
)

const RESPONSE_MAX_AGE = 108600

var propertyIdMap = map[string]string{
	"sahithyan.dev":    "356543625",
	"kalvi.lk":         "398313900",
	"s1.sahithyan.dev": "449008877",
}

func GetPageView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	targetDomain := query.Get("domain")
	if targetDomain == "" {
		originHeader := r.Header.Get("Origin")
		fmt.Println("origin", originHeader)

		if originHeader == "" {
			w.Header().Set("status", "400")
			io.WriteString(w, "Couldn't find target domain")
			return
		}

		targetDomain = strings.Split(strings.Replace(originHeader, "https://", "", 1), ":")[0]
	}
	pagePath := query.Get("path")

	targetPropertyId, isDomainPresent := propertyIdMap[targetDomain]
	if !isDomainPresent {
		w.Header().Set("status", "400")
		io.WriteString(w, fmt.Sprintf("Unsupported domain: %s", targetDomain))
		return
	}

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./blog-pageviews-1713811580440-51ea0fe058d1.json")
	ctx := context.Background()
	analyticsdataService, err := analyticsdata.NewService(ctx)

	if err != nil {
		fmt.Printf("Couldn't start analyticsdata service: %s\n", err)
		w.Header().Set("status", "500")
		io.WriteString(w, "Internal error occured")
		return
	}

	runReportRequest := &analyticsdata.RunReportRequest{
		Dimensions: []*analyticsdata.Dimension{
			{
				Name: "pagePath",
			},
		},
		DateRanges: []*analyticsdata.DateRange{
			{
				StartDate: "2022-03-05",
				EndDate:   "today",
			},
		},
		Metrics: []*analyticsdata.Metric{
			{
				Name: "screenPageViews",
			},
		},
		OrderBys: []*analyticsdata.OrderBy{
			{
				Desc: true,
				Metric: &analyticsdata.MetricOrderBy{
					MetricName: "screenPageViews",
				},
			},
		},
	}
	if pagePath != "" {
		runReportRequest.DimensionFilter = &analyticsdata.FilterExpression{
			Filter: &analyticsdata.Filter{
				FieldName: "pagePath",
				StringFilter: &analyticsdata.StringFilter{
					MatchType: "EXACT",
					Value:     pagePath,
				},
			},
		}

	}

	result, err := analyticsdataService.Properties.RunReport(fmt.Sprintf("properties/%s", targetPropertyId), runReportRequest).Do()

	if err != nil {
		fmt.Printf("Couldn't run report on analyticsdata: %s\n", err)
		w.Header().Set("status", "500")
		io.WriteString(w, "Internal error occured")
		return
	}

	pageViewMap := map[string]string{}

	for i := 0; i < int(result.RowCount); i++ {
		row := result.Rows[i]
		key := row.DimensionValues[0].Value
		value := row.MetricValues[0].Value

		if pagePath == key {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", RESPONSE_MAX_AGE))
			io.WriteString(w, fmt.Sprintf(`{"views": "%s", "dateRetrievedOn": "%s" }`, value, time.Now().UTC().Format(time.RFC3339)))
			return
		}

		pageViewMap[key] = value
	}
	output, err := json.Marshal(pageViewMap)
	if err != nil {
		fmt.Printf("Error occured while marshalling json: %s\n", err)
		w.Header().Set("status", "500")
		io.WriteString(w, "Internal error occured")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", RESPONSE_MAX_AGE))
	io.WriteString(w, string(output))
}
