package controller

import (
	"fmt"
	"log"
	"net/http"
	"webpaygo/api/models"
	"webpaygo/api/utils"

	"github.com/fenriz07/Golang-Transbank-WebPay-Rest/pkg/transaction"
)

// Definición de la estructura TransactionLog
type TransactionLog struct {
	NumberOrder string
	IdSession   string
	Response    string
	Error       string
}

func VerifTransaction(w http.ResponseWriter, r *http.Request) {
	/*
	 En caso del que pago sea anulado, comprobar si existe el parametro TBK_TOKEN.
	  Si existe el pago fue anulado por el usuario y debe comprobarse su estado con el Commit,
	  Si fue anulado adicionalmente tenemos los parametros TBK_ORDEN_COMPRA || TBK_ID_SESION
	*/

	log.Println("******************empieza*************")

	var token string = ""
	var numberOrder string = ""
	var idSession string = ""

	canceledToken := r.FormValue("TBK_TOKEN")

	if len(canceledToken) != 0 {
		token = canceledToken
		numberOrder = r.FormValue("TBK_ORDEN_COMPRA")
		idSession = r.FormValue("TBK_ID_SESION")

		log.Printf("Number Order: %s\n Id Session: %s\n", numberOrder, idSession)

	} else {
		token = r.FormValue("token_ws")
	}

	/*Commit de la transacción y resultado de la misma*/
	resp, err := transaction.Commit(token)

	if err != nil {
		fmt.Println(err)
	}

	// Crea un nuevo LogEntry
	newLogEntry := models.LogEntry{
		NumberOrder:       numberOrder,
		IdSession:         idSession,
		Status:            resp.Status,
		Amount:            resp.Amount,
		BuyOrder:          resp.BuyOrder,
		SessionID:         resp.SessionID,
		AccountingDate:    resp.AccountingDate,
		TransactionDate:   resp.TransactionDate,
		PaymentTypeCode:   resp.PaymentTypeCode,
		CardDetail:        resp.CardDetail,
		AuthorizationCode: resp.AuthorizationCode,
	}
	logTransactionData(newLogEntry)

	log.Println(resp)

	/*Obtención del status de la transacción*/
	resp2, err := transaction.GetStatus(token)

	log.Println(resp2)

	if err != nil {
		log.Println(err)
	}

	/*Anulación*/
	resp3, err := transaction.Refund(token, 1000)

	if err != nil {
		log.Println(err)
	}

	log.Println("Respuesta 3")
	log.Println(resp3)

}

func logTransactionData(logData models.LogEntry) {
	// Almacenar el log en la base de datos
	err := storeLogEntry(logData)
	if err != nil {
		log.Println("Error storing log in the database:", err)
	}
}

func storeLogEntry(logData models.LogEntry) error {
	db, err := utils.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	//log.Printf("Inserting log entry into database: %+v\n", logData)

	_, err = db.Exec(`
        INSERT INTO dato_factura (
            status, amount, buy_order, session_id, 
            accounting_date, transaction_date, payment_type_code, 
            card_number, authorization_code, number_order, id_session
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `,
		logData.Status, logData.Amount, logData.BuyOrder, logData.SessionID,
		logData.AccountingDate, logData.TransactionDate,
		logData.PaymentTypeCode, logData.CardDetail.CardNumber,
		logData.AuthorizationCode, logData.NumberOrder, logData.IdSession,
	)
	if err != nil {
		log.Println("Error database:", err)
		return err

	}

	return nil
}
