package apis

import (
	"encoding/json"
	"github.com/prometheus/common/model"
	"io/ioutil"

	"github.com/emicklei/go-restful"
	"github.com/prometheus/alertmanager/template"

	"github.com/rancher/webhook-receiver/pkg/options"
	"github.com/rancher/webhook-receiver/pkg/tmpl"
	log "github.com/sirupsen/logrus"
)

func RegisterAPIs() {
	alertWs := new(restful.WebService).
		Path("/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	alertWs.Route(alertWs.POST("/{receiver-name}").To(sendAlert))
	alertWs.Route(alertWs.GET("/healthz").To(reportLiveness))
	alertWs.Route(alertWs.GET("/state").To(reportState))
	restful.Add(alertWs)
}

func sendAlert(req *restful.Request, resp *restful.Response) {
	bodyData, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		log.Errorf("read req body err:%v", err)
		resp.WriteErrorString(400, err.Error())
		return
	}

	td := template.Data{}
	if err := json.Unmarshal(bodyData, &td); err != nil {
		// rancher test send
		rtd := model.Alerts{}
		if err := json.Unmarshal(bodyData, &rtd); err == nil {
			log.Info("test success")
			resp.WriteHeader(200)
			return
		}

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

	msg, err := tmpl.ExecuteTextString(td)
	if err != nil {
		log.Errorf("tmpl parse err: %v", err)
		resp.WriteErrorString(500, err.Error())
		return
	}

	log.Infof("receiver:%s,provider:%s,msg:%s\n", name, receiver.Provider, msg)
	if err := sender.Send(msg, receiver); err != nil {
		log.Errorf("send msg err:%v", err)
		resp.WriteErrorString(500, err.Error())
		return
	} else {
		resp.WriteHeader(200)
		log.Infof("send msg successful")
	}

	return
}

func reportLiveness(req *restful.Request, resp *restful.Response) {
	resp.WriteHeader(200)
	return
}

func reportState(req *restful.Request, resp *restful.Response) {
	if options.GetState() {
		resp.WriteHeader(200)
		return
	} else {
		resp.WriteHeader(500)
		return
	}
}
