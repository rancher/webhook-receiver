package options

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"

	"github.com/rancher/receiver/pkg/providers"
	"github.com/rancher/receiver/pkg/providers/alibaba"
	"github.com/rancher/receiver/pkg/providers/netease"
)

var (
	mut       sync.RWMutex
	receivers map[string]providers.Receiver
	senders   map[string]providers.Sender
)

// when occur error, it will panic directly
func Init(configPath string) {
	updateMemoryConfig(configPath)
	go syncMemoryConfig(configPath)
}

func GetReceiverAndSender(receiverName string) (providers.Receiver, providers.Sender, error) {
	mut.RLock()
	defer mut.RUnlock()

	receiver, exists := receivers[receiverName]
	if !exists {
		return providers.Receiver{}, nil, fmt.Errorf("error, receiver:%s is not exists\n", receiverName)
	}

	sender, exists := senders[receiver.Provider]
	if !exists {
		return providers.Receiver{}, nil, fmt.Errorf("error, provider:%s is not exists\n", receiver.Provider)
	}

	return receiver, sender, nil
}

func syncMemoryConfig(configPath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("new watcher err: ", err)
	}
	if err := watcher.Add(configPath); err != nil {
		log.Fatalf("watch file:%s err:%v", configPath, err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				log.Println("event:", event)
				log.Println("ok:", ok)
				if !ok {
					break
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					updateMemoryConfig(configPath)
				}

			case err, ok := <-watcher.Errors:
				log.Println(err)
				log.Println(ok)
			}
		}
		log.Println("break select")
	}()
}

func updateMemoryConfig(configPath string) {
	log.Println("info: ", configPath, " is update")
	opt, err := newOption(configPath)
	if err != nil {
		log.Fatal("new option err: ", err)
	}

	mut.Lock()
	defer mut.Unlock()

	receivers = make(map[string]providers.Receiver)
	for _, v := range opt.Receivers {
		receivers[v.Name] = v
	}

	senders = make(map[string]providers.Sender)
	for k, v := range opt.Providers {
		creator, err := getProviderCreator(k)
		if err != nil {
			log.Fatal(err)
		}
		sender, err := creator(v)
		if err != nil {
			log.Fatal(err)
		}
		senders[k] = sender
	}
}

func getProviderCreator(name string) (providers.Creator, error) {
	switch name {
	case alibaba.Name:
		return alibaba.New, nil
	case netease.Name:
		return netease.New, nil
	default:
		return nil, errors.New(fmt.Sprintf("provider :%s is not support", name))
	}
}

// for yaml parse
type option struct {
	Providers map[string]map[string]string
	Receivers []providers.Receiver
}

func newOption(configPath string) (option, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return option{}, err
	}
	opt := option{}
	if err := yaml.Unmarshal(data, &opt); err != nil {
		return option{}, err
	}

	return opt, nil
}
