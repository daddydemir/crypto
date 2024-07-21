package handler

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/coincap"
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/daddydemir/crypto/pkg/graphs/ma"
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
								<th> Name </th>
								<th> Symbol </th>
								<th data-sort="price"> Price </th>
								<th> RSI </th>
								<th> SMA </th>
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
	<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.5.1/jquery.min.js">
	<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz" crossorigin="anonymous"></script>
	<script src="https://cdn.datatables.net/1.10.22/js/jquery.dataTables.min.js" />
	<script src="https://cdn.datatables.net/1.10.22/js/dataTables.bootstrap4.min.js" />
	<script>
		$(document).onload(function() {
			$('#coin-table').DataTable();
		});
	</script>
</html>
`

	coins := coincap.ListCoins()

	content := `
<tr>
	<td> %s </td>
	<td> %s </td>
	<td> %.3f </td>
	<td> <a href="%s" target="_blank"> graph </a> </td>
	<td> <a href="%s" target="_blank"> graph </a> </td>
</tr>
`
	var contents string
	for _, coin := range coins {
		temp := content
		temp = fmt.Sprintf(temp, coin.Name, coin.Symbol, coin.PriceUsd, "/api/v1/graph/rsi/"+coin.Id, "/api/v1/graph/sma/"+coin.Id)
		contents += temp
	}

	html = fmt.Sprintf(html, contents)

	w.Write([]byte(html))
}

func rsiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)

	coin := vars["coin"]
	fmt.Printf("Coin : %v \n", coin)

	rsi := graphs.RSI{}
	histories := rsi.Calculate(coin)
	draw := rsi.Draw(histories)
	draw(w, r)
}

func smaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)

	coin := vars["coin"]
	fmt.Printf("Coin : %v \n", coin)

	draw := ma.Draw(coin)
	draw(w, r)
}
