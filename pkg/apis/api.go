package apis

import (
	"encoding/json"

	"github.com/emicklei/go-restful"
	"github.com/prometheus/alertmanager/template"

	"github.com/rancher/receiver/pkg/options"
	"github.com/rancher/receiver/pkg/tmpl"
	log "github.com/sirupsen/logrus"
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
		log.Errorf("webhook data parse err:%v", err)
		resp.WriteErrorString(400, err.Error())
		return
	}

	name := req.PathParameter("receiver-name")
	receiver, sender, err := options.GetReceiverAndSender(name)
	if err != nil {
		log.Errorf("get receiver name:%s err:%v", name, err)
		resp.WriteErrorString(500, err.Error())
		return
	}

	msg, err := tmpl.ExecuteTextString(data)
	if err != nil {
		log.Errorf("tmpl parse err: %v", err)
		resp.WriteErrorString(500, err.Error())
		return
	}

	log.Infof("receiver:%s,provider:%s,msg:%s\n", receiver.Name, receiver.Provider, msg)
	if err := sender.Send(msg, receiver); err != nil {
		log.Errorf("send msg err:%v", err)
		resp.WriteErrorString(500, err.Error())
		return
	}

	log.Infof("send msg successful")

	return
}
