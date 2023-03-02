package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Result struct {
	Company     string
	UsageVT     float64
	UsageNT     float64
	MainBreaker float64
	Phases      float64
	CostVT      string
	CostNT      string
	CostMonth   string
	Poze        string
	Total       string
	TotalMonth  string
}

func input(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/layout.html", "templates/content_input.html")
	t.ExecuteTemplate(w, "layout", "")
}

func cost(w http.ResponseWriter, r *http.Request) {
	company := r.FormValue("company")
	usageVT := convStrFloat(r.FormValue("usageVT"))
	priceVT := convStrFloat(r.FormValue("priceVT"))
	calVT := usageVT * priceVT
	valVT := fmt.Sprintf("%.2f", calVT)

	usageNT := convStrFloat(r.FormValue("usageNT"))
	priceNT := convStrFloat(r.FormValue("priceNT"))
	calNT := usageNT * priceNT
	valNT := fmt.Sprintf("%.2f", calNT)

	constPay := convStrFloat(r.FormValue("constPay"))
	priceInputBreaker := convStrFloat(r.FormValue("priceInputBreaker"))
	oTE := convStrFloat(r.FormValue("OTE"))
	calMonthly := 12 * (constPay + priceInputBreaker + oTE)
	valMonthly := fmt.Sprintf("%.2f", calMonthly)

	byConsumption := convStrFloat(r.FormValue("byConsumption"))
	byBreaker := convStrFloat(r.FormValue("byBreaker"))
	mainBreaker := convStrFloat(r.FormValue("mainBreaker"))
	phases := convStrFloat(r.FormValue("Phases"))
	aPoze := byConsumption * (usageNT + usageVT)
	bPoze := byBreaker * mainBreaker * phases
	calPoze := calPoze(aPoze, bPoze)
	valPoze := fmt.Sprintf("%.2f", calPoze)

	calTotal := calVT + calNT + calMonthly + calPoze
	valTotal := fmt.Sprintf("%.2f", calTotal)
	calTotalMonth := calTotal / 12
	valTotalMonth := fmt.Sprintf("%.2f", calTotalMonth)

	t, _ := template.ParseFiles("templates/layout.html", "templates/content_output.html")
	//data1 := []string{valVT, valNT, costMonthly, aPoze, bPoze}
	result := Result{Company: company, UsageVT: usageVT, UsageNT: usageNT, MainBreaker: mainBreaker, Phases: phases, CostVT: valVT, CostNT: valNT, CostMonth: valMonthly, Poze: valPoze, Total: valTotal, TotalMonth: valTotalMonth}

	t.ExecuteTemplate(w, "layout", result)
}

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	//mux.ServeFiles("/static/*filepath", http.Dir("static"))
	mux.HandleFunc("/input", input)
	mux.HandleFunc("/cost", cost)

	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}

func convStrFloat(s string) (f float64) {
	if s == "" {
		return
	} else {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			fmt.Println(err)
		}
		return f
	}
}

func calPoze(aPoze float64, bPoze float64) (poze float64) {
	if aPoze >= bPoze {
		poze = bPoze
	} else {
		poze = aPoze
	}
	return poze
}
