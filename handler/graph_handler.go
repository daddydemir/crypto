package handler

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/daddydemir/crypto/pkg/graphs/bollingerBands"
	"github.com/daddydemir/crypto/pkg/graphs/ma"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"github.com/daddydemir/crypto/pkg/service"
	"github.com/gorilla/mux"
	"net/http"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	html := `
<html>
	<head>
		<title> Coins </title>
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
	</head>
	<body>
		<h3> Coins </h3>
		<div class="container">
			<div class="row justify-content-center">
				<div class="col-md-auto">
					<table class="table table-hover table-striped table-bordered table-resposive" id="coin-table">
						<thead>
							<tr>
								<th> # </th>
								<th> Name </th>
								<th> Symbol </th>
								<th data-sort="price"> Price </th>
								<th> RSI </th>
								<th> Index </th>
								<th> SMA </th>
								<th> EMA </th>
								<th> MA </th>
								<th> Bollinger Bands </th>
							</tr>
						</thead>
						<tbody>
							%s
						</tbody>
					</table>
			</div>
		</div>
	</div>
	</body>
	<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.5.1/jquery.min.js"> </script>
	<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz" crossorigin="anonymous"></script>
	<script src="https://cdn.datatables.net/1.10.22/js/jquery.dataTables.min.js" > </script>
	<script src="https://cdn.datatables.net/1.10.22/js/dataTables.bootstrap4.min.js" > </script>
	<script>
		$('#coin-table').DataTable({
			paging: false,
			searching: false,
		});
	</script>
</html>
`

	coins := coincap.ListCoins()
	rsi := graphs.RSI{}

	content := `
<tr class='%s'>
	<td> %d </td>
	<td> %s </td>
	<td> %s </td>
	<td> %.3f </td>
	<td> <a href="%s" target="_blank"> graph </a> </td>
	<td> %.2f </td>
	<td> <a href="%s" target="_blank"> graph </a> </td>
	<td> <a href="%s" target="_blank"> graph </a> </td>
	<td> <a href="%s" target="_blank"> graph </a> </td>
	<td> <a href="%s" target="_blank"> graph </a> </td>
</tr>
`
	var contents string
	for i, coin := range coins {
		temp := content
		var index float32
		var class string
		if i <= 20 {
			index = rsi.Index(coin.Id)
			if index >= 70 {
				class = "table-success"
			} else if index <= 30 {
				class = "table-danger"
			} else {
				class = "table-warning"
			}
		}

		temp = fmt.Sprintf(temp, class, i, coin.Name, coin.Symbol, coin.PriceUsd,
			"/api/v1/graph/rsi/"+coin.Id,
			index,
			"/api/v1/graph/sma/"+coin.Id,
			"/api/v1/graph/ema/"+coin.Id,
			"/api/v1/graph/ma/"+coin.Id,
			"/api/v1/graph/bollingerBands/"+coin.Id,
		)
		contents += temp
	}

	html = fmt.Sprintf(html, contents)

	w.Write([]byte(html))
}

func rsiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)
	coin := vars["coin"]

	graphicService := service.NewRsiService(coin)
	function := graphicService.Draw()
	function(w, r)
}

func smaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)
	coin := vars["coin"]

	graphicService := service.NewSmaService(coin, 10)
	function := graphicService.Draw()
	function(w, r)
}

func emaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)
	coin := vars["coin"]

	graphicService := service.NewEmaService(coin, 25)
	function := graphicService.Draw()
	function(w, r)
}

func maHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)
	coin := vars["coin"]
	newMa := ma.NewMa(coin, 0)

	draw := newMa.Draw(newMa.Calculate())
	draw(w, r)
}

func bollingerBandsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)
	coin := vars["coin"]

	bands := bollingerBands.NewBollingerBands(coin, 20)
	list := bands.Calculate()

	function := bands.Draw(list)

	function(w, r)

}
