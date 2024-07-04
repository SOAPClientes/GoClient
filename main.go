package main

import (
	"encoding/xml"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
	"fmt"
	"strconv"

	"github.com/tiaguinho/gosoap"
)

type MensajeResponse struct {
	XMLName       xml.Name `xml:"MensajeResponse"`
	MensajeResult string   `xml:"MensajeResult"`
}

type InformacionResponse struct {
	XMLName          xml.Name `xml:"InformacionResponse"`
	InformacionResult string   `xml:"InformacionResult"`
}

type OperacionesResponse struct {
	XMLName            xml.Name `xml:"OperacionesResponse"`
	OperacionesResult string   `xml:"OperacionesResult"`
}

type TablaResponse struct {
	XMLName      xml.Name `xml:"TablaResponse"`
	TablaResult  []string `xml:"TablaResult>string"`
}

type EstudiantesResponse struct {
	XMLName             xml.Name `xml:"EstudiantesResponse"`
	EstudiantesResult   []string `xml:"EstudiantesResult>string"`
}

type Estudiantes2Response struct {
	XMLName           xml.Name `xml:"Estudiantes2Response"`
	Estudiantes2Result struct {
		String []string `xml:"string"`
	} `xml:"Estudiantes2Result"`
}


var templates = template.Must(template.ParseGlob("templates/*.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func mensajeHandler(w http.ResponseWriter, r *http.Request) {
	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}
	soap, err := gosoap.SoapClient("https://soapserviceserver.azurewebsites.net/WebServices/WebServer.asmx?WSDL", httpClient)
	if err != nil {
		log.Fatalf("SoapClient error: %s", err)
	}

	params := gosoap.Params{}
	res, err := soap.Call("Mensaje", params)
	if err != nil {
		log.Fatalf("Call error: %s", err)
	}

	var response MensajeResponse
	err = xml.Unmarshal(res.Body, &response)
	if err != nil {
		log.Fatalf("xml.Unmarshal error: %s", err)
	}

	templates.ExecuteTemplate(w, "mensaje.html", response.MensajeResult)
}

func informacionHandler(w http.ResponseWriter, r *http.Request) {
	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}
	soap, err := gosoap.SoapClient("https://soapserviceserver.azurewebsites.net/WebServices/WebServer.asmx?WSDL", httpClient)
	if err != nil {
		log.Fatalf("SoapClient error: %s", err)
	}

	params := gosoap.Params{}
	res, err := soap.Call("Informacion", params)
	if err != nil {
		log.Fatalf("Call error: %s", err)
	}

	var response InformacionResponse
	err = xml.Unmarshal(res.Body, &response)
	if err != nil {
		log.Fatalf("xml.Unmarshal error: %s", err)
	}

	templates.ExecuteTemplate(w, "informacion.html", response.InformacionResult)
}

func operacionesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		httpClient := &http.Client{
			Timeout: 15 * time.Second,
		}
		soap, err := gosoap.SoapClient("https://soapserviceserver.azurewebsites.net/WebServices/WebServer.asmx?WSDL", httpClient)
		if err != nil {
			log.Fatalf("SoapClient error: %s", err)
		}

		operacion := r.FormValue("operacion")
		valor1 := r.FormValue("n1")
		valor2 := r.FormValue("n2")

		params := gosoap.Params{
			"operacion": operacion,
			"valor1":    valor1,
			"valor2":    valor2,
		}
		res, err := soap.Call("Operaciones", params)
		if err != nil {
			log.Fatalf("Call error: %s", err)
		}

		var response OperacionesResponse
		err = xml.Unmarshal(res.Body, &response)
		if err != nil {
			log.Fatalf("xml.Unmarshal error: %s", err)
		}

		templates.ExecuteTemplate(w, "operaciones.html", response.OperacionesResult)
		return
	}
	templates.ExecuteTemplate(w, "operaciones.html", nil)
}

func tablaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		httpClient := &http.Client{
			Timeout: 15 * time.Second,
		}
		soap, err := gosoap.SoapClient("https://soapserviceserver.azurewebsites.net/WebServices/WebServer.asmx?WSDL", httpClient)
		if err != nil {
			log.Fatalf("SoapClient error: %s", err)
		}

		numero := r.FormValue("n1")
		params := gosoap.Params{
			"numero": numero,
		}
		res, err := soap.Call("Tabla", params)
		if err != nil {
			log.Fatalf("Call error: %s", err)
		}

		var response TablaResponse
		err = xml.Unmarshal(res.Body, &response)
		if err != nil {
			log.Fatalf("xml.Unmarshal error: %s", err)
		}

		templates.ExecuteTemplate(w, "tabla.html", response.TablaResult)
		return
	}
	templates.ExecuteTemplate(w, "tabla.html", nil)
}

func estudiantesHandler(w http.ResponseWriter, r *http.Request) {
	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}
	soap, err := gosoap.SoapClient("https://soapserviceserver.azurewebsites.net/WebServices/WebServer.asmx?WSDL", httpClient)
	if err != nil {
		log.Fatalf("SoapClient error: %s", err)
	}

	params := gosoap.Params{}
	res, err := soap.Call("Estudiantes", params)
	if err != nil {
		log.Fatalf("Call error: %s", err)
	}

	var response EstudiantesResponse
	err = xml.Unmarshal(res.Body, &response)
	if err != nil {
		log.Fatalf("xml.Unmarshal error: %s", err)
	}

	students := make([][]string, len(response.EstudiantesResult))
	for i, student := range response.EstudiantesResult {
		studentData := strings.Split(student, ",")
		if len(studentData) < 3 {

			for len(studentData) < 3 {
				studentData = append(studentData, "")
			}
		}

		studentWithNumber := append([]string{fmt.Sprintf("%d", i+1)}, studentData...)
		students[i] = studentWithNumber
	}

	for i, student := range students {
		nota := student[3]
		estado := "Desaprobado"
		if notaFloat, err := strconv.ParseFloat(nota, 64); err == nil {
			if notaFloat >= 10.5 {
				estado = "Aprobado"
			}
		}
		students[i] = append(student, estado)
	}

	templates.ExecuteTemplate(w, "estudiantes.html", students)
}

func estudiantes2Handler(w http.ResponseWriter, r *http.Request) {
	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}
	soap, err := gosoap.SoapClient("https://soapserviceserver.azurewebsites.net/WebServices/WebServer.asmx?WSDL", httpClient)
	if err != nil {
		log.Fatalf("SoapClient error: %s", err)
	}

	params := gosoap.Params{}
	res, err := soap.Call("Estudiantes2", params)
	if err != nil {
		log.Fatalf("Call error: %s", err)
	}

	var response Estudiantes2Response
	err = xml.Unmarshal(res.Body, &response)
	if err != nil {
		log.Fatalf("xml.Unmarshal error: %s", err)
	}

	log.Printf("Raw response: %+v\n", res.Body)
	log.Printf("Deserialized response: %+v\n", response)

	students := make([][]string, len(response.Estudiantes2Result.String))
	for i, student := range response.Estudiantes2Result.String {
		studentData := strings.Split(student, ",")
		if len(studentData) < 3 {
			for len(studentData) < 3 {
				studentData = append(studentData, "")
			}
		}
		studentWithNumber := append([]string{fmt.Sprintf("%d", i+1)}, studentData...)
		students[i] = studentWithNumber
	}

	for i, student := range students {
		nota := student[3]
		estado := "Desaprobado"
		if notaFloat, err := strconv.ParseFloat(nota, 64); err == nil {
			if notaFloat >= 10.5 {
				estado = "Aprobado"
			}
		}
		students[i] = append(student, estado)
	}

	log.Printf("Processed students: %+v\n", students)

	templates.ExecuteTemplate(w, "estudiantes2.html", students)
}



func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/mensaje", mensajeHandler)
	http.HandleFunc("/informacion", informacionHandler)
	http.HandleFunc("/operaciones", operacionesHandler)
	http.HandleFunc("/tabla", tablaHandler)
	http.HandleFunc("/estudiantes", estudiantesHandler)
	http.HandleFunc("/estudiantes2", estudiantes2Handler)

	log.Println("Server started at :8084")
	log.Fatal(http.ListenAndServe(":8084", nil))
}
