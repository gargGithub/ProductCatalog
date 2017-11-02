package main

import (
	"html/template"
	"net/http"
	"github.com/jinzhu/gorm"
	_"github.com/go-sql-driver/mysql"
	"log"
	"fmt"
	"github.com/gorilla/mux"
)

var tpl *template.Template
var db *gorm.DB
var err error


func main() {
    Connect()
	defer db.Close()
	router:=mux.NewRouter()
	router.HandleFunc("/",index)
	router.HandleFunc("/types",types)
	router.HandleFunc("/{type}",listCompanies)
	router.HandleFunc("/{type}/{company}",listCars)
	router.HandleFunc("/{type}/{company}/{carname}",listVariant)
	router.HandleFunc("/{type}/{company}/{carname}/{variant}",variantDetails)
	router.Handle("/favicon.ico",http.NotFoundHandler())
	http.ListenAndServe(":8080",router)
}



func init(){
 tpl = template.Must(template.ParseGlob("Templates/*.html"))
}



func index(w http.ResponseWriter, req *http.Request){

	tpl.ExecuteTemplate(w,"index.html",nil)
}


func types(w http.ResponseWriter, req *http.Request){
	tpl.ExecuteTemplate(w,"types.html",nil)
}


func listCompanies(w http.ResponseWriter, req *http.Request){

	var route = mux.Vars(req)
	if route["type"] == "hatch"{
		var h []HatchBack
		db.Debug().Table("hatch_backs").Select("DISTINCT company_name").Find(&h)
		var hNames = []string{}
		for _,v:=range h {
			hNames = append(hNames, v.CompanyName)
		}
		tpl.ExecuteTemplate(w, "CompaniesListHatch.html", hNames)
	}
	if route["type"] == "sedan"{
		var s []Sedan
		db.Debug().Table("sedans").Select("DISTINCT company_name").Find(&s)
		var sNames = []string{}
		for _,v:=range s{
			sNames = append(sNames,v.CompanyName)
		}

		tpl.ExecuteTemplate(w, "CompaniesListSedan.html", sNames)
	}
}


func listCars(w http.ResponseWriter, req *http.Request){

	var c = mux.Vars(req)
	route:="/"+c["type"]+"/"+c["company"]
	var UrlParts []string
	UrlParts = append(UrlParts,c["type"],c["company"])
	if c["type"]=="hatch" {
		var hCars []HatchBack

		db.Debug().Table("hatch_backs").Select("DISTINCT car_name").Where("company_name=?", c["company"]).Find(&hCars)
		for i:=0;i<len(hCars);i++{
			hCars[i].UrlString = route
			hCars[i].StripUrl = UrlParts
		}

		tpl.ExecuteTemplate(w, "ListCarsComp.html", hCars)
	}

	if c["type"]=="sedan" {
		var sCars []Sedan

		db.Debug().Table("sedans").Select("DISTINCT car_name").Where("company_name=?", c["company"]).Find(&sCars)
		for i:=0;i<len(sCars);i++{
			sCars[i].UrlString = route
			sCars[i].StripUrl = UrlParts

		}

		tpl.ExecuteTemplate(w, "ListCarsComp.html", sCars)
	}

}

func listVariant(w http.ResponseWriter, req *http.Request)  {
	c:=mux.Vars(req)
	route:="/"+c["type"]+"/"+c["company"]+"/"+c["carname"]

	if c["type"]=="hatch"{
		var hVariant []HatchBack

		db.Debug().Table("hatch_backs").Select("DISTINCT variant_name").Where("car_name=? AND company_name=?",c["carname"],c["company"]).Find(&hVariant)
		for i:=0;i<len(hVariant);i++{
			hVariant[i].UrlString = route
		}

		tpl.ExecuteTemplate(w,"VariantList.html",hVariant)

	}

	if c["type"]=="sedan"{
		var sVariant []Sedan

		db.Debug().Table("sedans").Select("DISTINCT variant_name").Where("car_name=? AND company_name=?",c["carname"],c["company"]).Find(&sVariant)
		for i:=0;i<len(sVariant);i++{
			sVariant[i].UrlString = route
		}

		tpl.ExecuteTemplate(w,"VariantList.html",sVariant)

	}

}

func variantDetails(w http.ResponseWriter, req *http.Request) {
	c:=mux.Vars(req)

	if c["type"]=="hatch"{
      var hJoin HatchJoin
		db.Debug().Raw("SELECT specs.*, hatch_backs.variant_name FROM specs INNER JOIN hatch_backs ON hatch_backs.s_id = specs.s_id Where hatch_backs.variant_name=? AND hatch_backs.car_name=?",c["variant"],c["carname"]).Scan(&hJoin)
		db.Debug().Raw("SELECT features.*, hatch_backs.variant_name FROM features INNER JOIN hatch_backs ON hatch_backs.f_id = features.f_id Where hatch_backs.variant_name=? AND hatch_backs.car_name=?",c["variant"],c["carname"]).Scan(&hJoin)
		fmt.Println(hJoin)
		tpl.ExecuteTemplate(w,"variantSpec.html",hJoin)

	}

	if c["type"]=="sedan"{
		var sJoin SedanJoin
		db.Raw("SELECT specs.*, sedans.variant_name FROM specs INNER JOIN sedans ON sedans.s_id = specs.s_id Where sedans.variant_name=? AND sedans.car_name=?", c["variant"],c["carname"]).Scan(&sJoin)
		db.Raw("SELECT features.*, sedans.variant_name FROM features INNER JOIN sedans ON sedans.f_id = features.f_id Where sedans.variant_name=? AND sedans.car_name=?", c["variant"],c["carname"]).Scan(&sJoin)
		tpl.ExecuteTemplate(w,"variantSpec.html",sJoin)

	}

}

func Connect() {

	db,err = gorm.Open("mysql","root:password@tcp(127.0.0.1:3306)/product_catalog?charset=utf8" +
		"&parseTime=True&loc=Local")
	if err!=nil {
		log.Fatal(err)
	} else{
	  fmt.Println("connected")
	}


	err = db.DB().Ping()
	if err!=nil{
		log.Fatal(err)
	}

}


type HatchBack struct {
	HatchID uint `gorm:"primary_key"`
	CompanyName string `gorm:"type:varchar(100)"`
	CarName string `gorm:"type:varchar(100)"`
	VariantName string `gorm:"type:varchar(100)"`
	SID uint
	FID uint
	UrlString string `gorm:"-"`
	StripUrl []string `gorm:"-"`
}


type Sedan struct {
	SedanID uint `gorm:"primary_key"`
	CompanyName string `gorm:"type:varchar(100)"`
	CarName string `gorm:"type:varchar(100)"`
	VariantName string `gorm:"type:varchar(100)"`
	SID uint
	FID uint
	UrlString string `gorm:"-"`
	StripUrl []string `gorm:"-"`
}


type Specs struct {
	SID uint `gorm:"primary_key"`
	Length string `gorm:"type:varchar(50)"`
	Width string `gorm:"type:varchar(50)"`
	Height string `gorm:"type:varchar(50)"`
	Bootspace string `gorm:"type:varchar(50)"`
	FuelTankCapacity string `gorm:"type:varchar(50)"`
	Mileage string `gorm:"type:varchar(50)"`
	Cylinders string `gorm:"type:varchar(50)"`
}


type Features struct {
	FID uint `gorm:"primary_key"`
	Airbags string `gorm:"type:varchar(50)"`
	ABS string `gorm:"type:varchar(50)"`
	FourWheelDrive string `gorm:"type:varchar(50)"`
	AirConditioner string `gorm:"type:varchar(50)"`
	CupHolders string `gorm:"type:varchar(50)"`
	PowerWindows string `gorm:"type:varchar(50)"`
	Tachometer string `gorm:"type:varchar(50)"`
}

type HatchJoin struct {
	HatchBack
	Specs
	Features
}
type SedanJoin struct {
	Sedan
	Specs
	Features
}