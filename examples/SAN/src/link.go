package src

import (
	lib "gosan"
)

type SANLink struct {
	*Mock
	*NamingProps
	*lib.Link     `json:"-"`
	*SANComponent `json:"-"`
	iob           *IOBalancer `json:"-"`
	srcDst        []interface{} `json:"-"`
}

//func NewSANLink(link *lib.Link, iob *IOBalancer) *SANLink {
//	naming := &NamingProps{
//		ID:   link.Name,
//		Name: link.Name,
//		Type: "link",
//	}
//
//	src, ok := iob.allComponents[link.Src.Id]
//	if !ok {
//		log.Panicf("No such source with id: %s", link.Src.Id)
//	}
//	dst, ok := iob.allComponents[link.Dst.Id]
//	if !ok {
//		log.Panicf("No such destination with id: %s", link.Dst.Id)
//	}
//
//	srcDst := []interface{}{src, dst}
//
//	return &SANLink{
//		NamingProps: naming,
//		Link:        link,
//		SANComponent: &SANComponent{
//			currentState: "default",
//		},
//		iob:    iob,
//		srcDst: srcDst,
//	}
//}
