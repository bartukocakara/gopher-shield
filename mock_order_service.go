package main
import ("fmt"; "net/http")
func main() {
    http.HandleFunc("/checkout", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Order Confirmed: Payment Processed Successfully")
    })
    http.ListenAndServe(":3000", nil)
}