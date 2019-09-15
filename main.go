package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"math"
	"strconv"
	"github.com/360EntSecGroup-Skylar/excelize"
)

type Response struct{
	Status int	`json:"status"`
	Message string `json:"message"`
	Total int	`json:"total"`
	Results []Results `json:"results"`
}

type Results struct {
	Name string `json:"name"`
	Address string `json:"address"`
	Province string `json:"province"`
	City string `json:"city"`
	Area string `json:"area"`
	Telephone string `json:"telephone"`
	DetailInfo DetailInfo `json:"detail_info"`
}

type DetailInfo struct {
	Image_num string `json:"image_num"`
	Detail_url string `json:"detail_url"`
}

func main()  {
	var total int
	allData := []Results{}
	params := map[string]string{
		"keyword": "电动车修理",	//搜索的关键字
		"location": "30.55269,104.075726",	//环形定位
		"radius": "25000",	//范围，单位（米）
		"ak": "XXXXXXXXXXXXXXXXXXXXXX",	//百度账号AK
	}
	url := "http://api.map.baidu.com/place/v2/search?query="+params["keyword"]+"&location="+params["location"]+"&radius="+params["radius"]+"&ak="+params["ak"]+"&output=json&scope=2&page_num=0&radius_limit=20&page_size=20"
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)

		rep := Response{}
		err := json.Unmarshal(r.Body, &rep)
		if err == nil {
			if total == 0 {
				total = rep.Total
			}
			for _, v := range rep.Results {
				allData = append(allData, v )
			}
		}
	})

	c.Visit(url)
	total = int(math.Ceil(float64(float64(total)/20.0) ) )
	for i := 1; i < total; i++ {
		c.Visit(fmt.Sprintf("%s?page_num=%d", url, i))
	}
	c.Wait()

	f := excelize.NewFile()
	// Set value of a cell.
	f.SetCellValue("Sheet1", "A1", "天府三街为中心25KM（电动车修理）")
	f.MergeCell("Sheet1", "A1", "H1")
	style, _ := f.NewStyle(`{"alignment":{"horizontal":"center"}}`)
	f.SetColWidth("Sheet1", "E", "E", 35)
	f.SetColWidth("Sheet1", "F", "F", 20)
	f.SetColWidth("Sheet1", "H", "H", 30)
	f.SetCellStyle("Sheet1", "A1", "H1", style)
	f.SetCellValue("Sheet1", "A2", "名称")
	f.SetCellValue("Sheet1", "B2", "省份")
	f.SetCellValue("Sheet1", "C2", "城市")
	f.SetCellValue("Sheet1", "D2", "区域")
	f.SetCellValue("Sheet1", "E2", "地址")
	f.SetCellValue("Sheet1", "F2", "联系电话")
	f.SetCellValue("Sheet1", "G2", "图片数量")
	f.SetCellValue("Sheet1", "H2", "详情")

	for k,v := range allData {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(k+3), v.Name)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(k+3), v.Province)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(k+3), v.City)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(k+3), v.Area)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(k+3), v.Address)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(k+3), v.Telephone)
		f.SetCellValue("Sheet1", "G"+strconv.Itoa(k+3), v.DetailInfo.Image_num)
		f.SetCellValue("Sheet1", "H"+strconv.Itoa(k+3), v.DetailInfo.Detail_url)
	}

	// Save xlsx file by the given path.
	err := f.SaveAs("./map.xlsx")
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(allData)
}
