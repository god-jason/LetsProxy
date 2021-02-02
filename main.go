package main

import (
	"flag"
	"fmt"
	"github.com/kardianos/service"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var (
	help       bool
	install    bool
	uninstall  bool
	configPath string
)

var serviceConfig = &service.Config{
	Name:        "LetsProxy",
	DisplayName: "LetsProxy",
	Description: "HTTPS反向代理 内嵌Letsencrypt",
	Arguments:   nil,
}

func init() {
	flag.BoolVar(&help, "h", false, "帮助")
	flag.StringVar(&configPath, "c", os.Args[0]+".yaml", "配置文件") //os.Args[0]+".yaml"
	flag.BoolVar(&install, "i", false, "安装服务")
	flag.BoolVar(&uninstall, "u", false, "卸载服务")

	rand.Seed(time.Now().UnixNano())
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	var err error
	configPath, err = filepath.Abs(configPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("配置文件", configPath)

	//加载配置
	err = LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	//传递参数到服务
	serviceConfig.Arguments = []string{"-c", configPath}

	// 构建服务对象
	prog := &Program{}
	s, err := service.New(prog, serviceConfig)
	if err != nil {
		log.Fatal(err)
	}

	// 用于记录系统日志
	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	if install {
		err = s.Install()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("安装成功")
		return
	}

	if uninstall {
		err = s.Uninstall()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("卸载成功")
		return
	}

	err = s.Run()
	if err != nil {
		_ = logger.Error(err)
	}

}

type Program struct{}

func (p *Program) Start(s service.Service) error {
	log.Println("开始服务")
	go p.run()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	log.Println("停止服务")
	return nil
}

func (p *Program) run() {

	// 此处编写具体的服务代码
	hup := make(chan os.Signal, 2)
	signal.Notify(hup, syscall.SIGHUP)
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt, os.Kill)

	go func() {
		for {
			select {
			case <-hup:
			case <-quit:
				os.Exit(0)
			}
		}
	}()

	Serve()
}

func Serve() {
	var ds = make([]string, 0)
	var dm = make(map[string][]*url.URL)

	//解析域名
	for domain, target := range config.Proxies {
		log.Println("反向代理", domain, "-->", target)
		targets := strings.Split(target, ",")

		urls := make([]*url.URL, 0)
		for _, t := range targets {
			u, e := url.Parse(t)
			if e != nil {
				log.Fatal(u)
				return
			}
			urls = append(urls, u)
		}

		domains := strings.Split(domain, ",")
		for _, d := range domains {
			d = strings.TrimSpace(d)
			ds = append(ds, d)
			dm[d] = urls
		}
	}

	//初始化autocert
	manager := &autocert.Manager{
		Cache:      autocert.DirCache(config.Cache),
		Email:      config.Email,
		HostPolicy: autocert.HostWhitelist(ds...),
		Prompt:     autocert.AcceptTOS,
	}

	handler := &httputil.ReverseProxy{Director: func(req *http.Request) {
		us, ok := dm[req.Host]
		if !ok {
			return
		}

		var u *url.URL
		if len(us) > 1 {
			u = us[rand.Int()%len(us)]
		} else {
			u = us[0]
		}

		//log.Println("Request", req.URL.String(), u.String())
		req.URL.Scheme = u.Scheme
		req.URL.Host = u.Host
		req.URL.Path = urlJoin(u.Path, req.URL.Path) //拼接路径

		//添加参数
		if u.RawQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = u.RawQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = u.RawQuery + "&" + req.URL.RawQuery
		}

		//设置User-Agent(
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}}
	//创建server
	svr := &http.Server{
		Addr:      "0.0.0.0:443",
		TLSConfig: manager.TLSConfig(),
		Handler:   handler,
	}

	//监听http
	if config.Http {
		go func() {
			server := &http.Server{
				Addr:    "0.0.0.0:80",
				Handler: &redirectHandler{},
			}
			log.Fatal(server.ListenAndServe())
		}()
	}

	//监听https
	log.Fatal(svr.ListenAndServeTLS("", ""))
}

func urlJoin(a, b string) string {
	if !strings.HasSuffix(a, "/") {
		a += "/"
	}
	if strings.HasPrefix(b, "/") {
		a += b[1:]
	} else {
		a += b
	}
	return a
}
