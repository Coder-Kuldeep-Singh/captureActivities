package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// generate random data for line chart
func generateBarItems(length *[]DailyGraph) []opts.LineData {
	items := make([]opts.LineData, 0)
	for _, item := range *length {
		items = append(items, opts.LineData{Value: item.Count, Name: item.Product})
	}
	return items
}

// GetDays renders the results day by day as counts
func GetDays(c *gin.Context) {
	dbType := "mysql"
	db := loadEnv().connect(dbType)
	length := getUsedProductPerDays(db)
	db.Close()
	titles := []string{}
	for _, title := range *length {
		titles = append(titles, title.Days)
	}
	bar := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "My first bar chart generated by go-echarts",
		Subtitle: "It's extremely easy to use, right?",
	}))
	bar.Tooltip.Show = true
	// Put data into instance
	bar.SetXAxis(titles).
		AddSeries("Days", generateBarItems(length))
	// Where the magic happens

	bar.Render(c.Writer)
}

// Daily renders the all clicks on today and the applications where user(I) clicked
func Daily(c *gin.Context) {
	dbType := "mysql"
	db := loadEnv().connect(dbType)
	length := getUsedProductPerDay(db)
	db.Close()
	titles := []string{}
	for _, title := range *length {
		titles = append(titles, title.Product)
	}
	// log.Println(titles)
	bar := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "My first bar chart generated by go-echarts",
		Subtitle: "It's extremely easy to use, right?",
	}))
	bar.Tooltip.Show = true
	// bar.Tooltip.Formatter = fmt.Sprintf("daily based code\n product: count:")
	// Put data into instance
	bar.SetXAxis(titles).
		AddSeries("Today's clicks", generateBarItems(length))

	bar.Render(c.Writer)
}
