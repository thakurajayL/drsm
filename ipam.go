package drsm

import (
	"context"
	"fmt"
	"log"
	//"strconv"
	ipam "github.com/metal-stack/go-ipam"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO : should have ability to create new instances of ipam
func (d *Drsm) initIpam(opt *Options) {
	url := d.db.Url
	//hosts := []string{url}
	dbOptions := &options.ClientOptions{}
	dbOptions = dbOptions.ApplyURI(url)
	dbConfig := ipam.MongoConfig{DatabaseName: d.db.Name, CollectionName: "ipaddress", MongoClientOptions: dbOptions}
	mo, err := ipam.NewMongo(context.TODO(), dbConfig)
	if err != nil {
		log.Println("ipmodule error. NewMongo error  ", err)
	}
	ipModule := ipam.NewWithStorage(mo)
	log.Println("ipmodule ", ipModule)

	// TODO : handle more than one pool
	for _, v := range opt.IpPool {
		prefix, err := ipModule.NewPrefix(context.TODO(), v)
		if err != nil {
			panic(err)
		}
		d.ipModule = ipModule
		d.prefix = prefix
		break
	}
}

func (d *Drsm) acquireIp(name string) (string, error) {
	ip, err := d.ipModule.AcquireIP(context.TODO(), d.prefix.Cidr)
	if err != nil {
		err := fmt.Errorf("No address")
		return "", err
	}
	log.Println("Acquired IP ", ip.IP)
	return ip.IP.String(), nil
}
