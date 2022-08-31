package grpc

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
	"rusprofile/internal/api"
	"strings"
)

type Server struct {
	api.UnimplementedRusprofileServer
}

func NewServer() *grpc.Server {
	gsrv := grpc.NewServer()
	srv := Server{}
	api.RegisterRusprofileServer(gsrv, &srv)
	return gsrv
}

func Run(serv *grpc.Server, address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server is listening at %s", address)
	if err := serv.Serve(lis); err != nil {
		log.Fatalf("gRPC server failed to start: %v", err)
	}
}

func (*Server) GetCompanyInfo(ctx context.Context, input *api.SearchingInfo) (*api.ReturnedCompanyInfo, error) {
	doc, err := getDoc(input.Inn)
	if err != nil {
		return nil, err
	}
	info := api.ReturnedCompanyInfo{
		Name:  doc.Find(".company-name").Text(),
		Inn:   doc.Find("#clip_inn").Text(),
		Kpp:   doc.Find("#clip_kpp").Text(),
		Chief: parsePerson(doc.Find("a[class='link-arrow gtm_main_fl']").Text()),
	}
	if len(info.Name) == 0 || len(info.Inn) == 0 || len(info.Kpp) == 0 || len(info.Chief.Surname) == 0 {
		return nil, status.Error(codes.NotFound, "Company was not found")
	}
	return &info, nil
}

func parsePerson(fullName string) *api.ReturnedCompanyInfo_Person {
	fields := strings.Fields(fullName)
	switch len(fields) {
	case 0:
		return &api.ReturnedCompanyInfo_Person{}
	case 1:
		return &api.ReturnedCompanyInfo_Person{Surname: fields[0]}
	case 2:
		return &api.ReturnedCompanyInfo_Person{
			Surname:   fields[0],
			FirstName: fields[1],
		}
	case 3:
		return &api.ReturnedCompanyInfo_Person{
			Surname:    fields[0],
			FirstName:  fields[1],
			MiddleName: fields[2],
		}
	default:
		return &api.ReturnedCompanyInfo_Person{
			Surname:    fields[0],
			FirstName:  fields[1],
			MiddleName: strings.Join(fields[2:], " "),
		}
	}
}

func getDoc(INN string) (*goquery.Document, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.rusprofile.ru/search?query=%s", INN))
	if err != nil {
		return nil, status.Error(codes.NotFound, "INN was not found")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, status.Error(codes.Unknown, "Response status code of external service is not 200")
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return doc, nil
}
