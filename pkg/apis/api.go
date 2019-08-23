package apis

import (
	"encoding/json"
	"log"

	"github.com/emicklei/go-restful"
	"github.com/prometheus/alertmanager/template"

	"github.com/rancher/receiver/pkg/options"
	"github.com/rancher/receiver/pkg/tmpl"
)

func RegisterAPIs() {
	alertWs := new(restful.WebService).
		Path("/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	alertWs.Route(alertWs.POST("/{receiver-name}").To(sendAlert))
	restful.Add(alertWs)
}

func sendAlert(req *restful.Request, resp *restful.Response) {
	// TODO http req response action
	data := template.Data{}
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		log.Println("err: ", err)
		resp.WriteErrorString(500, err.Error())
		return
	}

	name := req.PathParameter("receiver-name")
	receiver, sender, err := options.GetReceiverAndSender(name)
	if err != nil {
		log.Printf("get receiver and sender err:%v", err)
		return
	}

	msg, err := tmpl.ExecuteTextString(data)
	if err != nil {
		log.Println("tmpl parse err: ", err)
		return
	}

	log.Println("msg:", msg)
	log.Println("receiver name:", receiver.Name)
	log.Println("provider name:", receiver.Provider)
	sender.Send(msg, receiver)

	return
}

