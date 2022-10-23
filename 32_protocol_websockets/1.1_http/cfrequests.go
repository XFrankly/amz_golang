package main

import (
	"bytes"
	"net/http"
	"time"
)

/*
import CloudFlareClient
import sys

token = 'p123456'

try:
    domain = sys.argv[1]
except IndexError:
        exit('第一个参数需要输入域名')
try:
    IP_address = sys.argv[2]
except IndexError:
        exit('第二个参数需要输入服务器IP')


def main(domain,IP_address):
    zone_name = domain
    cf = CloudFlare.CloudFlare(token=token)
    try:
        zone_info = cf.zones.post(data={'jump_start': False, 'name': zone_name})
    except CloudFlare.exceptions.CloudFlareAPIError as e:
        exit('/zones.post %s - %d %s' % (zone_name, e, e))
    except Exception as e:
        exit('/zones.post %s - %s' % (zone_name, e))
    zone_id = zone_info['id']

    dns_record = {'type':'A', 'name': domain, 'content': IP_address, 'proxied': True}
    try:
        r = cf.zones.dns_records.post(zone_id, data=dns_record)
    except CloudFlare.exceptions.CloudFlareAPIError as e:
        exit('/zones.dns_records.post %s %s - %d %s' % (zone_name, dns_record['name'], e, e))
    print(r)

def set_always_https(domain):
    zone_name = domain
    cf = CloudFlare.CloudFlare(token=token)
    params = {'name': domain, 'per_page': 1, 'page': 1}
    try:
        zones = cf.zones.get(params=params)
    except CloudFlare.exceptions.CloudFlareAPIError as e:
        exit('/zones.get %d %s - api call failed' % (e, e))
    except Exception as e:
        exit('/zones - %s - api call failed' % (e))
    zone_id = zones[0]['id']

    try:
        r = cf.zones.settings.always_use_https.patch(zone_id, data={'value': 'on'})
    except CloudFlare.exceptions.CloudFlareAPIError as e:
        exit('/zones.settings.always_use_https.patch %d %s - api call failed' % (e, e))
    updated_value = r['value']
    print('always_use_https status now is %s' %updated_value)

def get_domain_list():
    cf = CloudFlare.CloudFlare(token=token)

    zones = cf.zones.get(params={'per_page':1,'page':1})
    print(zones)
    for zone in zones:
        print(zone['id'], zone['name'])


if __name__ == '__main__':
    main(domain, IP_address)
    set_always_https(domain)


*/
func (self *Tasks) Do_Request(method string, full_path string, body []byte) *http.Response {
	bytes_body := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, full_path, bytes_body) //bytes.NewBuffer(body))

	logg.Println(method, "<-- request ->> :", "request body:", bytes_body)
	if err != nil {
		//Handle Error
		logger.Fatal(method, "<-- URL ->> :", full_path, "err:", err, method, full_path)
	}

	req.Header.Set("X-Custom-Header", "FuckTheWorld.")
	req.Header.Set("Content-Type", "application/json")
	req.Header = *self.ClientConn.Header

	self.ClientConn.Client.Timeout = time.Second * 250 //超时时间 25s
	res, err := self.ClientConn.Client.Do(req)
	if err != nil {
		//Handle Error
		logger.Fatal(method, "<-- URL ->> :", full_path, "\n err:", err, method)
	}
	// res 属于 *http.Response 结构如下
	//Status     string // e.g. "200 OK"
	//	StatusCode int    // e.g. 200
	//	Proto      string // e.g. "HTTP/1.0"
	//	ProtoMajor int    // e.g. 1
	//	ProtoMinor int    // e.g. 0
	//Header Header
	//Body io.ReadCloser
	//ContentLength int64
	//TransferEncoding []string
	//Close bool
	//Uncompressed bool
	//Trailer Header
	//Request *Request
	//TLS *tls.ConnectionState

	return res

}
